package interpreter

import "github.com/indeedeng/iwf/service"

const globalChangeId = "global"
const startingVersionUsingGlobalVersioning = 1
const startingVersionOptimizedUpsertSearchAttribute = 2
const startingVersionRenamedStateApi = 3
const maxOfAllVersions = startingVersionRenamedStateApi

// GlobalVersioner see https://stackoverflow.com/questions/73941723/what-is-a-good-way-pattern-to-use-temporal-cadence-versioning-api
type GlobalVersioner struct {
	workflowProvider WorkflowProvider
	ctx              UnifiedContext
}

func NewGlobalVersioner(workflowProvider WorkflowProvider, ctx UnifiedContext) *GlobalVersioner {
	return &GlobalVersioner{
		workflowProvider: workflowProvider,
		ctx:              ctx,
	}
}

func (p *GlobalVersioner) IsAfterVersionOfUsingGlobalVersioning() bool {
	version := p.workflowProvider.GetVersion(p.ctx, globalChangeId, 0, maxOfAllVersions)
	return version >= startingVersionUsingGlobalVersioning
}

func (p *GlobalVersioner) IsAfterVersionOfOptimizedUpsertSearchAttribute() bool {
	version := p.workflowProvider.GetVersion(p.ctx, globalChangeId, 0, maxOfAllVersions)
	return version >= startingVersionOptimizedUpsertSearchAttribute
}

func (p *GlobalVersioner) IsAfterVersionOfRenamedStateApi() bool {
	version := p.workflowProvider.GetVersion(p.ctx, globalChangeId, 0, maxOfAllVersions)
	return version >= startingVersionRenamedStateApi
}

func (p *GlobalVersioner) UpsertGlobalVersionSearchAttribute() error {
	// TODO this bug in Cadence SDK may cause concurrent writes
	// https://github.com/uber-go/cadence-client/issues/1198
	if p.workflowProvider.GetBackendType() != service.BackendTypeCadence {
		return p.workflowProvider.UpsertSearchAttributes(p.ctx, map[string]interface{}{
			service.SearchAttributeGlobalVersion: maxOfAllVersions,
		})
	}
	return nil
}
