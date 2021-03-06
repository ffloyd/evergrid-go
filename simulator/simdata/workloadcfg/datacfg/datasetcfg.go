package datacfg

import "github.com/ffloyd/evergrid-go/global/types"

// DatasetCfgYAML is a representation of dataset segment in data yaml config
type DatasetCfgYAML struct {
	Name string
	Size int
}

// DatasetCfg is a representation of dataset in data config
type DatasetCfg struct {
	Name string
	Size int
}

// Info returns types.DatasetInfo representation
func (conf *DatasetCfg) Info() *types.DatasetInfo {
	return &types.DatasetInfo{
		UID:  conf.Name,
		Size: types.MByte(conf.Size * 1024),
	}
}
