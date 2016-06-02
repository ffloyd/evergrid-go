package networkcfg

import (
	log "github.com/Sirupsen/logrus"
)

// AgentCfgYAML is a representation of agent section in YAML infrastructure file
type AgentCfgYAML struct {
	Name        string
	Type        string
	ControlUnit string `yaml:"control_unit"`
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

	Node *NodeCfg // parent
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (agentYAML AgentCfgYAML) Parse(parent *NodeCfg) *AgentCfg {
	return &AgentCfg{
		Name:            agentYAML.Name,
		Type:            resolveAgentType(agentYAML.Type),
		ControlUnitName: agentYAML.ControlUnit,
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
