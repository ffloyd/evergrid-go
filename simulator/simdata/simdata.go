package simdata

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulator/simdata/workloadcfg"
	"gopkg.in/yaml.v2"
)

// YAML is a representation of simdata yaml file
type YAML struct {
	Name     string
	Network  string
	Workload string
}

// SimData represents all simulation data config needed to experiment
type SimData struct {
	Name     string
	Network  *networkcfg.NetworkCfg
	Workload *workloadcfg.WorkloadCfg
}

// Load for loading config from YAML to simdata.Config
func Load(simdataFilename string) *SimData {
	rawYAML, e := ioutil.ReadFile(simdataFilename)
	if e != nil {
		log.Fatalf("File error: %v", e)
	}

	configYAML := new(YAML)
	yaml.Unmarshal(rawYAML, configYAML)

	simdata := &SimData{
		Name: configYAML.Name,
	}

	absPath, e := filepath.Abs(simdataFilename)
	if e != nil {
		log.Fatalf("Filepath error: %v", e)
	}
	basePath := filepath.Dir(absPath)

	networkFilename := filepath.Join(basePath, "networks", configYAML.Network)
	simdata.Network = networkcfg.Load(networkFilename)

	workloadFilename := filepath.Join(basePath, "workloads", configYAML.Workload)
	simdata.Workload = workloadcfg.Load(workloadFilename)

	log.WithFields(log.Fields{
		"file": absPath,
		"name": simdata.Name,
	}).Info("Simulation config parsed and loaded")

	return simdata
}
