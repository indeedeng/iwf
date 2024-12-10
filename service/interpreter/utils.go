package interpreter

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"path/filepath"
	"runtime"
	"slices"
)

func caller(skip int) string {
	_, path, lineno, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v:%v", filepath.Base(path), lineno)
}

func LastCaller() string {
	return caller(3)
}

// DeterministicKeys returns the keys of a map in deterministic (sorted) order. To be used in for
// loops in workflows for deterministic iteration.
func DeterministicKeys[K constraints.Ordered, V any](m map[K]V) []K {
	// copy from https://github.com/temporalio/sdk-go/blob/7828e06cf517dd2d27881a33efaaf4ff985f2e14/workflow/workflow.go#L787
	// and example usage https://github.com/temporalio/samples-go/blob/c69dc0bacc78163a50465c6f80aa678739673a4d/safe_message_handler/workflow.go#L119
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	slices.Sort(r)
	return r
}
