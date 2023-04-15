# WorkflowStartRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**IwfWorkflowType** | **string** |  | 
**WorkflowTimeoutSeconds** | **int32** |  | 
**IwfWorkerUrl** | **string** |  | 
**StartStateId** | Pointer to **string** |  | [optional] 
**StateInput** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**StateOptions** | Pointer to [**WorkflowStateOptions**](WorkflowStateOptions.md) |  | [optional] 
**WorkflowStartOptions** | Pointer to [**WorkflowStartOptions**](WorkflowStartOptions.md) |  | [optional] 

## Methods

### NewWorkflowStartRequest

`func NewWorkflowStartRequest(workflowId string, iwfWorkflowType string, workflowTimeoutSeconds int32, iwfWorkerUrl string, ) *WorkflowStartRequest`

NewWorkflowStartRequest instantiates a new WorkflowStartRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStartRequestWithDefaults

`func NewWorkflowStartRequestWithDefaults() *WorkflowStartRequest`

NewWorkflowStartRequestWithDefaults instantiates a new WorkflowStartRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowStartRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowStartRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowStartRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetIwfWorkflowType

`func (o *WorkflowStartRequest) GetIwfWorkflowType() string`

GetIwfWorkflowType returns the IwfWorkflowType field if non-nil, zero value otherwise.

### GetIwfWorkflowTypeOk

`func (o *WorkflowStartRequest) GetIwfWorkflowTypeOk() (*string, bool)`

GetIwfWorkflowTypeOk returns a tuple with the IwfWorkflowType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIwfWorkflowType

`func (o *WorkflowStartRequest) SetIwfWorkflowType(v string)`

SetIwfWorkflowType sets IwfWorkflowType field to given value.


### GetWorkflowTimeoutSeconds

`func (o *WorkflowStartRequest) GetWorkflowTimeoutSeconds() int32`

GetWorkflowTimeoutSeconds returns the WorkflowTimeoutSeconds field if non-nil, zero value otherwise.

### GetWorkflowTimeoutSecondsOk

`func (o *WorkflowStartRequest) GetWorkflowTimeoutSecondsOk() (*int32, bool)`

GetWorkflowTimeoutSecondsOk returns a tuple with the WorkflowTimeoutSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowTimeoutSeconds

`func (o *WorkflowStartRequest) SetWorkflowTimeoutSeconds(v int32)`

SetWorkflowTimeoutSeconds sets WorkflowTimeoutSeconds field to given value.


### GetIwfWorkerUrl

`func (o *WorkflowStartRequest) GetIwfWorkerUrl() string`

GetIwfWorkerUrl returns the IwfWorkerUrl field if non-nil, zero value otherwise.

### GetIwfWorkerUrlOk

`func (o *WorkflowStartRequest) GetIwfWorkerUrlOk() (*string, bool)`

GetIwfWorkerUrlOk returns a tuple with the IwfWorkerUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIwfWorkerUrl

`func (o *WorkflowStartRequest) SetIwfWorkerUrl(v string)`

SetIwfWorkerUrl sets IwfWorkerUrl field to given value.


### GetStartStateId

`func (o *WorkflowStartRequest) GetStartStateId() string`

GetStartStateId returns the StartStateId field if non-nil, zero value otherwise.

### GetStartStateIdOk

`func (o *WorkflowStartRequest) GetStartStateIdOk() (*string, bool)`

GetStartStateIdOk returns a tuple with the StartStateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartStateId

`func (o *WorkflowStartRequest) SetStartStateId(v string)`

SetStartStateId sets StartStateId field to given value.

### HasStartStateId

`func (o *WorkflowStartRequest) HasStartStateId() bool`

HasStartStateId returns a boolean if a field has been set.

### GetStateInput

`func (o *WorkflowStartRequest) GetStateInput() EncodedObject`

GetStateInput returns the StateInput field if non-nil, zero value otherwise.

### GetStateInputOk

`func (o *WorkflowStartRequest) GetStateInputOk() (*EncodedObject, bool)`

GetStateInputOk returns a tuple with the StateInput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateInput

`func (o *WorkflowStartRequest) SetStateInput(v EncodedObject)`

SetStateInput sets StateInput field to given value.

### HasStateInput

`func (o *WorkflowStartRequest) HasStateInput() bool`

HasStateInput returns a boolean if a field has been set.

### GetStateOptions

`func (o *WorkflowStartRequest) GetStateOptions() WorkflowStateOptions`

GetStateOptions returns the StateOptions field if non-nil, zero value otherwise.

### GetStateOptionsOk

`func (o *WorkflowStartRequest) GetStateOptionsOk() (*WorkflowStateOptions, bool)`

GetStateOptionsOk returns a tuple with the StateOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateOptions

`func (o *WorkflowStartRequest) SetStateOptions(v WorkflowStateOptions)`

SetStateOptions sets StateOptions field to given value.

### HasStateOptions

`func (o *WorkflowStartRequest) HasStateOptions() bool`

HasStateOptions returns a boolean if a field has been set.

### GetWorkflowStartOptions

`func (o *WorkflowStartRequest) GetWorkflowStartOptions() WorkflowStartOptions`

GetWorkflowStartOptions returns the WorkflowStartOptions field if non-nil, zero value otherwise.

### GetWorkflowStartOptionsOk

`func (o *WorkflowStartRequest) GetWorkflowStartOptionsOk() (*WorkflowStartOptions, bool)`

GetWorkflowStartOptionsOk returns a tuple with the WorkflowStartOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowStartOptions

`func (o *WorkflowStartRequest) SetWorkflowStartOptions(v WorkflowStartOptions)`

SetWorkflowStartOptions sets WorkflowStartOptions field to given value.

### HasWorkflowStartOptions

`func (o *WorkflowStartRequest) HasWorkflowStartOptions() bool`

HasWorkflowStartOptions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


