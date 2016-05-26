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
