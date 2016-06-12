package worker_test

import (
	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/comm"
)

type SenderAgent struct {
	name   string
	fsm    simenv.AgentFSM
	simenv *simenv.SimEnv

	workerName string
	requests   []interface{}

	log *logrus.Entry
}

func NewSenderAgent(name string, worker string, requests []interface{}, logContext *logrus.Entry) *SenderAgent {
	return &SenderAgent{
		name:       name,
		requests:   requests,
		workerName: worker,
		log:        logContext,
	}
}

func (sa *SenderAgent) Name() string {
	return sa.name
}

func (sa *SenderAgent) Run(env *simenv.SimEnv) simenv.AgentChans {
	sa.simenv = env
	sa.fsm = *simenv.NewAgentFSM(sa.log.WithFields(logrus.Fields{
		"agent": sa.Name(),
		"tick":  env.CurrentTick(),
	}))
	go sa.work()
	return sa.fsm.Chans()
}

func (sa *SenderAgent) Send(msg interface{}) chan interface{} {
	panic("No implementation")
}

func (sa *SenderAgent) work() {
	reqIndex := 0
	worker := sa.simenv.Find(sa.workerName)

	for {
		sa.fsm.ToReady()
		sa.fsm.ToWorking()
		if reqIndex < len(sa.requests) {
			busyStatus := <-worker.Send(comm.WorkerBusy{})
			switch value := busyStatus.(type) {
			case bool:
				if !value {
					<-worker.Send(sa.requests[reqIndex])
					reqIndex++
				}
			default:
				sa.log.Panic("Invalid worker response")
			}
		} else {
			sa.fsm.SetStopFlag(true)
		}
		sa.fsm.ToIdle()
		<-sa.fsm.ToDoneChan()
	}
}
