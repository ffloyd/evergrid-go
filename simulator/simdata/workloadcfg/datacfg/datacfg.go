package datacfg

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// YAML is a representation of data yaml file
type YAML struct {
	Name        string
	Datasets    []DatasetCfgYAML
	Calculators []CalculatorCfgYAML
}

// DataCfg is a representation of data config
type DataCfg struct {
	Name        string
	Datasets    map[string]*DatasetCfg
	Calculators map[string]*CalculatorCfg
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
		Name:        configYAML.Name,
		Datasets:    make(map[string]*DatasetCfg),
		Calculators: make(map[string]*CalculatorCfg),
	}

	for _, datasetYAML := range configYAML.Datasets {
		datasetCfg := &DatasetCfg{
			Name: datasetYAML.Name,
			Size: datasetYAML.Size,
		}
		dataCfg.Datasets[datasetCfg.Name] = datasetCfg
	}

	for _, processorYAML := range configYAML.Calculators {
		processorCfg := &CalculatorCfg{
			Name:        processorYAML.Name,
			MFlopsPerMb: processorYAML.MFlopsPerMb,
		}
		dataCfg.Calculators[processorCfg.Name] = processorCfg
	}

	absPath, e := filepath.Abs(dataFilename)
	if e != nil {
		log.Fatalf("Filepath error: %v", e)
	}

	log.WithFields(log.Fields{
		"file": absPath,
		"name": dataCfg.Name,
	}).Info("Data config parsed and loaded")

	return dataCfg
}
