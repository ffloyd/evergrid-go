package infrastructure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
)

// ConfigYAML is a representation of YAML infrastructure file
type ConfigYAML struct {
	Name    string
	Network NetworkYAML
}

// Config for infrastructure creation
type Config struct {
	Name    string
	Network *Network
}

// LoadYAML for loading config from YAML
func LoadYAML(configFilename string) *ConfigYAML {
	rawYAML, e := ioutil.ReadFile(configFilename)
	if e != nil {
		log.Fatalf("File error: %v", e)
	}

	configYAML := new(ConfigYAML)
	yaml.Unmarshal(rawYAML, configYAML)

	return configYAML
}

// Parse transform unmarshalled config to internal config representation
// all validations must be performed on this stage
func (configYAML ConfigYAML) Parse() *Config {
	config := &Config{
		Name: configYAML.Name,
	}
	config.Network = configYAML.Network.Parse(config)

	return config
}
