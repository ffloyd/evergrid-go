package random

import (
	"math/rand"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/scheduler"
)

// Scheduler -
type Scheduler struct {
	infoChans    scheduler.InfoChans
	requestChans scheduler.RequestChans
	controlChans scheduler.ControlChans

	log *logrus.Entry
}

// NewScheduler -
func NewScheduler(logContext *logrus.Entry) *Scheduler {
	return &Scheduler{
		infoChans:    scheduler.NewInfoChans(),
		requestChans: scheduler.NewRequestChans(),
		controlChans: scheduler.NewControlChans(),

		log: logContext,
	}
}

// Name -
func (s *Scheduler) Name() string {
	return "Random Scheduler"
}

// Run -
func (s *Scheduler) Run() {
	go s.work()
}

// RequestChans -
func (s *Scheduler) RequestChans() scheduler.RequestChans {
	return s.requestChans
}

// ControlChans -
func (s *Scheduler) ControlChans() scheduler.ControlChans {
	return s.controlChans
}

// InfoChans -
func (s *Scheduler) InfoChans() scheduler.InfoChans {
	return s.infoChans
}

func (s *Scheduler) work() {
	chans := s.requestChans
	for {
		select {
		case request := <-chans.UploadDataset:
			if s.leadershipStatus() {
				s.processUploadDataset(request)
				chans.DelegateToLeader <- false
			} else {
				chans.DelegateToLeader <- true
			}

		case request := <-chans.RunExperiment:
			if s.leadershipStatus() {
				s.processRunExperiment(request)
				chans.DelegateToLeader <- false
			} else {
				chans.DelegateToLeader <- true
			}

		}
	}
}

func (s *Scheduler) leadershipStatus() bool {
	req := scheduler.NewGetLeadershipStatus()
	s.infoChans.LeadershipStatus <- req
	return <-req.Result
}

func (s *Scheduler) getRandomWorker() types.WorkerInfo {
	getNames := scheduler.NewGetWorkerNames()
	s.infoChans.WorkerNames <- getNames
	names := <-getNames.Result

	randomName := names[rand.Intn(len(names))]
	getWorker := scheduler.NewGetWorkerInfo(randomName)
	s.infoChans.WorkerInfo <- getWorker
	return *<-getWorker.Result
}

func (s *Scheduler) getDatasetInfo(datasetName string) types.DatasetInfo {
	getDataset := scheduler.NewGetDatasetInfo(datasetName)
	s.infoChans.DatasetInfo <- getDataset
	return *<-getDataset.Result
}

func (s *Scheduler) processUploadDataset(request scheduler.ReqUploadDataset) {
	worker := s.getRandomWorker()
	s.controlChans.UploadDataset <- scheduler.DoUploadDataset{
		Worker:  worker.UID,
		Dataset: request.Dataset.UID,
	}
	<-s.controlChans.Done

	s.log.WithFields(logrus.Fields{
		"dataset": request.Dataset.UID,
		"worker":  worker.UID,
	}).Info("Dataset uploading scheduled")
}

func (s *Scheduler) processRunExperiment(request scheduler.ReqRunExperiment) {
	dataset := s.getDatasetInfo(request.Dataset.UID)
	worker := append(dataset.Workers, dataset.EnqueuedOnWorkers...)[0]

	s.controlChans.BuildCalculator <- scheduler.DoBuildCalculator{
		Calculator: request.Calculator.UID,
		Worker:     worker,
	}
	<-s.controlChans.Done

	s.controlChans.RunCalculator <- scheduler.DoRunCalculator{
		Calculator: request.Calculator.UID,
		Dataset:    request.Dataset.UID,
		Worker:     worker,
	}
	<-s.controlChans.Done

	s.log.WithFields(logrus.Fields{
		"dataset":    request.Dataset.UID,
		"calculator": request.Calculator.UID,
		"worker":     worker,
	}).Info("Experiment run scheduled")
}
