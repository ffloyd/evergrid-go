package network

import log "github.com/Sirupsen/logrus"

// Network represents all network stats for all machines in simulation
type Network struct {
	segments []Segment
}

// New creates a new network
func New() *Network {
	defer log.Debug("New network initialized")
	return new(Network)
}
