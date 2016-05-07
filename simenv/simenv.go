package simenv

import (
	"github.com/ffloyd/evergrid-go/simenv/agent"
	"github.com/ffloyd/evergrid-go/simenv/network"
	"github.com/ffloyd/evergrid-go/simenv/ticker"
)

// Simenv represents whole simulation environment
type Simenv struct {
	ticker  *ticker.Ticker
	network *network.Network
	agents  []agent.Runner
}

// New generates new simulation environment
func New() *Simenv {
	se := new(Simenv)

	se.agents = []agent.Runner{
		agent.NewDummy("Dummy 1"),
		agent.NewDummy("Dummy 2"),
	}

	se.network = network.New()

	se.ticker = ticker.New(se.agents) // ticker initialization starts all agents

	return se
}

// Run starts simulation
func (se Simenv) Run() {
	se.ticker.Run()
}
