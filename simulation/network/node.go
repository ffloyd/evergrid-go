package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/loader"
)

// Node represents particular machine
type Node struct {
	name          string
	outerBandwith Bandwith // bandwith for communication with nodes outside the network segment
	agents        []agent.Runner
}

func newNode(config loader.Node) *Node {
	node := &Node{
		name: config.Name,
		outerBandwith: Bandwith{
			In:  config.OuterBandwith[0],
			Out: config.OuterBandwith[1],
		},
	}

	node.agents = make([]agent.Runner, len(config.Agents))
	for i, agentConfig := range config.Agents {
		node.agents[i] = agent.New(agentConfig)
	}

	log.WithField("name", node.name).Info("Network node initialized")
	return node
}
