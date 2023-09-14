# iWF project - main & server repo

[![Go Reference](https://pkg.go.dev/badge/github.com/indeedeng/iwf.svg)](https://pkg.go.dev/github.com/indeedeng/iwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/indeedeng/iwf)](https://goreportcard.com/report/github.com/indeedeng/iwf)
[![Coverage Status](https://codecov.io/github/indeedeng/iwf/coverage.svg?branch=main)](https://app.codecov.io/gh/indeedeng/iwf/branch/main)

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

**iWF will make you a 10x developer!**

iWF is an API orchestration platform for building resilient, fault-tolerant, scalable long-running applications. 
It offers an orchestration coding framework with abstractions for durable timers, async/background execution with backoff retry, 
KV storage,  RPC, and message queues. You will build long-running reliable processes faster than ever.

Related projects:

* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk) and [samples](https://github.com/indeedeng/iwf-java-samples) 
* [iWF Golang SDK](https://github.com/indeedeng/iwf-golang-sdk) and [samples](https://github.com/indeedeng/iwf-golang-samples)
* [iWF Python SDK](https://github.com/indeedeng/iwf-python-sdk) and [samples](https://github.com/indeedeng/iwf-python-samples)
* WIP [iWF TypeScript SDK](https://github.com/indeedeng/iwf-ts-sdk)

For support or any question, please post in our [Discussion](https://github.com/indeedeng/iwf/discussions), or raise an issue.
If you are interested in helping this project, check out our [CONTRIBUTING](https://github.com/indeedeng/iwf/blob/main/CONTRIBUTING.md) page.
Below is the basic and comprehensive documentation of iWF. There are some more details in the [wiki pages](https://github.com/indeedeng/iwf/wiki).

# What is [iWF](https://github.com/indeedeng/iwf/wiki)

## Use case study/examples
* [User sign-up/registry](https://github.com/indeedeng/iwf/wiki/Use-case-study-%E2%80%90%E2%80%90-user-signup-workflow)
* [Abstracted microservice orchestration]()
* Employer & JobSeeker engagement
  * [in Java](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/engagement)
  * [in Golang](https://github.com/indeedeng/iwf-golang-samples/blob/main/workflows/engagement)
* Subscription Workflow
  * [in Java](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/subscription)
  * [in Golang](https://github.com/indeedeng/iwf-golang-samples/blob/main/workflows/subscription)
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

See our [wiki pages](https://github.com/indeedeng/iwf/wiki/iWF-Application-Operations#troubleshoot--debugging).


## Operation

See our [wiki pages](https://github.com/indeedeng/iwf/wiki/iWF-Application-Operations).

# Posts & Articles & Reference

* [Cadence community spotlight](https://cadenceworkflow.io/blog/2023/01/31/community-spotlight-january-2023/)
* [A story of iWF](https://medium.com/@qlong/a-letter-to-cadence-temporal-and-workflow-tech-community-b32e9fa97a0c)
* iWF is an abstracted Temporal [framework](https://github.com/temporalio/awesome-temporal). Same for [Cadence](https://github.com/uber/cadence#cadence).