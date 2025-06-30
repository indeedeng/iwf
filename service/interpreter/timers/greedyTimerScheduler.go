package timers

import (
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter/cont"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"time"
)

type timerScheduler struct {
	// Timers requested by the workflow in desc order
	pendingScheduling []*service.TimerInfo
	// timers created through the workflow provider that are going to fire in desc order
	providerScheduledTimerUnixTs []int64
}

func (t *timerScheduler) addTimer(toAdd *service.TimerInfo) {

	if toAdd == nil || toAdd.Status != service.TimerPending {
		panic("invalid timer added")
	}

	insertIndex := 0
	for i, timer := range t.pendingScheduling {
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

	front := t.pendingScheduling[:insertIndex]
	var back []*service.TimerInfo
	if insertIndex >= len(t.pendingScheduling) {
		back = []*service.TimerInfo{toAdd}
	} else {
		back = append([]*service.TimerInfo{toAdd}, t.pendingScheduling[insertIndex:]...)
	}
	t.pendingScheduling = append(front, back...)
}

func (t *timerScheduler) removeTimer(toRemove *service.TimerInfo) {
	for i, timer := range t.pendingScheduling {
		if toRemove == timer {
			t.pendingScheduling = append(t.pendingScheduling[:i], t.pendingScheduling[i+1:]...)
			return
		}
	}
}

func (t *timerScheduler) pruneToNextTimer(pruneTo int64) *service.TimerInfo {
	index := len(t.providerScheduledTimerUnixTs)
	for i := len(t.providerScheduledTimerUnixTs) - 1; i >= 0; i-- {
		timerTime := t.providerScheduledTimerUnixTs[i]
		if timerTime > pruneTo {
			break
		}
		index = i
	}
	// If index is 0, it means all times are in the past
	if index == 0 {
		t.providerScheduledTimerUnixTs = nil
	} else {
		t.providerScheduledTimerUnixTs = t.providerScheduledTimerUnixTs[:index]
	}

	if len(t.pendingScheduling) == 0 {
		return nil
	}

	index = len(t.pendingScheduling)

	for i := len(t.pendingScheduling) - 1; i >= 0; i-- {
		timer := t.pendingScheduling[i]
		if timer.FiringUnixTimestampSeconds > pruneTo && timer.Status == service.TimerPending {
			break
		}
		index = i
	}

	// If index is 0, it means all timers are pruned
	if index == 0 {
		t.pendingScheduling = nil
		return nil
	}

	prunedTimer := t.pendingScheduling[index-1]
	t.pendingScheduling = t.pendingScheduling[:index]
	return prunedTimer
}

func startGreedyTimerScheduler(
	ctx interfaces.UnifiedContext,
	provider interfaces.WorkflowProvider,
	continueAsNewCounter *cont.ContinueAsNewCounter) *timerScheduler {

	t := timerScheduler{}
	provider.GoNamed(ctx, "greedy-timer-scheduler", func(ctx interfaces.UnifiedContext) {
		for {
			err := provider.Await(ctx, func() bool {
				now := provider.Now(ctx).Unix()
				next := t.pruneToNextTimer(now)
				return (next != nil && (len(t.providerScheduledTimerUnixTs) == 0 || next.FiringUnixTimestampSeconds < t.providerScheduledTimerUnixTs[len(t.providerScheduledTimerUnixTs)-1])) || continueAsNewCounter.IsThresholdMet()
			})

			if err != nil {
				break
			}

			if continueAsNewCounter.IsThresholdMet() {
				break
			}

			now := provider.Now(ctx).Unix()
			next := t.pruneToNextTimer(now)
			fireAt := next.FiringUnixTimestampSeconds
			duration := time.Duration(fireAt-now) * time.Second
			// This will create a new timer but not yield the goroutines awaiting the timer firing.
			// This works since when a timer fires, a new workflow task is created with the expectation that
			//   there is a goroutines awaiting some condition(some time has past) to continue,
			//   see WaitForTimerFiredOrSkipped.
			provider.NewTimer(ctx, duration)
			t.providerScheduledTimerUnixTs = append(t.providerScheduledTimerUnixTs, fireAt)
		}
	})

	return &t
}
