package interfaces

import (
	"context"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/errors"
)

type ActivityProvider interface {
	GetLogger(ctx context.Context) UnifiedLogger
	NewApplicationError(errType string, details interface{}) error
	GetActivityInfo(ctx context.Context) ActivityInfo
	RecordHeartbeat(ctx context.Context, details ...interface{})
}

type ActivityInfo struct {
	ScheduledTime     time.Time // Time of activity scheduled by a workflow
	Attempt           int32     // Attempt starts from 1, and increased by 1 for every retry if retry policy is specified.
	IsLocalActivity   bool      // Whether the activity is at local activity
	WorkflowExecution WorkflowExecution
}

var activityProviderRegistry = make(map[service.BackendType]ActivityProvider)

func RegisterActivityProvider(backendType service.BackendType, provider ActivityProvider) {
	if _, ok := activityProviderRegistry[backendType]; ok {
		panic("backend type " + backendType + " has been registered")
	}
	activityProviderRegistry[backendType] = provider
}

func GetActivityProviderByType(backendType service.BackendType) ActivityProvider {
	provider := activityProviderRegistry[backendType]
	if provider == nil {
		panic("not supported yet: " + backendType)
	}
	return provider
}

type UnifiedLogger interface {
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
	WorkflowExecution        WorkflowExecution
	WorkflowStartTime        time.Time
	WorkflowExecutionTimeout time.Duration
	FirstRunID               string
	CurrentRunID             string
}

type ActivityOptions struct {
	StartToCloseTimeout time.Duration
	HeartbeatTimeout    time.Duration
	RetryPolicy         *iwfidl.RetryPolicy
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

type TimerProcessor interface {
	Dump() []service.StaleSkipTimerSignal
	SkipTimer(stateExeId string, timerId string, timerIdx int) bool
	RetryStaleSkipTimer() bool
	WaitForTimerFiredOrSkipped(ctx UnifiedContext, stateExeId string, timerIdx int, cancelWaiting *bool) service.InternalTimerStatus
	RemovePendingTimersOfState(stateExeId string)
	AddTimers(stateExeId string, commands []iwfidl.TimerCommand, completedTimerCmds map[int]service.InternalTimerStatus)
	GetTimerInfos() map[string][]*service.TimerInfo
	GetTimerStartedUnixTimestamps() []int64
}

type WorkflowProvider interface {
	NewApplicationError(errType string, details interface{}) error
	IsApplicationError(err error) bool
	GetWorkflowInfo(ctx UnifiedContext) WorkflowInfo
	GetSearchAttributes(
		ctx UnifiedContext, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType,
	) (map[string]iwfidl.SearchAttribute, error)
	UpsertSearchAttributes(ctx UnifiedContext, attributes map[string]interface{}) error
	UpsertMemo(ctx UnifiedContext, memo map[string]iwfidl.EncodedObject) error
	SetQueryHandler(ctx UnifiedContext, queryType string, handler interface{}) error
	SetRpcUpdateHandler(
		ctx UnifiedContext, updateType string, validator UnifiedRpcValidator, handler UnifiedRpcHandler,
	) error
	ExtendContextWithValue(parent UnifiedContext, key string, val interface{}) UnifiedContext
	GoNamed(ctx UnifiedContext, name string, f func(ctx UnifiedContext))
	GetThreadCount() int
	GetPendingThreadNames() map[string]int
	Await(ctx UnifiedContext, condition func() bool) error
	WithActivityOptions(ctx UnifiedContext, options ActivityOptions) UnifiedContext
	ExecuteActivity(
		valuePtr interface{}, optimizeByLocalActivity bool, ctx UnifiedContext, activity interface{},
		args ...interface{},
	) (err error)
	ExecuteLocalActivity(
		valuePtr interface{}, ctx UnifiedContext, activity interface{}, args ...interface{},
	) (err error)
	Now(ctx UnifiedContext) time.Time
	IsReplaying(ctx UnifiedContext) bool
	Sleep(ctx UnifiedContext, d time.Duration) (err error)
	NewTimer(ctx UnifiedContext, d time.Duration) Future
	GetSignalChannel(ctx UnifiedContext, signalName string) (receiveChannel ReceiveChannel)
	GetContextValue(ctx UnifiedContext, key string) interface{}
	GetVersion(ctx UnifiedContext, changeID string, minSupported, maxSupported int) int
	GetUnhandledSignalNames(ctx UnifiedContext) []string
	GetBackendType() service.BackendType
	GetLogger(ctx UnifiedContext) UnifiedLogger
	NewInterpreterContinueAsNewError(ctx UnifiedContext, input service.InterpreterWorkflowInput) error
}

type ReceiveChannel interface {
	ReceiveAsync(valuePtr interface{}) (ok bool)
	ReceiveBlocking(ctx UnifiedContext, valuePtr interface{}) (ok bool)
}

type Future interface {
	Get(ctx UnifiedContext, valuePtr interface{}) error
	IsReady() bool
}

type HandlerOutput struct {
	RpcOutput   *iwfidl.WorkflowRpcResponse
	StatusError *errors.ErrorAndStatus
}
type InvokeRpcActivityOutput struct {
	RpcOutput   *iwfidl.WorkflowWorkerRpcResponse
	StatusError *errors.ErrorAndStatus
}
type UnifiedRpcHandler func(ctx UnifiedContext, input iwfidl.WorkflowRpcRequest) (*HandlerOutput, error)
type UnifiedRpcValidator func(ctx UnifiedContext, input iwfidl.WorkflowRpcRequest) error
