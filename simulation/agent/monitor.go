package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/scheduler"
)

// Monitor - part of ControlUnit. Processing all sensor's requests from scheduler.
type Monitor struct {
	agentName   string
	env         *Environ
	sensorChans *scheduler.SensorChans
}

func startMonitor(scheduler *scheduler.Scheduler, env *Environ, agentName string) *Monitor {
	monitor := &Monitor{
		agentName:   agentName,
		env:         env,
		sensorChans: scheduler.Chans.Sensors,
	}

	go monitor.workerIsLeader()
	go monitor.workerGlobalState()

	log.WithFields(log.Fields{
		"agent": agentName,
	}).Info("Monitor started on Control Unit")

	return monitor
}

func (monitor *Monitor) workerIsLeader() {
	for {
		monitor.sensorChans.IsLeader <- (monitor.env.LeaderControlUnit().Name() == monitor.agentName)
	}
}
