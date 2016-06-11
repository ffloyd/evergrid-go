package gendata

// Config - is a config for generate new simdata files
type Config struct {
	DestDir string
	Name    string

	DatsetsCount   int
	MinDatasetSize int // in GBytes
	MaxDatasetSize int

	ProcessorsCount int // in MFlopsPerMb
	MinSpeed        int
	MaxSpeed        int

	ProcessorRuns  int
	RunProbability float64 // probability of task execution per tick

	NetworkSegments   int
	MinNodesInSegment int
	MaxNodesInSegment int
	MinNodeSpeed      int // MFlops
	MaxNodeSpeed      int
	MinDiskSize       int // GBytes
	MaxDiskSize       int
	MinPricePerTick   float64
	MaxPricePerTick   float64
}
