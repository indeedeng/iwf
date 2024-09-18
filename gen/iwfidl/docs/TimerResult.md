# TimerResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | **string** |  | 
**TimerStatus** | [**TimerStatus**](TimerStatus.md) |  | 

## Methods

### NewTimerResult

`func NewTimerResult(commandId string, timerStatus TimerStatus, ) *TimerResult`

NewTimerResult instantiates a new TimerResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTimerResultWithDefaults

`func NewTimerResultWithDefaults() *TimerResult`

NewTimerResultWithDefaults instantiates a new TimerResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *TimerResult) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *TimerResult) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *TimerResult) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.


### GetTimerStatus

`func (o *TimerResult) GetTimerStatus() TimerStatus`

GetTimerStatus returns the TimerStatus field if non-nil, zero value otherwise.

### GetTimerStatusOk

`func (o *TimerResult) GetTimerStatusOk() (*TimerStatus, bool)`

GetTimerStatusOk returns a tuple with the TimerStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimerStatus

`func (o *TimerResult) SetTimerStatus(v TimerStatus)`

SetTimerStatus sets TimerStatus field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


