# iWF project - main & server repo

[![Go Reference](https://pkg.go.dev/badge/github.com/indeedeng/iwf.svg)](https://pkg.go.dev/github.com/indeedeng/iwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/indeedeng/iwf)](https://goreportcard.com/report/github.com/indeedeng/iwf)
[![Coverage Status](https://codecov.io/github/indeedeng/iwf/coverage.svg?branch=main)](https://app.codecov.io/gh/indeedeng/iwf/branch/main)

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

**iWF will make you a 10x developer!**

iWF is an all-in-one platform for developing long-running business processes. It offers a convenient abstraction layer
for utilizing databases, ElasticSearch, message queues, durable timers, and more, with a clean, simple, and
user-friendly interface.

iWF is a versatile WorkflowAsCode engine that is both simple and powerful. By
utilizing [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal)
as an interpreter on the backend, iWF preserves all the capabilities of Cadence/Temporal, while maintaining the same
level of scalability and reliability.

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
    - [If you are familiar with Cadence/Temporal/AWS SWF/Azure Durable Functions](#if-you-are-familiar-with-cadencetemporalaws-swfazure-durable-functions)
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
- [Posts & Articles](#posts--articles--reference)

# Community & Help

* [Slack Channels](https://join.slack.com/t/iwfglobal/shared_invite/zt-1mgi9q6gw-aog6KBuTHFu1KolBerBaLA)
* [Github Discussion](https://github.com/indeedeng/iwf/discussions)
* [StackOverflow](https://stackoverflow.com/questions/tagged/iwf)
* [Github Issues](https://github.com/indeedeng/iwf/issues)

# What is iWF

## Architecture

An iWF application is composed of several iWF workflow workers. These workers host two REST APIs for
WorkflowState `start` and `decide`.
The application utilizes the iWF SDKs to communicate with an iWF server and perform actions on workflow executions, such
as starting, stopping,
signaling, and retrieving results

The iWF server provides the APIs, which are also RESTful, as a iWF API service. Internally, this API service
communicates
with the Cadence/Temporal service as its backend.

In addition to hosting the iWF API service, the iWF server includes Cadence/Temporal workers that
host [an interpreter workflow](https://github.com/indeedeng/iwf/blob/main/service/interpreter/workflowImpl.go).
This interpreter workflow interprets any iWF workflows into the Cadence/Temporal workflow. It then invokes the two
application
worker APIs (WorkflowState `start` and `decide`) through Cadence/Temporal activities. As a result, all REST API requests
and responses
are recorded as history events, which can be useful for debugging and troubleshooting purposes.
This means that there's no need to replay the application workflow code.

![architecture diagram](https://user-images.githubusercontent.com/4523955/207514928-56fea636-c711-4f20-9e90-94ddd1c9844d.png)

* See [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit) for more
  details.

## Basic Concepts

### Workflow and WorkflowState definition

A _long-running process_ is called **`Workflow`**.

iWF enables the building of long-running applications by implementing the Workflow interface in either
[Golang](https://github.com/iworkflowio/iwf-golang-sdk/blob/main/iwf/workflow.go) or
[Java](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/Workflow.java).
An implementation of the interface is referred to as a `WorkflowDefinition`.

A WorkflowDefinition consists of multiple WorkflowStates, which can be implemented using either
the [Java WorkflowState interface](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/WorkflowState.java)
or [Golang WorkflowState interface](https://github.com/iworkflowio/iwf-golang-sdk/blob/main/iwf/workflow_state.go).
A WorkflowState is implemented using two APIs: the `start` API and the `decide` API:

* The `start` API is invoked as soon as a WorkflowState is started. It returns a set of Commands to the server,
  and once these commands are completed, the `decide` API is triggered. The number of commands can range from zero to
  multiple.
* The `decide` API determines the next set of states to be executed. The next states can range from zero to multiple.

The same WorkflowState can be re-executed as different stateExecutions.

![workflow diagram example](https://user-images.githubusercontent.com/4523955/218195868-17818b58-0d00-4523-8cc6-df4c04526c0d.png)

### Workflow execution and WorkflowState execution

An application can initiate a workflow instance for any WorkflowDefinition using a `workflowId`. The instance of a
workflow is referred
to as a `WorkflowExecution`. The iWF server returns a `runId`, which is a universally unique identifier (UUID), as the
identifier for the WorkflowExecution. The runId is guaranteed to be globally unique.

:warning: Note:
> Depending on the context, the term "workflow" may refer to a WorkflowExecution (most commonly), a WorkflowDefinition,
> or both.

For a running WorkflowExecution, there must be at least one WorkflowState being executed, and if there are none, the
WorkflowExecution
will be marked as completed. An instance of a WorkflowState's execution is referred to as a `StateExecution` and is
identified by a
`StateExecutionId`. The StateExecutionId is formatted as `<StateId>-<Number>`, where the `StateId` is defined by the
WorkflowState definition
and the `Number` represents the number of times the StateId has been started. The StateExecutionId is unique only within
the context
of a specific WorkflowExecution.

### Commands

The following are the three types of commands:

* `SignalCommand`: will wait for a signal to be published to the workflow signal channel. External application can use
  SignalWorkflow API to signal a workflow.
* `TimerCommand`: will wait for a **durable timer** to fire.
* `InterStateChannelCommand`: will wait for a value to be published from another state in the same workflow execution

The start API can return multiple commands and choose a different DeciderTriggerType to trigger the decide API.
The available options for the DeciderTriggerType are:

* `AllCommandCompleted`: This option waits for all commands to be completed.

* `AnyCommandCompleted`: This option waits for any of the commands to be completed.

* `AnyCommandCombinationCompleted`: This option waits for any combination of the commands in a specified list to be
  completed.

### Persistence

iWF offers a highly simplified persistence abstraction, which eliminates the need for developers to interact with
any database systems to register or maintain schemas. The only schema that needs to be defined is in the workflow code.

* `DataObject`
    * are used for sharing data values across the workflow.
    * can be retrieved by external application using GetDataObjects API
    * can be viewed in Cadence/Temporal WebUI in QueryHandler tab
* `SearchAttribute`:
    * are used for sharing data values across the workflow similarly to DataObjects
    * can be retrieved by external application using GetSearchAttributes API
    * are used for searching for workflows by external application using `SearchWorkflow` API
    * are used for searching for workflows in Cadence/Temporal WebUI in Advanced tab
    * any search attribute type must be registered in Cadence/Temporal server before using for searching because it is
      backed up ElasticSearch
    * See [Temporal doc](https://docs.temporal.io/concepts/what-is-a-search-attribute)
      and [Cadence doc](https://cadenceworkflow.io/docs/concepts/search-workflows/) to understand more about
      SearchAttribute
* `StateLocals`
    * are used for passing data values from the `start` API to the `decide` API within the same StateExecution,
      thereby reducing the need to use DataObjects
* `RecordEvents`
    * are used to record events within the state execution and are useful for debugging using the Workflow history. They
      can be used to record the input and output of dependency RPC calls, for example

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

iWF has two primary communication mechanisms:

* `SignalChannel`: is used for receiving input from external sources asynchronously. It is employed with the
  SignalCommand.
* `InterStateChannel`: is used for communication between state executions. It is employed with the
  InterStateChannelCommand.

### iWF workflow design diagram

When creating an iWF workflow, utilizing state diagrams similar to this template can be beneficial for visualizing the
process.

![state diagram template](https://user-images.githubusercontent.com/4523955/218195877-9c99f3ff-bba9-41db-86c6-e7352ed8b0f1.png)

For example, the subscription workflow diagram:

* [Java sample](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/subscription)
* [Golang sample](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/subscription)

![subscription state diagram](https://user-images.githubusercontent.com/4523955/218195883-6d8c23ea-130a-481b-bb80-3e5bb3354176.png)

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

## Advanced Concepts

### WorkflowOptions

iWF let you deeply customize the workflow behaviors with the below options.

#### IdReusePolicy for WorkflowId

At any given time, there can be only one WorkflowExecution running for a specific workflowId.
A new WorkflowExecution can be initiated using the same workflowId by setting the appropriate `IdReusePolicy` in
WorkflowOptions.

* `ALLOW_IF_NO_RUNNING` or `ALLOW_DUPLICATE`
    * Allow starting workflow if there is no execution running with the workflowId
    * This is the **default policy** if not specified in WorkflowOptions
* `ALLOW_IF_PREVIOUS_EXISTS_ABNORMALLY` or `ALLOW_DUPLICATE_FAILED_ONLY`
    * Allow starting workflow if a previous Workflow Execution with the same Workflow Id does not have a Completed
      status.
      Use this policy when there is a need to re-execute a Failed, Timed Out, Terminated or Cancelled workflow
      execution.
* `DISALLOW_REUSE` or `REJECT_DUPLICATE`
    * Not allow to start a new workflow execution with the same workflowId.
* `TERMINATE_IF_RUNNING`
    * Always allow starting workflow no matter what -- iWF server will terminate the current running one if it exists.

NOTE: the names `ALLOW_DUPLICATE`/`ALLOW_DUPLICATE_FAILED_ONLY`/`REJECT_DUPLICATE` are inherited from Cadence/Temporal
but iWF provides more accurate names as alternatives.

#### CRON Schedule

iWF allows you to start a workflow with a fixed cron schedule like below

```text
// CronSchedule - Optional cron schedule for workflow. If a cron schedule is specified, the workflow will run
// as a cron based on the schedule. The scheduling will be based on UTC time. The schedule for the next run only happens
// after the current run is completed/failed/timeout. If a RetryPolicy is also supplied, and the workflow failed
// or timed out, the workflow will be retried based on the retry policy. While the workflow is retrying, it won't
// schedule its next run. If the next schedule is due while the workflow is running (or retrying), then it will skip
that
// schedule. Cron workflow will not stop until it is terminated or cancelled (by returning cadence.CanceledError).
// The cron spec is as follows:
// ┌───────────── minute (0 - 59)
// │ ┌───────────── hour (0 - 23)
// │ │ ┌───────────── day of the month (1 - 31)
// │ │ │ ┌───────────── month (1 - 12)
// │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
// │ │ │ │ │
// │ │ │ │ │
// * * * * *
```

NOTE:

* iWF also
  supports [more advanced cron expressions](https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format)
* The [crontab guru](https://crontab.guru/) site is useful for testing your cron expressions.
* To cancel a cron schedule, use terminate of cancel type to stop the workflow execution.
* By default, there is no cron schedule.

#### RetryPolicy for workflow

Workflow execution can have a backoff retry policy which will retry on failed or timeout.

By default, there is no retry policy.

#### Initial Search Attributes

Client can specify some initial search attributes when starting the workflow.

By default, there is no initial search attributes.

### WorkflowStateOptions

Similarly, users can customize the WorkflowState

#### Start/Decide API timeout and retry policy

Users can customize the API timeout and retry policy for WorkflowState Start and Decide API.

By default, the API timeout is 30s with infinite backoff retry:

- InitialIntervalSeconds: 1
- MaxInternalSeconds:100
- MaximumAttempts: 0 # zero means infinite attempts
- BackoffCoefficient: 2

#### Persistence loading policy

When a state API loads DataObjects/SearchAttributes, by default it will load everything which could cause size limit
error
for Cadence/Temporal activity input/output limit(2MB by default). User can use other loading
policy `LOAD_PARTIAL_WITHOUT_LOCKING`
to specify certain DataObjects/SearchAttributes only to load for this WorkflowState.

`WITHOUT_LOCKING` here means if multiple StateExecutions try to upsert the same DataObject/SearchAttribute, it can be
done in parallel without locking.
iWF will provide more advanced policy to allow loading with "locking" in the future.

#### Start API failure policy

By default, the workflow execution will fail when Start/Decide API max out the retry attempts. In some cases that
workflow want to ignore the errors.
A new future is [WIP](https://github.com/indeedeng/iwf/issues/148) to introduce a `StartApiFailurePolicy` to allow this.

Alternatively, WorkflowState can utilize `attempts` or `firstAttemptTime` from the context to decide ignore the
exception/error.

# Why iWF

## If you are familiar with Cadence/Temporal/AWS SWF/Azure Durable Functions

Check [iWF vs Cadence/Temporal](https://medium.com/@qlong/iwf-vs-cadence-temporal-1e11b35960fe) for
comparison with Cadence/Temporal.

The article should still apply to AWS SWF and Azure Durable Functions:

* AWS SWF is the predecessor of Cadence/Temporal and shares the same API, but its capabilities and features are more
  limited when compared.
    * For a comparison between SWF and Cadence, refer to [this post](https://news.ycombinator.com/item?id=19733880.)
* Azure Durable Functions shared the same programming model(replay based workflow execution engine) but also with
  limited features compared to Cadence/Temporal.
    * Additionally, it is recommended to read
      this [article](https://medium.com/@cgillum/common-pitfalls-with-durable-execution-frameworks-like-durable-functions-or-temporal-eaf635d4a8bb)
      about the pitfall about the programming model

## If you are not

Check out this [article](https://medium.com/@qlong/iwf-vs-other-general-purposed-workflow-engines-f8f3e3d8993d) to
understand difference between iWF and other workflow engines.

**TL;DR:**

* WorkflowAsCode for highly flexible/customizable business logic, highly testable and easy to maintain
* Parallel execution of multiple threads of business
* Persistence storage for intermediate states stored as "dataObjects"
* Persistence searchable attributes that can be used for flexible searching, even full text searching, backed by
  ElasticSearch
* Receiving data from external system by Signal
* Durable timer, and cron job scheduling
* Reset workflow to let you recover the workflows from bad states easily
* Troubleshooting/debugging is easy
* Scalability/reliability
* ...

# How to run this server

## Using docker image & docker-compose

Checkout this repo, go to the docker-compose folder and run it:

```shell
cd docker-compose && docker-compose up
```

This by default will run Temporal server with it.
And it will also register a `default` Temporal namespace and required search attributes by iWF.
Link to the Temporal WebUI: http://localhost:8233/namespaces/default/workflows

By default, iWF server is serving port **8801**, server URL is http://localhost:8801/ )

NOTE:

Use `docker pull iworkflowio/iwf-server:latest` to update the latest image.Or update the docker-compose file to specify
the version tag.

## How to build & run locally

* Run `make bins` to build the binary `iwf-server`
* Make sure you have registered the system search attributes required by iWF server:
    * Keyword: IwfWorkflowType
    * Int: IwfGlobalWorkflowVersion
    * Keyword: IwfExecutingStateIds
    * See [Contribution](./CONTRIBUTING.md) for more detailed commands.
    * For Cadence without advancedVisibility enabled,
      set [disableSystemSearchAttributes](https://github.com/indeedeng/iwf/blob/main/config/development_cadence.yaml#L8)
      to true
* Then run  `./iwf-server start` to run the service . This defaults to serve workflows APIs with Temporal interpreter
  implementation. It requires to have local Temporal setup. See Run with local Temporal.
* Alternatively, run `./iwf-server --config config/development_cadence.yaml start` to run with local Cadence. See below
  instructions for setting up local Cadence.

## How to use in production

You can customize the docker image, or just use the [api](https://github.com/indeedeng/iwf/tree/main/service/api)
and [interpreter](https://github.com/indeedeng/iwf/tree/main/service/interpreter) that are exposed as the api service
and workflow service.

* Also make sure you have registered the system search attributes required by iWF server:
    * Keyword: IwfWorkflowType
    * Int: IwfGlobalWorkflowVersion
    * Keyword: IwfExecutingStateIds
    * See [Contribution](./CONTRIBUTING.md) for more detailed commands.
    * For Cadence without advancedVisibility enabled,
      set [disableSystemSearchAttributes](https://github.com/indeedeng/iwf/blob/main/config/development_cadence.yaml#L8)
      to true

For more info, contact qlong.seattle@gmail.com

# Monitoring and Operations

## iWF server

There are two components for iWF server: API service and interpreter worker service.

For API service, set up monitors/dashboards:

* API availability
* API latency

The interpreter worker service is just a standard Cadence/Temporal workflow application. Follow the developer guides.

* For
  [Cadence to set up monitor/dashboards ](https://cadenceworkflow.io/docs/operation-guide/monitor/#cadence-application-monitoring)

* For [Temporal to set up monitor/dashboards](https://github.com/temporalio/dashboards)
  and [metrics definition](https://docs.temporal.io/references/sdk-metrics)

## iWF application

As you may realize, iWF application is a typical REST microservice. You just need the standard ways to operate it.

Usually, you need to set up monitors/dashboards:

* API availability
* API latency

## Troubleshooting

When something goes wrong in your applications, here are the tips:

* Use query handlers like (`DumpAllInternal` or `GetCurrentTimerInfos`) in Cadence/Temporal WebUI to quickly understand
  the current status of the workflows.
    * DumpAllInternal will return all the internal status or the pending states
    * GetCurrentTimerInfos will return all the timers of the pending states
* Let your worker service return error stacktrace as the response body to iWF server. E.g.
  like [this example of Spring Boot using ExceptionHandler](https://github.com/indeedeng/iwf-java-samples/blob/2d500093e2aaecf2d728f78366fee776a73efd29/src/main/java/io/iworkflow/controller/IwfWorkerApiController.java#L51)
  .
* If you return the full stacktrace in response body, the pending activity view will show it to you! Then use
  Cadence/Temporal WebUI to debug your application.
* All the input/output to your workflow are stored in the activity input/output of history event. The input is
  in `ActivityTaskScheduledEvent`, output is in `ActivityTaskCompletedEvent` or in pending activity view if having
  errors.

## Operation

In additional of using Cadence/Temporal CLI, you can just
use [some HTTP script like this](./script/http/local/home.http) to operate on workflows to:

* Start a workflow
* Stop a workflow
* Reset a workflow
* Skip a timer
* etc, any APIs supported by the [iWF server API schema](https://github.com/indeedeng/iwf-idl/blob/main/iwf.yaml)

# How to migrate from Cadence/Temporal

Check this [wiki](https://github.com/indeedeng/iwf/wiki/How-to-migrate-from-Cadence-Temporal) for how to migrate from
Cadence/Temporal.

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

### 1.3

- [x] Support failing workflow with results
- [x] Support differentiate different uncompleted workflow closed status for GetWorkflow

### Future

- [ ] Auto continueAsNew([WIP](https://github.com/indeedeng/iwf/issues/107))
- [ ] WaitForMoreResults in StateDecision
- [ ] LongRunningActivityCommand
- [ ] More Decider trigger type
- [ ] Failing workflow details
- [ ] StateOption.PersistenceLoadingPolicy: LOAD_ALL_WITH_EXCLUSIVE_LOCK and LOAD_PARTIAL_WITH_EXCLUSIVE_LOCK

# Some history

AWS introduced SWF in 2012, but later switched to Step Functions in 2016 because they found it difficult to support.
Cadence and Temporal picked up where SWF left off and extend with more features and more robust, but keeping the same
programming models. Programming with SWF/Cadence/Temporal is challenging because it exposes too many internal details.
iWF was created to offer the same level of power as Cadence and Temporal, but with a clean and simple API that hides
all the underlying complexity.

For more information, please see
the [document](https://docs.google.com/document/d/1zyCKvy4S2l7XBVJzZuS65OIsqV9CRPPYJY3OBbuWrPE).

<img width="916" alt="history diagram" src="https://user-images.githubusercontent.com/4523955/201188875-32e1d070-ab53-4ac5-92fd-bb8ed16dd7dc.png">

# Posts & Articles & Reference

* Temporal adopted
  as [the first community drive DSL framework/abstraction](https://github.com/temporalio/awesome-temporal) of Temporal
* Cadence adopted in its [README](https://github.com/uber/cadence#cadence)
  , [official documentation](https://cadenceworkflow.io/docs/get-started/#what-s-next)
  and [Cadence community spotlight](https://cadenceworkflow.io/blog/2023/01/31/community-spotlight-january-2023/)
* [A Letter to Cadence/Temporal, and Workflow Tech Community](https://medium.com/@qlong/a-letter-to-cadence-temporal-and-workflow-tech-community-b32e9fa97a0c)
* [iWF vs Cadence/Temporal](https://medium.com/@qlong/iwf-vs-cadence-temporal-1e11b35960fe)
* [iWF vs other general purposed workflow Engines](https://medium.com/@qlong/iwf-vs-other-general-purposed-workflow-engines-f8f3e3d8993d)
* [Cadence® iWF](https://www.instaclustr.com/blog/cadence-iwf/?utm_content=1669999382&utm_medium=linkedin&utm_source=organicsocial)

