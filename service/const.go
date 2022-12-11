package service

const (
	TaskQueue                         = "Interpreter_DEFAULT"
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

	GetDataObjectsWorkflowQueryType = "GetDataObjects"

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

	SearchAttributeGlobalVersion     = "IwfGlobalWorkflowVersion"
	SearchAttributeExecutingStateIds = "IwfExecutingStateIds"
	SearchAttributeIwfWorkflowType   = "IwfWorkflowType"

	WorkflowIDReusePolicyAllowDuplicateFailedOnly = "ALLOW_DUPLICATE_FAILED_ONLY"
	WorkflowIDReusePolicyAllowDuplicate           = "ALLOW_DUPLICATE"
	WorkflowIDReusePolicyRejectDuplicate          = "REJECT_DUPLICATE"
	WorkflowIDReusePolicyTerminateIfRunning       = "TERMINATE_IF_RUNNING"
)

type BackendType string

const BackendTypeCadence BackendType = "cadence"
const BackendTypeTemporal BackendType = "temporal"

type ResetType string

const ResetTypeHistoryEventId ResetType = "HISTORY_EVENT_ID"
const ResetTypeBeginning ResetType = "BEGINNING"
const ResetTypeHistoryEventTime ResetType = "HISTORY_EVENT_TIME"
