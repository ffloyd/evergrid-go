package types

// WorkerInfo - represents current worker status
type WorkerInfo struct {
	UID            string
	Busy           bool
	MFlops         MFlop
	TotalDiskSpace MByte
	FreeDiskSpace  MByte
	PricePerTick   float64
	Datasets       []string
	Calculators    []string
	ControlUnit    string

	QueueLength        int
	DatasetsInQueue    []string
	CalculatorsInQueue []string
}
