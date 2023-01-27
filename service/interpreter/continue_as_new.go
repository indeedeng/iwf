package interpreter

import (
	"github.com/indeedeng/iwf/service"
)

func setQueryHandlersForContinueAsNew(ctx UnifiedContext, provider WorkflowProvider, interStateChannel *InterStateChannel) error {
	err := provider.SetQueryHandler(ctx, service.DumpAllInternalQueryType, func(request service.DumpAllInfoRequest) (*service.DumpAllInfoResponse, error) {
		// TODO use request for pagination
		interStateChannalReceived := interStateChannel.ReadReceived(nil)
		return &service.DumpAllInfoResponse{
			InterStateChannelReceived: interStateChannalReceived,
		}, nil
	})
	if err != nil {
		return err
	}
	return nil
}
