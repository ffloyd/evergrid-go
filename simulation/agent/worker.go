package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
)

// Worker is an agent which represents processes
type Worker struct {
	Base
	ControlUnit *ControlUnit
	State       types.WorkerInfo
}

// NewWorker creates new worker agent
func NewWorker(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *Worker {
	worker := &Worker{
		Base: *NewBase(config, net, env),
		State: types.WorkerInfo{
			UID:            types.UID(config.Name),
			Busy:           false,
			MFlops:         config.WorkerMFlops,
			TotalDiskSpace: config.WorkerDisk,
			FreeDiskSpace:  config.WorkerDisk,
		},
	}
	env.Workers[worker.Name()] = worker

	worker.ControlUnit = env.ControlUnits[config.ControlUnitName]
	worker.ControlUnit.workers = append(worker.ControlUnit.workers, worker)

	log.WithFields(log.Fields{
		"agent":        worker.Name(),
		"node":         worker.Node(),
		"control_unit": worker.ControlUnit.Name(),
	}).Info("Worker agent initialized")
	return worker
}

func (worker Worker) run() {
	for {
		worker.sync.toReady()
		worker.sync.toIdle()
		<-worker.sync.toDoneCallback()
	}
}

// Run is implementation of agent.Runner iface
func (worker Worker) Run() *Synchronizer {
	go worker.run()
	return worker.sync
}
