package network

// Segment represents a local connected scope of machines. As example - if they are part of same DigitalOcean region.
type Segment struct {
	InnerBandwith Bandwith // bandwith for communication inside this segment

	nodes []Node
}
