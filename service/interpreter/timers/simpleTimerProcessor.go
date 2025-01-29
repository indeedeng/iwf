package timers

import (
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type SimpleTimerProcessor struct {
	stateExecutionCurrentTimerInfos map[string][]*service.TimerInfo
	staleSkipTimerSignals           []service.StaleSkipTimerSignal
	provider                        interfaces.WorkflowProvider
	logger                          interfaces.UnifiedLogger
}

func NewSimpleTimerProcessor(
	ctx interfaces.UnifiedContext, provider interfaces.WorkflowProvider, staleSkipTimerSignals []service.StaleSkipTimerSignal,
) *SimpleTimerProcessor {
	tp := &SimpleTimerProcessor{
		provider:                        provider,
		stateExecutionCurrentTimerInfos: map[string][]*service.TimerInfo{},
		logger:                          provider.GetLogger(ctx),
		staleSkipTimerSignals:           staleSkipTimerSignals,
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

func (t *SimpleTimerProcessor) Dump() []service.StaleSkipTimerSignal {
	return t.staleSkipTimerSignals
}

// SkipTimer will attempt to skip a timer, return false if no valid timer found
func (t *SimpleTimerProcessor) SkipTimer(stateExeId, timerId string, timerIdx int) bool {
	timer, valid := service.ValidateTimerSkipRequest(t.stateExecutionCurrentTimerInfos, stateExeId, timerId, timerIdx)
	if !valid {
		// since we have checked it before sending signals, this should only happen in some vary rare cases for racing condition
		t.logger.Warn("cannot process timer skip request, maybe state is already closed...putting into a stale skip timer queue", stateExeId, timerId, timerIdx)

		t.staleSkipTimerSignals = append(t.staleSkipTimerSignals, service.StaleSkipTimerSignal{
			StateExecutionId:  stateExeId,
			TimerCommandId:    timerId,
			TimerCommandIndex: timerIdx,
		})
		return false
	}
	timer.Status = service.TimerSkipped
	return true
}

func (t *SimpleTimerProcessor) RetryStaleSkipTimer() bool {
	for i, staleSkip := range t.staleSkipTimerSignals {
		found := t.SkipTimer(staleSkip.StateExecutionId, staleSkip.TimerCommandId, staleSkip.TimerCommandIndex)
		if found {
			newList := removeElement(t.staleSkipTimerSignals, i)
			t.staleSkipTimerSignals = newList
			return true
		}
	}
	return false
}

// WaitForTimerFiredOrSkipped waits for timer completed(fired or skipped),
// return true when the timer is fired or skipped
// return false if the waitingCommands is canceled by cancelWaiting bool pointer(when the trigger type is completed, or continueAsNew)
func (t *SimpleTimerProcessor) WaitForTimerFiredOrSkipped(
	ctx interfaces.UnifiedContext, stateExeId string, timerIdx int, cancelWaiting *bool,
) service.InternalTimerStatus {
	timerInfos := t.stateExecutionCurrentTimerInfos[stateExeId]
	if len(timerInfos) == 0 {
		if *cancelWaiting {
			// The waiting thread is later than the timer execState thread
			// The execState thread got completed early and call RemovePendingTimersOfState to remove the timerInfos
			// returning pending here
			return service.TimerPending
		} else {
			panic("bug: this shouldn't happen")
		}
	}
	timer := timerInfos[timerIdx]
	if timer.Status == service.TimerFired || timer.Status == service.TimerSkipped {
		return timer.Status
	}
	skippedByStaleSkip := t.RetryStaleSkipTimer()
	if skippedByStaleSkip {
		t.logger.Warn("timer skipped by stale skip signal", stateExeId, timerIdx)
		return service.TimerSkipped
	}
	now := t.provider.Now(ctx).Unix()
	fireAt := timer.FiringUnixTimestampSeconds
	duration := time.Duration(fireAt-now) * time.Second
	future := t.provider.NewTimer(ctx, duration)
	_ = t.provider.Await(ctx, func() bool {
		return future.IsReady() || timer.Status == service.TimerSkipped || *cancelWaiting
	})
	if timer.Status == service.TimerSkipped {
		return service.TimerSkipped
	}
	if future.IsReady() {
		return service.TimerFired
	}
	// otherwise *cancelWaiting should return false to indicate that this timer isn't completed(fired or skipped)
	return service.TimerPending
}

// RemovePendingTimersOfState is for when a state is completed, remove all its pending timers
func (t *SimpleTimerProcessor) RemovePendingTimersOfState(stateExeId string) {
	delete(t.stateExecutionCurrentTimerInfos, stateExeId)
}

func (t *SimpleTimerProcessor) AddTimers(
	stateExeId string, commands []iwfidl.TimerCommand, completedTimerCmds map[int]service.InternalTimerStatus,
) {
	timers := make([]*service.TimerInfo, len(commands))
	for idx, cmd := range commands {
		var timer service.TimerInfo
		if status, ok := completedTimerCmds[idx]; ok {
			timer = service.TimerInfo{
				CommandId:                  cmd.CommandId,
				FiringUnixTimestampSeconds: cmd.GetFiringUnixTimestampSeconds(),
				Status:                     status,
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
