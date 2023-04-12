package compatibility

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/ptr"
)

func GetDeciderTriggerType(commandRequest iwfidl.CommandRequest) iwfidl.DeciderTriggerType {
	if commandRequest.HasCommandWaitingType() {
		newType := commandRequest.GetCommandWaitingType()
		switch newType {
		case iwfidl.ALL_COMMAND_COMPLETED:
			return ptr.Any(iwfidl.ALL_COMMAND_COMPLETED)
		}
	}
}
