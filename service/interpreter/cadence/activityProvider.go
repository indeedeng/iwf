package cadence

import (
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/interpreter"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

type activityProvider struct{}

func init() {
	interpreter.RegisterActivityProvider(service.BackendTypeCadence, &activityProvider{})
}

type activityLogger struct {
	zlogger *zap.Logger
}

func buildZapFields(keyvals []interface{}) []zap.Field {
	var fields []zap.Field
	if len(keyvals)%2 != 0 {
		panic(fmt.Sprintf("invalid lenght for keyvals %v", len(keyvals)))
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			panic(fmt.Sprintf("invalid keyvals for logging at %d %v ", i, keyvals))
		}
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}
	return fields
}

func (a *activityLogger) Debug(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Debug(msg, fields...)
}

func (a *activityLogger) Info(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Info(msg, fields...)
}

func (a *activityLogger) Warn(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Warn(msg, fields...)
}

func (a *activityLogger) Error(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Error(msg, fields...)
}

func (a *activityProvider) GetLogger(ctx context.Context) interpreter.ActivityLogger {
	zLogger := activity.GetLogger(ctx)
	return &activityLogger{
		zlogger: zLogger,
	}
}

func (a *activityProvider) NewApplicationError(message, errType string, details ...interface{}) error {
	return fmt.Errorf("application error: error type: %v, message %v, details %v", errType, message, details)
}
