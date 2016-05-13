package datacfg

// ProcessorCfgYAML is a representation of processor segment in data yaml file
type ProcessorCfgYAML struct {
	Name        string
	GflopsPerMb float64 `yaml:"gflops_per_mb"`
}

// ProcessorCfg is a representation of processor form data config
type ProcessorCfg struct {
	ProcessorCfgYAML
}
