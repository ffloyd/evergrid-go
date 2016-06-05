package types

// ProcessorInfo - represents current worker status
type ProcessorInfo struct {
	UID         UID
	MFlopsPerMb MFlop
	Workers     []UID
}
