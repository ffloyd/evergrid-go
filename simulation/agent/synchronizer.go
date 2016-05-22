package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation/ticker"
)

// Synchronizer is struct for communication between Ticker and agents
type Synchronizer struct {
	ticksChan      chan int // for ticks broadcasting (ticker -> agent)
	statusChan     chan ticker.SyncableStatus
	startWorkChan  chan bool
	finishWorkChan chan bool

	tick   int
	status ticker.SyncableStatus

	agentName string
}

// NewSynchronizer initializes correct Chans instanse
func NewSynchronizer(agentName string) *Synchronizer {
	return &Synchronizer{
		ticksChan:      make(chan int),
		statusChan:     make(chan ticker.SyncableStatus),
		startWorkChan:  make(chan bool),
		finishWorkChan: make(chan bool),
		status:         ticker.StatusDone,
		agentName:      agentName,
	}
}

// CreateStatusChan is an implementation of ticker.Syncable interface
func (sync *Synchronizer) CreateStatusChan() chan ticker.SyncableStatus {
	return sync.statusChan
}

// CreateTicksChan is an implementation of ticker.Syncable interface
func (sync *Synchronizer) CreateTicksChan() chan int {
	return sync.ticksChan
}

// CreateStartWorkChan is an implementation of ticker.Syncable interface
func (sync *Synchronizer) CreateStartWorkChan() chan bool {
	return sync.startWorkChan
}

// CreateFinishWorkChan is an implementation of ticker.Syncable interface
func (sync *Synchronizer) CreateFinishWorkChan() chan bool {
	return sync.finishWorkChan
}

// to[Status] functions. Define state machine for these statuses

// receives new tick besides of changing status form Done to Ready
func (sync *Synchronizer) toReady() {
	if sync.status != ticker.StatusDone {
		log.Panicf("Synchronizer: invalid transition from %v to %v", sync.status, ticker.StatusReady)
	}

	sync.tick = <-sync.ticksChan
	log.WithFields(log.Fields{
		"agent": sync.agentName,
		"tick":  sync.tick,
	}).Debug("Agent receive tick")

	oldStatus := sync.status
	sync.status = ticker.StatusReady
	log.WithFields(log.Fields{
		"agent": sync.agentName,
		"tick":  sync.tick,
		"from":  oldStatus,
		"to":    sync.status,
	}).Debug("Agent change status")
	sync.statusChan <- ticker.StatusReady

	<-sync.startWorkChan
	log.WithFields(log.Fields{
		"agent": sync.agentName,
		"tick":  sync.tick,
	}).Debug("Agent got permission to start work")
}

// waits for startWorkChan message besides of changing status from Done to Ready
func (sync *Synchronizer) toWorking() {
	if sync.status != ticker.StatusReady && sync.status != ticker.StatusIdle {
		log.Panicf("Synchronizer: invalid transition from %v to %v", sync.status, ticker.StatusWorking)
	}

	oldStatus := sync.status
	sync.status = ticker.StatusWorking
	log.WithFields(log.Fields{
		"agent": sync.agentName,
		"tick":  sync.tick,
		"from":  oldStatus,
		"to":    sync.status,
	}).Debug("Agent change status")
	sync.statusChan <- ticker.StatusWorking
}

func (sync *Synchronizer) toIdle() {
	if sync.status != ticker.StatusWorking && sync.status != ticker.StatusReady {
		log.Panicf("Synchronizer: invalid transition from %v to %v", sync.status, ticker.StatusIdle)
	}

	oldStatus := sync.status
	sync.status = ticker.StatusIdle
	log.WithFields(log.Fields{
		"agent": sync.agentName,
		"tick":  sync.tick,
		"from":  oldStatus,
		"to":    sync.status,
	}).Debug("Agent change status")
	sync.statusChan <- ticker.StatusIdle
}

// don't use directly - use toDoneCallback instead
func (sync *Synchronizer) toDone() {
	if sync.status != ticker.StatusIdle {
		log.Panicf("Synchronizer: invalid transition from %v to %v", sync.status, ticker.StatusDone)
	}

	oldStatus := sync.status
	sync.status = ticker.StatusDone
	log.WithFields(log.Fields{
		"agent": sync.agentName,
		"tick":  sync.tick,
		"from":  oldStatus,
		"to":    sync.status,
	}).Debug("Agent change status")
	sync.statusChan <- ticker.StatusDone
}

// creates one-off channel which reacts on finishWorkChan message
func (sync *Synchronizer) toDoneCallback() chan bool {
	result := make(chan bool)

	go func() {
		<-sync.finishWorkChan
		sync.toDone()
		result <- true
	}()

	return result
}
