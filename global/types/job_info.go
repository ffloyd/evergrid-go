package types

// JobInfo - information of current job on worker
type JobInfo struct {
	UID         UID
	Type        JobType
	Worker      UID
	Processor   UID
	Dataset     UID
	MFlopsTotal MFlop
	MFlopsDone  MFlop
	Completed   bool
}

// JobType is a enum for possible job types
type JobType int

// JobType possible values
const (
	JobUploadDataset JobType = iota
	JobBuildProcessor
	JobRunProcessor
)
