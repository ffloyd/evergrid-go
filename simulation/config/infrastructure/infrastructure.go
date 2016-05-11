package infrastructure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
)

// InfrastuctureYAML is a representation of YAML infrastructure file
type InfrastuctureYAML struct {
	Name    string
	Network NetworkYAML
}

// Infrastucture struct is a config for creation
type Infrastucture struct {
	Name    string
	Network *Network
}

// LoadYAML for loading config from YAML
func LoadYAML(configFilename string) *InfrastuctureYAML {
	data, e := ioutil.ReadFile(configFilename)
	if e != nil {
		log.Fatalf("File error: %v", e)
	}

	result := new(InfrastuctureYAML)
	yaml.Unmarshal(data, result)

	return result
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (yamlData InfrastuctureYAML) Parse() *Infrastucture {
	result := &Infrastucture{
		Name: yamlData.Name,
	}
	result.Network = yamlData.Network.Parse(result)

	return result
}
