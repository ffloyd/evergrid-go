package gendata

/*
Config это описание конфигурации для генератора

В данный момент все случайные величины генерируются равномерным распределением
*/
type Config struct {
	DestDir string // корневая директория для файлов
	Name    string // имя генерируемого сценария работы

	DatsetsCount   int // количество различных датасетов
	MinDatasetSize int // минимальный размер в гигабайтах
	MaxDatasetSize int // максимальный размер в гигабайтах

	CalculatorsCount        int // количество различных вычислителей
	MinCalculatorComplexity int // минимальная сложность в мегафлопсах на мегабайт данных
	MaxCalculatorComplexity int // максимальная сложность в мегафлопсах на мегабайт данных

	CalculatorRuns int     // количество запусков экспериметов
	RunProbability float64 // веротность того, в рамках текщего тика произойдет запуск эксперимента

	NetworkSegments   int     // количество сегментов в сети
	MinNodesInSegment int     // минимальное количество машин в сегменте
	MaxNodesInSegment int     // максимальное количество машин в сегменте
	MinNodeSpeed      int     // минимальная производительность воркера в мегафлопсах
	MaxNodeSpeed      int     // максимальная производительность воркера в мегафлопсах
	MinDiskSize       int     // минимальный размер диска в гигабайтах
	MaxDiskSize       int     // максимальный размер диска в гигабайтах
	MinPricePerTick   float64 // минимальная цена за 1 минуту работы в условных единицах
	MaxPricePerTick   float64 // максимальная цена за 1 минуту работы в условных единицах
}
