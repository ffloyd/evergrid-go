package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
)

// ControlUnit is a representation of control unit app
type ControlUnit struct {
	Base

	workers []*Worker
}

// NewControlUnit creates a new control unit
func NewControlUnit(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *ControlUnit {
	unit := &ControlUnit{
		Base: *NewBase(config, net, env),
	}
	env.ControlUnits[unit.Name()] = unit

	log.WithFields(log.Fields{
		"name": unit.Name(),
		"node": unit.Node(),
	}).Info("Control Unit agent initialized")
	return unit
}

func (unit ControlUnit) run() {
	for {
		log.WithFields(log.Fields{
			"tick":  <-unit.tickerChans.Ticks,
			"agent": unit,
		}).Debug("received tick")
		unit.tickerChans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (unit ControlUnit) Run() *TickerChans {
	go unit.run()
	return unit.tickerChans
}
