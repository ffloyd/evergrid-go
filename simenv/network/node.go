package network

// Node represents particular machine
type Node struct {
	OuterBandwith Bandwith

	ActiveTransfers []*Transfer
}
