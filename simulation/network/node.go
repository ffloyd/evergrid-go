package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
)

// Node represents particular machine
type Node struct {
	name    string
	segment *Segment
	agents  map[string]agent.Agent
}

func newNode(config *infrastructure.Node, parent *Segment) *Node {
	node := &Node{
		name:    config.Name,
		segment: parent,
		agents:  make(map[string]agent.Agent),
	}

	log.WithField("name", node.name).Info("Network node initialized")
	return node
}

// AttachAgent adds agent to node's agents list
func (node *Node) AttachAgent(agent agent.Agent) {
	node.agents[agent.Name()] = agent
}
