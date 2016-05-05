package simenv

// AgentChans is struct for communication between GlobalTimer and agents
type AgentChans struct {
	ready chan bool // for incoming ready status
	ticks chan int  // for ticks broadcasting
}

// NewAgentChans initializes AgentChans instanse
func NewAgentChans() *AgentChans {
	return &AgentChans{make(chan bool), make(chan int)}
}
