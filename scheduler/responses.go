package scheduler

// RespUploadDataset - response for UploadDataset action
type RespUploadDataset struct {
	DelegateToLeader      chan RespDelegateToLeader
	UploadDatasetToWorker chan RespUploadDatasetToWorker
	Done                  chan RespDone
}

func newRespUploadDataset() *RespUploadDataset {
	return &RespUploadDataset{
		DelegateToLeader:      make(chan RespDelegateToLeader),
		UploadDatasetToWorker: make(chan RespUploadDatasetToWorker),
		Done: make(chan RespDone),
	}
}

// RespRunProcessorOnDataset - response for ReqRunProcessorOnDataset
type RespRunProcessorOnDataset struct {
	DelegateToLeader      chan RespDelegateToLeader
	UploadDatasetToWorker chan RespUploadDatasetToWorker
	BuildProcessor        chan RespBuildProcessor
	RunProcessor          chan RespRunProcessor
	Done                  chan RespDone
}

func newRespRunProcessorOnDataset() *RespRunProcessorOnDataset {
	return &RespRunProcessorOnDataset{
		DelegateToLeader:      make(chan RespDelegateToLeader),
		UploadDatasetToWorker: make(chan RespUploadDatasetToWorker),
		BuildProcessor:        make(chan RespBuildProcessor),
		RunProcessor:          make(chan RespRunProcessor),
		Done:                  make(chan RespDone),
	}
}
