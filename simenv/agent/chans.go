package agent

// Chans is struct for communication between Ticker and agents
type Chans struct {
	Ready chan bool // for incoming ready status
	Ticks chan int  // for ticks broadcasting
}

// NewChans initializes correct Chans instanse
func NewChans() *Chans {
	return &Chans{make(chan bool), make(chan int)}
}
