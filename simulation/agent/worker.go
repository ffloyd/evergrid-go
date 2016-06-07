package agent

import (
	log "github.com/Sirupsen/logrus"
	workerUtil "github.com/ffloyd/evergrid-go/simulation/agent/worker"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
)

// Worker is an agent which represents worker machine
type Worker struct {
	Base
	ControlUnit *ControlUnit
	State       *workerUtil.State

	UploadChan  chan workerUtil.ReqUpload
	BuildChan   chan workerUtil.ReqBuild
	ExecuteChan chan workerUtil.ReqExecute

	uploader *workerUtil.Uploader
	builder  *workerUtil.Builder
	executor *workerUtil.Executor
}

// NewWorker creates new worker agent
func NewWorker(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *Worker {
	worker := &Worker{
		Base:        *NewBase(config, net, env),
		State:       workerUtil.NewState(config.Name, config.WorkerDisk, config.WorkerMFlops),
		UploadChan:  make(chan workerUtil.ReqUpload),
		BuildChan:   make(chan workerUtil.ReqBuild),
		ExecuteChan: make(chan workerUtil.ReqExecute),
	}
	env.Workers[worker.Name()] = worker

	worker.ControlUnit = env.ControlUnits[config.ControlUnitName]
	worker.ControlUnit.workers = append(worker.ControlUnit.workers, worker)

	worker.uploader = workerUtil.NewUploader(worker.State)
	worker.builder = workerUtil.NewBuilder(worker.State)
	worker.executor = workerUtil.NewExecutor(worker.State)

	log.WithFields(log.Fields{
		"agent":        worker.Name(),
		"node":         worker.Node(),
		"control_unit": worker.ControlUnit.Name(),
	}).Info("Worker agent initialized")
	return worker
}

func (worker *Worker) run() {
	for {
		worker.sync.toReady()
		worker.sync.toWorking()
		worker.uploader.Process()
		worker.builder.Process()
		worker.executor.Process()
		worker.sync.toIdle()
		doneCh := worker.sync.toDoneCallback()

	SelectLoop:
		for {
			select {
			case request := <-worker.UploadChan:
				worker.uploader.Prepare(request)
			case request := <-worker.BuildChan:
				worker.builder.Prepare(request)
			case request := <-worker.ExecuteChan:
				worker.executor.Prepare(request)
			case <-doneCh:
				worker.sync.SetStopFlag(!worker.State.IsBusy())
				break SelectLoop
			}
		}
	}
}

// Run is implementation of agent.Runner iface
func (worker *Worker) Run() *Synchronizer {
	go worker.run()
	return worker.sync
}
