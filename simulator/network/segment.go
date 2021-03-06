package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
)

// Segment represents a local connected scope of machines. As example - if they are part of same DigitalOcean region.
type Segment struct {
	name          string
	innerBandwith Bandwith // bandwith for communication inside this segment
	outerBandwith Bandwith // bandwith for communication with nodes outside the network segment
	network       *Network

	nodes map[string]*Node
}

func newSegment(config *networkcfg.SegmentCfg, parent *Network) *Segment {
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
		network: parent,
		nodes:   make(map[string]*Node),
	}

	for _, nodeConfig := range config.Nodes {
		segment.nodes[nodeConfig.Name] = newNode(nodeConfig, segment)
	}

	log.WithField("name", segment.name).Info("Network segment initialized")

	return segment
}

// AgentNames returns all agent names inside segment
func (seg *Segment) AgentNames() []string {
	var result []string

	for _, node := range seg.nodes {
		for _, agent := range node.agents {
			result = append(result, agent.Name())
		}
	}

	return result
}

// Bandwith return inner or outer bandwith depends on argument
func (seg *Segment) Bandwith(inner bool) Bandwith {
	if inner {
		return seg.innerBandwith
	}

	return seg.outerBandwith
}
