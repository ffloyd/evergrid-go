package scheduler

import "github.com/ffloyd/evergrid-go/global/types"

// Algorithm is an enum for algorithms represented in system
type Algorithm int

const (
	// Random is a scheduler, who must work worse than any meaningful other
	Random Algorithm = iota
	// FIFO is a simple "first in first out" implementation of scheduler
	FIFO
)

// Scheduler is just a scheduler
type Scheduler struct {
	ID        string
	algorithm Algorithm
	Chans     *Chans
}

// New creates new scheduler which used specified algorithm
func New(alg Algorithm, agentName string) *Scheduler {
	sched := &Scheduler{
		ID:        agentName,
		algorithm: alg,
		Chans:     newChans(),
	}

	return sched
}

// Run starts scheduler work
func (sched *Scheduler) Run() {
	switch sched.algorithm {
	case Random:
		impl := &randomScheduler{
			base:             sched,
			datasetLocations: make(map[types.UID]types.UID),
		}
		impl.run()
	case FIFO:
		impl := &fifoScheduler{
			base: sched,
		}
		impl.run()
	}
}
