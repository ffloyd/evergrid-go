package types

/*
CalculatorInfo это описание известной информации о вычислителе. В изначальной спецификации Evergrid предполагается,
что это некоторый Doceker-контейнер с зашитым в нем алгоритмом обработки данных. Но для внутренного
наименования было выбрано слово Calculator, т. к. неочевидно, что container должен что-то считать.

Конечное местоположение Calculator'а - некоторый Worker. Перед использованием на Worker'е Calculator
должен быть собран (build).
*/
type CalculatorInfo struct {
	UID               string   // Уникальный идентификатор
	MFlopsPerMb       MFlop    // Количество MFlop расходуемых на один мегабайт данных
	Workers           []string // Список имен Worker'ов, на которых _уже_ собран этот Calculator
	EnqueuedOnWorkers []string // Список имен Worker'ов, на которых _будет_ собран этот Calculator
}
