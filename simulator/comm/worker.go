package comm

import "github.com/ffloyd/evergrid-go/global/types"

// WorkerUploadDataset -
type WorkerUploadDataset struct {
	Dataset types.DatasetInfo
}

// WorkerBuildCalculator -
type WorkerBuildCalculator struct {
	Calculator types.CalculatorInfo
}

// WorkerBusy -
type WorkerBusy struct{}
