# WorkflowDumpResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Checksum** | **string** |  | 
**TotalPages** | **int32** |  | 
**JsonData** | **string** |  | 

## Methods

### NewWorkflowDumpResponse

`func NewWorkflowDumpResponse(checksum string, totalPages int32, jsonData string, ) *WorkflowDumpResponse`

NewWorkflowDumpResponse instantiates a new WorkflowDumpResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowDumpResponseWithDefaults

`func NewWorkflowDumpResponseWithDefaults() *WorkflowDumpResponse`

NewWorkflowDumpResponseWithDefaults instantiates a new WorkflowDumpResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChecksum

`func (o *WorkflowDumpResponse) GetChecksum() string`

GetChecksum returns the Checksum field if non-nil, zero value otherwise.

### GetChecksumOk

`func (o *WorkflowDumpResponse) GetChecksumOk() (*string, bool)`

GetChecksumOk returns a tuple with the Checksum field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChecksum

`func (o *WorkflowDumpResponse) SetChecksum(v string)`

SetChecksum sets Checksum field to given value.


### GetTotalPages

`func (o *WorkflowDumpResponse) GetTotalPages() int32`

GetTotalPages returns the TotalPages field if non-nil, zero value otherwise.

### GetTotalPagesOk

`func (o *WorkflowDumpResponse) GetTotalPagesOk() (*int32, bool)`

GetTotalPagesOk returns a tuple with the TotalPages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalPages

`func (o *WorkflowDumpResponse) SetTotalPages(v int32)`

SetTotalPages sets TotalPages field to given value.


### GetJsonData

`func (o *WorkflowDumpResponse) GetJsonData() string`

GetJsonData returns the JsonData field if non-nil, zero value otherwise.

### GetJsonDataOk

`func (o *WorkflowDumpResponse) GetJsonDataOk() (*string, bool)`

GetJsonDataOk returns a tuple with the JsonData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJsonData

`func (o *WorkflowDumpResponse) SetJsonData(v string)`

SetJsonData sets JsonData field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


