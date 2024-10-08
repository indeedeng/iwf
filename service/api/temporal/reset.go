package temporal

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/timeparser"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/converter"
	"strings"
)

func getResetEventIDByType(ctx context.Context, resetType iwfidl.WorkflowResetType,
	namespace, wid, rid string,
	frontendClient workflowservice.WorkflowServiceClient, converter converter.DataConverter,
	historyEventId int32, earliestHistoryTimeStr string, stateId, stateExecutionId string,
) (resetBaseRunID string, workflowTaskFinishID int64, err error) {
	// default to the same runID
	resetBaseRunID = rid

	switch resetType {
	case iwfidl.HISTORY_EVENT_ID:
		workflowTaskFinishID = int64(historyEventId)
		return
	case iwfidl.HISTORY_EVENT_TIME:
		var earliestTimeUnixNano int64
		earliestTimeUnixNano, err = timeparser.ParseTime(earliestHistoryTimeStr)
		if err != nil {
			return
		}
		workflowTaskFinishID, err = getEarliestDecisionEventID(ctx, namespace, wid, rid, earliestTimeUnixNano, frontendClient)
		if err != nil {
			return
		}
	case iwfidl.BEGINNING:
		resetBaseRunID, workflowTaskFinishID, err = getFirstWorkflowTaskEventID(ctx, namespace, wid, rid, frontendClient)
		if err != nil {
			return
		}
	case iwfidl.STATE_ID, iwfidl.STATE_EXECUTION_ID:
		workflowTaskFinishID, err = getDecisionEventIDByStateOrStateExecutionId(ctx, namespace, wid, rid, stateId, stateExecutionId, frontendClient, converter)
		if err != nil {
			return
		}
	default:
		panic("not supported resetType")
	}
	return
}

func getFirstWorkflowTaskEventID(ctx context.Context, namespace, wid, rid string, frontendClient workflowservice.WorkflowServiceClient) (resetBaseRunID string, workflowTaskEventID int64, err error) {
	resetBaseRunID = rid
	req := &workflowservice.GetWorkflowExecutionHistoryRequest{
		Namespace: namespace,
		Execution: &common.WorkflowExecution{
			WorkflowId: wid,
			RunId:      rid,
		},
		MaximumPageSize: 1000,
		NextPageToken:   nil,
	}
	for {
		var resp *workflowservice.GetWorkflowExecutionHistoryResponse
		resp, err = frontendClient.GetWorkflowExecutionHistory(ctx, req)
		if err != nil {
			return
		}
		for _, e := range resp.GetHistory().GetEvents() {
			if e.GetEventType() == enums.EVENT_TYPE_WORKFLOW_TASK_COMPLETED {
				workflowTaskEventID = e.GetEventId()
				return resetBaseRunID, workflowTaskEventID, nil
			}
			if e.GetEventType() == enums.EVENT_TYPE_WORKFLOW_TASK_SCHEDULED {
				if workflowTaskEventID == 0 {
					workflowTaskEventID = e.GetEventId() + 1
				}
			}
		}
		if len(resp.NextPageToken) != 0 {
			req.NextPageToken = resp.NextPageToken
		} else {
			break
		}
	}
	if workflowTaskEventID == 0 {
		err = fmt.Errorf("unable to find any scheduled or completed task")
		return
	}
	return
}

func getEarliestDecisionEventID(
	ctx context.Context,
	namespace string, wid string,
	rid string, earliestTime int64,
	frontendClient workflowservice.WorkflowServiceClient,
) (decisionFinishID int64, err error) {
	req := &workflowservice.GetWorkflowExecutionHistoryRequest{
		Namespace: namespace,
		Execution: &common.WorkflowExecution{
			WorkflowId: wid,
			RunId:      rid,
		},
		MaximumPageSize: 1000,
		NextPageToken:   nil,
	}

OuterLoop:
	for {
		var resp *workflowservice.GetWorkflowExecutionHistoryResponse
		resp, err = frontendClient.GetWorkflowExecutionHistory(ctx, req)
		if err != nil {
			return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
		}
		for _, e := range resp.GetHistory().GetEvents() {
			if e.GetEventType() == enums.EVENT_TYPE_WORKFLOW_TASK_COMPLETED {
				if e.GetEventTime().GetSeconds() >= earliestTime {
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
	namespace string, wid string,
	rid string, stateId, stateExecutionId string,
	frontendClient workflowservice.WorkflowServiceClient, converter converter.DataConverter,
) (decisionFinishID int64, err error) {
	req := &workflowservice.GetWorkflowExecutionHistoryRequest{
		Namespace: namespace,
		Execution: &common.WorkflowExecution{
			WorkflowId: wid,
			RunId:      rid,
		},
		MaximumPageSize: 1000,
		NextPageToken:   nil,
	}

	for {
		var resp *workflowservice.GetWorkflowExecutionHistoryResponse
		resp, err = frontendClient.GetWorkflowExecutionHistory(ctx, req)
		if err != nil {
			return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
		}
		for _, e := range resp.GetHistory().GetEvents() {
			if e.GetEventType() == enums.EVENT_TYPE_WORKFLOW_TASK_COMPLETED {
				decisionFinishID = e.GetEventId()
			}
			if e.GetEventType() == enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED {
				typeName := e.GetActivityTaskScheduledEventAttributes().GetActivityType().GetName()
				if strings.Contains(typeName, "StateStart") || strings.Contains(typeName, "StateApiWaitUntil") {
					var backendType service.BackendType
					var input service.StateStartActivityInput
					err = converter.FromPayloads(e.GetActivityTaskScheduledEventAttributes().Input, &backendType, &input)
					if err != nil {
						return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
					}
					if input.Request.WorkflowStateId == stateId || input.Request.Context.GetStateExecutionId() == stateExecutionId {
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
