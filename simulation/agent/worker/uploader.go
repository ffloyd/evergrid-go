package worker

import (
	"fmt"
	"math"

	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
)

// Uploader is an internal component for worker which manages dataset uploads
type Uploader struct {
	state *State

	uploading bool
	dataset   *types.DatasetInfo
	uploaded  types.MByte
	speed     types.MBit
}

// NewUploader creates a new builder instance
func NewUploader(state *State) *Uploader {
	return &Uploader{
		state: state,
	}
}

// Prepare initate uploading process
func (uploader *Uploader) Prepare(request ReqUpload) {
	// Check if dataset already uploaded
	if uploader.state.HasDataset(request.Dataset) {
		log.WithFields(log.Fields{
			"agent":   uploader.state.info.UID,
			"dataset": request.Dataset.UID,
		}).Info("Dataset already presents on this worker")
		return
	}

	uploader.state.Busy()
	uploader.uploading = true
	uploader.dataset = request.Dataset
	uploader.uploaded = 0

	// TODO: enable correct speed calculation
	// Check if dataset presents in current segment
	// internalComm := false
	// segmentAgentNames := worker.Node().Segment().AgentNames()
	// for _, agentName := range segmentAgentNames {
	// 	closeWorker, ok := worker.env.Workers[agentName]
	// 	if ok {
	// 		if closeWorker.State.Datasets[datasetUID] != nil {
	// 			internalComm = true
	// 			break
	// 		}
	// 	}
	// }
	//
	// // Reolve upload speed
	// bandwith := worker.Node().Segment().Bandwith(internalComm)
	// if bandwith.In < bandwith.Out {
	// 	upload.speed = types.MBit(bandwith.In)
	// } else {
	// 	upload.speed = types.MBit(bandwith.Out)
	// }

	uploader.speed = 100

	log.WithFields(log.Fields{
		"agent":   uploader.state.info.UID,
		"dataset": uploader.dataset.UID,
	}).Info("Initiate dataset upload")
}

// Process performs upload activity
func (uploader *Uploader) Process() {
	if !uploader.uploading {
		return
	}

	// 1 tick = 1 minute
	mbytesDownloaded := types.MByte(uploader.speed * 60 / 8)
	uploader.uploaded += mbytesDownloaded

	if uploader.uploaded >= uploader.dataset.Size {
		uploader.uploading = false
		uploader.state.AddDataset(uploader.dataset)
		uploader.state.Free()

		log.WithFields(log.Fields{
			"agent":   uploader.state.info.UID,
			"dataset": uploader.dataset.UID,
		}).Info("Dataset uploaded")
	} else {
		progress := math.Min(1.0, float64(uploader.uploaded)/float64(uploader.dataset.Size))

		log.WithFields(log.Fields{
			"agent":    uploader.state.info.UID,
			"dataset":  uploader.dataset.UID,
			"progress": fmt.Sprintf("%d%%", int(progress*100)),
		}).Info("Uploading dataset")
	}
}
