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
	se.ticker = ticker.New()
	se.network = network.New()

	se.agents = []agent.Runner{
		agent.Dummy{Name: "Agent 1"},
		agent.Dummy{Name: "Agent 2"},
	}

	for _, ag := range se.agents {
		se.ticker.AddAgent(ag)
	}

	return se
}

// Run starts simulation
func (se Simenv) Run() {
	se.ticker.Run()
}
