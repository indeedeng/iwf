package temporal

import (
	"errors"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/retry"
	"github.com/indeedeng/iwf/service/interpreter"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type workflowProvider struct {
	threadCount        int
	pendingThreadNames map[string]int
}

func newTemporalWorkflowProvider() interpreter.WorkflowProvider {
	return &workflowProvider{
		pendingThreadNames: map[string]int{},
	}
}

func (w *workflowProvider) GetBackendType() service.BackendType {
	return service.BackendTypeTemporal
}

func (w *workflowProvider) NewApplicationError(errType string, details interface{}) error {
	return temporal.NewApplicationError("", errType, details)
}

func (w *workflowProvider) IsApplicationError(err error) bool {
	var applicationError *temporal.ApplicationError
	return errors.As(err, &applicationError)
}

func (w *workflowProvider) NewInterpreterContinueAsNewError(ctx interpreter.UnifiedContext, input service.InterpreterWorkflowInput) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.NewContinueAsNewError(wfCtx, Interpreter, input)
}

func (w *workflowProvider) UpsertSearchAttributes(ctx interpreter.UnifiedContext, attributes map[string]interface{}) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.UpsertSearchAttributes(wfCtx, attributes)
}

func (w *workflowProvider) UpsertMemo(ctx interpreter.UnifiedContext, rawMemo map[string]iwfidl.EncodedObject) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}

	memo := map[string]interface{}{}
	dataConverter, shouldEncrypt := env.CheckAndGetTemporalMemoEncryptionDataConverter()
	if shouldEncrypt {
		for k, v := range rawMemo {
			pl, err := dataConverter.ToPayload(v)
			if err != nil {
				return err
			}
			memo[k] = pl
		}
	} else {
		for k, v := range rawMemo {
			memo[k] = v
		}
	}

	return workflow.UpsertMemo(wfCtx, memo)
}

func (w *workflowProvider) NewTimer(ctx interpreter.UnifiedContext, d time.Duration) interpreter.Future {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	f := workflow.NewTimer(wfCtx, d)
	return &futureImpl{
		future: f,
	}
}

func (w *workflowProvider) GetWorkflowInfo(ctx interpreter.UnifiedContext) interpreter.WorkflowInfo {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	info := workflow.GetInfo(wfCtx)
	return interpreter.WorkflowInfo{
		WorkflowExecution: interpreter.WorkflowExecution{
			ID:    info.WorkflowExecution.ID,
			RunID: info.WorkflowExecution.RunID,
		},
		WorkflowStartTime:        info.WorkflowStartTime,
		WorkflowExecutionTimeout: info.WorkflowExecutionTimeout,
	}
}

func (w *workflowProvider) SetQueryHandler(ctx interpreter.UnifiedContext, queryType string, handler interface{}) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.SetQueryHandler(wfCtx, queryType, handler)
}

func (w *workflowProvider) SetRpcUpdateHandler(
	ctx interpreter.UnifiedContext, updateType string, validator interpreter.UnifiedRpcValidator, handler interpreter.UnifiedRpcHandler,
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	v2 := func(ctx workflow.Context, input iwfidl.WorkflowRpcRequest) error {
		ctx2 := interpreter.NewUnifiedContext(ctx)
		return validator(ctx2, input)
	}
	h2 := func(ctx workflow.Context, input iwfidl.WorkflowRpcRequest) (*interpreter.HandlerOutput, error) {
		ctx2 := interpreter.NewUnifiedContext(ctx)
		return handler(ctx2, input)
	}
	return workflow.SetUpdateHandlerWithOptions(
		wfCtx,
		updateType,
		h2,
		workflow.UpdateHandlerOptions{Validator: v2},
	)
}

func (w *workflowProvider) ExtendContextWithValue(parent interpreter.UnifiedContext, key string, val interface{}) interpreter.UnifiedContext {
	wfCtx, ok := parent.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return interpreter.NewUnifiedContext(workflow.WithValue(wfCtx, key, val))
}

func (w *workflowProvider) GoNamed(ctx interpreter.UnifiedContext, name string, f func(ctx interpreter.UnifiedContext)) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	f2 := func(ctx workflow.Context) {
		ctx2 := interpreter.NewUnifiedContext(ctx)
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

func (w *workflowProvider) GetPendingThreadNames() map[string]int {
	return w.pendingThreadNames
}

func (w *workflowProvider) GetThreadCount() int {
	return w.threadCount
}

func (w *workflowProvider) Await(ctx interpreter.UnifiedContext, condition func() bool) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Await(wfCtx, condition)
}

func (w *workflowProvider) WithActivityOptions(ctx interpreter.UnifiedContext, options interpreter.ActivityOptions) interpreter.UnifiedContext {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}

	// in Temporal, scheduled to close timeout is the timeout include all retries
	scheduledToCloseTimeout := time.Duration(0)
	if options.RetryPolicy.GetMaximumAttemptsDurationSeconds() > 0 {
		scheduledToCloseTimeout = time.Second * time.Duration(options.RetryPolicy.GetMaximumAttemptsDurationSeconds())
	}

	wfCtx2 := workflow.WithActivityOptions(wfCtx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: scheduledToCloseTimeout,
		StartToCloseTimeout:    options.StartToCloseTimeout,
		RetryPolicy:            retry.ConvertTemporalActivityRetryPolicy(options.RetryPolicy),
	})
	return interpreter.NewUnifiedContext(wfCtx2)
}

type futureImpl struct {
	future workflow.Future
}

func (t *futureImpl) IsReady() bool {
	return t.future.IsReady()
}

func (t *futureImpl) Get(ctx interpreter.UnifiedContext, valuePtr interface{}) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}

	return t.future.Get(wfCtx, valuePtr)
}

func (w *workflowProvider) ExecuteActivity(ctx interpreter.UnifiedContext, activity interface{}, args ...interface{}) (future interpreter.Future) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	f := workflow.ExecuteActivity(wfCtx, activity, args...)
	return &futureImpl{
		future: f,
	}
}

func (w *workflowProvider) Now(ctx interpreter.UnifiedContext) time.Time {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Now(wfCtx)
}

func (w *workflowProvider) Sleep(ctx interpreter.UnifiedContext, d time.Duration) (err error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Sleep(wfCtx, d)
}

func (w *workflowProvider) IsReplaying(ctx interpreter.UnifiedContext) bool {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.IsReplaying(wfCtx)
}

func (w *workflowProvider) GetVersion(ctx interpreter.UnifiedContext, changeID string, minSupported, maxSupported int) int {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}

	version := workflow.GetVersion(wfCtx, changeID, workflow.Version(minSupported), workflow.Version(maxSupported))
	return int(version)
}

type temporalReceiveChannel struct {
	channel workflow.ReceiveChannel
}

func (t *temporalReceiveChannel) ReceiveAsync(valuePtr interface{}) (ok bool) {
	return t.channel.ReceiveAsync(valuePtr)
}

func (t *temporalReceiveChannel) ReceiveBlocking(ctx interpreter.UnifiedContext, valuePtr interface{}) (ok bool) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}

	return t.channel.Receive(wfCtx, valuePtr)
}

func (w *workflowProvider) GetSignalChannel(ctx interpreter.UnifiedContext, signalName string) interpreter.ReceiveChannel {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	wfChan := workflow.GetSignalChannel(wfCtx, signalName)
	return &temporalReceiveChannel{
		channel: wfChan,
	}
}

func (w *workflowProvider) GetContextValue(ctx interpreter.UnifiedContext, key string) interface{} {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return wfCtx.Value(key)
}

func (w *workflowProvider) GetLogger(ctx interpreter.UnifiedContext) interpreter.UnifiedLogger {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.GetLogger(wfCtx)
}

func (w *workflowProvider) GetUnhandledSignalNames(ctx interpreter.UnifiedContext) []string {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.GetUnhandledSignalNames(wfCtx)
}
