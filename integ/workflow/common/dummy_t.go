package common

// DummyT uses assert.ElementsMatch for comparing slices, but with a bool result.
type DummyT struct{}

func (t DummyT) Errorf(string, ...interface{}) {}
