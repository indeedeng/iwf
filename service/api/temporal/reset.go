package temporal

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/timeparser"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
)

func getResetEventIDByType(ctx context.Context, resetType service.ResetType,
	namespace, wid, rid string,
	frontendClient workflowservice.WorkflowServiceClient,
	historyEventId int32, earliestHistoryTimeStr string,
) (resetBaseRunID string, workflowTaskFinishID int64, err error) {
	// default to the same runID
	resetBaseRunID = rid

	switch resetType {
	case service.ResetTypeHistoryEventId:
		workflowTaskFinishID = int64(historyEventId)
		return
	case service.ResetTypeHistoryEventTime:
		var earliestTimeUnixNano int64
		earliestTimeUnixNano, err = timeparser.ParseTime(earliestHistoryTimeStr)
		if err != nil {
			return
		}
		workflowTaskFinishID, err = getEarliestDecisionEventID(ctx, namespace, wid, rid, earliestTimeUnixNano, frontendClient)
		if err != nil {
			return
		}
	case service.ResetTypeBeginning:
		resetBaseRunID, workflowTaskFinishID, err = getFirstWorkflowTaskEventID(ctx, namespace, wid, rid, frontendClient)
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
		resp, err := frontendClient.GetWorkflowExecutionHistory(ctx, req)
		if err != nil {
			return 0, composeErrorWithMessage("GetWorkflowExecutionHistory failed", err)
		}
		for _, e := range resp.GetHistory().GetEvents() {
			if e.GetEventType() == enums.EVENT_TYPE_WORKFLOW_TASK_COMPLETED {
				if e.GetEventTime().UnixNano() >= earliestTime {
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
		return 0, composeErrorWithMessage("Get DecisionFinishID failed", fmt.Errorf("no DecisionFinishID"))
	}
	return
}

func composeErrorWithMessage(msg string, err error) error {
	err = fmt.Errorf("%v, %v", msg, err)
	return err
}
