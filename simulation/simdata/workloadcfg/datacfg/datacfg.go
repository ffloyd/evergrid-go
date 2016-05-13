package datacfg

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// YAML is a representation of data yaml file
type YAML struct {
	Datasets   []DatasetCfgYAML
	Processors []ProcessorCfgYAML
}

// DataCfg is a representation of data config
type DataCfg struct {
	Datasets   map[string]*DatasetCfg
	Processors map[string]*ProcessorCfg
}

// Load parses data yaml file
func Load(dataFilename string) *DataCfg {
	rawYAML, e := ioutil.ReadFile(dataFilename)
	if e != nil {
		log.Fatalf("File error: %v", e)
	}

	configYAML := new(YAML)
	yaml.Unmarshal(rawYAML, configYAML)

	dataCfg := &DataCfg{
		Datasets:   make(map[string]*DatasetCfg),
		Processors: make(map[string]*ProcessorCfg),
	}

	for _, datasetYAML := range configYAML.Datasets {
		datasetCfg := &DatasetCfg{datasetYAML}
		dataCfg.Datasets[datasetCfg.Name] = datasetCfg
	}

	for _, processorYAML := range configYAML.Processors {
		processorCfg := &ProcessorCfg{processorYAML}
		dataCfg.Processors[processorCfg.Name] = processorCfg
	}

	return dataCfg
}
