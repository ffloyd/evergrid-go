package controlunit

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/scheduler"
	"github.com/ffloyd/evergrid-go/simulator/comm"
)

type localQueue struct {
	cu *ControlUnit

	queues map[string][]interface{}
	mutex  sync.Mutex
}

func newLocalQueue(cu *ControlUnit) *localQueue {
	return &localQueue{
		cu:     cu,
		queues: make(map[string][]interface{}),
	}
}

func (lq *localQueue) Push(task interface{}) {
	lq.mutex.Lock()
	var workerName string
	sd := lq.cu.sharedData
	sd.Mutex.Lock()

	switch value := task.(type) {
	case scheduler.DoUploadDataset:
		workerName = value.Worker
		dataset := sd.Datasets[value.Dataset]
		worker := sd.Workers[workerName]
		dataset.EnqueuedOnWorkers = append(dataset.EnqueuedOnWorkers, workerName)
		worker.QueueLength++
		worker.DatasetsInQueue = append(worker.DatasetsInQueue, dataset.UID)

		lq.cu.log.WithFields(logrus.Fields{
			"worker":  workerName,
			"dataset": dataset.UID,
		}).Info("Enqueued dataset uploading")
	case scheduler.DoBuildCalculator:
		workerName = value.Worker
		calculator := sd.Calculators[value.Calculator]
		worker := sd.Workers[workerName]
		calculator.EnqueuedOnWorkers = append(calculator.EnqueuedOnWorkers, workerName)
		worker.QueueLength++
		worker.CalculatorsInQueue = append(worker.CalculatorsInQueue, calculator.UID)

		lq.cu.log.WithFields(logrus.Fields{
			"worker":     workerName,
			"calculator": calculator.UID,
		}).Info("Enqueued calculator building")
	case scheduler.DoRunCalculator:
		workerName = value.Worker
		lq.cu.log.WithFields(logrus.Fields{
			"worker":     workerName,
			"calculator": value.Calculator,
			"dataset":    value.Dataset,
		}).Info("Enqueued calculator run")
	default:
		lq.cu.log.Panicf("Invalid task type: %v", value)
	}

	sd.Mutex.Unlock()
	lq.queues[workerName] = append(lq.queues[workerName], task)
	lq.mutex.Unlock()
}

func (lq *localQueue) Process() {
	lq.mutex.Lock()
	lq.updateWorkersInfo()

	sd := lq.cu.sharedData
	for workerName := range lq.queues {
		sd.Mutex.Lock()
		if sd.Workers[workerName].Busy {
			sd.Mutex.Unlock()
			continue
		}
		sd.Mutex.Unlock()

		if len(lq.queues[workerName]) == 0 {
			continue
		}

		sd.Mutex.Lock()
		worker := sd.Workers[workerName]
		sd.Mutex.Unlock()

		var request interface{}
		request, lq.queues[workerName] = lq.queues[workerName][0], lq.queues[workerName][1:len(lq.queues[workerName])]

		switch value := request.(type) {
		case scheduler.DoUploadDataset:
			sd.Mutex.Lock()
			dataset := sd.Datasets[value.Dataset]
			dataset.EnqueuedOnWorkers = dataset.EnqueuedOnWorkers[1:len(dataset.EnqueuedOnWorkers)]
			dataset.Workers = append(dataset.Workers, workerName)

			worker.DatasetsInQueue = worker.DatasetsInQueue[1:len(worker.DatasetsInQueue)]
			worker.Datasets = append(worker.Datasets, dataset.UID)
			worker.QueueLength--
			sd.Mutex.Unlock()

			<-lq.cu.workers[workerName].Send(comm.WorkerUploadDataset{
				Dataset: *dataset,
			})
		case scheduler.DoBuildCalculator:
			sd.Mutex.Lock()
			calculator := sd.Calculators[value.Calculator]
			calculator.EnqueuedOnWorkers = calculator.EnqueuedOnWorkers[1:len(calculator.EnqueuedOnWorkers)]
			calculator.Workers = append(calculator.Workers, workerName)

			worker.CalculatorsInQueue = worker.CalculatorsInQueue[1:len(worker.CalculatorsInQueue)]
			worker.Calculators = append(worker.Calculators, calculator.UID)
			worker.QueueLength--
			sd.Mutex.Unlock()

			<-lq.cu.workers[workerName].Send(comm.WorkerBuildCalculator{
				Calculator: *calculator,
			})
		case scheduler.DoRunCalculator:
			<-lq.cu.workers[workerName].Send(comm.WorkerRunCalculator{
				Calculator: value.Calculator,
				Dataset:    value.Dataset,
			})
		}
	}
	lq.mutex.Unlock()
}

func (lq *localQueue) Empty() bool {
	lq.mutex.Lock()
	sum := 0
	for _, arr := range lq.queues {
		sum += len(arr)
	}
	lq.mutex.Unlock()
	return sum == 0
}

func (lq *localQueue) updateWorkersInfo() {
	for name, agent := range lq.cu.workers {
		switch value := (<-agent.Send(comm.WorkerInfo{})).(type) {
		case types.WorkerInfo:
			lq.cu.sharedData.Mutex.Lock()
			old := lq.cu.sharedData.Workers[name]

			if old != nil {
				value.QueueLength = old.QueueLength
				value.DatasetsInQueue = old.DatasetsInQueue
				value.CalculatorsInQueue = old.CalculatorsInQueue
				value.Datasets = old.Datasets
				value.Calculators = old.Calculators
			}

			lq.cu.sharedData.Workers[name] = &value
			lq.cu.sharedData.Mutex.Unlock()
		}
	}
}
