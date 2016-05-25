package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/scheduler"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulation/simdata/workloadcfg"
)

// ControlUnit is a representation of control unit app
type ControlUnit struct {
	Base

	incomingRequests    chan *workloadcfg.RequestCfg
	requestConfirmation chan bool
	workers             []*Worker

	scheduler *scheduler.Scheduler
	monitor   *Monitor
}

// NewControlUnit creates a new control unit
func NewControlUnit(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *ControlUnit {
	unit := &ControlUnit{
		Base:                *NewBase(config, net, env),
		incomingRequests:    make(chan *workloadcfg.RequestCfg),
		requestConfirmation: make(chan bool),
	}
	env.ControlUnits[unit.Name()] = unit

	log.WithFields(log.Fields{
		"agent": unit.Name(),
		"node":  unit.Node(),
	}).Info("Control Unit agent initialized")
	return unit
}

func (unit *ControlUnit) processRequest(request *workloadcfg.RequestCfg) {
	switch request.Type {
	case "upload_dataset":
		unit.processDataUpload(request)
	default:
		log.Fatalf("Unknown request type: %s", request.Type)
	}
}

func (unit *ControlUnit) processDataUpload(request *workloadcfg.RequestCfg) {
	schedReq := &scheduler.ReqUploadDataset{
		DatasetID: request.Dataset.Name,
		Response:  make(chan *scheduler.RespUploadDataset),
	}

	unit.scheduler.Chans.Requests.UploadDataset <- schedReq
	response := <-schedReq.Response

	if response.DelegateToLeader {
		leader := unit.env.LeaderControlUnit()
		log.WithFields(log.Fields{
			"agent":   unit,
			"dataset": request.Dataset.Name,
			"leader":  leader,
		}).Info("Redirecting upload dataset request to leader")
		leader.incomingRequests <- request
		<-leader.requestConfirmation
	} else {
		log.WithFields(log.Fields{
			"agent":   unit,
			"dataset": request.Dataset.Name,
		}).Info("Upload dataset request processed")
	}
}

func (unit *ControlUnit) startScheduler() {
	log.WithFields(log.Fields{
		"agent":     unit,
		"algorithm": "FIFO",
	}).Info("Starting scheduler on Control Unit")

	unit.scheduler = scheduler.New(scheduler.FIFO, unit.Name())
	unit.monitor = startMonitor(unit.scheduler, unit.env, unit.Name())
	go unit.scheduler.Run()

	<-unit.scheduler.Chans.Alive

	log.WithFields(log.Fields{
		"agent":     unit,
		"algorithm": "FIFO",
	}).Info("Scheduler started on Control Unit")
}

func (unit *ControlUnit) run() {
	unit.startScheduler()

	for {
		unit.sync.toReady()
		unit.sync.toIdle()
		doneCh := unit.sync.toDoneCallback()

	SelectLoop:
		for {
			select {
			case request := <-unit.incomingRequests:
				unit.sync.toWorking()
				log.WithFields(log.Fields{
					"agent": unit,
					"tick":  unit.sync.tick,
					"type":  request.Type,
				}).Info("Control unit received request")
				unit.requestConfirmation <- true
				unit.processRequest(request)
				unit.sync.toIdle()
			case <-doneCh:
				break SelectLoop
			}
		}
	}
}

// Run is implementation of agent.Runner iface
func (unit *ControlUnit) Run() *Synchronizer {
	go unit.run()
	return unit.sync
}
