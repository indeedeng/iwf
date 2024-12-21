package config

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
)

type WorkflowConfiger struct {
	config iwfidl.WorkflowConfig
}

func NewWorkflowConfiger(config iwfidl.WorkflowConfig) *WorkflowConfiger {
	return &WorkflowConfiger{
		config: config,
	}
}

func (wc *WorkflowConfiger) Get() iwfidl.WorkflowConfig {
	return wc.config
}

func (wc *WorkflowConfiger) ShouldOptimizeActivity() bool {
	return wc.config.GetOptimizeActivity()
}

func (wc *WorkflowConfiger) UpdateByAPI(config iwfidl.WorkflowConfig) {
	if config.DisableSystemSearchAttribute != nil {
		wc.config.DisableSystemSearchAttribute = config.DisableSystemSearchAttribute
	}
	if config.ExecutingStateIdMode != nil {
		wc.config.ExecutingStateIdMode = config.ExecutingStateIdMode
	}
	if config.ContinueAsNewPageSizeInBytes != nil {
		wc.config.ContinueAsNewPageSizeInBytes = config.ContinueAsNewPageSizeInBytes
	}
	if config.ContinueAsNewThreshold != nil {
		wc.config.ContinueAsNewThreshold = config.ContinueAsNewThreshold
	}
	if config.OptimizeActivity != nil {
		wc.config.OptimizeActivity = config.OptimizeActivity
	}
}
