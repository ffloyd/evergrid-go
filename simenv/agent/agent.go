package agent

// Agent is a common part for all types of agents in simulation
type Agent struct {
	name  string
	chans *Chans
}

// String for implement Stringer interface
func (agent Agent) String() string {
	return agent.name
}
