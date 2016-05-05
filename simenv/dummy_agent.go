package simenv

import log "github.com/Sirupsen/logrus"

// DummyAgent is a simple agent which only writes in log currentTick
type DummyAgent struct {
	Name string
}

func (agent DummyAgent) run(chans *AgentChans) {
	for {
		log.WithFields(log.Fields{
			"tick":  <-chans.ticks,
			"agent": agent.Name,
		}).Debug("received tick")
		chans.ready <- true
	}
}

// AgentRun is implementation of AgentRunner iface
func (agent DummyAgent) AgentRun() *AgentChans {
	chans := NewAgentChans()

	go agent.run(chans)
	return chans
}
