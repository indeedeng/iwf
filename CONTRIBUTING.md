# Develop iWF server

Any contribution is welcome. Even just a fix for typo in a code comment, or README/wiki.

See [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit) for how it works.

Here is the repository layout if you are interested to learn about it:

* `cmd/` the code to bootstrap the server -- loading config and connect to Cadence/Temporal service, and start iWF API
  and interpreter service
* `config/` the config to start the server, and also config template to start the Docker image
* `docker-compose/` the docker compose file to start a full iWF server with Temporal dependency
* `gen/` the generated code from iwf-idl (Open API definition/Swagger)
* `integ/` the end to end integration tests.
    * `workflow/` the iWF workflows that are written without SDK(just implemented the REST APIs)
    * `*.go` the tests
* `iwf-idl/` the idl submodule
* `script/` some scripts
    * `http/` some example HTTP scripts to call server, like REST API
    * `start-server.sh` the script to start iWF server in Docker image
* `service/` iWF implementation
    * `api/` API service implementation
        * `cadence/` the Cadence abstraction of `UnifiedClient`
        * `temporal/` the Temporal abstraction of `UnifiedClient`
        * `*.go` the implementation of API service using `UnifiedClient` so that it works for both Cadence and Temporal
    * `interpreter/` interpreter worker service implementation
        * `cadence/` the Cadence abstraction of `ActivityProvider` and `WorkflowProvider`
        * `temporal/` the Temporal abstraction of `ActivityProvider` and `WorkflowProvider`
        * `*.go` the implementation of interpreter workflow service using `ActivityProvider` and `WorkflowProvider` so
          that it works for both Cadence and Temporal
            * `workflowImpl.go` the core workflow implementation
    * `common/` some common libraries between `api` and `interpreter`
    * `*.go` some common definitions between `api` and `interpreter`

## How to update IDL and the generated code

1. Install openapi-generator using Homebrew if you haven't. See
   more [documentation](https://openapi-generator.tech/docs/installation)
2. Check out the idl submodule by running the command: `git submodule update --init --recursive`
3. Run the command `git submodule update --remote --merge` to update IDL to the latest commit
4. Run `make idl-code-gen` to refresh the generated code. The command requires to have `openapi-generator` CLI.See
   the [openapi-generator doc](https://openapi-generator.tech/docs/installation/) for how to install it. And you may
   also need to upgrade it to the latest if it's older than what we are currently using.

An easy way to install openapi-generator CLI is to use Homebrew:

```
brew install openapi-generator

```

And to upgrade it:

```
brew update && brew upgrade openapi-generator

```

# How to run server or integration test

## Prepare Cadence/Temporal environment
iWF server depends on Cadence or Temporal. You need at least one to be ready for running with iWF .
Or maybe both just for testing to ensure the code works for both Cadence and Temporal. 

### Option 1: Run with our docker-compose file (Recommended)

Simply run `docker compose -f docker-compose/integ-dependencies.yml up -` will:

* Start both Cadence & Temporal as dependencies
* Set up required system search attributes
* Set up customized search attributes for integration test(`persistence_test.go`)
* Temporal WebUI:  http://localhost:8233/
* Cadence WebUI:  http://localhost:8088/

### Option 2: Run with your own Temporal service

First of all, you need a Temporal service if you haven't had it:

Option 1 (recommended): use [Temporal CLI](https://github.com/temporalio/cli) -- `temporal server start-dev`

Option 2: use [temporal docker-compose](https://github.com/temporalio/docker-compose)


Assuming you are using `default` namespace:

1. Make sure you have registered system search attributes required by iWF server

```shell
  temporal  operator search-attribute  create -name IwfWorkflowType -type Keyword
  temporal  operator search-attribute  create -name IwfGlobalWorkflowVersion -type Int 
  temporal  operator search-attribute  create -name IwfExecutingStateIds -type KeywordList 
  temporal  operator search-attribute  create -name IwfFirstRunId -type Keyword
```

2. For `persistence_test.go` integTests, you need to register below custom search attributes.

```shell
  temporal  operator search-attribute  create -name CustomKeywordField -type Keyword
  temporal  operator search-attribute  create -name CustomIntField -type Int
  temporal  operator search-attribute  create -name CustomBoolField -type Bool
  temporal  operator search-attribute  create -name CustomDoubleField -type Double
  temporal  operator search-attribute  create -name CustomDatetimeField -type Datetime
  temporal  operator search-attribute  create -name CustomStringField -type Text
```

3. If you run into any issues with Search Attributes registration, use the below command to check the existing Search
   attributes:`temporal operator search-attribute list`

### Option 3: Run with your own Cadence service

1. You can run a local Cadence server following the [instructions](https://github.com/uber/cadence/tree/master/docker)

```
docker-compose -f docker-compose-es-v7.yml up
```

2. Register a new domain if not haven `cadence --do default domain register`
3. Register system search attributes required by iWF server

```
cadence adm cl asa --search_attr_key IwfGlobalWorkflowVersion --search_attr_type 2
cadence adm cl asa --search_attr_key IwfExecutingStateIds --search_attr_type 1
cadence adm cl asa --search_attr_key IwfWorkflowType --search_attr_type 1
cadence adm cl asa --search_attr_key IwfFirstRunId --search_attr_type 1
```

After registering, it may
take [up 60s](https://github.com/uber/cadence/blob/d618e32ac5ea05c411cca08c3e4859e800daa1e0/docker/config_template.yaml#L286)
because of this [issue](https://github.com/uber/cadence/issues/5076). for Cadence to load the new search attributes. If
you run the test too early, you may see error:  `"IwfWorkflowType is not a valid search attribute key"`.

4. For Cadence docker compose, go to Cadence http://localhost:8088/domains/default/workflows?range=last-30-days

5. If not running by Cadence docker-compose, you must register those custom search attributes yourself.
   `CustomKeywordField, CustomIntField, CustomBoolField, CustomBoolField, CustomDoubleField, CustomDatetimeField, CustomStringField`

6. If you run into any issues with Search Attributes registration, use the below command to check the existing Search
   attributes:  `cadence cl get-search-attr`

## Run the server

The first step you may want to explore is to run it locally!

To run the server with Temporal
* If you are in an IDE, you can run the main function in `./cmd/main.go` with argument `start`.
* Or in terminal `go run cmd/server/main.go start`
* Or build the binary and run it by `make bins` and then run `./iwf-server start`

To run with Cadence, make sure you specify the cadence config `--config config/development_cadence.yaml start`:
* In an IDE, you can run the main function in `./cmd/main.go` with argument ` --config config/development_cadence.yaml start`.
* Or in terminal `go run cmd/server/main.go --config config/development_cadence.yaml start`
* Or build the binary and run it by`make bins` and then run `./iwf-server --config config/development_cadence.yaml start`

## Run the integration tests
For development, you may want to run the test locally for debugging, especially your PR has failed the tests in CI pipeline.

:warning: NOTE: When running with local Cadence, you may need to wait for up to 60s for Search attributes to be ready, because of
this [issue](https://github.com/uber/cadence/issues/5076).

* To run the whole integ test suite against Cadence+Temporal service by this command `make integTests`
* To run the whole suite for Temporal only `make temporalIntegTests` 
* To run the whole suite for Cadence only `make cadenceIntegTests`
* To run a specify test case or a test file, you can utilize the IDE or `go test` command.

To debug the failed test, search for `--- FAIL` in the output logs (in GitHub Action, click "view raw logs"") 
