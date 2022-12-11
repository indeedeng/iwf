# PersistenceLoadingPolicy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PersistenceLoadingType** | Pointer to **string** |  | [optional] 
**PartialLoadingKeys** | Pointer to **[]string** |  | [optional] 

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

`func (o *PersistenceLoadingPolicy) GetPersistenceLoadingType() string`

GetPersistenceLoadingType returns the PersistenceLoadingType field if non-nil, zero value otherwise.

### GetPersistenceLoadingTypeOk

`func (o *PersistenceLoadingPolicy) GetPersistenceLoadingTypeOk() (*string, bool)`

GetPersistenceLoadingTypeOk returns a tuple with the PersistenceLoadingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPersistenceLoadingType

`func (o *PersistenceLoadingPolicy) SetPersistenceLoadingType(v string)`

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


