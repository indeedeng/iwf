package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/indeedeng/iwf/service/interpreter/versions"
)

const globalChangeId = "global"

// GlobalVersioner see https://stackoverflow.com/questions/73941723/what-is-a-good-way-pattern-to-use-temporal-cadence-versioning-api
type GlobalVersioner struct {
	workflowProvider  WorkflowProvider
	ctx               UnifiedContext
	version           int
	OmitVersionMarker bool // indicate the version marker and upsertSearchAttribute is already set at the start of the workflow
}

func NewGlobalVersioner(
	workflowProvider WorkflowProvider, omitVersionMarker bool, ctx UnifiedContext,
) (*GlobalVersioner, error) {
	sas, err := workflowProvider.GetSearchAttributes(ctx, []iwfidl.SearchAttributeKeyAndType{
		{Key: ptr.Any(service.SearchAttributeGlobalVersion),
			ValueType: ptr.Any(iwfidl.INT)},
	})
	if err != nil {
		return nil, err
	}
	version := 0
	if omitVersionMarker {
		// TODO: future improvement https://github.com/indeedeng/iwf/issues/369
		attribute, ok := sas[service.SearchAttributeGlobalVersion]
		if !ok {
			panic("search attribute global version is not found")
		}
		version = int(attribute.GetIntegerValue())
		if versions.MaxOfAllVersions < version {
			panic("requesting for a version that is not supported, panic to retry in next workflow task")
		}
	} else {
		version = workflowProvider.GetVersion(ctx, globalChangeId, 0, versions.MaxOfAllVersions)
	}

	return &GlobalVersioner{
		workflowProvider:  workflowProvider,
		ctx:               ctx,
		version:           version,
		OmitVersionMarker: omitVersionMarker,
	}, nil
}

func (p *GlobalVersioner) IsAfterVersionOfContinueAsNewOnNoStates() bool {
	return p.version >= versions.StartingVersionContinueAsNewOnNoStates
}

func (p *GlobalVersioner) IsAfterVersionOfUsingGlobalVersioning() bool {
	return p.version >= versions.StartingVersionUsingGlobalVersioning
}

func (p *GlobalVersioner) IsAfterVersionOfOptimizedUpsertSearchAttribute() bool {
	return p.version >= versions.StartingVersionOptimizedUpsertSearchAttribute
}

func (p *GlobalVersioner) IsAfterVersionOfRenamedStateApi() bool {
	return p.version >= versions.StartingVersionRenamedStateApi
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
			service.SearchAttributeGlobalVersion: versions.MaxOfAllVersions,
		})
	}
	return nil
}
