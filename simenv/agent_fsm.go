package simenv

import "github.com/Sirupsen/logrus"

// AgentFSM -
type AgentFSM struct {
	chans    AgentChans
	state    AgentState
	stopFlag bool

	logContext *logrus.Entry
}

// NewAgentFSM -
func NewAgentFSM(logContext *logrus.Entry) AgentFSM {
	return AgentFSM{
		chans:      NewAgentChans(),
		state:      StateDone,
		stopFlag:   false,
		logContext: logContext,
	}
}

func (fsm AgentFSM) logf(format string, args ...interface{}) {
	if fsm.logContext != nil {
		fsm.logContext.Debugf(format, args...)
	}
}

func (fsm AgentFSM) logState() {
	if fsm.logContext != nil {
		fsm.logContext.WithField("state", fsm.state).Debug("Agent state changed")
	}
}

// Chans -
func (fsm AgentFSM) Chans() AgentChans {
	return fsm.chans
}

// State -
func (fsm AgentFSM) State() AgentState {
	return fsm.state
}

// StopFlag -
func (fsm AgentFSM) StopFlag() bool {
	return fsm.stopFlag
}

// ToReady -
func (fsm *AgentFSM) ToReady() {
	if fsm.state != StateDone {
		fsm.logContext.Panicf("Wrong state: %v", fsm.state)
	}

	fsm.chans.statusChan <- StateReady
	fsm.state = StateReady

	fsm.logState()

	<-fsm.chans.startWorkChan

	fsm.logf("Agent allowed to work")
}

// ToIdle -
func (fsm *AgentFSM) ToIdle() {
	if fsm.state != StateReady && fsm.state != StateWorking {
		fsm.logContext.Panicf("Wrong state: %v", fsm.state)
	}

	fsm.chans.statusChan <- StateIdle
	fsm.state = StateIdle

	fsm.logState()
}

// ToWorking -
func (fsm *AgentFSM) ToWorking() {
	if fsm.state != StateReady && fsm.state != StateIdle {
		fsm.logContext.Panicf("Wrong state: %v", fsm.state)
	}

	fsm.chans.statusChan <- StateWorking
	fsm.state = StateWorking

	fsm.logState()
}

func (fsm *AgentFSM) toDone() {
	if fsm.state != StateIdle {
		fsm.logContext.Panicf("Wrong state: %v", fsm.state)
	}

	fsm.chans.statusChan <- StateDone
	fsm.state = StateDone

	fsm.logState()
}

// ToDoneChan -
func (fsm *AgentFSM) ToDoneChan() chan Ok {
	result := make(chan Ok)

	go func() {
		<-fsm.chans.finishWorkChan
		fsm.toDone()
		result <- Ok{}
	}()

	return result
}

// SetStopFlag -
func (fsm *AgentFSM) SetStopFlag(value bool) {
	if fsm.stopFlag != value {
		fsm.stopFlag = value
		if fsm.logContext != nil {
			fsm.logContext.WithField("stopFlag", value).Debug("Agent changed stopFlag")
		}
		fsm.chans.stopFlagChan <- value
	}
}
