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
func New(infraFilename string) *Simulation {
	infrastructure := loader.LoadInfrastructure(infraFilename)

	sim := &Simulation{
		infrastructureName: infrastructure.Name,
		network:            network.New(infrastructure.Network),
	}

	sim.agents = sim.network.Agents()
	sim.ticker = ticker.New(sim.agents)

	log.WithFields(log.Fields{
		"infrastructure": sim.infrastructureName,
	}).Info("Simulation initialized")

	return sim
}

// Run starts simulation
func (se Simulation) Run() {
	se.ticker.Run()
}
