package agent

import (
	log "github.com/Sirupsen/logrus"

	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
)

// Base is a common part for all types of agents in simulation
// Also, only this part exported to network package via interface
type Base struct {
	name        string
	node        *network.Node
	tickerChans *TickerChans
	env         *Environ
	tick        int // current tick
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

func (agent *Base) startTick() {
	agent.tick = <-agent.tickerChans.Ticks
	log.WithFields(log.Fields{
		"tick":  agent.tick,
		"agent": agent,
	}).Debug("agent ready for tick")
	agent.tickerChans.Ready <- true // tick begin sync
}

func (agent Base) finishTick() {
	log.WithFields(log.Fields{
		"tick":  agent.tick,
		"agent": agent,
	}).Debug("agent finished tick")
	agent.tickerChans.Ready <- true // tick end sync
}

// NewBase is common initialization part all agents
func NewBase(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *Base {
	node := net.Node(config.Node.Name)
	base := &Base{
		name:        config.Name,
		node:        node,
		env:         env,
		tickerChans: NewTickerChans(),
	}
	node.AttachAgent(base)
	return base
}
