package network

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
)

// Node represents particular machine
type Node struct {
	name    string
	segment *Segment
	agents  map[string]Agent
}

func newNode(config *networkcfg.NodeCfg, parent *Segment) *Node {
	node := &Node{
		name:    config.Name,
		segment: parent,
		agents:  make(map[string]Agent),
	}

	log.WithField("name", node.name).Info("Network node initialized")
	return node
}

// AttachAgent adds agent to node's agents list
func (node *Node) AttachAgent(agent Agent) {
	node.agents[agent.Name()] = agent
}

// Segment returns node's network segment
func (node *Node) Segment() *Segment {
	return node.segment
}

// String implements fmt.Stringer interface
func (node Node) String() string {
	return fmt.Sprintf("%s (%s)", node.name, node.segment.name)
}
