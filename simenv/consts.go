package simenv

// Ok -
type Ok struct{}

// AgentState - represent statuses for Syncable
type AgentState int

const (
	// StateReady - means that new tick received
	StateReady AgentState = iota

	// StateWorking - means that current entity doing some work now
	StateWorking

	// StateIdle - means that current entity active and waits for interactions from other entities
	StateIdle

	// StateDone - means that current entity finished all possible work for current tick
	StateDone
)

func (state AgentState) String() string {
	switch state {
	case StateReady:
		return "StateReady"
	case StateWorking:
		return "StateWorking"
	case StateIdle:
		return "StateIdle"
	case StateDone:
		return "StateDone"
	default:
		panic("Unknown state value")
	}
}
