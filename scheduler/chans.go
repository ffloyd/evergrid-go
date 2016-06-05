package scheduler

import "github.com/ffloyd/evergrid-go/global/types"

// Chans is a set of chans for communicate with scheduler
type Chans struct {
	Alive    chan bool
	Requests *RequestChans
	Sensors  *SensorChans
}

func newChans() *Chans {
	return &Chans{
		Alive:    make(chan bool),
		Requests: newRequestChans(),
		Sensors:  newSensorChans(),
	}
}

// RequestChans - chans for requests for scheduler
type RequestChans struct {
	UploadDataset         chan *ReqUploadDataset
	RunProcessorOnDataset chan *ReqRunProcessorOnDataset
}

func newRequestChans() *RequestChans {
	return &RequestChans{
		UploadDataset:         make(chan *ReqUploadDataset),
		RunProcessorOnDataset: make(chan *ReqRunProcessorOnDataset),
	}
}

// SensorChans is used by scheduler to determine global state
type SensorChans struct {
	IsLeader    chan bool
	GlobalState chan chan *types.GlobalState // because we need lazy calculation of result
}

func newSensorChans() *SensorChans {
	return &SensorChans{
		IsLeader:    make(chan bool),
		GlobalState: make(chan chan *types.GlobalState),
	}
}
