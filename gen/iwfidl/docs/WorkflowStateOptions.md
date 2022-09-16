# WorkflowStateOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SearchAttributesLoadingPolicy** | Pointer to [**AttributesLoadingPolicy**](AttributesLoadingPolicy.md) |  | [optional] 
**QueryAttributesLoadingPolicy** | Pointer to [**AttributesLoadingPolicy**](AttributesLoadingPolicy.md) |  | [optional] 
**CommandCarryOverPolicy** | Pointer to [**CommandCarryOverPolicy**](CommandCarryOverPolicy.md) |  | [optional] 

## Methods

### NewWorkflowStateOptions

`func NewWorkflowStateOptions() *WorkflowStateOptions`

NewWorkflowStateOptions instantiates a new WorkflowStateOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStateOptionsWithDefaults

`func NewWorkflowStateOptionsWithDefaults() *WorkflowStateOptions`

NewWorkflowStateOptionsWithDefaults instantiates a new WorkflowStateOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSearchAttributesLoadingPolicy

`func (o *WorkflowStateOptions) GetSearchAttributesLoadingPolicy() AttributesLoadingPolicy`

GetSearchAttributesLoadingPolicy returns the SearchAttributesLoadingPolicy field if non-nil, zero value otherwise.

### GetSearchAttributesLoadingPolicyOk

`func (o *WorkflowStateOptions) GetSearchAttributesLoadingPolicyOk() (*AttributesLoadingPolicy, bool)`

GetSearchAttributesLoadingPolicyOk returns a tuple with the SearchAttributesLoadingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributesLoadingPolicy

`func (o *WorkflowStateOptions) SetSearchAttributesLoadingPolicy(v AttributesLoadingPolicy)`

SetSearchAttributesLoadingPolicy sets SearchAttributesLoadingPolicy field to given value.

### HasSearchAttributesLoadingPolicy

`func (o *WorkflowStateOptions) HasSearchAttributesLoadingPolicy() bool`

HasSearchAttributesLoadingPolicy returns a boolean if a field has been set.

### GetQueryAttributesLoadingPolicy

`func (o *WorkflowStateOptions) GetQueryAttributesLoadingPolicy() AttributesLoadingPolicy`

GetQueryAttributesLoadingPolicy returns the QueryAttributesLoadingPolicy field if non-nil, zero value otherwise.

### GetQueryAttributesLoadingPolicyOk

`func (o *WorkflowStateOptions) GetQueryAttributesLoadingPolicyOk() (*AttributesLoadingPolicy, bool)`

GetQueryAttributesLoadingPolicyOk returns a tuple with the QueryAttributesLoadingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueryAttributesLoadingPolicy

`func (o *WorkflowStateOptions) SetQueryAttributesLoadingPolicy(v AttributesLoadingPolicy)`

SetQueryAttributesLoadingPolicy sets QueryAttributesLoadingPolicy field to given value.

### HasQueryAttributesLoadingPolicy

`func (o *WorkflowStateOptions) HasQueryAttributesLoadingPolicy() bool`

HasQueryAttributesLoadingPolicy returns a boolean if a field has been set.

### GetCommandCarryOverPolicy

`func (o *WorkflowStateOptions) GetCommandCarryOverPolicy() CommandCarryOverPolicy`

GetCommandCarryOverPolicy returns the CommandCarryOverPolicy field if non-nil, zero value otherwise.

### GetCommandCarryOverPolicyOk

`func (o *WorkflowStateOptions) GetCommandCarryOverPolicyOk() (*CommandCarryOverPolicy, bool)`

GetCommandCarryOverPolicyOk returns a tuple with the CommandCarryOverPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandCarryOverPolicy

`func (o *WorkflowStateOptions) SetCommandCarryOverPolicy(v CommandCarryOverPolicy)`

SetCommandCarryOverPolicy sets CommandCarryOverPolicy field to given value.

### HasCommandCarryOverPolicy

`func (o *WorkflowStateOptions) HasCommandCarryOverPolicy() bool`

HasCommandCarryOverPolicy returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


