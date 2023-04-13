package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/compatibility"
	"github.com/indeedeng/iwf/service/common/mapper"
)

type PersistenceManager struct {
	dataObjects      map[string]iwfidl.KeyValue
	searchAttributes map[string]iwfidl.SearchAttribute
	provider         WorkflowProvider
}

func NewPersistenceManager(provider WorkflowProvider, initSearchAttributes []iwfidl.SearchAttribute) *PersistenceManager {
	searchAttributes := make(map[string]iwfidl.SearchAttribute)
	for _, sa := range initSearchAttributes {
		searchAttributes[sa.GetKey()] = sa
	}
	return &PersistenceManager{
		dataObjects:      make(map[string]iwfidl.KeyValue),
		searchAttributes: searchAttributes,
		provider:         provider,
	}
}

func RebuildPersistenceManager(provider WorkflowProvider,
	dolist []iwfidl.KeyValue, salist []iwfidl.SearchAttribute,
) *PersistenceManager {
	dataObjects := make(map[string]iwfidl.KeyValue)
	searchAttributes := make(map[string]iwfidl.SearchAttribute)
	for _, do := range dolist {
		dataObjects[do.GetKey()] = do
	}
	for _, sa := range salist {
		searchAttributes[sa.GetKey()] = sa
	}
	return &PersistenceManager{
		dataObjects:      dataObjects,
		searchAttributes: searchAttributes,
		provider:         provider,
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
	if stateOptions != nil && compatibility.GetDataObjectsLoadingPolicy(stateOptions) != nil {
		policy := compatibility.GetDataObjectsLoadingPolicy(stateOptions)
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

func (am *PersistenceManager) ProcessUpsertSearchAttribute(ctx UnifiedContext, attributes []iwfidl.SearchAttribute) error {
	if len(attributes) == 0 {
		return nil
	}

	for _, attr := range attributes {
		am.searchAttributes[attr.GetKey()] = attr
	}
	attrsToUpsert, err := mapper.MapToInternalSearchAttributes(attributes)
	if err != nil {
		return err
	}
	return am.provider.UpsertSearchAttributes(ctx, attrsToUpsert)
}

func (am *PersistenceManager) ProcessUpsertDataObject(attributes []iwfidl.KeyValue) error {
	for _, attr := range attributes {
		am.dataObjects[attr.GetKey()] = attr
	}
	return nil
}
