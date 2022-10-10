package temporal

import (
	"github.com/cadence-oss/iwf-server/service/interpreter"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type workflowProvider struct{}

var DefaultWorkflowProvider = getWorkflowProvider()

func getWorkflowProvider() interpreter.WorkflowProvider {
	return &workflowProvider{}
}

func (w *workflowProvider) NewApplicationError(message, errType string, details ...interface{}) error {
	return temporal.NewApplicationError(message, errType, details...)
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

func (w *workflowProvider) ExtendContextWithValue(parent interface{}, key interface{}, val interface{}) interface{} {
	wfCtx, ok := parent.(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.WithValue(wfCtx, key, val)
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

func (w *workflowProvider) WithActivityOptions(ctx interpreter.UnifiedContext, options interpreter.ActivityOptions) interface{} {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.WithActivityOptions(wfCtx, workflow.ActivityOptions{
		StartToCloseTimeout: options.StartToCloseTimeout,
	})
}

func (w *workflowProvider) ExecuteActivity(ctx interpreter.UnifiedContext, activity interface{}, args ...interface{}) (future interface{}) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.ExecuteActivity(wfCtx, activity, args...)
}

func (w workflowProvider) Now(ctx interpreter.UnifiedContext) time.Time {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Now(wfCtx)
}

func (w workflowProvider) Sleep(ctx interpreter.UnifiedContext, d time.Duration) (err error) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return workflow.Sleep(wfCtx, d)
}

type temporalReceiveChannel struct {
	channel workflow.ReceiveChannel
}

func (t *temporalReceiveChannel) Receive(ctx interpreter.UnifiedContext, valuePtr interface{}) (more bool) {
	wfCtx, ok := ctx.GetContext().(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	return t.channel.Receive(wfCtx, valuePtr)
}

func (w workflowProvider) GetSignalChannel(ctx interpreter.UnifiedContext, signalName string) interpreter.ReceiveChannel {
	wfCtx, ok := ctx.(workflow.Context)
	if !ok {
		panic("cannot convert to temporal workflow context")
	}
	wfChan := workflow.GetSignalChannel(wfCtx, signalName)
	return &temporalReceiveChannel{
		channel: wfChan,
	}
}
