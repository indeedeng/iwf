# ActivityResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | **string** |  | 
**ActivityType** | **string** |  | 
**Output** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**ActivityStatus** | **string** |  | 
**TimeoutType** | Pointer to **string** |  | [optional] 

## Methods

### NewActivityResult

`func NewActivityResult(commandId string, activityType string, activityStatus string, ) *ActivityResult`

NewActivityResult instantiates a new ActivityResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActivityResultWithDefaults

`func NewActivityResultWithDefaults() *ActivityResult`

NewActivityResultWithDefaults instantiates a new ActivityResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *ActivityResult) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *ActivityResult) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *ActivityResult) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.


### GetActivityType

`func (o *ActivityResult) GetActivityType() string`

GetActivityType returns the ActivityType field if non-nil, zero value otherwise.

### GetActivityTypeOk

`func (o *ActivityResult) GetActivityTypeOk() (*string, bool)`

GetActivityTypeOk returns a tuple with the ActivityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActivityType

`func (o *ActivityResult) SetActivityType(v string)`

SetActivityType sets ActivityType field to given value.


### GetOutput

`func (o *ActivityResult) GetOutput() EncodedObject`

GetOutput returns the Output field if non-nil, zero value otherwise.

### GetOutputOk

`func (o *ActivityResult) GetOutputOk() (*EncodedObject, bool)`

GetOutputOk returns a tuple with the Output field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutput

`func (o *ActivityResult) SetOutput(v EncodedObject)`

SetOutput sets Output field to given value.

### HasOutput

`func (o *ActivityResult) HasOutput() bool`

HasOutput returns a boolean if a field has been set.

### GetActivityStatus

`func (o *ActivityResult) GetActivityStatus() string`

GetActivityStatus returns the ActivityStatus field if non-nil, zero value otherwise.

### GetActivityStatusOk

`func (o *ActivityResult) GetActivityStatusOk() (*string, bool)`

GetActivityStatusOk returns a tuple with the ActivityStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActivityStatus

`func (o *ActivityResult) SetActivityStatus(v string)`

SetActivityStatus sets ActivityStatus field to given value.


### GetTimeoutType

`func (o *ActivityResult) GetTimeoutType() string`

GetTimeoutType returns the TimeoutType field if non-nil, zero value otherwise.

### GetTimeoutTypeOk

`func (o *ActivityResult) GetTimeoutTypeOk() (*string, bool)`

GetTimeoutTypeOk returns a tuple with the TimeoutType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeoutType

`func (o *ActivityResult) SetTimeoutType(v string)`

SetTimeoutType sets TimeoutType field to given value.

### HasTimeoutType

`func (o *ActivityResult) HasTimeoutType() bool`

HasTimeoutType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


