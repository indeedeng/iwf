# iwf-server
interpreter workflow engine for Cadence/Temporal

## How to build & run
* Run `make bins` to build the binary `iwf-server`
* Then run  `./iwf-server start` to run the service

## Development

### Update IDL and generated code
1. Install openapi-generator using Homebrew if you haven't. See more [documentation](https://openapi-generator.tech/docs/installation) 
2. Check out the idl submodule by running the command: `git submodule update --init --recursive`
3. Run the command `git submodule update --remote --merge` to update IDL to the latest commit
4. Run `make idl-code-gen` to refresh the generated code 