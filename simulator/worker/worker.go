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
	builder  builder
	executor executor

	stats Stats
}

// New -
func New(cfg networkcfg.AgentCfg, logContext *logrus.Entry) *Worker {
	if cfg.Type != networkcfg.AgentWorker {
		logContext.Panic("Wrong agent type in config")
	}

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
	worker.log = worker.log.WithFields(logrus.Fields{
		"agent": worker.Name(),
		"tick":  env.CurrentTick(),
	})

	worker.simenv = env
	worker.controlUnit = env.Find(worker.controlUnitName)
	worker.fsm = simenv.NewAgentFSM(worker.log)

	worker.uploader = newUploader(worker, worker.log)
	worker.builder = newBuilder(worker, worker.log)
	worker.executor = newExecutor(worker, worker.log)

	worker.sendLock.Lock()
	go worker.work()
	return worker.fsm.Chans()
}

// Send -
func (worker *Worker) Send(msg interface{}) chan interface{} {
	worker.sendLock.Lock()
	worker.fsm.ToWorking()

	var responseMsg interface{}
	responseMsg = simenv.Ok{}

	switch request := msg.(type) {
	case comm.WorkerUploadDataset:
		worker.busyCheck()
		worker.uploader.Prepare(request)
	case comm.WorkerBuildCalculator:
		worker.busyCheck()
		worker.builder.Prepare(request)
	case comm.WorkerRunCalculator:
		worker.busyCheck()
		worker.executor.Prepare(request)
	case comm.WorkerBusy:
		responseMsg = worker.busy
	case comm.WorkerInfo:
		responseMsg = worker.getInfo()
	default:
		worker.log.Panicf("Unknown request type: %v", request)
	}

	responseChan := make(chan interface{})
	go func() {
		worker.fsm.ToIdle()
		responseChan <- responseMsg
		worker.sendLock.Unlock()
	}()
	return responseChan
}

func (worker *Worker) work() {
	worker.fsm.SetStopFlag(true)
	for {
		worker.fsm.ToReady()

		worker.fsm.ToWorking()
		worker.uploader.Process()
		worker.builder.Process()
		worker.executor.Process()
		worker.fsm.ToIdle()

		worker.sendLock.Unlock()
		<-worker.fsm.ToDoneChan()

		worker.sendLock.Lock()
	}
}

func (worker *Worker) busyCheck() {
	if worker.busy {
		worker.log.Panic("Incorrect request to busy worker")
	}
}

func (worker *Worker) getInfo() types.WorkerInfo {
	return types.WorkerInfo{
		UID:            worker.Name(),
		Busy:           worker.busy,
		MFlops:         worker.performance,
		TotalDiskSpace: worker.totalSpace,
		FreeDiskSpace:  worker.freeSpace,
		PricePerTick:   worker.pricePerTick,
		ControlUnit:    worker.controlUnitName,
	}
}
