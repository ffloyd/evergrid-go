package simenv

// AgentGroup -
type AgentGroup struct {
	agents map[string]Agent
}

// NewAgentGroup -
func NewAgentGroup() *AgentGroup {
	return &AgentGroup{
		agents: make(map[string]Agent),
	}
}

// Add -
func (group *AgentGroup) Add(agent Agent) Error {
	return Error{}
}

// Run -
func (group *AgentGroup) Run() AgentChans {
	return AgentChans{}
}
