package main_test

import "testing"

func TestSetFieldRunner(t *testing.T) {
	TestSetField(t)
}

func TestSetFieldIntRunner(t *testing.T) {
	TestSetFieldInt(t)
}

func TestSetFieldUnknownRunner(t *testing.T) {
	TestSetFieldUnknown(t)
}

func TestPhase1FinalizeMissingFieldsRunner(t *testing.T) {
	TestPhase1FinalizeMissingFields(t)
}

func TestPhase1FinalizeWritesYAMLRunner(t *testing.T) {
	TestPhase1FinalizeWritesYAML(t)
}

func TestHandleSlashShowRunner(t *testing.T) {
	TestHandleSlashShow(t)
}

func TestHandleSlashClearResetsRunner(t *testing.T) {
	TestHandleSlashClearResets(t)
}

func TestHandleSlashLoadMalformedRunner(t *testing.T) {
	TestHandleSlashLoadMalformed(t)
}

func TestHandleSlashLoadValidJumpsToPhase2Runner(t *testing.T) {
	TestHandleSlashLoadValidJumpsToPhase2(t)
}

func TestHandleSlashUnknownCommandRunner(t *testing.T) {
	TestHandleSlashUnknownCommand(t)
}

func TestHandleSlashQuitRunner(t *testing.T) {
	TestHandleSlashQuit(t)
}

func TestHandleSlashDebugTogglesRunner(t *testing.T) {
	TestHandleSlashDebugToggles(t)
}

func TestHandleSlashLogsReadsSnapshotRunner(t *testing.T) {
	TestHandleSlashLogsReadsSnapshot(t)
}

func TestWriteYAMLRoundTripRunner(t *testing.T) {
	TestWriteYAMLRoundTrip(t)
}
