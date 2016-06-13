// Package simulator отвечает за запуск симуляции архитектуры Evergrid
package simulator

import (
	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/scheduler"
	"github.com/ffloyd/evergrid-go/scheduler/random"
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/controlunit"
	"github.com/ffloyd/evergrid-go/simulator/core"
	"github.com/ffloyd/evergrid-go/simulator/network"
	"github.com/ffloyd/evergrid-go/simulator/simdata"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulator/worker"
)

/*
Simulator - это струтура содержащая всю информацию о симуляции.
*/
type Simulator struct {
	simData *simdata.SimData

	network *network.Network
	simenv  *simenv.SimEnv

	logContext *logrus.Entry
	sharedData *controlunit.SharedData

	cuNames   []string
	cuWorkers map[string][]string
}

// New - инициализация симуляции на основе корневого файла сценария.
func New(simdataFilename string) *Simulator {
	sim := &Simulator{
		simData:   simdata.Load(simdataFilename),
		simenv:    simenv.New(),
		cuWorkers: make(map[string][]string),
	}

	sim.logContext = logrus.WithField("simulation", sim.simData.Name)
	sim.network = network.New(sim.simData.Network)
	sim.sharedData = controlunit.NewSharedData()

	// fill shared data
	for name, calcCfg := range sim.simData.Workload.Data.Calculators {
		sim.sharedData.Calculators[name] = calcCfg.Info()
	}
	for name, datasetCfg := range sim.simData.Workload.Data.Datasets {
		sim.sharedData.Datasets[name] = datasetCfg.Info()
	}

	for _, agentCfg := range sim.simData.Network.Agents {
		switch agentCfg.Type {
		case networkcfg.AgentControlUnit:
			sim.cuNames = append(sim.cuNames, agentCfg.Name)
		case networkcfg.AgentWorker:
			sim.cuWorkers[agentCfg.ControlUnitName] = append(sim.cuWorkers[agentCfg.ControlUnitName], agentCfg.Name)
		}
	}

	for _, agentCfg := range sim.simData.Network.Agents {
		sim.addAgent(agentCfg)
	}

	return sim
}

// Run - запуск симуляции
func (sim *Simulator) Run() {
	sim.simenv.Run()
}

func (sim *Simulator) addAgent(cfg *networkcfg.AgentCfg) {
	var agent simenv.Agent

	switch cfg.Type {
	case networkcfg.AgentControlUnit:
		agent = controlunit.New(*cfg, sim.cuWorkers[cfg.Name], sim.sharedData, sim.genScheduler, sim.logContext)
	case networkcfg.AgentCore:
		agent = core.New(*cfg, sim.simData.Workload.Requests, sim.cuNames, sim.logContext)
	case networkcfg.AgentWorker:
		agent = worker.New(*cfg, sim.logContext)
	}

	sim.simenv.Add(agent)
}

func (sim *Simulator) genScheduler(logContext *logrus.Entry) scheduler.Scheduler {
	return random.NewScheduler(logContext)
}
