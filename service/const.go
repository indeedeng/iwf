package service

type (
	BackendType string
)

const (
	TaskQueue                         = "Interpreter_DEFAULT"
	GracefulCompletingWorkflowStateId = "_SYS_GRACEFUL_COMPLETING_WORKFLOW"
	ForceCompletingWorkflowStateId    = "_SYS_FORCE_COMPLETING_WORKFLOW"
	ForceFailingWorkflowStateId       = "_SYS_FORCE_FAILING_WORKFLOW"

	StateStartApi  = "/api/v1/workflowState/start"
	StateDecideApi = "/api/v1/workflowState/decide"

	GetDataObjectsWorkflowQueryType = "GetDataObjects"
	GetCurrentTimerInfosQueryType   = "GetCurrentTimerInfos"
	DumpAllInternalQueryType        = "DumpAllInternal"

	SkipTimerSignalChannelName = "IwfSystemSkipTimerChannel"

	SearchAttributeGlobalVersion     = "IwfGlobalWorkflowVersion"
	SearchAttributeExecutingStateIds = "IwfExecutingStateIds"
	SearchAttributeIwfWorkflowType   = "IwfWorkflowType"

	BackendTypeCadence  BackendType = "cadence"
	BackendTypeTemporal BackendType = "temporal"
)
