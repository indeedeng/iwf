package helpers

import "testing"

func FailTestWithError(error error, t *testing.T) {
	t.Errorf("%s - Test failed with error: %v", t.Name(), error)
}

func FailTestWithErrorMessage(errorMessage string, t *testing.T) {
	t.Errorf("%s - Test failed with error: %s", t.Name(), errorMessage)
}
