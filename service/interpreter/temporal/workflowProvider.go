package temporal

import (
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type workflowProvider struct{}

var defaultWorkflowProvider interpreter.WorkflowProvider = &workflowProvider{}

func (w *workflowProvider) GetBackendType() service.BackendType {
	return service.BackendTypeTemporal
}

func (w *workflowProvider) NewApplicationError(message, errType string, details ...interface{}) error {
	return temporal.NewApplicationError(message, errType, details...)
}

func (w *workflowProvider) UpsertSearchAttributes(ctx interpreter.UnifiedContext, attributes map[string]interface{}) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.UpsertSearchAttributes(wfCtx, attributes)
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
		WorkflowStartTime: info.WorkflowStartTime,
	}
}

func (w *workflowProvider) SetQueryHandler(ctx interpreter.UnifiedContext, queryType string, handler interface{}) error {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.SetQueryHandler(wfCtx, queryType, handler)
}

func (w *workflowProvider) ExtendContextWithValue(parent interpreter.UnifiedContext, key string, val interface{}) interpreter.UnifiedContext {
	wfCtx, ok := parent.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return interpreter.NewUnifiedContext(workflow.WithValue(wfCtx, key, val))
}

func (w workflowProvider) GoNamed(ctx interpreter.UnifiedContext, name string, f func(ctx interpreter.UnifiedContext)) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	f2 := func(ctx workflow.Context) {
		ctx2 := interpreter.NewUnifiedContext(ctx)
		f(ctx2)
	}
	workflow.GoNamed(wfCtx, name, f2)
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
	wfCtx2 := workflow.WithActivityOptions(wfCtx, workflow.ActivityOptions{
		StartToCloseTimeout: options.StartToCloseTimeout,
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

func (t *temporalReceiveChannel) Receive(ctx interpreter.UnifiedContext, valuePtr interface{}) (more bool) {
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
