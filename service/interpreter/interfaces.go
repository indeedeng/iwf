package interpreter

import (
	"context"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
)

type ActivityProvider interface {
	GetLogger(ctx context.Context) ActivityLogger
	NewApplicationError(message, errType string, details ...interface{}) error
}

type ActivityLogger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

func getActivityProviderByType(backendType service.BackendType) ActivityProvider {
	if backendType == service.BackendTypeTemporal {
		return temporal.DefaultActivityProvider
	}
	panic("not supported yet: " + backendType)
}
