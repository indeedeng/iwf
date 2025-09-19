package temporal

import (
	"errors"
	"time"

	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/retry"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type workflowProvider struct {
	threadCount        int
	pendingThreadNames map[string]int
}

func newTemporalWorkflowProvider() interfaces.WorkflowProvider {
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

func (w *workflowProvider) NewInterpreterContinueAsNewError(
	ctx interfaces.UnifiedContext, input service.InterpreterWorkflowInput,
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.NewContinueAsNewError(wfCtx, Interpreter, input)
}

func (w *workflowProvider) UpsertSearchAttributes(
	ctx interfaces.UnifiedContext, attributes map[string]interface{},
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.UpsertSearchAttributes(wfCtx, attributes)
}

func (w *workflowProvider) UpsertMemo(ctx interfaces.UnifiedContext, rawMemo map[string]iwfidl.EncodedObject) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}

	memo := map[string]interface{}{}
	dataConverter, shouldEncrypt := env.CheckAndGetTemporalMemoEncryptionDataConverter()
	if shouldEncrypt {
		for _, key := range workflow.DeterministicKeys(rawMemo) {
			pl, err := dataConverter.ToPayload(rawMemo[key])
			if err != nil {
				return err
			}
			memo[key] = pl
		}
	} else {
		for _, key := range workflow.DeterministicKeys(rawMemo) {
			memo[key] = rawMemo[key]
		}
	}

	return workflow.UpsertMemo(wfCtx, memo)
}

func (w *workflowProvider) NewTimer(ctx interfaces.UnifiedContext, d time.Duration) interfaces.Future {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	f := workflow.NewTimer(wfCtx, d)
	return &futureImpl{
		future: f,
	}
}

func (w *workflowProvider) GetWorkflowInfo(ctx interfaces.UnifiedContext) interfaces.WorkflowInfo {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	info := workflow.GetInfo(wfCtx)
	return interfaces.WorkflowInfo{
		WorkflowExecution: interfaces.WorkflowExecution{
			ID:    info.WorkflowExecution.ID,
			RunID: info.WorkflowExecution.RunID,
		},
		WorkflowStartTime:        info.WorkflowStartTime,
		WorkflowExecutionTimeout: info.WorkflowExecutionTimeout,
		FirstRunID:               info.FirstRunID,
		CurrentRunID:             info.WorkflowExecution.RunID,
	}
}

func (w *workflowProvider) GetSearchAttributes(
	ctx interfaces.UnifiedContext, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType,
) (map[string]iwfidl.SearchAttribute, error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	sas := workflow.GetInfo(wfCtx).SearchAttributes

	return mapper.MapTemporalToIwfSearchAttributes(sas, requestedSearchAttributes)
}

func (w *workflowProvider) SetQueryHandler(
	ctx interfaces.UnifiedContext, queryType string, handler interface{},
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.SetQueryHandler(wfCtx, queryType, handler)
}

func (w *workflowProvider) SetRpcUpdateHandler(
	ctx interfaces.UnifiedContext, updateType string, validator interfaces.UnifiedRpcValidator,
	handler interfaces.UnifiedRpcHandler,
) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	v2 := func(ctx workflow.Context, input iwfidl.WorkflowRpcRequest) error {
		ctx2 := interfaces.NewUnifiedContext(ctx)
		return validator(ctx2, input)
	}
	h2 := func(ctx workflow.Context, input iwfidl.WorkflowRpcRequest) (*interfaces.HandlerOutput, error) {
		ctx2 := interfaces.NewUnifiedContext(ctx)
		return handler(ctx2, input)
	}
	return workflow.SetUpdateHandlerWithOptions(
		wfCtx,
		updateType,
		h2,
		workflow.UpdateHandlerOptions{Validator: v2},
	)
}

func (w *workflowProvider) ExtendContextWithValue(
	parent interfaces.UnifiedContext, key string, val interface{},
) interfaces.UnifiedContext {
	wfCtx, ok := parent.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return interfaces.NewUnifiedContext(workflow.WithValue(wfCtx, key, val))
}

func (w *workflowProvider) GoNamed(
	ctx interfaces.UnifiedContext, name string, f func(ctx interfaces.UnifiedContext),
) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
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

func (w *workflowProvider) GetPendingThreadNames() map[string]int {
	return w.pendingThreadNames
}

func (w *workflowProvider) GetThreadCount() int {
	return w.threadCount
}

func (w *workflowProvider) Await(ctx interfaces.UnifiedContext, condition func() bool) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Await(wfCtx, condition)
}

func (w *workflowProvider) WithActivityOptions(
	ctx interfaces.UnifiedContext, options interfaces.ActivityOptions,
) interfaces.UnifiedContext {
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
		HeartbeatTimeout:       options.HeartbeatTimeout,
	})

	// support local activity optimization
	wfCtx3 := workflow.WithLocalActivityOptions(wfCtx2, workflow.LocalActivityOptions{
		// set the LA timeout to 7s to make sure the workflow will not need a heartbeat
		ScheduleToCloseTimeout: time.Second * 7,
		RetryPolicy:            retry.ConvertTemporalActivityRetryPolicy(options.RetryPolicy),
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
		panic("cannot convert to temporal workflow context")
	}

	return t.future.Get(wfCtx, valuePtr)
}

func (w *workflowProvider) ExecuteActivity(
	valuePtr interface{}, optimizeByLocalActivity bool,
	ctx interfaces.UnifiedContext, activity interface{}, args ...interface{},
) (err error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
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
		panic("cannot convert to temporal workflow context")
	}

	f := workflow.ExecuteLocalActivity(wfCtx, activity, args...)
	return f.Get(wfCtx, valuePtr)
}

func (w *workflowProvider) Now(ctx interfaces.UnifiedContext) time.Time {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Now(wfCtx)
}

func (w *workflowProvider) Sleep(ctx interfaces.UnifiedContext, d time.Duration) (err error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Sleep(wfCtx, d)
}

func (w *workflowProvider) IsReplaying(ctx interfaces.UnifiedContext) bool {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.IsReplaying(wfCtx)
}

func (w *workflowProvider) GetVersion(
	ctx interfaces.UnifiedContext, changeID string, minSupported, maxSupported int,
) int {
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

func (t *temporalReceiveChannel) ReceiveBlocking(ctx interfaces.UnifiedContext, valuePtr interface{}) (ok bool) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}

	return t.channel.Receive(wfCtx, valuePtr)
}

func (w *workflowProvider) GetSignalChannel(
	ctx interfaces.UnifiedContext, signalName string,
) interfaces.ReceiveChannel {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	wfChan := workflow.GetSignalChannel(wfCtx, signalName)
	return &temporalReceiveChannel{
		channel: wfChan,
	}
}

func (w *workflowProvider) GetContextValue(ctx interfaces.UnifiedContext, key string) interface{} {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return wfCtx.Value(key)
}

func (w *workflowProvider) GetLogger(ctx interfaces.UnifiedContext) interfaces.UnifiedLogger {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.GetLogger(wfCtx)
}

func (w *workflowProvider) GetUnhandledSignalNames(ctx interfaces.UnifiedContext) []string {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.GetUnhandledSignalNames(wfCtx)
}
