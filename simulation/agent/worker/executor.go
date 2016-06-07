package worker

import (
	"fmt"
	"math"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

// Executor is an internal component for worker which manages processors execution
type Executor struct {
	state *State

	processor  *types.ProcessorInfo
	dataset    *types.DatasetInfo
	mflopsDone types.MFlop
	executing  bool
}

// NewExecutor creates a new builder instance
func NewExecutor(state *State) *Executor {
	return &Executor{
		state: state,
	}
}

// Prepare initate execution process
func (executor *Executor) Prepare(request ReqExecute) {
	executor.state.Busy()
	executor.dataset, executor.processor = request.Dataset, request.Processor
	executor.mflopsDone = 0
	executor.executing = true

	if !executor.state.HasProcessor(executor.processor) {
		log.Panic("Processor must be build on worker before execution")
	}

	log.WithFields(log.Fields{
		"agent":     executor.state.info.UID,
		"processor": executor.processor.UID,
	}).Info("Initiate processor run")
}

// Process performs processor execution on worker
func (executor *Executor) Process() {
	if !executor.executing {
		return
	}

	// 1 tick == 1 minute
	executor.mflopsDone += executor.state.info.MFlops * 60

	currentMFlops := executor.mflopsDone
	totalMFlops := executor.processor.MFlopsPerMb * types.MFlop(executor.dataset.Size)

	if currentMFlops >= totalMFlops {
		executor.executing = false
		executor.state.Free()

		log.WithFields(log.Fields{
			"agent":     executor.state.info.UID,
			"processor": executor.processor.UID,
		}).Info("Processor executed")
	} else {
		progress := math.Min(1.0, float64(currentMFlops)/float64(totalMFlops))

		log.WithFields(log.Fields{
			"agent":     executor.state.info.UID,
			"processor": executor.processor.UID,
			"progress":  fmt.Sprintf("%d%%", int(progress*100)),
		}).Info("Processing...")
	}
}
