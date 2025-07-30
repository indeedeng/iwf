# PublishToInternalChannelRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**Messages** | Pointer to [**[]InterStateChannelPublishing**](InterStateChannelPublishing.md) |  | [optional] 

## Methods

### NewPublishToInternalChannelRequest

`func NewPublishToInternalChannelRequest(workflowId string, ) *PublishToInternalChannelRequest`

NewPublishToInternalChannelRequest instantiates a new PublishToInternalChannelRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPublishToInternalChannelRequestWithDefaults

`func NewPublishToInternalChannelRequestWithDefaults() *PublishToInternalChannelRequest`

NewPublishToInternalChannelRequestWithDefaults instantiates a new PublishToInternalChannelRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *PublishToInternalChannelRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *PublishToInternalChannelRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *PublishToInternalChannelRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetMessages

`func (o *PublishToInternalChannelRequest) GetMessages() []InterStateChannelPublishing`

GetMessages returns the Messages field if non-nil, zero value otherwise.

### GetMessagesOk

`func (o *PublishToInternalChannelRequest) GetMessagesOk() (*[]InterStateChannelPublishing, bool)`

GetMessagesOk returns a tuple with the Messages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessages

`func (o *PublishToInternalChannelRequest) SetMessages(v []InterStateChannelPublishing)`

SetMessages sets Messages field to given value.

### HasMessages

`func (o *PublishToInternalChannelRequest) HasMessages() bool`

HasMessages returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


