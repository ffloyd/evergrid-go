package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

// State is a worker state struct. Only its methods should change WorkerInfo fields
type State struct {
	info *types.WorkerInfo
}

// NewState creates a new state for worker
func NewState(name string, diskSpace types.MByte, performance types.MFlop) *State {
	return &State{
		info: &types.WorkerInfo{
			UID:            types.UID(name),
			Busy:           false,
			MFlops:         performance,
			TotalDiskSpace: diskSpace,
			FreeDiskSpace:  diskSpace,
			Datasets:       make(map[types.UID]*types.DatasetInfo),
			Processors:     make(map[types.UID]*types.ProcessorInfo),
		},
	}
}

// IsBusy returns true when worker doing something
func (state *State) IsBusy() bool {
	return state.info.Busy
}

// Busy switch state to busy. Panics if already busy
func (state *State) Busy() {
	if state.IsBusy() {
		log.WithFields(log.Fields{
			"agent": state.info.UID,
		}).Info("Worker must be free, but it busy now =(")
	}
	state.info.Busy = true
}

// Free switch state to free and allows worker to receive next job.
func (state *State) Free() {
	state.info.Busy = false
}

// HasDataset returns true if dataset is already uploaded to worker
func (state *State) HasDataset(dataset *types.DatasetInfo) bool {
	return state.info.Datasets[dataset.UID] != nil
}

// HasProcessor returns true if processor is already built on worker
func (state *State) HasProcessor(processor *types.ProcessorInfo) bool {
	return state.info.Processors[processor.UID] != nil
}

// AddDataset adds dataset to worker
func (state *State) AddDataset(dataset *types.DatasetInfo) {
	state.info.Datasets[dataset.UID] = dataset
}

// AddProcessor adds processor to worker
func (state *State) AddProcessor(processor *types.ProcessorInfo) {
	state.info.Processors[processor.UID] = processor
}

// Info returns correct instanse of WorkerInfo
func (state *State) Info() *types.WorkerInfo {
	return state.info
}
