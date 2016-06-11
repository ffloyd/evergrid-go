package worker_test

import (
	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv"
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
	sa.fsm = simenv.NewAgentFSM(sa.log.WithFields(logrus.Fields{
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

	for {
		sa.fsm.ToReady()
		sa.fsm.ToWorking()
		if reqIndex < len(sa.requests) {
			<-sa.simenv.Find(sa.workerName).Send(sa.requests[reqIndex])
			reqIndex++
		} else {
			sa.fsm.SetStopFlag(true)
		}
		sa.fsm.ToIdle()
		<-sa.fsm.ToDoneChan()
	}
}
