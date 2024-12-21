package timers

import (
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

func removeElement(s []service.StaleSkipTimerSignal, i int) []service.StaleSkipTimerSignal {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// FixTimerCommandFromActivityOutput converts the durationSeconds to firingUnixTimestampSeconds
// doing it right after the activity output so that we don't need to worry about the time drift after continueAsNew
func FixTimerCommandFromActivityOutput(now time.Time, request iwfidl.CommandRequest) iwfidl.CommandRequest {
	var timerCommands []iwfidl.TimerCommand
	for _, cmd := range request.GetTimerCommands() {
		if cmd.HasDurationSeconds() {
			timerCommands = append(timerCommands, iwfidl.TimerCommand{
				CommandId:                  cmd.CommandId,
				FiringUnixTimestampSeconds: iwfidl.PtrInt64(now.Unix() + int64(cmd.GetDurationSeconds())),
			})
		} else {
			timerCommands = append(timerCommands, cmd)
		}
	}
	request.TimerCommands = timerCommands
	return request
}
