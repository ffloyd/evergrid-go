package comm

import "github.com/ffloyd/evergrid-go/global/types"

// ControlUnitUploadDataset - запрос на загрузку датасета
// Может приходить Control Unit'у либо от Core, либо от другого Control Unit'а (в случае делегирования лидеру)
type ControlUnitUploadDataset struct {
	Dataset types.DatasetInfo
}

// ControlUnitRunExperiment - запрос на запуск вычислителя с указанным датасетом
// Может приходить Control Unit'у либо от Core, либо от другого Control Unit'а (в случае делегирования лидеру)
type ControlUnitRunExperiment struct {
	Calculator types.CalculatorInfo
	Dataset    types.DatasetInfo
}
