package types

// GlobalState - represents system-wide state
type GlobalState struct {
	ControlUnits map[UID]*ControlUnitInfo
	Datasets     map[UID]*DatasetInfo
	Processors   map[UID]*ProcessorInfo
	Workers      map[UID]*WorkerInfo

	ActiveJobs map[UID]*JobInfo
}
