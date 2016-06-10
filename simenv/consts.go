package simenv

// Ok -
type Ok struct{}

// AgentStatus - represent statuses for Syncable
type AgentStatus int

const (
	// StatusReady - means that new tick received
	StatusReady AgentStatus = iota

	// StatusWorking - means that current entity doing some work now
	StatusWorking

	// StatusIdle - means that current entity active and waits for interactions from other entities
	StatusIdle

	// StatusDone - means that current entity finished all possible work for current tick
	StatusDone
)
