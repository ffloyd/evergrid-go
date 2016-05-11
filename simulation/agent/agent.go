package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/config/infrastructure"
)

// Agent is a common part for all types of agents in simulation
type Agent struct {
	name  string
	chans *Chans
}

// String for implement Stringer interface
func (agent Agent) String() string {
	return agent.name
}

// New initializes agent from config
func New(config *infrastructure.Agent) Runner {
	var agent Runner
	switch config.Type {
	case "dummy":
		agent = NewDummy(config.Name)
	default:
		log.Fatalf("Unknown agent type %s", config.Type)
	}

	return agent
}
