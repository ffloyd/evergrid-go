package scheduler

import "github.com/ffloyd/evergrid-go/global/types"

// RespDelegateToLeader - delegate this task to leader. Terminal response.
type RespDelegateToLeader struct{}

// RespDone - standart terminal response
type RespDone struct{}

// RespUploadDatasetToWorker - adds uploading dataset to worker to queue
type RespUploadDatasetToWorker struct {
	Worker  types.UID
	Dataset types.UID
}

// RespBuildProcessor - adds build stage of processor to worker's queue
type RespBuildProcessor struct {
	Worker    types.UID
	Processor types.UID
}

// RespRunProcessor - adds execution of processor to worker's queue
type RespRunProcessor struct {
	Worker    types.UID
	Processor types.UID
	Dataset   types.UID
}
