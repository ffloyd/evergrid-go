package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/loader"
)

// Node represents particular machine
type Node struct {
	name   string
	agents []agent.Runner
}

func newNode(config loader.Node) *Node {
	node := &Node{
		name: config.Name,
	}

	node.agents = make([]agent.Runner, len(config.Agents))
	for i, agentConfig := range config.Agents {
		node.agents[i] = agent.New(agentConfig)
	}

	log.WithField("name", node.name).Info("Network node initialized")
	return node
}
