package agent

import (
	"fmt"
	"math"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
)

// Worker is an agent which represents worker machine
type Worker struct {
	Base
	ControlUnit *ControlUnit
	State       types.WorkerInfo

	NewUpload         chan *types.DatasetInfo
	NewProcessorBuild chan *types.ProcessorInfo
	NewProcessorRun   chan WorkerRunProcessorRequest

	currentUpload *jobProcessUpload
	currentBuild  *jobBuildProcessor
	currentRun    *jobRunProcessor
}

type jobProcessUpload struct {
	Dataset  *types.DatasetInfo
	speed    types.MBit
	uploaded types.MByte
}

type jobBuildProcessor struct {
	Processor *types.ProcessorInfo
	progress  int // from 0 to 100
}

type jobRunProcessor struct {
	Processor  *types.ProcessorInfo
	Dataset    *types.DatasetInfo
	mflopsDone types.MFlop
}

// WorkerRunProcessorRequest - request format for one of worker's channels
type WorkerRunProcessorRequest struct {
	Processor *types.ProcessorInfo
	Dataset   *types.DatasetInfo
}

// NewWorker creates new worker agent
func NewWorker(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *Worker {
	worker := &Worker{
		Base: *NewBase(config, net, env),
		State: types.WorkerInfo{
			UID:            types.UID(config.Name),
			Busy:           false,
			MFlops:         config.WorkerMFlops,
			TotalDiskSpace: config.WorkerDisk,
			FreeDiskSpace:  config.WorkerDisk,
			Datasets:       make(map[types.UID]*types.DatasetInfo),
			Processors:     make(map[types.UID]*types.ProcessorInfo),
		},
		NewUpload:         make(chan *types.DatasetInfo),
		NewProcessorBuild: make(chan *types.ProcessorInfo),
		NewProcessorRun:   make(chan WorkerRunProcessorRequest),
	}
	env.Workers[worker.Name()] = worker

	worker.ControlUnit = env.ControlUnits[config.ControlUnitName]
	worker.ControlUnit.workers = append(worker.ControlUnit.workers, worker)

	log.WithFields(log.Fields{
		"agent":        worker.Name(),
		"node":         worker.Node(),
		"control_unit": worker.ControlUnit.Name(),
	}).Info("Worker agent initialized")
	return worker
}

func (worker *Worker) startUpload(dataset *types.DatasetInfo) {
	if worker.State.Busy {
		log.Panicf("Worker '%s' is busy now", worker.Name())
	}

	upload := &jobProcessUpload{
		Dataset:  dataset,
		uploaded: 0,
	}

	datasetUID := upload.Dataset.UID

	// Check if dataset already uploaded
	if worker.State.Datasets[datasetUID] != nil {
		worker.currentUpload = nil
		log.WithFields(log.Fields{
			"agent":   worker.Name(),
			"dataset": datasetUID,
		}).Info("Dataset already presents on this worker")
		return
	}

	// Check if dataset presents in current segment
	internalComm := false
	segmentAgentNames := worker.Node().Segment().AgentNames()
	for _, agentName := range segmentAgentNames {
		closeWorker, ok := worker.env.Workers[agentName]
		if ok {
			if closeWorker.State.Datasets[datasetUID] != nil {
				internalComm = true
				break
			}
		}
	}

	// Reolve upload speed
	bandwith := worker.Node().Segment().Bandwith(internalComm)
	if bandwith.In < bandwith.Out {
		upload.speed = types.MBit(bandwith.In)
	} else {
		upload.speed = types.MBit(bandwith.Out)
	}

	worker.currentUpload = upload
	worker.State.Busy = true
	log.WithFields(log.Fields{
		"agent":   worker.Name(),
		"dataset": upload.Dataset.UID,
	}).Info("Initiate dataset upload")
}

func (worker *Worker) processUpload() {
	upload := worker.currentUpload

	if upload == nil {
		return
	}

	// 1 tick = 1 minute
	mbytesDownloaded := types.MByte(upload.speed * 60 / 8)
	upload.uploaded += mbytesDownloaded

	if upload.uploaded >= upload.Dataset.Size {
		worker.currentUpload = nil
		worker.State.Busy = false

		log.WithFields(log.Fields{
			"agent":   worker.Name(),
			"dataset": upload.Dataset.UID,
		}).Info("Dataset uploaded")
	} else {
		progress := math.Min(1.0, float64(upload.uploaded)/float64(upload.Dataset.Size))

		log.WithFields(log.Fields{
			"agent":    worker.Name(),
			"dataset":  upload.Dataset.UID,
			"progress": fmt.Sprintf("%d%%", int(progress*100)),
		}).Info("Uploading dataset")
	}
}

func (worker *Worker) startBuildProcessor(processor *types.ProcessorInfo) {
	if worker.State.Busy {
		log.Panicf("Worker '%s' is busy now", worker.Name())
	}

	build := &jobBuildProcessor{
		Processor: processor,
		progress:  0,
	}

	processorUID := build.Processor.UID

	// check for "already bult" situation
	if worker.State.Processors[processorUID] != nil {
		worker.currentBuild = nil
		log.WithFields(log.Fields{
			"agent":     worker.Name(),
			"processor": processorUID,
		}).Info("Processor already presents on this worker")
		return
	}

	worker.currentBuild = build
	worker.State.Busy = true
	log.WithFields(log.Fields{
		"agent":     worker.Name(),
		"processor": processorUID,
	}).Info("Initiate processor build")
}

func (worker *Worker) processBuildProcessor() {
	build := worker.currentBuild
	if build == nil {
		return
	}

	// 1 minute build for everything for now
	worker.currentBuild = nil
	worker.State.Busy = false
	worker.State.Processors[build.Processor.UID] = build.Processor

	log.WithFields(log.Fields{
		"agent":     worker.Name(),
		"processor": build.Processor.UID,
	}).Info("Processor built")
}

func (worker *Worker) startRunProcessor(req WorkerRunProcessorRequest) {
	if worker.State.Busy {
		log.Panicf("Worker '%s' is busy now", worker.Name())
	}

	run := &jobRunProcessor{
		Processor:  req.Processor,
		Dataset:    req.Dataset,
		mflopsDone: 0,
	}

	processorUID := run.Processor.UID
	if worker.State.Processors[processorUID] == nil {
		log.Panic("Processor must be build on worker before execution")
	}

	worker.currentRun = run
	worker.State.Busy = true
	log.WithFields(log.Fields{
		"agent":     worker.Name(),
		"processor": processorUID,
	}).Info("Initiate processor run")
}

func (worker *Worker) processRunProcessor() {
	run := worker.currentRun
	if run == nil {
		return
	}

	// 1 tick == 1 minute
	worker.currentRun.mflopsDone += worker.State.MFlops * 60

	currentMFlops := run.mflopsDone
	totalMFlops := run.Processor.MFlopsPerMb * types.MFlop(run.Dataset.Size)

	if currentMFlops >= totalMFlops {
		worker.currentRun = nil
		worker.State.Busy = false

		log.WithFields(log.Fields{
			"agent":     worker.Name(),
			"processor": run.Processor.UID,
		}).Info("Processor executed")
	} else {
		progress := math.Min(1.0, float64(currentMFlops)/float64(totalMFlops))

		log.WithFields(log.Fields{
			"agent":     worker.Name(),
			"processor": run.Processor.UID,
			"progress":  fmt.Sprintf("%d%%", int(progress*100)),
		}).Info("Processing...")
	}
}

func (worker *Worker) run() {
	for {
		worker.sync.toReady()
		worker.sync.toWorking()
		worker.processUpload()
		worker.processBuildProcessor()
		worker.processRunProcessor()
		worker.sync.toIdle()
		doneCh := worker.sync.toDoneCallback()

	SelectLoop:
		for {
			select {
			case dataset := <-worker.NewUpload:
				worker.startUpload(dataset)
			case processor := <-worker.NewProcessorBuild:
				worker.startBuildProcessor(processor)
			case req := <-worker.NewProcessorRun:
				worker.startRunProcessor(req)
			case <-doneCh:
				break SelectLoop
			}
		}
	}
}

// Run is implementation of agent.Runner iface
func (worker Worker) Run() *Synchronizer {
	go worker.run()
	return worker.sync
}
