package scheduler

// RespDelegateToLeader - delegate this task to leader. Terminal response.
type RespDelegateToLeader struct{}

// RespDone - standart terminal response
type RespDone struct{}

// RespUploadDatasetToWorker - adds uploading dataset to worker to queue
type RespUploadDatasetToWorker struct {
	Worker string
}
