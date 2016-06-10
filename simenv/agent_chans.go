package simenv

/*
AgentChans -
*/
type AgentChans struct {
	ticksChan      chan int
	statusChan     chan AgentStatus
	startWorkChan  chan Ok
	finishWorkChan chan Ok
	stopFlagChan   chan bool
}

// NewAgentChans -
func NewAgentChans() AgentChans {
	return AgentChans{
		ticksChan:      make(chan int),
		statusChan:     make(chan AgentStatus),
		startWorkChan:  make(chan Ok),
		finishWorkChan: make(chan Ok),
		stopFlagChan:   make(chan bool),
	}
}
