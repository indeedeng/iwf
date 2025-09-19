package cadence

import (
	"fmt"
	"time"

	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/retry"
	"go.uber.org/cadence"
	"go.uber.org/cadence/workflow"
)

type workflowProvider struct {
	threadCount        int
	pendingThreadNames map[string]int
}

func newCadenceWorkflowProvider() interfaces.WorkflowProvider {
	return &workflowProvider{
		pendingThreadNames: map[string]int{},
	}
}

func (w *workflowProvider) GetBackendType() service.BackendType {
	return service.BackendTypeCadence
}

func (w *workflowProvider) NewApplicationError(errType string, details interface{}) error {
	return cadence.NewCustomError(errType, details)
}

func (w *workflowProvider) IsApplicationError(err error) bool {
	_, ok := err.(*cadence.CustomError)
	return ok
}

func (w *workflowProvider) NewInterpreterContinueAsNewError(
	ctx interfaces.UnifiedContext, input service.InterpreterWorkflowInput,
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.NewContinueAsNewError(wfCtx, Interpreter, input)
}

func (w *workflowProvider) UpsertSearchAttributes(
	ctx interfaces.UnifiedContext, attributes map[string]interface{},
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.UpsertSearchAttributes(wfCtx, attributes)
}

func (w *workflowProvider) UpsertMemo(ctx interfaces.UnifiedContext, memo map[string]iwfidl.EncodedObject) error {
	return fmt.Errorf("upsert memo is not supported in Cadence")
}

func (w *workflowProvider) NewTimer(ctx interfaces.UnifiedContext, d time.Duration) interfaces.Future {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	f := workflow.NewTimer(wfCtx, d)
	return &futureImpl{
		future: f,
	}
}

func (w *workflowProvider) GetWorkflowInfo(ctx interfaces.UnifiedContext) interfaces.WorkflowInfo {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	info := workflow.GetInfo(wfCtx)
	return interfaces.WorkflowInfo{
		WorkflowExecution: interfaces.WorkflowExecution{
			ID:    info.WorkflowExecution.ID,
			RunID: info.WorkflowExecution.RunID,
		},
		WorkflowStartTime:        time.UnixMilli(0), // TODO need support from Cadence client: https://github.com/uber-go/cadence-client/issues/1204
		WorkflowExecutionTimeout: time.Duration(info.ExecutionStartToCloseTimeoutSeconds) * time.Second,
		FirstRunID:               info.WorkflowExecution.RunID, // Cadence does not provide FirstRunID TODO https://github.com/uber-go/cadence-client/issues/1371 use firstRunID when available
		CurrentRunID:             info.WorkflowExecution.RunID,
	}
}

func (w *workflowProvider) GetSearchAttributes(
	ctx interfaces.UnifiedContext, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType,
) (map[string]iwfidl.SearchAttribute, error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	sas := workflow.GetInfo(wfCtx).SearchAttributes

	return mapper.MapCadenceToIwfSearchAttributes(sas, requestedSearchAttributes)
}

func (w *workflowProvider) SetQueryHandler(
	ctx interfaces.UnifiedContext, queryType string, handler interface{},
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.SetQueryHandler(wfCtx, queryType, handler)
}

func (w *workflowProvider) SetRpcUpdateHandler(
	ctx interfaces.UnifiedContext, updateType string, validator interfaces.UnifiedRpcValidator,
	handler interfaces.UnifiedRpcHandler,
) error {
	// NOTE: this feature is not available in Cadence
	return nil
}

func (w *workflowProvider) ExtendContextWithValue(
	parent interfaces.UnifiedContext, key string, val interface{},
) interfaces.UnifiedContext {
	wfCtx, ok := parent.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return interfaces.NewUnifiedContext(workflow.WithValue(wfCtx, key, val))
}

func (w *workflowProvider) GoNamed(
	ctx interfaces.UnifiedContext, name string, f func(ctx interfaces.UnifiedContext),
) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	f2 := func(ctx workflow.Context) {
		ctx2 := interfaces.NewUnifiedContext(ctx)
		w.pendingThreadNames[name]++
		w.threadCount++
		f(ctx2)
		w.pendingThreadNames[name]--
		if w.pendingThreadNames[name] == 0 {
			delete(w.pendingThreadNames, name)
		}
		w.threadCount--
	}
	workflow.GoNamed(wfCtx, name, f2)
}

func (w *workflowProvider) WaitForThreadByName(ctx interfaces.UnifiedContext, name string) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Await(wfCtx, func() bool {
		_, ok := w.pendingThreadNames[name]
		// when the threadName doesn't exist in the map anymore
		// it means the thread has completed
		return !ok
	})
}

func (w *workflowProvider) GetPendingThreadNames() map[string]int {
	return w.pendingThreadNames
}

func (w *workflowProvider) GetThreadCount() int {
	return w.threadCount
}

func (w *workflowProvider) Await(ctx interfaces.UnifiedContext, condition func() bool) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.Await(wfCtx, condition)
}

func (w *workflowProvider) WithActivityOptions(
	ctx interfaces.UnifiedContext, options interfaces.ActivityOptions,
) interfaces.UnifiedContext {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}

	unlimited := time.Hour * 24 * 365
	startToCloseTimeout := options.StartToCloseTimeout
	if startToCloseTimeout == 0 {
		// unlimited to match Temporal for default
		startToCloseTimeout = unlimited
	}

	wfCtx2 := workflow.WithActivityOptions(wfCtx, workflow.ActivityOptions{
		StartToCloseTimeout:    startToCloseTimeout,
		ScheduleToStartTimeout: time.Second * 10,
		HeartbeatTimeout:       options.HeartbeatTimeout,
		RetryPolicy:            retry.ConvertCadenceActivityRetryPolicy(options.RetryPolicy),
	})

	// support local activity optimization
	wfCtx3 := workflow.WithLocalActivityOptions(wfCtx2, workflow.LocalActivityOptions{
		// set the LA timeout to 7s to make sure the workflow will not need a heartbeat
		ScheduleToCloseTimeout: time.Second * 7,
		RetryPolicy:            retry.ConvertCadenceActivityRetryPolicy(options.RetryPolicy),
	})
	return interfaces.NewUnifiedContext(wfCtx3)
}

type futureImpl struct {
	future workflow.Future
}

func (t *futureImpl) IsReady() bool {
	return t.future.IsReady()
}

func (t *futureImpl) Get(ctx interfaces.UnifiedContext, valuePtr interface{}) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}

	return t.future.Get(wfCtx, valuePtr)
}

func (w *workflowProvider) ExecuteActivity(
	valuePtr interface{}, optimizeByLocalActivity bool,
	ctx interfaces.UnifiedContext, activity interface{}, args ...interface{},
) (err error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	if optimizeByLocalActivity {
		f := workflow.ExecuteLocalActivity(wfCtx, activity, args...)
		err = f.Get(wfCtx, valuePtr)
		if err != nil {
			f = workflow.ExecuteActivity(wfCtx, activity, args...)
			return f.Get(wfCtx, valuePtr)
		}
		return err
	}

	f := workflow.ExecuteActivity(wfCtx, activity, args...)
	return f.Get(wfCtx, valuePtr)
}

func (w *workflowProvider) ExecuteLocalActivity(
	valuePtr interface{}, ctx interfaces.UnifiedContext, activity interface{}, args ...interface{},
) (err error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}

	f := workflow.ExecuteLocalActivity(wfCtx, activity, args...)
	return f.Get(wfCtx, valuePtr)
}

func (w *workflowProvider) Now(ctx interfaces.UnifiedContext) time.Time {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.Now(wfCtx)
}

func (w *workflowProvider) IsReplaying(ctx interfaces.UnifiedContext) bool {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.IsReplaying(wfCtx)
}

func (w *workflowProvider) Sleep(ctx interfaces.UnifiedContext, d time.Duration) (err error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.Sleep(wfCtx, d)
}

func (w *workflowProvider) GetVersion(
	ctx interfaces.UnifiedContext, changeID string, minSupported, maxSupported int,
) int {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}

	version := workflow.GetVersion(wfCtx, changeID, workflow.Version(minSupported), workflow.Version(maxSupported))
	return int(version)
}

type cadenceReceiveChannel struct {
	channel workflow.Channel
}

func (t *cadenceReceiveChannel) ReceiveAsync(valuePtr interface{}) (ok bool) {
	return t.channel.ReceiveAsync(valuePtr)
}

func (t *cadenceReceiveChannel) ReceiveBlocking(ctx interfaces.UnifiedContext, valuePtr interface{}) (ok bool) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}

	return t.channel.Receive(wfCtx, valuePtr)
}

func (w *workflowProvider) GetSignalChannel(
	ctx interfaces.UnifiedContext, signalName string,
) interfaces.ReceiveChannel {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	wfChan := workflow.GetSignalChannel(wfCtx, signalName)
	return &cadenceReceiveChannel{
		channel: wfChan,
	}
}

func (w *workflowProvider) GetContextValue(ctx interfaces.UnifiedContext, key string) interface{} {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return wfCtx.Value(key)
}

func (w *workflowProvider) GetLogger(ctx interfaces.UnifiedContext) interfaces.UnifiedLogger {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}

	zLogger := workflow.GetLogger(wfCtx)
	return &loggerImpl{
		zlogger: zLogger,
	}
}

func (w *workflowProvider) GetUnhandledSignalNames(ctx interfaces.UnifiedContext) []string {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to cadence workflow context")
	}
	return workflow.GetUnhandledSignalNames(wfCtx)
}
