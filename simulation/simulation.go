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
	// log.SetLevel(log.DebugLevel)

	sim := &Simulation{
		simData: simdata.Load(simdataFilename),
		agents:  agent.NewEnviron(),
	}

	sim.network = network.New(sim.simData.Network)

	for _, agentConfig := range sim.simData.Network.Agents {
		sim.addAgent(agentConfig)
	}

	sim.initDatasetsAndProcessors()

	sim.ticker = ticker.New(sim.agents.SyncGroup())

	log.WithFields(log.Fields{
		"name": sim.simData.Name,
	}).Info("Simulation initialized")

	return sim
}

func (sim *Simulation) addAgent(agentConfig *networkcfg.AgentCfg) {
	switch agentConfig.Type {
	case networkcfg.AgentDummy:
		agent.NewDummy(agentConfig, sim.network, sim.agents)
	case networkcfg.AgentWorker:
		agent.NewWorker(agentConfig, sim.network, sim.agents)
	case networkcfg.AgentControlUnit:
		agent.NewControlUnit(agentConfig, sim.network, sim.agents)
	case networkcfg.AgentCore:
		agent.NewCore(agentConfig, sim.network, sim.agents, sim.simData.Workload)
	default:
		log.Fatalf("Unknown agent type")
	}
}

func (sim *Simulation) initDatasetsAndProcessors() {
	data := sim.simData.Workload.Data

	for _, datasetConf := range data.Datasets {
		dataset := datasetConf.Info()
		sim.agents.Datasets[string(dataset.UID)] = dataset
	}

	for _, processorConf := range data.Processors {
		processor := processorConf.Info()
		sim.agents.Processors[string(processor.UID)] = processor
	}
}

func (sim *Simulation) report() {
	log.WithField("value", sim.ticker.CurrentTick()).Info("Report: total ticks")
}

// Run starts simulation
func (sim *Simulation) Run() {
	sim.ticker.Run()
	sim.report()
}
