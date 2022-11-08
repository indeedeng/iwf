package service

const (
	TaskQueue                         = "Interpreter"
	GracefulCompletingWorkflowStateId = "_SYS_GRACEFUL_COMPLETING_WORKFLOW"
	ForceCompletingWorkflowStateId    = "_SYS_FORCE_COMPLETING_WORKFLOW"
	ForceFailingWorkflowStateId       = "_SYS_FORCE_FAILING_WORKFLOW"

	StateStartApi  = "/api/v1/workflowState/start"
	StateDecideApi = "/api/v1/workflowState/decide"

	DeciderTypeAllCommandCompleted = "ALL_COMMAND_COMPLETED"
	DeciderTypeAnyCommandCompleted = "ANY_COMMAND_COMPLETED"

	TimerStatusFired     = "FIRED"
	TimerStatusScheduled = "SCHEDULED"

	SignalStatusWaiting  = "WAITING"
	SignalStatusReceived = "RECEIVED"

	InternStateChannelCommandStatusWaiting = "WAITING"
	InternStateChannelCommandReceived      = "RECEIVED"

	SearchAttributeValueTypeKeyword = "KEYWORD"
	SearchAttributeValueTypeInt     = "INT"

	AttributeQueryType = "GetQueryAttributes"

	WorkflowErrorTypeUserWorkflowDecision = "UserWorkflowDecision"
	WorkflowErrorTypeUserWorkflowError    = "UserWorkflowError"
	WorkflowErrorTypeUserInternalError    = "InternalError"

	WorkflowStatusRunning       = "RUNNING"
	WorkflowStatusCompleted     = "COMPLETED"
	WorkflowStatusFailed        = "FAILED"
	WorkflowStatusTimeout       = "TIMEOUT"
	WorkflowStatusTerminated    = "TERMINATED"
	WorkflowStatusCanceled      = "CANCELED"
	WorkflowStatusContinueAsNew = "CONTINUED_AS_NEW"

	SearchAttributeGlobalVersion     = "GlobalWorkflowVersion"
	SearchAttributeExecutingStateIds = "ExecutingStateIds"
	SearchAttributeIwfWorkflowType   = "IwfWorkflowType"
)

type BackendType string

const BackendTypeCadence BackendType = "cadence"
const BackendTypeTemporal BackendType = "temporal"

type ResetType string

const ResetTypeHistoryEventId ResetType = "HISTORY_EVENT_ID"
const ResetTypeFirstDecisionCompleted ResetType = "FIRST_DECISION_COMPLETED"
const ResetTypeLastDecisionCompleted ResetType = "LAST_DECISION_COMPLETED"
const ResetTypeLastContinuedAsNew ResetType = "LAST_CONTINUED_AS_NEW"
const ResetTypeBadBinary ResetType = "BAD_BINARY"
const ResetTypeDecisionCompletedTime ResetType = "DECISION_COMPLETED_TIME"
const ResetTypeFirstDecisionScheduled ResetType = "FIRST_DECISION_SCHEDULED"
const ResetTypeLastDecisionScheduled ResetType = "LAST_DECISION_SCHEDULED"
