package infrastructure

// AgentYAML is a representation of agent section in YAML infrastructure file
type AgentYAML struct {
	Name string
	Type string
}

// Agent is a struct needed to create new agent.Agent instance
type Agent struct {
	Name string
	Type string

	Node *Node // parent
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (agentYAML AgentYAML) Parse(parent *Node) *Agent {
	return &Agent{
		Name: agentYAML.Name,
		Type: agentYAML.Type,
		Node: parent,
	}
}
