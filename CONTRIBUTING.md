# Develop iWF server

Any contribution is welcome. Even just a fix for typo in a code comment, or README/wiki.

See [Design doc](https://docs.google.com/document/d/1BpJuHf67ibaOWmN_uWw_pbrBVyb6U1PILXyzohxA5Ms/edit) for how it works.

Here is the repository layout if you are interested to learn about it:

* `cmd/` the code to bootstrap the server -- loading config and connect to Cadence/Temporal service, and start iWF API and interpreter service
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
        * `*.go` the implementation of interpreter workflow service using `ActivityProvider` and `WorkflowProvider` so that it works for both Cadence and Temporal
            * `workflowImpl.go` the core workflow implementation
    * `common/` some common libraries between `api` and `interpreter`
    * `*.go` some common definitions between `api` and `interpreter`

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
4 For `persistence_test.go` integTests, you need to register search attributes into Temporalite. Unlike Temporal docker, those Search Attributes are not provided by default:
```shell
tctl adm cl asa -n CustomKeywordField -t Keyword
tctl adm cl asa -n CustomIntField -t Int
tctl adm cl asa -n CustomBoolField -t Bool
tctl adm cl asa -n CustomDoubleField -t Double
tctl adm cl asa -n CustomDatetimeField -t Datetime
tctl adm cl asa -n CustomStringField -t text
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
After registering, it may take [up 60s](https://github.com/uber/cadence/blob/d618e32ac5ea05c411cca08c3e4859e800daa1e0/docker/config_template.yaml#L286) for Cadence to load the new search attributes. If you run the test too early, you may see error:
 like `"IwfWorkflowType is not a valid search attribute key"`
4. Go to Cadence http://localhost:8088/domains/default/workflows?range=last-30-days

If you run into any issues with Search Attributes registration, use the below command to check the existing Search attributes:

* Cadence: `cadence cl get-search-attr`
* Temporal: `tctl adm cl get-search-attributes`
