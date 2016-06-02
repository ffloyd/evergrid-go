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
			amILeader := <-chans.Sensors.IsLeader
			if amILeader {
				log.WithFields(log.Fields{
					"ID": sched.base.ID,
				}).Info("FIFO scheduler: processing upload_dataset request")
				sched.processRequest(<-<-chans.Sensors.GlobalState)
				udReqest.Response.Done <- RespDone{}
			} else {
				log.WithFields(log.Fields{
					"ID": sched.base.ID,
				}).Info("FIFO scheduler: redirect upload dataset request to leader")
				udReqest.Response.DelegateToLeader <- RespDelegateToLeader{}
			}
		}
	}
}

func (sched *fifoScheduler) processRequest(state *types.GlobalState) {
}
