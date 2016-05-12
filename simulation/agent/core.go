package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
	"github.com/ffloyd/evergrid-go/simulation/network"
)

// Core represents a source of requests to system
// There is should be only one Core in the network.
type Core struct {
	Base
}

// NewCore creates a new core agent
func NewCore(config *infrastructure.Agent, net *network.Network, env *Environ) *Core {
	core := &Core{
		Base: *NewBase(config, net, env),
	}
	env.Cores[core.Name()] = core

	log.WithFields(log.Fields{
		"name": core.Name(),
		"node": core.Node(),
	}).Info("Core agent initialized")
	return core
}

func (core Core) run() {
	for {
		log.WithFields(log.Fields{
			"tick":  <-core.tickerChans.Ticks,
			"agent": core,
		}).Debug("received tick")
		core.tickerChans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (core Core) Run() *TickerChans {
	go core.run()
	return core.tickerChans
}