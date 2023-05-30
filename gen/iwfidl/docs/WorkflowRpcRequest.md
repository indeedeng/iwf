# WorkflowRpcRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**RpcName** | **string** |  | 
**Input** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**SearchAttributesLoadingPolicy** | Pointer to [**PersistenceLoadingPolicy**](PersistenceLoadingPolicy.md) |  | [optional] 
**DataAttributesLoadingPolicy** | Pointer to [**PersistenceLoadingPolicy**](PersistenceLoadingPolicy.md) |  | [optional] 
**TimeoutSeconds** | Pointer to **int32** |  | [optional] 
**UseMemoForDataAttributes** | Pointer to **bool** |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttributeKeyAndType**](SearchAttributeKeyAndType.md) |  | [optional] 

## Methods

### NewWorkflowRpcRequest

`func NewWorkflowRpcRequest(workflowId string, rpcName string, ) *WorkflowRpcRequest`

NewWorkflowRpcRequest instantiates a new WorkflowRpcRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowRpcRequestWithDefaults

`func NewWorkflowRpcRequestWithDefaults() *WorkflowRpcRequest`

NewWorkflowRpcRequestWithDefaults instantiates a new WorkflowRpcRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowRpcRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowRpcRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowRpcRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowRpcRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowRpcRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowRpcRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowRpcRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetRpcName

`func (o *WorkflowRpcRequest) GetRpcName() string`

GetRpcName returns the RpcName field if non-nil, zero value otherwise.

### GetRpcNameOk

`func (o *WorkflowRpcRequest) GetRpcNameOk() (*string, bool)`

GetRpcNameOk returns a tuple with the RpcName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRpcName

`func (o *WorkflowRpcRequest) SetRpcName(v string)`

SetRpcName sets RpcName field to given value.


### GetInput

`func (o *WorkflowRpcRequest) GetInput() EncodedObject`

GetInput returns the Input field if non-nil, zero value otherwise.

### GetInputOk

`func (o *WorkflowRpcRequest) GetInputOk() (*EncodedObject, bool)`

GetInputOk returns a tuple with the Input field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInput

`func (o *WorkflowRpcRequest) SetInput(v EncodedObject)`

SetInput sets Input field to given value.

### HasInput

`func (o *WorkflowRpcRequest) HasInput() bool`

HasInput returns a boolean if a field has been set.

### GetSearchAttributesLoadingPolicy

`func (o *WorkflowRpcRequest) GetSearchAttributesLoadingPolicy() PersistenceLoadingPolicy`

GetSearchAttributesLoadingPolicy returns the SearchAttributesLoadingPolicy field if non-nil, zero value otherwise.

### GetSearchAttributesLoadingPolicyOk

`func (o *WorkflowRpcRequest) GetSearchAttributesLoadingPolicyOk() (*PersistenceLoadingPolicy, bool)`

GetSearchAttributesLoadingPolicyOk returns a tuple with the SearchAttributesLoadingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributesLoadingPolicy

`func (o *WorkflowRpcRequest) SetSearchAttributesLoadingPolicy(v PersistenceLoadingPolicy)`

SetSearchAttributesLoadingPolicy sets SearchAttributesLoadingPolicy field to given value.

### HasSearchAttributesLoadingPolicy

`func (o *WorkflowRpcRequest) HasSearchAttributesLoadingPolicy() bool`

HasSearchAttributesLoadingPolicy returns a boolean if a field has been set.

### GetDataAttributesLoadingPolicy

`func (o *WorkflowRpcRequest) GetDataAttributesLoadingPolicy() PersistenceLoadingPolicy`

GetDataAttributesLoadingPolicy returns the DataAttributesLoadingPolicy field if non-nil, zero value otherwise.

### GetDataAttributesLoadingPolicyOk

`func (o *WorkflowRpcRequest) GetDataAttributesLoadingPolicyOk() (*PersistenceLoadingPolicy, bool)`

GetDataAttributesLoadingPolicyOk returns a tuple with the DataAttributesLoadingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataAttributesLoadingPolicy

`func (o *WorkflowRpcRequest) SetDataAttributesLoadingPolicy(v PersistenceLoadingPolicy)`

SetDataAttributesLoadingPolicy sets DataAttributesLoadingPolicy field to given value.

### HasDataAttributesLoadingPolicy

`func (o *WorkflowRpcRequest) HasDataAttributesLoadingPolicy() bool`

HasDataAttributesLoadingPolicy returns a boolean if a field has been set.

### GetTimeoutSeconds

`func (o *WorkflowRpcRequest) GetTimeoutSeconds() int32`

GetTimeoutSeconds returns the TimeoutSeconds field if non-nil, zero value otherwise.

### GetTimeoutSecondsOk

`func (o *WorkflowRpcRequest) GetTimeoutSecondsOk() (*int32, bool)`

GetTimeoutSecondsOk returns a tuple with the TimeoutSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeoutSeconds

`func (o *WorkflowRpcRequest) SetTimeoutSeconds(v int32)`

SetTimeoutSeconds sets TimeoutSeconds field to given value.

### HasTimeoutSeconds

`func (o *WorkflowRpcRequest) HasTimeoutSeconds() bool`

HasTimeoutSeconds returns a boolean if a field has been set.

### GetUseMemoForDataAttributes

`func (o *WorkflowRpcRequest) GetUseMemoForDataAttributes() bool`

GetUseMemoForDataAttributes returns the UseMemoForDataAttributes field if non-nil, zero value otherwise.

### GetUseMemoForDataAttributesOk

`func (o *WorkflowRpcRequest) GetUseMemoForDataAttributesOk() (*bool, bool)`

GetUseMemoForDataAttributesOk returns a tuple with the UseMemoForDataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseMemoForDataAttributes

`func (o *WorkflowRpcRequest) SetUseMemoForDataAttributes(v bool)`

SetUseMemoForDataAttributes sets UseMemoForDataAttributes field to given value.

### HasUseMemoForDataAttributes

`func (o *WorkflowRpcRequest) HasUseMemoForDataAttributes() bool`

HasUseMemoForDataAttributes returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *WorkflowRpcRequest) GetSearchAttributes() []SearchAttributeKeyAndType`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *WorkflowRpcRequest) GetSearchAttributesOk() (*[]SearchAttributeKeyAndType, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *WorkflowRpcRequest) SetSearchAttributes(v []SearchAttributeKeyAndType)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *WorkflowRpcRequest) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


