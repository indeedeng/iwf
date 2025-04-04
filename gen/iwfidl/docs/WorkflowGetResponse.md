# WorkflowGetResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowRunId** | **string** |  | 
**WorkflowStatus** | [**WorkflowStatus**](WorkflowStatus.md) |  | 
**WorkflowType** | **string** |  | 
**Results** | Pointer to [**[]StateCompletionOutput**](StateCompletionOutput.md) |  | [optional] 
**ErrorType** | Pointer to [**WorkflowErrorType**](WorkflowErrorType.md) |  | [optional] 
**ErrorMessage** | Pointer to **string** |  | [optional] 

## Methods

### NewWorkflowGetResponse

`func NewWorkflowGetResponse(workflowRunId string, workflowStatus WorkflowStatus, workflowType string, ) *WorkflowGetResponse`

NewWorkflowGetResponse instantiates a new WorkflowGetResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowGetResponseWithDefaults

`func NewWorkflowGetResponseWithDefaults() *WorkflowGetResponse`

NewWorkflowGetResponseWithDefaults instantiates a new WorkflowGetResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowRunId

`func (o *WorkflowGetResponse) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowGetResponse) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowGetResponse) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.


### GetWorkflowStatus

`func (o *WorkflowGetResponse) GetWorkflowStatus() WorkflowStatus`

GetWorkflowStatus returns the WorkflowStatus field if non-nil, zero value otherwise.

### GetWorkflowStatusOk

`func (o *WorkflowGetResponse) GetWorkflowStatusOk() (*WorkflowStatus, bool)`

GetWorkflowStatusOk returns a tuple with the WorkflowStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowStatus

`func (o *WorkflowGetResponse) SetWorkflowStatus(v WorkflowStatus)`

SetWorkflowStatus sets WorkflowStatus field to given value.


### GetWorkflowType

`func (o *WorkflowGetResponse) GetWorkflowType() string`

GetWorkflowType returns the WorkflowType field if non-nil, zero value otherwise.

### GetWorkflowTypeOk

`func (o *WorkflowGetResponse) GetWorkflowTypeOk() (*string, bool)`

GetWorkflowTypeOk returns a tuple with the WorkflowType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowType

`func (o *WorkflowGetResponse) SetWorkflowType(v string)`

SetWorkflowType sets WorkflowType field to given value.


### GetResults

`func (o *WorkflowGetResponse) GetResults() []StateCompletionOutput`

GetResults returns the Results field if non-nil, zero value otherwise.

### GetResultsOk

`func (o *WorkflowGetResponse) GetResultsOk() (*[]StateCompletionOutput, bool)`

GetResultsOk returns a tuple with the Results field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResults

`func (o *WorkflowGetResponse) SetResults(v []StateCompletionOutput)`

SetResults sets Results field to given value.

### HasResults

`func (o *WorkflowGetResponse) HasResults() bool`

HasResults returns a boolean if a field has been set.

### GetErrorType

`func (o *WorkflowGetResponse) GetErrorType() WorkflowErrorType`

GetErrorType returns the ErrorType field if non-nil, zero value otherwise.

### GetErrorTypeOk

`func (o *WorkflowGetResponse) GetErrorTypeOk() (*WorkflowErrorType, bool)`

GetErrorTypeOk returns a tuple with the ErrorType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorType

`func (o *WorkflowGetResponse) SetErrorType(v WorkflowErrorType)`

SetErrorType sets ErrorType field to given value.

### HasErrorType

`func (o *WorkflowGetResponse) HasErrorType() bool`

HasErrorType returns a boolean if a field has been set.

### GetErrorMessage

`func (o *WorkflowGetResponse) GetErrorMessage() string`

GetErrorMessage returns the ErrorMessage field if non-nil, zero value otherwise.

### GetErrorMessageOk

`func (o *WorkflowGetResponse) GetErrorMessageOk() (*string, bool)`

GetErrorMessageOk returns a tuple with the ErrorMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorMessage

`func (o *WorkflowGetResponse) SetErrorMessage(v string)`

SetErrorMessage sets ErrorMessage field to given value.

### HasErrorMessage

`func (o *WorkflowGetResponse) HasErrorMessage() bool`

HasErrorMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


