package worker

import (
	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simulator/comm"
)

type builder struct {
	worker *Worker

	calculator types.CalculatorInfo
	building   bool

	builtCalculators map[string]types.CalculatorInfo

	log *logrus.Entry
}

func newBuilder(w *Worker, logContext *logrus.Entry) builder {
	return builder{
		worker:           w,
		log:              logContext,
		builtCalculators: make(map[string]types.CalculatorInfo),
	}
}

func (b *builder) Calculators() map[string]types.CalculatorInfo {
	result := make(map[string]types.CalculatorInfo)
	for k, v := range b.builtCalculators {
		result[k] = v
	}
	return result
}

func (b *builder) Prepare(request comm.WorkerBuildCalculator) {
	b.calculator = request.Calculator

	_, hasCalculator := b.builtCalculators[b.calculator.UID]
	if hasCalculator {
		b.log.WithField("calculator", b.calculator.UID).Info("Calculator already built on this worker")
		return
	}

	b.worker.busy = true
	b.worker.fsm.SetStopFlag(false)
	b.building = true

	b.log.WithField("calculator", b.calculator.UID).Info("Calculator building initiated")
}

func (b *builder) Process() {
	if !b.building {
		return
	}

	b.worker.stats.BuildingTicks++

	b.building = false
	b.builtCalculators[b.calculator.UID] = b.calculator
	b.worker.busy = false
	b.worker.fsm.SetStopFlag(true)
	b.log.WithField("calculator", b.calculator.UID).Info("Calculator building done")
}
