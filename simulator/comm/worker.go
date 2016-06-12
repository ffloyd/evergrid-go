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

// WorkerRunCalculator -
type WorkerRunCalculator struct {
	Calculator string
	Dataset    string
}

// WorkerBusy -
type WorkerBusy struct{}

// WorkerInfo -
type WorkerInfo struct{}
