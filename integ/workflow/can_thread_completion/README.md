# CAN Thread Completion Integration Test

## Purpose

This integration test validates the fix for the continue-as-new (CAN) thread completion bug where internal channel signals were lost during workflow continuation.

**What This Test Validates:**
- ✅ **All three command thread types complete** before state snapshotting (timer, signal, internal channel)
- ✅ **Data is preserved** across continue-as-new transitions
- ✅ **Thread completion works** both with multiple commands (State1) and in isolation (State3)
- ✅ **`AddPotentialStateExecutionToResume()`** waits for all threads before snapshotting

## Bug Description

**Original Bug:**
When continue-as-new was triggered, `AddPotentialStateExecutionToResume()` was called BEFORE waiting for command threads (timer, signal, internal channel) to complete. This caused a race condition where:

1. A goroutine would call `InternalChannel.Retrieve()` to move data from the channel
2. But before the data was assigned to `completedInterStateChannelCmds[idx]`
3. `AddPotentialStateExecutionToResume()` would snapshot the state
4. Result: The signal was removed from `InternalChannel` but not captured in the snapshot → **data loss**

**The Fix:**
Wait for all command threads to complete BEFORE calling `AddPotentialStateExecutionToResume()`. This ensures all `Retrieve()` operations complete and data is fully populated in the `completedXXXCmds` maps before snapshotting.

## Test Workflow Structure

### Workflow Flow Diagram

```
State1 (ALL_COMMAND_COMPLETED)
  ├─ Timer (2s)      → fires → completes
  ├─ Signal          → ext signal → completes  
  └─ Channel         → published → completes
         │
         ├──────────────┬──────────────┐
         │              │              │
      State2         State3         (CAN may trigger here)
    (Channel Test)  (Timer Test)
         │              │
    Dead-end       Completes Workflow
```

### Detailed State Descriptions

### State1 (All Three Thread Types Test)
- **WaitUntil API:**
  - Sets up **timer command** (2 seconds)
  - Sets up **signal command** (waits for external signal)
  - Sets up **internal channel command** (waits for channel data)
  - **Immediately publishes** to internal channel (so it's available)
  - Uses `ALL_COMMAND_COMPLETED` trigger type ← **All three threads must complete!**
  
- **Execute API:**
  - ✅ **Verifies timer fired** (tests timer thread completion path)
  - ✅ **Verifies signal received** (tests signal thread completion path)
  - ✅ **Verifies channel received** (tests internal channel thread completion path)
  - Publishes to "test-channel2" for State2
  - **Moves to BOTH State2 AND State3** (parallel execution)

**External Signal:** Test code sends a signal after 500ms to trigger signal thread completion

### State2 (Continue-As-New Preservation Test)
- **WaitUntil API:**
  - Waits for "test-channel2" published by State1 Execute
  
- **Execute API:**
  - **Verifies channel was received** ← This proves data preserved through continue-as-new
  - Ends in **dead-end** (lets State3 complete the workflow)

### State3 (Isolated Timer Thread Test)
- **WaitUntil API:**
  - **Only** sets up a timer command (2 seconds)
  - No signal or channel commands
  
- **Execute API:**
  - **Verifies timer fired**
  - **Completes workflow gracefully**
  
**Purpose:** Provides an isolated test of the timer thread path without other command types, complementing the multi-command test in State1

## Test Scenarios

The test runs in two configurations:

1. **Normal execution** (`config=nil`)
   - Workflow may or may not trigger continue-as-new
   - Validates basic functionality

2. **Forced continue-as-new** (`ContinueAsNewThreshold=1`)
   - Forces continue-as-new to trigger very quickly
   - Ensures the fix works when CAN happens during state transitions
   - This is the critical scenario that exposed the bug

## Key Assertions

### State1 - All Three Thread Types Complete Together

1. ✅ **Timer thread** (`s1_timer_fired == true`)
   - Tests timer thread completion path (lines 690-716 in workflowImpl.go)
   - Timer fires → stores in `completedTimerCmds` → `waitForThreads[threadName] = true`

2. ✅ **Signal thread** (`s1_signal_received == true`)
   - Tests signal thread completion path (lines 719-752 in workflowImpl.go)
   - External signal sent → `Retrieve()` → stores in `completedSignalCmds` → `waitForThreads[threadName] = true`

3. ✅ **Internal channel thread** (`s1_channel_received == true`)
   - Tests internal channel thread completion path (lines 755-791 in workflowImpl.go)
   - Channel published → `Retrieve()` → stores in `completedInterStateChannelCmds` → `waitForThreads[threadName] = true`

### State2 - Continue-As-New Preservation

4. ⚠️ **CRITICAL** (`s2_channel_received == true`)
   - Confirms S2 received the channel published by S1's Execute
   - **If this fails** → Bug exists! Data was lost during continue-as-new because `AddPotentialStateExecutionToResume()` was called before threads completed
   - **If this passes** → Bug is fixed! Data was preserved because we waited for all threads before snapshotting

5. **Channel value validation** (`s2_channel_value`)
   - Verifies the received data matches what was published (encoding and data content)

### State3 - Isolated Timer Thread Test

6. ✅ **Timer thread (isolated)** (`s3_timer_fired == true`)
   - Validates timer thread works correctly without other command types
   - Provides complementary test to the multi-command scenario in State1

### Overall Validation

7. **State execution history**: Verifies all three states (S1, S2, S3) executed exactly once
8. **Workflow completion**: Ensures workflow completes successfully after all validations

## Code Changes Required

The fix is in `service/interpreter/workflowImpl.go`:

```go
// BEFORE (buggy):
continueAsNewer.AddPotentialStateExecutionToResume(...)  // ← Called too early!
_ = provider.Await(ctx, func() bool {
    return IsDeciderTriggerConditionMet(...) || continueAsNewCounter.IsThresholdMet()
})
commandReqDoneOrCanceled = true
// Wait for threads... (but too late!)

// AFTER (fixed):
// Wait for all threads FIRST
if globalVersioner.IsAfterVersionOfWaitingCommandThreads() {
    if err := provider.Await(ctx, func() bool {
        for _, completed := range waitForThreads {
            if !completed { return false }
        }
        return true
    }); err != nil {
        return nil, service.WaitingCommandsStateExecutionStatus, err
    }
}
continueAsNewer.AddPotentialStateExecutionToResume(...)  // ← Now it's safe!
```

## Running the Test

```bash
# Temporal backend
go test -v ./integ -run TestCANThreadCompletionTemporal

# Cadence backend  
go test -v ./integ -run TestCANThreadCompletionCadence
```

## Expected Outcome

✅ **With the fix:** All assertions pass, including:
- State1: All three thread types complete (`s1_timer_fired`, `s1_signal_received`, `s1_channel_received`)
- State2: Channel preserved through CAN (`s2_channel_received == true`)
- State3: Isolated timer works (`s3_timer_fired == true`)

❌ **Without the fix:** The test would fail with:
```
CONTINUE-AS-NEW PRESERVATION: Internal channel signal was lost during continue-as-new!
The channel published by State1 Execute should have been received by State2.
This validates that AddPotentialStateExecutionToResume waits for all threads to complete
before snapshotting state, ensuring no data is lost across continue-as-new.
```

## Related Files

- Test: `integ/can_thread_completion_test.go`
- Workflow Handler: `integ/workflow/can_thread_completion/routers.go`  
- Bug Fix: `service/interpreter/workflowImpl.go` (lines ~795-810)
- Replay Test: `replayTests/history/eval.json` (historical bug evidence)

