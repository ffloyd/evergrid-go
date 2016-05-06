package network

// Transfer represents transferring data between two nodes
type Transfer struct {
	SrcNode  *Node
	DestNode *Node

	Speed       int // in megabits/sec
	DataSize    int // data size in megabytes
	Transferred int // megabytes transferred
}
