package networkcfg

// NodeCfgYAML is a representation of network node section in YAML infrastructure file
type NodeCfgYAML struct {
	Name   string
	Agents []AgentCfgYAML
}

// NodeCfg is a struct needed to create new network.Node instance
type NodeCfg struct {
	Name   string
	Agents []*AgentCfg

	Segment *SegmentCfg // parent
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (nodeYAML NodeCfgYAML) Parse(parent *SegmentCfg) *NodeCfg {
	node := &NodeCfg{
		Name:    nodeYAML.Name,
		Segment: parent,
	}

	agents := make([]*AgentCfg, len(nodeYAML.Agents))
	for i, agentYAML := range nodeYAML.Agents {
		agents[i] = agentYAML.Parse(node)
	}

	node.Agents = agents

	return node
}
