package simulator

import (
	"github.com/ffloyd/evergrid-go/simenv"
	"github.com/ffloyd/evergrid-go/simulator/network"
	"github.com/ffloyd/evergrid-go/simulator/simdata"
)

// Simulator -
type Simulator struct {
	simData *simdata.SimData

	network *network.Network
	simenv  *simenv.SimEnv
}

// New -
func New(simdataFilename string) *Simulator {
	sim := &Simulator{
		simData: simdata.Load(simdataFilename),
		simenv:  simenv.New(),
	}

	sim.network = network.New(sim.simData.Network)

	return sim
}

// Run -
func (sim *Simulator) Run() {
	sim.simenv.Run()
}
