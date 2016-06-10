package simenv

/*
AgentChans -
*/
type AgentChans struct {
	statusChan     chan AgentState
	startWorkChan  chan Ok
	finishWorkChan chan Ok
	stopFlagChan   chan bool
}

// NewAgentChans -
func NewAgentChans() AgentChans {
	return AgentChans{
		statusChan:     make(chan AgentState),
		startWorkChan:  make(chan Ok),
		finishWorkChan: make(chan Ok),
		stopFlagChan:   make(chan bool),
	}
}
