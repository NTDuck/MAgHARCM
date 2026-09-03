package agents_test

import "testing"

func TestPlateauDetectorEmptyRunner(t *testing.T) {
	TestPlateauDetectorEmpty(t)
}

func TestPlateauDetectorSingleSampleRunner(t *testing.T) {
	TestPlateauDetectorSingleSample(t)
}

func TestPlateauDetectorBelowThresholdRunner(t *testing.T) {
	TestPlateauDetectorBelowThreshold(t)
}

func TestPlateauDetectorAboveThresholdRunner(t *testing.T) {
	TestPlateauDetectorAboveThreshold(t)
}

func TestPlateauDetectorMixedThenRecoverRunner(t *testing.T) {
	TestPlateauDetectorMixedThenRecover(t)
}

func TestPlateauDetectorWindowOneRunner(t *testing.T) {
	TestPlateauDetectorWindowOne(t)
}

func TestCoverageDeltaRunner(t *testing.T) {
	TestCoverageDelta(t)
}

func TestPlateauDetectorSummaryRunner(t *testing.T) {
	TestPlateauDetectorSummary(t)
}

func TestPlateauDetectorRecordReturnsPlateauStateRunner(t *testing.T) {
	TestPlateauDetectorRecordReturnsPlateauState(t)
}
