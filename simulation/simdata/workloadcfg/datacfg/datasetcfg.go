package datacfg

// DatasetCfgYAML is a representation of dataset segment in data yaml config
type DatasetCfgYAML struct {
	Name string
	Size int
}

// DatasetCfg is a representation of dataset in data config
type DatasetCfg struct {
	DatasetCfgYAML
}
