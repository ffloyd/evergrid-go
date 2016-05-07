package loader

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

// Infrastucture is a representation of json infrastucture file
type Infrastucture struct {
	Name    string
	Network Network
}

// Network is a representation of network section in json infrastructure file
type Network struct {
	Name     string
	Segments []Segment
}

// Segment is a representation of network segment section in json infrastructure file
type Segment struct {
	Name  string
	Nodes []Node
}

// Node is a representation of network node section in json infrastructure file
type Node struct {
	Type   string
	Name   string
	Agents []Agent
}

// Agent is a representation of agent section in json infrastructure file
type Agent struct {
	Name string
	Type string
}

// LoadInfrastructure just parses config file and validates it
func LoadInfrastructure(filename string) *Infrastucture {
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Fatalf("File error: %v", e)
	}

	parsed := new(Infrastucture)
	json.Unmarshal(data, parsed)

	return parsed
}
