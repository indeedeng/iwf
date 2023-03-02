package retry

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/cadence/workflow"
	"time"
)

func ConvertCadenceWorkflowRetryPolicy(policy *iwfidl.WorkflowRetryPolicy) *workflow.RetryPolicy {
	if policy == nil {
		return nil
	}
	fillWorkflowRetryPolicyDefault(policy)

	return &workflow.RetryPolicy{
		InitialInterval:    time.Second * time.Duration(policy.GetInitialIntervalSeconds()),
		MaximumInterval:    time.Second * time.Duration(policy.GetMaximumIntervalSeconds()),
		MaximumAttempts:    policy.GetMaximumAttempts(),
		BackoffCoefficient: float64(policy.GetBackoffCoefficient()),
	}
}

func ConvertCadenceActivityRetryPolicy(policy *iwfidl.RetryPolicy) *workflow.RetryPolicy {
	if policy == nil {
		return nil
	}
	fillActivityRetryPolicyDefault(policy)

	// in Cadence, ExpirationInterval is the timeout include all retries
	expirationInterval := time.Duration(0)
	if policy.GetMaximumAttemptsDurationSeconds() > 0 {
		expirationInterval = time.Second * time.Duration(policy.GetMaximumAttemptsDurationSeconds())
	}

	return &workflow.RetryPolicy{
		InitialInterval:    time.Second * time.Duration(policy.GetInitialIntervalSeconds()),
		MaximumInterval:    time.Second * time.Duration(policy.GetMaximumIntervalSeconds()),
		MaximumAttempts:    policy.GetMaximumAttempts(),
		BackoffCoefficient: float64(policy.GetBackoffCoefficient()),
		ExpirationInterval: expirationInterval,
	}
}

func ConvertTemporalWorkflowRetryPolicy(policy *iwfidl.WorkflowRetryPolicy) *temporal.RetryPolicy {
	if policy == nil {
		return nil
	}
	fillWorkflowRetryPolicyDefault(policy)

	return &temporal.RetryPolicy{
		InitialInterval:    time.Second * time.Duration(policy.GetInitialIntervalSeconds()),
		MaximumInterval:    time.Second * time.Duration(policy.GetMaximumIntervalSeconds()),
		MaximumAttempts:    policy.GetMaximumAttempts(),
		BackoffCoefficient: float64(policy.GetBackoffCoefficient()),
	}
}

func ConvertTemporalActivityRetryPolicy(policy *iwfidl.RetryPolicy) *temporal.RetryPolicy {
	if policy == nil {
		return nil
	}
	fillActivityRetryPolicyDefault(policy)

	return &temporal.RetryPolicy{
		InitialInterval:    time.Second * time.Duration(policy.GetInitialIntervalSeconds()),
		MaximumInterval:    time.Second * time.Duration(policy.GetMaximumIntervalSeconds()),
		MaximumAttempts:    policy.GetMaximumAttempts(),
		BackoffCoefficient: float64(policy.GetBackoffCoefficient()),
	}
}

func fillActivityRetryPolicyDefault(policy *iwfidl.RetryPolicy) {
	if policy.InitialIntervalSeconds == nil {
		policy.InitialIntervalSeconds = iwfidl.PtrInt32(1)
	}
	if policy.BackoffCoefficient == nil {
		policy.BackoffCoefficient = iwfidl.PtrFloat32(2)
	}
	if policy.MaximumIntervalSeconds == nil {
		policy.MaximumIntervalSeconds = iwfidl.PtrInt32(100)
	}
	if policy.MaximumAttempts == nil {
		policy.MaximumAttempts = iwfidl.PtrInt32(0)
	}
	if policy.MaximumAttemptsDurationSeconds == nil {
		policy.MaximumAttemptsDurationSeconds = iwfidl.PtrInt32(0)
	}
}

func fillWorkflowRetryPolicyDefault(policy *iwfidl.WorkflowRetryPolicy) {
	if policy.InitialIntervalSeconds == nil {
		policy.InitialIntervalSeconds = iwfidl.PtrInt32(1)
	}
	if policy.BackoffCoefficient == nil {
		policy.BackoffCoefficient = iwfidl.PtrFloat32(2)
	}
	if policy.MaximumIntervalSeconds == nil {
		policy.MaximumIntervalSeconds = iwfidl.PtrInt32(100)
	}
	if policy.MaximumAttempts == nil {
		policy.MaximumAttempts = iwfidl.PtrInt32(0)
	}
}
