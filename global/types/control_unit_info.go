package types

// ControlUnitInfo - represents current control unit status
type ControlUnitInfo struct {
	UID     UID
	Workers []WorkerInfo
}
