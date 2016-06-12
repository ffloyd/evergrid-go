package types

// DatasetInfo - represents current status of dataset
type DatasetInfo struct {
	UID               string
	Size              MByte
	Workers           []string
	EnqueuedOnWorkers []string
}
