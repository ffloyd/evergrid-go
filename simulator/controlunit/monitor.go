package controlunit

import "github.com/ffloyd/evergrid-go/scheduler"

type monitor struct {
	cu *ControlUnit
}

func newMonitor(cu *ControlUnit) monitor {
	return monitor{
		cu: cu,
	}
}

func (mon *monitor) Run() {
	go mon.work()
}

func (mon *monitor) work() {
	chans := mon.cu.scheduler.InfoChans()
	for {
		select {
		case request := <-chans.WorkerNames:
			go mon.processWorkerNames(request)
		case request := <-chans.WorkerInfo:
			go mon.processWorkerInfo(request)
		case request := <-chans.DatasetInfo:
			go mon.processDatasetInfo(request)
		case request := <-chans.CalculatorInfo:
			go mon.processCalculatorInfo(request)
		case request := <-chans.LeadershipStatus:
			go mon.processLeadershipStatus(request)
		}
	}
}

func (mon *monitor) processWorkerNames(request scheduler.GetWorkerNames) {
	mon.cu.sharedData.Mutex.Lock()
	names := make([]string, len(mon.cu.sharedData.Workers))

	i := 0
	for name := range mon.cu.sharedData.Workers {
		names[i] = name
		i++
	}
	mon.cu.sharedData.Mutex.Unlock()

	request.Result <- names
}

func (mon *monitor) processWorkerInfo(request scheduler.GetWorkerInfo) {
	mon.cu.sharedData.Mutex.Lock()
	request.Result <- mon.cu.sharedData.Workers[request.WorkerUID]
	mon.cu.sharedData.Mutex.Unlock()
}

func (mon *monitor) processDatasetInfo(request scheduler.GetDatasetInfo) {
	mon.cu.sharedData.Mutex.Lock()
	request.Result <- mon.cu.sharedData.Datasets[request.DatasetUID]
	mon.cu.sharedData.Mutex.Unlock()
}

func (mon *monitor) processCalculatorInfo(request scheduler.GetCalculatorInfo) {
	mon.cu.sharedData.Mutex.Lock()
	request.Result <- mon.cu.sharedData.Calculators[request.CalculatorUID]
	mon.cu.sharedData.Mutex.Unlock()
}

func (mon *monitor) processLeadershipStatus(request scheduler.GetLeadershipStatus) {
	mon.cu.sharedData.Mutex.Lock()
	request.Result <- (mon.cu.sharedData.LeaderControlUnit.Name() == mon.cu.Name())
	mon.cu.sharedData.Mutex.Unlock()
}
