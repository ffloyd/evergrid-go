package networkcfg

// AgentCfgYAML is a representation of agent section in YAML infrastructure file
type AgentCfgYAML struct {
	Name string
	Type string
}

// AgentCfg is a struct needed to create new agent.Agent instance
type AgentCfg struct {
	Name string
	Type string

	Node *NodeCfg // parent
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (agentYAML AgentCfgYAML) Parse(parent *NodeCfg) *AgentCfg {
	return &AgentCfg{
		Name: agentYAML.Name,
		Type: agentYAML.Type,
		Node: parent,
	}
}
