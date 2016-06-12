package controlunit

import (
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
	workers     []simenv.Agent
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

	cu.workers = make([]simenv.Agent, len(cu.workerNames))
	for i, workerName := range cu.workerNames {
		cu.workers[i] = env.Find(workerName)
	}

	cu.monitor.Run()
	cu.scheduler.Run()
	go cu.work()
	return cu.fsm.Chans()
}

// Send - respond means that request arrived to proper scheduler
func (cu *ControlUnit) Send(msg interface{}) chan interface{} {
	cu.sendLock.Lock()
	cu.fsm.ToWorking()

	schedChans := cu.scheduler.RequestChans()

	switch request := msg.(type) {
	case comm.ControlUnitUploadDataset:
		schedChans.UploadDataset <- scheduler.ReqUploadDataset{
			Dataset: request.Dataset,
		}

		if <-schedChans.DelegateToLeader {
			<-cu.sharedData.LeaderControlUnit.Send(request)
		}
	case comm.ControlUnitRunExperiment:
		schedChans.RunExperiment <- scheduler.ReqRunExperiment{
			Calculator: request.Calculator,
			Dataset:    request.Dataset,
		}

		if <-schedChans.DelegateToLeader {
			<-cu.sharedData.LeaderControlUnit.Send(request)
		}
	case scheduler.DoUploadDataset:
		cu.processUploadDataset(request)
	case scheduler.DoBuildCalculator:
		cu.processBuildCalculator(request)
	case scheduler.DoRunCalculator:
		cu.processRunCalculator(request)
	default:
		cu.log.Panicf("Unknown request type: %v", request)
	}

	cu.fsm.ToIdle()
	response := make(chan interface{})
	go func() {
		response <- simenv.Ok{}
	}()
	cu.sendLock.Unlock()
	return response
}

func (cu *ControlUnit) work() {
	cu.fsm.SetStopFlag(true)
	chans := cu.scheduler.ControlChans()
	cu.sendLock.Lock()
	for {
		cu.fsm.ToReady()
		cu.fsm.ToIdle()
		cu.sendLock.Unlock()

		doneChan := cu.fsm.ToDoneChan()

	SelectLoop:
		for {
			select {
			case request := <-chans.UploadDataset:
				cu.processUploadDataset(request)
			case request := <-chans.BuildCalculator:
				cu.processBuildCalculator(request)
			case request := <-chans.RunCalculator:
				cu.processRunCalculator(request)
			case <-doneChan:
				cu.sendLock.Lock()
				break SelectLoop
			}
		}
	}
}

func (cu *ControlUnit) processUploadDataset(request scheduler.DoUploadDataset) {
	cu.log.Info(request)
}

func (cu *ControlUnit) processBuildCalculator(request scheduler.DoBuildCalculator) {
	cu.log.Info(request)
}

func (cu *ControlUnit) processRunCalculator(request scheduler.DoRunCalculator) {
	cu.log.Info(request)
}

func (cu *ControlUnit) amILeader() bool {
	return cu.Name() == cu.sharedData.LeaderControlUnit.Name()
}
