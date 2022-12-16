# iWF project - main & server repo 

iWF is a platform providing an all-in-one tooling for building long-running business application. It provides an abstraction for persistence(database, elasticSearch) and more! It aims to provide clean, simple and easy to use interface, like an iPhone. 

**It will not make you a 10x developer...but you may feel like one!**

We call _long running application_ **`Workflow`**. 

It's a simple and powerful WorkflowAsCode general purpose workflow engine.

The server is back by [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal) as an interpreter.

Related projects:
* [API definition between SDKs and server](https://github.com/indeedeng/iwf-idl). Any languages can be supported as long as implementing it. 
* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk) 
* [iWF Java Samples](https://github.com/indeedeng/iwf-java-samples)
* [iWF Golang SDK](https://github.com/cadence-oss/iwf-golang-sdk), WIP, Contribution is welcome.
* More SDKs? Contribution is welcome.

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
- [Development](#development)
  - [How to update IDL and the generated code](#how-to-update-idl-and-the-generated-code)
  - [Run with local Temporalite](#run-with-local-temporalite)
  - [Run with local Cadence](#run-with-local-cadence)
  - [Development Plan](#development-plan)
- [Some history](#some-history)

# Community & Help
* [Github Discussion](https://github.com/indeedeng/iwf/discussions)
  * Best for Q&A, support/help, general discusion, and annoucement 
* [StackOverflow](https://stackoverflow.com/questions/tagged/iwf)
  * Best for Q&A and general discusion
* [Github Issues](https://github.com/indeedeng/iwf/issues)
  * Best for reporting bugs and feature requests
* For any questions & consultant email to: qlong.seattle@gmail.com

# What is iWF

## Architecture
A iWF application will host a set of iWF workflow workers. The workers host two REST APIs of WorkflowState `start` and `decide` using iWF SDKs.
The application will call iWF server to interact with workflow executions -- start, stop, signal, get results, etc, using iWF SDKs.

iWF server hosts those APIs(also REST) as a iWF API service. The API service will call Cadence/Temporal service as the backend.

iWF server also hosts Cadence/Temporal workers which hosts [an interpreter workflow](https://github.com/indeedeng/iwf/blob/main/service/interpreter/workflowImpl.go).
Any iWF workflows are interpreted into this Cadence/Temporal workflow. The interpreter workflow will invoke the two iWF APIs of
the application workflow workers. Internally, the two APIs are executed by Cadence/Temporal activity. Therefore, all the REST API request/response with the worker are stored in history events which are useful for debugging/troubleshooting. 

![architecture diagram](https://user-images.githubusercontent.com/4523955/207514928-56fea636-c711-4f20-9e90-94ddd1c9844d.png)

## Basic Concepts

### Workflow and WorkflowState definition
iWF lets you build long-running applications by implementing the workflow interface, e.g. [Java Workflow interface](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/Workflow.java).
An instance of the interface is a `WorkflowDefinition`. User applications use `iwfWorkflowType` to differentiate WorkflowDefinitions.    

A WorkflowDefinition contains several `WorkflowState` e.g. [Java WorkflowState interface](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/WorkflowState.java). 
A WorkflowState is implemented with two APIs: `start` and `decide`. 
* `start` API is invoked immediately when a WorkflowState is started. It will return some `Commands` to server. When the requested `Commands` are completed, `decide` API will be triggered. 
* `decide` API will decide next states to execute. Next states be multiple, and can be re-executed as different `stateExecutions`. 

### Workflow execution and WorkflowState execution
Application can start a workflow instance with a `workflowId` for any workflow definition. A workflow instance is called `WorkflowExecution`. 
iWF server returns `runId` of UUID as the identifier of the WorkflowExecution. The runId is globally unique.  

WorkflowId uniqueness: At anytime, there must be at most one WorkflowExecution running with the same workflowId. However, after a previous WorkflowExecution finished running (in any closed status), 
application may start a new WorkflowExecutions with the same workflowId using appropriate `IdReusePolicy`. 

There must be at least one WorkflowState being executed for a running WorkflowExecution. The instance of WorkflowState is called `StateExecution`.     

:warning: Note:
> Depends on the context, the only word `workflow` may mean WorkflowExecution(most commonly), WorkflowDefinition or both.  

### Commands
These are the three command types:
* `SignalCommand`: will be waiting for a signal from external to the workflow signal channel. External application can use SignalWorkflow API to signal a workflow. 
* `TimerCommand`: will be waiting for a **durable timer** to fire.
* `InterStateChannelCommand`: will be waiting for a value being published from another state execution(internally in the same workflow execution)

Note that `start` API can return multiple commands, and choose different DeciderTriggerType for triggering decide API:
* `AllCommandCompleted`: this will wait for all command completed
* `AnyCommandCompleted`: this will wait for any command completed

### Persistence
iWF provides super simple persistence abstraction for workflow to use. Developers don't need to touch any database system to register/maintain the schemas. The only schema is defined in the workflow code.
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
* Get workflow: get basic information about a workflow like status
* Get workflow data objects: get the dataObjects of a workflow execution
* Get workflow search attributes: get the search attributes of a workflow execution
* Reset workflow: reset a workflow to previous states
* Get workflow results: get the workflow completion results (with or without waiting for completion)

# Why iWF

## If you are familiar with Cadence/Temporal
* See [Slide deck](https://docs.google.com/presentation/d/1CpsroSf6NeVce_XyUhFTkd9bLHN8UHRtM9NavPCMhj8/edit#slide=id.gfe2f455492_0_56) for what problems it is solving
* See [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit) for how it works  

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

## How to build & run locally
* Run `make bins` to build the binary `iwf-server`
* Then run  `./iwf-server start` to run the service . This defaults to serve workflows APIs with Temporal interpreter implementation. It requires to have local Temporal setup. See Run with local Temporal.
* Alternatively, run `./iwf-server --config config/development_cadence.yaml start` to run with local Cadence. See below instructions for setting up local Cadence. 
* Run `make integTests` to run all integration tests. This by default requires to have both local Cadence and Temporal to be set up.

## How to use in production

You can customize the docker image, or just use the [api](https://github.com/indeedeng/iwf/tree/main/service/api) and [interpreter](https://github.com/indeedeng/iwf/tree/main/service/interpreter) that are exposed as the api service and workflow service.

For more info, contact qlong.seattle@gmail.com

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
tctl adm cl asa -n IwfGlobalWorkflowVersion -t Int
tctl adm cl asa -n IwfExecutingStateIds -t Keyword

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
cadence adm cl asa --search_attr_key IwfGlobalWorkflowVersion --search_attr_type 2
cadence adm cl asa --search_attr_key IwfExecutingStateIds --search_attr_type 0
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
- [x] DataObjectRW
- [x] StateLocal
- [x] Signal workflow API
- [x] Get DataObjects/SearchAttributes API
- [x] Get workflow info API
- [x] Search workflow API
- [x] Stop workflow API

### 1.1
- [x] Reset workflow API (Cadence only, TODO for Temporal)
- [x] Command type(s) for inter-state communications (e.g. internal channel)
- [x] AnyCommandCompleted Decider trigger type
- [x] More workflow start options: IdReusePolicy, cron schedule, retry
- [x] StateOption: Start/Decide API timeout and retry policy
- [ ] More workflow start options: initial search attributes/memo
- [ ] Reset workflow by stateId

### 1.2
- [ ] Decider trigger type: AnyCommandClosed
- [ ] WaitForMoreResults in StateDecision
- [ ] Skip timer API for testing/operation
- [ ] LongRunningActivityCommand
- [ ] Failing workflow details
- [ ] Auto ContinueAsNew 
- [ ] StateOption: more PersistenceLoadingPolicy
- [ ] StateOption: more CommandCarryOverPolicy

# Some history
AWS published SWF in 2012 and then moved to Step Functions in 2016 because they found it’s too hard to support SWF.
Cadence & Temporal continued the idea of SWF and became much more powerful.
However, AWS is right that the programming of SWF/Cadence/Temporal is hard to adopt because of leaking too many internals.
Inspired by Step Function, iWF is created to provide equivalent power of Cadence/Temporal, but hiding all the internal details
and provide clean and simple API to use. 

Read this [doc](https://docs.google.com/document/d/1zyCKvy4S2l7XBVJzZuS65OIsqV9CRPPYJY3OBbuWrPE) for more.

<img width="916" alt="history diagram" src="https://user-images.githubusercontent.com/4523955/201188875-32e1d070-ab53-4ac5-92fd-bb8ed16dd7dc.png">
