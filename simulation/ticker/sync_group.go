package ticker

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
)

// SyncGroup is a group of syncables which works as one syncable
type SyncGroup []Syncable

// CreateStatusChan is an implementation of Syncable interface
func (group SyncGroup) CreateStatusChan() chan SyncableStatus {
	result := make(chan SyncableStatus)
	go group.statusChanWorker(result)
	return result
}

// CreateTicksChan is an implementation of Syncable interface
func (group SyncGroup) CreateTicksChan() chan int {
	result := make(chan int)
	go group.ticksChanWorker(result)
	return result
}

// CreateStartWorkChan is an implementation of Syncable interface
func (group SyncGroup) CreateStartWorkChan() chan bool {
	result := make(chan bool)
	go group.startWorkChanWorker(result)
	return result
}

// CreateFinishWorkChan is an implementation of Syncable interface
func (group SyncGroup) CreateFinishWorkChan() chan bool {
	result := make(chan bool)
	go group.finishWorkChanWorker(result)
	return result
}

func (group SyncGroup) statusChanWorker(parentChan chan SyncableStatus) {
	statuses := make([]SyncableStatus, len(group))
	for i := range statuses {
		statuses[i] = StatusDone
	}

	cases := make([]reflect.SelectCase, len(group))
	for i, sync := range group {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(sync.CreateStatusChan()),
		}
	}

	isSimilar := func() bool {
		previous := statuses[0]
		for i := 1; i < len(statuses); i++ {
			if previous != statuses[i] {
				return false
			}
			previous = statuses[i]
		}
		return true
	}

	for {
		chosen, rawValue, ok := reflect.Select(cases)
		if ok != true {
			log.Panic("SyncGroup fail")
		}

		similarBefore := isSimilar()

		value := SyncableStatus(rawValue.Int())
		statuses[chosen] = value

		similarAfter := isSimilar()

		if similarBefore == false && similarAfter == true {
			parentChan <- value
		}
	}
}

func (group SyncGroup) ticksChanWorker(parentChan chan int) {
	nestedChans := make([]chan int, len(group))
	for i, sync := range group {
		nestedChans[i] = sync.CreateTicksChan()
	}

	for {
		newTick := <-parentChan
		for _, nestedChan := range nestedChans {
			nestedChan <- newTick
		}
	}
}

func (group SyncGroup) startWorkChanWorker(parentChan chan bool) {
	nestedChans := make([]chan bool, len(group))
	for i, sync := range group {
		nestedChans[i] = sync.CreateStartWorkChan()
	}

	for {
		message := <-parentChan
		for _, nestedChan := range nestedChans {
			nestedChan <- message
		}
	}
}

func (group SyncGroup) finishWorkChanWorker(parentChan chan bool) {
	nestedChans := make([]chan bool, len(group))
	for i, sync := range group {
		nestedChans[i] = sync.CreateFinishWorkChan()
	}

	for {
		message := <-parentChan
		for _, nestedChan := range nestedChans {
			nestedChan <- message
		}
	}
}
