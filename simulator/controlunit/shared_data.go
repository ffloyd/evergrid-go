package controlunit

import (
	"sync"

	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simenv"
)

// SharedData используется ControlUnit'ом как упрощения
// моделирования задачи получения глобального состояние системы и
// задачи leader election среди ControlUnit'ов
type SharedData struct {
	Mutex sync.Mutex // для синхронизации доступа к данным

	LeaderElection    sync.Once    // обеспечивает механизм leader election
	LeaderControlUnit simenv.Agent // текущий лидер среди ControlUnit'ов

	Datasets    map[string]*types.DatasetInfo    // актуальная информация о датасетах
	Calculators map[string]*types.CalculatorInfo // актуальная информация о вычислителях
	Workers     map[string]*types.WorkerInfo     // актуальная информация о воркерах
}

// NewSharedData - инициализирует пустой SharedData
func NewSharedData() *SharedData {
	return &SharedData{
		Datasets:    make(map[string]*types.DatasetInfo),
		Calculators: make(map[string]*types.CalculatorInfo),
		Workers:     make(map[string]*types.WorkerInfo),
	}
}
