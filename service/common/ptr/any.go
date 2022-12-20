package ptr

func Any[T any](obj T) *T {
	return &obj
}
