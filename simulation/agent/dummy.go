package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
	"github.com/ffloyd/evergrid-go/simulation/network"
)

// Dummy is a simple agent which only writes currentTick in log
type Dummy struct {
	Base
}

// NewDummy creates new dummy agent
func NewDummy(config *infrastructure.Agent, net *network.Network, env *Environ) *Dummy {
	dummy := &Dummy{
		Base: *NewBase(config, net, env),
	}
	env.Dummies[dummy.Name()] = dummy

	log.WithField("name", dummy).Info("Dummy agent initialized")
	return dummy
}

func (agent Dummy) run() {
	for {
		log.WithFields(log.Fields{
			"tick":  <-agent.chans.Ticks,
			"agent": agent,
		}).Debug("received tick")
		agent.chans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (agent Dummy) Run() *Chans {
	go agent.run()
	return agent.chans
}
