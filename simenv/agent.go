package simenv

// Agent -
type Agent interface {
	Name() string
	Run(*AgentGroup) AgentChans
	Send(interface{}) chan interface{}
}
