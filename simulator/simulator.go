// Package simulator отвечает за запуск симуляции архитектуры Evergrid
package simulator

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/scheduler"
	"github.com/ffloyd/evergrid-go/scheduler/naivecheap"
	"github.com/ffloyd/evergrid-go/scheduler/naivefast"
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
	schedulerName string

	simData *simdata.SimData

	network *network.Network
	simenv  *simenv.SimEnv

	logContext *logrus.Entry
	sharedData *controlunit.SharedData

	cuNames   []string
	cuWorkers map[string][]string

	workers []*worker.Worker
}

// New - инициализация симуляции на основе корневого файла сценария.
func New(simdataFilename string, schedulerName string, jsonLog string) *Simulator {
	if jsonLog != "" {
		filename := jsonLog
		logrus.SetFormatter(&logrus.JSONFormatter{})
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	}

	sim := &Simulator{
		schedulerName: schedulerName,
		simData:       simdata.Load(simdataFilename),
		simenv:        simenv.New(),
		cuWorkers:     make(map[string][]string),
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

// Run - запуск симуляции и вывод общей статистики
func (sim *Simulator) Run() {
	sim.simenv.Run()
	worker.StatsReport(sim.workers, sim.logContext)
}

func (sim *Simulator) addAgent(cfg *networkcfg.AgentCfg) {
	var agent simenv.Agent

	switch cfg.Type {
	case networkcfg.AgentControlUnit:
		agent = controlunit.New(*cfg, sim.cuWorkers[cfg.Name], sim.sharedData, sim.genScheduler, sim.logContext)
	case networkcfg.AgentCore:
		agent = core.New(*cfg, sim.simData.Workload.Requests, sim.cuNames, sim.logContext)
	case networkcfg.AgentWorker:
		worker := worker.New(*cfg, sim.logContext)
		sim.workers = append(sim.workers, worker)
		agent = worker
	}

	sim.simenv.Add(agent)
}

func (sim *Simulator) genScheduler(logContext *logrus.Entry) scheduler.Scheduler {
	switch sim.schedulerName {
	case "random":
		return random.NewScheduler(logContext)
	case "naivefast":
		return naivefast.NewScheduler(logContext)
	case "naivecheap":
		return naivecheap.NewScheduler(logContext)
	default:
		panic("Unknown scheduler type")
	}
}
