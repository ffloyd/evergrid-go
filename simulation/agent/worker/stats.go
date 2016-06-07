package worker

// Stats is a statistics aggregation structure for worker
type Stats struct {
	ExecutionTicks int
	UploadTicks    int
	BuildTicks     int
}

// NewStats returns correctly initialized Stats structure
func NewStats() *Stats {
	return &Stats{}
}
