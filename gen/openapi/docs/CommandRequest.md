# CommandRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DeciderTriggerType** | Pointer to **string** |  | [optional] 
**ActivityCommands** | Pointer to [**[]ActivityCommand**](ActivityCommand.md) |  | [optional] 

## Methods

### NewCommandRequest

`func NewCommandRequest() *CommandRequest`

NewCommandRequest instantiates a new CommandRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommandRequestWithDefaults

`func NewCommandRequestWithDefaults() *CommandRequest`

NewCommandRequestWithDefaults instantiates a new CommandRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDeciderTriggerType

`func (o *CommandRequest) GetDeciderTriggerType() string`

GetDeciderTriggerType returns the DeciderTriggerType field if non-nil, zero value otherwise.

### GetDeciderTriggerTypeOk

`func (o *CommandRequest) GetDeciderTriggerTypeOk() (*string, bool)`

GetDeciderTriggerTypeOk returns a tuple with the DeciderTriggerType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeciderTriggerType

`func (o *CommandRequest) SetDeciderTriggerType(v string)`

SetDeciderTriggerType sets DeciderTriggerType field to given value.

### HasDeciderTriggerType

`func (o *CommandRequest) HasDeciderTriggerType() bool`

HasDeciderTriggerType returns a boolean if a field has been set.

### GetActivityCommands

`func (o *CommandRequest) GetActivityCommands() []ActivityCommand`

GetActivityCommands returns the ActivityCommands field if non-nil, zero value otherwise.

### GetActivityCommandsOk

`func (o *CommandRequest) GetActivityCommandsOk() (*[]ActivityCommand, bool)`

GetActivityCommandsOk returns a tuple with the ActivityCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActivityCommands

`func (o *CommandRequest) SetActivityCommands(v []ActivityCommand)`

SetActivityCommands sets ActivityCommands field to given value.

### HasActivityCommands

`func (o *CommandRequest) HasActivityCommands() bool`

HasActivityCommands returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


