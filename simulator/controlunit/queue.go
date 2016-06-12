package controlunit

import "github.com/ffloyd/evergrid-go/scheduler"

type localQueue struct {
	cu *ControlUnit

	queues map[string][]interface{}
}

func newLocalQueue(cu *ControlUnit) localQueue {
	return localQueue{
		cu:     cu,
		queues: make(map[string][]interface{}),
	}
}

func (lq *localQueue) Push(task interface{}) {
	var worker string
	switch value := task.(type) {
	case scheduler.DoUploadDataset:
		worker = value.Worker
	case scheduler.DoBuildCalculator:
		worker = value.Worker
	case scheduler.DoRunCalculator:
		worker = value.Worker
	default:
		lq.cu.log.Panicf("Invalid task type: %v", value)
	}

	lq.queues[worker] = append(lq.queues[worker], task)
}
