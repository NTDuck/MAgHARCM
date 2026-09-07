package runner_test

import "testing"

func TestRunRejectsNilConfigRunner(t *testing.T) {
	TestRunRejectsNilConfig(t)
}

func TestRunRejectsEmptyDirsRunner(t *testing.T) {
	TestRunRejectsEmptyDirs(t)
}

func TestSuccessNilStateRunner(t *testing.T) {
	TestSuccessNilState(t)
}
