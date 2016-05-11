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
	agents  map[string]agent.Runner
}

func newNode(config *infrastructure.Node, parent *Segment) *Node {
	node := &Node{
		name:    config.Name,
		segment: parent,
		agents:  make(map[string]agent.Runner),
	}

	log.WithField("name", node.name).Info("Network node initialized")
	return node
}

// AttachAgent adds agent to node's agents list
func (node *Node) AttachAgent(name string, agent agent.Runner) {
	node.agents[name] = agent
}
