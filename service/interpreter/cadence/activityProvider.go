package cadence

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter"
	"go.uber.org/cadence/activity"
)

type activityProvider struct{}

func init() {
	interpreter.RegisterActivityProvider(service.BackendTypeCadence, &activityProvider{})
}

func (a *activityProvider) NewApplicationError(message, errType string, details ...interface{}) error {
	return fmt.Errorf("application error: error type: %v, message %v, details %v", errType, message, details)
}

func (a *activityProvider) GetLogger(ctx context.Context) interpreter.UnifiedLogger {
	zLogger := activity.GetLogger(ctx)
	return &loggerImpl{
		zlogger: zLogger,
	}
}
