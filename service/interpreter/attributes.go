package interpreter

import (
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
)

type AttributeManager struct {
	queryAttributes         map[string]iwfidl.EncodedObject
	searchAttributes        map[string]iwfidl.SearchAttribute
	searchAttributeUpserter func(attributes map[string]interface{}) error
}

func NewAttributeManager(searchAttributeUpserter func(attributes map[string]interface{}) error) *AttributeManager {
	return &AttributeManager{
		queryAttributes:         make(map[string]iwfidl.EncodedObject),
		searchAttributes:        make(map[string]iwfidl.SearchAttribute),
		searchAttributeUpserter: searchAttributeUpserter,
	}
}

func (am *AttributeManager) GetQueryAttributesByKey(request service.QueryAttributeRequest) service.QueryAttributeResponse {
	return service.QueryAttributeResponse{}
}

func (am *AttributeManager) GetAllSearchAttributes() []iwfidl.SearchAttribute {
	return nil
}

func (am *AttributeManager) GetAllQueryAttributes() []iwfidl.KeyValue {
	return nil
}

func (am *AttributeManager) ProcessUpsertSearchAttribute(attributes []iwfidl.SearchAttribute) {

}

func (am *AttributeManager) ProcessUpsertQueryAttribute(attributes []iwfidl.KeyValue) {

}
