# iWF project - main & server repo

[![Go Reference](https://pkg.go.dev/badge/github.com/indeedeng/iwf.svg)](https://pkg.go.dev/github.com/indeedeng/iwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/indeedeng/iwf)](https://goreportcard.com/report/github.com/indeedeng/iwf)
[![Coverage Status](https://codecov.io/github/indeedeng/iwf/coverage.svg?branch=main)](https://app.codecov.io/gh/indeedeng/iwf/branch/main)

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

**iWF will make you a 10x developer!**

iWF is a platform for developing resilient, fault-tolerant, scalable long-running applications. 
It offers a convenient abstraction for durable timers, background execution with backoff retry, 
customized persisted data (with optional caching, and indexing), message queues, RPC, and more. You will build long-running reliable processes faster than ever. 

iWF is built on top of [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal).

Related projects:

* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk)
* [iWF Java Samples](https://github.com/indeedeng/iwf-java-samples)
* [iWF Golang SDK](https://github.com/indeedeng/iwf-golang-sdk)
* [iWF Golang Samples](https://github.com/indeedeng/iwf-golang-samples)
* WIP [iWF Python SDK](https://github.com/indeedeng/iwf-python-sdk)
* WIP [iWF TypeScript SDK](https://github.com/indeedeng/iwf-ts-sdk)
# What is iWF

## Example: microservice orchestration

### Problem
![1](https://github.com/indeedeng/iwf/assets/4523955/e0c7001e-2c8f-4a93-92d7-37e50a248c26)

As above diagram, you want to:
* Orchestrate 4 APIs as a workflow
* Each API needs backoff retry
* The data from topic 1 needs to be passed through
* API2 and API3+4 need to be in different parallel threads
* Need to wait for a signal from topic 2 for a day before calling API3
* If not ready after a day, call API4

This is a very abstracted example. It could be applied into any real-world scenario like refund process:
* API1: create a refund request object in DB
* API2: notify different users refund is created
* topic2: wait for approval
* API3: process refund after approval
* API4: notify timeout and expired

### Some existing solutions

With some other existing technologies, you solve it using message queue(like SQS which has timer) + Database like below:
![2](https://github.com/indeedeng/iwf/assets/4523955/babfca50-c605-4fae-b146-18d2aad79c6e)

* Using visibility timeout for backoff retry
  * Need to re-enqueue the message for larger backoff
* Using visibility timeout for durable timer
  * Need to re-enqueue the message for once to have 24 hours timer
* Need to create one topic for every step
* Need additional storage for processing ready signal
* Only go to 3 or 4 if both conditions are met
* Also need DLQ and build tooling around

It's complicated and hard to maintain and extend.   

### iWF solution
![3](https://github.com/indeedeng/iwf/assets/4523955/3428523e-c3d9-4fd6-8d10-c19b91ac7ecd)

The solution with iWF:
* All in one single dependency
* WorkflowAsCode
* Natural to represent business
* Builtin support/operation tooling

```java
public class OrchestrationWorkflow implements ObjectWorkflow {

    public static final String DA_DATA1 = "SomeData";
    public static final String READY_SIGNAL = "Ready";

    private List<StateDef> stateDefs;

    public OrchestrationWorkflow() {
        this.stateDefs = Arrays.asList(
                StateDef.startingState(new State1()),
                StateDef.nonStartingState(new State2()),
                StateDef.nonStartingState(new State3()),
                StateDef.nonStartingState(new State4())
        );
    }

    @Override
    public List<StateDef> getWorkflowStates() {
        return stateDefs;
    }

    @Override
    public List<PersistenceFieldDef> getPersistenceSchema() {
        return Arrays.asList(
                DataAttributeDef.create(String.class, DA_DATA1)
        );
    }
}

class State1 implements WorkflowState<String> {

    @Override
    public Class<String> getInputType() {
        return String.class;
    }

    @Override
    public StateDecision execute(final Context context, final String input, final CommandResults commandResults, Persistence persistence, final Communication communication) {
        persistence.setDataAttribute(DA_DATA1, input);
        System.out.println("call API1 with backoff retry in this method..");
        return StateDecision.multiNextStates(State2.class, State3.class);
    }
}

class State2 implements WorkflowState<Void> {

    @Override
    public Class<Void> getInputType() {
        return Void.class;
    }

    @Override
    public StateDecision execute(final Context context, final Void input, final CommandResults commandResults, Persistence persistence, final Communication communication) {
        String someData = persistence.getDataAttribute(DA_DATA1, String.class);
        System.out.println("call API2 with backoff retry in this method..");
        return StateDecision.deadEnd();
    }
}

class State3 implements WorkflowState<Void> {

    @Override
    public Class<Void> getInputType() {
        return Void.class;
    }

    @Override
    public CommandRequest waitUntil(final Context context, final Void input, final Persistence persistence, final Communication communication) {
        return CommandRequest.forAnyCommandCompleted(
                TimerCommand.createByDuration(Duration.ofHours(24)),
                SignalCommand.create(READY_SIGNAL)
        );
    }

    @Override
    public StateDecision execute(final Context context, final Void input, final CommandResults commandResults, final Persistence persistence, final Communication communication) {
        if (commandResults.getAllTimerCommandResults().get(0).getTimerStatus() == TimerStatus.FIRED) {
            return StateDecision.singleNextState(State4.class);
        }
        
        String someData = persistence.getDataAttribute(DA_DATA1, String.class);
        System.out.println("call API3 with backoff retry in this method..");
        return StateDecision.gracefulCompleteWorkflow();
    }
}

class State4 implements WorkflowState<Void> {

    @Override
    public Class<Void> getInputType() {
        return Void.class;
    }

    @Override
    public StateDecision execute(final Context context, final Void input, final CommandResults commandResults, Persistence persistence, final Communication communication) {
        String someData = persistence.getDataAttribute(DA_DATA1, String.class);
        System.out.println("call API4 with backoff retry in this method..");
        return StateDecision.gracefulCompleteWorkflow();
    }
}
```

## Basic Concepts


The top-level concept is **`ObjectWorkflow`** -- nearly any "object" can be an ObjectWorkflow, as long as it's long-lasting, at least a few seconds. 

A user application creates an ObjectWorkflow by implementing the Workflow interface, in one of the supported languages e.g.
[Java](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/ObjectWorkflow.java)
or [Golang](https://github.com/indeedeng/iwf-golang-sdk/blob/main/iwf/workflow.go).
An implementation of the interface is referred to as a `WorkflowDefinition` and consists of the components shown below:

| Name                                                                     | Description                                                                                                                                               | 
|:-------------------------------------------------------------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------| 
| [Data Attribute](#persistence)                                           | Persistence field to storing data                                                                                                                         | 
| [Search Attribute](#persistence)                                         | "Searchable data attribute" -- attribute data is persisted and also indexed in search engine backed by ElasticSearch or OpenSearch                        |
| [Workflow State](#workflow-state)                                        | A background execution unit. State is super powerful like a small workflow of two steps: *waitUntil* (optional) and *execute* with default infinite retry |
| [RPC](#rpc)                                                              | Remote procedure call. Invoked by a client, executed in worker, and can interact with data/search attributes, internal channel, and state execution       |
| [Signal Channel](#signal-channel-vs-rpc)                                 | Asynchronous message queue for the workflow object to receive messages from external sources                                                              |
| [Internal Channel](#internalchannel-synchronization-for-multi-threading) | "Internal Signal Channel" -- An internal message queue for workflow states/RPC                                                                            |

A workflow definition can be visualized like this:

![Example workflow diagram](https://user-images.githubusercontent.com/4523955/234424825-ff3673c0-af23-4eb7-887d-b1f421f3aaa4.png)

These are all the concepts that you need to build a super complicated workflow.
See this engagement workflow example in [Java](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/engagement)
or [Golang](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/engagement)
for how it looks in practice!

Below are detailed explanations of the concepts. 
These concepts are powerful, and also extremely simple to learn and use (as is the philosophy of iWF).

## Persistence

iWF let you store customized data as a database during the workflow execution. This eliminates the need to depend on a database to implement your workflow.

Your data are stored as Data Attributes and Search Attributes. Together both define the "persistence schema". 
The persistence schema is defined and maintained in the code along with other business logic.

Search Attributes work like infinite indexes in a traditional database. You
only need to specify which attributes should be indexed, without worrying about complications you might be used to in
a traditional database like the number of shards, and the order of the fields in an index.

Logically, the workflow definition displayed in the example workflow diagram will have a persistence schema as follows:

| Workflow Execution   | Search Attr A | Search Attr B | Data Attr C | Data Attr D |
|----------------------|---------------|:-------------:|------------:|------------:|
| Workflow Execution 1 | val 1         |     val 2     |       val 3 |       val 4 |
| Workflow Execution 2 | val 5         |     val 6     |       val 7 |       val 8 |
| ...                  | ...           |      ...      |         ... |         ... |

With Search attributes, you can write [customized SQL-like queries to find any workflow execution(s)](https://docs.temporal.io/visibility#search-attribute), just like using a database query.

Note that after workflows are closed(completed, timeout, terminated, canceled, failed), all the data retained in your persistence schema will be deleted once the configured retention period elapses.  

The iWF persistence is mainly for storing the workflow intermediate states/data. 
**It is important to not abuse iWF persistence for things like permanent storage, or for tracking/analytics purpose.**  

### Caching 
By default, remote procedure calls (RPCs) will load data/search attributes with the Cadence/Temporal [query API](https://docs.temporal.io/workflows#query), 
which is not optimized for very high request volume (~>100 requests per second) on a single workflow execution. Such request volumes could cause
too many history replays, especially when workflows are closed. This could in turn produce undesirable latency and load.

You can enable **caching** to support those high-volume requests.  

Note:
* With caching enabled read-after-write access will become *eventually consistent*, unless `bypassCachingForStrongConsistency=true` is set in RPC options
* Caching will introduce an extra event in history (upsertMemo operation for WorkflowPropertiesModified event) for updating the persisted data attributes
* Caching will be more useful for read-only RPC (no persistence.SetXXX API or communication API calls in RPC implementation) or GetDataAttributes API.
  * A read-only RPC can still invoke any other RPCs (like calling other microservices, or DB operation) in the RPC implementation
* Caching is currently only supported if the backend is Temporal, because [Cadence doesn't support mutable memo](https://github.com/uber/cadence/issues/3729)

## Workflow State
WorkflowState is how you implement your asynchronous process as a "workflow".  
It will run in the background, with infinite backoff retry by default. 
 
A WorkflowState is itself like “a small workflow” of 1 or 2 steps:

**[ `waitUntil` ] → `execute`**

The `waitUntil` API returns "commands" to wait for. When the commands are completed, the `execute` API will be invoked.
Both the `waitUntil` and `execute` APIs have access to read/write the persistence schema defined in the workflow.

The full execution flow looks like this:

![Workflow State diagram](https://user-images.githubusercontent.com/4523955/234921554-587d8ad4-84f5-4987-b838-959869293465.png)

The `waitUntil` API is optional. If not defined, then the `execute` API will be invoked immediately when the Workflow State is started.

Note: the `waitUntil` and `execute` APIs are invoked by the iWF service with infinite backoff retry by default. See the [WorkflowStateOptions](#WorkflowStateOptions) section for customization.  

The execute API will return a StateDecision:
* For a single next state a decision can 
  * Go to a different state
  * Go to the same state, i.e. a loop
  * Go to the previous state, i.e. a loop
* Dead end -- Just stop the thread
* Graceful complete -- Stop the thread, and also will stop the workflow when all other threads are stopped
* Force complete -- Stop the workflow immediately
* Force fail  -- Stop the workflow immediately with failure
* For multiple next states, these are executed in parallel as multiple threads

With StateDecisions, the workflow definitions can have flows like these:

![decision flow1](https://user-images.githubusercontent.com/4523955/234919901-f327dfb6-5b38-4440-a2eb-5d1c832b694e.png)

or 

![decision flow2](https://user-images.githubusercontent.com/4523955/234919896-30db8628-daeb-4f1d-bd2b-7bf826989c75.png)

or as complex as needed for any use case!


### Commands for WorkflowState's WaitUntil API

iWF provides three types of commands:

* `SignalCommand` -- Wait for a signal to be published to the workflow signal channel. External applications can use
  SignalWorkflow API to signal a workflow.
* `TimerCommand` -- Wait for a **durable timer** to fire.
* `InternalChannelCommand` -- Wait for a message from InternalChannel.

The `waitUntil` API can return multiple commands along with a `CommandWaitingType`:

* `AllCommandCompleted` -- Wait for all commands to be completed.
* `AnyCommandCompleted` -- Wait for any of the commands to be completed.
* `AnyCommandCombinationCompleted` -- Wait for any combination of the commands in a specified list to be completed.

### InternalChannel: synchronization for multi-threading
When there are multiple threads of workflow states running in parallel, you may want to have them wait on each other to ensure some particular ordering.

For example, in your problem space, WorkflowStates 1,2,3 need to be completed before WorkflowState 4. 

In this case, you need to utilize the "InternalChannel". WorkflowState 4 should be waiting on an "InternalChannel" for 3 messages via the `waitUntil` API. 
WorkflowState 1,2,3 will each publish a message when completing. This ensures propper ordering.  

## RPC

RPC stands for "Remote Procedure Call". Allows external systems to interact with the workflow execution.

It's invoked by client, executed in workflow worker, and then respond back the results to client. 

RPC can have access to not only persistence read/write API, but also interact with WorkflowStates using InternalChannel, 
or trigger a new WorkflowState execution in a new thread.

### Atomicity of RPC APIs

It's important to note that in addition to read/write persistence fields, a RPC can **trigger new state executions, and publish message to InternalChannel, all atomically.**

Atomically sending internal channel, or triggering state executions is an important pattern to ensure consistency across dependencies for critical business – this 
solves a very common problem in many existing distributed system applications. Because most RPCs (like REST/gRPC/GraphQL) don't provide a way to invoke 
background execution when updating persistence. People sometimes have to use complicated design to acheive this. 

**But in iWF, it's all builtin, and user application just needs a few lines of code!** 

![flow with RPC](https://user-images.githubusercontent.com/4523955/234930263-40b98ca7-4401-44fa-af8a-32d5ae075438.png)

Note that the read+write locking in RPC is [WIP](https://github.com/indeedeng/iwf/issues/305).

### Signal Channel vs RPC

There are two major ways for external clients to interact with workflows: Signal and RPC. So what are the difference? 

They are completely different:
* Signal is sent to iWF service without waiting for response of the processing
* RPC will wait for worker to process the RPC request synchronously
* Signal will be held in a signal channel until a workflow state consumes it
* RPC will be processed by worker immediately

![signals vs rpc](https://user-images.githubusercontent.com/4523955/234932674-b0d062b2-e5dd-4dbe-93b5-1b9863acc5e0.png)

So choose based on the situations/requirements

|                |        Availability        |                                        Latency | Workflow Requirement                                     |
|----------------|:-------------------------- |:----------------------------------------------- |:---------------------------------------------------------|
| Signal Channel |            High            |                                            Low | Requires a WorkflowState to process                      |
| RPC            | Depends on workflow worker | Higher than signal, depends on workflow worker | No WorkflowState required                                |

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
* `ALLOW_IF_PREVIOUS_EXITS_ABNORMALLY`
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

By default, the API timeout is 30s with infinite backoff retry. 
Users can customize the API timeout and retry policy:

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

When a workflowState/RPC API loads DataAttributes/SearchAttributes, by default it will use `LOAD_ALL_WITOUT_LOCKING` to load everything.

For WorkflowState, there is a 2MB limit by default to load data. User can use another loading policy `LOAD_PARTIAL_WITHOUT_LOCKING`
to specify certain DataAttributes/SearchAttributes only to load.

`WITHOUT_LOCKING` here means if multiple StateExecutions/RPC try to upsert the same DataAttribute/SearchAttribute, they can be
done in parallel without locking.

If racing conditions could be a problem, using`PARTIAL_WITH_EXCLUSIVE_LOCK` allows specifying some keys to be locked during the execution.

However, `PARTIAL_WITH_EXCLUSIVE_LOCK` is not supported in RPC yet. The feature is WIP with waiting for the Temporal "update" being production ready.
As a workaround, RPC can kick off a WorkflowState to update the persistence data using the locking policy.

#### WaitUntil API failure policy

By default, the workflow execution will fail when API max out the retry attempts. In some cases that
workflow want to ignore the errors.

Using `PROCEED_ON_API_FAILURE` for `WaitUntilApiFailurePolicy` will let workflow continue to invoke `execute`
API when the API fails with maxing out all the retry attempts.

Alternatively, WorkflowState can utilize `attempts` or `firstAttemptTime` from the context to decide ignore the exception/error.

## Limitation

Though iWF can be used for a very wide range of use case even just CRUD, iWF is NOT for everything. It is not suitable for use cases like:

* High performance transaction( e.g. within 10ms)
* High frequent writes on a single workflow execution(like a single record in database) for hot partition issue
  * High frequent reads on a single workflow execution is okay if using memo for data attributes
* Join operation across different workflows
* Transaction for operation across multiple workflows


# Architecture

An iWF application is composed of several iWF workflow workers. These workers host REST APIs as "worker APIs" for server to call. This callback pattern similar to AWS Step Functions invoking Lambdas, if you are familiar with.

An application also perform actions on workflow executions, such as starting, stopping, signaling, and retrieving results 
by calling iWF service APIs as "service APIs".

The service APIs are provided by the "API service" in iWF server. Internally, this API service communicates with the Cadence/Temporal service as its backend.

In addition, the iWF server also runs the Cadence/Temporal workers as "worker service". The worker service
hosts [an interpreter workflow](https://github.com/indeedeng/iwf/blob/main/service/interpreter/workflowImpl.go).
This workflow implements all the core features as described above, and also things like "Auto ContinueAsNew" to let you use 
iWF without any scaling limitation. 

![architecture diagram](https://user-images.githubusercontent.com/4523955/234935630-e69c648e-7714-4672-beb2-d9867bedf940.png)

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

* All the input/output to your workflow are stored in the activity input/output of history event. The input is
  in `ActivityTaskScheduledEvent`, output is in `ActivityTaskCompletedEvent` or in pending activity view if having
  errors.
* Use query handlers like (`GetDataObjects` or `GetCurrentTimerInfos`) in Cadence/Temporal WebUI to quickly understand
  the current status of the workflows.
    * DumpAllInternal will return all the internal status or the pending states
    * GetCurrentTimerInfos will return all the timers of the pending states
* Let your worker service return error stacktrace as the response body to iWF server. E.g.
  like [this example of Spring Boot using ExceptionHandler](https://github.com/indeedeng/iwf-java-samples/blob/2d500093e2aaecf2d728f78366fee776a73efd29/src/main/java/io/iworkflow/controller/IwfWorkerApiController.java#L51)
  .
* If you return the full stacktrace in response body, the pending activity view will show it to you! Then use
  Cadence/Temporal WebUI to debug your application.


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
