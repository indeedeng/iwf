package retry

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/cadence/workflow"
	"time"
)

func ConvertCadenceRetryPolicy(policy *iwfidl.RetryPolicy) *workflow.RetryPolicy {
	if policy == nil {
		return nil
	}
	fillRetryPolicyDefault(policy)

	return &workflow.RetryPolicy{
		InitialInterval:    time.Second * time.Duration(policy.GetInitialIntervalSeconds()),
		MaximumInterval:    time.Second * time.Duration(policy.GetMaximumIntervalSeconds()),
		MaximumAttempts:    policy.GetMaximumAttempts(),
		BackoffCoefficient: float64(policy.GetBackoffCoefficient()),
	}
}

func ConvertTemporalRetryPolicy(policy *iwfidl.RetryPolicy) *temporal.RetryPolicy {
	if policy == nil {
		return nil
	}
	fillRetryPolicyDefault(policy)

	return &temporal.RetryPolicy{
		InitialInterval:    time.Second * time.Duration(policy.GetInitialIntervalSeconds()),
		MaximumInterval:    time.Second * time.Duration(policy.GetMaximumIntervalSeconds()),
		MaximumAttempts:    policy.GetMaximumAttempts(),
		BackoffCoefficient: float64(policy.GetBackoffCoefficient()),
	}
}

func fillRetryPolicyDefault(policy *iwfidl.RetryPolicy) {
	if policy.InitialIntervalSeconds == nil {
		policy.InitialIntervalSeconds = iwfidl.PtrInt32(1)
	}
	if policy.BackoffCoefficient == nil {
		policy.BackoffCoefficient = iwfidl.PtrFloat32(2)
	}
	if policy.MaximumIntervalSeconds == nil {
		policy.MaximumAttempts = iwfidl.PtrInt32(100)
	}
	if policy.MaximumAttempts == nil {
		policy.MaximumAttempts = iwfidl.PtrInt32(0)
	}
}
