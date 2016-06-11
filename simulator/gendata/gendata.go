package gendata

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulator/simdata"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulator/simdata/workloadcfg"
	"github.com/ffloyd/evergrid-go/simulator/simdata/workloadcfg/datacfg"
)

type genDataState struct {
	// directories
	baseDir      string
	networksDir  string
	workloadsDir string
	dataDir      string

	// filesnames (without path and extension)
	baseFilename     string
	networkFilename  string
	workloadFilename string
	dataFilename     string

	// structs to export
	dataYAML     datacfg.YAML
	workloadYAML workloadcfg.YAML
	networkYAML  networkcfg.YAML
	simdataYAML  simdata.YAML
}

// GenData generates simdata files using given config
func GenData(config Config) {
	state := genDataState{}
	state.prepareDirs(config)
	state.genFilenames(config)
	state.genData(config)
	state.genWorkload(config)
	state.genNetwork(config)
	state.genBase(config)
	state.save()
}

func (state *genDataState) prepareDirs(config Config) {
	var err error

	state.baseDir, err = filepath.Abs(config.DestDir)
	if err != nil {
		log.Panicf("Filepath error: %s", err)
	}

	state.networksDir = state.baseDir + string(filepath.Separator) + "networks"
	state.workloadsDir = state.baseDir + string(filepath.Separator) + "workloads"
	state.dataDir = state.workloadsDir + string(filepath.Separator) + "data"

	initDir := func(dirName string) {
		log.Infof("Initialize directory %s", dirName)
		os.MkdirAll(dirName, 0777)
	}
	initDir(state.baseDir)
	initDir(state.networksDir)
	initDir(state.workloadsDir)
	initDir(state.dataDir)
}

func (state *genDataState) genFilenames(config Config) {
	state.baseFilename = config.Name + ".yaml"
	state.networkFilename = config.Name + ".yaml"
	state.workloadFilename = config.Name + ".yaml"
	state.dataFilename = config.Name + ".yaml"
}

func (state *genDataState) genData(config Config) {
	log.WithFields(log.Fields{
		"processors":          config.ProcessorsCount,
		"datasets":            config.DatsetsCount,
		"min_dataset_size":    config.MinDatasetSize,
		"max_dataset_size":    config.MaxDatasetSize,
		"min_processor_speed": config.MinSpeed,
		"max_processor_speed": config.MaxSpeed,
	}).Info("Generating data...")

	data := &state.dataYAML

	data.Name = state.dataFilename
	data.Datasets = make([]datacfg.DatasetCfgYAML, config.DatsetsCount)
	data.Processors = make([]datacfg.ProcessorCfgYAML, config.ProcessorsCount)

	for i := 0; i < config.DatsetsCount; i++ {
		data.Datasets[i] = datacfg.DatasetCfgYAML{
			Name: fmt.Sprintf("Dataset %d", i),
			Size: uniformDistr(config.MinDatasetSize, config.MaxDatasetSize),
		}
	}

	for i := 0; i < config.ProcessorsCount; i++ {
		data.Processors[i] = datacfg.ProcessorCfgYAML{
			Name:        fmt.Sprintf("Processor %d", i),
			MFlopsPerMb: float64(uniformDistr(config.MinSpeed, config.MaxSpeed)),
		}
	}

	log.Info("Generating data done")
}

func (state *genDataState) genWorkload(config Config) {
	log.WithFields(log.Fields{
		"processor_runs":  config.ProcessorRuns,
		"run_probability": config.RunProbability,
	}).Info("Generating workload...")

	workload := &state.workloadYAML
	workload.Name = state.workloadFilename
	workload.Data = state.dataFilename
	workload.Requests = make(map[int][]workloadcfg.RequestCfgYAML)

	datasetUploaded := make([]bool, config.DatsetsCount)
	processorsUploaded := make([]bool, config.ProcessorsCount)

	addRequest := func(tick int, request workloadcfg.RequestCfgYAML) {
		workload.Requests[tick] = append(workload.Requests[tick], request)
	}

	eventsLeft := config.ProcessorRuns
	tick := 1
	for ; eventsLeft > 0; tick++ {
		if !coin(config.RunProbability) {
			continue
		}

		datasetIndex, processorIndex := rand.Intn(config.DatsetsCount), rand.Intn(config.ProcessorsCount)
		if !datasetUploaded[datasetIndex] {
			addRequest(tick, workloadcfg.RequestCfgYAML{
				Type:    "upload_dataset",
				Dataset: state.dataYAML.Datasets[datasetIndex].Name,
			})
			datasetUploaded[datasetIndex] = true
		}

		addRequest(tick, workloadcfg.RequestCfgYAML{
			Type:      "run_expirement",
			Processor: state.dataYAML.Processors[processorIndex].Name,
			Dataset:   state.dataYAML.Datasets[datasetIndex].Name,
		})
		processorsUploaded[processorIndex] = true

		eventsLeft--
	}

	log.Info("Generating workload done")
}

func (state *genDataState) genNetwork(config Config) {
	log.WithFields(log.Fields{}).Info("Generating network...")

	network := &state.networkYAML
	network.Name = state.networkFilename
	network.Segments = make([]networkcfg.SegmentCfgYAML, config.NetworkSegments)

	for i := 0; i < config.NetworkSegments; i++ {
		network.Segments[i] = *genNetworkSegment(i, config, i == 0)
	}

	log.Info("Generating network done")
}

func genNetworkSegment(index int, config Config, withCore bool) *networkcfg.SegmentCfgYAML {
	nodeCount := uniformDistr(config.MinNodesInSegment, config.MaxNodesInSegment) + 1

	if withCore {
		nodeCount++
	}

	result := &networkcfg.SegmentCfgYAML{
		Name:          fmt.Sprintf("Segment %d", index),
		InnerBandwith: []int{100, 100},
		OuterBandwith: []int{50, 50},
		Nodes:         make([]networkcfg.NodeCfgYAML, nodeCount),
	}

	controlUnit := genControlUnitNode(fmt.Sprintf("CU%d", index), config)
	result.Nodes[0] = *controlUnit

	startFrom := 1
	if withCore {
		result.Nodes[1] = *genCoreNode("X", config)
		startFrom++
	}

	for i := startFrom; i < nodeCount; i++ {
		suffix := fmt.Sprintf("%d-%d", index, i)
		cuName := controlUnit.Agents[0].Name
		result.Nodes[i] = *genWorkerNode(suffix, cuName, config)
	}

	return result
}

func genWorkerNode(suffix string, controlUnit string, config Config) *networkcfg.NodeCfgYAML {
	result := &networkcfg.NodeCfgYAML{
		Name:   fmt.Sprintf("Node %s", suffix),
		Agents: make([]networkcfg.AgentCfgYAML, 1),
	}

	agentYAML := networkcfg.AgentCfgYAML{
		Name:         fmt.Sprintf("Worker %s", suffix),
		Type:         "worker",
		ControlUnit:  controlUnit,
		WorkerDisk:   uniformDistr(config.MinDiskSize, config.MaxDiskSize),
		WorkerMFlops: uniformDistr(config.MinNodeSpeed, config.MaxNodeSpeed),
		PricePerTick: uniformDistrF64(config.MinPricePerTick, config.MaxPricePerTick),
	}
	result.Agents[0] = agentYAML

	return result
}

func genControlUnitNode(suffix string, config Config) *networkcfg.NodeCfgYAML {
	result := &networkcfg.NodeCfgYAML{
		Name:   fmt.Sprintf("Node %s", suffix),
		Agents: make([]networkcfg.AgentCfgYAML, 1),
	}

	agentYAML := networkcfg.AgentCfgYAML{
		Name: fmt.Sprintf("ControlUnit %s", suffix),
		Type: "control_unit",
	}
	result.Agents[0] = agentYAML

	return result
}

func genCoreNode(suffix string, config Config) *networkcfg.NodeCfgYAML {
	result := &networkcfg.NodeCfgYAML{
		Name:   fmt.Sprintf("Node %s", suffix),
		Agents: make([]networkcfg.AgentCfgYAML, 1),
	}

	agentYAML := networkcfg.AgentCfgYAML{
		Name: fmt.Sprintf("Core %s", suffix),
		Type: "core",
	}
	result.Agents[0] = agentYAML

	return result
}

func (state *genDataState) genBase(config Config) {
	log.WithFields(log.Fields{}).Info("Generating simdata main file...")
	state.simdataYAML = simdata.YAML{
		Name:     config.Name,
		Network:  state.networkFilename,
		Workload: state.workloadFilename,
	}
	log.WithFields(log.Fields{}).Info("Simdata main file generated")
}

func (state *genDataState) save() {
	log.Info("Saving result")

	checkErr := func(err error) {
		if err != nil {
			log.Panicf("Error when saving data: %s", err)
		}
	}

	saveYAML := func(yamlStruct interface{}, dir string, name string) {
		buffer, err := yaml.Marshal(yamlStruct)
		checkErr(err)

		dataFilepath := dir + string(filepath.Separator) + name
		err = ioutil.WriteFile(dataFilepath, buffer, 0644)
		checkErr(err)
	}

	saveYAML(state.dataYAML, state.dataDir, state.dataFilename)
	saveYAML(state.workloadYAML, state.workloadsDir, state.workloadFilename)
	saveYAML(state.networkYAML, state.networksDir, state.networkFilename)
	saveYAML(state.simdataYAML, state.baseDir, state.baseFilename)

	log.Info("Saving done")
}
