package logger_test

import "testing"

func TestEmitWritesExpectedLinesRunner(t *testing.T) {
	TestEmitWritesExpectedLines(t)
}

func TestEmitThreadsafeRunner(t *testing.T) {
	TestEmitThreadsafe(t)
}
