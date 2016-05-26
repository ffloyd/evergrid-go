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
