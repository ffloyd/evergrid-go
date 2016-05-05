package simenv

import (
	log "github.com/Sirupsen/logrus"
)

// AgentRunner must be implemented for interaction with GlobalTimer
type AgentRunner interface {
	AgentRun() *AgentChans
}

// GlobalTimer is a tool for sync multiple simulation agents (like NetLogo's one)
type GlobalTimer struct {
	currentTick int
	agents      []*AgentChans
}

// Init initializes GlobalTimer
func (gt *GlobalTimer) Init() {
	// noting to do right now
	log.Info("Simulation initialized")
}

// AddAgent adds an agent to watch
func (gt *GlobalTimer) AddAgent(agent AgentRunner) {
	gt.agents = append(gt.agents, agent.AgentRun())
}

// Run simulation
func (gt *GlobalTimer) Run() {
	for {
		// send new tick to all agents
		gt.currentTick++
		log.WithField("tick", gt.currentTick).Debug("new tick")
		for _, agentChans := range gt.agents {
			agentChans.ticks <- gt.currentTick
		}

		// wait for ready status from all agents
		for _, agentChans := range gt.agents {
			_ = <-agentChans.ready
		}
	}
}
