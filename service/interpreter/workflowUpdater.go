package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"time"
)

type WorkflowUpdater struct {
	persistenceManager   *PersistenceManager
	provider             WorkflowProvider
	continueAsNewer      *ContinueAsNewer
	continueAsNewCounter *ContinueAsNewCounter
	interStateChannel    *InterStateChannel
	stateRequestQueue    *StateRequestQueue
	logger               UnifiedLogger
	basicInfo            service.BasicInfo
}

func NewWorkflowUpdater(ctx UnifiedContext, provider WorkflowProvider, persistenceManager *PersistenceManager, stateRequestQueue *StateRequestQueue,
	continueAsNewer *ContinueAsNewer, continueAsNewCounter *ContinueAsNewCounter, interStateChannel *InterStateChannel, basicInfo service.BasicInfo,
) (*WorkflowUpdater, error) {
	updater := &WorkflowUpdater{
		persistenceManager:   persistenceManager,
		continueAsNewer:      continueAsNewer,
		continueAsNewCounter: continueAsNewCounter,
		interStateChannel:    interStateChannel,
		stateRequestQueue:    stateRequestQueue,
		basicInfo:            basicInfo,
		provider:             provider,
		logger:               provider.GetLogger(ctx),
	}
	err := provider.SetRpcUpdateHandler(ctx, service.ExecuteOptimisticLockingRpcUpdateType, updater.validator, updater.handler)
	if err != nil {
		return nil, err
	}
	return updater, nil
}

func (u *WorkflowUpdater) handler(ctx UnifiedContext, input iwfidl.WorkflowRpcRequest) (output *HandlerOutput, err error) {

	u.continueAsNewer.IncreaseInflightOperation()
	defer u.continueAsNewer.DecreaseInflightOperation()

	info := u.provider.GetWorkflowInfo(ctx)
	rpcPrep := service.PrepareRpcQueryResponse{
		DataObjects:              u.persistenceManager.LoadDataObjects(ctx, input.DataAttributesLoadingPolicy),
		SearchAttributes:         u.persistenceManager.LoadSearchAttributes(ctx, input.SearchAttributesLoadingPolicy),
		WorkflowRunId:            info.WorkflowExecution.RunID,
		WorkflowStartedTimestamp: info.WorkflowStartTime.Unix(),
		IwfWorkflowType:          u.basicInfo.IwfWorkflowType,
		IwfWorkerUrl:             u.basicInfo.IwfWorkerUrl,
	}

	activityOptions := ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttemptsDurationSeconds: input.TimeoutSeconds,
			MaximumAttempts:                iwfidl.PtrInt32(3),
		},
	}
	ctx = u.provider.WithActivityOptions(ctx, activityOptions)
	var activityOutput InvokeRpcActivityOutput
	err = u.provider.ExecuteActivity(ctx, InvokeWorkerRpc, u.provider.GetBackendType(), rpcPrep, input).Get(ctx, &activityOutput)
	u.persistenceManager.UnlockPersistence(input.DataAttributesLoadingPolicy, input.SearchAttributesLoadingPolicy)

	if err != nil {
		return nil, u.provider.NewApplicationError(string(iwfidl.SERVER_INTERNAL_ERROR_TYPE), "activity invocation failure:"+err.Error())
	}

	handlerOutput := &HandlerOutput{
		StatusError: activityOutput.StatusError,
	}
	rpcOutput := activityOutput.RpcOutput
	if rpcOutput != nil {
		handlerOutput.RpcOutput = &iwfidl.WorkflowRpcResponse{
			Output: rpcOutput.Output,
		}
		u.continueAsNewCounter.IncSyncUpdateReceived()
		_ = u.persistenceManager.ProcessUpsertDataObject(ctx, rpcOutput.UpsertDataAttributes)
		_ = u.persistenceManager.ProcessUpsertSearchAttribute(ctx, rpcOutput.UpsertSearchAttributes)
		u.interStateChannel.ProcessPublishing(rpcOutput.PublishToInterStateChannel)
		if rpcOutput.StateDecision != nil {
			u.stateRequestQueue.AddStateStartRequests(rpcOutput.StateDecision.NextStates)
		}
	}

	return handlerOutput, nil
}

func (u *WorkflowUpdater) validator(_ UnifiedContext, input iwfidl.WorkflowRpcRequest) error {
	var daKeys, saKeys []string
	if input.HasDataAttributesLoadingPolicy() {
		daKeys = input.DataAttributesLoadingPolicy.LockingKeys
	}
	if input.HasSearchAttributesLoadingPolicy() {
		saKeys = input.SearchAttributesLoadingPolicy.LockingKeys
	}
	keysUnlocked := u.persistenceManager.CheckDataAndSearchAttributesKeysAreUnlocked(daKeys, saKeys)
	if keysUnlocked {
		return nil
	} else {
		return u.provider.NewApplicationError(string(iwfidl.RPC_ACQUIRE_LOCK_FAILURE), "requested data or search attributes are being locked by other operations")
	}
}
