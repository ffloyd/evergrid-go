package simenv

import "sync"

/*
SimEnv - это ядро симуляции.

Именно через эту структуру происходит синхронзация работы всех агентов.

Симуляция разбита на дискретные "тики" в рамках которой агенты меняют свои состояния.
*/
type SimEnv struct {
	tick          int
	agents        map[string]Agent
	tickBroadcast []chan int
	inProgress    bool
}

// New - инициализатор
func New() *SimEnv {
	return &SimEnv{
		agents: make(map[string]Agent),
	}
}

// Add добавляет агентов в симуляцию
func (simenv *SimEnv) Add(agents ...Agent) {
	for _, agent := range agents {
		simenv.agents[agent.Name()] = agent
	}
}

// Find ищет агента по имени
func (simenv SimEnv) Find(agentName string) Agent {
	return simenv.agents[agentName]
}

// Run запускает всех агентов и симуляцию
func (simenv *SimEnv) Run() error {
	group := runAgentGroup(simenv)

	for {
		simenv.tick++
		for _, tickChan := range simenv.tickBroadcast {
			tickChan <- simenv.tick
		}

		group.WaitForState(StateReady)
		group.StartWork()
		group.WaitForState(StateIdle)
		group.FinishWork()

		if group.StopFlag() {
			break
		}
	}

	return nil
}

// CurrentTick возвращает структуру, которая всегда содержит текущий тик симуляции
func (simenv *SimEnv) CurrentTick() *CurrentTick {
	ct := &CurrentTick{simenv.tick, new(sync.Mutex)}
	channel := make(chan int)

	simenv.tickBroadcast = append(simenv.tickBroadcast, channel)

	go ct.connect(channel)

	return ct
}
