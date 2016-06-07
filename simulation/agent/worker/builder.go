package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

// Builder is an internal component for worker which manages processor builds
type Builder struct {
	state     *State
	processor *types.ProcessorInfo
	progress  int
	building  bool
}

// NewBuilder creates a new builder instance
func NewBuilder(state *State) *Builder {
	return &Builder{
		state: state,
	}
}

// Prepare initate building process
func (builder *Builder) Prepare(request ReqBuild) {
	if builder.state.HasProcessor(request.Processor) {
		log.WithFields(log.Fields{
			"agent":     builder.state.info.UID,
			"processor": request.Processor.UID,
		}).Info("Processor already presents on this worker")
		return
	}

	builder.state.Busy()
	builder.processor = request.Processor
	builder.progress = 0
	builder.building = true

	log.WithFields(log.Fields{
		"agent":     builder.state.info.UID,
		"processor": request.Processor.UID,
	}).Info("Initiate processor build")
}

// Process performs build on worker
func (builder *Builder) Process() {
	if !builder.building {
		return
	}

	// 1 (initialization) tick build for everything for now
	builder.building = false
	builder.state.AddProcessor(builder.processor)
	builder.state.Free()

	log.WithFields(log.Fields{
		"agent":     builder.state.info.UID,
		"processor": builder.processor.UID,
	}).Info("Processor built")
}
