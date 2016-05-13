package workloadcfg

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/yaml.v2"

	"github.com/ffloyd/evergrid-go/simulation/simdata/workloadcfg/datacfg"
)

// YAML is a representation of workload yaml file
type YAML struct {
	Name     string
	Data     string
	Requests map[int][]RequestCfgYAML
}

// WorkloadCfg is a representation of workload config
type WorkloadCfg struct {
	Name     string
	Data     *datacfg.DataCfg
	Requests map[int][]*RequestCfg
}

// Load parses workload config
func Load(workloadFilename string) *WorkloadCfg {
	rawYAML, e := ioutil.ReadFile(workloadFilename)
	if e != nil {
		log.Fatalf("File error: %v", e)
	}

	configYAML := new(YAML)
	yaml.Unmarshal(rawYAML, configYAML)

	workloadCfg := &WorkloadCfg{
		Name:     configYAML.Name,
		Requests: make(map[int][]*RequestCfg),
	}

	absPath, e := filepath.Abs(workloadFilename)
	if e != nil {
		log.Fatalf("Filepath error: %v", e)
	}
	basePath := filepath.Dir(absPath)

	dataFilename := filepath.Join(basePath, "data", configYAML.Data)
	workloadCfg.Data = datacfg.Load(dataFilename)

	for tick, requestsYAML := range configYAML.Requests {
		workloadCfg.Requests[tick] = make([]*RequestCfg, len(requestsYAML))
		for i, reqYAML := range requestsYAML {
			workloadCfg.Requests[tick][i] = reqYAML.Parse(workloadCfg.Data)
		}
	}

	log.WithFields(log.Fields{
		"file": absPath,
		"name": workloadCfg.Name,
	}).Info("Workload config parsed and loaded")

	return workloadCfg
}
