package worker_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/comm"
	"github.com/ffloyd/evergrid-go/simulator/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulator/worker"
)

func testWorker(requests []interface{}, context string) *worker.Worker {
	config := networkcfg.AgentCfg{
		Name:            "Worker 1",
		Type:            networkcfg.AgentWorker,
		ControlUnitName: "",
		WorkerDisk:      100 * 1024,
		WorkerMFlops:    100,
		PricePerTick:    10.0,
	}

	logContext := logrus.WithField("ctx", context)

	worker := worker.New(config, logContext)

	sender := NewSenderAgent("Sender 1", "Worker 1", requests, logContext)

	env := simenv.New()
	env.Add(worker, sender)
	env.Run()

	return worker
}

func TestUploading(t *testing.T) {
	requests := []interface{}{comm.WorkerUploadDataset{
		Dataset: types.DatasetInfo{
			UID:  "Dataset 1",
			Size: 1000,
		},
	}}

	worker := testWorker(requests, "upload_dataset")
	if worker.Stats().UploadingTicks != 2 {
		t.Fail()
	}
}

func TestBuilding(t *testing.T) {
	requests := []interface{}{
		comm.WorkerBuildCalculator{
			Calculator: types.CalculatorInfo{
				UID:         "Calculator 1",
				MFlopsPerMb: 1000,
			},
		},
	}

	worker := testWorker(requests, "build_calculator")
	if worker.Stats().BuildingTicks != 1 {
		t.Fail()
	}
}
