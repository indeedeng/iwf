package interpreter

import "github.com/indeedeng/iwf/service"

const globalChangeId = "global"
const startingVersionUsingGlobalVersioning = 1
const startingVersionOptimizedUpsertSearchAttribute = 2
const maxOfAllVersions = startingVersionOptimizedUpsertSearchAttribute

// see https://stackoverflow.com/questions/73941723/what-is-a-good-way-pattern-to-use-temporal-cadence-versioning-api
type globalVersioner struct {
	workflowProvider WorkflowProvider
	ctx              UnifiedContext
}

func NewGlobalVersionProvider(workflowProvider WorkflowProvider, ctx UnifiedContext) *globalVersioner {
	return &globalVersioner{
		workflowProvider: workflowProvider,
		ctx:              ctx,
	}
}

func (p *globalVersioner) IsAfterVersionOfUsingGlobalVersioning() bool {
	version := p.workflowProvider.GetVersion(p.ctx, globalChangeId, 0, maxOfAllVersions)
	return version >= startingVersionUsingGlobalVersioning
}

func (p *globalVersioner) IsAfterVersionOfOptimizedUpsertSearchAttribute() bool {
	version := p.workflowProvider.GetVersion(p.ctx, globalChangeId, 0, maxOfAllVersions)
	return version >= startingVersionOptimizedUpsertSearchAttribute
}

func (p *globalVersioner) UpsertGlobalVersionSearchAttribute() error {
	// TODO this bug in Cadence SDK may cause concurrent writes
	// https://github.com/uber-go/cadence-client/issues/1198
	if p.workflowProvider.GetBackendType() != service.BackendTypeCadence {
		return p.workflowProvider.UpsertSearchAttributes(p.ctx, map[string]interface{}{
			service.SearchAttributeGlobalVersion: maxOfAllVersions,
		})
	}
	return nil
}
