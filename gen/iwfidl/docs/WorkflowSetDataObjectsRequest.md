# WorkflowSetDataObjectsRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**Objects** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 

## Methods

### NewWorkflowSetDataObjectsRequest

`func NewWorkflowSetDataObjectsRequest(workflowId string, ) *WorkflowSetDataObjectsRequest`

NewWorkflowSetDataObjectsRequest instantiates a new WorkflowSetDataObjectsRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowSetDataObjectsRequestWithDefaults

`func NewWorkflowSetDataObjectsRequestWithDefaults() *WorkflowSetDataObjectsRequest`

NewWorkflowSetDataObjectsRequestWithDefaults instantiates a new WorkflowSetDataObjectsRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowSetDataObjectsRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowSetDataObjectsRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowSetDataObjectsRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowSetDataObjectsRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowSetDataObjectsRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowSetDataObjectsRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowSetDataObjectsRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetObjects

`func (o *WorkflowSetDataObjectsRequest) GetObjects() []KeyValue`

GetObjects returns the Objects field if non-nil, zero value otherwise.

### GetObjectsOk

`func (o *WorkflowSetDataObjectsRequest) GetObjectsOk() (*[]KeyValue, bool)`

GetObjectsOk returns a tuple with the Objects field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjects

`func (o *WorkflowSetDataObjectsRequest) SetObjects(v []KeyValue)`

SetObjects sets Objects field to given value.

### HasObjects

`func (o *WorkflowSetDataObjectsRequest) HasObjects() bool`

HasObjects returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


