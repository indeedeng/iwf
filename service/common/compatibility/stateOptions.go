package compatibility

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/ptr"
)

func GetStartApiTimeoutSeconds(stateOptions *iwfidl.WorkflowStateOptions) int32 {
	if stateOptions == nil {
		return 0
	}
	if stateOptions.HasStartApiTimeoutSeconds() {
		return stateOptions.GetStartApiTimeoutSeconds()
	}
	return stateOptions.GetWaitUntilApiTimeoutSeconds()
}

func GetDecideApiTimeoutSeconds(stateOptions *iwfidl.WorkflowStateOptions) int32 {
	if stateOptions == nil {
		return 0
	}
	if stateOptions.HasDecideApiTimeoutSeconds() {
		return stateOptions.GetDecideApiTimeoutSeconds()
	}
	return stateOptions.GetExecuteApiTimeoutSeconds()
}

func GetStartApiRetryPolicy(stateOptions *iwfidl.WorkflowStateOptions) *iwfidl.RetryPolicy {
	if stateOptions == nil {
		return nil
	}
	if stateOptions.HasStartApiRetryPolicy() {
		return stateOptions.StartApiRetryPolicy
	}
	return stateOptions.WaitUntilApiRetryPolicy
}

func GetDecideApiRetryPolicy(stateOptions *iwfidl.WorkflowStateOptions) *iwfidl.RetryPolicy {
	if stateOptions == nil {
		return nil
	}
	if stateOptions.HasDecideApiRetryPolicy() {
		return stateOptions.DecideApiRetryPolicy
	}
	return stateOptions.ExecuteApiRetryPolicy
}

func GetWaitUntilApiDataObjectsLoadingPolicy(stateOptions *iwfidl.WorkflowStateOptions) *iwfidl.PersistenceLoadingPolicy {
	if stateOptions == nil {
		return nil
	}

	if stateOptions.HasWaitUntilApiDataAttributesLoadingPolicy() {
		return stateOptions.WaitUntilApiDataAttributesLoadingPolicy
	}

	if stateOptions.HasDataAttributesLoadingPolicy() {
		return stateOptions.DataAttributesLoadingPolicy
	}

	return stateOptions.DataObjectsLoadingPolicy
}

func GetExecuteApiDataObjectsLoadingPolicy(stateOptions *iwfidl.WorkflowStateOptions) *iwfidl.PersistenceLoadingPolicy {
	if stateOptions == nil {
		return nil
	}

	if stateOptions.HasExecuteApiDataAttributesLoadingPolicy() {
		return stateOptions.ExecuteApiDataAttributesLoadingPolicy
	}

	if stateOptions.HasDataAttributesLoadingPolicy() {
		return stateOptions.DataAttributesLoadingPolicy
	}

	return stateOptions.DataObjectsLoadingPolicy
}

func GetWaitUntilApiSearchAttributesLoadingPolicy(stateOptions *iwfidl.WorkflowStateOptions) *iwfidl.PersistenceLoadingPolicy {
	if stateOptions == nil {
		return nil
	}

	if stateOptions.HasWaitUntilApiSearchAttributesLoadingPolicy() {
		return stateOptions.WaitUntilApiSearchAttributesLoadingPolicy

	}

	return stateOptions.SearchAttributesLoadingPolicy
}

func GetExecuteApiSearchAttributesLoadingPolicy(stateOptions *iwfidl.WorkflowStateOptions) *iwfidl.PersistenceLoadingPolicy {
	if stateOptions == nil {
		return nil
	}

	if stateOptions.HasExecuteApiSearchAttributesLoadingPolicy() {
		return stateOptions.ExecuteApiSearchAttributesLoadingPolicy

	}

	return stateOptions.SearchAttributesLoadingPolicy
}

func GetStartApiFailurePolicy(stateOptions *iwfidl.WorkflowStateOptions) *iwfidl.StartApiFailurePolicy {
	if stateOptions.HasStartApiFailurePolicy() {
		return stateOptions.StartApiFailurePolicy
	}
	if stateOptions.HasWaitUntilApiFailurePolicy() {
		newPolicy := stateOptions.GetWaitUntilApiFailurePolicy()
		switch newPolicy {
		case iwfidl.FAIL_WORKFLOW_ON_FAILURE:
			return ptr.Any(iwfidl.FAIL_WORKFLOW_ON_START_API_FAILURE)
		case iwfidl.PROCEED_ON_FAILURE:
			return ptr.Any(iwfidl.PROCEED_TO_DECIDE_ON_START_API_FAILURE)
		default:
			panic("invalid policy to convert " + string(newPolicy))
		}
	}
	return nil
}

func GetSkipWaitUntilApi(stateOptions *iwfidl.WorkflowStateOptions) bool {
	if stateOptions.HasSkipStartApi() {
		return stateOptions.GetSkipStartApi()
	}
	return stateOptions.GetSkipWaitUntil()
}
