package infrastructure

// NetworkYAML is a representation of network section in YAML infrastructure file
type NetworkYAML struct {
	Name     string
	Segments []SegmentYAML
}

// Network is a struct needed to create new network.Network instance
type Network struct {
	Name     string
	Segments []*Segment

	Infrastucture *Infrastucture // parent

	Nodes  []*Node  // all nodes inside network
	Agents []*Agent // all agents inside network
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (yamlData NetworkYAML) Parse(parent *Infrastucture) *Network {
	result := &Network{
		Name:          yamlData.Name,
		Infrastucture: parent,
	}

	segments, nodesCount, agentsCount := make([]*Segment, len(yamlData.Segments)), 0, 0
	for i, segmentYAML := range yamlData.Segments {
		segments[i] = segmentYAML.Parse(result)
		nodesCount += len(segments[i].Nodes)
		agentsCount += len(segments[i].Agents)
	}

	nodes, agents := make([]*Node, 0, nodesCount), make([]*Agent, 0, agentsCount)
	for _, segment := range segments {
		nodes, agents = append(nodes, segment.Nodes...), append(agents, segment.Agents...)
	}

	result.Segments = segments
	result.Nodes = nodes
	result.Agents = agents

	return result
}
