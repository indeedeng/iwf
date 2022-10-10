package interpreter

import (
	"context"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"time"
)

type ActivityProvider interface {
	GetLogger(ctx context.Context) ActivityLogger
	NewApplicationError(message, errType string, details ...interface{}) error
}

func getActivityProviderByType(backendType service.BackendType) ActivityProvider {
	if backendType == service.BackendTypeTemporal {
		return temporal.DefaultActivityProvider
	}
	panic("not supported yet: " + backendType)
}

type ActivityLogger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

// WorkflowExecution details.
type WorkflowExecution struct {
	ID    string
	RunID string
}

// WorkflowInfo information about currently executing workflow
type WorkflowInfo struct {
	WorkflowExecution WorkflowExecution
	WorkflowStartTime time.Time
}

type ActivityOptions struct {
	StartToCloseTimeout time.Duration
}

type UnifiedContext interface {
	GetContext() interface{}
}

type contextHolder struct {
	ctx interface{}
}

func (c *contextHolder) GetContext() interface{} {
	return c.ctx
}

func NewUnifiedContext(ctx interface{}) UnifiedContext {
	return &contextHolder{
		ctx: ctx,
	}
}

type WorkflowProvider interface {
	NewApplicationError(message, errType string, details ...interface{}) error
	GetWorkflowInfo(ctx UnifiedContext) WorkflowInfo
	UpsertSearchAttributes(ctx UnifiedContext, attributes map[string]interface{}) error
	SetQueryHandler(ctx UnifiedContext, queryType string, handler interface{}) error
	ExtendContextWithValue(parent interface{}, key string, val interface{}) UnifiedContext
	GoNamed(ctx UnifiedContext, name string, f func(ctx UnifiedContext))
	Await(ctx UnifiedContext, condition func() bool) error
	WithActivityOptions(ctx UnifiedContext, options ActivityOptions) UnifiedContext
	ExecuteActivity(ctx UnifiedContext, activity interface{}, args ...interface{}) (future Future)
	Now(ctx UnifiedContext) time.Time
	Sleep(ctx UnifiedContext, d time.Duration) (err error)
	GetSignalChannel(ctx UnifiedContext, signalName string) (receiveChannel ReceiveChannel)
	GetContextValue(ctx UnifiedContext, key string) interface{}
}

type ReceiveChannel interface {
	Receive(ctx UnifiedContext, valuePtr interface{}) (more bool)
}

type Future interface {
	Get(ctx UnifiedContext, valuePtr interface{}) error
}
