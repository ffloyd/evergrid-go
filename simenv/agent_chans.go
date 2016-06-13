package simenv

/*
AgentChans это каналы, по которым агент общается с ядром.

Не должны создаваться или использоваться напрямую - интерфейс работы с ними реализован в AgentFSM
*/
type AgentChans struct {
	statusChan     chan AgentState
	startWorkChan  chan Ok
	finishWorkChan chan Ok
	stopFlagChan   chan bool
}

// NewAgentChans -
func newAgentChans() AgentChans {
	return AgentChans{
		statusChan:     make(chan AgentState),
		startWorkChan:  make(chan Ok),
		finishWorkChan: make(chan Ok),
		stopFlagChan:   make(chan bool),
	}
}
