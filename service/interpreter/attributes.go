package interpreter

import (
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
)

type AttributeManager struct {
	queryAttributes         map[string]iwfidl.KeyValue
	searchAttributes        map[string]iwfidl.SearchAttribute
	searchAttributeUpserter func(attributes map[string]interface{}) error
}

func NewAttributeManager(searchAttributeUpserter func(attributes map[string]interface{}) error) *AttributeManager {
	return &AttributeManager{
		queryAttributes:         make(map[string]iwfidl.KeyValue),
		searchAttributes:        make(map[string]iwfidl.SearchAttribute),
		searchAttributeUpserter: searchAttributeUpserter,
	}
}

func (am *AttributeManager) GetQueryAttributesByKey(request service.QueryAttributeRequest) service.QueryAttributeResponse {
	all := false
	if len(request.Keys) == 0 {
		all = true
	}
	var res []iwfidl.KeyValue
	keyMap := map[string]bool{}
	for _, k := range request.Keys {
		keyMap[k] = true
	}
	for key, value := range am.queryAttributes {
		if keyMap[key] || all {
			res = append(res, value)
		}
	}
	return service.QueryAttributeResponse{
		AttributeValues: res,
	}
}

func (am *AttributeManager) GetAllSearchAttributes() []iwfidl.SearchAttribute {
	var res []iwfidl.SearchAttribute
	for _, value := range am.searchAttributes {
		res = append(res, value)
	}
	return res
}

func (am *AttributeManager) GetAllQueryAttributes() []iwfidl.KeyValue {
	var res []iwfidl.KeyValue
	for _, value := range am.queryAttributes {
		res = append(res, value)
	}
	return res
}

func (am *AttributeManager) ProcessUpsertSearchAttribute(attributes []iwfidl.SearchAttribute) error {
	if len(attributes) == 0 {
		return nil
	}
	attrsToUpsert := map[string]interface{}{}
	for _, attr := range attributes {
		am.searchAttributes[attr.GetKey()] = attr
		switch attr.GetValueType() {
		case service.SearchAttributeValueTypeKeyword:
			attrsToUpsert[attr.GetKey()] = attr.GetStringValue()
		case service.SearchAttributeValueTypeInt:
			num := attr.GetIntegerValue()
			attrsToUpsert[attr.GetKey()] = num
		default:
			return fmt.Errorf("unsupported search attribute value type %v", attr.GetValueType())
		}
	}
	return am.searchAttributeUpserter(attrsToUpsert)
}

func (am *AttributeManager) ProcessUpsertQueryAttribute(attributes []iwfidl.KeyValue) error {
	for _, attr := range attributes {
		am.queryAttributes[attr.GetKey()] = attr
	}
	return nil
}
