package networkcfg

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// YAML is a representation of network section in YAML infrastructure file
type YAML struct {
	Name     string
	Segments []SegmentCfgYAML
}

// NetworkCfg is a struct needed to create new network.Network instance
type NetworkCfg struct {
	Name     string
	Segments []*SegmentCfg

	Nodes  []*NodeCfg  // all nodes inside network
	Agents []*AgentCfg // all agents inside network
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (networkYAML YAML) Parse() *NetworkCfg {
	networkCfg := &NetworkCfg{
		Name: networkYAML.Name,
	}

	segments, nodesCount, agentsCount := make([]*SegmentCfg, len(networkYAML.Segments)), 0, 0
	for i, segmentYAML := range networkYAML.Segments {
		segments[i] = segmentYAML.Parse(networkCfg)
		nodesCount += len(segments[i].Nodes)
		agentsCount += len(segments[i].Agents)
	}

	nodes, agents := make([]*NodeCfg, 0, nodesCount), make([]*AgentCfg, 0, agentsCount)
	for _, segment := range segments {
		nodes, agents = append(nodes, segment.Nodes...), append(agents, segment.Agents...)
	}

	networkCfg.Segments = segments
	networkCfg.Nodes = nodes
	networkCfg.Agents = agents

	return networkCfg
}

// Load parses network configuration from yaml file
func Load(configFilename string) *NetworkCfg {
	rawYAML, e := ioutil.ReadFile(configFilename)
	if e != nil {
		log.Fatalf("File error: %v", e)
	}

	networkYAML := new(YAML)
	yaml.Unmarshal(rawYAML, networkYAML)

	return networkYAML.Parse()
}
