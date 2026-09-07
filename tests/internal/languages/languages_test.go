package languages_test

import "testing"

func TestLanguageRegistryLookupRunner(t *testing.T) {
	TestLanguageRegistryLookup(t)
}

func TestDynamicLanguageLoaderRunner(t *testing.T) {
	TestDynamicLanguageLoader(t)
}

func TestRuntimePuregoDlopenDlsymRunner(t *testing.T) {
	TestRuntimePuregoDlopenDlsym(t)
}

func TestExtractFileStructureMultiLanguageRunner(t *testing.T) {
	TestExtractFileStructureMultiLanguage(t)
}
