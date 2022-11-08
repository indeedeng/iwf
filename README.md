# iWF project - main & server repo 

For most long running applications, iWF is a All-In-One solution to replace Database, MessageQueue and ElasticSearch. 

**It will not make you a 10x developer...but you may feel like one!**

It's a simple and powerful WorkflowAsCode general purpose workflow engine.

Back by [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal) as an interpreter.

Related projects:
* [API definition between SDKs and server](https://github.com/indeedeng/iwf-idl). Any languages can be supported as long as implementing it. 
* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk) 
* [iWF Java Samples](https://github.com/indeedeng/iwf-java-samples)
* [iWF Golang SDK](https://github.com/cadence-oss/iwf-golang-sdk)

Contribution is welcome.   

# Table of contents
* [Why you would need iWF](#why-you-would-need-iwf)
  * [If you are familar with <a href="https://github.com/uber/cadence">Cadence</a>/<a href="https://github.com/temporalio/temporal">Temporal</a>](#if-you-are-familar-with-cadencetemporal)
  * [If you are not](#if-you-are-not)
* [What is iWF](#what-is-iwf)
  * [Advanced Concepts &amp; Usage](#advanced-concepts--usage)
    * [WorkflowStartOption](#workflowstartoption)
    * [WorkflowStateOption](#workflowstateoption)
    * [Reset Workflow](#reset-workflow)
* [How to run](#how-to-run)
  * [Using docker image](#using-docker-image)
  * [How to build &amp; run locally](#how-to-build--run-locally)
  * [How to use in production](#how-to-use-in-production)
    * [Option 1: use as library to customize your startup](#option-1-use-as-library-to-customize-your-startup)
* [Development](#development)
  * [How to update IDL and the generated code](#how-to-update-idl-and-the-generated-code)
  * [Run with local Temporalite](#run-with-local-temporalite)
  * [Run with local Cadence](#run-with-local-cadence)
  * [Development Plan](#development-plan)
    * [1.0](#10)
    * [1.1](#11)
    * [1.2](#12)

# Why you would need iWF

## If you are familar with [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal)
* See [Slide deck](https://docs.google.com/presentation/d/1CpsroSf6NeVce_XyUhFTkd9bLHN8UHRtM9NavPCMhj8/edit#slide=id.gfe2f455492_0_56) for what problems it is solving
* See [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit) for how it works  

## If you are not
* Check out this [doc](https://docs.google.com/document/d/1zyCKvy4S2l7XBVJzZuS65OIsqV9CRPPYJY3OBbuWrPE) to understand some history

iWF is a application platform that provides you a comprehensive tooling:
* WorkflowAsCode for highly flexibile/customizable business logic
* Parallel execution of multiple threads of business
* Intermidiate states stored as "QueryAttributes" and can be retrived by APIs
* Receiving data from external system by Signal
* Searchable attributes that can be used for flexible searching, even full text searching, backed by ElasticSearch 
* Durable timer, and cron job scheduling
* Reset workflow to let you recover the workflows from bad states easily 
* ...
* Basically with iWF, you mostly don't need any database, messageQueue or even ElasticSearch to build your applications

# What is iWF

iWF lets you build long-running applications by implementing the workflow interface, e.g. [Java Workflow interface](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/github/cadenceoss/iwf/core/Workflow.java).

The key elements of a workflow are `WorkflowState`. A workflow can contain any number of WorkflowStates. 

A WorkflowState is implemented with two APIs: `start` and `decide`. `start` API is invoked immediately when a WorkflowState is started. It will return some `Commands` to server. When the requested `Commands` are completed, `decide` API will be triggered.

There are several types of commands:
* `SignalCommand`: will be waiting for a signal from external to the workflow signal channel. 
* `TimerCommand`: will be waiting for a durable timer to fire.
* `InterStateChannelCommand`: will be waiting for a value being published from another state(internally in the same workflow)
* `LongRunninngActivityCommand`: will schedule a Cadence/Temporal activity. This is only necessary for long-running activity like hours/days. 

Note that `start` API can return multiple commands, and choose different DeciderTriggerType for triggering decide API:
* `AllCommandCompleted`: this will wait for all command completed
* `AnyCommandCompleted`: this will wait for any command completed

iWF provides the below primitives when implementing the WorkflowState:
* `StateLocals` is for 
  * passing some data values from state API to decide API in the same WorkflowState execution
  * recording some events that can be useful for debugging using Workflow history. Usually you may want to record the input/output of the dependency RPC calls.  
* `QueryAttribute` is for 
  * sharing some data values across the workflow
  * can be retrieved by external application using GetQueryAttributes API
  * can be viewed in Cadence/Temporal WebUI in QueryHandler tab
* `SearchAttribute` is similarly:
  * sharing some data values across the workflow
  * can be retrieved by external application using GetSearchAttributes API
  * search for workflows by external application using `SearchWorkflow` API
  * search for workflows in Cadence/Temporal WebUI in Advanced tab
  * search attribute type must be registered in Cadence/Temporal server before using for searching because it is backed up ElasticSearch
  * the data types supported are limited as server has to understand the value for indexing

## Advanced Concepts & Usage
On top of the above basic concepts, you may want to deeply customize your workflow by using the below features.

### WorkflowStartOption
* IdReusePolicy
* CronSchedule
* RetryPolicy
### WorkflowStateOption
* AttributeLoadingPolicy
* API timeout & retry
### Reset Workflow



# How to run

##  Using docker image
TODO

## How to build & run locally
* Run `make bins` to build the binary `iwf-server`
* Then run  `./iwf-server start` to run the service . This defaults to serve workflows APIs with Temporal interpreter implementation. It requires to have local Temporal setup. See Run with local Temporal.
* Alternatively, run `./iwf-server --config config/development_cadence.yaml start` to run with local Cadence. See below instructions for setting up local Cadence. 
* Run `make integTests` to run all integration tests. This by default requires to have both local Cadence and Temporal to be set up.

## How to use in production

### Option 1: use as library to customize your startup 

Particularly, use the [api](https://github.com/indeedeng/iwf/tree/main/service/api) and [interpreter](https://github.com/indeedeng/iwf/tree/main/service/interpreter) that are exposed as the api service and workflow service.

# Development

Any contribution is welcome.

### How to update IDL and the generated code
1. Install openapi-generator using Homebrew if you haven't. See more [documentation](https://openapi-generator.tech/docs/installation) 
2. Check out the idl submodule by running the command: `git submodule update --init --recursive`
3. Run the command `git submodule update --remote --merge` to update IDL to the latest commit
4. Run `make idl-code-gen` to refresh the generated code


### Run with local Temporalite
1. Run a local Temporalite following the [instruction](https://github.com/temporalio/temporalite). If you see error `error setting up schema`, try use command `temporalite start --namespace default -f my_test.db` instead to start. 
2. Register a default namespace
```shell
tctl --ns default n re
```
3. Go to http://localhost:8233/ for Temporal WebUI

NOTE: alternatively, go to [Temporal-dockercompose](https://github.com/temporalio/docker-compose) to run with docker

3. Register system search attributes required by iWF server
```shell
tctl adm cl asa -n IwfWorkflowType -t Keyword
tctl adm cl asa -n GlobalWorkflowVersion -t Int
tctl adm cl asa -n StateExecutionStatus -t Keyword

```
4 For `attribute_test.go` integTests, you need to register search attributes:
```shell
tctl adm cl asa -n CustomKeywordField -t Keyword
tctl adm cl asa -n CustomIntField -t Int
```

### Run with local Cadence
1. Run a local Cadence server following the [instructions](https://github.com/uber/cadence/tree/master/docker)
```
docker-compose -f docker-compose-es-v7.yml up
```
2. Register a new domain if not haven `cadence --do default domain register`
3. Register system search attributes required by iWF server
```
cadence adm cl asa --search_attr_key GlobalWorkflowVersion --search_attr_type 2
cadence adm cl asa --search_attr_key StateExecutionStatus --search_attr_type 0
cadence adm cl asa --search_attr_key IwfWorkflowType --search_attr_type 0
```
4. Go to Cadence http://localhost:8088/domains/default/workflows?range=last-30-days

## Development Plan
### 1.0
- [x] Start workflow API
- [x] Executing `start`/`decide` APIs and completing workflow
- [x] Parallel execution of multiple states 
- [x] Timer command
- [x] Signal command
- [x] SearchAttributeRW
- [x] QueryAttributeRW
- [x] StateLocalAttribute
- [x] Signal workflow API
- [x] Query workflow API
- [x] Get workflow API
- [x] Search workflow API
- [x] Cancel workflow API

### 1.1
- [x] Reset workflow API (Cadence only, TODO for Temporal)
- [x] Command type(s) for inter-state communications (e.g. internal channel)
- [x] AnyCommandCompleted Decider trigger type
- [ ] More workflow start options: IdReusePolicy, initial earch attributes, cron schedule, retry, etc
- [ ] StateOption: Start/Decide API timeout and retry
- [ ] Reset workflow by stateId

### 1.2
- [ ] Decider trigger type: AnyCommandClosed
- [ ] WaitForMoreResults in StateDecision
- [ ] Skip timer API for testing/operation
- [ ] LongRunningActivityCommand
- [ ] Failing workflow details
- [ ] Auto ContinueAsNew 
- [ ] StateOption: more AttributeLoadingPolicy
- [ ] StateOption: more CommandCarryOverPolicy
