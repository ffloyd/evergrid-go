package scheduler

import (
	log "github.com/Sirupsen/logrus"
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
			amILeader := <-chans.Sensors.IsLeader
			if amILeader {
				log.WithFields(log.Fields{
					"ID": sched.base.ID,
				}).Info("FIFO scheduler: processing upload_dataset request")
				udReqest.Response <- &RespUploadDataset{
					DelegateToLeader: false,
					UploadToWorkers:  make([]string, 0),
				}
			} else {
				log.WithFields(log.Fields{
					"ID": sched.base.ID,
				}).Info("FIFO scheduler: redirect upload dataset request to leader")
				udReqest.Response <- &RespUploadDataset{
					DelegateToLeader: true,
				}
			}
		}
	}
}
