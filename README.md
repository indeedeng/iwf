atest# iWF project - main & server repo

[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](http://iworkflow-slack.work)
[![Go Reference](https://pkg.go.dev/badge/github.com/indeedeng/iwf.svg)](https://pkg.go.dev/github.com/indeedeng/iwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/indeedeng/iwf)](https://goreportcard.com/report/github.com/indeedeng/iwf)
[![Coverage Status](https://codecov.io/github/indeedeng/iwf/coverage.svg?branch=main)](https://app.codecov.io/gh/indeedeng/iwf/branch/main)

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

**iWF will make you a 10x developer!**

iWF is an API orchestration platform for building resilient, fault-tolerant, scalable long-running applications. 
It offers an orchestration **coding framework** with abstractions for durable timers, async/background execution with backoff retry, 
KV storage,  RPC, and message queues. You will build long-running reliable processes faster than ever.


# What is [iWF](https://github.com/indeedeng/iwf/wiki)

## Use case study/examples
* [SAGA pattern](https://medium.com/@qlong/saga-pattern-deep-dive-with-indeed-workflow-engine-b7e82c59e51f?sk=672abd70b0e092d4cda7788276c5a241)
  * [Java samples](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/money/transfer), [Golang samples](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/moneytransfer), [Python samples](https://github.com/indeedeng/iwf-python-samples/tree/main/moneytransfer)
* [User sign-up/registry in Python/Java](https://github.com/indeedeng/iwf/wiki/Use-case-study-%E2%80%90%E2%80%90-user-signup-workflow)
* [Abstracted microservice orchestration in Java/Golang](https://github.com/indeedeng/iwf/wiki/Use-case-study-%E2%80%90%E2%80%90-Microservice-Orchestration)
* Employer & JobSeeker engagement in [Java](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/engagement) or [Golang](https://github.com/indeedeng/iwf-golang-samples/blob/main/workflows/engagement)
* Subscription Workflow in [Java](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/subscription) or [Golang](https://github.com/indeedeng/iwf-golang-samples/blob/main/workflows/subscription)

## Basic concepts
* [Basic concepts overview](https://github.com/indeedeng/iwf/wiki/Basic-concepts-overview)
* [WorkflowState](https://github.com/indeedeng/iwf/wiki/WorkflowState)
* [RPC](https://github.com/indeedeng/iwf/wiki/RPC)
* [Persistence](https://github.com/indeedeng/iwf/wiki/Persistence)

## Advanced concepts
* [WorkflowOptions](https://github.com/indeedeng/iwf/wiki/WorkflowOptions)
* [WorkflowStateOptions](https://github.com/indeedeng/iwf/wiki/WorkflowStateOptions)
* [Persistence Caching](https://github.com/indeedeng/iwf/wiki/Persistence-Caching)

# How to use

As a coding framework, iWF provides three SDKs to use with:

* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk) and [samples](https://github.com/indeedeng/iwf-java-samples)
* [iWF Golang SDK](https://github.com/indeedeng/iwf-golang-sdk) and [samples](https://github.com/indeedeng/iwf-golang-samples)
* [iWF Python SDK](https://github.com/indeedeng/iwf-python-sdk) and [samples](https://github.com/indeedeng/iwf-python-samples)

The iWF SDKs required to run with the server:

## Using all-in-one docker image

Run the docker command to start the container for:
* IWF service: http://localhost:8801/
* Temporal WebUI: http://localhost:8233/
* Temporal service: localhost:7233
```shell
docker run -p 8801:8801 -p 7233:7233 -p 8233:8233 -e AUTO_FIX_WORKER_URL=host.docker.internal --add-host host.docker.internal:host-gateway -it iworkflowio/iwf-server-lite:latest
```

To update the server version, use `docker pull iworkflowio/iwf-server-lite:latest` to get the latest image. Or change `latest` to specify a version tag.

## Using docker image & docker-compose

This is the simpler(preferred) option to run the server locally for development.

Checkout this repo, go to the docker-compose folder and run it:

```shell
cd docker-compose && docker-compose up
```

This by default will run Temporal server with it, again:
* IWF service: http://localhost:8801/
* Temporal WebUI: http://localhost:8233/
* Temporal service: localhost:7233

To update the server version, use `docker pull iworkflowio/iwf-server:latest` to get the latest image. Or update the docker-compose file to specify
a version tag.

## How to build & run locally

* Run `make bins` to build the binary `iwf-server`
* Make sure you have registered the system search attributes required by iWF server:
    * Keyword: IwfWorkflowType
    * Int: IwfGlobalWorkflowVersion
    * Keyword: IwfExecutingStateIds
    * See [Contribution](./CONTRIBUTING.md) for more detailed commands.
    * For Cadence without advancedVisibility enabled,
      set [disableSystemSearchAttributes](https://github.com/indeedeng/iwf/blob/main/config/development_cadence.yaml#L8)
      to true and [executingStateIdMode](https://github.com/indeedeng/iwf/blob/main/config/development_cadence.yaml#L9)
      to DISABLED
* Then run  `./iwf-server start` to run the service . This defaults to serve workflows APIs with Temporal interpreter
  implementation. It requires to have local Temporal setup. See Run with local Temporal.
* Alternatively, run `./iwf-server --config config/development_cadence.yaml start` to run with local Cadence. See below
  instructions for setting up local Cadence.


# Support

Join our Slack channel! [![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](http://iworkflow-slack.work)

You can also post in our [Discussion](https://github.com/indeedeng/iwf/discussions), or raise an issue.

# Contributing

Check out our [CONTRIBUTING](https://github.com/indeedeng/iwf/blob/main/CONTRIBUTING.md) page.


# Posts & Articles & Reference
* [Why I created Indeed Workflow Engine](https://medium.com/@qlong/a-letter-to-cadence-temporal-and-workflow-tech-community-b32e9fa97a0c)
* [Durable Timers in Indeed Workflow Engine](https://medium.com/@qlong/iwf-deep-dive-workflowstate-durable-timer-1-0bb89e6d6fd4?sk=6095e32b5ad677b6ea4f3c604751ece5)
* Cadence community spotlights
  * [#1](https://cadenceworkflow.io/blog/2023/01/31/community-spotlight-january-2023/)
  * [#2](https://cadenceworkflow.io/blog/2023/11/30/community-spotlight-update-november-2023/)
  * [#3](https://cadenceworkflow.io/blog/2023/08/31/community-spotlight-august-2023/)
* iWF is an abstracted Temporal [framework](https://github.com/temporalio/awesome-temporal). Same for [Cadence](https://github.com/uber/cadence#cadence).
* [How ContinueAsNew is built in iWF](https://medium.com/@qlong/guide-to-continueasnew-in-cadence-temporal-workflow-using-iwf-as-an-example-part-1-c24ae5266f07)
