package interpreter

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

func (wc *WorkflowConfiger) SetIfPresent(config iwfidl.WorkflowConfig) {
	if config.DisableSystemSearchAttribute != nil {
		wc.config.DisableSystemSearchAttribute = config.DisableSystemSearchAttribute
	}
	if config.ContinueAsNewPageSizeInBytes != nil {
		wc.config.ContinueAsNewPageSizeInBytes = config.ContinueAsNewPageSizeInBytes
	}
	if config.ContinueAsNewThreshold != nil {
		wc.config.ContinueAsNewThreshold = config.ContinueAsNewThreshold
	}
}
