package pattern

// Must wraps a call to a function returning (T, error) and panics if the error is non-nil.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// Must0 wraps a call to a function returning error and panics if the error is non-nil.
func Must0(err error) {
	if err != nil {
		panic(err)
	}
}
