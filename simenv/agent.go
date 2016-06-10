package simenv

// Agent -
type Agent interface {
	Name() string
	Run(*SimEnv) AgentChans
	Send(interface{}) chan interface{}
}
