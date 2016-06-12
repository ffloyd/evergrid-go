package scheduler

// Scheduler -
type Scheduler interface {
	Name() string

	Run()

	RequestChans() RequestChans
	ControlChans() ControlChans
	InfoChans() InfoChans
}
