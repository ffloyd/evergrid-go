package simenv

/*
SimEnv -
*/
type SimEnv struct {
}

// NewSimEnv -
func NewSimEnv(agentGroup *AgentGroup) *SimEnv {
	return &SimEnv{}
}

// Run -
func (simenv *SimEnv) Run() {
}

// CurrentTick -
func (simenv *SimEnv) CurrentTick() CurrentTick {
	return CurrentTick(0)
}
