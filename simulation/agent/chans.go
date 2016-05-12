package agent

// TickerChans is struct for communication between Ticker and agents
type TickerChans struct {
	Ready chan bool // for ready status (agent -> ticker)
	Ticks chan int  // for ticks broadcasting (ticker -> agent)
}

// NewTickerChans initializes correct Chans instanse
func NewTickerChans() *TickerChans {
	return &TickerChans{make(chan bool), make(chan int)}
}
