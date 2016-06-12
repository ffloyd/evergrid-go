package random

import (
	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/scheduler"
)

// Scheduler -
type Scheduler struct {
	infoChans    scheduler.InfoChans
	requestChans scheduler.RequestChans
	controlChans scheduler.ControlChans

	log *logrus.Entry
}

// NewScheduler -
func NewScheduler(logContext *logrus.Entry) *Scheduler {
	return &Scheduler{
		infoChans:    scheduler.NewInfoChans(),
		requestChans: scheduler.NewRequestChans(),
		controlChans: scheduler.NewControlChans(),

		log: logContext,
	}
}

// Name -
func (s *Scheduler) Name() string {
	return "Random Scheduler"
}

// Run -
func (s *Scheduler) Run() {
	go s.work()
}

// RequestChans -
func (s *Scheduler) RequestChans() scheduler.RequestChans {
	return s.requestChans
}

// ControlChans -
func (s *Scheduler) ControlChans() scheduler.ControlChans {
	return s.controlChans
}

// InfoChans -
func (s *Scheduler) InfoChans() scheduler.InfoChans {
	return s.infoChans
}

func (s *Scheduler) work() {
	chans := s.requestChans
	for {
		select {
		case request := <-chans.UploadDataset:
			s.log.Info(request)
			chans.DelegateToLeader <- false
		case request := <-chans.RunExperiment:
			s.log.Info(request)
			chans.DelegateToLeader <- false
		}
	}
}
