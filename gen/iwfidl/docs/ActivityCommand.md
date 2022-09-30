# ActivityCommand

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | **string** |  | 
**ActivityType** | **string** |  | 
**Input** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**ActivityOptions** | Pointer to [**ActivityOptions**](ActivityOptions.md) |  | [optional] 

## Methods

### NewActivityCommand

`func NewActivityCommand(commandId string, activityType string, ) *ActivityCommand`

NewActivityCommand instantiates a new ActivityCommand object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActivityCommandWithDefaults

`func NewActivityCommandWithDefaults() *ActivityCommand`

NewActivityCommandWithDefaults instantiates a new ActivityCommand object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *ActivityCommand) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *ActivityCommand) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *ActivityCommand) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.


### GetActivityType

`func (o *ActivityCommand) GetActivityType() string`

GetActivityType returns the ActivityType field if non-nil, zero value otherwise.

### GetActivityTypeOk

`func (o *ActivityCommand) GetActivityTypeOk() (*string, bool)`

GetActivityTypeOk returns a tuple with the ActivityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActivityType

`func (o *ActivityCommand) SetActivityType(v string)`

SetActivityType sets ActivityType field to given value.


### GetInput

`func (o *ActivityCommand) GetInput() EncodedObject`

GetInput returns the Input field if non-nil, zero value otherwise.

### GetInputOk

`func (o *ActivityCommand) GetInputOk() (*EncodedObject, bool)`

GetInputOk returns a tuple with the Input field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInput

`func (o *ActivityCommand) SetInput(v EncodedObject)`

SetInput sets Input field to given value.

### HasInput

`func (o *ActivityCommand) HasInput() bool`

HasInput returns a boolean if a field has been set.

### GetActivityOptions

`func (o *ActivityCommand) GetActivityOptions() ActivityOptions`

GetActivityOptions returns the ActivityOptions field if non-nil, zero value otherwise.

### GetActivityOptionsOk

`func (o *ActivityCommand) GetActivityOptionsOk() (*ActivityOptions, bool)`

GetActivityOptionsOk returns a tuple with the ActivityOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActivityOptions

`func (o *ActivityCommand) SetActivityOptions(v ActivityOptions)`

SetActivityOptions sets ActivityOptions field to given value.

### HasActivityOptions

`func (o *ActivityCommand) HasActivityOptions() bool`

HasActivityOptions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


