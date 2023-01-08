package interpreter

import (
	"github.com/indeedeng/iwf/service"
	"time"
)

type timerProcessor struct {
	stateExecutionCurrentTimerInfos map[string][]*service.TimerInfo
	provider                        WorkflowProvider
	logger                          UnifiedLogger
}

func NewTimerProcessor(provider WorkflowProvider) *timerProcessor {
	return &timerProcessor{
		provider:                        provider,
		stateExecutionCurrentTimerInfos: map[string][]*service.TimerInfo{},
	}
}

func (t *timerProcessor) GetCurrentTimerInfos() map[string][]*service.TimerInfo {
	return t.stateExecutionCurrentTimerInfos
}

func (t *timerProcessor) SkipTimer(stateExeId, timerId string, timerIdx int) {
	timer, valid := service.ValidateTimerSkipRequest(t.stateExecutionCurrentTimerInfos, stateExeId, timerId, timerIdx)
	if !valid {
		// since we have checked it before sending signals, this should only happen in some vary rare cases for racing condition
		t.logger.Error("invalid timer skip request received!", stateExeId, timerId, timerIdx)
		return
	}
	timer.Status = service.TimerSkipped
}

// WaitForTimerCompleted waits for timer completed(fired or skipped), return false if the waiting is canceled by cancelWaiting bool pointer
func (t *timerProcessor) WaitForTimerCompleted(ctx UnifiedContext, stateExeId string, timerIdx int, cancelWaiting *bool) bool {
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
