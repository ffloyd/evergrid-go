package datacfg

import "github.com/ffloyd/evergrid-go/global/types"

// ProcessorCfgYAML is a representation of processor segment in data yaml file
type CalculatorCfgYAML struct {
	Name        string
	MFlopsPerMb float64 `yaml:"mflops_per_mb"`
}

// ProcessorCfg is a representation of processor form data config
type CalculatorCfg struct {
	Name        string
	MFlopsPerMb float64
}

// Info returns types.ProcessorInfo representation
func (conf *CalculatorCfg) Info() *types.CalculatorInfo {
	return &types.CalculatorInfo{
		UID:         conf.Name,
		MFlopsPerMb: types.MFlop(conf.MFlopsPerMb),
	}
}
