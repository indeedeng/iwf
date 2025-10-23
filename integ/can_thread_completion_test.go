package integ

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/can_thread_completion"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
)

// CAN = Continue-As-New
// TestCANThreadCompletionTemporal tests that all command threads complete before continue-as-new
// snapshots state. This validates the fix for the bug where internal channel signals were lost
// during continue-as-new.
//
// The bug scenario:
//  1. State1 WaitUntil sets up timer, signal, and internal channel commands
//  2. State1 Execute publishes to internal channel and moves to State2
//  3. Continue-as-new is triggered
//  4. BEFORE THE FIX: The internal channel signal would be lost because AddPotentialStateExecutionToResume
//     was called before waiting for command threads to complete
//  5. AFTER THE FIX: We wait for all command threads to complete before snapshotting, so the signal is preserved
func TestCANThreadCompletionTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestCANThreadCompletion(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()

		// Test with continue-as-new threshold set very low to force CAN during the workflow
		doTestCANThreadCompletion(t, service.BackendTypeTemporal, &iwfidl.WorkflowConfig{
			ContinueAsNewThreshold: ptr.Any(int32(1)), // Trigger CAN very quickly
		})
		smallWaitForFastTest()
	}
}

func TestCANThreadCompletionCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestCANThreadCompletion(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()

		// Test with continue-as-new threshold set very low to force CAN during the workflow
		doTestCANThreadCompletion(t, service.BackendTypeCadence, &iwfidl.WorkflowConfig{
			ContinueAsNewThreshold: ptr.Any(int32(1)), // Trigger CAN very quickly
		})
		smallWaitForFastTest()
	}
}

// TestAnyCommandCompletedWithCANTemporal validates the critical fix:
// With ANY_COMMAND_COMPLETED, if one command completes and CAN is triggered,
// we should NOT wait for other unfinished commands before proceeding to execute.
// This test ensures we only wait for threads that have retrieved data.
func TestAnyCommandCompletedWithCANTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	doTestAnyCommandCompletedWithCAN(t, service.BackendTypeTemporal)
}

func TestAnyCommandCompletedWithCANCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	doTestAnyCommandCompletedWithCAN(t, service.BackendTypeCadence)
}

func doTestAnyCommandCompletedWithCAN(t *testing.T, backendType service.BackendType) {
	// Start test workflow server
	wfHandler := can_thread_completion.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

	// Start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	wfId := "any_cmd_can_test_" + strconv.Itoa(int(time.Now().UnixNano()))

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startTime := time.Now()
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        can_thread_completion.WorkflowType,
		WorkflowTimeoutSeconds: 60,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(can_thread_completion.StateAnyCmd),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: &iwfidl.WorkflowConfig{
				ContinueAsNewThreshold: ptr.Any(int32(1)), // Force CAN immediately
			},
		},
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// Send signal quickly (timer will take much longer)
	go func() {
		time.Sleep(500 * time.Millisecond)
		signalReq := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
		httpResp, err := signalReq.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
			WorkflowId:        wfId,
			SignalChannelName: "any-cmd-signal",
			SignalValue: &iwfidl.EncodedObject{
				Encoding: ptr.Any("json"),
				Data:     ptr.Any("signal-data"),
			},
		}).Execute()
		if err != nil {
			t.Logf("Warning: Failed to send signal: %v", err)
		}
		if httpResp != nil && httpResp.StatusCode != 200 {
			t.Logf("Warning: Signal returned non-200 status: %d", httpResp.StatusCode)
		}
	}()

	assertions := assert.New(t)

	// Wait for workflow to complete (with longer timeout since we're testing timing)
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId:      wfId,
		WaitTimeSeconds: ptr.Any(int32(20)), // Wait up to 20 seconds
	}).Execute()
	if err != nil {
		assertions.False(strings.Contains(err.Error(), "420"), "Workflow took too long to complete, which means that ANY_COMMAND_COMPLETED command was not completed as expected. Most likely the command was also waiting for the timer command, when it can complete after receiving a signal.")
	}
	failTestAtHttpError(err, httpResp, t)

	elapsedTime := time.Since(startTime)

	// Get test results
	history, data := wfHandler.GetTestResult()

	// Verify StateAnyCmd executed
	assertions.Equalf(map[string]int64{
		"StateAnyCmd_start":  1,
		"StateAnyCmd_decide": 1,
	}, history, "State execution history mismatch: %v", history)

	// Verify workflow completed successfully
	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus(),
		"Workflow should complete successfully")

	// CRITICAL: Verify signal was received (proves ANY_COMMAND_COMPLETED triggered)
	signalReceived, ok := data["any_cmd_signal_received"].(bool)
	assertions.True(ok, "any_cmd_signal_received data should be present")
	assertions.True(signalReceived,
		"ANY_COMMAND_COMPLETED: Signal should have been received, triggering the Execute API")

	// Verify we didn't wait for the long timer
	// The timer is set for 20 seconds, but with ANY_COMMAND_COMPLETED + signal,
	// the workflow should complete in ~1-2 seconds, not 20+
	assertions.Less(elapsedTime, 5*time.Second,
		"Workflow took %v, which suggests we waited for the long timer. "+
			"With ANY_COMMAND_COMPLETED, we should proceed as soon as the signal is received, "+
			"not wait for all threads.", elapsedTime)

	t.Logf("✅ ANY_COMMAND_COMPLETED + CAN test passed in %v (expected < 5s, timer was 20s)", elapsedTime)
}

func doTestCANThreadCompletion(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// Start test workflow server
	wfHandler := can_thread_completion.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

	// Start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	wfId := can_thread_completion.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        can_thread_completion.WorkflowType,
		WorkflowTimeoutSeconds: 30,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(can_thread_completion.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// Send a signal after a short delay to trigger the signal thread completion
	// This ensures all three thread types (timer, signal, channel) complete their full path
	go func() {
		time.Sleep(500 * time.Millisecond)
		signalReq := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
		signalValue := iwfidl.EncodedObject{
			Encoding: ptr.Any("json"),
			Data:     ptr.Any("signal-data"),
		}
		httpResp, err := signalReq.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
			WorkflowId:        wfId,
			SignalChannelName: "test-signal",
			SignalValue:       &signalValue,
		}).Execute()
		if err != nil {
			t.Logf("Warning: Failed to send signal: %v", err)
		}
		if httpResp != nil && httpResp.StatusCode != 200 {
			t.Logf("Warning: Signal returned non-200 status: %d", httpResp.StatusCode)
		}
	}()

	// Wait for workflow to complete
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// Get test results
	history, data := wfHandler.GetTestResult()
	assertions := assert.New(t)

	// Verify all three states were executed
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
		"S3_start":  1,
		"S3_decide": 1,
	}, history, "CAN thread completion test failed - state execution history mismatch: %v", history)

	// Verify workflow completed successfully
	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus(),
		"Workflow should complete successfully")

	// ========== State1 Assertions: All three thread types complete ==========

	// 1. Timer thread completion
	s1TimerFired, ok := data["s1_timer_fired"].(bool)
	assertions.True(ok, "s1_timer_fired data should be present")
	assertions.True(s1TimerFired,
		"TIMER THREAD: Timer should have fired in State1. "+
			"This tests the timer thread completion path (lines 690-716 in workflowImpl.go)")

	// 2. Signal thread completion
	s1SignalReceived, ok := data["s1_signal_received"].(bool)
	assertions.True(ok, "s1_signal_received data should be present")
	assertions.True(s1SignalReceived,
		"SIGNAL THREAD: Signal should have been RECEIVED in State1. "+
			"This tests the signal thread completion path (lines 719-752 in workflowImpl.go)")

	// 3. Internal channel thread completion
	s1ChannelReceived, ok := data["s1_channel_received"].(bool)
	assertions.True(ok, "s1_channel_received data should be present")
	assertions.True(s1ChannelReceived,
		"CHANNEL THREAD: Internal channel should have been received in State1. "+
			"This tests the internal channel thread completion path (lines 755-791 in workflowImpl.go)")

	// ========== State2 Assertion: Continue-as-new preservation ==========

	// THE KEY ASSERTION: Verify that the internal channel signal was preserved through CAN
	s2ChannelReceived, ok := data["s2_channel_received"].(bool)
	assertions.True(ok, "s2_channel_received data should be present")
	assertions.True(s2ChannelReceived,
		"CONTINUE-AS-NEW PRESERVATION: Internal channel signal was lost during continue-as-new! "+
			"The channel published by State1 Execute should have been received by State2. "+
			"This validates that AddPotentialStateExecutionToResume waits for all threads to complete "+
			"before snapshotting state, ensuring no data is lost across continue-as-new.")

	// Verify channel value was correctly received
	if s2ChannelReceived {
		channelValue, ok := data["s2_channel_value"].(iwfidl.EncodedObject)
		assertions.True(ok, "s2_channel_value should be present when channel is received")
		if ok {
			assertions.Equal("json", *channelValue.Encoding, "Channel encoding should match")
			assertions.Equal("channel-data", *channelValue.Data, "Channel data should match")
		}
	}

	// ========== State3 Assertion: Timer thread in isolation (bonus test) ==========

	s3TimerFired, ok := data["s3_timer_fired"].(bool)
	assertions.True(ok, "s3_timer_fired data should be present")
	assertions.True(s3TimerFired,
		"TIMER THREAD (ISOLATED): Timer should have fired in State3. "+
			"This provides an isolated test of the timer thread path without other commands, "+
			"complementing the multi-command test in State1.")
}
