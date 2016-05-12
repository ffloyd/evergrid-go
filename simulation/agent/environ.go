package agent

// Environ is a set of all agents in the system
type Environ struct {
	Dummies map[string]*Dummy
	Workers map[string]*Worker
}

// NewEnviron is a simple initializer
func NewEnviron() *Environ {
	return &Environ{
		Dummies: make(map[string]*Dummy),
		Workers: make(map[string]*Worker),
	}
}

// AllAgents returns slice of all agents
func (env Environ) AllAgents() []Agent {
	agentsCount := len(env.Dummies) + len(env.Workers)
	agents := make([]Agent, agentsCount)
	i := 0

	for _, dummy := range env.Dummies {
		agents[i] = dummy
		i++
	}

	for _, worker := range env.Workers {
		agents[i] = worker
	}

	return agents
}
