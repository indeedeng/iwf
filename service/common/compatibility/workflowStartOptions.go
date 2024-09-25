package compatibility

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
)

func GetWorkflowIdReusePolicy(options iwfidl.WorkflowStartOptions) *iwfidl.WorkflowIDReusePolicy {
	if options.HasIdReusePolicy() {
		newType := options.GetIdReusePolicy()
		switch newType {
		case iwfidl.ALLOW_IF_NO_RUNNING:
			return iwfidl.ALLOW_DUPLICATE.Ptr()
		case iwfidl.ALLOW_IF_PREVIOUS_EXISTS_ABNORMALLY, iwfidl.ALLOW_IF_PREVIOUS_EXITS_ABNORMALLY:
			// Keeping typo enum for backwards compatibility. Both old and corrected enums return the same result.
			return iwfidl.ALLOW_DUPLICATE_FAILED_ONLY.Ptr()
		case iwfidl.DISALLOW_REUSE:
			return iwfidl.REJECT_DUPLICATE.Ptr()
		case iwfidl.ALLOW_TERMINATE_IF_RUNNING:
			return iwfidl.TERMINATE_IF_RUNNING.Ptr()
		default:
			panic("invalid id reuse policy to convert:" + string(newType))
		}
	}
	return options.WorkflowIDReusePolicy
}
