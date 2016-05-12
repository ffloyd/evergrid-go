package agent

import (
	"fmt"

	"github.com/ffloyd/evergrid-go/simulation/network"
)

// Agent interface must be implemented for every agent
type Agent interface {
	fmt.Stringer
	network.Agent
	Run() *Chans
	Node() *network.Node
}
