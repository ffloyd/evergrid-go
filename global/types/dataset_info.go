package types

// DatasetInfo - represents current status of dataset
type DatasetInfo struct {
	UID     UID
	Size    MByte
	Workers []UID
}
