package ticker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv/agent"
)

// Ticker is a global timer like NetLogo's one
type Ticker struct {
	currentTick int
	agents      []*agent.Chans
}

// New creates a new Ticker
func New() *Ticker {
	defer log.Info("New ticker initialized")
	return &Ticker{}
}

// AddAgent adds an agent to watch
func (ticker *Ticker) AddAgent(newAgent agent.Runner) {
	ticker.agents = append(ticker.agents, newAgent.Run())
}

// Run simulation
func (ticker *Ticker) Run() {
	for {
		// send new tick to all agents
		ticker.currentTick++
		log.WithField("tick", ticker.currentTick).Debug("new tick")
		for _, chans := range ticker.agents {
			chans.Ticks <- ticker.currentTick
		}

		// wait for ready status from all agents
		for _, chans := range ticker.agents {
			_ = <-chans.Ready
		}
	}
}
