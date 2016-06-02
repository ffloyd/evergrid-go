package agent

import "github.com/ffloyd/evergrid-go/global/types"

type globalStateConstructor struct {
	*types.GlobalState
	env *Environ
}

func (monitor *Monitor) workerGlobalState() {
	nextState := func() *types.GlobalState {
		next := globalStateConstructor{
			GlobalState: types.NewGlobalState(),
		}
		next.Load(monitor.env)
		return next.GlobalState
	}

	innerChan := make(chan *types.GlobalState)

	for {
		monitor.sensorChans.GlobalState <- innerChan
		innerChan <- nextState()
	}
}

func (state *globalStateConstructor) Load(env *Environ) {
	state.env = env

	// load ControlUnits
	for uid, cu := range state.env.ControlUnits {
		cuInfo := &types.ControlUnitInfo{
			UID:     types.UID(uid),
			Workers: make([]*types.WorkerInfo, len(cu.workers)),
		}

		for i, worker := range cu.workers {
			workerInfo := &types.WorkerInfo{
				UID: types.UID(worker.Name()),
			}
			cuInfo.Workers[i] = workerInfo
		}

		state.ControlUnits[types.UID(uid)] = cuInfo
	}

	// load workers to toplevel
	for _, cuInfo := range state.ControlUnits {
		for _, workerInfo := range cuInfo.Workers {
			state.Workers[workerInfo.UID] = workerInfo
		}
	}

	// load Datasets, Processors and Jobs from workers
	for _, workerInfo := range state.Workers {
		for _, datasetInfo := range workerInfo.Datasets {
			if state.Datasets[datasetInfo.UID] == nil {
				state.Datasets[datasetInfo.UID] = datasetInfo
			}
		}

		for _, processorInfo := range workerInfo.Processors {
			if state.Processors[processorInfo.UID] == nil {
				state.Processors[processorInfo.UID] = processorInfo
			}
		}

		if workerInfo.CurrentJob != nil {
			state.ActiveJobs[workerInfo.CurrentJob.UID] = workerInfo.CurrentJob
		}
	}
}
