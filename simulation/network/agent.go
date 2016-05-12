package network

// Agent interface defines methods which essential for agent persistance in Node struct
type Agent interface {
	Name() string
}
