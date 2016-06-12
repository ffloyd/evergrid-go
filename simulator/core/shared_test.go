package core_test

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv"
)

type ReceiverAgent struct {
	name   string
	fsm    simenv.AgentFSM
	simenv *simenv.SimEnv
	log    *logrus.Entry

	Requests  []interface{}
	writeLock sync.Mutex
}

func NewReceiverAgent(name string, logContex *logrus.Entry) *ReceiverAgent {
	return &ReceiverAgent{
		name: name,
		log:  logContex,
	}
}

func (ra *ReceiverAgent) Name() string {
	return ra.name
}

func (ra *ReceiverAgent) Run(env *simenv.SimEnv) simenv.AgentChans {
	ra.log = ra.log.WithFields(logrus.Fields{
		"agent": ra.Name(),
		"tick":  env.CurrentTick(),
	})

	ra.simenv = env
	ra.fsm = *simenv.NewAgentFSM(ra.log)

	go ra.work()
	return ra.fsm.Chans()
}

func (ra *ReceiverAgent) Send(req interface{}) chan interface{} {
	ra.writeLock.Lock()
	ra.Requests = append(ra.Requests, req)
	ra.writeLock.Unlock()

	response := make(chan interface{})
	go func() {
		response <- simenv.Ok{}
	}()
	return response
}

func (ra *ReceiverAgent) work() {
	ra.fsm.SetStopFlag(true)

	for {
		ra.fsm.ToReady()
		ra.fsm.ToIdle()
		<-ra.fsm.ToDoneChan()
	}
}
