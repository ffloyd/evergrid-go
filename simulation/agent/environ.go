package agent

import (
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simulation/ticker"
)

// Environ is a set of all agents, known datasets and processors in the system
type Environ struct {
	Dummies      map[string]*Dummy
	Workers      map[string]*Worker
	ControlUnits map[string]*ControlUnit
	Cores        map[string]*Core

	Datasets   map[string]*types.DatasetInfo
	Processors map[string]*types.ProcessorInfo

	leaderControlUnit *ControlUnit
}

// NewEnviron is a simple initializer
func NewEnviron() *Environ {
	return &Environ{
		Dummies:      make(map[string]*Dummy),
		Workers:      make(map[string]*Worker),
		ControlUnits: make(map[string]*ControlUnit),
		Cores:        make(map[string]*Core),
		Datasets:     make(map[string]*types.DatasetInfo),
		Processors:   make(map[string]*types.ProcessorInfo),
	}
}

// AllAgents returns slice of all agents
func (env Environ) AllAgents() []Agent {
	agentsCount := len(env.Dummies) + len(env.Workers) + len(env.ControlUnits) + len(env.Cores)
	agents := make([]Agent, agentsCount)
	i := 0

	for _, dummy := range env.Dummies {
		agents[i] = dummy
		i++
	}

	for _, worker := range env.Workers {
		agents[i] = worker
		i++
	}

	for _, unit := range env.ControlUnits {
		agents[i] = unit
		i++
	}

	for _, core := range env.Cores {
		agents[i] = core
		i++
	}

	return agents
}

// SyncGroup return a ticker.SyncGroup with all agents inside
func (env Environ) SyncGroup() ticker.SyncGroup {
	agents := env.AllAgents()
	result := make(ticker.SyncGroup, len(agents))
	for i, agent := range agents {
		result[i] = agent.Run()
	}
	return result
}

// LeaderControlUnit returns a current leader between control units
func (env *Environ) LeaderControlUnit() *ControlUnit {
	if env.leaderControlUnit == nil {
		var firstCU *ControlUnit
		for _, cu := range env.ControlUnits {
			firstCU = cu
			break
		}
		env.leaderControlUnit = firstCU
	}

	return env.leaderControlUnit
}
