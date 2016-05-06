package agent

import log "github.com/Sirupsen/logrus"

// Dummy is a simple agent which only writes currentTick in log
type Dummy struct {
	Name string
}

func (agent Dummy) run(chans *Chans) {
	for {
		log.WithFields(log.Fields{
			"tick":  <-chans.Ticks,
			"agent": agent.Name,
		}).Debug("received tick")
		chans.Ready <- true
	}
}

// Run is implementation of agent.Runner iface
func (agent Dummy) Run() *Chans {
	chans := NewChans()

	go agent.run(chans)
	return chans
}
