# WorkflowWorkerRpcRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Context** | [**Context**](Context.md) |  | 
**WorkflowType** | **string** |  | 
**RpcName** | **string** |  | 
**Input** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 
**DataAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 

## Methods

### NewWorkflowWorkerRpcRequest

`func NewWorkflowWorkerRpcRequest(context Context, workflowType string, rpcName string, ) *WorkflowWorkerRpcRequest`

NewWorkflowWorkerRpcRequest instantiates a new WorkflowWorkerRpcRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowWorkerRpcRequestWithDefaults

`func NewWorkflowWorkerRpcRequestWithDefaults() *WorkflowWorkerRpcRequest`

NewWorkflowWorkerRpcRequestWithDefaults instantiates a new WorkflowWorkerRpcRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContext

`func (o *WorkflowWorkerRpcRequest) GetContext() Context`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *WorkflowWorkerRpcRequest) GetContextOk() (*Context, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *WorkflowWorkerRpcRequest) SetContext(v Context)`

SetContext sets Context field to given value.


### GetWorkflowType

`func (o *WorkflowWorkerRpcRequest) GetWorkflowType() string`

GetWorkflowType returns the WorkflowType field if non-nil, zero value otherwise.

### GetWorkflowTypeOk

`func (o *WorkflowWorkerRpcRequest) GetWorkflowTypeOk() (*string, bool)`

GetWorkflowTypeOk returns a tuple with the WorkflowType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowType

`func (o *WorkflowWorkerRpcRequest) SetWorkflowType(v string)`

SetWorkflowType sets WorkflowType field to given value.


### GetRpcName

`func (o *WorkflowWorkerRpcRequest) GetRpcName() string`

GetRpcName returns the RpcName field if non-nil, zero value otherwise.

### GetRpcNameOk

`func (o *WorkflowWorkerRpcRequest) GetRpcNameOk() (*string, bool)`

GetRpcNameOk returns a tuple with the RpcName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRpcName

`func (o *WorkflowWorkerRpcRequest) SetRpcName(v string)`

SetRpcName sets RpcName field to given value.


### GetInput

`func (o *WorkflowWorkerRpcRequest) GetInput() EncodedObject`

GetInput returns the Input field if non-nil, zero value otherwise.

### GetInputOk

`func (o *WorkflowWorkerRpcRequest) GetInputOk() (*EncodedObject, bool)`

GetInputOk returns a tuple with the Input field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInput

`func (o *WorkflowWorkerRpcRequest) SetInput(v EncodedObject)`

SetInput sets Input field to given value.

### HasInput

`func (o *WorkflowWorkerRpcRequest) HasInput() bool`

HasInput returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *WorkflowWorkerRpcRequest) GetSearchAttributes() []SearchAttribute`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *WorkflowWorkerRpcRequest) GetSearchAttributesOk() (*[]SearchAttribute, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *WorkflowWorkerRpcRequest) SetSearchAttributes(v []SearchAttribute)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *WorkflowWorkerRpcRequest) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.

### GetDataAttributes

`func (o *WorkflowWorkerRpcRequest) GetDataAttributes() []KeyValue`

GetDataAttributes returns the DataAttributes field if non-nil, zero value otherwise.

### GetDataAttributesOk

`func (o *WorkflowWorkerRpcRequest) GetDataAttributesOk() (*[]KeyValue, bool)`

GetDataAttributesOk returns a tuple with the DataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataAttributes

`func (o *WorkflowWorkerRpcRequest) SetDataAttributes(v []KeyValue)`

SetDataAttributes sets DataAttributes field to given value.

### HasDataAttributes

`func (o *WorkflowWorkerRpcRequest) HasDataAttributes() bool`

HasDataAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


