package interpreter

import (
	"github.com/indeedeng/iwf/service"
)

const globalChangeId = "global"

const StartingVersionUsingGlobalVersioning = 1
const StartingVersionOptimizedUpsertSearchAttribute = 2
const StartingVersionRenamedStateApi = 3
const StartingVersionContinueAsNewOnNoStates = 4
const StartingVersionTemporal26SDK = 5
const StartingVersionExecutingStateIdMode = 6
const StartingVersionNoIwfGlobalVersionSearchAttribute = 7
const MaxOfAllVersions = StartingVersionNoIwfGlobalVersionSearchAttribute

// GlobalVersioner see https://stackoverflow.com/questions/73941723/what-is-a-good-way-pattern-to-use-temporal-cadence-versioning-api
type GlobalVersioner struct {
	workflowProvider  WorkflowProvider
	ctx               UnifiedContext
	version           int
	OmitVersionMarker bool // indicate the version marker and upsertSearchAttribute is already set at the start of the workflow
}

func NewGlobalVersioner(
	workflowProvider WorkflowProvider, ctx UnifiedContext,
) (*GlobalVersioner, error) {
	version := workflowProvider.GetVersion(ctx, globalChangeId, 0, MaxOfAllVersions)

	return &GlobalVersioner{
		workflowProvider: workflowProvider,
		ctx:              ctx,
		version:          version,
	}, nil
}

// methods checking version number

func (p *GlobalVersioner) IsAfterVersionOfContinueAsNewOnNoStates() bool {
	return p.version >= StartingVersionContinueAsNewOnNoStates
}

func (p *GlobalVersioner) IsAfterVersionOfUsingGlobalVersioning() bool {
	return p.version >= StartingVersionUsingGlobalVersioning
}

func (p *GlobalVersioner) IsAfterVersionOfOptimizedUpsertSearchAttribute() bool {
	return p.version >= StartingVersionOptimizedUpsertSearchAttribute
}

func (p *GlobalVersioner) IsAfterVersionOfExecutingStateIdMode() bool {
	return p.version >= StartingVersionExecutingStateIdMode
}

func (p *GlobalVersioner) IsAfterVersionOfRenamedStateApi() bool {
	return p.version >= StartingVersionRenamedStateApi
}

func (p *GlobalVersioner) IsAfterVersionOfTemporal26SDK() bool {
	return p.version >= StartingVersionTemporal26SDK
}

func (p *GlobalVersioner) IsAfterVersionOfNoIwfGlobalVersionSearchAttribute() bool {
	return p.version >= StartingVersionNoIwfGlobalVersionSearchAttribute
}

// methods checking feature/functionality availability

func (p *GlobalVersioner) IsUsingGlobalVersionSearchAttribute() bool {
	return p.version >= StartingVersionUsingGlobalVersioning && p.version < StartingVersionNoIwfGlobalVersionSearchAttribute
}

func (p *GlobalVersioner) UpsertGlobalVersionSearchAttribute() error {
	if p.OmitVersionMarker {
		// the search attribute is already set when starting the workflow
		return nil
	}
	// TODO this bug in Cadence SDK may cause concurrent writes
	// https://github.com/uber-go/cadence-client/issues/1198
	if p.workflowProvider.GetBackendType() != service.BackendTypeCadence {
		return p.workflowProvider.UpsertSearchAttributes(p.ctx, map[string]interface{}{
			service.SearchAttributeGlobalVersion: MaxOfAllVersions,
		})
	}
	return nil
}
