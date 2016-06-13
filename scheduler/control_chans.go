package scheduler

/*
Done - тип для сообщения не содержащего информации.

Теперь вместо

	ch := make(chan struct{})
	// some code
	ch <- struct{}{}

можно писать более осмысленное

	ch := make(chan Done)
	ch <- Done{}
*/
type Done struct{}

/*
ControlChans это каналы, через которые планировщик управляет внешней средой путем отправки сообщений.

Использование этой группы каналов должно быть синхронным: после отправки сообщения
по любому из трех первых каналов надо дождаться подтверждения его обработки:

	x.UploadDataset <- command
	<-x.Done // wait for confirmation

Можно изменить эту структуру и структуры команд так, чтобы их использование было
асинхронным, но это усложнит код ценой незначительного на ранних этапах роста производительности.
*/
type ControlChans struct {
	UploadDataset   chan DoUploadDataset
	BuildCalculator chan DoBuildCalculator
	RunCalculator   chan DoRunCalculator

	Done chan Done
}

// NewControlChans инициализатор для ControlChans
func NewControlChans() ControlChans {
	return ControlChans{
		UploadDataset:   make(chan DoUploadDataset),
		BuildCalculator: make(chan DoBuildCalculator),
		RunCalculator:   make(chan DoRunCalculator),

		Done: make(chan Done),
	}
}

// DoUploadDataset команда на добавление загрузки датасета в очередь задач воркера.
type DoUploadDataset struct {
	Dataset string
	Worker  string
}

// DoBuildCalculator команда на добавление сборки вычислителя в очередь задач воркера.
type DoBuildCalculator struct {
	Calculator string
	Worker     string
}

// DoRunCalculator команда на добавление запуска вычислителя с указанным датасетом в очередь задач воркера.
type DoRunCalculator struct {
	Calculator string
	Dataset    string
	Worker     string
}
