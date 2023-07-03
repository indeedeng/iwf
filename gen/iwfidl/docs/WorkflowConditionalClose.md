# WorkflowConditionalClose

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConditionalCloseType** | Pointer to [**WorkflowConditionalCloseType**](WorkflowConditionalCloseType.md) |  | [optional] 
**InternalChannelName** | Pointer to **string** |  | [optional] 
**CloseInput** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**ConditionUnmetDecision** | Pointer to [**[]StateMovement**](StateMovement.md) |  | [optional] 

## Methods

### NewWorkflowConditionalClose

`func NewWorkflowConditionalClose() *WorkflowConditionalClose`

NewWorkflowConditionalClose instantiates a new WorkflowConditionalClose object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowConditionalCloseWithDefaults

`func NewWorkflowConditionalCloseWithDefaults() *WorkflowConditionalClose`

NewWorkflowConditionalCloseWithDefaults instantiates a new WorkflowConditionalClose object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditionalCloseType

`func (o *WorkflowConditionalClose) GetConditionalCloseType() WorkflowConditionalCloseType`

GetConditionalCloseType returns the ConditionalCloseType field if non-nil, zero value otherwise.

### GetConditionalCloseTypeOk

`func (o *WorkflowConditionalClose) GetConditionalCloseTypeOk() (*WorkflowConditionalCloseType, bool)`

GetConditionalCloseTypeOk returns a tuple with the ConditionalCloseType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditionalCloseType

`func (o *WorkflowConditionalClose) SetConditionalCloseType(v WorkflowConditionalCloseType)`

SetConditionalCloseType sets ConditionalCloseType field to given value.

### HasConditionalCloseType

`func (o *WorkflowConditionalClose) HasConditionalCloseType() bool`

HasConditionalCloseType returns a boolean if a field has been set.

### GetInternalChannelName

`func (o *WorkflowConditionalClose) GetInternalChannelName() string`

GetInternalChannelName returns the InternalChannelName field if non-nil, zero value otherwise.

### GetInternalChannelNameOk

`func (o *WorkflowConditionalClose) GetInternalChannelNameOk() (*string, bool)`

GetInternalChannelNameOk returns a tuple with the InternalChannelName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalChannelName

`func (o *WorkflowConditionalClose) SetInternalChannelName(v string)`

SetInternalChannelName sets InternalChannelName field to given value.

### HasInternalChannelName

`func (o *WorkflowConditionalClose) HasInternalChannelName() bool`

HasInternalChannelName returns a boolean if a field has been set.

### GetCloseInput

`func (o *WorkflowConditionalClose) GetCloseInput() EncodedObject`

GetCloseInput returns the CloseInput field if non-nil, zero value otherwise.

### GetCloseInputOk

`func (o *WorkflowConditionalClose) GetCloseInputOk() (*EncodedObject, bool)`

GetCloseInputOk returns a tuple with the CloseInput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloseInput

`func (o *WorkflowConditionalClose) SetCloseInput(v EncodedObject)`

SetCloseInput sets CloseInput field to given value.

### HasCloseInput

`func (o *WorkflowConditionalClose) HasCloseInput() bool`

HasCloseInput returns a boolean if a field has been set.

### GetConditionUnmetDecision

`func (o *WorkflowConditionalClose) GetConditionUnmetDecision() []StateMovement`

GetConditionUnmetDecision returns the ConditionUnmetDecision field if non-nil, zero value otherwise.

### GetConditionUnmetDecisionOk

`func (o *WorkflowConditionalClose) GetConditionUnmetDecisionOk() (*[]StateMovement, bool)`

GetConditionUnmetDecisionOk returns a tuple with the ConditionUnmetDecision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditionUnmetDecision

`func (o *WorkflowConditionalClose) SetConditionUnmetDecision(v []StateMovement)`

SetConditionUnmetDecision sets ConditionUnmetDecision field to given value.

### HasConditionUnmetDecision

`func (o *WorkflowConditionalClose) HasConditionUnmetDecision() bool`

HasConditionUnmetDecision returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


