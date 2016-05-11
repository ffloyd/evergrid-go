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

	Infrastucture *Config // parent

	Nodes  []*Node  // all nodes inside network
	Agents []*Agent // all agents inside network
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (networkYAML NetworkYAML) Parse(parent *Config) *Network {
	networkCfg := &Network{
		Name:          networkYAML.Name,
		Infrastucture: parent,
	}

	segments, nodesCount, agentsCount := make([]*Segment, len(networkYAML.Segments)), 0, 0
	for i, segmentYAML := range networkYAML.Segments {
		segments[i] = segmentYAML.Parse(networkCfg)
		nodesCount += len(segments[i].Nodes)
		agentsCount += len(segments[i].Agents)
	}

	nodes, agents := make([]*Node, 0, nodesCount), make([]*Agent, 0, agentsCount)
	for _, segment := range segments {
		nodes, agents = append(nodes, segment.Nodes...), append(agents, segment.Agents...)
	}

	networkCfg.Segments = segments
	networkCfg.Nodes = nodes
	networkCfg.Agents = agents

	return networkCfg
}
