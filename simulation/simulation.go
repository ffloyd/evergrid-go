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
	infrastructureName string

	ticker  *ticker.Ticker
	network *network.Network
	agents  []agent.Runner
}

// New generates new simulation environment
func New() *Simulation {
	se := new(Simulation)

	infrastructure := loader.LoadInfrastructure("simdata/infrastructure/small.json")
	se.infrastructureName = infrastructure.Name

	se.agents = []agent.Runner{
		agent.NewDummy("Dummy 1"),
		agent.NewDummy("Dummy 2"),
	}

	se.network = network.New()

	se.ticker = ticker.New(se.agents) // ticker initialization starts all agents

	log.WithFields(log.Fields{
		"infrastructure": se.infrastructureName,
	}).Info("Simulation initialized")

	return se
}

// Run starts simulation
func (se Simulation) Run() {
	se.ticker.Run()
}
