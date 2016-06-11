package networkcfg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

// AgentCfgYAML is a representation of agent section in YAML infrastructure file
type AgentCfgYAML struct {
	Name         string
	Type         string
	ControlUnit  string  `yaml:"control_unit,omitempty"`
	WorkerDisk   int     `yaml:"worker_disk,omitempty"`
	WorkerMFlops int     `yaml:"worker_mflops,omitempty"`
	PricePerTick float64 `yaml:"price_per_tick,omitempty"`
}

// AgentType is an enum for agent types
type AgentType int

// AgentType is an enum for agent types
const (
	AgentCore AgentType = iota
	AgentControlUnit
	AgentWorker
	AgentDummy
)

// AgentCfg is a struct needed to create new agent.Agent instance
type AgentCfg struct {
	Name            string
	Type            AgentType
	ControlUnitName string
	WorkerDisk      types.MByte
	WorkerMFlops    types.MFlop
	PricePerTick    float64

	Node *NodeCfg // parent
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (agentYAML AgentCfgYAML) Parse(parent *NodeCfg) *AgentCfg {
	return &AgentCfg{
		Name:            agentYAML.Name,
		Type:            resolveAgentType(agentYAML.Type),
		ControlUnitName: agentYAML.ControlUnit,
		WorkerDisk:      types.MByte(agentYAML.WorkerDisk * 1024),
		WorkerMFlops:    types.MFlop(agentYAML.WorkerMFlops),
		PricePerTick:    agentYAML.PricePerTick,
		Node:            parent,
	}
}

func resolveAgentType(name string) AgentType {
	switch name {
	case "core":
		return AgentCore
	case "control_unit":
		return AgentControlUnit
	case "worker":
		return AgentWorker
	case "dummy":
		return AgentDummy
	default:
		log.Panicf("Unknown agent type: %s", name)
		return -1
	}
}
