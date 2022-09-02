# iwf-server
interpreter workflow engine for Cadence/Temporal

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