package simenv

type SimpleAgent struct {
	simenv *SimEnv

	fsm  AgentFSM
	name string

	MessageCount int
}

func NewSimpleAgent(name string) *SimpleAgent {
	return &SimpleAgent{
		name: name,
	}
}

func (agent *SimpleAgent) Name() string {
	return agent.name
}

func (agent *SimpleAgent) Run(simenv *SimEnv) AgentChans {
	agent.simenv = simenv
	agent.fsm = *NewAgentFSM(nil)
	go agent.work()
	return agent.fsm.Chans()
}

func (agent *SimpleAgent) Send(msg interface{}) chan interface{} {
	agent.MessageCount++
	return nil
}

func (agent *SimpleAgent) work() {
	i := 0

	for {
		agent.fsm.ToReady()
		agent.fsm.ToWorking()
		if i < 3 {
			for name, anotherAgent := range agent.simenv.agents {
				if name == agent.name {
					continue
				}

				anotherAgent.Send("message")
			}
		} else {
			agent.fsm.SetStopFlag(true)
		}

		agent.fsm.ToIdle()

		<-agent.fsm.ToDoneChan()

		i++
	}

}
