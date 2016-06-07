package scheduler

import (
	"math"
	"math/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

type naiveCheapFifoScheduler struct {
	base *Scheduler
}

func (sched *naiveCheapFifoScheduler) run() {
	chans := sched.base.Chans
	for {
		select {
		case chans.Alive <- true:
		case udReqest := <-chans.Requests.UploadDataset:
			sched.processUploadDataset(udReqest)
		case rdRequest := <-chans.Requests.RunProcessorOnDataset:
			sched.processRunProcessorOnDataset(rdRequest)
		}
	}
}

func (sched *naiveCheapFifoScheduler) processUploadDataset(request *ReqUploadDataset) {
	sensors := sched.base.Chans.Sensors
	if !<-sensors.IsLeader {
		log.WithFields(log.Fields{
			"ID": sched.base.ID,
		}).Info("FIFO scheduler: redirect upload dataset request to leader")
		request.Response.DelegateToLeader <- RespDelegateToLeader{}
		return
	}

	log.WithFields(log.Fields{
		"ID": sched.base.ID,
	}).Info("FIFO scheduler: processing upload_dataset request")

	status := <-<-sensors.GlobalState

	workers := make([]*types.WorkerInfo, len(status.Workers))
	i := 0
	for _, worker := range status.Workers {
		workers[i] = worker
		i++
	}

	randomWorker := workers[rand.Intn(len(workers))]

	request.Response.UploadDatasetToWorker <- RespUploadDatasetToWorker{
		Dataset: types.UID(request.DatasetID),
		Worker:  randomWorker.UID,
	}

	request.Response.Done <- RespDone{}
}

func (sched *naiveCheapFifoScheduler) processRunProcessorOnDataset(request *ReqRunProcessorOnDataset) {
	sensors := sched.base.Chans.Sensors
	if !<-sensors.IsLeader {
		log.WithFields(log.Fields{
			"ID": sched.base.ID,
		}).Info("FIFO scheduler: redirect upload dataset request to leader")
		request.Response.DelegateToLeader <- RespDelegateToLeader{}
		return
	}

	log.WithFields(log.Fields{
		"ID": sched.base.ID,
	}).Info("FIFO scheduler: processing run_processor request")

	status := <-<-sensors.GlobalState

	workers := make([]*types.WorkerInfo, len(status.Workers))
	i := 0
	for _, worker := range status.Workers {
		workers[i] = worker
		i++
	}

	minQueueLen := math.MaxInt32
	minPrice := 1e10
	var chosenWorker *types.WorkerInfo
	for _, pretendent := range workers {
		if pretendent.QueueLength < minQueueLen {
			minQueueLen = pretendent.QueueLength
			chosenWorker = pretendent
			continue
		}

		if (pretendent.QueueLength == minQueueLen) && (pretendent.PricePerTick < minPrice) {
			minPrice = pretendent.PricePerTick
			chosenWorker = pretendent
		}
	}

	request.Response.UploadDatasetToWorker <- RespUploadDatasetToWorker{
		Dataset: types.UID(request.DatasetID),
		Worker:  chosenWorker.UID,
	}

	request.Response.BuildProcessor <- RespBuildProcessor{
		Processor: types.UID(request.ProcessorID),
		Worker:    chosenWorker.UID,
	}

	request.Response.RunProcessor <- RespRunProcessor{
		Processor: types.UID(request.ProcessorID),
		Worker:    chosenWorker.UID,
		Dataset:   types.UID(request.DatasetID),
	}

	request.Response.Done <- RespDone{}
}
