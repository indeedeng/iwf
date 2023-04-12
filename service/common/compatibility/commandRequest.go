package compatibility

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
)

func GetDeciderTriggerType(commandRequest iwfidl.CommandRequest) iwfidl.DeciderTriggerType {
	if commandRequest.HasCommandWaitingType() {
		newType := commandRequest.GetCommandWaitingType()
		switch newType {
		case iwfidl.ALL_COMPLETED:
			return iwfidl.ALL_COMMAND_COMPLETED
		case iwfidl.ANY_COMPLETED:
			return iwfidl.ANY_COMMAND_COMPLETED
		case iwfidl.ANY_COMBINATION_COMPLETED:
			return iwfidl.ANY_COMMAND_COMBINATION_COMPLETED
		default:
			panic("invalid waiting type to convert:" + string(newType))
		}
	}
	return commandRequest.GetDeciderTriggerType()
}
