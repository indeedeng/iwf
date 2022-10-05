package service

const (
	TaskQueue                         = "Interpreter"
	GracefulCompletingWorkflowStateId = "_SYS_GRACEFUL_COMPLETING_WORKFLOW"
	ForceCompletingWorkflowStateId    = "_SYS_FORCE_COMPLETING_WORKFLOW"
	ForceFailingWorkflowStateId       = "_SYS_FORCE_FAILING_WORKFLOW"

	StateStartApi  = "/api/v1/workflowState/start"
	StateDecideApi = "/api/v1/workflowState/decide"

	DeciderTypeAllCommandCompleted = "ALL_COMMAND_COMPLETED"

	TimerStatusFired     = "FIRED"
	TimerStatusScheduled = "SCHEDULED"

	SignalStatusWaiting  = "WAITING"
	SignalStatusReceived = "RECEIVED"

	SearchAttributeValueTypeKeyword = "KEYWORD"
	SearchAttributeValueTypeInt     = "INT"

	AttributeQueryType = "GetQueryAttributes"

	WorkflowErrorTypeUserWorkflowDecision = "UserWorkflowDecision"
	WorkflowErrorTypeUserWorkflowError    = "UserWorkflowError"
	WorkflowErrorTypeUserInternalError    = "InternalError"
)
