package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/mapper"
)

type PersistenceManager struct {
	dataObjects      map[string]iwfidl.KeyValue
	searchAttributes map[string]iwfidl.SearchAttribute
	provider         WorkflowProvider

	lockedDataObjectKeys      map[string]bool
	lockedSearchAttributeKeys map[string]bool

	useMemo bool
}

func NewPersistenceManager(provider WorkflowProvider, initSearchAttributes []iwfidl.SearchAttribute, useMemo bool) *PersistenceManager {
	searchAttributes := make(map[string]iwfidl.SearchAttribute)
	for _, sa := range initSearchAttributes {
		searchAttributes[sa.GetKey()] = sa
	}
	return &PersistenceManager{
		dataObjects:      make(map[string]iwfidl.KeyValue),
		searchAttributes: searchAttributes,
		provider:         provider,

		lockedDataObjectKeys:      make(map[string]bool),
		lockedSearchAttributeKeys: make(map[string]bool),

		useMemo: useMemo,
	}
}

func RebuildPersistenceManager(provider WorkflowProvider,
	dolist []iwfidl.KeyValue, salist []iwfidl.SearchAttribute,
	useMemo bool,
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

		// locks will not be carried over during continueAsNew
		lockedDataObjectKeys:      make(map[string]bool),
		lockedSearchAttributeKeys: make(map[string]bool),

		useMemo: useMemo,
	}
}

func (am *PersistenceManager) GetDataObjectsByKey(request service.GetDataObjectsQueryRequest) service.GetDataObjectsQueryResponse {
	var dataAttributes []iwfidl.KeyValue

	if len(request.Keys) == 0 {
		dataAttributes = am.GetAllDataObjects()
	} else {
		dataAttributes = LoadDataAttributes(request.Keys, am.GetAllDataObjects())
	}

	return service.GetDataObjectsQueryResponse{
		DataObjects: dataAttributes,
	}
}

func (am *PersistenceManager) LoadSearchAttributes(ctx UnifiedContext, loadingPolicy *iwfidl.PersistenceLoadingPolicy) []iwfidl.SearchAttribute {
	var loadingType iwfidl.PersistenceLoadingType
	var partialLoadingKeys []string
	if loadingPolicy != nil {
		loadingType = loadingPolicy.GetPersistenceLoadingType()
		partialLoadingKeys = loadingPolicy.PartialLoadingKeys

		if loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
			am.awaitAndLockForKeys(ctx, am.lockedSearchAttributeKeys, loadingPolicy.GetLockingKeys())
		}
	}

	if loadingType == "" || loadingType == iwfidl.ALL_WITHOUT_LOCKING {
		return am.GetAllSearchAttributes()
	} else if loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING || loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
		return LoadSearchAttributes(partialLoadingKeys, am.GetAllSearchAttributes())
	} else if loadingType == iwfidl.NONE {
		return []iwfidl.SearchAttribute{}
	} else {
		panic("not supported loading type " + loadingType)
	}
}

func (am *PersistenceManager) LoadDataObjects(ctx UnifiedContext, loadingPolicy *iwfidl.PersistenceLoadingPolicy) []iwfidl.KeyValue {
	var loadingType iwfidl.PersistenceLoadingType
	var partialLoadingKeys []string
	if loadingPolicy != nil {
		loadingType = loadingPolicy.GetPersistenceLoadingType()
		partialLoadingKeys = loadingPolicy.PartialLoadingKeys

		if loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
			am.awaitAndLockForKeys(ctx, am.lockedDataObjectKeys, loadingPolicy.GetLockingKeys())
		}
	}

	if loadingType == "" || loadingType == iwfidl.ALL_WITHOUT_LOCKING {
		return am.GetAllDataObjects()
	} else if loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING || loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
		return LoadDataAttributes(partialLoadingKeys, am.GetAllDataObjects())
	} else if loadingType == iwfidl.NONE {
		return []iwfidl.KeyValue{}
	} else {
		panic("not supported loading type " + loadingType)
	}
}

func LoadSearchAttributes(loadingKeys []string, allSearchAttributes []iwfidl.SearchAttribute) []iwfidl.SearchAttribute {
	if len(loadingKeys) == 0 {
		return []iwfidl.SearchAttribute{}
	}

	keyMap := map[string]bool{}
	for _, loadingKey := range loadingKeys {
		keyMap[loadingKey] = true
	}

	var results []iwfidl.SearchAttribute

	for _, sa := range allSearchAttributes {
		if keyMap[sa.GetKey()] {
			results = append(results, sa)
		}
	}

	return results
}

func LoadDataAttributes(loadingKeys []string, allDataAttributes []iwfidl.KeyValue) []iwfidl.KeyValue {
	if len(loadingKeys) == 0 {
		return []iwfidl.KeyValue{}
	}

	keyMap := map[string]bool{}
	for _, loadingKey := range loadingKeys {
		keyMap[loadingKey] = true
	}

	var results []iwfidl.KeyValue

	for _, da := range allDataAttributes {
		if keyMap[da.GetKey()] {
			results = append(results, da)
		}
	}

	return results
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

func (am *PersistenceManager) ProcessUpsertDataObject(ctx UnifiedContext, attributes []iwfidl.KeyValue) error {
	if len(attributes) == 0 {
		return nil
	}
	for _, attr := range attributes {
		am.dataObjects[attr.GetKey()] = attr
	}
	if am.useMemo {
		memo := map[string]iwfidl.EncodedObject{}
		for _, att := range attributes {
			memo[att.GetKey()] = att.GetValue()
		}
		return am.provider.UpsertMemo(ctx, memo)
	}
	return nil
}

func (am *PersistenceManager) awaitAndLockForKeys(ctx UnifiedContext, lockedKeys map[string]bool, keysToLock []string) {
	// wait until all keys are not locked
	err := am.provider.Await(ctx, func() bool {
		for _, k := range keysToLock {
			if lockedKeys[k] {
				return false
			}
		}
		return true
	})
	if err != nil {
		return
	}
	// then lock the keys
	for _, k := range keysToLock {
		lockedKeys[k] = true
	}
}

func (am *PersistenceManager) unlockKeys(lockedKeys map[string]bool, keysToUnlock []string) {
	for _, k := range keysToUnlock {
		delete(lockedKeys, k)
	}
}

func (am *PersistenceManager) UnlockPersistence(saPolicy *iwfidl.PersistenceLoadingPolicy, daPolicy *iwfidl.PersistenceLoadingPolicy) {
	if saPolicy != nil && saPolicy.GetPersistenceLoadingType() == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
		am.unlockKeys(am.lockedSearchAttributeKeys, saPolicy.GetLockingKeys())
	}

	if daPolicy != nil && daPolicy.GetPersistenceLoadingType() == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
		am.unlockKeys(am.lockedDataObjectKeys, daPolicy.GetLockingKeys())
	}
}
