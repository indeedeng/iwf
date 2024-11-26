package temporal

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/indeedeng/iwf/config"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/api/serviceerror"
	"testing"
)

func TestAlreadyStartedErrorForWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRealTemporalClient := NewMockClient(ctrl)
	mockDataConverter := NewMockDataConverter(ctrl)

	client := NewTemporalClient(mockRealTemporalClient, "test-ns", mockDataConverter, false, &config.QueryWorkflowFailedRetryPolicy{
		InitialIntervalSeconds: 0,
		MaximumAttempts:        0,
	})

	err := &serviceerror.WorkflowExecutionAlreadyStarted{}
	assert.Equal(t, true, client.IsWorkflowAlreadyStartedError(err))
}

func TestAlreadyStartedErrorForCronWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRealTemporalClient := NewMockClient(ctrl)
	mockDataConverter := NewMockDataConverter(ctrl)

	client := NewTemporalClient(mockRealTemporalClient, "test-ns", mockDataConverter, false, &config.QueryWorkflowFailedRetryPolicy{
		InitialIntervalSeconds: 0,
		MaximumAttempts:        0,
	})

	err := errors.New("schedule with this ID is already registered")

	assert.Equal(t, true, client.IsWorkflowAlreadyStartedError(err))
}
