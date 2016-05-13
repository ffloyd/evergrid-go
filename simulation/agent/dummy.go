package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
)

// Dummy is a simple agent which only writes currentTick in log
type Dummy struct {
	Base
}

// NewDummy creates new dummy agent
func NewDummy(config *networkcfg.Agent, net *network.Network, env *Environ) *Dummy {
	dummy := &Dummy{
		Base: *NewBase(config, net, env),
	}
	env.Dummies[dummy.Name()] = dummy

	log.WithFields(log.Fields{
		"name": dummy.Name(),
		"node": dummy.Node(),
	}).Info("Dummy agent initialized")
	return dummy
}

func (agent Dummy) run() {
	for {
		log.WithFields(log.Fields{
			"tick":  <-agent.tickerChans.Ticks,
			"agent": agent,
		}).Debug("received tick")
		agent.tickerChans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (agent Dummy) Run() *TickerChans {
	go agent.run()
	return agent.tickerChans
}
