package ticker

/*
SyncChans -
*/
type SyncChans struct {
	ticksChan      chan int
	statusChan     chan SyncableStatus
	startWorkChan  chan bool
	finishWorkChan chan bool
	stopFlagChan   chan bool
}

// NewSyncChans -
func NewSyncChans() SyncChans {
	return SyncChans{
		ticksChan:      make(chan int),
		statusChan:     make(chan SyncableStatus),
		startWorkChan:  make(chan bool),
		finishWorkChan: make(chan bool),
		stopFlagChan:   make(chan bool),
	}
}
