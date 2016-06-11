package datacfg

import "github.com/ffloyd/evergrid-go/global/types"

// ProcessorCfgYAML is a representation of processor segment in data yaml file
type ProcessorCfgYAML struct {
	Name        string
	MFlopsPerMb float64 `yaml:"mflops_per_mb"`
}

// ProcessorCfg is a representation of processor form data config
type ProcessorCfg struct {
	ProcessorCfgYAML
}

// Info returns types.ProcessorInfo representation
func (conf *ProcessorCfg) Info() *types.ProcessorInfo {
	return &types.ProcessorInfo{
		UID:         types.UID(conf.Name),
		MFlopsPerMb: types.MFlop(conf.MFlopsPerMb),
	}
}
