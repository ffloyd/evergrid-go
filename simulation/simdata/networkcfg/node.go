package networkcfg

// NodeYAML is a representation of network node section in YAML infrastructure file
type NodeYAML struct {
	Name   string
	Agents []AgentYAML
}

// Node is a struct needed to create new network.Node instance
type Node struct {
	Name   string
	Agents []*Agent

	Segment *Segment // parent
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (nodeYAML NodeYAML) Parse(parent *Segment) *Node {
	node := &Node{
		Name:    nodeYAML.Name,
		Segment: parent,
	}

	agents := make([]*Agent, len(nodeYAML.Agents))
	for i, agentYAML := range nodeYAML.Agents {
		agents[i] = agentYAML.Parse(node)
	}

	node.Agents = agents

	return node
}
