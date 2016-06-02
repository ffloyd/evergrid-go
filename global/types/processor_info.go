package types

// ProcessorInfo - represents current worker status
type ProcessorInfo struct {
	UID     UID
	Workers []*WorkerInfo
}
