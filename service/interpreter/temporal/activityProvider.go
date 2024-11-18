package temporal

import (
	"context"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type activityProvider struct{}

func init() {
	interpreter.RegisterActivityProvider(service.BackendTypeTemporal, &activityProvider{})
}

func (a *activityProvider) GetLogger(ctx context.Context) interpreter.UnifiedLogger {
	return activity.GetLogger(ctx)
}

func (a *activityProvider) NewApplicationError(errType string, details interface{}) error {
	return temporal.NewApplicationError("", errType, details)
}

func (a *activityProvider) GetActivityInfo(ctx context.Context) interpreter.ActivityInfo {
	info := activity.GetInfo(ctx)
	return interpreter.ActivityInfo{
		ScheduledTime:   info.ScheduledTime,
		Attempt:         info.Attempt,
		IsLocalActivity: info.IsLocalActivity,
		WorkflowExecution: interpreter.WorkflowExecution{
			ID:    info.WorkflowExecution.ID,
			RunID: info.WorkflowExecution.RunID,
		},
	}
}
