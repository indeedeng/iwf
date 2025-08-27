package cadence

import (
	"context"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"go.uber.org/cadence"
	"go.uber.org/cadence/activity"
)

type activityProvider struct{}

func init() {
	interfaces.RegisterActivityProvider(service.BackendTypeCadence, &activityProvider{})
}

func (a *activityProvider) NewApplicationError(errType string, details interface{}) error {
	return cadence.NewCustomError(errType, details)
}

func (a *activityProvider) GetLogger(ctx context.Context) interfaces.UnifiedLogger {
	zLogger := activity.GetLogger(ctx)
	return &loggerImpl{
		zlogger: zLogger,
	}
}

func (a *activityProvider) GetActivityInfo(ctx context.Context) interfaces.ActivityInfo {
	info := activity.GetInfo(ctx)
	return interfaces.ActivityInfo{
		ScheduledTime:   info.ScheduledTimestamp,
		Attempt:         info.Attempt + 1, // NOTE increase by one to match Temporal
		IsLocalActivity: false,            // TODO cadence doesn't support this yet
		WorkflowExecution: interfaces.WorkflowExecution{
			ID:    info.WorkflowExecution.ID,
			RunID: info.WorkflowExecution.RunID,
		},
	}
}

func (a *activityProvider) RecordHeartbeat(ctx context.Context, details ...interface{}) {
	activity.RecordHeartbeat(ctx, details...)
}
