# WorkflowDumpRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | **string** |  | 
**PageSizeInBytes** | **int32** |  | 
**PageNum** | **int32** |  | 

## Methods

### NewWorkflowDumpRequest

`func NewWorkflowDumpRequest(workflowId string, workflowRunId string, pageSizeInBytes int32, pageNum int32, ) *WorkflowDumpRequest`

NewWorkflowDumpRequest instantiates a new WorkflowDumpRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowDumpRequestWithDefaults

`func NewWorkflowDumpRequestWithDefaults() *WorkflowDumpRequest`

NewWorkflowDumpRequestWithDefaults instantiates a new WorkflowDumpRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowDumpRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowDumpRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowDumpRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowDumpRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowDumpRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowDumpRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.


### GetPageSizeInBytes

`func (o *WorkflowDumpRequest) GetPageSizeInBytes() int32`

GetPageSizeInBytes returns the PageSizeInBytes field if non-nil, zero value otherwise.

### GetPageSizeInBytesOk

`func (o *WorkflowDumpRequest) GetPageSizeInBytesOk() (*int32, bool)`

GetPageSizeInBytesOk returns a tuple with the PageSizeInBytes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSizeInBytes

`func (o *WorkflowDumpRequest) SetPageSizeInBytes(v int32)`

SetPageSizeInBytes sets PageSizeInBytes field to given value.


### GetPageNum

`func (o *WorkflowDumpRequest) GetPageNum() int32`

GetPageNum returns the PageNum field if non-nil, zero value otherwise.

### GetPageNumOk

`func (o *WorkflowDumpRequest) GetPageNumOk() (*int32, bool)`

GetPageNumOk returns a tuple with the PageNum field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageNum

`func (o *WorkflowDumpRequest) SetPageNum(v int32)`

SetPageNum sets PageNum field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


