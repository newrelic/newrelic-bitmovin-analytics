package integration

import (
	"newrelic/multienv/integration/bitmovin"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/model"
)

// InitRecv Integration Receiver Initializer
func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	return bitmovin.InitRecv(pipeConfig)
}

// InitProc Integration Processor Initializer
func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	return bitmovin.InitProc(pipeConfig)
}

// Proc Integration Processor
func Proc(data any) []model.MeltModel {
	return bitmovin.Proc(data)
}
