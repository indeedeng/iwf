package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/event"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"time"
)

type WorkflowUpdater struct {
	persistenceManager   *PersistenceManager
	provider             interfaces.WorkflowProvider
	continueAsNewer      *ContinueAsNewer
	continueAsNewCounter *ContinueAsNewCounter
	internalChannel      *InternalChannel
	signalReceiver       *SignalReceiver
	stateRequestQueue    *StateRequestQueue
	configer             *WorkflowConfiger
	logger               interfaces.UnifiedLogger
	basicInfo            service.BasicInfo
	globalVersioner      *GlobalVersioner
}

func NewWorkflowUpdater(
	ctx interfaces.UnifiedContext, provider interfaces.WorkflowProvider, persistenceManager *PersistenceManager,
	stateRequestQueue *StateRequestQueue,
	continueAsNewer *ContinueAsNewer, continueAsNewCounter *ContinueAsNewCounter, configer *WorkflowConfiger,
	internalChannel *InternalChannel, signalReceiver *SignalReceiver, basicInfo service.BasicInfo,
	globalVersioner *GlobalVersioner,
) (*WorkflowUpdater, error) {
	updater := &WorkflowUpdater{
		persistenceManager:   persistenceManager,
		continueAsNewer:      continueAsNewer,
		continueAsNewCounter: continueAsNewCounter,
		internalChannel:      internalChannel,
		signalReceiver:       signalReceiver,
		stateRequestQueue:    stateRequestQueue,
		configer:             configer,
		basicInfo:            basicInfo,
		provider:             provider,
		logger:               provider.GetLogger(ctx),
		globalVersioner:      globalVersioner,
	}
	if globalVersioner.IsAfterVersionOfTemporal26SDK() {
		err := provider.SetRpcUpdateHandler(ctx, service.ExecuteOptimisticLockingRpcUpdateType, updater.validator, updater.handler)
		if err != nil {
			return nil, err
		}
	}
	return updater, nil
}

func (u *WorkflowUpdater) handler(
	ctx interfaces.UnifiedContext, input iwfidl.WorkflowRpcRequest,
) (output *interfaces.HandlerOutput, err error) {
	u.continueAsNewer.IncreaseInflightOperation()
	defer u.continueAsNewer.DecreaseInflightOperation()

	info := u.provider.GetWorkflowInfo(ctx)

	defer func() {
		if !u.provider.IsReplaying(ctx) {
			event.Handle(iwfidl.IwfEvent{
				EventType:        iwfidl.RPC_EXECUTION_EVENT,
				RpcName:          &input.RpcName,
				WorkflowType:     u.basicInfo.IwfWorkflowType,
				WorkflowId:       info.WorkflowExecution.ID,
				SearchAttributes: u.persistenceManager.GetAllSearchAttributes(),
			})
		}
	}()

	rpcPrep := service.PrepareRpcQueryResponse{
		DataObjects:              u.persistenceManager.LoadDataObjects(ctx, input.DataAttributesLoadingPolicy),
		SearchAttributes:         u.persistenceManager.LoadSearchAttributes(ctx, input.SearchAttributesLoadingPolicy),
		WorkflowRunId:            info.WorkflowExecution.RunID,
		WorkflowStartedTimestamp: info.WorkflowStartTime.Unix(),
		IwfWorkflowType:          u.basicInfo.IwfWorkflowType,
		IwfWorkerUrl:             u.basicInfo.IwfWorkerUrl,
		SignalChannelInfo:        u.signalReceiver.GetInfos(),
		InternalChannelInfo:      u.internalChannel.GetInfos(),
	}

	activityOptions := interfaces.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttemptsDurationSeconds: input.TimeoutSeconds,
			MaximumAttempts:                iwfidl.PtrInt32(3),
		},
	}
	ctx = u.provider.WithActivityOptions(ctx, activityOptions)
	var activityOutput interfaces.InvokeRpcActivityOutput
	err = u.provider.ExecuteActivity(&activityOutput, u.configer.ShouldOptimizeActivity(), ctx,
		InvokeWorkerRpc, u.provider.GetBackendType(), rpcPrep, input)
	u.persistenceManager.UnlockPersistence(input.SearchAttributesLoadingPolicy, input.DataAttributesLoadingPolicy)

	if err != nil {
		return nil, u.provider.NewApplicationError(string(iwfidl.SERVER_INTERNAL_ERROR_TYPE), "activity invocation failure:"+err.Error())
	}

	handlerOutput := &interfaces.HandlerOutput{
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
		u.internalChannel.ProcessPublishing(rpcOutput.PublishToInterStateChannel)
		if rpcOutput.StateDecision != nil {
			u.stateRequestQueue.AddStateStartRequests(rpcOutput.StateDecision.NextStates)
		}
	}

	return handlerOutput, nil
}

func (u *WorkflowUpdater) validator(_ interfaces.UnifiedContext, input iwfidl.WorkflowRpcRequest) error {
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
