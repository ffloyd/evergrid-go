package types

// JobInfo - information of current job on worker
type JobInfo struct {
	UID         UID
	Worker      *WorkerInfo
	Processor   *ProcessorInfo
	MFlopsTotal MFlop
	MFlopsDone  MFlop
	Completed   bool
}
