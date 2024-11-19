package temporal

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/api/serviceerror"
	"testing"
)

func TestAlreadyStartedErrorForWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRealTemporalClient := NewMockClient(ctrl)
	mockDataConverter := NewMockDataConverter(ctrl)

	client := NewTemporalClient(mockRealTemporalClient, "test-ns", mockDataConverter, false)

	err := &serviceerror.WorkflowExecutionAlreadyStarted{}
	assert.Equal(t, true, client.IsWorkflowAlreadyStartedError(err))
}

func TestAlreadyStartedErrorForCronWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRealTemporalClient := NewMockClient(ctrl)
	mockDataConverter := NewMockDataConverter(ctrl)

	client := NewTemporalClient(mockRealTemporalClient, "test-ns", mockDataConverter, false)

	err := errors.New("schedule with this ID is already registered")

	assert.Equal(t, true, client.IsWorkflowAlreadyStartedError(err))
}
