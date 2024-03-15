# \DefaultAPI

All URIs are relative to *http://petstore.swagger.io/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1WorkflowConfigUpdatePost**](DefaultAPI.md#ApiV1WorkflowConfigUpdatePost) | **Post** /api/v1/workflow/config/update | update the config of a workflow
[**ApiV1WorkflowDataobjectsGetPost**](DefaultAPI.md#ApiV1WorkflowDataobjectsGetPost) | **Post** /api/v1/workflow/dataobjects/get | get workflow data objects
[**ApiV1WorkflowGetPost**](DefaultAPI.md#ApiV1WorkflowGetPost) | **Post** /api/v1/workflow/get | get a workflow&#39;s status and results(if completed &amp; requested)
[**ApiV1WorkflowGetWithWaitPost**](DefaultAPI.md#ApiV1WorkflowGetWithWaitPost) | **Post** /api/v1/workflow/getWithWait | get a workflow&#39;s status and results(if completed &amp; requested), wait if the workflow is still running
[**ApiV1WorkflowInternalDumpPost**](DefaultAPI.md#ApiV1WorkflowInternalDumpPost) | **Post** /api/v1/workflow/internal/dump | dump internal info of a workflow
[**ApiV1WorkflowResetPost**](DefaultAPI.md#ApiV1WorkflowResetPost) | **Post** /api/v1/workflow/reset | reset a workflow
[**ApiV1WorkflowRpcPost**](DefaultAPI.md#ApiV1WorkflowRpcPost) | **Post** /api/v1/workflow/rpc | execute an RPC of a workflow
[**ApiV1WorkflowSearchPost**](DefaultAPI.md#ApiV1WorkflowSearchPost) | **Post** /api/v1/workflow/search | search for workflows by a search attribute query
[**ApiV1WorkflowSearchattributesGetPost**](DefaultAPI.md#ApiV1WorkflowSearchattributesGetPost) | **Post** /api/v1/workflow/searchattributes/get | get workflow search attributes
[**ApiV1WorkflowSignalPost**](DefaultAPI.md#ApiV1WorkflowSignalPost) | **Post** /api/v1/workflow/signal | signal a workflow
[**ApiV1WorkflowStartPost**](DefaultAPI.md#ApiV1WorkflowStartPost) | **Post** /api/v1/workflow/start | start a workflow
[**ApiV1WorkflowStateDecidePost**](DefaultAPI.md#ApiV1WorkflowStateDecidePost) | **Post** /api/v1/workflowState/decide | for invoking WorkflowState.decide API
[**ApiV1WorkflowStateStartPost**](DefaultAPI.md#ApiV1WorkflowStateStartPost) | **Post** /api/v1/workflowState/start | for invoking WorkflowState.start API
[**ApiV1WorkflowStopPost**](DefaultAPI.md#ApiV1WorkflowStopPost) | **Post** /api/v1/workflow/stop | stop a workflow
[**ApiV1WorkflowTimerSkipPost**](DefaultAPI.md#ApiV1WorkflowTimerSkipPost) | **Post** /api/v1/workflow/timer/skip | skip the timer of a workflow
[**ApiV1WorkflowWaitForStateCompletionPost**](DefaultAPI.md#ApiV1WorkflowWaitForStateCompletionPost) | **Post** /api/v1/workflow/waitForStateCompletion | 
[**ApiV1WorkflowWorkerRpcPost**](DefaultAPI.md#ApiV1WorkflowWorkerRpcPost) | **Post** /api/v1/workflowWorker/rpc | for invoking workflow RPC API in the worker
[**InfoHealthcheckGet**](DefaultAPI.md#InfoHealthcheckGet) | **Get** /info/healthcheck | return health info of the server



## ApiV1WorkflowConfigUpdatePost

> ApiV1WorkflowConfigUpdatePost(ctx).WorkflowConfigUpdateRequest(workflowConfigUpdateRequest).Execute()

update the config of a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowConfigUpdateRequest := *openapiclient.NewWorkflowConfigUpdateRequest("WorkflowId_example", *openapiclient.NewWorkflowConfig()) // WorkflowConfigUpdateRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.ApiV1WorkflowConfigUpdatePost(context.Background()).WorkflowConfigUpdateRequest(workflowConfigUpdateRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowConfigUpdatePost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowConfigUpdatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowConfigUpdateRequest** | [**WorkflowConfigUpdateRequest**](WorkflowConfigUpdateRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowDataobjectsGetPost

> WorkflowGetDataObjectsResponse ApiV1WorkflowDataobjectsGetPost(ctx).WorkflowGetDataObjectsRequest(workflowGetDataObjectsRequest).Execute()

get workflow data objects

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowGetDataObjectsRequest := *openapiclient.NewWorkflowGetDataObjectsRequest("WorkflowId_example") // WorkflowGetDataObjectsRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowDataobjectsGetPost(context.Background()).WorkflowGetDataObjectsRequest(workflowGetDataObjectsRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowDataobjectsGetPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowDataobjectsGetPost`: WorkflowGetDataObjectsResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowDataobjectsGetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowDataobjectsGetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowGetDataObjectsRequest** | [**WorkflowGetDataObjectsRequest**](WorkflowGetDataObjectsRequest.md) |  | 

### Return type

[**WorkflowGetDataObjectsResponse**](WorkflowGetDataObjectsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowGetPost

> WorkflowGetResponse ApiV1WorkflowGetPost(ctx).WorkflowGetRequest(workflowGetRequest).Execute()

get a workflow's status and results(if completed & requested)

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowGetRequest := *openapiclient.NewWorkflowGetRequest("WorkflowId_example") // WorkflowGetRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowGetPost(context.Background()).WorkflowGetRequest(workflowGetRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowGetPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowGetPost`: WorkflowGetResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowGetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowGetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowGetRequest** | [**WorkflowGetRequest**](WorkflowGetRequest.md) |  | 

### Return type

[**WorkflowGetResponse**](WorkflowGetResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowGetWithWaitPost

> WorkflowGetResponse ApiV1WorkflowGetWithWaitPost(ctx).WorkflowGetRequest(workflowGetRequest).Execute()

get a workflow's status and results(if completed & requested), wait if the workflow is still running

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowGetRequest := *openapiclient.NewWorkflowGetRequest("WorkflowId_example") // WorkflowGetRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowGetWithWaitPost(context.Background()).WorkflowGetRequest(workflowGetRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowGetWithWaitPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowGetWithWaitPost`: WorkflowGetResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowGetWithWaitPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowGetWithWaitPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowGetRequest** | [**WorkflowGetRequest**](WorkflowGetRequest.md) |  | 

### Return type

[**WorkflowGetResponse**](WorkflowGetResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowInternalDumpPost

> WorkflowDumpResponse ApiV1WorkflowInternalDumpPost(ctx).WorkflowDumpRequest(workflowDumpRequest).Execute()

dump internal info of a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowDumpRequest := *openapiclient.NewWorkflowDumpRequest("WorkflowId_example", "WorkflowRunId_example", int32(123), int32(123)) // WorkflowDumpRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowInternalDumpPost(context.Background()).WorkflowDumpRequest(workflowDumpRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowInternalDumpPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowInternalDumpPost`: WorkflowDumpResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowInternalDumpPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowInternalDumpPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowDumpRequest** | [**WorkflowDumpRequest**](WorkflowDumpRequest.md) |  | 

### Return type

[**WorkflowDumpResponse**](WorkflowDumpResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowResetPost

> WorkflowResetResponse ApiV1WorkflowResetPost(ctx).WorkflowResetRequest(workflowResetRequest).Execute()

reset a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowResetRequest := *openapiclient.NewWorkflowResetRequest("WorkflowId_example", openapiclient.WorkflowResetType("HISTORY_EVENT_ID")) // WorkflowResetRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowResetPost(context.Background()).WorkflowResetRequest(workflowResetRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowResetPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowResetPost`: WorkflowResetResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowResetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowResetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowResetRequest** | [**WorkflowResetRequest**](WorkflowResetRequest.md) |  | 

### Return type

[**WorkflowResetResponse**](WorkflowResetResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowRpcPost

> WorkflowRpcResponse ApiV1WorkflowRpcPost(ctx).WorkflowRpcRequest(workflowRpcRequest).Execute()

execute an RPC of a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowRpcRequest := *openapiclient.NewWorkflowRpcRequest("WorkflowId_example", "RpcName_example") // WorkflowRpcRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowRpcPost(context.Background()).WorkflowRpcRequest(workflowRpcRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowRpcPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowRpcPost`: WorkflowRpcResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowRpcPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowRpcPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowRpcRequest** | [**WorkflowRpcRequest**](WorkflowRpcRequest.md) |  | 

### Return type

[**WorkflowRpcResponse**](WorkflowRpcResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowSearchPost

> WorkflowSearchResponse ApiV1WorkflowSearchPost(ctx).WorkflowSearchRequest(workflowSearchRequest).Execute()

search for workflows by a search attribute query

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowSearchRequest := *openapiclient.NewWorkflowSearchRequest("Query_example") // WorkflowSearchRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowSearchPost(context.Background()).WorkflowSearchRequest(workflowSearchRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowSearchPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowSearchPost`: WorkflowSearchResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowSearchPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowSearchPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowSearchRequest** | [**WorkflowSearchRequest**](WorkflowSearchRequest.md) |  | 

### Return type

[**WorkflowSearchResponse**](WorkflowSearchResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowSearchattributesGetPost

> WorkflowGetSearchAttributesResponse ApiV1WorkflowSearchattributesGetPost(ctx).WorkflowGetSearchAttributesRequest(workflowGetSearchAttributesRequest).Execute()

get workflow search attributes

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowGetSearchAttributesRequest := *openapiclient.NewWorkflowGetSearchAttributesRequest("WorkflowId_example") // WorkflowGetSearchAttributesRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowSearchattributesGetPost(context.Background()).WorkflowGetSearchAttributesRequest(workflowGetSearchAttributesRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowSearchattributesGetPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowSearchattributesGetPost`: WorkflowGetSearchAttributesResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowSearchattributesGetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowSearchattributesGetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowGetSearchAttributesRequest** | [**WorkflowGetSearchAttributesRequest**](WorkflowGetSearchAttributesRequest.md) |  | 

### Return type

[**WorkflowGetSearchAttributesResponse**](WorkflowGetSearchAttributesResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowSignalPost

> ApiV1WorkflowSignalPost(ctx).WorkflowSignalRequest(workflowSignalRequest).Execute()

signal a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowSignalRequest := *openapiclient.NewWorkflowSignalRequest("WorkflowId_example", "SignalChannelName_example") // WorkflowSignalRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.ApiV1WorkflowSignalPost(context.Background()).WorkflowSignalRequest(workflowSignalRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowSignalPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowSignalPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowSignalRequest** | [**WorkflowSignalRequest**](WorkflowSignalRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowStartPost

> WorkflowStartResponse ApiV1WorkflowStartPost(ctx).WorkflowStartRequest(workflowStartRequest).Execute()

start a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowStartRequest := *openapiclient.NewWorkflowStartRequest("WorkflowId_example", "IwfWorkflowType_example", int32(123), "IwfWorkerUrl_example") // WorkflowStartRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowStartPost(context.Background()).WorkflowStartRequest(workflowStartRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowStartPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowStartPost`: WorkflowStartResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowStartPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowStartPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowStartRequest** | [**WorkflowStartRequest**](WorkflowStartRequest.md) |  | 

### Return type

[**WorkflowStartResponse**](WorkflowStartResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowStateDecidePost

> WorkflowStateDecideResponse ApiV1WorkflowStateDecidePost(ctx).WorkflowStateDecideRequest(workflowStateDecideRequest).Execute()

for invoking WorkflowState.decide API

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowStateDecideRequest := *openapiclient.NewWorkflowStateDecideRequest(*openapiclient.NewContext("WorkflowId_example", "WorkflowRunId_example", int64(123)), "WorkflowType_example", "WorkflowStateId_example") // WorkflowStateDecideRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowStateDecidePost(context.Background()).WorkflowStateDecideRequest(workflowStateDecideRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowStateDecidePost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowStateDecidePost`: WorkflowStateDecideResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowStateDecidePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowStateDecidePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowStateDecideRequest** | [**WorkflowStateDecideRequest**](WorkflowStateDecideRequest.md) |  | 

### Return type

[**WorkflowStateDecideResponse**](WorkflowStateDecideResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowStateStartPost

> WorkflowStateStartResponse ApiV1WorkflowStateStartPost(ctx).WorkflowStateStartRequest(workflowStateStartRequest).Execute()

for invoking WorkflowState.start API

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowStateStartRequest := *openapiclient.NewWorkflowStateStartRequest(*openapiclient.NewContext("WorkflowId_example", "WorkflowRunId_example", int64(123)), "WorkflowType_example", "WorkflowStateId_example") // WorkflowStateStartRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowStateStartPost(context.Background()).WorkflowStateStartRequest(workflowStateStartRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowStateStartPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowStateStartPost`: WorkflowStateStartResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowStateStartPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowStateStartPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowStateStartRequest** | [**WorkflowStateStartRequest**](WorkflowStateStartRequest.md) |  | 

### Return type

[**WorkflowStateStartResponse**](WorkflowStateStartResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowStopPost

> ApiV1WorkflowStopPost(ctx).WorkflowStopRequest(workflowStopRequest).Execute()

stop a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowStopRequest := *openapiclient.NewWorkflowStopRequest("WorkflowId_example") // WorkflowStopRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.ApiV1WorkflowStopPost(context.Background()).WorkflowStopRequest(workflowStopRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowStopPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowStopPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowStopRequest** | [**WorkflowStopRequest**](WorkflowStopRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowTimerSkipPost

> ApiV1WorkflowTimerSkipPost(ctx).WorkflowSkipTimerRequest(workflowSkipTimerRequest).Execute()

skip the timer of a workflow

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowSkipTimerRequest := *openapiclient.NewWorkflowSkipTimerRequest("WorkflowId_example", "WorkflowStateExecutionId_example") // WorkflowSkipTimerRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DefaultAPI.ApiV1WorkflowTimerSkipPost(context.Background()).WorkflowSkipTimerRequest(workflowSkipTimerRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowTimerSkipPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowTimerSkipPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowSkipTimerRequest** | [**WorkflowSkipTimerRequest**](WorkflowSkipTimerRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowWaitForStateCompletionPost

> WorkflowWaitForStateCompletionResponse ApiV1WorkflowWaitForStateCompletionPost(ctx).WorkflowWaitForStateCompletionRequest(workflowWaitForStateCompletionRequest).Execute()



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowWaitForStateCompletionRequest := *openapiclient.NewWorkflowWaitForStateCompletionRequest("WorkflowId_example", "StateExecutionId_example") // WorkflowWaitForStateCompletionRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowWaitForStateCompletionPost(context.Background()).WorkflowWaitForStateCompletionRequest(workflowWaitForStateCompletionRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowWaitForStateCompletionPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowWaitForStateCompletionPost`: WorkflowWaitForStateCompletionResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowWaitForStateCompletionPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowWaitForStateCompletionPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowWaitForStateCompletionRequest** | [**WorkflowWaitForStateCompletionRequest**](WorkflowWaitForStateCompletionRequest.md) |  | 

### Return type

[**WorkflowWaitForStateCompletionResponse**](WorkflowWaitForStateCompletionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1WorkflowWorkerRpcPost

> WorkflowWorkerRpcResponse ApiV1WorkflowWorkerRpcPost(ctx).WorkflowWorkerRpcRequest(workflowWorkerRpcRequest).Execute()

for invoking workflow RPC API in the worker

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {
    workflowWorkerRpcRequest := *openapiclient.NewWorkflowWorkerRpcRequest(*openapiclient.NewContext("WorkflowId_example", "WorkflowRunId_example", int64(123)), "WorkflowType_example", "RpcName_example") // WorkflowWorkerRpcRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.ApiV1WorkflowWorkerRpcPost(context.Background()).WorkflowWorkerRpcRequest(workflowWorkerRpcRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiV1WorkflowWorkerRpcPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiV1WorkflowWorkerRpcPost`: WorkflowWorkerRpcResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiV1WorkflowWorkerRpcPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1WorkflowWorkerRpcPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowWorkerRpcRequest** | [**WorkflowWorkerRpcRequest**](WorkflowWorkerRpcRequest.md) |  | 

### Return type

[**WorkflowWorkerRpcResponse**](WorkflowWorkerRpcResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InfoHealthcheckGet

> HealthInfo InfoHealthcheckGet(ctx).Execute()

return health info of the server

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/indeedeng/iwf-idl"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.InfoHealthcheckGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.InfoHealthcheckGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `InfoHealthcheckGet`: HealthInfo
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.InfoHealthcheckGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInfoHealthcheckGetRequest struct via the builder pattern


### Return type

[**HealthInfo**](HealthInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

