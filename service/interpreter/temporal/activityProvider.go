package temporal

import (
	"context"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type activityProvider struct{}

func init() {
	interfaces.RegisterActivityProvider(service.BackendTypeTemporal, &activityProvider{})
}

func (a *activityProvider) GetLogger(ctx context.Context) interfaces.UnifiedLogger {
	return activity.GetLogger(ctx)
}

func (a *activityProvider) NewApplicationError(errType string, details interface{}) error {
	return temporal.NewApplicationError("", errType, details)
}

func (a *activityProvider) GetActivityInfo(ctx context.Context) interfaces.ActivityInfo {
	info := activity.GetInfo(ctx)
	return interfaces.ActivityInfo{
		ScheduledTime:   info.ScheduledTime,
		Attempt:         info.Attempt,
		IsLocalActivity: info.IsLocalActivity,
		WorkflowExecution: interfaces.WorkflowExecution{
			ID:    info.WorkflowExecution.ID,
			RunID: info.WorkflowExecution.RunID,
		},
	}
}
