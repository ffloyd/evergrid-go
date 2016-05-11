package agent

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
)

// AgentBase is a common part for all types of agents in simulation
type agentBase struct {
	name  string
	chans *Chans
}

// String for implement Stringer interface
func (agent agentBase) String() string {
	return agent.name
}

// Name needed for Agent interface implementation
func (agent agentBase) Name() string {
	return agent.name
}

// Agent interface must be implemented for every agent
type Agent interface {
	fmt.Stringer
	Name() string
	Run() *Chans
}

// New initializes agent from config
func New(config *infrastructure.Agent) Agent {
	var agent Agent
	switch config.Type {
	case "dummy":
		agent = NewDummy(config.Name)
	default:
		log.Fatalf("Unknown agent type %s", config.Type)
	}

	return agent
}
