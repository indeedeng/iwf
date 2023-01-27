package interpreter

import (
	"github.com/indeedeng/iwf/service"
)

func SetQueryHandlersForContinueAsNew(ctx UnifiedContext, provider WorkflowProvider, interStateChannel *InterStateChannel) error {
	err := provider.SetQueryHandler(ctx, service.DumpAllInternalQueryType, func(request service.DumpAllInternalRequest) (*service.DumpAllInternalResponse, error) {
		// TODO use request for pagination
		interStateChannelReceived := interStateChannel.ReadReceived(nil)
		return &service.DumpAllInternalResponse{
			InterStateChannelReceived: interStateChannelReceived,
		}, nil
	})
	if err != nil {
		return err
	}
	return nil
}
