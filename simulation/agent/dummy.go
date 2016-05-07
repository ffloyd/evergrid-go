package agent

import log "github.com/Sirupsen/logrus"

// Dummy is a simple agent which only writes currentTick in log
type Dummy struct {
	Agent
}

// NewDummy creates new dummy agent
func NewDummy(name string) *Dummy {
	dummy := new(Dummy)
	dummy.name = name
	dummy.chans = NewChans()

	log.WithField("name", dummy).Info("Dummy agent initialized")
	return dummy
}

func (agent Dummy) run() {
	for {
		log.WithFields(log.Fields{
			"tick":  <-agent.chans.Ticks,
			"agent": agent,
		}).Debug("received tick")
		agent.chans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (agent Dummy) Run() *Chans {
	go agent.run()
	return agent.chans
}
