package scheduler

// ReqUploadDataset - defines request to upload new dataset
type ReqUploadDataset struct {
	DatasetID string
	Response  *RespUploadDataset
}

// NewReqUploadDataset returns a new prepared request to scheduler
func NewReqUploadDataset(datasetID string) *ReqUploadDataset {
	return &ReqUploadDataset{
		DatasetID: datasetID,
		Response:  newRespUploadDataset(),
	}
}

// ReqRunProcessorOnDataset - defines request to run processor on given dataset
type ReqRunProcessorOnDataset struct {
	DatasetID   string
	ProcessorID string
	Response    *RespRunProcessorOnDataset
}

// NewReqRunProcessorOnDataset returns a new prepared request to scheduler
func NewReqRunProcessorOnDataset(datasetID string, processorID string) *ReqRunProcessorOnDataset {
	return &ReqRunProcessorOnDataset{
		DatasetID:   datasetID,
		ProcessorID: processorID,
		Response:    newRespRunProcessorOnDataset(),
	}
}
