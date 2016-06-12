package core

import (
	"math/rand"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/comm"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulator/simdata/workloadcfg"
)

// Core -
type Core struct {
	name   string
	fsm    simenv.AgentFSM
	simenv *simenv.SimEnv
	log    *logrus.Entry

	requests         map[int][]*workloadcfg.RequestCfg
	controlUnitNames []string
	controlUnits     []simenv.Agent
	currentTick      *simenv.CurrentTick
}

// New -
func New(cfg networkcfg.AgentCfg, requests map[int][]*workloadcfg.RequestCfg, cuNames []string, logContext *logrus.Entry) *Core {
	if cfg.Type != networkcfg.AgentCore {
		logContext.Panic("Wrong agent type in config")
	}
	return &Core{
		name: cfg.Name,
		log:  logContext,

		requests:         requests,
		controlUnitNames: cuNames,
	}
}

// Name -
func (core *Core) Name() string {
	return core.name
}

// Run -
func (core *Core) Run(env *simenv.SimEnv) simenv.AgentChans {
	core.currentTick = env.CurrentTick()

	core.log = core.log.WithFields(logrus.Fields{
		"agent": core.Name(),
		"tick":  core.currentTick,
	})

	core.simenv = env
	core.fsm = *simenv.NewAgentFSM(core.log)

	core.controlUnits = make([]simenv.Agent, len(core.controlUnitNames))
	for i, cuName := range core.controlUnitNames {
		core.controlUnits[i] = env.Find(cuName)
	}

	go core.work()
	return core.fsm.Chans()
}

// Send -
func (core *Core) Send(msg interface{}) chan interface{} {
	core.log.Panic("Core cannot receive requests")
	return nil
}

func (core *Core) work() {
	activeTicksProcessed := 0

	for {
		core.fsm.ToReady()
		core.fsm.ToWorking()
		controlUnit := core.controlUnits[rand.Intn(len(core.controlUnits))]
		tick := core.currentTick.Int()

		for _, request := range core.requests[tick] {
			core.log.WithFields(logrus.Fields{
				"control_unit": controlUnit.Name(),
				"type":         request.Type,
			}).Info("Core sending request to Control Unit")
			<-controlUnit.Send(core.convertRequest(*request))
		}

		if core.requests[tick] != nil {
			activeTicksProcessed++
		}

		core.fsm.SetStopFlag(activeTicksProcessed == len(core.requests))
		core.fsm.ToIdle()
		<-core.fsm.ToDoneChan()
	}
}

func (core *Core) convertRequest(request workloadcfg.RequestCfg) interface{} {
	switch request.Type {
	case "upload_dataset":
		return comm.ControlUnitUploadDataset{
			Dataset: *request.Dataset.Info(),
		}
	case "run_experiment":
		return comm.ControlUnitRunExperiment{
			Calculator: *request.Calculator.Info(),
			Dataset:    *request.Dataset.Info(),
		}
	default:
		core.log.Panicf("Unknown request type: %v", request.Type)
		return nil
	}
}
