### Replay Tests

### Why
Replay tests are special and unique in Cadence/Temporal programming model.
It's ensuring the [determinism](https://docs.temporal.io/workflows#deterministic-constraints) is not being broken by the changes in the workflow logic.

See more about replay tests in the [Cadence documentation](https://cadenceworkflow.io/docs/go-client/workflow-replay-shadowing/#workflow-replayer)
and Temporal documentation [here](https://docs.temporal.io/develop/go/testing-suite#replay).

To simplify the work, we only use Temporal replay tests in iWF.

### Global versioning design pattern
In iWF, we are using the [global versioning design pattern](https://medium.com/@qlong/how-to-overcome-some-maintenance-challenges-of-temporal-cadence-workflow-versioning-f893815dd18d) 
to ensure the determinism of the workflow.
The pattern makes it simple to manage the workflow versioning and replay tests.

* For every new [global version](../service/interpreter/versions/versions.go), we add at least a new [history file](./history) in the replay_test.
* For each version, we may need to have multiple history files to cover different scenarios(code paths).
* To get the JSON history file, start and run a workflow that will use the code path that you want to protect the determinism. Then download the JSON from WebUI.
* Usually, the workflow is from an integration test 