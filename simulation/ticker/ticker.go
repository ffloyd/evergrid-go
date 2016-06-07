package ticker

import (
	log "github.com/Sirupsen/logrus"
)

// Ticker is a global timer like NetLogo's one
type Ticker struct {
	currentTick int

	ticksChan      chan int
	statusChan     chan SyncableStatus
	startWorkChan  chan bool
	finishWorkChan chan bool
	stopFlagChan   chan bool
}

// New creates a new Ticker.
func New(sync Syncable) *Ticker {
	defer log.Info("New ticker initialized")

	return &Ticker{
		ticksChan:      sync.CreateTicksChan(),
		statusChan:     sync.CreateStatusChan(),
		startWorkChan:  sync.CreateStartWorkChan(),
		finishWorkChan: sync.CreateFinishWorkChan(),
		stopFlagChan:   sync.CreateStopFlagChan(),
	}
}

// Run starts ticker
func (ticker *Ticker) Run() {
MainLoop:
	for {
		// send new tick
		ticker.currentTick++
		log.WithField("tick", ticker.currentTick).Debug("Ticker: new tick")
		ticker.ticksChan <- ticker.currentTick

		// confirm receiving
		readyStatus := <-ticker.statusChan
		if readyStatus != StatusReady {
			panic("Invalid status for syncable")
		}
		log.WithField("tick", ticker.currentTick).Debug("Ticker: all agents are ready")

		// initiate work
		ticker.startWorkChan <- true
		log.WithField("tick", ticker.currentTick).Debug("Ticker: start agent's work")

		// wait for Idle status
		for {
			status := <-ticker.statusChan
			if status == StatusIdle {
				break
			}
		}
		log.WithField("tick", ticker.currentTick).Debug("Ticker: all agents became idle")

		// stop work for this tick
		ticker.finishWorkChan <- true

		// confirm that work finished
		doneStatus := <-ticker.statusChan
		if doneStatus != StatusDone {
			panic("Invalid status for syncable")
		}
		log.WithField("tick", ticker.currentTick).Debug("Ticker: all agents finish work")

		// check for stopFlag
		select {
		case stopFlag := <-ticker.stopFlagChan:
			if stopFlag {
				break MainLoop
			}
		default:
		}
	}
	log.WithField("tick", ticker.currentTick).Info("Ticker: experiment finished")
}
