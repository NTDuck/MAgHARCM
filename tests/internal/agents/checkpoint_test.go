package agents_test

import "testing"

func TestCheckpointSaveLoadRoundTripRunner(t *testing.T) {
	TestCheckpointSaveLoadRoundTrip(t)
}

func TestCheckpointLoadLatestEmptyRunner(t *testing.T) {
	TestCheckpointLoadLatestEmpty(t)
}

func TestCheckpointLoadLatestPicksHighestIterationRunner(t *testing.T) {
	TestCheckpointLoadLatestPicksHighestIteration(t)
}

func TestCheckpointSaveEmptyRunIDReturnsErrorRunner(t *testing.T) {
	TestCheckpointSaveEmptyRunIDReturnsError(t)
}

func TestCheckpointCleanupRemovesDirectoryRunner(t *testing.T) {
	TestCheckpointCleanupRemovesDirectory(t)
}

func TestCheckpointVersionMismatchReturnsErrorRunner(t *testing.T) {
	TestCheckpointVersionMismatchReturnsError(t)
}
