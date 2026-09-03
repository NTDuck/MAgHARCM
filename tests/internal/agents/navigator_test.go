package agents_test

import "testing"

func TestNavigatorLookupSymbolNoProviderRunner(t *testing.T) {
	TestNavigatorLookupSymbolNoProvider(t)
}

func TestNavigatorLookupSymbolWithMockProviderRunner(t *testing.T) {
	TestNavigatorLookupSymbolWithMockProvider(t)
}

func TestNavigatorLookupSymbolPartialFailureRunner(t *testing.T) {
	TestNavigatorLookupSymbolPartialFailure(t)
}

func TestNavigatorLookupSymbolAllFailuresRunner(t *testing.T) {
	TestNavigatorLookupSymbolAllFailures(t)
}

func TestNavigatorLookupSymbolsRunner(t *testing.T) {
	TestNavigatorLookupSymbols(t)
}

func TestNavigatorNewNavigatorNilRunner(t *testing.T) {
	TestNavigatorNewNavigatorNil(t)
}

func TestRefCountRunner(t *testing.T) {
	TestRefCount(t)
}

func TestProjectDirOrDotRunner(t *testing.T) {
	TestProjectDirOrDot(t)
}

func TestNewNavigatorAgentDisabledRunner(t *testing.T) {
	TestNewNavigatorAgentDisabled(t)
}

func TestNavigatorAgentNilStateRunner(t *testing.T) {
	TestNavigatorAgentNilState(t)
}

func TestNavigatorAgentResolvesSymbolsRunner(t *testing.T) {
	TestNavigatorAgentResolvesSymbols(t)
}

func TestNavigatorAgentEmptyMappingRunner(t *testing.T) {
	TestNavigatorAgentEmptyMapping(t)
}
