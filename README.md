# iWF project - main & server repo 

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

iWF is a platform providing an all-in-one tooling for building long-running business application. It provides an 
abstraction for persistence(database, elasticSearch) and more! It aims to provide clean, simple and easy to use interface, 
like an iPhone. 

**It will not make you a 10x developer...but you may feel like one!**

A _long-running process_ is called **`Workflow`**. 

It's a simple and powerful WorkflowAsCode general purpose workflow engine.

The server is back by [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal) as an interpreter.

Related projects:
* [OpenAPI definition between SDKs and server](https://github.com/indeedeng/iwf-idl). 
* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk) 
* [iWF Java Samples](https://github.com/indeedeng/iwf-java-samples)
* [iWF Golang SDK](https://github.com/iworkflowio/iwf-golang-sdk)
* [iWF Golang Samples](https://github.com/iworkflowio/iwf-golang-samples) 
* More SDKs? Contribution is welcome. Any languages can be supported as long as implementing the [OpenAPI IDL](https://github.com/indeedeng/iwf-idl/blob/main/iwf.yaml).

# Table of contents

- [Community & Help](#community--help)
- [What is iWF](#what-is-iwf)
  - [Architecture](#architecture)
  - [Basic Concepts](#basic-concepts)
    - [Workflow and WorkflowState definition](#workflow-and-workflowstate-definition)
    - [Workflow execution and WorkflowState execution](#workflow-execution-and-workflowstate-execution)
    - [Commands](#commands)
    - [Persistence](#persistence)
    - [Communication](#communication)
  - [Client APIs](#client-apis)
- [Why iWF](#why-iwf)
  - [If you are familiar with Cadence/Temporal](#if-you-are-familiar-with-cadencetemporal)
  - [If you are not](#if-you-are-not)
- [How to run this server](#how-to-run-this-server)
  - [Using docker image & docker-compose](#using-docker-image--docker-compose)
  - [How to build & run locally](#how-to-build--run-locally)
  - [How to use in production](#how-to-use-in-production)
- [How to migrate from Cadence/Temporal](#how-to-migrate-from-cadencetemporal)
  - [Activity](#activity)
  - [Signal](#signal)
  - [Timer](#timer)
  - [Query](#query)
  - [Search Attribute](#search-attribute)
  - [Versioning and change compatibility](#versioning-and-change-compatibility)
  - [Parallel execution with synchronization](#parallel-execution-with-synchronization)
  - [ContinueAsNew](#continueasnew)
  - [Non-workflow code](#non-workflow-code)
  - [Testing](#testing)
  - [Anything else](#anything-else)
- [Monitoring and Operations](#monitoring-and-operations)
  - [iWF server](#iwf-server)
  - [iWF application](#iwf-application)
- [Development Plan](#development-plan)
- [Some history](#some-history)
- [Contribution](CONTRIBUTING.md)
- [Posts & Articles](#posts--articles)

# Community & Help
* [Slack Channels](https://join.slack.com/t/iwfglobal/shared_invite/zt-1mgi9q6gw-aog6KBuTHFu1KolBerBaLA)
* [Github Discussion](https://github.com/indeedeng/iwf/discussions)
* [StackOverflow](https://stackoverflow.com/questions/tagged/iwf)
* [Github Issues](https://github.com/indeedeng/iwf/issues)

# What is iWF

## Architecture
An iWF application includes a set of iWF workflow workers which host two REST APIs of WorkflowState `start` and `decide`.
The application calls an iWF server to operate on workflow executions -- start, stop, signal, get results, etc, using iWF SDKs.

The iWF server hosts those APIs(also REST) as a iWF API service. Internally the API service will call Cadence/Temporal 
service as the backend.

The iWF server also includes Cadence/Temporal workers which host [an interpreter workflow](https://github.com/indeedeng/iwf/blob/main/service/interpreter/workflowImpl.go).
Internally, any iWF workflows are interpreted into this Cadence/Temporal workflow. The interpreter workflow will invoke 
the two application worker APIs(`start` and `decide`).
The API invocations are implemented by Cadence/Temporal activities. Therefore, all the REST API request/response with the worker are 
stored in history events which are useful for debugging/troubleshooting, and no replay is needed for application workflow code.

![architecture diagram](https://user-images.githubusercontent.com/4523955/207514928-56fea636-c711-4f20-9e90-94ddd1c9844d.png)

* See [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit) for more details. 
 
## Basic Concepts

### Workflow and WorkflowState definition
iWF lets you build long-running applications by implementing the workflow interface, e.g. 
[Java Workflow interface](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/Workflow.java) 
or [Golang Workflow interface](https://github.com/iworkflowio/iwf-golang-sdk/blob/main/iwf/workflow.go).
An instance of the interface is a `WorkflowDefinition`. User applications use `IwfWorkflowType` to differentiate WorkflowDefinitions.    

A WorkflowDefinition contains several `WorkflowState` e.g. 
[Java WorkflowState interface](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/WorkflowState.java) 
or [Golang WorkflowState interface](https://github.com/iworkflowio/iwf-golang-sdk/blob/main/iwf/workflow_state.go). 
A WorkflowState is implemented with two APIs: `start` and `decide`. 
* `start` API is invoked immediately when a WorkflowState is started. It will return some `Commands` to server. When the 
requested `Commands` are completed, `decide` API will be triggered. The number of commands can be zero, one or multiple.
* `decide` API will decide next states to execute. Next states can be zero, one or multiple, and can be re-executed as different `stateExecutions`. 

### Workflow execution and WorkflowState execution
Application can start a workflow instance with a `workflowId` for any workflow definition. A workflow instance is called `WorkflowExecution`. 
iWF server returns `runId` of UUID as the identifier of the WorkflowExecution. The runId is globally unique.  

:warning: Note:
> Depends on the context, the only word `workflow` may mean WorkflowExecution(most commonly), WorkflowDefinition or both.

For a running WorkflowExecution, there must be at least one `WorkflowState` being executed, otherwise the workflow execution will complete. 
An execution instance of WorkflowState is called `StateExecution`, which by identified `StateExecutionId`. A `StateExecutionId` is formatted
as `<StateId>-<Number>`. `StateId` is defined by workflow state definition, while `Number` is how many times this `StateId` has been executed.
StateExecutionId is only unique within the workflow execution.


WorkflowId uniqueness and reuse: For the same workflowId, there must be at most one WorkflowExecution running at anytime. However,
after a previous WorkflowExecution finished running (in any closed status),
application may start a new WorkflowExecution with the same workflowId using appropriate `IdReusePolicy`.


### Commands
These are the three command types:
* `SignalCommand`: will be waiting for a signal from external to the workflow signal channel. External application can use SignalWorkflow API to signal a workflow. 
* `TimerCommand`: will be waiting for a **durable timer** to fire.
* `InterStateChannelCommand`: will be waiting for a value being published from another state in the same workflow execution

Note that `start` API can return multiple commands, and choose different DeciderTriggerType for triggering decide API:
* `AllCommandCompleted`: this will wait for all command completed
* `AnyCommandCompleted`: this will wait for any command completed
* `AnyCommandCombinationCompleted`: this will wait for a list of command combinations on any combination completed

### Persistence
iWF provides super simple persistence abstraction. Developers don't need to touch any database system to register/maintain the schemas. 
The only schema is defined in the workflow code.
* `DataObject` is  
  * sharing some data values across the workflow
  * can be retrieved by external application using GetDataObjects API
  * can be viewed in Cadence/Temporal WebUI in QueryHandler tab
* `SearchAttribute` is similarly:
  * sharing some data values across the workflow
  * can be retrieved by external application using GetSearchAttributes API
  * search for workflows by external application using `SearchWorkflow` API
  * search for workflows in Cadence/Temporal WebUI in Advanced tab
  * search attribute type must be registered in Cadence/Temporal server before using for searching because it is backed up ElasticSearch
  * the data types supported are limited as server has to understand the value for indexing
  * See [Temporal doc](https://docs.temporal.io/concepts/what-is-a-search-attribute) and [Cadence doc](https://cadenceworkflow.io/docs/concepts/search-workflows/) to understand more about SearchAttribute 
* `StateLocal` is for
  * passing some data values from state API to decide API in the same WorkflowState execution
* `RecordEvent` is for
  * recording some events within the state execution. They are useful for debugging using Workflow history. Usually you may want to record the input/output of the dependency RPC calls.

Logically, each workflow type will have a persistence schema like below:
```text
+-------------+-------+-----------------+-----------------+----------------------+----------------------+-----+
| workflowId  | runId | dataObject key1 | dataObject key2 | searchAttribute key1 | searchAttribute key2 | ... |
+-------------+-------+-----------------+-----------------+----------------------+----------------------+-----+
| your-wf-id1 | uuid1 | valu1           | value2          | keyword-value1       | 123(integer)         | ... |
+-------------+-------+-----------------+-----------------+----------------------+----------------------+-----+
| your-wf-id1 | uuid2 | value3          | value4          | keyword-value2       | 456(integer)         | ... |
+-------------+-------+-----------------+-----------------+----------------------+----------------------+-----+
| your-wf-id2 | uuid3 | value5          | value5          | keyword-value3       | 789(integer)         | ... |
+-------------+-------+-----------------+-----------------+----------------------+----------------------+-----+
| ...         | ...   | ...             | ...             | ...                  | ...                  | ... |
+-------------+-------+-----------------+-----------------+----------------------+----------------------+-----+
```
  
### Communication
There are two major communication mechanism in iWF:
* `SignalChannel` is for receiving input from external asynchronously. It's used with `SignalCommand`.
* `InterStateChannel`: for interaction between state executions. It's used with `InterStateChannelCommand`.

## Client APIs
Client APIs are hosted by iWF server for user workflow application to interact with their workflow executions. 
* Start workflow: start a new workflow execution
* Stop workflow: stop a workflow execution
* Signal workflow: send a signal to a workflow execution 
* Search workflow: search for workflows using a query language like SQL with search attributes
* Get workflow: get basic information about a workflow like status and results(if completed or waiting for completed)
* Get workflow data objects: get the dataObjects of a workflow execution
* Get workflow search attributes: get the search attributes of a workflow execution
* Reset workflow: reset a workflow to previous states
* Skip timer: skip a timer of a workflow (usually for testing or operation)

# Why iWF

## If you are familiar with Cadence/Temporal
* See [Slide deck](https://docs.google.com/presentation/d/1CpsroSf6NeVce_XyUhFTkd9bLHN8UHRtM9NavPCMhj8/edit#slide=id.gfe2f455492_0_56) for 
comparison with Cadence/Temporal.

## If you are not
* Check out this [doc](https://docs.google.com/document/d/1zyCKvy4S2l7XBVJzZuS65OIsqV9CRPPYJY3OBbuWrPE) to understand some history

iWF is an application platform that provides you a comprehensive tooling:
* WorkflowAsCode for highly flexibile/customizable business logic
* Parallel execution of multiple threads of business
* Persistence storage for intermediate states stored as "dataObjects"
* Persistence searchable attributes that can be used for flexible searching, even full text searching, backed by ElasticSearch
* Receiving data from external system by Signal
* Durable timer, and cron job scheduling
* Reset workflow to let you recover the workflows from bad states easily 
* Highly testable and easy to maintain
* ...

# How to run this server

##  Using docker image & docker-compose
Checkout this repo, go to the docker-compose folder and run it:
```shell
cd docker-compose && docker-compose up
```
This by default will run Temporal server with it. 
And it will also register a `default` namespace and required search attributes by iWF.
Link to WebUI: http://localhost:8233/namespaces/default/workflows

By default, iWF server is serving port **8801**, server URL is http://localhost:8801/ )

NOTE:

Use `docker pull iworkflowio/iwf-server:latest` to update the latest image.Or update the docker-compose file to specify the version tag.

## How to build & run locally
* Run `make bins` to build the binary `iwf-server`
* Make sure you have registered the system search attributes required by iWF server:
  * Keyword: IwfWorkflowType 
  * Int: IwfGlobalWorkflowVersion
  * Keyword: IwfExecutingStateIds 
  * See [Contribution](./CONTRIBUTING.md) for more detailed commands.
* Then run  `./iwf-server start` to run the service . This defaults to serve workflows APIs with Temporal interpreter implementation. It requires to have local Temporal setup. See Run with local Temporal.
* Alternatively, run `./iwf-server --config config/development_cadence.yaml start` to run with local Cadence. See below instructions for setting up local Cadence. 
 

## How to use in production

You can customize the docker image, or just use the [api](https://github.com/indeedeng/iwf/tree/main/service/api) and [interpreter](https://github.com/indeedeng/iwf/tree/main/service/interpreter) that are exposed as the api service and workflow service.

For more info, contact qlong.seattle@gmail.com

# How to migrate from Cadence/Temporal
Migrating from Cadence/Temporal is simple and easy. It's only possible to migrate new workflow executions. Let your applications to only start new workflows in iWF. For the existing running workflows in Cadence/Temporal, keep the Cadence/Temporal workers until they are finished.

## Activity
Wait, what? **There is no activity at all in iWF?**

Yes, iWF workflows are essentially a REST service and all the activity code in Cadence/Temporal can just move in iWF workflow code -- start or decide API of WorkflowState.

A main reason that many people use Cadence/Temporal activity is to take advantage of the history showing input/output in WebUI. This is handy for debugging/troubleshooting.
iWF provides a `RecordEvent` API to mimic. It can be called with any arbitrary data, and they will be recorded into history just for debugging/troubleshooting.  

## Signal
Depends on different SDKs of Cadence/Temporal, there are different APIs like SignalMethod/SignalChannel/SignalHandler etc.
In iWF, just use SignalCommand as equivalent. 

In some use cases, you may have multiple signals commands and use `AnyCommandCompleted` decider trigger type to wait for any command completed.

## Timer
There are different timer APIs in Cadence/Temporal depends on which SDK:
* workflow.Sleep(duration)
* workflow.Await(duration, condition)
* workflow.NewTimer(duration)
* ...

In iWF, just use TimerCommand as equivalent.

Again in some use cases, you may combine signal/timer commands and use `AnyCommandCompleted` decider trigger type to wait for any command completed.

## Query
Depends on different SDKs of Cadence/Temporal, there are different APIs like QueryHandler/QueryMethod/etc. 

In iWF, use DataObjects as equivalent. Unlike Cadence/Temporal, DataObjects should be explicitly defined in WorkflowDefinition.

Note that by default all the DataObjects and SearchAttributes will be loaded into any WorkflowState as `LOAD_ALL_WITHOUT_LOCKING` persistence loading policy. 
This could be a performance issue if there are too many big items. Consider using different loading policy like `LOAD_PARTIAL_WITHOUT_LOCKING` to improve by customizing the WorkflowStateOptions.

Also note that DataObjects are not just for returning data to API, but also for sharing data across different StateExecutions. But if it's just to share data from start API to decide API in the same StateExecution, using `StateLocal` is preferred for efficiency reason.

## Search Attribute
iWF has the same concepts of Search Attribute.
Unlike Cadence/Temporal, SearchAttribute should be explicitly defined in WorkflowDefinition.

## Versioning and change compatibility
There is no versioning API at all in iWF! 

As there is no replay at all for iWF workflow applications, there is no non-deterministic errors or versioning API. All workflow state executions are stored in Cadence/Temporal activities of the interpreter workflow activities.

Workflow code change will always apply to any running existing and new workflow executions once deployed. This gives superpower and flexibility to maintain long-running business applications.

However, making workflow code change will still have backward-compatibility issue like all other microservice applications. 
Below are the standard ways to address the issues:

1) It's rare but if you don't want old workflows to execute the new code, use a flag in new executions to branch out. For example, if changing flow `StateA->StateB` to `StateA->StateC` only for new workflows, then set a new flag in the new workflow so that StateA can decide go to StateB or StateC. 
2) Removing state code could cause errors(state not found) if there is any state execution still running.  For example, after changed `StateA->StateB` to `StateA->StateC`, you may want to delete StateB. If a StateExecution stays at StateB(most commonly waiting for commands to complete), deleting StateB will cause a not found error when StateB is executed.
   1) The error will be gone if you add the StateB back. Because by default, all State APIs will be backoff retried forever.
   2) If you want to delete StateB as early as possible, use `IwfWorkflowType` and `IwfExecutingStateIds` search attributes to confirm if there is any workflows still running at the state. These are built-in search attributes from iWF server.  

## Parallel execution with synchronization
In Cadence/Temporal, multi-threading is powerful for complicated applications. But the APIs are hard to understand, to use, and to debug. Especially each language/SDK has its own set of APIs without much consistency.

In iWF, there are just a few concepts that are very straightforward:
1) The `decide` API can go to multiple next states. All next states will be executed in parallel
2) `decide` API can also go back to any previous StateId, or the same StateId, to form a loop. The StateExecutionId is the unique identifier. 
3) Use `InterStateChannel` for synchronization communication. It's just like a signal channel that works internally.

Some notes:

1) Any state can decide to complete or fail the workflow, or just go to a dead end(no next state).
2) Because of above, there could be zero, or more than one state completing with data as workflow results. 
3) To get multiple state results from a workflow execution, use the special API `getComplexWorkflowResult` of client API.

## ContinueAsNew

There is NO ContinueAsNew API exposed to user workflow!
ContinueAsNew of Cadence/Temporal is a purely leaked technical details. It's due to the replay model conflicting with the underlying storage limit/performance.
As iWF is built on Cadence/Temporal, it will be implemented in a way that is transparent to user workflows. 

Internally the interpreter workflow can continueAsNew without letting iWF user workflow to know. This is called "auto continueAsNew" --
 
After exceeding the history threshold(defined by numOfStateExecutionCompleted) auto continueAsNew,
AutoContinueAsNew will carry over the pending states, along with all the internal states like DataObjects, interStateChannels, searchAttributes.

_This feature is WIP._ See [this issue](https://github.com/indeedeng/iwf/issues/107) for progress.

For now, the recommended workaround is to do a manual "continueAsNew" by using "TerminateIfRunning" IdReusePolicy to start a new 
workflow execution with the same workflowId. This is almost same as Cadence/Temporal's ContinueAsNew API(an atomic operation), except for some minor difference:
* It won't carry over Search Attributes automatically. You have to carry over them using initial search attributes by WorkflowOptions
* The old execution will be terminated, rather than in "ContinuedAsNew" status

Same as using Cadence/Temporal's ContinueAsNew API, user must ensure all the signals are drained otherwise signals could be lost.

## Non-workflow code
Check [Client APIs](#client-apis) for all the APIs that are equivalent to Cadence/Temporal client APIs.

Features like `IdReusePolicy`, `CronSchedule`, `RetryPolicy` are also supported in iWF.

What's more, there are features that are impossible in Cadence/Temporal are provided like reset workflow by StateId or StateExecutionId. 
Because WorkflowState are explicitly defined, resetting API is a lot more friendly to use. 

## Testing
For unit testing, user code should be mocking all the dependencies in `WorkflowState` implementation, including the input/output of the 
`start` and `decide` API. Users should be able to use any standard testing frameworks/libraries. 

For integration test, iWF provides a SkipTimer API to fire any timer immediately. Although users can always implement this themselves,
the API provides a standard way and saves the effort of re-inventing the wheels.

## Anything else
Is that all? For now yes. We believe these are all you need to migrate to iWF from Cadence/Temporal.

The main philosophy of iWF is providing simple and easy to understand APIs to users(as minimist), as apposed to the complicated and also huge number APIs in Cadence/Temporal. 

So what about something else like:
* Timeout and backoff retry: State start/decide APIs have default timeout and infinite backoff retry. You can customize in StateOptions.  
* ChildWorkflow can be replaced with regular workflow + signal. See this [StackOverflow](https://stackoverflow.com/questions/74494134/should-i-use-child-workflow-or-use-activity-to-start-new-workflow) for why.
* SignalWithStart: Use start + signal API will be the same except for more exception handling work. We have seen a lot of people don't know how to use it correctly in Cadence/Temporal. We will consider provide it in a better way in the future.
* Long-running activity with stateful recovery(heartbeat details): this is indeed a good one that we want to add. But this is not very commonly used yet. Please let us know if you really need it. For now a workaround is use a workflow instead.

If you believe there is something else you really need, open a [ticket](https://github.com/indeedeng/iwf/issues) or join us in the [discussion](https://github.com/indeedeng/iwf/discussions).


# Monitoring and Operations
## iWF server 
There are two components for iWF server: API service and interpreter worker service.

For API service, set up monitors/dashboards:
* API availability
* API latency

The interpreter worker service is just a standard Cadence/Temporal workflow application. Follow the developer guides:
* For [Cadence to set up monitor/dashboards ](https://cadenceworkflow.io/docs/operation-guide/monitor/#cadence-application-monitoring)
* For [Temporal to set up monitor/dashboards](https://github.com/temporalio/dashboards) and [metrics definition](https://docs.temporal.io/references/sdk-metrics)

## iWF application
As you may realize, iWF application is a typical REST microservice. You just need the standard ways to operate it. 

Usually, you need to set up monitors/dashboards:
* API availability 
* API latency

When something goes wrong in your applications, here are the tips:
* Let your worker service return error stacktrace as the response body to iWF server. E.g. like [this example of Spring Boot using ExceptionHandler](https://github.com/indeedeng/iwf-java-samples/blob/2d500093e2aaecf2d728f78366fee776a73efd29/src/main/java/io/iworkflow/controller/IwfWorkerApiController.java#L51). 
* If you return the full stacktrace in response body, the pending activity view will show it to you! Then use Cadence/Temporal WebUI to debug your application.
* All the input/output to your workflow are stored in the activity input/output of history event. The input is in `ActivityTaskScheduledEvent`, output is in `ActivityTaskCompletedEvent` or in pending activity view if having errors.

# Development Plan
### 1.0
- [x] Start workflow API
- [x] Executing `start`/`decide` APIs and completing workflow
- [x] Parallel execution of multiple states 
- [x] Timer command
- [x] Signal command
- [x] SearchAttributeRW
- [x] DataObjectRW
- [x] StateLocal
- [x] Signal workflow API
- [x] Get DataObjects/SearchAttributes API
- [x] Get workflow info API
- [x] Search workflow API
- [x] Stop workflow API
- [x] Reset workflow API 
- [x] Command type(s) for inter-state communications (e.g. internal channel)
- [x] AnyCommandCompleted Decider trigger type
- [x] More workflow start options: IdReusePolicy, cron schedule, retry
- [x] StateOption: Start/Decide API timeout and retry policy
- [x] Reset workflow by stateId or stateExecutionId 
- [x] StateOption.PersistenceLoadingPolicy: LOAD_PARTIAL_WITHOUT_LOCKING

### 1.1
- [x] More Search attribute types: Datetime, double, bool, keyword array, text
- [x] More workflow start options: initial search attributes

### 1.2 
- [x] Skip timer API for testing/operation
- [x] Decider trigger type: any command combination 

### Future
- [ ] Auto continueAsNew([WIP](https://github.com/indeedeng/iwf/issues/107))
- [ ] WaitForMoreResults in StateDecision
- [ ] LongRunningActivityCommand
- [ ] More Decider trigger type
- [ ] Failing workflow details
- [ ] StateOption.PersistenceLoadingPolicy: LOAD_ALL_WITH_EXCLUSIVE_LOCK and LOAD_PARTIAL_WITH_EXCLUSIVE_LOCK

# Some history
AWS published SWF in 2012 and then moved to Step Functions in 2016 because they found it’s too hard to support SWF.
Cadence & Temporal continued the idea of SWF and became much more powerful.
However, AWS is right that the programming of SWF/Cadence/Temporal is hard to adopt because of leaking too many internals.
Inspired by Step Function, iWF is created to provide equivalent power of Cadence/Temporal, but hiding all the internal details
and provide clean and simple API to use. 

Read this [doc](https://docs.google.com/document/d/1zyCKvy4S2l7XBVJzZuS65OIsqV9CRPPYJY3OBbuWrPE) for more.

<img width="916" alt="history diagram" src="https://user-images.githubusercontent.com/4523955/201188875-32e1d070-ab53-4ac5-92fd-bb8ed16dd7dc.png">

# Posts & Articles  
* [A Letter to Cadence/Temporal, and Workflow Tech Community](https://medium.com/@qlong/a-letter-to-cadence-temporal-and-workflow-tech-community-b32e9fa97a0c)
* [iWF vs Cadence/Temporal](https://medium.com/@qlong/iwf-vs-cadence-temporal-1e11b35960fe)
* [iWF vs other general purposed workflow Engines](https://medium.com/@qlong/iwf-vs-other-general-purposed-workflow-engines-f8f3e3d8993d)
* [Cadence® iWF](https://www.instaclustr.com/blog/cadence-iwf/?utm_content=1669999382&utm_medium=linkedin&utm_source=organicsocial)
