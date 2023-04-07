package cadence

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/timeparser"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/encoded"
	"strings"
)

func getResetIDsByType(
	ctx context.Context,
	resetType iwfidl.WorkflowResetType,
	domain, wid, rid string,
	frontendClient workflowserviceclient.Interface, converter encoded.DataConverter,
	historyEventId int32, earliestHistoryTimeStr string, stateId, stateExecutionId string,
) (resetBaseRunID string, decisionFinishID int64, err error) {
	// default to the same runID
	resetBaseRunID = rid

	switch resetType {
	case iwfidl.HISTORY_EVENT_ID:
		decisionFinishID = int64(historyEventId)
		return
	case iwfidl.BEGINNING:
		decisionFinishID, err = getFirstDecisionTaskByType(ctx, domain, wid, rid, frontendClient, shared.EventTypeDecisionTaskCompleted)
		if err != nil {
			return
		}
	case iwfidl.HISTORY_EVENT_TIME:
		var earliestTimeUnixNano int64
		earliestTimeUnixNano, err = timeparser.ParseTime(earliestHistoryTimeStr)
		if err != nil {
			return
		}
		decisionFinishID, err = getEarliestDecisionID(ctx, domain, wid, rid, earliestTimeUnixNano, frontendClient)
		if err != nil {
			return
		}
	case iwfidl.STATE_ID, iwfidl.STATE_EXECUTION_ID:
		decisionFinishID, err = getDecisionEventIDByStateOrStateExecutionId(ctx, domain, wid, rid, stateId, stateExecutionId, frontendClient, converter)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("not supported resetType")
	}
	return
}

func getFirstDecisionTaskByType(
	ctx context.Context,
	domain string,
	workflowID string,
	runID string,
	frontendClient workflowserviceclient.Interface,
	decisionType shared.EventType,
) (decisionFinishID int64, err error) {

	req := &shared.GetWorkflowExecutionHistoryRequest{
		Domain: &domain,
		Execution: &shared.WorkflowExecution{
			WorkflowId: &workflowID,
			RunId:      &runID,
		},
		MaximumPageSize: iwfidl.PtrInt32(1000),
		NextPageToken:   nil,
	}

	for {
		var resp *shared.GetWorkflowExecutionHistoryResponse
		resp, err = frontendClient.GetWorkflowExecutionHistory(ctx, req)
		if err != nil {
			return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
		}

		for _, e := range resp.GetHistory().GetEvents() {
			if e.GetEventType() == decisionType {
				decisionFinishID = e.GetEventId()
				return decisionFinishID, nil
			}
		}

		if len(resp.NextPageToken) != 0 {
			req.NextPageToken = resp.NextPageToken
		} else {
			break
		}
	}
	if decisionFinishID == 0 {
		return 0, composeErrorWithMessage("Get historyEventId failed", fmt.Errorf("no historyEventId"))
	}
	return
}

func getEarliestDecisionID(
	ctx context.Context,
	domain string, wid string,
	rid string, earliestTime int64,
	frontendClient workflowserviceclient.Interface,
) (decisionFinishID int64, err error) {
	req := &shared.GetWorkflowExecutionHistoryRequest{
		Domain: &domain,
		Execution: &shared.WorkflowExecution{
			WorkflowId: &wid,
			RunId:      &rid,
		},
		MaximumPageSize: iwfidl.PtrInt32(1000),
		NextPageToken:   nil,
	}

OuterLoop:
	for {
		var resp *shared.GetWorkflowExecutionHistoryResponse
		resp, err = frontendClient.GetWorkflowExecutionHistory(ctx, req)
		if err != nil {
			return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
		}
		for _, e := range resp.GetHistory().GetEvents() {
			if e.GetEventType() == shared.EventTypeDecisionTaskCompleted {
				if e.GetTimestamp() >= earliestTime {
					decisionFinishID = e.GetEventId()
					break OuterLoop
				}
			}
		}
		if len(resp.NextPageToken) != 0 {
			req.NextPageToken = resp.NextPageToken
		} else {
			break
		}
	}
	if decisionFinishID == 0 {
		return 0, composeErrorWithMessage("Get historyEventId failed", fmt.Errorf("no historyEventId"))
	}
	return
}

func getDecisionEventIDByStateOrStateExecutionId(
	ctx context.Context,
	domain string, wid string,
	rid string, stateId, stateExecutionId string,
	frontendClient workflowserviceclient.Interface,
	converter encoded.DataConverter,
) (decisionFinishID int64, err error) {
	req := &shared.GetWorkflowExecutionHistoryRequest{
		Domain: &domain,
		Execution: &shared.WorkflowExecution{
			WorkflowId: &wid,
			RunId:      &rid,
		},
		MaximumPageSize: iwfidl.PtrInt32(1000),
		NextPageToken:   nil,
	}

	for {
		var resp *shared.GetWorkflowExecutionHistoryResponse
		resp, err = frontendClient.GetWorkflowExecutionHistory(ctx, req)
		if err != nil {
			return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
		}
		for _, e := range resp.GetHistory().GetEvents() {
			if e.GetEventType() == shared.EventTypeDecisionTaskCompleted {
				decisionFinishID = e.GetEventId()
			}
			if e.GetEventType() == shared.EventTypeActivityTaskScheduled {
				typeName := e.GetActivityTaskScheduledEventAttributes().GetActivityType().GetName()
				if strings.Contains(typeName, "StateStart") || strings.Contains(typeName, "StateApiWaitUntil") {
					var backendType service.BackendType
					var input service.StateStartActivityInput
					err = converter.FromData(e.GetActivityTaskScheduledEventAttributes().Input, &backendType, &input)
					if err != nil {
						return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
					}
					if input.Request.WorkflowStateId == stateId || input.Request.Context.StateExecutionId == stateExecutionId {
						if decisionFinishID == 0 {
							return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", fmt.Errorf("invalid history or something goes very wrong"))
						}
						return
					}
				}
			}
		}
		if len(resp.NextPageToken) != 0 {
			req.NextPageToken = resp.NextPageToken
		} else {
			break
		}
	}
	return 0, composeErrorWithMessage("Get historyEventId failed", fmt.Errorf("no historyEventId"))
}

func composeErrorWithMessage(msg string, err error) error {
	err = fmt.Errorf("%v, %v", msg, err)
	return err
}
