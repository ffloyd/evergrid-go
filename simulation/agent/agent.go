package agent

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
	"github.com/ffloyd/evergrid-go/simulation/network"
)

// Agent interface must be implemented for every agent
type Agent interface {
	fmt.Stringer
	network.Agent
	Run() *Chans
	Node() *network.Node
}

// New initializes agent from config
func New(config *infrastructure.Agent, net *network.Network) Agent {
	var agent Agent

	switch config.Type {
	case "dummy":
		agent = newDummy(config, net)
	default:
		log.Fatalf("Unknown agent type %s", config.Type)
	}

	return agent
}
