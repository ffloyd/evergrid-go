package scheduler

import "github.com/ffloyd/evergrid-go/global/types"

// RequestChans -
type RequestChans struct {
	UploadDataset    chan ReqUploadDataset
	RunExperiment    chan ReqRunExperiment
	DelegateToLeader chan bool
}

// NewRequestChans -
func NewRequestChans() RequestChans {
	return RequestChans{
		UploadDataset:    make(chan ReqUploadDataset),
		RunExperiment:    make(chan ReqRunExperiment),
		DelegateToLeader: make(chan bool),
	}
}

// ReqUploadDataset -
type ReqUploadDataset struct {
	Dataset types.DatasetInfo
}

// ReqRunExperiment -
type ReqRunExperiment struct {
	Calculator types.CalculatorInfo
	Dataset    types.DatasetInfo
}
