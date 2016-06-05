package scheduler

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

type fifoScheduler struct {
	base *Scheduler
}

func (sched *fifoScheduler) run() {
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

func (sched *fifoScheduler) processUploadDataset(request *ReqUploadDataset) {
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

	var firstWorker *types.WorkerInfo
	for _, worker := range status.Workers {
		firstWorker = worker
		break
	}

	request.Response.UploadDatasetToWorker <- RespUploadDatasetToWorker{
		Dataset: types.UID(request.DatasetID),
		Worker:  firstWorker.UID,
	}

	request.Response.Done <- RespDone{}
}

func (sched *fifoScheduler) processRunProcessorOnDataset(request *ReqRunProcessorOnDataset) {
	log.WithFields(log.Fields{
		"ID": sched.base.ID,
	}).Info("FIFO scheduler: processing run_processor request")
	request.Response.Done <- RespDone{}
}
