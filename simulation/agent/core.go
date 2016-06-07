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
		"agent": core.Name(),
		"node":  core.Node(),
	}).Info("Core agent initialized")
	return core
}

func (core *Core) getControlUnit() *ControlUnit {
	count := len(core.env.ControlUnits)
	controlUnits := make([]*ControlUnit, count)
	cuIndex := 0
	for _, cu := range core.env.ControlUnits {
		controlUnits[cuIndex] = cu
		cuIndex++
	}

	return controlUnits[rand.Intn(count)]
}

func (core *Core) run() {
	activeTicksProcessed := 0

	for {
		core.sync.toReady()
		core.sync.toWorking()

		for _, request := range core.workload.Requests[core.sync.tick] {
			controlUnit := core.getControlUnit()
			log.WithFields(log.Fields{
				"tick":         core.sync.tick,
				"agent":        core,
				"control unit": controlUnit,
				"type":         request.Type,
			}).Info("Core seding request to Control Unit")
			controlUnit.incomingRequests <- request
			<-controlUnit.requestConfirmation
		}

		if core.workload.Requests[core.sync.tick] != nil {
			activeTicksProcessed++
		}

		core.sync.toIdle()
		<-core.sync.toDoneCallback()
		core.sync.SetStopFlag(activeTicksProcessed == len(core.workload.Requests))
	}
}

// Run is implementation of agent.Runner iface
func (core *Core) Run() *Synchronizer {
	go core.run()
	return core.sync
}
