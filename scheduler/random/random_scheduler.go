// Package random содержит тривиальную реализацию планировщика, который распределеяет нагрузку в случайном порядке
package random

import (
	"math/rand"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/scheduler"
)

/*
Scheduler - это тривиальная реализация планировщика.

При запросе на загрузку датасета выбирается один случайный воркер и датасет загружается на него.

При запросе на выполнение эксперимента - эксперимент запускается на воркере, с уже загруженным датасетом.
*/
type Scheduler struct {
	infoChans    scheduler.InfoChans
	requestChans scheduler.RequestChans
	controlChans scheduler.ControlChans

	log *logrus.Entry
}

// NewScheduler - реализация планировщика, который распределеяет нагрузку в случайном порядке
func NewScheduler(logContext *logrus.Entry) *Scheduler {
	return &Scheduler{
		infoChans:    scheduler.NewInfoChans(),
		requestChans: scheduler.NewRequestChans(),
		controlChans: scheduler.NewControlChans(),

		log: logContext,
	}
}

// Name возвращает назавние планировщика: Random Scheduler
func (s *Scheduler) Name() string {
	return "Random Scheduler"
}

/*
Run запускает планировщик.
*/
func (s *Scheduler) Run() {
	go s.work()
}

// RequestChans - каналы для запросов к планировщику
func (s *Scheduler) RequestChans() scheduler.RequestChans {
	return s.requestChans
}

// ControlChans - каналы управления для планировщика
func (s *Scheduler) ControlChans() scheduler.ControlChans {
	return s.controlChans
}

// InfoChans - каналы мониторинга для планировщика
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
	names := s.InfoChans().GetWorkerNames()

	randomName := names[rand.Intn(len(names))]
	return s.InfoChans().GetWorkerInfo(randomName)
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
	dataset := s.InfoChans().GetDatasetInfo(request.Dataset.UID)
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
