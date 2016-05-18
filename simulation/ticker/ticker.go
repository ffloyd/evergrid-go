package ticker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/agent"
)

// Ticker is a global timer like NetLogo's one
type Ticker struct {
	currentTick int
	agentChans  []*agent.TickerChans
}

// New creates a new Ticker. Also it runs all agents because it essential for correct work of ticker.
func New(agents []agent.Agent) *Ticker {
	defer log.Info("New ticker initialized")

	ticker := new(Ticker)

	ticker.agentChans = make([]*agent.TickerChans, len(agents))
	for i, agent := range agents {
		ticker.agentChans[i] = agent.Run()
	}

	return ticker
}

// Run starts ticker
func (ticker *Ticker) Run() {
	for {
		// send new tick to all agents
		ticker.currentTick++
		log.WithField("tick", ticker.currentTick).Debug("new tick")
		for _, chans := range ticker.agentChans {
			chans.Ticks <- ticker.currentTick
		}

		// tick start sync
		for _, chans := range ticker.agentChans {
			_ = <-chans.Ready
		}

		// tick end sync
		for _, chans := range ticker.agentChans {
			_ = <-chans.Ready
		}
	}
}
