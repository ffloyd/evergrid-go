package infrastructure

// SegmentYAML is a representation of network segment section in YAML infrastructure file
type SegmentYAML struct {
	Name          string
	InnerBandwith [2]int
	OuterBandwith [2]int
	Nodes         []NodeYAML
}

// Segment is a struct needed to create new network.Segment instance
type Segment struct {
	Name          string
	InnerBandwith [2]int
	OuterBandwith [2]int
	Nodes         []*Node

	Network *Network // parent

	Agents []*Agent // all agents inside segment
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (yamlData SegmentYAML) Parse(parent *Network) *Segment {
	result := &Segment{
		Name:          yamlData.Name,
		InnerBandwith: yamlData.InnerBandwith,
		OuterBandwith: yamlData.OuterBandwith,
		Network:       parent,
	}

	nodes, agentsCount := make([]*Node, len(yamlData.Nodes)), 0
	for i, nodeYAML := range yamlData.Nodes {
		nodes[i] = nodeYAML.Parse(result)
		agentsCount += len(nodes[i].Agents)
	}

	agents := make([]*Agent, 0, agentsCount)
	for _, node := range nodes {
		agents = append(agents, node.Agents...)
	}

	result.Nodes = nodes
	result.Agents = agents

	return result
}
