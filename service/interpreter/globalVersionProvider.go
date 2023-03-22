package interpreter

import "github.com/indeedeng/iwf/service"

const globalChangeId = "global"
const startingVersionUsingGlobalVersioning = 1
const startingVersionUsingMemoForDataObjects = 2
const MaxOfAllVersions = startingVersionUsingMemoForDataObjects

// see https://stackoverflow.com/questions/73941723/what-is-a-good-way-pattern-to-use-temporal-cadence-versioning-api
type globalVersioner struct {
	workflowProvider WorkflowProvider
}

func NewGlobalVersionProvider(workflowProvider WorkflowProvider) *globalVersioner {
	return &globalVersioner{
		workflowProvider: workflowProvider,
	}
}

func (p *globalVersioner) IsAfterVersionOfUsingMemoForDataObjects(ctx UnifiedContext) bool {
	version := p.workflowProvider.GetVersion(ctx, globalChangeId, 0, MaxOfAllVersions)
	return version >= startingVersionUsingMemoForDataObjects
}

func (p *globalVersioner) IsAfterVersionOfUsingGlobalVersioning(ctx UnifiedContext) bool {
	version := p.workflowProvider.GetVersion(ctx, globalChangeId, 0, MaxOfAllVersions)
	return version >= startingVersionUsingGlobalVersioning
}

func (p *globalVersioner) UpsertGlobalVersionSearchAttribute(ctx UnifiedContext) error {
	// TODO this bug in Cadence SDK may cause concurrent writes
	// https://github.com/uber-go/cadence-client/issues/1198
	if p.workflowProvider.GetBackendType() != service.BackendTypeCadence {
		return p.workflowProvider.UpsertSearchAttributes(ctx, map[string]interface{}{
			service.SearchAttributeGlobalVersion: MaxOfAllVersions,
		})
	}
	return nil
}
