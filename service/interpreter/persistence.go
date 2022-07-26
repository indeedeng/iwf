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

func (am *PersistenceManager) LoadSearchAttributes(stateOptions *iwfidl.WorkflowStateOptions) []iwfidl.SearchAttribute {
	var loadingType iwfidl.PersistenceLoadingType
	var partialLoadingKeys []string
	if stateOptions != nil && stateOptions.SearchAttributesLoadingPolicy != nil {
		policy := stateOptions.GetSearchAttributesLoadingPolicy()
		loadingType = policy.GetPersistenceLoadingType()
		partialLoadingKeys = policy.PartialLoadingKeys
	}
	if loadingType == "" || loadingType == iwfidl.ALL_WITHOUT_LOCKING {
		return am.GetAllSearchAttributes()
	} else if loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING {
		var res []iwfidl.SearchAttribute
		keyMap := map[string]bool{}
		for _, k := range partialLoadingKeys {
			keyMap[k] = true
		}
		for key, value := range am.searchAttributes {
			if keyMap[key] {
				res = append(res, value)
			}
		}
		return res
	} else {
		panic("not supported loading type " + loadingType)
	}
}

func (am *PersistenceManager) LoadDataObjects(stateOptions *iwfidl.WorkflowStateOptions) []iwfidl.KeyValue {
	var loadingType iwfidl.PersistenceLoadingType
	var partialLoadingKeys []string
	if stateOptions != nil && stateOptions.DataObjectsLoadingPolicy != nil {
		policy := stateOptions.GetDataObjectsLoadingPolicy()
		loadingType = policy.GetPersistenceLoadingType()
		partialLoadingKeys = policy.PartialLoadingKeys
	}

	if loadingType == "" || loadingType == iwfidl.ALL_WITHOUT_LOCKING {
		return am.GetAllDataObjects()
	} else if loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING {
		res := am.GetDataObjectsByKey(service.GetDataObjectsQueryRequest{
			Keys: partialLoadingKeys,
		})
		return res.DataObjects
	} else {
		panic("not supported loading type " + loadingType)
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
		case iwfidl.KEYWORD:
			attrsToUpsert[attr.GetKey()] = attr.GetStringValue()
		case iwfidl.INT:
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
