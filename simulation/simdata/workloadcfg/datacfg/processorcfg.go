package datacfg

// ProcessorCfgYAML is a representation of processor segment in data yaml file
type ProcessorCfgYAML struct {
	Name        string
	GflopsPerMb float64
}

// ProcessorCfg is a representation of processor form data config
type ProcessorCfg struct {
	ProcessorCfgYAML
}
