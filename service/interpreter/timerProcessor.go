package interpreter

import (
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type TimerProcessor struct {
	stateExecutionCurrentTimerInfos map[string][]*service.TimerInfo
	provider                        WorkflowProvider
	logger                          UnifiedLogger
}

func NewTimerProcessor(ctx UnifiedContext, provider WorkflowProvider) *TimerProcessor {
	tp := &TimerProcessor{
		provider:                        provider,
		stateExecutionCurrentTimerInfos: map[string][]*service.TimerInfo{},
	}
	provider.GoNamed(ctx, "skip-timer-signal-handler", func(ctx UnifiedContext) {
		for {
			ch := provider.GetSignalChannel(ctx, service.SkipTimerSignalChannelName)
			val := service.SkipTimerSignalRequest{}

			err := provider.Await(ctx, func() bool {
				return ch.ReceiveAsync(&val)
			})
			if err != nil {
				// break the loop to prevent goroutine leakage
				break
			}

			tp.SkipTimer(val.StateExecutionId, val.CommandId, val.CommandIndex)
		}
	})
	err := provider.SetQueryHandler(ctx, service.GetCurrentTimerInfosQueryType, func() (service.GetCurrentTimerInfosQueryResponse, error) {
		return service.GetCurrentTimerInfosQueryResponse{
			StateExecutionCurrentTimerInfos: tp.stateExecutionCurrentTimerInfos,
		}, nil
	})
	if err != nil {
		panic("cannot set query handler")
	}
	return tp
}

func (t *TimerProcessor) GetCurrentTimerInfos() map[string][]*service.TimerInfo {
	return t.stateExecutionCurrentTimerInfos
}

func (t *TimerProcessor) SkipTimer(stateExeId, timerId string, timerIdx int) {
	timer, valid := service.ValidateTimerSkipRequest(t.stateExecutionCurrentTimerInfos, stateExeId, timerId, timerIdx)
	if !valid {
		// since we have checked it before sending signals, this should only happen in some vary rare cases for racing condition
		t.logger.Error("invalid timer skip request received!", stateExeId, timerId, timerIdx)
		return
	}
	timer.Status = service.TimerSkipped
}

// WaitForTimerCompleted waits for timer completed(fired or skipped), return false if the waiting is canceled by cancelWaiting bool pointer
func (t *TimerProcessor) WaitForTimerCompleted(ctx UnifiedContext, stateExeId string, timerIdx int, cancelWaiting *bool) bool {
	timer := t.stateExecutionCurrentTimerInfos[stateExeId][timerIdx]
	now := t.provider.Now(ctx).Unix()
	fireAt := timer.FiringUnixTimestampSeconds
	duration := time.Duration(fireAt-now) * time.Second
	future := t.provider.NewTimer(ctx, duration)
	_ = t.provider.Await(ctx, func() bool {
		return future.IsReady() || timer.Status == service.TimerSkipped || *cancelWaiting
	})
	if *cancelWaiting {
		return false
	}
	return true
}

func (t *TimerProcessor) FinishProcessing(stateExeId string) {
	delete(t.stateExecutionCurrentTimerInfos, stateExeId)
}

func (t *TimerProcessor) StartProcessing(stateExeId string, commands []iwfidl.TimerCommand) {
	timers := make([]*service.TimerInfo, len(commands))
	for idx, cmd := range commands {
		timer := service.TimerInfo{
			CommandId:                  cmd.CommandId,
			FiringUnixTimestampSeconds: cmd.GetFiringUnixTimestampSeconds(),
			Status:                     service.TimerPending,
		}
		timers[idx] = &timer
	}
	t.stateExecutionCurrentTimerInfos[stateExeId] = timers
}
