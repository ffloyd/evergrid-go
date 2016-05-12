package agent

// Environ is a set of all agents in the system
type Environ struct {
	Dummies      map[string]*Dummy
	Workers      map[string]*Worker
	ControlUnits map[string]*ControlUnit
	Cores        map[string]*Core
}

// NewEnviron is a simple initializer
func NewEnviron() *Environ {
	return &Environ{
		Dummies:      make(map[string]*Dummy),
		Workers:      make(map[string]*Worker),
		ControlUnits: make(map[string]*ControlUnit),
		Cores:        make(map[string]*Core),
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
