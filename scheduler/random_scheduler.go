package scheduler

import (
	"math/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

type randomScheduler struct {
	base             *Scheduler
	datasetLocations map[types.UID]types.UID
}

func (sched *randomScheduler) run() {
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

func (sched *randomScheduler) processUploadDataset(request *ReqUploadDataset) {
	sensors := sched.base.Chans.Sensors
	if !<-sensors.IsLeader {
		log.WithFields(log.Fields{
			"ID": sched.base.ID,
		}).Info("Random scheduler: redirect upload dataset request to leader")
		request.Response.DelegateToLeader <- RespDelegateToLeader{}
		return
	}

	log.WithFields(log.Fields{
		"ID": sched.base.ID,
	}).Info("Random scheduler: processing upload_dataset request")

	// upload dataset to a random worker
	status := <-<-sensors.GlobalState
	workerUIDs := make([]types.UID, len(status.Workers))
	i := 0
	for _, worker := range status.Workers {
		workerUIDs[i] = worker.UID
		i++
	}

	chosenWorker := workerUIDs[rand.Intn(len(workerUIDs))]
	sched.datasetLocations[types.UID(request.DatasetID)] = chosenWorker

	request.Response.UploadDatasetToWorker <- RespUploadDatasetToWorker{
		Dataset: types.UID(request.DatasetID),
		Worker:  chosenWorker,
	}

	request.Response.Done <- RespDone{}
}

func (sched *randomScheduler) processRunProcessorOnDataset(request *ReqRunProcessorOnDataset) {
	sensors := sched.base.Chans.Sensors
	if !<-sensors.IsLeader {
		log.WithFields(log.Fields{
			"ID": sched.base.ID,
		}).Info("Random scheduler: redirect upload dataset request to leader")
		request.Response.DelegateToLeader <- RespDelegateToLeader{}
		return
	}

	log.WithFields(log.Fields{
		"ID": sched.base.ID,
	}).Info("Random scheduler: processing run_processor request")

	datasetUID := types.UID(request.DatasetID)

	workerUID := sched.datasetLocations[datasetUID]

	request.Response.BuildProcessor <- RespBuildProcessor{
		Processor: types.UID(request.ProcessorID),
		Worker:    workerUID,
	}

	request.Response.RunProcessor <- RespRunProcessor{
		Processor: types.UID(request.ProcessorID),
		Worker:    workerUID,
		Dataset:   datasetUID,
	}

	request.Response.Done <- RespDone{}
}
