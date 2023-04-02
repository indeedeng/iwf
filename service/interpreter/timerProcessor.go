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
		logger:                          provider.GetLogger(ctx),
	}

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

// WaitForTimerFiredOrSkipped waits for timer completed(fired or skipped),
// return true when the timer is fired or skipped
// return false if the waitingCommands is canceled by cancelWaiting bool pointer(when the trigger type is completed, or continueAsNew)
func (t *TimerProcessor) WaitForTimerFiredOrSkipped(ctx UnifiedContext, stateExeId string, timerIdx int, cancelWaiting *bool) bool {
	timer := t.stateExecutionCurrentTimerInfos[stateExeId][timerIdx]
	if timer.Status == service.TimerFired || timer.Status == service.TimerSkipped {
		return true
	}
	now := t.provider.Now(ctx).Unix()
	fireAt := timer.FiringUnixTimestampSeconds
	duration := time.Duration(fireAt-now) * time.Second
	future := t.provider.NewTimer(ctx, duration)
	_ = t.provider.Await(ctx, func() bool {
		return future.IsReady() || timer.Status == service.TimerSkipped || *cancelWaiting
	})
	if future.IsReady() || timer.Status == service.TimerSkipped {
		return true
	}
	// otherwise *cancelWaiting should return false to indicate that this timer isn't completed(fired or skipped)
	return false
}

// RemovePendingTimersOfState is for when a state is completed, remove all its pending timers
func (t *TimerProcessor) RemovePendingTimersOfState(stateExeId string) {
	delete(t.stateExecutionCurrentTimerInfos, stateExeId)
}

func (t *TimerProcessor) AddTimers(stateExeId string, commands []iwfidl.TimerCommand, completedTimerCmds map[int]bool) {
	timers := make([]*service.TimerInfo, len(commands))
	for idx, cmd := range commands {
		var timer service.TimerInfo
		if completedTimerCmds[idx] {
			timer = service.TimerInfo{
				CommandId:                  cmd.CommandId,
				FiringUnixTimestampSeconds: cmd.GetFiringUnixTimestampSeconds(),
				Status:                     service.TimerFired, // TODO differentiate skipped and fire
			}
		} else {
			timer = service.TimerInfo{
				CommandId:                  cmd.CommandId,
				FiringUnixTimestampSeconds: cmd.GetFiringUnixTimestampSeconds(),
				Status:                     service.TimerPending,
			}
		}

		timers[idx] = &timer
	}
	t.stateExecutionCurrentTimerInfos[stateExeId] = timers
}
