package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
)

// Network represents all network stats for all machines in simulation.
type Network struct {
	name     string
	segments map[string]*Segment

	nodes map[string]*Node
}

// New creates a new network basing on config
func New(config *networkcfg.NetworkCfg) *Network {
	network := &Network{
		name:     config.Name,
		segments: make(map[string]*Segment),
		nodes:    make(map[string]*Node),
	}

	for _, segmentConfig := range config.Segments {
		segment := newSegment(segmentConfig, network)
		network.segments[segmentConfig.Name] = segment
		for name, node := range segment.nodes {
			network.nodes[name] = node
		}
	}

	log.WithField("name", network.name).Info("Network initialized")
	return network
}

// Node returns node by its name
func (network Network) Node(name string) *Node {
	return network.nodes[name]
}
