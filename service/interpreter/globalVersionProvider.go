package interpreter

import "github.com/indeedeng/iwf/service"

const globalChangeId = "global"
const startingVersionUsingGlobalVersioning = 1
const maxOfAllVersions = startingVersionUsingGlobalVersioning

// see https://stackoverflow.com/questions/73941723/what-is-a-good-way-pattern-to-use-temporal-cadence-versioning-api
type globalVersioner struct {
	workflowProvider WorkflowProvider
}

func newGlobalVersionProvider(workflowProvider WorkflowProvider) *globalVersioner {
	return &globalVersioner{
		workflowProvider: workflowProvider,
	}
}

func (p *globalVersioner) isAfterVersionOfUsingGlobalVersioning(ctx UnifiedContext) bool {
	version := p.workflowProvider.GetVersion(ctx, globalChangeId, 0, maxOfAllVersions)
	return version >= startingVersionUsingGlobalVersioning
}

func (p *globalVersioner) upsertGlobalVersionSearchAttribute(ctx UnifiedContext) error {
	return p.workflowProvider.UpsertSearchAttributes(ctx, map[string]interface{}{
		service.SearchAttributeGlobalVersion: maxOfAllVersions,
	})
}
