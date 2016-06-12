package controlunit

import (
	"sync"

	"github.com/ffloyd/evergrid-go/global/types"
	"github.com/ffloyd/evergrid-go/simenv"
)

// SharedData -
type SharedData struct {
	Mutex sync.Mutex

	LeaderElection    sync.Once
	LeaderControlUnit simenv.Agent

	Datasets    map[string]*types.DatasetInfo
	Calculators map[string]*types.CalculatorInfo
	Workers     map[string]*types.WorkerInfo
}

// NewSharedData -
func NewSharedData() *SharedData {
	return &SharedData{
		Datasets:    make(map[string]*types.DatasetInfo),
		Calculators: make(map[string]*types.CalculatorInfo),
		Workers:     make(map[string]*types.WorkerInfo),
	}
}
