package worker

import (
	"fmt"
	"math"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simulator/comm"
)

type executor struct {
	worker *Worker

	executing  bool
	calculator types.CalculatorInfo
	dataset    types.DatasetInfo
	mflopsDone types.MFlop

	log *logrus.Entry
}

func newExecutor(w *Worker, logContext *logrus.Entry) executor {
	return executor{
		worker: w,
		log:    logContext,
	}
}

func (ex *executor) Prepare(request comm.WorkerRunCalculator) {
	calculator, hasCalculator := ex.worker.builder.builtCalculators[request.Calculator]
	dataset, hasDataset := ex.worker.uploader.uploadedDatasets[request.Dataset]

	if !hasCalculator {
		ex.log.WithField("calculator", calculator.UID).Panic("Worker hasn't calculator to run")
	}

	if !hasDataset {
		ex.log.WithField("dataset", dataset.UID).Panic("Worker hasn't dataset to process")
	}

	ex.executing = true
	ex.calculator, ex.dataset = calculator, dataset
	ex.worker.busy = true
	ex.worker.fsm.SetStopFlag(false)

	ex.log.WithFields(logrus.Fields{
		"calculator": ex.calculator.UID,
		"dataset":    ex.dataset.UID,
	}).Info("Initiate calculator execution")
}

func (ex *executor) Process() {
	if !ex.executing {
		return
	}

	ex.worker.stats.CalculatingTicks++

	ex.mflopsDone += ex.worker.performance * 60
	totalMFlops := ex.calculator.MFlopsPerMb * types.MFlop(ex.dataset.Size)

	if ex.mflopsDone >= totalMFlops {
		ex.executing = false
		ex.worker.busy = false
		ex.worker.fsm.SetStopFlag(true)

		ex.log.WithFields(logrus.Fields{
			"calculator": ex.calculator.UID,
			"dataset":    ex.dataset.UID,
		}).Info("Calculation done")
	} else {
		progress := math.Min(1.0, float64(ex.mflopsDone)/float64(totalMFlops))

		ex.log.WithFields(logrus.Fields{
			"calculator": ex.calculator.UID,
			"dataset":    ex.dataset.UID,
			"progress":   fmt.Sprintf("%d%%", int(progress*100)),
		}).Info("Calculating...")
	}
}
