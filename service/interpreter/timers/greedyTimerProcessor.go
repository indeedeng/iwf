package timers

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter/cont"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
)

type GreedyTimerProcessor struct {
	timerManager                    *timerScheduler
	stateExecutionCurrentTimerInfos map[string][]*service.TimerInfo
	staleSkipTimerSignals           []service.StaleSkipTimerSignal
	provider                        interfaces.WorkflowProvider
	logger                          interfaces.UnifiedLogger
}

func NewGreedyTimerProcessor(
	ctx interfaces.UnifiedContext,
	provider interfaces.WorkflowProvider,
	continueAsNewCounter *cont.ContinueAsNewCounter,
	staleSkipTimerSignals []service.StaleSkipTimerSignal,
) *GreedyTimerProcessor {

	// start some single thread that manages pendingScheduling
	scheduler := startGreedyTimerScheduler(ctx, provider, continueAsNewCounter)

	tp := &GreedyTimerProcessor{
		provider:                        provider,
		timerManager:                    scheduler,
		stateExecutionCurrentTimerInfos: map[string][]*service.TimerInfo{},
		logger:                          provider.GetLogger(ctx),
		staleSkipTimerSignals:           staleSkipTimerSignals,
	}

	return tp
}

func (t *GreedyTimerProcessor) Dump() []service.StaleSkipTimerSignal {
	return t.staleSkipTimerSignals
}

func (t *GreedyTimerProcessor) GetTimerInfos() map[string][]*service.TimerInfo {
	return t.stateExecutionCurrentTimerInfos
}

func (t *GreedyTimerProcessor) GetTimerStartedUnixTimestamps() []int64 {
	return t.timerManager.providerScheduledTimerUnixTs
}

// SkipTimer will attempt to skip a timer, return false if no valid timer found
func (t *GreedyTimerProcessor) SkipTimer(stateExeId, timerId string, timerIdx int) bool {
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

func (t *GreedyTimerProcessor) RetryStaleSkipTimer() bool {
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
func (t *GreedyTimerProcessor) WaitForTimerFiredOrSkipped(
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
		timer.Status = service.TimerSkipped
		return service.TimerSkipped
	}

	_ = t.provider.Await(ctx, func() bool {
		// This is trigger when one of the timers scheduled by the timerScheduler fires, scheduling a
		//   new workflow task that evaluates the workflow's goroutines
		return timer.Status == service.TimerFired || timer.Status == service.TimerSkipped || timer.FiringUnixTimestampSeconds <= t.provider.Now(ctx).Unix() || *cancelWaiting
	})

	if timer.Status == service.TimerSkipped {
		return service.TimerSkipped
	}

	if timer.FiringUnixTimestampSeconds <= t.provider.Now(ctx).Unix() {
		timer.Status = service.TimerFired
		return service.TimerFired
	}

	// otherwise *cancelWaiting should return false to indicate that this timer isn't completed(fired or skipped)
	t.timerManager.removeTimer(timer)
	return service.TimerPending
}

// RemovePendingTimersOfState is for when a state is completed, remove all its pending pendingScheduling
func (t *GreedyTimerProcessor) RemovePendingTimersOfState(stateExeId string) {

	timers := t.stateExecutionCurrentTimerInfos[stateExeId]

	for _, timer := range timers {
		t.timerManager.removeTimer(timer)
	}

	delete(t.stateExecutionCurrentTimerInfos, stateExeId)
}

func (t *GreedyTimerProcessor) AddTimers(
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
		if timer.Status == service.TimerPending {
			t.timerManager.addTimer(&timer)
		}
		timers[idx] = &timer
	}
	t.stateExecutionCurrentTimerInfos[stateExeId] = timers
}
