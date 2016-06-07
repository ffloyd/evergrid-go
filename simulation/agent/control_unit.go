package agent

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/scheduler"
	workerUtil "github.com/ffloyd/evergrid-go/simulation/agent/worker"
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
	cuQueue   *cuQueue
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
	case "run_expirement":
		unit.processRunExperiment(request)
	default:
		log.Fatalf("Unknown request type: %s", request.Type)
	}
}

func (unit *ControlUnit) processRunExperiment(request *workloadcfg.RequestCfg) {
	schedReq := scheduler.NewReqRunProcessorOnDataset(request.Dataset.Name, request.Processor.Name)

	unit.scheduler.Chans.Requests.RunProcessorOnDataset <- schedReq
	response := schedReq.Response

SelectLoop:
	for {
		select {
		case <-response.Done:
			log.WithFields(log.Fields{
				"agent":     unit,
				"dataset":   request.Dataset.Name,
				"processor": request.Processor.Name,
			}).Info("Run processor request processed")
			break SelectLoop
		case <-response.DelegateToLeader:
			leader := unit.env.LeaderControlUnit()
			log.WithFields(log.Fields{
				"agent":   unit,
				"dataset": request.Dataset.Name,
				"leader":  leader,
			}).Info("Redirecting run processor request to leader")
			leader.incomingRequests <- request
			<-leader.requestConfirmation
			break SelectLoop
		case resp := <-response.UploadDatasetToWorker:
			jobUID := fmt.Sprintf("Upload '%s' to worker '%s'", resp.Dataset, resp.Worker)
			job := types.JobInfo{
				UID:     types.UID(jobUID),
				Type:    types.JobUploadDataset,
				Worker:  resp.Worker,
				Dataset: resp.Dataset,
			}

			queue := unit.env.Workers[string(resp.Worker)].ControlUnit.cuQueue
			queue.forWorker(string(resp.Worker)).push(job)
		case resp := <-response.BuildProcessor:
			jobUID := fmt.Sprintf("Build processor '%s' on worker '%s'", resp.Processor, resp.Worker)
			job := types.JobInfo{
				UID:       types.UID(jobUID),
				Type:      types.JobBuildProcessor,
				Worker:    resp.Worker,
				Processor: resp.Processor,
			}

			queue := unit.env.Workers[string(resp.Worker)].ControlUnit.cuQueue
			queue.forWorker(string(resp.Worker)).push(job)
		case resp := <-response.RunProcessor:
			jobUID := fmt.Sprintf("Build processor '%s' on worker '%s'", resp.Processor, resp.Worker)
			job := types.JobInfo{
				UID:       types.UID(jobUID),
				Type:      types.JobRunProcessor,
				Worker:    resp.Worker,
				Processor: resp.Processor,
				Dataset:   resp.Dataset,
			}

			queue := unit.env.Workers[string(resp.Worker)].ControlUnit.cuQueue
			queue.forWorker(string(resp.Worker)).push(job)
		}
	}
}

func (unit *ControlUnit) processDataUpload(request *workloadcfg.RequestCfg) {
	schedReq := scheduler.NewReqUploadDataset(request.Dataset.Name)

	unit.scheduler.Chans.Requests.UploadDataset <- schedReq
	response := schedReq.Response

SelectLoop:
	for {
		select {
		case <-response.Done:
			log.WithFields(log.Fields{
				"agent":   unit,
				"dataset": request.Dataset.Name,
			}).Info("Upload dataset request processed")
			break SelectLoop
		case <-response.DelegateToLeader:
			leader := unit.env.LeaderControlUnit()
			log.WithFields(log.Fields{
				"agent":   unit,
				"dataset": request.Dataset.Name,
				"leader":  leader,
			}).Info("Redirecting upload dataset request to leader")
			leader.incomingRequests <- request
			<-leader.requestConfirmation
			break SelectLoop
		case resp := <-response.UploadDatasetToWorker:
			log.WithFields(log.Fields{
				"agent":   unit,
				"dataset": request.Dataset.Name,
				"worker":  resp.Worker,
			}).Info("Enqueue uploading to worker")
			jobUID := fmt.Sprintf("Upload '%s' to worker '%s'", resp.Dataset, resp.Worker)
			job := types.JobInfo{
				UID:     types.UID(jobUID),
				Type:    types.JobUploadDataset,
				Worker:  resp.Worker,
				Dataset: resp.Dataset,
			}

			queue := unit.env.Workers[string(resp.Worker)].ControlUnit.cuQueue
			queue.forWorker(string(resp.Worker)).push(job)
		}
	}
}

func (unit *ControlUnit) startScheduler() {
	log.WithFields(log.Fields{
		"agent":     unit,
		"algorithm": "Random",
	}).Info("Starting scheduler on Control Unit")

	unit.scheduler = scheduler.New(scheduler.Random, unit.Name())
	unit.monitor = startMonitor(unit.scheduler, unit.env, unit.Name())
	go unit.scheduler.Run()

	<-unit.scheduler.Chans.Alive

	log.WithFields(log.Fields{
		"agent":     unit,
		"algorithm": "FIFO",
	}).Info("Scheduler started on Control Unit")
}

func (unit *ControlUnit) initQueues() {
	workerNames := make([]string, len(unit.workers))
	for i, worker := range unit.workers {
		workerNames[i] = worker.Name()
	}

	unit.cuQueue = newCUQueue(workerNames)
}

func (unit *ControlUnit) processQueues() {
	for workerName, queue := range unit.cuQueue.workersQueues {
		worker := unit.env.Workers[workerName]
		if worker.State.IsBusy() {
			continue
		}

		nextJob := queue.pop()
		if nextJob == nil {
			continue
		}

		switch nextJob.Type {
		case types.JobUploadDataset:
			worker.UploadChan <- workerUtil.ReqUpload{
				Dataset: unit.env.Datasets[string(nextJob.Dataset)],
			}
		case types.JobBuildProcessor:
			worker.BuildChan <- workerUtil.ReqBuild{
				Processor: unit.env.Processors[string(nextJob.Processor)],
			}
		case types.JobRunProcessor:
			worker.ExecuteChan <- workerUtil.ReqExecute{
				Processor: unit.env.Processors[string(nextJob.Processor)],
				Dataset:   unit.env.Datasets[string(nextJob.Dataset)],
			}
		default:
			log.Panic("Unknown job type")
		}
	}
}

func (unit *ControlUnit) stopFlagUpdate() {
	// stop flag is true if there are no stuff in queues
	queues := unit.cuQueue.workersQueues

	totalTasks := 0
	for _, que := range queues {
		totalTasks += len(que.queue)
	}

	unit.sync.SetStopFlag(totalTasks == 0)
}

func (unit *ControlUnit) run() {
	unit.initQueues()
	unit.startScheduler()

	for {
		unit.sync.toReady()
		unit.sync.toWorking()
		unit.processQueues()
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
				unit.stopFlagUpdate()
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
