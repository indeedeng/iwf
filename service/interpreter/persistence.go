package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type PersistenceManager struct {
	dataObjects             map[string]iwfidl.KeyValue
	searchAttributes        map[string]iwfidl.SearchAttribute
	searchAttributeUpserter func(attributes map[string]interface{}) error
}

func NewPersistenceManager(searchAttributeUpserter func(attributes map[string]interface{}) error) *PersistenceManager {
	return &PersistenceManager{
		dataObjects:             make(map[string]iwfidl.KeyValue),
		searchAttributes:        make(map[string]iwfidl.SearchAttribute),
		searchAttributeUpserter: searchAttributeUpserter,
	}
}

func (am *PersistenceManager) GetDataObjectsByKey(request service.GetDataObjectsQueryRequest) service.GetDataObjectsQueryResponse {
	all := false
	if len(request.Keys) == 0 {
		all = true
	}
	var res []iwfidl.KeyValue
	keyMap := map[string]bool{}
	for _, k := range request.Keys {
		keyMap[k] = true
	}
	for key, value := range am.dataObjects {
		if keyMap[key] || all {
			res = append(res, value)
		}
	}
	return service.GetDataObjectsQueryResponse{
		DataObjects: res,
	}
}

func (am *PersistenceManager) GetAllSearchAttributes() []iwfidl.SearchAttribute {
	var res []iwfidl.SearchAttribute
	for _, value := range am.searchAttributes {
		res = append(res, value)
	}
	return res
}

func (am *PersistenceManager) GetAllDataObjects() []iwfidl.KeyValue {
	var res []iwfidl.KeyValue
	for _, value := range am.dataObjects {
		res = append(res, value)
	}
	return res
}

func (am *PersistenceManager) ProcessUpsertSearchAttribute(attributes []iwfidl.SearchAttribute) error {
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

func (am *PersistenceManager) ProcessUpsertDataObject(attributes []iwfidl.KeyValue) error {
	for _, attr := range attributes {
		am.dataObjects[attr.GetKey()] = attr
	}
	return nil
}
