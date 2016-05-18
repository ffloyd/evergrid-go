package agent

import (
	"math/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulation/simdata/workloadcfg"
)

// Core represents a source of requests to system
// There is should be only one Core in the network.
type Core struct {
	Base
	workload *workloadcfg.WorkloadCfg
}

// NewCore creates a new core agent
func NewCore(config *networkcfg.AgentCfg, net *network.Network, env *Environ, workload *workloadcfg.WorkloadCfg) *Core {
	core := &Core{
		Base:     *NewBase(config, net, env),
		workload: workload,
	}
	env.Cores[core.Name()] = core

	log.WithFields(log.Fields{
		"name": core.Name(),
		"node": core.Node(),
	}).Info("Core agent initialized")
	return core
}

func (core Core) getControlUnit() *ControlUnit {
	count := len(core.env.ControlUnits)
	controlUnits := make([]*ControlUnit, count)
	cuIndex := 0
	for _, cu := range core.env.ControlUnits {
		controlUnits[cuIndex] = cu
		cuIndex++
	}

	return controlUnits[rand.Intn(count)]
}

func (core Core) run() {
	for {
		tick := <-core.tickerChans.Ticks
		core.onNewTick(tick)

		for _, request := range core.workload.Requests[tick] {
			controlUnit := core.getControlUnit()
			controlUnit.incomingRequests <- request
			log.WithFields(log.Fields{
				"tick":         tick,
				"core":         core,
				"control unit": controlUnit,
				"type":         request.Type,
			}).Info("Core sent request to control unit")
		}

		for _, controlUnit := range core.env.ControlUnits {
			log.WithField("control unit", controlUnit).Debug("Send noMoreRequests to control unit")
			controlUnit.noMoreRequests <- true
		}

		core.tickerChans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (core Core) Run() *TickerChans {
	go core.run()
	return core.tickerChans
}
