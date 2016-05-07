package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/loader"
)

// Network represents all network stats for all machines in simulation.
type Network struct {
	name     string
	segments []*Segment
}

// New creates a new network basing on config
func New(config loader.Network) *Network {
	network := new(Network)
	network.name = config.Name

	network.segments = make([]*Segment, len(config.Segments))
	for i, segmentConfig := range config.Segments {
		network.segments[i] = newSegment(segmentConfig)
	}

	log.WithField("name", network.name).Info("Network initialized")
	return network
}

// Agents return list of all agents inside network
func (net Network) Agents() []agent.Runner {
	var result []agent.Runner
	for _, segment := range net.segments {
		result = append(result, segment.agents()...)
	}
	return result
}
