package ticker

// SyncableStatus - represent statuses for Syncable
type SyncableStatus int

const (
	// StatusReady - means that new tick received
	StatusReady SyncableStatus = iota

	// StatusWorking - means that current entity doing some work now
	StatusWorking

	// StatusIdle - means that current entity active and waits for interactions from other entities
	StatusIdle

	// StatusDone - means that current entity finished all possible work for current tick
	StatusDone
)

// Syncable - specification of entity which can work with ticker
type Syncable interface {
	CreateStatusChan() chan SyncableStatus
	CreateTicksChan() chan int
	CreateStartWorkChan() chan bool
	CreateFinishWorkChan() chan bool
}

// String in implementation of stringer interface
func (status SyncableStatus) String() string {
	switch status {
	case StatusReady:
		return "StatusReady"
	case StatusWorking:
		return "StatusWorking"
	case StatusIdle:
		return "StatusIdle"
	case StatusDone:
		return "StatusDone"
	default:
		panic("Unknown constant value")
	}
}
