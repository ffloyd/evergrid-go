package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
	"github.com/ffloyd/evergrid-go/simulation/network"
)

// Worker is an agent which represents processes
type Worker struct {
	Base
}

// NewWorker creates new worker agent
func NewWorker(config *infrastructure.Agent, net *network.Network, env *Environ) *Worker {
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
		log.WithFields(log.Fields{
			"tick":  <-worker.tickerChans.Ticks,
			"agent": worker,
		}).Debug("received tick")
		worker.tickerChans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (worker Worker) Run() *TickerChans {
	go worker.run()
	return worker.tickerChans
}
