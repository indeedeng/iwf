# iWF project - main & server repo 

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

**iWF will make you a 10x developer!**

iWF is a platform providing all-in-one tooling for building long-running business application. It provides an 
abstraction for persistence(database, elasticSearch) and more, with clean, simple and easy to use interface.

It's a simple and powerful WorkflowAsCode general purpose workflow engine. The server is back by [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal) as an interpreter,
preserved the same power of Cadence/Temporal(including scalability/reliability).

Related projects:
* [OpenAPI definition between SDKs and server](https://github.com/indeedeng/iwf-idl). 
* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk) 
* [iWF Java Samples](https://github.com/indeedeng/iwf-java-samples)
  * [Product use case example: subscription](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/subscription)
* [iWF Golang SDK](https://github.com/iworkflowio/iwf-golang-sdk)
* [iWF Golang Samples](https://github.com/iworkflowio/iwf-golang-samples)
  * [Product use case example: subscription](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/subscription)

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
    - [Workflow Diagram](#iwf-workflow-design-diagram)
  - [Client APIs](#client-apis)
- [Why iWF](#why-iwf)
  - [If you are familiar with Cadence/Temporal](#if-you-are-familiar-with-cadencetemporal)
  - [If you are not](#if-you-are-not)
- [How to run this server](#how-to-run-this-server)
  - [Using docker image & docker-compose](#using-docker-image--docker-compose)
  - [How to build & run locally](#how-to-build--run-locally)
  - [How to use in production](#how-to-use-in-production)
- [Monitoring and Operations](#monitoring-and-operations)
  - [iWF server](#iwf-server)
  - [iWF application](#iwf-application)
  - [Debug & Troubleshooting](#troubleshooting)
  - [Operation](#operation)
- [Development Plan](#development-plan)
- [How to migrate from Cadence/Temporal](#how-to-migrate-from-cadencetemporal)
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

A _long-running process_ is called **`Workflow`**.

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

![workflow diagram example](https://user-images.githubusercontent.com/4523955/217110563-ee3f52a0-5a86-440c-af04-30416f29e5db.png)


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

### iWF workflow design diagram

When designing an iWF workflow, it's useful to use iWF state diagrams like this template for visualization.

![state diagram template](https://user-images.githubusercontent.com/4523955/217110210-56631d35-d353-4ecc-8c0c-b826a212a403.png)

For example, the subscription workflow diagram:
* [Java sample](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/subscription)
* [Golang sample](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/subscription)

![subscription state diagram](https://user-images.githubusercontent.com/4523955/217110240-5dfe1d33-0b7c-49f2-8c12-b0d91c4eb970.png)


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
* Check [iWF vs Cadence/Temporal](https://medium.com/@qlong/iwf-vs-cadence-temporal-1e11b35960fe) for 
comparison with Cadence/Temporal.

## If you are not
* Check out this [article](https://medium.com/@qlong/iwf-vs-other-general-purposed-workflow-engines-f8f3e3d8993d) to understand difference between iWF and other workflow engines.

iWF is an application platform that provides you a comprehensive tooling:
* WorkflowAsCode for highly flexible/customizable business logic
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

## Troubleshooting
When something goes wrong in your applications, here are the tips:
* Use query handlers like (`DumpAllInternal` or `GetCurrentTimerInfos`) in Cadence/Temporal WebUI to quickly understand the current status of the workflows.
  * DumpAllInternal will return all the internal status or the pending states
  * GetCurrentTimerInfos will return all the timers of the pending states
* Let your worker service return error stacktrace as the response body to iWF server. E.g. like [this example of Spring Boot using ExceptionHandler](https://github.com/indeedeng/iwf-java-samples/blob/2d500093e2aaecf2d728f78366fee776a73efd29/src/main/java/io/iworkflow/controller/IwfWorkerApiController.java#L51). 
* If you return the full stacktrace in response body, the pending activity view will show it to you! Then use Cadence/Temporal WebUI to debug your application.
* All the input/output to your workflow are stored in the activity input/output of history event. The input is in `ActivityTaskScheduledEvent`, output is in `ActivityTaskCompletedEvent` or in pending activity view if having errors.

## Operation

In additional of using Cadence/Temporal CLI, you can just use [some HTTP script like this](./script/http/local/home.http) to operate on workflows to:
* Start a workflow
* Stop a workflow
* Reset a workflow
* Skip a timer 
* etc

# How to migrate from Cadence/Temporal
Check this [wiki](https://github.com/indeedeng/iwf/wiki/How-to-migrate-from-Cadence-Temporal) for how to migrate from Cadence/Temporal.

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
