package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

func setQueryHandlersForContinueAsNew(ctx UnifiedContext, provider WorkflowProvider, interStateChannel *InterStateChannel) error {
	err := provider.SetQueryHandler(ctx, service.GetInterStateChannelDataQueryType, func(channelNames []string) (map[string][]*iwfidl.EncodedObject, error) {
		return interStateChannel.ReadData(channelNames), nil
	})
	if err != nil {
		return err
	}
	return nil
}
