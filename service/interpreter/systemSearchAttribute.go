package interpreter

import "github.com/indeedeng/iwf/service"

func upsertInitialSystemSearchAttributes(ctx UnifiedContext, input service.InterpreterWorkflowInput, globalVersioner *GlobalVersioner, provider WorkflowProvider) error {
	var err error
	if globalVersioner.IsAfterVersionOfUsingGlobalVersioning() {
		err = globalVersioner.UpsertGlobalVersionSearchAttribute()
		if err != nil {
			return err
		}
	}

	if !input.Config.GetDisableSystemSearchAttribute() {
		if !globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() {
			// we have stopped upsert here in new versions, because it's done in start workflow request
			err = provider.UpsertSearchAttributes(ctx, map[string]interface{}{
				service.SearchAttributeIwfWorkflowType: input.IwfWorkflowType,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
