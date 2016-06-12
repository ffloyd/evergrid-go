package types

// CalculatorInfo - represents current worker status
type CalculatorInfo struct {
	UID               string
	MFlopsPerMb       MFlop
	Workers           []string
	EnqueuedOnWorkers []string
}
