# WorkflowStateStartRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Context** | Pointer to [**Context**](Context.md) |  | [optional] 
**WorkflowType** | Pointer to **string** |  | [optional] 
**WorkflowStateId** | Pointer to **string** |  | [optional] 
**StateInput** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 
**QueryAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 

## Methods

### NewWorkflowStateStartRequest

`func NewWorkflowStateStartRequest() *WorkflowStateStartRequest`

NewWorkflowStateStartRequest instantiates a new WorkflowStateStartRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStateStartRequestWithDefaults

`func NewWorkflowStateStartRequestWithDefaults() *WorkflowStateStartRequest`

NewWorkflowStateStartRequestWithDefaults instantiates a new WorkflowStateStartRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContext

`func (o *WorkflowStateStartRequest) GetContext() Context`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *WorkflowStateStartRequest) GetContextOk() (*Context, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *WorkflowStateStartRequest) SetContext(v Context)`

SetContext sets Context field to given value.

### HasContext

`func (o *WorkflowStateStartRequest) HasContext() bool`

HasContext returns a boolean if a field has been set.

### GetWorkflowType

`func (o *WorkflowStateStartRequest) GetWorkflowType() string`

GetWorkflowType returns the WorkflowType field if non-nil, zero value otherwise.

### GetWorkflowTypeOk

`func (o *WorkflowStateStartRequest) GetWorkflowTypeOk() (*string, bool)`

GetWorkflowTypeOk returns a tuple with the WorkflowType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowType

`func (o *WorkflowStateStartRequest) SetWorkflowType(v string)`

SetWorkflowType sets WorkflowType field to given value.

### HasWorkflowType

`func (o *WorkflowStateStartRequest) HasWorkflowType() bool`

HasWorkflowType returns a boolean if a field has been set.

### GetWorkflowStateId

`func (o *WorkflowStateStartRequest) GetWorkflowStateId() string`

GetWorkflowStateId returns the WorkflowStateId field if non-nil, zero value otherwise.

### GetWorkflowStateIdOk

`func (o *WorkflowStateStartRequest) GetWorkflowStateIdOk() (*string, bool)`

GetWorkflowStateIdOk returns a tuple with the WorkflowStateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowStateId

`func (o *WorkflowStateStartRequest) SetWorkflowStateId(v string)`

SetWorkflowStateId sets WorkflowStateId field to given value.

### HasWorkflowStateId

`func (o *WorkflowStateStartRequest) HasWorkflowStateId() bool`

HasWorkflowStateId returns a boolean if a field has been set.

### GetStateInput

`func (o *WorkflowStateStartRequest) GetStateInput() EncodedObject`

GetStateInput returns the StateInput field if non-nil, zero value otherwise.

### GetStateInputOk

`func (o *WorkflowStateStartRequest) GetStateInputOk() (*EncodedObject, bool)`

GetStateInputOk returns a tuple with the StateInput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateInput

`func (o *WorkflowStateStartRequest) SetStateInput(v EncodedObject)`

SetStateInput sets StateInput field to given value.

### HasStateInput

`func (o *WorkflowStateStartRequest) HasStateInput() bool`

HasStateInput returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *WorkflowStateStartRequest) GetSearchAttributes() []SearchAttribute`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *WorkflowStateStartRequest) GetSearchAttributesOk() (*[]SearchAttribute, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *WorkflowStateStartRequest) SetSearchAttributes(v []SearchAttribute)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *WorkflowStateStartRequest) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.

### GetQueryAttributes

`func (o *WorkflowStateStartRequest) GetQueryAttributes() []KeyValue`

GetQueryAttributes returns the QueryAttributes field if non-nil, zero value otherwise.

### GetQueryAttributesOk

`func (o *WorkflowStateStartRequest) GetQueryAttributesOk() (*[]KeyValue, bool)`

GetQueryAttributesOk returns a tuple with the QueryAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueryAttributes

`func (o *WorkflowStateStartRequest) SetQueryAttributes(v []KeyValue)`

SetQueryAttributes sets QueryAttributes field to given value.

### HasQueryAttributes

`func (o *WorkflowStateStartRequest) HasQueryAttributes() bool`

HasQueryAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


