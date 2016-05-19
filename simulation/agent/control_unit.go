package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/scheduler"
	"github.com/ffloyd/evergrid-go/simulation/network"
	"github.com/ffloyd/evergrid-go/simulation/simdata/networkcfg"
	"github.com/ffloyd/evergrid-go/simulation/simdata/workloadcfg"
)

// ControlUnit is a representation of control unit app
type ControlUnit struct {
	Base

	incomingRequests chan *workloadcfg.RequestCfg
	noMoreRequests   chan bool
	workers          []*Worker

	scheduler *scheduler.Scheduler
}

// NewControlUnit creates a new control unit
func NewControlUnit(config *networkcfg.AgentCfg, net *network.Network, env *Environ) *ControlUnit {
	unit := &ControlUnit{
		Base:             *NewBase(config, net, env),
		incomingRequests: make(chan *workloadcfg.RequestCfg),
		noMoreRequests:   make(chan bool),
	}
	env.ControlUnits[unit.Name()] = unit

	log.WithFields(log.Fields{
		"name": unit.Name(),
		"node": unit.Node(),
	}).Info("Control Unit agent initialized")
	return unit
}

func (unit ControlUnit) processRequest(request *workloadcfg.RequestCfg) {

}

func (unit *ControlUnit) startScheduler() {
	log.WithFields(log.Fields{
		"control_unit": unit,
		"algorithm":    "FIFO",
	}).Info("Starting scheduler")

	unit.scheduler = scheduler.New(scheduler.FIFO)
	go unit.scheduler.Run()

	<-unit.scheduler.Chans.Alive

	log.WithFields(log.Fields{
		"control_unit": unit,
		"algorithm":    "FIFO",
	}).Info("Scheduler started")
}

func (unit ControlUnit) run() {
	unit.startScheduler()

	for {
		unit.startTick()

	SelectLoop:
		for {
			select {
			case request := <-unit.incomingRequests:
				log.WithFields(log.Fields{
					"control_unit": unit,
					"tick":         unit.tick,
					"type":         request.Type,
				}).Info("Control unit received request")
				unit.processRequest(request)
			case <-unit.noMoreRequests:
				unit.finishTick()
				break SelectLoop
			}
		}
	}
}

// Run is implementation of agent.Runner iface
func (unit ControlUnit) Run() *TickerChans {
	go unit.run()
	return unit.tickerChans
}
