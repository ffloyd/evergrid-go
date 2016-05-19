package scheduler

// Chans is a set of chans for communicate with scheduler
type Chans struct {
	Alive chan bool
}

func newChans() *Chans {
	return &Chans{
		Alive: make(chan bool),
	}
}
