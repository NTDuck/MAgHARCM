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
