package simenv

/*
SimEnv -
*/
type SimEnv struct {
	tick          int
	agents        map[string]Agent
	tickBroadcast []chan int
	inProgress    bool
}

// NewSimEnv -
func NewSimEnv() *SimEnv {
	return &SimEnv{
		agents: make(map[string]Agent),
	}
}

// Add -
func (simenv *SimEnv) Add(agents ...Agent) {
	for _, agent := range agents {
		simenv.agents[agent.Name()] = agent
	}
}

// Run -
func (simenv *SimEnv) Run() error {
	group := runAgentGroup(simenv)

	for {
		simenv.tick++
		for _, tickChan := range simenv.tickBroadcast {
			tickChan <- simenv.tick
		}

		group.WaitForState(StateReady)
		group.StartWork()
		group.WaitForState(StateIdle)
		group.FinishWork()

		if group.StopFlag() {
			break
		}
	}

	return nil
}

// CurrentTick -
func (simenv *SimEnv) CurrentTick() *CurrentTick {
	ct := &CurrentTick{simenv.tick}
	channel := make(chan int)

	simenv.tickBroadcast = append(simenv.tickBroadcast, channel)

	go ct.connect(channel)

	return ct
}
