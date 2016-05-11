package simulation

import (
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/ticker"

	log "github.com/Sirupsen/logrus"
)

// Simulation represents whole simulation environment
type Simulation struct {
	infrastructureConfig *infrastructure.Infrastucture

	ticker  *ticker.Ticker
	network *network.Network
	agents  []agent.Runner
}

// New generates new simulation environment
func New(infraFilename string) *Simulation {
	sim := &Simulation{
		infrastructureConfig: infrastructure.LoadYAML(infraFilename).Parse(),
	}

	sim.network = network.New(sim.infrastructureConfig.Network)

	sim.agents = make([]agent.Runner, len(sim.infrastructureConfig.Network.Agents))
	for i, agentConfig := range sim.infrastructureConfig.Network.Agents {
		sim.agents[i] = agent.New(agentConfig)
		sim.network.Node(agentConfig.Node.Name).AttachAgent(agentConfig.Name, sim.agents[i])
	}

	sim.ticker = ticker.New(sim.agents)

	log.WithFields(log.Fields{
		"infrastructure": sim.infrastructureConfig.Name,
	}).Info("Simulation initialized")

	return sim
}

// Run starts simulation
func (se Simulation) Run() {
	se.ticker.Run()
}
