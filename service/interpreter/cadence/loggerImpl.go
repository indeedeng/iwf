package cadence

import (
	"fmt"
	"go.uber.org/zap"
)

type loggerImpl struct {
	zlogger *zap.Logger
}

func buildZapFields(keyvals []interface{}) []zap.Field {
	var fields []zap.Field
	if len(keyvals)%2 != 0 {
		panic(fmt.Sprintf("invalid length for keyvals %v", len(keyvals)))
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

func (a *loggerImpl) Debug(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Debug(msg, fields...)
}

func (a *loggerImpl) Info(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Info(msg, fields...)
}

func (a *loggerImpl) Warn(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Warn(msg, fields...)
}

func (a *loggerImpl) Error(msg string, keyvals ...interface{}) {
	fields := buildZapFields(keyvals)
	a.zlogger.Error(msg, fields...)
}
