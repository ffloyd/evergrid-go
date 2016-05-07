package simulation

import (
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/loader"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/ticker"

	log "github.com/Sirupsen/logrus"
)

// Simulation represents whole simulation environment
type Simulation struct {
	ticker  *ticker.Ticker
	network *network.Network
	agents  []agent.Runner
}

// New generates new simulation environment
func New() *Simulation {
	infrastructure := loader.LoadInfrastructure("simdata/infrastructure/small.json")
	log.Info(infrastructure)

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
