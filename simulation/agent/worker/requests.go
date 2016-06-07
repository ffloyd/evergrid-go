package worker

import "github.com/ffloyd/evergrid-go/global/types"

// ReqUpload is a request form to worker
type ReqUpload struct {
	Dataset *types.DatasetInfo
}

// ReqBuild is a request form to worker
type ReqBuild struct {
	Processor *types.ProcessorInfo
}

// ReqExecute is a request form to worker
type ReqExecute struct {
	Processor *types.ProcessorInfo
	Dataset   *types.DatasetInfo
}
