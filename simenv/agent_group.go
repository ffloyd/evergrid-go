package simenv

import "reflect"

type agentGroup struct {
	chans []AgentChans
	names []string

	statuses  []AgentState
	stopFlags []bool

	stateChange chan AgentState
}

func runAgentGroup(simenv *SimEnv) *agentGroup {
	agentsCount := len(simenv.agents)

	result := &agentGroup{
		chans:       make([]AgentChans, agentsCount),
		names:       make([]string, agentsCount),
		statuses:    make([]AgentState, agentsCount),
		stopFlags:   make([]bool, agentsCount),
		stateChange: make(chan AgentState),
	}

	i := 0
	for name, agent := range simenv.agents {
		result.chans[i] = agent.Run(simenv)
		result.names[i] = name
		result.statuses[i] = StateDone
		result.stopFlags[i] = false
		i++
	}

	go result.statusChanWorker()
	go result.stopFlagsWorker()

	return result
}

func (group *agentGroup) WaitForState(status AgentState) {
	for <-group.stateChange != status {
	}
}

func (group *agentGroup) StartWork() {
	for _, agentChans := range group.chans {
		agentChans.startWorkChan <- Ok{}
	}
}

func (group *agentGroup) FinishWork() {
	for _, agentChans := range group.chans {
		agentChans.finishWorkChan <- Ok{}
	}
}

func (group *agentGroup) StopFlag() bool {
	for _, value := range group.stopFlags {
		if !value {
			return false
		}
	}

	return true
}

func (group *agentGroup) statusChanWorker() {
	cases := make([]reflect.SelectCase, len(group.chans))
	for i, agentChans := range group.chans {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(agentChans.statusChan),
		}
	}

	isSimilar := func() bool {
		firstStatus := group.statuses[0]

		for i := 1; i < len(group.statuses); i++ {
			if group.statuses[i] != firstStatus {
				return false
			}
		}

		return true
	}

	for {
		chosen, rawValue, ok := reflect.Select(cases)
		if ok != true {
			panic("agentGroup fail")
		}

		similarBefore := isSimilar()

		value := AgentState(rawValue.Int())
		group.statuses[chosen] = value

		similarAfter := isSimilar()

		if similarBefore == false && similarAfter == true {
			group.stateChange <- value
		}
	}
}

func (group *agentGroup) stopFlagsWorker() {
	cases := make([]reflect.SelectCase, len(group.chans))
	for i, agentChans := range group.chans {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(agentChans.stopFlagChan),
		}
	}

	for {
		chosen, rawValue, ok := reflect.Select(cases)
		if ok != true {
			panic("agentGroup fail")
		}

		group.stopFlags[chosen] = rawValue.Bool()
	}
}
