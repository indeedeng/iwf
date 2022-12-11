# WorkflowGetDataObjectsRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**Keys** | Pointer to **[]string** |  | [optional] 

## Methods

### NewWorkflowGetDataObjectsRequest

`func NewWorkflowGetDataObjectsRequest(workflowId string, ) *WorkflowGetDataObjectsRequest`

NewWorkflowGetDataObjectsRequest instantiates a new WorkflowGetDataObjectsRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowGetDataObjectsRequestWithDefaults

`func NewWorkflowGetDataObjectsRequestWithDefaults() *WorkflowGetDataObjectsRequest`

NewWorkflowGetDataObjectsRequestWithDefaults instantiates a new WorkflowGetDataObjectsRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowGetDataObjectsRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowGetDataObjectsRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowGetDataObjectsRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowGetDataObjectsRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowGetDataObjectsRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowGetDataObjectsRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowGetDataObjectsRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetKeys

`func (o *WorkflowGetDataObjectsRequest) GetKeys() []string`

GetKeys returns the Keys field if non-nil, zero value otherwise.

### GetKeysOk

`func (o *WorkflowGetDataObjectsRequest) GetKeysOk() (*[]string, bool)`

GetKeysOk returns a tuple with the Keys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeys

`func (o *WorkflowGetDataObjectsRequest) SetKeys(v []string)`

SetKeys sets Keys field to given value.

### HasKeys

`func (o *WorkflowGetDataObjectsRequest) HasKeys() bool`

HasKeys returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


