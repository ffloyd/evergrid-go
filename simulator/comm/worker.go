package comm

import "github.com/ffloyd/evergrid-go/global/types"

// WorkerUploadDataset - указание ControlUnit'а Worker'у загрузить датасет
type WorkerUploadDataset struct {
	Dataset types.DatasetInfo
}

// WorkerBuildCalculator - указание ControlUnit'а Worker'у собрать вычислитель
type WorkerBuildCalculator struct {
	Calculator types.CalculatorInfo
}

// WorkerRunCalculator  - указание ControlUnit'а Worker'у запустить вычислитель с указанным датасетом
type WorkerRunCalculator struct {
	Calculator string
	Dataset    string
}

// WorkerBusy - запрос ControlUnit'а Worker'у. Если в данный момент выполняется какая-либо задача будет возращено true
type WorkerBusy struct{}

// WorkerInfo - запрос ControlUnit'а Worker'у. Возвращается частично заполненный types.WorkerInfo.
//
// Частичная заполненность обуславливается, например, тем, что Worker не знает о том, сколько задачи у него в очереди на выполнение.
type WorkerInfo struct{}
