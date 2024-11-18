package cadence

import (
	"context"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter"
	"go.uber.org/cadence"
	"go.uber.org/cadence/activity"
)

type activityProvider struct{}

func init() {
	interpreter.RegisterActivityProvider(service.BackendTypeCadence, &activityProvider{})
}

func (a *activityProvider) NewApplicationError(errType string, details interface{}) error {
	return cadence.NewCustomError(errType, details)
}

func (a *activityProvider) GetLogger(ctx context.Context) interpreter.UnifiedLogger {
	zLogger := activity.GetLogger(ctx)
	return &loggerImpl{
		zlogger: zLogger,
	}
}

func (a *activityProvider) GetActivityInfo(ctx context.Context) interpreter.ActivityInfo {
	info := activity.GetInfo(ctx)
	return interpreter.ActivityInfo{
		ScheduledTime:   info.ScheduledTimestamp,
		Attempt:         info.Attempt + 1, // NOTE increase by one to match Temporal
		IsLocalActivity: false,            // TODO cadence doesn't support this yet
		WorkflowExecution: interpreter.WorkflowExecution{
			ID:    info.WorkflowExecution.ID,
			RunID: info.WorkflowExecution.RunID,
		},
	}
}
