# ErrorResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Detail** | Pointer to **string** |  | [optional] 
**SubStatus** | Pointer to [**ErrorSubStatus**](ErrorSubStatus.md) |  | [optional] 
**OriginalWorkerErrorDetail** | Pointer to **string** |  | [optional] 
**OriginalWorkerErrorType** | Pointer to **string** |  | [optional] 
**OriginalWorkerErrorStatus** | Pointer to **int32** |  | [optional] 

## Methods

### NewErrorResponse

`func NewErrorResponse() *ErrorResponse`

NewErrorResponse instantiates a new ErrorResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorResponseWithDefaults

`func NewErrorResponseWithDefaults() *ErrorResponse`

NewErrorResponseWithDefaults instantiates a new ErrorResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDetail

`func (o *ErrorResponse) GetDetail() string`

GetDetail returns the Detail field if non-nil, zero value otherwise.

### GetDetailOk

`func (o *ErrorResponse) GetDetailOk() (*string, bool)`

GetDetailOk returns a tuple with the Detail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetail

`func (o *ErrorResponse) SetDetail(v string)`

SetDetail sets Detail field to given value.

### HasDetail

`func (o *ErrorResponse) HasDetail() bool`

HasDetail returns a boolean if a field has been set.

### GetSubStatus

`func (o *ErrorResponse) GetSubStatus() ErrorSubStatus`

GetSubStatus returns the SubStatus field if non-nil, zero value otherwise.

### GetSubStatusOk

`func (o *ErrorResponse) GetSubStatusOk() (*ErrorSubStatus, bool)`

GetSubStatusOk returns a tuple with the SubStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubStatus

`func (o *ErrorResponse) SetSubStatus(v ErrorSubStatus)`

SetSubStatus sets SubStatus field to given value.

### HasSubStatus

`func (o *ErrorResponse) HasSubStatus() bool`

HasSubStatus returns a boolean if a field has been set.

### GetOriginalWorkerErrorDetail

`func (o *ErrorResponse) GetOriginalWorkerErrorDetail() string`

GetOriginalWorkerErrorDetail returns the OriginalWorkerErrorDetail field if non-nil, zero value otherwise.

### GetOriginalWorkerErrorDetailOk

`func (o *ErrorResponse) GetOriginalWorkerErrorDetailOk() (*string, bool)`

GetOriginalWorkerErrorDetailOk returns a tuple with the OriginalWorkerErrorDetail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginalWorkerErrorDetail

`func (o *ErrorResponse) SetOriginalWorkerErrorDetail(v string)`

SetOriginalWorkerErrorDetail sets OriginalWorkerErrorDetail field to given value.

### HasOriginalWorkerErrorDetail

`func (o *ErrorResponse) HasOriginalWorkerErrorDetail() bool`

HasOriginalWorkerErrorDetail returns a boolean if a field has been set.

### GetOriginalWorkerErrorType

`func (o *ErrorResponse) GetOriginalWorkerErrorType() string`

GetOriginalWorkerErrorType returns the OriginalWorkerErrorType field if non-nil, zero value otherwise.

### GetOriginalWorkerErrorTypeOk

`func (o *ErrorResponse) GetOriginalWorkerErrorTypeOk() (*string, bool)`

GetOriginalWorkerErrorTypeOk returns a tuple with the OriginalWorkerErrorType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginalWorkerErrorType

`func (o *ErrorResponse) SetOriginalWorkerErrorType(v string)`

SetOriginalWorkerErrorType sets OriginalWorkerErrorType field to given value.

### HasOriginalWorkerErrorType

`func (o *ErrorResponse) HasOriginalWorkerErrorType() bool`

HasOriginalWorkerErrorType returns a boolean if a field has been set.

### GetOriginalWorkerErrorStatus

`func (o *ErrorResponse) GetOriginalWorkerErrorStatus() int32`

GetOriginalWorkerErrorStatus returns the OriginalWorkerErrorStatus field if non-nil, zero value otherwise.

### GetOriginalWorkerErrorStatusOk

`func (o *ErrorResponse) GetOriginalWorkerErrorStatusOk() (*int32, bool)`

GetOriginalWorkerErrorStatusOk returns a tuple with the OriginalWorkerErrorStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginalWorkerErrorStatus

`func (o *ErrorResponse) SetOriginalWorkerErrorStatus(v int32)`

SetOriginalWorkerErrorStatus sets OriginalWorkerErrorStatus field to given value.

### HasOriginalWorkerErrorStatus

`func (o *ErrorResponse) HasOriginalWorkerErrorStatus() bool`

HasOriginalWorkerErrorStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


