package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter/config"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
)

func SetQueryHandlers(
	ctx interfaces.UnifiedContext, provider interfaces.WorkflowProvider, persistenceManager *PersistenceManager,
	internalChannel *InternalChannel, signalReceiver *SignalReceiver,
	continueAsNewer *ContinueAsNewer,
	workflowConfiger *config.WorkflowConfiger, basicInfo service.BasicInfo,
) error {
	err := provider.SetQueryHandler(ctx, service.GetDataAttributesWorkflowQueryType, func(req service.GetDataAttributesQueryRequest) (service.GetDataAttributesQueryResponse, error) {
		dos := persistenceManager.GetDataObjectsByKey(req)
		return dos, nil
	})
	if err != nil {
		return err
	}
	err = provider.SetQueryHandler(ctx, service.GetSearchAttributesWorkflowQueryType, func() ([]iwfidl.SearchAttribute, error) {
		return persistenceManager.GetAllSearchAttributes(), nil
	})
	if err != nil {
		return err
	}
	err = continueAsNewer.SetQueryHandlersForContinueAsNew(ctx)
	if err != nil {
		return err
	}
	err = provider.SetQueryHandler(ctx, service.DebugDumpQueryType, func() (*service.DebugDumpResponse, error) {
		return &service.DebugDumpResponse{
			Config:   workflowConfiger.Get(),
			Snapshot: continueAsNewer.GetSnapshot(),
		}, nil
	})
	if err != nil {
		return err
	}
	err = provider.SetQueryHandler(ctx, service.PrepareRpcQueryType, func(req service.PrepareRpcQueryRequest) (service.PrepareRpcQueryResponse, error) {
		info := provider.GetWorkflowInfo(ctx) // TODO use firstRunId instead

		return service.PrepareRpcQueryResponse{
			DataObjects:              persistenceManager.LoadDataObjects(ctx, req.DataObjectsLoadingPolicy),
			SearchAttributes:         persistenceManager.LoadSearchAttributes(ctx, req.SearchAttributesLoadingPolicy),
			WorkflowRunId:            info.WorkflowExecution.RunID,
			WorkflowStartedTimestamp: info.WorkflowStartTime.Unix(),
			IwfWorkflowType:          basicInfo.IwfWorkflowType,
			IwfWorkerUrl:             basicInfo.IwfWorkerUrl,
			SignalChannelInfo:        signalReceiver.GetInfos(),
			InternalChannelInfo:      internalChannel.GetInfos(),
		}, nil
	})
	if err != nil {
		return err
	}
	return nil
}
