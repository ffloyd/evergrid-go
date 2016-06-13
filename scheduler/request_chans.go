package scheduler

import "github.com/ffloyd/evergrid-go/global/types"

/*
RequestChans - это каналы, по которым в планировщик приходят запросы.

DelegateToLeader - это канал по которому планировщик говорит о том, что:

  в случае отправки true - надо лелегировать этот запрос лидеру
  в случае отправки false - что запрос принят и можно присылать следующий

Данная группа каналов является "синхронной" - т. е. запросы в планировщик должны
приходить последовательно.
*/
type RequestChans struct {
	UploadDataset chan ReqUploadDataset
	RunExperiment chan ReqRunExperiment

	DelegateToLeader chan bool
}

// NewRequestChans - инициализатор для RequestChans
func NewRequestChans() RequestChans {
	return RequestChans{
		UploadDataset:    make(chan ReqUploadDataset),
		RunExperiment:    make(chan ReqRunExperiment),
		DelegateToLeader: make(chan bool),
	}
}

// ReqUploadDataset - запрос на загрузку датасета в систему.
type ReqUploadDataset struct {
	Dataset types.DatasetInfo
}

// ReqRunExperiment - запрос на запуск вычислителя с данным датасетом.
type ReqRunExperiment struct {
	Calculator types.CalculatorInfo
	Dataset    types.DatasetInfo
}
