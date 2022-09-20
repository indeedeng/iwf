# iwf-server
interpreter workflow engine for Cadence/Temporal

# What & Why
See our
* [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit)  
* [Slide deck](https://docs.google.com/presentation/d/1CpsroSf6NeVce_XyUhFTkd9bLHN8UHRtM9NavPCMhj8/edit#slide=id.gfe2f455492_0_56)

## How to build & run
* Run `make bins` to build the binary `iwf-server`
* Then run  `./iwf-server start` to run the service . This defaults to test API + Temporal interpreter implementation, require to have local Temporal setup. See Run with local Temporal.
* Hit http://localhost:8801/ to trigger a test workflow

## Development

### Update IDL and generated code
1. Install openapi-generator using Homebrew if you haven't. See more [documentation](https://openapi-generator.tech/docs/installation) 
2. Check out the idl submodule by running the command: `git submodule update --init --recursive`
3. Run the command `git submodule update --remote --merge` to update IDL to the latest commit
4. Run `make idl-code-gen` to refresh the generated code

### Run with local Temporal
1. Run a local Temporalite following the [instruction](https://github.com/temporalio/temporalite). If you see error `error setting up schema`, try use command `temporalite start --namespace default -f my_test.db` instead to start. 

### Run with local Cadence
TODO

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

## 1.1
- [ ] Signal workflow API
- [ ] Query workflow API
- [ ] List/search workflow API
- [ ] StateOption: AttributeLoadingPolicy
- [ ] StateOption: CommandCarryOverPolicy
- [ ] StateOption: Start/Decide API timeout and retry
- [ ] LongRunningActivityCommand
- [ ] Decider trigger types: AnyCommandClosed
- [ ] Support IdReusePolicy when starting workflow
- [ ] Unit tests

## 1.2
- [ ] AnyCommandCompleted Decider trigger type and WaitForMoreResults in StateDecision
- [ ] Skip timer API for testing/operation
- [ ] Failing workflow details
- [ ] Auto ContinueAsNew 
