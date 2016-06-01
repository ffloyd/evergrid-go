package datacfg

// ProcessorCfgYAML is a representation of processor segment in data yaml file
type ProcessorCfgYAML struct {
	Name        string
	MFlopsPerMb float64 `yaml:"mflops_per_mb"`
}

// ProcessorCfg is a representation of processor form data config
type ProcessorCfg struct {
	ProcessorCfgYAML
}
