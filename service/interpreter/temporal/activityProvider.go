package temporal

import (
	"context"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/interpreter"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type activityProvider struct{}

func init() {
	interpreter.RegisterActivityProvider(service.BackendTypeTemporal, &activityProvider{})
}

func (a *activityProvider) GetLogger(ctx context.Context) interpreter.ActivityLogger {
	return activity.GetLogger(ctx)
}

func (a *activityProvider) NewApplicationError(message, errType string, details ...interface{}) error {
	return temporal.NewApplicationError(message, errType, details...)
}
