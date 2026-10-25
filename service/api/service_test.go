package api

import (
	"context"
	"testing"

	"github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/blobstore"
	"github.com/indeedeng/iwf/service/common/log"
)

type recordingBlobStore struct {
	blobstore.BlobStore
	written string
}

func (s *recordingBlobStore) WriteObject(_ context.Context, _, data string) (string, string, error) {
	s.written = data
	return "store", "path", nil
}

type recordingClient struct {
	uclient.UnifiedClient
	input   service.InterpreterWorkflowInput
	started bool
}

func (c *recordingClient) StartInterpreterWorkflow(_ context.Context, _ uclient.StartWorkflowOptions, args ...interface{}) (string, error) {
	c.input = args[0].(service.InterpreterWorkflowInput)
	c.started = true
	return "run", nil
}

func TestStartWorkflowWithOptionalStateInputEncoding(t *testing.T) {
	store := &recordingBlobStore{}
	client := &recordingClient{}
	svc, err := NewApiService(config.Config{ExternalStorage: config.ExternalStorageConfig{
		Enabled: true, ThresholdInBytes: 5,
	}}, client, "task-queue", log.NewNoop(), store)
	if err != nil {
		t.Fatal(err)
	}

	data := "0123456789"
	input := &iwfidl.EncodedObject{Data: &data}
	response, apiErr := svc.ApiV1WorkflowStartPost(context.Background(), iwfidl.WorkflowStartRequest{
		WorkflowId:      "workflow-id",
		IwfWorkflowType: "test.Workflow",
		IwfWorkerUrl:    "http://localhost:8080",
		StateInput:      input,
	})
	if apiErr != nil {
		t.Fatalf("start workflow: %+v", apiErr)
	}
	if response.GetWorkflowRunId() != "run" || !client.started {
		t.Fatalf("workflow was not started: response=%+v", response)
	}
	if store.written != data {
		t.Fatalf("stored data = %q, want %q", store.written, data)
	}
	got := client.input.StateInput
	if got == nil || got.Encoding != nil || got.Data != nil || got.GetExtStoreId() != "store" || got.GetExtPath() != "path" {
		t.Fatalf("start input = %+v, want external reference with omitted encoding", got)
	}
	if input.GetData() != data || input.Encoding != nil {
		t.Fatalf("original input changed: %+v", input)
	}
}
