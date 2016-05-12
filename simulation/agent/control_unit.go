package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
	"github.com/ffloyd/evergrid-go/simulation/network"
)

// ControlUnit is a representation of control unit app
type ControlUnit struct {
	Base

	workers []*Worker
}

// NewControlUnit creates a new control unit
func NewControlUnit(config *infrastructure.Agent, net *network.Network, env *Environ) *ControlUnit {
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
