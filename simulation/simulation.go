package simulation

import (
	"github.com/ffloyd/evergrid-go/simulation/agent"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulation/ticker"

	log "github.com/Sirupsen/logrus"
)

// Simulation represents whole simulation environment
type Simulation struct {
	simData *simdata.SimData

	ticker  *ticker.Ticker
	network *network.Network

	agents *agent.Environ
}

// New generates new simulation environment
func New(simdataFilename string) *Simulation {
	sim := &Simulation{
		simData: simdata.Load(simdataFilename),
		agents:  agent.NewEnviron(),
	}

	sim.network = network.New(sim.simData.Network)

	for _, agentConfig := range sim.simData.Network.Agents {
		sim.addAgent(agentConfig)
	}

	sim.ticker = ticker.New(sim.agents.AllAgents())

	log.WithFields(log.Fields{
		"infrastructure": sim.simData.Name,
	}).Info("Simulation initialized")

	return sim
}

func (sim *Simulation) addAgent(agentConfig *networkcfg.AgentCfg) {
	switch agentConfig.Type {
	case "dummy":
		agent.NewDummy(agentConfig, sim.network, sim.agents)
	case "worker":
		agent.NewWorker(agentConfig, sim.network, sim.agents)
	case "control_unit":
		agent.NewControlUnit(agentConfig, sim.network, sim.agents)
	case "core":
		agent.NewCore(agentConfig, sim.network, sim.agents)
	default:
		log.Fatalf("Unknown agent type: %s", agentConfig.Type)
	}
}

// Run starts simulation
func (sim Simulation) Run() {
	sim.ticker.Run()
}
