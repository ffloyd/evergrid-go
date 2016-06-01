package types

import (
	"time"
)

// Prediction - represents execution prediction result
type Prediction struct {
	Estimate  time.Duration
	Worker    *WorkerInfo
	Dataset   *DatasetInfo
	Processor *ProcessorInfo
}
