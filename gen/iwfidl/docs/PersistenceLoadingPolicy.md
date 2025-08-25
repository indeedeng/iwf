# PersistenceLoadingPolicy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PersistenceLoadingType** | Pointer to [**PersistenceLoadingType**](PersistenceLoadingType.md) |  | [optional] 
**PartialLoadingKeys** | Pointer to **[]string** |  | [optional] 
**LockingKeys** | Pointer to **[]string** |  | [optional] 
**UseKeyAsPrefix** | Pointer to **bool** |  | [optional] 
**LazyLoadingLargeDataAttributes** | Pointer to **bool** |  | [optional] 

## Methods

### NewPersistenceLoadingPolicy

`func NewPersistenceLoadingPolicy() *PersistenceLoadingPolicy`

NewPersistenceLoadingPolicy instantiates a new PersistenceLoadingPolicy object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPersistenceLoadingPolicyWithDefaults

`func NewPersistenceLoadingPolicyWithDefaults() *PersistenceLoadingPolicy`

NewPersistenceLoadingPolicyWithDefaults instantiates a new PersistenceLoadingPolicy object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPersistenceLoadingType

`func (o *PersistenceLoadingPolicy) GetPersistenceLoadingType() PersistenceLoadingType`

GetPersistenceLoadingType returns the PersistenceLoadingType field if non-nil, zero value otherwise.

### GetPersistenceLoadingTypeOk

`func (o *PersistenceLoadingPolicy) GetPersistenceLoadingTypeOk() (*PersistenceLoadingType, bool)`

GetPersistenceLoadingTypeOk returns a tuple with the PersistenceLoadingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPersistenceLoadingType

`func (o *PersistenceLoadingPolicy) SetPersistenceLoadingType(v PersistenceLoadingType)`

SetPersistenceLoadingType sets PersistenceLoadingType field to given value.

### HasPersistenceLoadingType

`func (o *PersistenceLoadingPolicy) HasPersistenceLoadingType() bool`

HasPersistenceLoadingType returns a boolean if a field has been set.

### GetPartialLoadingKeys

`func (o *PersistenceLoadingPolicy) GetPartialLoadingKeys() []string`

GetPartialLoadingKeys returns the PartialLoadingKeys field if non-nil, zero value otherwise.

### GetPartialLoadingKeysOk

`func (o *PersistenceLoadingPolicy) GetPartialLoadingKeysOk() (*[]string, bool)`

GetPartialLoadingKeysOk returns a tuple with the PartialLoadingKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartialLoadingKeys

`func (o *PersistenceLoadingPolicy) SetPartialLoadingKeys(v []string)`

SetPartialLoadingKeys sets PartialLoadingKeys field to given value.

### HasPartialLoadingKeys

`func (o *PersistenceLoadingPolicy) HasPartialLoadingKeys() bool`

HasPartialLoadingKeys returns a boolean if a field has been set.

### GetLockingKeys

`func (o *PersistenceLoadingPolicy) GetLockingKeys() []string`

GetLockingKeys returns the LockingKeys field if non-nil, zero value otherwise.

### GetLockingKeysOk

`func (o *PersistenceLoadingPolicy) GetLockingKeysOk() (*[]string, bool)`

GetLockingKeysOk returns a tuple with the LockingKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLockingKeys

`func (o *PersistenceLoadingPolicy) SetLockingKeys(v []string)`

SetLockingKeys sets LockingKeys field to given value.

### HasLockingKeys

`func (o *PersistenceLoadingPolicy) HasLockingKeys() bool`

HasLockingKeys returns a boolean if a field has been set.

### GetUseKeyAsPrefix

`func (o *PersistenceLoadingPolicy) GetUseKeyAsPrefix() bool`

GetUseKeyAsPrefix returns the UseKeyAsPrefix field if non-nil, zero value otherwise.

### GetUseKeyAsPrefixOk

`func (o *PersistenceLoadingPolicy) GetUseKeyAsPrefixOk() (*bool, bool)`

GetUseKeyAsPrefixOk returns a tuple with the UseKeyAsPrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseKeyAsPrefix

`func (o *PersistenceLoadingPolicy) SetUseKeyAsPrefix(v bool)`

SetUseKeyAsPrefix sets UseKeyAsPrefix field to given value.

### HasUseKeyAsPrefix

`func (o *PersistenceLoadingPolicy) HasUseKeyAsPrefix() bool`

HasUseKeyAsPrefix returns a boolean if a field has been set.

### GetLazyLoadingLargeDataAttributes

`func (o *PersistenceLoadingPolicy) GetLazyLoadingLargeDataAttributes() bool`

GetLazyLoadingLargeDataAttributes returns the LazyLoadingLargeDataAttributes field if non-nil, zero value otherwise.

### GetLazyLoadingLargeDataAttributesOk

`func (o *PersistenceLoadingPolicy) GetLazyLoadingLargeDataAttributesOk() (*bool, bool)`

GetLazyLoadingLargeDataAttributesOk returns a tuple with the LazyLoadingLargeDataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLazyLoadingLargeDataAttributes

`func (o *PersistenceLoadingPolicy) SetLazyLoadingLargeDataAttributes(v bool)`

SetLazyLoadingLargeDataAttributes sets LazyLoadingLargeDataAttributes field to given value.

### HasLazyLoadingLargeDataAttributes

`func (o *PersistenceLoadingPolicy) HasLazyLoadingLargeDataAttributes() bool`

HasLazyLoadingLargeDataAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


