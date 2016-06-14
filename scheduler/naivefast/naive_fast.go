// Package naivefast содержит тривиальную реализацию планировщика с акцентом на быстрое выполнение задач.
package naivefast

import (
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/scheduler"
)

/*
Scheduler - это тривиальная реализация планировщика.

При запросе на загрузку датасета выбираются 3 наиболее производительных воркера с размером очереди меньше пяти,
либо просто три наиболее производительных воркера.

При запросе на выполнение эксперимента - среди воркеров с загруженным (или с запланированной загрузкой) датасетом
выбирается наиболее быстрый с размером очереди меньше 5-и, либо просто наиболее быстрый, если это невозможно.
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

// Name возвращает назавние планировщика: Naive Fast Scheduler
func (s *Scheduler) Name() string {
	return "Naive Fast Scheduler"
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
			if s.infoChans.GetLeadershipStatus() {
				s.processUploadDataset(request)
				chans.DelegateToLeader <- false
			} else {
				chans.DelegateToLeader <- true
			}

		case request := <-chans.RunExperiment:
			if s.infoChans.GetLeadershipStatus() {
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

func (s *Scheduler) getWorkersForUpload() []string {
	workerNames := s.infoChans.GetWorkerNames()
	workers := make([]types.WorkerInfo, len(workerNames))
	for i, workerName := range workerNames {
		workers[i] = s.infoChans.GetWorkerInfo(workerName)
	}

	var workersWithShortQueue []types.WorkerInfo
	for _, worker := range workers {
		if worker.QueueLength < 5 {
			workersWithShortQueue = append(workersWithShortQueue, worker)
		}
	}

	sort.Sort(byMFlopsDesc(workersWithShortQueue))

	var result []string

	if len(workersWithShortQueue) >= 3 {
		for _, worker := range workersWithShortQueue[0:3] {
			result = append(result, worker.UID)
		}
	} else {
		sort.Stable(byMFlopsDesc(workers))
		sort.Stable(byQueueAsc(workers))
		for _, worker := range workers[0:3] {
			result = append(result, worker.UID)
		}
	}

	return result
}

func (s *Scheduler) getWorkerForRun(datasetName string) string {
	datasetInfo := s.infoChans.GetDatasetInfo(datasetName)
	workerNames := append(datasetInfo.Workers, datasetInfo.EnqueuedOnWorkers...)
	workers := make([]types.WorkerInfo, len(workerNames))
	for i, workerName := range workerNames {
		workers[i] = s.infoChans.GetWorkerInfo(workerName)
	}

	var filtered []types.WorkerInfo
	for _, worker := range workers {
		if worker.QueueLength < 5 {
			filtered = append(filtered, worker)
		}
	}

	var result string
	if len(filtered) > 0 {
		sort.Sort(byMFlopsDesc(filtered))
		result = filtered[0].UID
	} else {
		sort.Sort(byMFlopsDesc(workers))
		sort.Stable(byQueueAsc(workers))
		result = workers[0].UID
	}

	return result
}

func (s *Scheduler) processUploadDataset(request scheduler.ReqUploadDataset) {
	for _, worker := range s.getWorkersForUpload() {
		s.controlChans.UploadDataset <- scheduler.DoUploadDataset{
			Worker:  worker,
			Dataset: request.Dataset.UID,
		}
		<-s.controlChans.Done

		s.log.WithFields(logrus.Fields{
			"dataset": request.Dataset.UID,
			"worker":  worker,
		}).Info("Dataset uploading scheduled")
	}
}

func (s *Scheduler) processRunExperiment(request scheduler.ReqRunExperiment) {
	worker := s.getWorkerForRun(request.Dataset.UID)

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
