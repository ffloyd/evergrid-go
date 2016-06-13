package types

/*
DatasetInfo это описание известной информации о датасете.

Датасет - это некоторый объем данных, которые являются входными
параметрами для вычислителей (Calculators).
*/
type DatasetInfo struct {
	UID               string   // Уникальный идентификатор
	Size              MByte    // Размер датасета в мегабайтах
	Workers           []string // Воркеры, на которые _уже_ загружен датасет
	EnqueuedOnWorkers []string // Воркеры, на которые _будет_ загружен датасет
}
