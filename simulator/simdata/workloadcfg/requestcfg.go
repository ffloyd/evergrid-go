package workloadcfg

import (
	"github.com/ffloyd/evergrid-go/simulator/simdata/workloadcfg/datacfg"
)

// RequestCfgYAML is for parse requests segments in workload YAML config
type RequestCfgYAML struct {
	Type       string
	Dataset    string
	Calculator string
}

// RequestCfg is a representation of request in workload config
type RequestCfg struct {
	Type       string
	Dataset    *datacfg.DatasetCfg
	Calculator *datacfg.CalculatorCfg
}

// Parse transforms RequestCfgYAML to RequestCfg
func (requestYAML RequestCfgYAML) Parse(dataCfg *datacfg.DataCfg) *RequestCfg {
	requestCfg := &RequestCfg{
		Type: requestYAML.Type,
	}

	if len(requestYAML.Dataset) > 0 {
		requestCfg.Dataset = dataCfg.Datasets[requestYAML.Dataset]
	}

	if len(requestYAML.Calculator) > 0 {
		requestCfg.Calculator = dataCfg.Calculators[requestYAML.Calculator]
	}

	return requestCfg
}
