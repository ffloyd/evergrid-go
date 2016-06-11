package comm

import "github.com/ffloyd/evergrid-go/global/types"

// ControlUnitUploadDataset -
type ControlUnitUploadDataset struct {
	Dataset types.DatasetInfo
}

// ControlUnitRunExperiment -
type ControlUnitRunExperiment struct {
	Calculator types.CalculatorInfo
	Dataset    types.DatasetInfo
}
