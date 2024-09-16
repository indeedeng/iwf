package service

type (
	BackendType string
)

const (
	EnvNameDebugMode = "DEBUG_MODE"

	DefaultContinueAsNewPageSizeInBytes = 1024 * 1024

	// below are special unofficial code for special use case

	// HttpStatusCodeSpecial4xxError1 is for poll timeout, RPC worker execution error
	HttpStatusCodeSpecial4xxError1 = 420
	// HttpStatusCodeSpecial4xxError2 is for RPC acquire locking failure
	HttpStatusCodeSpecial4xxError2 = 450

	TaskQueue = "Interpreter_DEFAULT"

	StateStartApi        = "/api/v1/workflowState/start"
	StateDecideApi       = "/api/v1/workflowState/decide"
	WorkflowWorkerRpcApi = "/api/v1/workflowWorker/rpc"

	GetDataAttributesWorkflowQueryType   = "GetDataAttributes"
	GetSearchAttributesWorkflowQueryType = "GetSearchAttributes"
	GetCurrentTimerInfosQueryType        = "GetCurrentTimerInfos"
	ContinueAsNewDumpQueryType           = "ContinueAsNewDump"
	DebugDumpQueryType                   = "DebugNewDump"
	PrepareRpcQueryType                  = "PrepareRpcQueryType"

	ExecuteOptimisticLockingRpcUpdateType = "ExecuteOptimisticLockingRpcUpdate"

	SearchAttributeGlobalVersion     = "IwfGlobalWorkflowVersion"
	SearchAttributeExecutingStateIds = "IwfExecutingStateIds"
	SearchAttributeIwfWorkflowType   = "IwfWorkflowType"

	BackendTypeCadence  BackendType = "cadence"
	BackendTypeTemporal BackendType = "temporal"

	IwfSystemConstPrefix = "__IwfSystem_"

	SkipTimerSignalChannelName            = IwfSystemConstPrefix + "SkipTimerChannel"
	FailWorkflowSignalChannelName         = IwfSystemConstPrefix + "FailWorkflowChannel"
	UpdateConfigSignalChannelName         = IwfSystemConstPrefix + "UpdateWorkflowConfig"
	ExecuteRpcSignalChannelName           = IwfSystemConstPrefix + "ExecuteRpc"
	StateCompletionSignalChannelName      = IwfSystemConstPrefix + "StateCompletion"
	TriggerContinueAsNewSignalChannelName = IwfSystemConstPrefix + "TriggerContinueAsNew"

	WorkerUrlMemoKey            = IwfSystemConstPrefix + "WorkerUrl"
	UseMemoForDataAttributesKey = IwfSystemConstPrefix + "UseMemoForDataAttributes"
)

var ValidIwfSystemSignalNames = map[string]bool{
	SkipTimerSignalChannelName:    true,
	FailWorkflowSignalChannelName: true,
	UpdateConfigSignalChannelName: true,
	ExecuteRpcSignalChannelName:   true,
}

const (
	GracefulCompletingWorkflowStateId = "_SYS_GRACEFUL_COMPLETING_WORKFLOW"
	ForceCompletingWorkflowStateId    = "_SYS_FORCE_COMPLETING_WORKFLOW"
	ForceFailingWorkflowStateId       = "_SYS_FORCE_FAILING_WORKFLOW"
	DeadEndWorkflowStateId            = "_SYS_DEAD_END"
)

var ValidClosingWorkflowStateId = map[string]bool{
	GracefulCompletingWorkflowStateId: true,
	ForceCompletingWorkflowStateId:    true,
	ForceFailingWorkflowStateId:       true,
	DeadEndWorkflowStateId:            true,
}
