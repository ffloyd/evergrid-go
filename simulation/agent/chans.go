package agent

// Chans is struct for communication between Ticker and agents
type Chans struct {
	Ready chan bool // for ready status (agent -> ticker)
	Ticks chan int  // for ticks broadcasting (ticker -> agent)
}

// NewChans initializes correct Chans instanse
func NewChans() *Chans {
	return &Chans{make(chan bool), make(chan int)}
}
