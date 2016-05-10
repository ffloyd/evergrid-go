package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/loader"
)

// Segment represents a local connected scope of machines. As example - if they are part of same DigitalOcean region.
type Segment struct {
	name          string
	innerBandwith Bandwith // bandwith for communication inside this segment
	outerBandwith Bandwith // bandwith for communication with nodes outside the network segment

	nodes []*Node
}

func newSegment(config loader.Segment) *Segment {
	segment := &Segment{
		name: config.Name,
		innerBandwith: Bandwith{
			In:  config.InnerBandwith[0],
			Out: config.InnerBandwith[1],
		},
		outerBandwith: Bandwith{
			In:  config.OuterBandwith[0],
			Out: config.OuterBandwith[1],
		},
	}

	segment.nodes = make([]*Node, len(config.Nodes))
	for i, nodeConfig := range config.Nodes {
		segment.nodes[i] = newNode(nodeConfig)
	}

	log.WithField("name", segment.name).Info("Network segment initialized")

	return segment
}

// agents return list of all agents inside segment
func (segment Segment) agents() []agent.Runner {
	var result []agent.Runner
	for _, node := range segment.nodes {
		result = append(result, node.agents...)
	}
	return result
}
