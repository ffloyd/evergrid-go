package infrastructure

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
func (yamlData NodeYAML) Parse(parent *Segment) *Node {
	result := &Node{
		Name:    yamlData.Name,
		Segment: parent,
	}

	agents := make([]*Agent, len(yamlData.Agents))
	for i, agentYAML := range yamlData.Agents {
		agents[i] = agentYAML.Parse(result)
	}

	result.Agents = agents

	return result
}
