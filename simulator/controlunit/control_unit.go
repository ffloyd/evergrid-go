package controlunit

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/scheduler"
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/comm"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
)

type schedulerGenerator func(logContext *logrus.Entry) scheduler.Scheduler

// ControlUnit -
type ControlUnit struct {
	name   string
	fsm    simenv.AgentFSM
	simenv *simenv.SimEnv
	log    *logrus.Entry

	sharedData *SharedData

	schedGen  schedulerGenerator
	scheduler scheduler.Scheduler

	localQueue localQueue
	monitor    monitor

	workerNames []string
	workers     map[string]simenv.Agent
	sendLock    sync.Mutex
}

// New -
func New(cfg networkcfg.AgentCfg, workerNames []string, sharedData *SharedData, schedGen schedulerGenerator, logContext *logrus.Entry) *ControlUnit {
	return &ControlUnit{
		name:       cfg.Name,
		log:        logContext,
		schedGen:   schedGen,
		sharedData: sharedData,

		workerNames: workerNames,
		workers:     make(map[string]simenv.Agent),
	}
}

// Name -
func (cu *ControlUnit) Name() string {
	return cu.name
}

// Run -
func (cu *ControlUnit) Run(env *simenv.SimEnv) simenv.AgentChans {
	cu.log = cu.log.WithFields(logrus.Fields{
		"agent": cu.Name(),
		"tick":  env.CurrentTick(),
	})

	cu.simenv = env
	cu.fsm = simenv.NewAgentFSM(cu.log)

	cu.scheduler = cu.schedGen(cu.log.WithField("context", "scheduler"))

	cu.localQueue = newLocalQueue(cu)
	cu.monitor = newMonitor(cu)

	// Leader election
	cu.sharedData.LeaderElection.Do(func() {
		cu.sharedData.LeaderControlUnit = cu
		cu.log.Info("Become leader")
	})

	for _, workerName := range cu.workerNames {
		cu.workers[workerName] = env.Find(workerName)
	}

	cu.monitor.Run()
	cu.scheduler.Run()
	go cu.work()
	return cu.fsm.Chans()
}

// Send - respond means that request arrived to proper scheduler
func (cu *ControlUnit) Send(msg interface{}) chan interface{} {
	schedChans := cu.scheduler.RequestChans()

	switch request := msg.(type) {
	case comm.ControlUnitUploadDataset:
		cu.sendLock.Lock()
		cu.fsm.ToWorking()
		schedChans.UploadDataset <- scheduler.ReqUploadDataset{
			Dataset: request.Dataset,
		}

		if <-schedChans.DelegateToLeader {
			<-cu.sharedData.LeaderControlUnit.Send(request)
		}
		cu.fsm.ToIdle()
		cu.sendLock.Unlock()
	case comm.ControlUnitRunExperiment:
		cu.sendLock.Lock()
		cu.fsm.ToWorking()
		schedChans.RunExperiment <- scheduler.ReqRunExperiment{
			Calculator: request.Calculator,
			Dataset:    request.Dataset,
		}

		if <-schedChans.DelegateToLeader {
			<-cu.sharedData.LeaderControlUnit.Send(request)
		}
		cu.fsm.ToIdle()
		cu.sendLock.Unlock()
	case scheduler.DoUploadDataset:
		cu.processWorkerAction(request)
	case scheduler.DoBuildCalculator:
		cu.processWorkerAction(request)
	case scheduler.DoRunCalculator:
		cu.processWorkerAction(request)
	default:
		cu.log.Panicf("Unknown request type: %v", request)
	}

	response := make(chan interface{})
	go func() {
		response <- simenv.Ok{}
	}()
	return response
}

func (cu *ControlUnit) work() {
	cu.fsm.SetStopFlag(true)
	chans := cu.scheduler.ControlChans()
	for {
		cu.sendLock.Lock()
		cu.fsm.ToReady()
		cu.fsm.ToWorking()
		cu.localQueue.Process()
		cu.fsm.ToIdle()
		cu.sendLock.Unlock()

		doneChan := cu.fsm.ToDoneChan()

	SelectLoop:
		for {
			select {
			case request := <-chans.UploadDataset:
				cu.processWorkerAction(request)
				chans.Done <- scheduler.Done{}
			case request := <-chans.BuildCalculator:
				cu.processWorkerAction(request)
				chans.Done <- scheduler.Done{}
			case request := <-chans.RunCalculator:
				cu.processWorkerAction(request)
				chans.Done <- scheduler.Done{}
			case <-doneChan:
				break SelectLoop
			}
		}
	}
}

func (cu *ControlUnit) processWorkerAction(request interface{}) {
	var workerName string
	var actionType string
	switch value := request.(type) {
	case scheduler.DoUploadDataset:
		workerName = value.Worker
		actionType = fmt.Sprintf("%T", value)
	case scheduler.DoRunCalculator:
		workerName = value.Worker
		actionType = fmt.Sprintf("%T", value)
	case scheduler.DoBuildCalculator:
		workerName = value.Worker
		actionType = fmt.Sprintf("%T", value)
	}

	cu.log.WithFields(logrus.Fields{
		"worker":      workerName,
		"action_type": actionType,
	}).Info("Process worker action")

	if cu.workers[workerName] != nil {
		cu.localQueue.Push(request)
	} else {
		cu.sharedData.Mutex.Lock()
		correctCUName := cu.sharedData.Workers[workerName].ControlUnit
		cu.sharedData.Mutex.Unlock()
		cu.log.Info(correctCUName)
		<-cu.simenv.Find(correctCUName).Send(request)
	}
}
