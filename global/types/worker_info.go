package types

// WorkerInfo - represents current worker status
type WorkerInfo struct {
	UID            UID
	Busy           bool
	MFlops         MFlop
	TotalDiskSpace MByte
	FreeDiskSpace  MByte
	PricePerTick   float64
	CurrentJob     *JobInfo
	Datasets       map[UID]*DatasetInfo
	Processors     map[UID]*ProcessorInfo
	QueueLength    int
}
