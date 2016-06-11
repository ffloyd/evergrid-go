package worker

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/comm"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
)

// Worker -
type Worker struct {
	name   string
	fsm    simenv.AgentFSM
	simenv *simenv.SimEnv
	log    *logrus.Entry

	controlUnitName string
	controlUnit     simenv.Agent

	totalSpace   types.MByte
	freeSpace    types.MByte
	performance  types.MFlop
	pricePerTick float64

	sendLock sync.Mutex
	busy     bool

	uploader uploader

	stats Stats
}

// New -
func New(cfg networkcfg.AgentCfg, logContext *logrus.Entry) *Worker {
	return &Worker{
		name: cfg.Name,
		log:  logContext,

		controlUnitName: cfg.ControlUnitName,

		totalSpace:   cfg.WorkerDisk,
		freeSpace:    cfg.WorkerDisk,
		performance:  cfg.WorkerMFlops,
		pricePerTick: cfg.PricePerTick,
	}
}

// Name -
func (worker *Worker) Name() string {
	return worker.name
}

// Stats -
func (worker *Worker) Stats() Stats {
	return worker.stats
}

// Run -
func (worker *Worker) Run(env *simenv.SimEnv) simenv.AgentChans {
	logContext := worker.log.WithFields(logrus.Fields{
		"agent": worker.Name(),
		"tick":  env.CurrentTick(),
	})

	worker.simenv = env
	worker.controlUnit = env.Find(worker.controlUnitName)
	worker.fsm = simenv.NewAgentFSM(logContext)

	worker.uploader = newUploader(worker, logContext)

	worker.sendLock.Lock()
	go worker.work()
	return worker.fsm.Chans()
}

// Send -
func (worker *Worker) Send(msg interface{}) chan interface{} {
	worker.sendLock.Lock()
	worker.fsm.ToWorking()
	switch request := msg.(type) {
	case comm.WorkerUploadDataset:
		worker.uploader.Prepare(request)
	default:
		worker.log.Panicf("Unknown request type: %v", request)
	}

	response := make(chan interface{})
	go func() {
		worker.fsm.ToIdle()
		response <- simenv.Ok{}
		worker.sendLock.Unlock()
	}()
	return response
}

func (worker *Worker) work() {
	worker.fsm.SetStopFlag(true)
	for {
		worker.fsm.ToReady()
		worker.fsm.ToWorking()
		worker.uploader.Process()
		worker.fsm.ToIdle()
		worker.sendLock.Unlock()
		<-worker.fsm.ToDoneChan()
		worker.sendLock.Lock()
	}
}
