# WorkflowSignalRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**SignalChannelName** | **string** |  | 
**SignalValue** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewWorkflowSignalRequest

`func NewWorkflowSignalRequest(workflowId string, signalChannelName string, ) *WorkflowSignalRequest`

NewWorkflowSignalRequest instantiates a new WorkflowSignalRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowSignalRequestWithDefaults

`func NewWorkflowSignalRequestWithDefaults() *WorkflowSignalRequest`

NewWorkflowSignalRequestWithDefaults instantiates a new WorkflowSignalRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowSignalRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowSignalRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowSignalRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowSignalRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowSignalRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowSignalRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowSignalRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetSignalChannelName

`func (o *WorkflowSignalRequest) GetSignalChannelName() string`

GetSignalChannelName returns the SignalChannelName field if non-nil, zero value otherwise.

### GetSignalChannelNameOk

`func (o *WorkflowSignalRequest) GetSignalChannelNameOk() (*string, bool)`

GetSignalChannelNameOk returns a tuple with the SignalChannelName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalChannelName

`func (o *WorkflowSignalRequest) SetSignalChannelName(v string)`

SetSignalChannelName sets SignalChannelName field to given value.


### GetSignalValue

`func (o *WorkflowSignalRequest) GetSignalValue() EncodedObject`

GetSignalValue returns the SignalValue field if non-nil, zero value otherwise.

### GetSignalValueOk

`func (o *WorkflowSignalRequest) GetSignalValueOk() (*EncodedObject, bool)`

GetSignalValueOk returns a tuple with the SignalValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalValue

`func (o *WorkflowSignalRequest) SetSignalValue(v EncodedObject)`

SetSignalValue sets SignalValue field to given value.

### HasSignalValue

`func (o *WorkflowSignalRequest) HasSignalValue() bool`

HasSignalValue returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


