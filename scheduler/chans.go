package scheduler

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
	UploadDataset chan *ReqUploadDataset
}

func newRequestChans() *RequestChans {
	return &RequestChans{
		UploadDataset: make(chan *ReqUploadDataset),
	}
}

// SensorChans is used by scheduler to determine global state
type SensorChans struct {
	IsLeader chan bool
}

func newSensorChans() *SensorChans {
	return &SensorChans{
		IsLeader: make(chan bool),
	}
}

// ReqUploadDataset - defines request to upload new dataset
type ReqUploadDataset struct {
	DatasetID string
	Response  chan *RespUploadDataset
}

// RespUploadDataset - response for UploadDataset action
type RespUploadDataset struct {
	DelegateToLeader bool
	UploadToWorkers  []string // contain worker names
}
