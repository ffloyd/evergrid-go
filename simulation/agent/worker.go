package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
)

// Worker is an agent which represents processes
type Worker struct {
	Base
}

// NewWorker creates new worker agent
func NewWorker(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *Worker {
	worker := &Worker{
		Base: *NewBase(config, net, env),
	}
	env.Workers[worker.Name()] = worker

	log.WithFields(log.Fields{
		"name": worker.Name(),
		"node": worker.Node(),
	}).Info("Worker agent initialized")
	return worker
}

func (worker Worker) run() {
	for {
		worker.startTick()
		worker.finishTick()
	}
}

// Run is implementation of agent.Runner iface
func (worker Worker) Run() *TickerChans {
	go worker.run()
	return worker.tickerChans
}
