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

	GetDataObjectsWorkflowQueryType        = "GetDataObjects"
	GetSearchAttributesWorkflowQueryType   = "GetSearchAttributes"
	GetCurrentTimerInfosQueryType          = "GetCurrentTimerInfos"
	DumpAllInternalQueryType               = "DumpAllInternal"
	DumpAllInternalWithPaginationQueryType = "DumpAllInternalWithPagination"

	SearchAttributeGlobalVersion     = "IwfGlobalWorkflowVersion"
	SearchAttributeExecutingStateIds = "IwfExecutingStateIds"
	SearchAttributeIwfWorkflowType   = "IwfWorkflowType"

	BackendTypeCadence  BackendType = "cadence"
	BackendTypeTemporal BackendType = "temporal"

	IwfSystemSignalPrefix          = "__IwfSystem_"
	SkipTimerSignalChannelName     = IwfSystemSignalPrefix + "SkipTimerChannel"
	FailWorkflowSignalChanncelName = IwfSystemSignalPrefix + "FailWorkflowChannel"

	DefaultContinueAsNewPageSizeInBytes = 1024 * 1024
)

var ValidIwfSystemSignalNames map[string]bool = map[string]bool{
	SkipTimerSignalChannelName:     true,
	FailWorkflowSignalChanncelName: true,
}
