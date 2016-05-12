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
	infrastructureConfig *infrastructure.Config

	ticker  *ticker.Ticker
	network *network.Network

	agents *agent.Environ
}

// New generates new simulation environment
func New(infrastructureFile string) *Simulation {
	sim := &Simulation{
		infrastructureConfig: infrastructure.LoadYAML(infrastructureFile).Parse(),
		agents:               agent.NewEnviron(),
	}

	sim.network = network.New(sim.infrastructureConfig.Network)

	for _, agentConfig := range sim.infrastructureConfig.Network.Agents {
		sim.addAgent(agentConfig)
	}

	sim.ticker = ticker.New(sim.agents.AllAgents())

	log.WithFields(log.Fields{
		"infrastructure": sim.infrastructureConfig.Name,
	}).Info("Simulation initialized")

	return sim
}

func (sim *Simulation) addAgent(agentConfig *infrastructure.Agent) {
	switch agentConfig.Type {
	case "dummy":
		agent.NewDummy(agentConfig, sim.network, sim.agents)
	case "worker":
		agent.NewWorker(agentConfig, sim.network, sim.agents)
	default:
		log.Fatalf("Unknown agent type: %s", agentConfig.Type)
	}
}

// Run starts simulation
func (sim Simulation) Run() {
	sim.ticker.Run()
}
