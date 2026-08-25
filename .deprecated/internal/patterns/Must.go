package patterns

// https://dev.to/dubjay18/gos-must-pattern-streamline-your-error-handling-27ff
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
