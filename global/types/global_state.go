package types

// GlobalState - represents system-wide state
type GlobalState struct {
	ControlUnits map[UID]*ControlUnitInfo
	Datasets     map[UID]*DatasetInfo
	Processors   map[UID]*ProcessorInfo
	Workers      map[UID]*WorkerInfo

	ActiveJobs map[UID]*JobInfo
}

// NewGlobalState returns global state with initialized maps
func NewGlobalState() *GlobalState {
	return &GlobalState{
		ControlUnits: make(map[UID]*ControlUnitInfo),
		Datasets:     make(map[UID]*DatasetInfo),
		Processors:   make(map[UID]*ProcessorInfo),
		Workers:      make(map[UID]*WorkerInfo),

		ActiveJobs: make(map[UID]*JobInfo),
	}
}
