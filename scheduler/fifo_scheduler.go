package scheduler

import log "github.com/Sirupsen/logrus"

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

	request.Response.Done <- RespDone{}
}
