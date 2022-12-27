# WorkflowSearchRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Query** | **string** |  | 
**PageSize** | Pointer to **int32** |  | [optional] 
**NextPageToken** | Pointer to **string** |  | [optional] 

## Methods

### NewWorkflowSearchRequest

`func NewWorkflowSearchRequest(query string, ) *WorkflowSearchRequest`

NewWorkflowSearchRequest instantiates a new WorkflowSearchRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowSearchRequestWithDefaults

`func NewWorkflowSearchRequestWithDefaults() *WorkflowSearchRequest`

NewWorkflowSearchRequestWithDefaults instantiates a new WorkflowSearchRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetQuery

`func (o *WorkflowSearchRequest) GetQuery() string`

GetQuery returns the Query field if non-nil, zero value otherwise.

### GetQueryOk

`func (o *WorkflowSearchRequest) GetQueryOk() (*string, bool)`

GetQueryOk returns a tuple with the Query field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuery

`func (o *WorkflowSearchRequest) SetQuery(v string)`

SetQuery sets Query field to given value.


### GetPageSize

`func (o *WorkflowSearchRequest) GetPageSize() int32`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *WorkflowSearchRequest) GetPageSizeOk() (*int32, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *WorkflowSearchRequest) SetPageSize(v int32)`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *WorkflowSearchRequest) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.

### GetNextPageToken

`func (o *WorkflowSearchRequest) GetNextPageToken() string`

GetNextPageToken returns the NextPageToken field if non-nil, zero value otherwise.

### GetNextPageTokenOk

`func (o *WorkflowSearchRequest) GetNextPageTokenOk() (*string, bool)`

GetNextPageTokenOk returns a tuple with the NextPageToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPageToken

`func (o *WorkflowSearchRequest) SetNextPageToken(v string)`

SetNextPageToken sets NextPageToken field to given value.

### HasNextPageToken

`func (o *WorkflowSearchRequest) HasNextPageToken() bool`

HasNextPageToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


