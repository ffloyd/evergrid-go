package simulation

import (
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/ticker"
)

// Simulation represents whole simulation environment
type Simulation struct {
	ticker  *ticker.Ticker
	network *network.Network
	agents  []agent.Runner
}

// New generates new simulation environment
func New() *Simulation {
	se := new(Simulation)

	se.agents = []agent.Runner{
		agent.NewDummy("Dummy 1"),
		agent.NewDummy("Dummy 2"),
	}

	se.network = network.New()

	se.ticker = ticker.New(se.agents) // ticker initialization starts all agents

	return se
}

// Run starts simulation
func (se Simulation) Run() {
	se.ticker.Run()
}
