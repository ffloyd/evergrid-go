package networkcfg

// SegmentCfgYAML is a representation of network segment section in YAML infrastructure file
type SegmentCfgYAML struct {
	Name          string
	InnerBandwith [2]int
	OuterBandwith [2]int
	Nodes         []NodeCfgYAML
}

// SegmentCfg is a struct needed to create new network.Segment instance
type SegmentCfg struct {
	Name          string
	InnerBandwith [2]int
	OuterBandwith [2]int
	Nodes         []*NodeCfg

	Network *NetworkCfg // parent

	Agents []*AgentCfg // all agents inside segment
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (segmentYAML SegmentCfgYAML) Parse(parent *NetworkCfg) *SegmentCfg {
	segment := &SegmentCfg{
		Name:          segmentYAML.Name,
		InnerBandwith: segmentYAML.InnerBandwith,
		OuterBandwith: segmentYAML.OuterBandwith,
		Network:       parent,
	}

	nodes, agentsCount := make([]*NodeCfg, len(segmentYAML.Nodes)), 0
	for i, nodeYAML := range segmentYAML.Nodes {
		nodes[i] = nodeYAML.Parse(segment)
		agentsCount += len(nodes[i].Agents)
	}

	agents := make([]*AgentCfg, 0, agentsCount)
	for _, node := range nodes {
		agents = append(agents, node.Agents...)
	}

	segment.Nodes = nodes
	segment.Agents = agents

	return segment
}
