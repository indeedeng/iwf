package service

const (
	TaskQueue                 = "Interpreter"
	CompletingWorkflowStateId = "_SYS_COMPLETING_WORKFLOW"
	FailingWorkflowStateId    = "_SYS_FAILING_WORKFLOW"

	StateStartApi  = "/api/v1/workflowState/start"
	StateDecideApi = "/api/v1/workflowState/decide"

	DeciderTypeAllCommandCompleted = "ALL_COMMAND_COMPLETED"

	TimerStatusFired     = "FIRED"
	TimerStatusScheduled = "SCHEDULED"

	SignalStatusWaiting  = "WAITING"
	SignalStatusReceived = "RECEIVED"

	SearchAttributeValueTypeKeyword = "KEYWORD"
	SearchAttributeValueTypeInt     = "INT"
)
