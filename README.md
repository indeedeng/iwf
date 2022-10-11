# iWF project - main & server repo 

A simple and powerful WorkflowAsCode general purpose workflow engine 

Back by Cadence/Temporal as an interpreter.

Related projects:
* [iWF Java SDK](https://github.com/cadence-oss/iwf-java) and [Java SDK API preview](https://docs.google.com/document/d/15CETNk9ewiP7M_6N9s7jo-Wm57WG977hch9kTVnaExA)
* [iWF Java Samples](https://github.com/cadence-oss/iwf-java-samples)
* [API schema](https://github.com/cadence-oss/iwf-idl)

# What & Why
* See [Slide deck](https://docs.google.com/presentation/d/1CpsroSf6NeVce_XyUhFTkd9bLHN8UHRtM9NavPCMhj8/edit#slide=id.gfe2f455492_0_56) for what problems it is solving
* See [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit) for how it works  


## How to build & run
* Run `make bins` to build the binary `iwf-server`
* Then run  `./iwf-server start` to run the service . This defaults to serve workflows APIs with Temporal interpreter implementation. It requires to have local Temporal setup. See Run with local Temporal.
* Alternatively, run `./iwf-server --config config/development_cadence.yaml start` to run with local Cadence. See below instructions for setting up local Cadence. 
* Run `make integTests` to run all integration tests. This by default requires to have both local Cadence and Temporal to be set up.

## Development

### Update IDL and generated code
1. Install openapi-generator using Homebrew if you haven't. See more [documentation](https://openapi-generator.tech/docs/installation) 
2. Check out the idl submodule by running the command: `git submodule update --init --recursive`
3. Run the command `git submodule update --remote --merge` to update IDL to the latest commit
4. Run `make idl-code-gen` to refresh the generated code


### Run with local Temporalite
1. Run a local Temporalite following the [instruction](https://github.com/temporalio/temporalite). If you see error `error setting up schema`, try use command `temporalite start --namespace default -f my_test.db` instead to start. 
2. Go to http://localhost:8233/ for Temporal WebUI

NOTE: alternatively, go to [Temporal-dockercompose](https://github.com/temporalio/docker-compose) to run with docker

### Run with local Cadence
1. Run a local Cadence server following the [instructions](https://github.com/uber/cadence/tree/master/docker)
2. Register a new domain if not haven `cadence --do default domain register`
3. Go to Cadence http://localhost:8088/domains/default/workflows?range=last-30-days

# Development Plan
## 1.0
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

## 1.1
- [x] Reset workflow API
- [ ] StateOption: AttributeLoadingPolicy
- [ ] StateOption: CommandCarryOverPolicy
- [ ] StateOption: Start/Decide API timeout and retry
- [ ] LongRunningActivityCommand
- [ ] Decider trigger types: AnyCommandClosed
- [ ] Support IdReusePolicy when starting workflow

## 1.2
- [ ] AnyCommandCompleted Decider trigger type and WaitForMoreResults in StateDecision
- [ ] Skip timer API for testing/operation
- [ ] Failing workflow details
- [ ] Auto ContinueAsNew 
