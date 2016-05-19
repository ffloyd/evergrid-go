package scheduler

// Algorithm is an enum for algorithms represented in system
type Algorithm int

const (
	// FIFO is a simple "first in first out" implementation of scheduler
	FIFO Algorithm = iota
)

// Scheduler is just a scheduler
type Scheduler struct {
	algorithm Algorithm
	Chans     *Chans
}

// New creates new scheduler which used specified algorithm
func New(alg Algorithm) *Scheduler {
	sched := &Scheduler{
		algorithm: alg,
		Chans:     newChans(),
	}

	return sched
}

// Run starts scheduler work
func (sched *Scheduler) Run() {
	switch sched.algorithm {
	case FIFO:
		impl := &fifoScheduler{
			base: sched,
		}
		impl.run()
	}
}
