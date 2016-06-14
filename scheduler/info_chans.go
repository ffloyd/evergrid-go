package scheduler

import "github.com/ffloyd/evergrid-go/global/types"

/*
InfoChans это каналы, по которым планировщик узнает о внешнем состоянии системы.

Запрос происходит путем отправки сообщения. Используемые структуры содержат канал
для получения ответа.

Пример использования:

	request := NewGetWorkerNames()
	ic.WorkerNames <- request
	names := <-request.Result

Данные каналы могут использоваться асинхронно.
*/
type InfoChans struct {
	WorkerNames      chan GetWorkerNames
	WorkerInfo       chan GetWorkerInfo
	DatasetInfo      chan GetDatasetInfo
	CalculatorInfo   chan GetCalculatorInfo
	LeadershipStatus chan GetLeadershipStatus
}

// NewInfoChans - инициализатор для InfoChans
func NewInfoChans() InfoChans {
	return InfoChans{
		WorkerNames:    make(chan GetWorkerNames),
		WorkerInfo:     make(chan GetWorkerInfo),
		DatasetInfo:    make(chan GetDatasetInfo),
		CalculatorInfo: make(chan GetCalculatorInfo),

		LeadershipStatus: make(chan GetLeadershipStatus),
	}
}

// GetWorkerNames реализует соответствующий запрос с блокировкой до получения результата
func (ic InfoChans) GetWorkerNames() []string {
	request := NewGetWorkerNames()
	ic.WorkerNames <- request
	return <-request.Result
}

// GetWorkerInfo реализует соответствующий запрос с блокировкой до получения результата
func (ic InfoChans) GetWorkerInfo(uid string) types.WorkerInfo {
	request := NewGetWorkerInfo(uid)
	ic.WorkerInfo <- request
	return *(<-request.Result)
}

// GetDatasetInfo реализует соответствующий запрос с блокировкой до получения результата
func (ic InfoChans) GetDatasetInfo(uid string) types.DatasetInfo {
	request := NewGetDatasetInfo(uid)
	ic.DatasetInfo <- request
	return *(<-request.Result)
}

// GetCalculatorInfo реализует соответствующий запрос с блокировкой до получения результата
func (ic InfoChans) GetCalculatorInfo(uid string) types.CalculatorInfo {
	request := NewGetCalculatorInfo(uid)
	ic.CalculatorInfo <- request
	return *(<-request.Result)
}

// GetLeadershipStatus реализует соответствующий запрос с блокировкой до получения результата
func (ic *InfoChans) GetLeadershipStatus() bool {
	request := NewGetLeadershipStatus()
	ic.LeadershipStatus <- request
	return <-request.Result
}

// GetWorkerNames возвращает идентификаторы всех воркеров в системе.
type GetWorkerNames struct {
	Result chan []string
}

// GetWorkerInfo возвращает актуальную информацию о воркере
type GetWorkerInfo struct {
	WorkerUID string
	Result    chan *types.WorkerInfo
}

// GetDatasetInfo возвращает актуальную информацию о датасете
type GetDatasetInfo struct {
	DatasetUID string
	Result     chan *types.DatasetInfo
}

// GetCalculatorInfo возвращает актуальную информацию о вычислителе
type GetCalculatorInfo struct {
	CalculatorUID string
	Result        chan *types.CalculatorInfo
}

// GetLeadershipStatus возврашает true если Control Unit, на котором находится scheduler является лидером.
type GetLeadershipStatus struct {
	Result chan bool
}

//
// Initializers
//

// NewGetWorkerNames - функция для удобной инициализации структуры
func NewGetWorkerNames() GetWorkerNames {
	return GetWorkerNames{
		Result: make(chan []string),
	}
}

// NewGetWorkerInfo - функция для удобной инициализации структуры
func NewGetWorkerInfo(uid string) GetWorkerInfo {
	return GetWorkerInfo{
		WorkerUID: uid,
		Result:    make(chan *types.WorkerInfo),
	}
}

// NewGetDatasetInfo - функция для удобной инициализации структуры
func NewGetDatasetInfo(uid string) GetDatasetInfo {
	return GetDatasetInfo{
		DatasetUID: uid,
		Result:     make(chan *types.DatasetInfo),
	}
}

// NewGetCalculatorInfo - функция для удобной инициализации структуры
func NewGetCalculatorInfo(uid string) GetCalculatorInfo {
	return GetCalculatorInfo{
		CalculatorUID: uid,
		Result:        make(chan *types.CalculatorInfo),
	}
}

// NewGetLeadershipStatus - функция для удобной инициализации структуры
func NewGetLeadershipStatus() GetLeadershipStatus {
	return GetLeadershipStatus{
		Result: make(chan bool),
	}
}
