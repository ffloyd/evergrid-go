package scheduler

// Done -
type Done struct{}

// ControlChans -
type ControlChans struct {
	UploadDataset   chan DoUploadDataset
	BuildCalculator chan DoBuildCalculator
	RunCalculator   chan DoRunCalculator

	Done chan Done
}

// NewControlChans -
func NewControlChans() ControlChans {
	return ControlChans{
		UploadDataset:   make(chan DoUploadDataset),
		BuildCalculator: make(chan DoBuildCalculator),
		RunCalculator:   make(chan DoRunCalculator),

		Done: make(chan Done),
	}
}

// DoUploadDataset -
type DoUploadDataset struct {
	Dataset string
	Worker  string
}

// DoBuildCalculator -
type DoBuildCalculator struct {
	Calculator string
	Worker     string
}

// DoRunCalculator -
type DoRunCalculator struct {
	Calculator string
	Dataset    string
	Worker     string
}
