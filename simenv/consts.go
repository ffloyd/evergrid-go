package simenv

/*
Ok - воспомагательная "косметическая" структура

При реализации метода Send интерфейса Agent могут быть случаи, когда в качестве ответа
надо лишь подтвердить получение или успешную обработку запроса. В таком случае можно писать:

	responseChan <- Ok{}
*/
type Ok struct{}

// AgentState - это тип для состояний агентов
type AgentState int

// Набор констант для описания состояний агентов
const (
	StateReady AgentState = iota
	StateWorking
	StateIdle
	StateDone
)

func (state AgentState) String() string {
	switch state {
	case StateReady:
		return "StateReady"
	case StateWorking:
		return "StateWorking"
	case StateIdle:
		return "StateIdle"
	case StateDone:
		return "StateDone"
	default:
		panic("Unknown state value")
	}
}
