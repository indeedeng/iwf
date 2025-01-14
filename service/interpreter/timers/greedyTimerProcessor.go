package timers

import (
	"github.com/indeedeng/iwf/service/interpreter/cont"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type TimerManager struct {
	// Does not map to the timers actually created by the workflow provider
	PendingScheduling []*service.TimerInfo
	// timers created through the workflow provider that are going to fire
	ScheduledTimerTimes []int64
}

type GreedyTimerProcessor struct {
	timerManger                     TimerManager
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

	tp := &GreedyTimerProcessor{
		provider:                        provider,
		timerManger:                     TimerManager{},
		stateExecutionCurrentTimerInfos: map[string][]*service.TimerInfo{},
		logger:                          provider.GetLogger(ctx),
		staleSkipTimerSignals:           staleSkipTimerSignals,
	}

	// start some single thread that manages PendingScheduling
	tp.createGreedyTimerScheduler(ctx, continueAsNewCounter)

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

func (t *TimerManager) addTimer(toAdd *service.TimerInfo) {

	if toAdd == nil || toAdd.Status != service.TimerPending {
		panic("invalid timer added")
	}

	insertIndex := 0
	for i, timer := range t.PendingScheduling {
		if toAdd.FiringUnixTimestampSeconds >= timer.FiringUnixTimestampSeconds {
			// don't want dupes. Makes remove simpler
			if toAdd == timer {
				return
			}
			insertIndex = i
			break
		}
		insertIndex = i + 1
	}
	t.PendingScheduling = append(
		t.PendingScheduling[:insertIndex],
		append([]*service.TimerInfo{toAdd}, t.PendingScheduling[insertIndex:]...)...)
}

func (t *TimerManager) removeTimer(toRemove *service.TimerInfo) {
	for i, timer := range t.PendingScheduling {
		if toRemove == timer {
			t.PendingScheduling = append(t.PendingScheduling[:i], t.PendingScheduling[i+1:]...)
			return
		}
	}
}

func (t *TimerManager) pruneToNextTimer(pruneTo int64) *service.TimerInfo {

	if len(t.PendingScheduling) == 0 {
		return nil
	}

	index := len(t.PendingScheduling)

	for i := len(t.PendingScheduling) - 1; i >= 0; i-- {
		timer := t.PendingScheduling[i]
		if timer.FiringUnixTimestampSeconds > pruneTo && timer.Status == service.TimerPending {
			break
		}
		index = i
	}

	// If index is 0, it means all timers are pruned
	if index == 0 {
		t.PendingScheduling = nil
		return nil
	}

	prunedTimer := t.PendingScheduling[index-1]
	t.PendingScheduling = t.PendingScheduling[:index]
	return prunedTimer
}

func (t *GreedyTimerProcessor) createGreedyTimerScheduler(
	ctx interfaces.UnifiedContext,
	continueAsNewCounter *cont.ContinueAsNewCounter) {

	t.provider.GoNamed(ctx, "greedy-timer-scheduler", func(ctx interfaces.UnifiedContext) {
		for {
			_ = t.provider.Await(ctx, func() bool {
				// remove fired PendingScheduling
				now := t.provider.Now(ctx).Unix()
				for i := len(t.timerManger.ScheduledTimerTimes) - 1; i >= 0; i-- {
					if t.timerManger.ScheduledTimerTimes[i] > now {
						t.timerManger.ScheduledTimerTimes = t.timerManger.ScheduledTimerTimes[:i+1]
						break
					}
				}
				next := t.timerManger.pruneToNextTimer(now)
				return (next != nil && (len(t.timerManger.ScheduledTimerTimes) == 0 || next.FiringUnixTimestampSeconds < t.timerManger.ScheduledTimerTimes[len(t.timerManger.ScheduledTimerTimes)-1])) || continueAsNewCounter.IsThresholdMet()
			})

			if continueAsNewCounter.IsThresholdMet() {
				break
			}

			now := t.provider.Now(ctx).Unix()
			next := t.timerManger.pruneToNextTimer(now)
			//next := t.timerManger.getEarliestTimer()
			// only create a new timer when a pending timer exists before the next existing timer fires
			if next != nil && (len(t.timerManger.ScheduledTimerTimes) == 0 || next.FiringUnixTimestampSeconds < t.timerManger.ScheduledTimerTimes[len(t.timerManger.ScheduledTimerTimes)-1]) {
				fireAt := next.FiringUnixTimestampSeconds
				duration := time.Duration(fireAt-now) * time.Second
				t.provider.NewTimer(ctx, duration)
				t.timerManger.ScheduledTimerTimes = append(t.timerManger.ScheduledTimerTimes, fireAt)
			}
		}
	})
}

func (t *GreedyTimerProcessor) Dump() []service.StaleSkipTimerSignal {
	return t.staleSkipTimerSignals
}

func (t *GreedyTimerProcessor) GetCurrentTimerInfos() map[string][]*service.TimerInfo {
	return t.stateExecutionCurrentTimerInfos
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
	t.timerManger.removeTimer(timer)
	return service.TimerPending
}

// RemovePendingTimersOfState is for when a state is completed, remove all its pending PendingScheduling
func (t *GreedyTimerProcessor) RemovePendingTimersOfState(stateExeId string) {

	timers := t.stateExecutionCurrentTimerInfos[stateExeId]

	for _, timer := range timers {
		t.timerManger.removeTimer(timer)
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
			t.timerManger.addTimer(&timer)
		}
		timers[idx] = &timer
	}
	t.stateExecutionCurrentTimerInfos[stateExeId] = timers
}
