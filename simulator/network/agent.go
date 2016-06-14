package network

// Agent - в рамках данного пакета нас итересуют только имена агентов
type Agent interface {
	Name() string
}
