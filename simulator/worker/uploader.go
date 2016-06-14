package worker

import (
	"fmt"
	"math"

	"github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simulator/comm"
)

type uploader struct {
	worker *Worker

	uploading bool
	dataset   types.DatasetInfo
	uploaded  types.MByte
	speed     types.MBit

	uploadedDatasets map[string]types.DatasetInfo

	log *logrus.Entry
}

func newUploader(w *Worker, logContext *logrus.Entry) uploader {
	return uploader{
		worker:           w,
		uploading:        false,
		uploadedDatasets: make(map[string]types.DatasetInfo),
		log:              logContext,
	}
}

func (up *uploader) Datasets() map[string]types.DatasetInfo {
	result := make(map[string]types.DatasetInfo)
	for k, v := range up.uploadedDatasets {
		result[k] = v
	}
	return result
}

func (up *uploader) Prepare(request comm.WorkerUploadDataset) {
	dataset := request.Dataset

	_, hasDataset := up.uploadedDatasets[dataset.UID]
	if hasDataset {
		up.log.WithField("dataset", dataset.UID).Info("Dataset already present on worker")
		return
	}

	up.worker.busy = true
	up.worker.fsm.SetStopFlag(false)
	up.uploading = true
	up.dataset = dataset
	up.uploaded = 0
	up.speed = 100

	up.log.WithField("dataset", dataset.UID).Info("Dataset uploading initiated")
}

func (up *uploader) Process() {
	if !up.uploading {
		return
	}

	up.worker.stats.UploadingTicks++
	up.worker.stats.MoneySpent += up.worker.pricePerTick

	// 1 tick = 1 minute
	mbytesDownloaded := types.MByte(up.speed * 60 / 8)
	up.uploaded += mbytesDownloaded

	if up.uploaded >= up.dataset.Size {
		up.uploading = false
		up.uploadedDatasets[up.dataset.UID] = up.dataset
		up.worker.busy = false
		up.worker.fsm.SetStopFlag(true)

		up.log.WithField("dataset", up.dataset.UID).Info("Dataset uploaded")
	} else {
		progress := math.Min(1.0, float64(up.uploaded)/float64(up.dataset.Size))

		up.log.WithFields(logrus.Fields{
			"dataset":  up.dataset.UID,
			"progress": fmt.Sprintf("%d%%", int(progress*100)),
		}).Info("Uploading dataset")
	}
}
