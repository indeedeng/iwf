# iWF project - main & server repo

[![Go Reference](https://pkg.go.dev/badge/github.com/indeedeng/iwf.svg)](https://pkg.go.dev/github.com/indeedeng/iwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/indeedeng/iwf)](https://goreportcard.com/report/github.com/indeedeng/iwf)
[![Coverage Status](https://codecov.io/github/indeedeng/iwf/coverage.svg?branch=main)](https://app.codecov.io/gh/indeedeng/iwf/branch/main)

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

**iWF will make you a 10x developer!**

iWF is an all-in-one platform for developing long-running business processes. It offers a convenient abstraction for durable timers, 
background execution with backoff retry, persistence, indexing, message queues, RPC, and more. You will build reliable/scalable backend applications
much faster than ever. 

iWF is built on top of [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal).

Related projects:

* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk)
* [iWF Java Samples](https://github.com/indeedeng/iwf-java-samples)
* [iWF Golang SDK](https://github.com/iworkflowio/iwf-golang-sdk) (V2 is WIP)
* [iWF Golang Samples](https://github.com/iworkflowio/iwf-golang-samples)

# What is iWF


## Basic Concepts


The top level concept is **`ObjectWorkflow`** -- anything can be an ObjectWorkflow, as long as it's long-lasting, at least a few seconds. 

User application creates ObjectWorkflow by implementing the Workflow interface, e.g. in
[Java](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/ObjectWorkflow.java).
An implementation of the interface is referred to as a `WorkflowDefinition`, consisting below components:

| Name             |                                                                 Description                                                                  | 
|:------------------|:-------------------------------------------------------------------------------------------------------------------------------------------- | 
| Data Attribute   |                                                      Persistence field to storing data                                                       | 
| Search Attribute |                                                         "Searchable data attribute"                                                          | 
| Signal Channel   |                                       Asynchronous message queue for the workflow object for external                                        |
| Internal Channel |                                              An internal message queue for workflow states/RPC                                               |
| Workflow State   |           A background execution unit. State is super powerful like a small workflow of two steps: waitUntil(optional) and execute           |
| RPC              | Remote procedure call. Invoked by client, executed in worker, and interact with data/search attributes, internal channel and state execution |

You can use a diagram to outline a workflow definition like this:

![Example workflow diagram](https://user-images.githubusercontent.com/4523955/234424825-ff3673c0-af23-4eb7-887d-b1f421f3aaa4.png)

Logically, this workflow definition will have a persistence schema like below:

| Workflow Execution   | Search Attr A | Search Attr B | Data Attr C | Data Attr D |
|----------------------|---------------|:-------------:|------------:|------------:|
| Workflow Execution 1 | val 1         |     val 2     |       val 3 |       val 4 |
| Workflow Execution 2 | val 5         |     val 6     |       val 7 |       val 8 |
| ...                  | ...           |      ...      |         ... |         ... |

And the schema just defined and maintained in your code along with other business logic.

## Workflow State
A workflow state is like “a small workflow” of 1~2 steps:

**[ waitUntil ] → execute**

The full detailed execution flow is like this:

![Workflow State diagram](https://user-images.githubusercontent.com/4523955/234427642-0a9e9332-0587-44f5-a71d-175ebb03c170.png)

The execute API will return some StateDecision:
* Single next state 
  * Go to to different state
  * Go to the same state as a loop
  * Go the the previous state as a loop
* Multiple next states, executing as multi threads in parallel
* Dead end -- Just stop the thread
* Graceful complete -- Stop the thread, and also will stop the workflow when all other threads are stopped
* Force complete -- Stop the workflow immediately
* Force fail  -- Stop the workflow immediately with failure

With decisions, a "complex" workflow definitions can have a flow like this:

![decision flow1](https://user-images.githubusercontent.com/4523955/234428066-629453a6-e385-47cf-9408-835f5aaf4b3a.png)

or

![decision flow2](https://user-images.githubusercontent.com/4523955/234428082-649be7f4-a699-406c-91cc-d8d25a41ae60.png)

If combining with some commands, it can be like this:

![decision flow3](https://user-images.githubusercontent.com/4523955/234428326-a697cc35-31d6-4b94-9d4c-fbf65474ecf6.png)


### Commands for WorkflowState's WaitUntil API

iWF provides three types of commands:

* `SignalCommand`: will wait for a signal to be published to the workflow signal channel. External application can use
  SignalWorkflow API to signal a workflow.
* `TimerCommand`: will wait for a **durable timer** to fire.
* `InternalChannelCommand`: will wait for a message from InternalChannel.

The waitUntil API can return multiple commands along with a `CommandWaitingType`:

* `AllCommandCompleted`: This option waits for all commands to be completed.

* `AnyCommandCompleted`: This option waits for any of the commands to be completed.

* `AnyCommandCombinationCompleted`: This option waits for any combination of the commands in a specified list to be
  completed.

## RPC

In addition to read/write persistence fields, a RPC can **trigger new state executions, and publish message to InternalChannel, all atomically.**


`RPC` triggering state executions  is an important pattern to ensure consistency across dependencies for critical business – this 
solves a very common problem in many existing distributed systems, almost everywhere.

![flow with RPC](https://user-images.githubusercontent.com/4523955/234428514-0dfaba96-91c6-4aa2-9fbb-f3e1904f3c24.png)

### Signal Channel vs RPC

They are completely different:
* Signal is sent to iWF service without waiting for response of the processing
* RPC will wait for worker to process the RPC request synchronously
* Signal will be held in a signal channel until a workflow state consumes it
* RPC will be processed by worker immediately

![signals vs rpc](https://user-images.githubusercontent.com/4523955/234428638-a0075124-1992-4d54-a15b-69a037b4f8fa.png)

| vs             |        Availability        |                                        Latency |                                    Workflow Requirement |
|----------------|:-------------------------- |:----------------------------------------------- |:-------------------------------------------------------- |
| Signal Channel |            High            |                                            Low |                     Requires a WorkflowState to process |
| RPC            | Depends on workflow worker | Higher than signal, depends on workflow worker |                               No WorkflowState required |

## Advanced Customization

### WorkflowOptions

iWF let you deeply customize the workflow behaviors with the below options.

#### IdReusePolicy for WorkflowId

At any given time, there can be only one WorkflowExecution running for a specific workflowId.
A new WorkflowExecution can be initiated using the same workflowId by setting the appropriate `IdReusePolicy` in
WorkflowOptions.

* `ALLOW_IF_NO_RUNNING` 
    * Allow starting workflow if there is no execution running with the workflowId
    * This is the **default policy** if not specified in WorkflowOptions
* `ALLOW_IF_PREVIOUS_EXISTS_ABNORMALLY`
    * Allow starting workflow if a previous Workflow Execution with the same Workflow Id does not have a Completed
      status.
      Use this policy when there is a need to re-execute a Failed, Timed Out, Terminated or Cancelled workflow
      execution.
* `DISALLOW_REUSE` 
    * Not allow to start a new workflow execution with the same workflowId.
* `ALLOW_TERMINATE_IF_RUNNING`
    * Always allow starting workflow no matter what -- iWF server will terminate the current running one if it exists.

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

#### WorkflowState WaitUntil/Execute API timeout and retry policy

Users can customize the API timeout and retry policy.

By default, the API timeout is 30s with infinite backoff retry:

- InitialIntervalSeconds: 1
- MaxInternalSeconds:100
- MaximumAttempts: 0
- MaximumAttemptsDurationSeconds: 0
- BackoffCoefficient: 2

Where zero means infinite attempts.

Both MaximumAttempts and MaximumAttemptsDurationSeconds are used for controlling the maximum attempts for the retry
policy.
MaximumAttempts is directly by number of attempts, where MaximumAttemptsDurationSeconds is by the total time duration of
all attempts including retries. It will be capped to the minimum if both are provided.

#### Persistence loading policy

When a state API loads DataObjects/SearchAttributes, by default it will load everything which could cause size limit
error
for Cadence/Temporal activity input/output limit(2MB by default). User can use other loading
policy `LOAD_PARTIAL_WITHOUT_LOCKING`
to specify certain DataObjects/SearchAttributes only to load for this WorkflowState.

`WITHOUT_LOCKING` here means if multiple StateExecutions try to upsert the same DataObject/SearchAttribute, they can be
done in parallel without locking.

#### WaitUntil API failure policy

By default, the workflow execution will fail when API max out the retry attempts. In some cases that
workflow want to ignore the errors.

Using `PROCEED_ON_API_FAILURE` for `WaitUntilApiFailurePolicy` will let workflow continue to execute decide
API when the API fails with maxing out all the retry attempts (therefore, you should override the default infinite
retry attempts to a different number).

Alternatively, WorkflowState can utilize `attempts` or `firstAttemptTime` from the context to decide ignore the
exception/error.

## Limitation

Though iWF can be used for a very wide range of use case even just CRUD, iWF is NOT for everything. It is not suitable for use cases like:

* High performance transaction( within 10ms)
* High throughput for a single object(like a single record in database) for hot partition issue
* Join operation across different workflows
* Transaction for operation across multiple workflows


# Architecture

An iWF application is composed of several iWF workflow workers. These workers host REST APIs for server to call. 
An application also perform actions on workflow executions, such as starting, stopping, signaling, and retrieving results 
by calling iWF service APIs.

The iWF server provides the APIs. Internally, this API service communicates with the Cadence/Temporal service as its backend.

In addition to hosting the iWF API service, the iWF server includes Cadence/Temporal workers that
host [an interpreter workflow](https://github.com/indeedeng/iwf/blob/main/service/interpreter/workflowImpl.go).
This interpreter workflow interprets any iWF workflows into the Cadence/Temporal workflow. 

![architecture diagram](https://user-images.githubusercontent.com/4523955/207514928-56fea636-c711-4f20-9e90-94ddd1c9844d.png)

# How to use

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

# Posts & Articles & Reference

* Temporal adopted
  as [the first community drive DSL framework/abstraction](https://github.com/temporalio/awesome-temporal) of Temporal
* Cadence adopted in its [README](https://github.com/uber/cadence#cadence)
  , [official documentation](https://cadenceworkflow.io/docs/get-started/#what-s-next)
  and [Cadence community spotlight](https://cadenceworkflow.io/blog/2023/01/31/community-spotlight-january-2023/)
