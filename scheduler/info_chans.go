package scheduler

import "github.com/ffloyd/evergrid-go/global/types"

// InfoChans -
type InfoChans struct {
	WorkerNames      chan GetWorkerNames
	WorkerInfo       chan GetWorkerInfo
	DatasetInfo      chan GetDatasetInfo
	CalculatorInfo   chan GetCalculatorInfo
	LeadershipStatus chan GetLeadershipStatus
}

// NewInfoChans -
func NewInfoChans() InfoChans {
	return InfoChans{
		WorkerNames:    make(chan GetWorkerNames),
		WorkerInfo:     make(chan GetWorkerInfo),
		DatasetInfo:    make(chan GetDatasetInfo),
		CalculatorInfo: make(chan GetCalculatorInfo),

		LeadershipStatus: make(chan GetLeadershipStatus),
	}
}

// GetWorkerNames -
type GetWorkerNames struct {
	Result chan []string
}

// GetWorkerInfo -
type GetWorkerInfo struct {
	WorkerUID string
	Result    chan *types.WorkerInfo
}

// GetDatasetInfo -
type GetDatasetInfo struct {
	DatasetUID string
	Result     chan *types.DatasetInfo
}

// GetCalculatorInfo -
type GetCalculatorInfo struct {
	CalculatorUID string
	Result        chan *types.CalculatorInfo
}

// GetLeadershipStatus -
type GetLeadershipStatus struct {
	Result chan bool
}

//
// Initializers
//

// NewGetWorkerNames -
func NewGetWorkerNames() GetWorkerNames {
	return GetWorkerNames{
		Result: make(chan []string),
	}
}

// NewGetWorkerInfo -
func NewGetWorkerInfo(uid string) GetWorkerInfo {
	return GetWorkerInfo{
		WorkerUID: uid,
		Result:    make(chan *types.WorkerInfo),
	}
}

// NewGetDatasetInfo -
func NewGetDatasetInfo(uid string) GetDatasetInfo {
	return GetDatasetInfo{
		DatasetUID: uid,
		Result:     make(chan *types.DatasetInfo),
	}
}

// NewGetCalculatorInfo -
func NewGetCalculatorInfo(uid string) GetCalculatorInfo {
	return GetCalculatorInfo{
		CalculatorUID: uid,
		Result:        make(chan *types.CalculatorInfo),
	}
}

// NewGetLeadershipStatus -
func NewGetLeadershipStatus() GetLeadershipStatus {
	return GetLeadershipStatus{
		Result: make(chan bool),
	}
}
