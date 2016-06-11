package core_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/core"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulator/simdata/workloadcfg"
	"github.com/ffloyd/evergrid-go/simulator/simdata/workloadcfg/datacfg"
)

func uploadDatasetRequest(datasetName string, datasetSize int) *workloadcfg.RequestCfg {
	return &workloadcfg.RequestCfg{
		Type: "upload_dataset",
		Dataset: &datacfg.DatasetCfg{
			Name: datasetName,
			Size: datasetSize,
		},
	}
}

func runExperimentRequest(datasetName string, datasetSize int, caluclatorName string, mflopsPerMb float64) *workloadcfg.RequestCfg {
	return &workloadcfg.RequestCfg{
		Type: "run_experiment",
		Dataset: &datacfg.DatasetCfg{
			Name: datasetName,
			Size: datasetSize,
		},
		Calculator: &datacfg.CalculatorCfg{
			Name:        caluclatorName,
			MFlopsPerMb: mflopsPerMb,
		},
	}
}

func TestCore(t *testing.T) {
	workload := map[int][]*workloadcfg.RequestCfg{
		1: []*workloadcfg.RequestCfg{
			uploadDatasetRequest("Dataset 1", 1),
		},
		10: []*workloadcfg.RequestCfg{
			runExperimentRequest("Dataset 1", 1, "Calc 1", 100),
			runExperimentRequest("Dataset 1", 1, "Calc 2", 200),
		},
	}
	coreCfg := networkcfg.AgentCfg{
		Type: networkcfg.AgentCore,
		Name: "Core",
	}

	log := logrus.WithField("ctx", "test_core")

	pseudoCu1 := NewReceiverAgent("CU 1", log)
	pseudoCu2 := NewReceiverAgent("CU 2", log)

	cuNames := []string{
		pseudoCu1.Name(),
		pseudoCu2.Name(),
	}

	core := core.New(coreCfg, workload, cuNames, log)

	sim := simenv.New()
	sim.Add(core, pseudoCu1, pseudoCu2)
	sim.Run()

	if len(pseudoCu1.Requests)+len(pseudoCu2.Requests) != 3 {
		t.Fail()
	}
}
