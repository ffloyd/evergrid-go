package agent

import (
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
	"github.com/ffloyd/evergrid-go/simulation/network"
)

// Base is a common part for all types of agents in simulation
// Also, only this part exported to network package via interface
type Base struct {
	name  string
	node  *network.Node
	chans *Chans
	env   *Environ
}

// String for implement Stringer interface
func (agent Base) String() string {
	return agent.name
}

// Name needed for network.Agent interface implementation
func (agent Base) Name() string {
	return agent.name
}

// Node needed for agent.Agent interface
func (agent Base) Node() *network.Node {
	return agent.node
}

// NewBase is common initialization part all agents
func NewBase(config *infrastructure.Agent, net *network.Network, env *Environ) *Base {
	node := net.Node(config.Node.Name)
	base := &Base{
		name:  config.Name,
		node:  node,
		env:   env,
		chans: NewChans(),
	}
	node.AttachAgent(base)
	return base
}