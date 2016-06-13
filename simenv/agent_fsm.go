package simenv

import (
	"sync"

	"github.com/Sirupsen/logrus"
)

/*
AgentFSM (Agent's finite state machine)- это структура контроллирующая состояние агента.

Тривиальный жизненный цикл агента выглядит следующим образом:

	// инициализация агента

	for {
		// действия необходимые перед началом тика
		fsm.ToReady()
		fsm.ToWorking()
		// полезная нагрузка, общение с другими агентами
		fsm.ToIdle()
		if noMoreWorkInFuture {
			fsm.SetStopFlag(true)
		}
		<-fsm.ToDoneChan()

	}

Текущий статус агента определяется его состоянием (Done, Ready, Working, Idle) и значением флага stopFlag.

Тик завершается, когда все агенты находятся в состоянии Idle. Симуляция завершается, когда у всех агентов stopFlag равен true.
*/
type AgentFSM struct {
	chans AgentChans
	state AgentState

	stopFlag      bool
	stopFlagMutex sync.Mutex

	logContext *logrus.Entry
}

/*
NewAgentFSM - инициализатор AgentFSM

logContext - это контекст для логов. На уровне debug логируется смена статусов.
*/
func NewAgentFSM(logContext *logrus.Entry) *AgentFSM {
	return &AgentFSM{
		chans:      newAgentChans(),
		state:      StateDone,
		stopFlag:   false,
		logContext: logContext,
	}
}

func (fsm *AgentFSM) logf(format string, args ...interface{}) {
	if fsm.logContext != nil {
		fsm.logContext.Debugf(format, args...)
	}
}

func (fsm *AgentFSM) logState() {
	if fsm.logContext != nil {
		fsm.logContext.WithField("state", fsm.state).Debug("Agent state changed")
	}
}

// Chans возвращает каналы для общения с агентом. Это нужно для реализации метода Run интерфейса Agent.
func (fsm *AgentFSM) Chans() AgentChans {
	return fsm.chans
}

// State возвращает текущее состояние агента.
func (fsm *AgentFSM) State() AgentState {
	return fsm.state
}

// StopFlag возвращает текущее значение флага stopFlag
func (fsm *AgentFSM) StopFlag() bool {
	fsm.stopFlagMutex.Lock()
	defer fsm.stopFlagMutex.Unlock()
	return fsm.stopFlag
}

/*
ToReady переводит агент в состояние Ready и ожидает, когда остальные агенты дотигнут этого состояния.

Переход в состояние Ready возможен только из состояния Done.
*/
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

/*
ToIdle переводит агент в состояние Idle.

Переход возможен только из состояний Ready и Working.
*/
func (fsm *AgentFSM) ToIdle() {
	if fsm.state != StateReady && fsm.state != StateWorking {
		fsm.logContext.Panicf("Wrong state: %v", fsm.state)
	}

	fsm.chans.statusChan <- StateIdle
	fsm.state = StateIdle

	fsm.logState()
}

/*
ToWorking переводит агент в состояние Working.

Переход возможен только из состояний Ready и Idle
*/
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

/*
ToDoneChan возвращает канал, по которому будет отправлен Ok{} когда агент перейдет в состояние Done.

Переход в состояние Done происходит автоматически в тот момент, когда все агенты будут находиться в состоянии Idle.
*/
func (fsm *AgentFSM) ToDoneChan() chan Ok {
	result := make(chan Ok)

	go func() {
		<-fsm.chans.finishWorkChan
		fsm.toDone()
		result <- Ok{}
	}()

	return result
}

/*
SetStopFlag устанавливает значение stopFlag.

Когда stopFlag всех агентов будет равен true - симуляция завершится по достижению конца тика.
*/
func (fsm *AgentFSM) SetStopFlag(value bool) {
	fsm.stopFlagMutex.Lock()
	if fsm.stopFlag != value {
		fsm.stopFlag = value
		if fsm.logContext != nil {
			fsm.logContext.WithField("stopFlag", value).Debug("Agent changed stopFlag")
		}
		fsm.chans.stopFlagChan <- value
	}
	fsm.stopFlagMutex.Unlock()
}
