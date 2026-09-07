package agents_test

import "testing"

func TestManifestRewriterBasicRewriteRunner(t *testing.T) {
	TestManifestRewriterBasicRewrite(t)
}

func TestManifestRewriterTestScopeRunner(t *testing.T) {
	TestManifestRewriterTestScope(t)
}

func TestManifestRewriterEmptyDepsRunner(t *testing.T) {
	TestManifestRewriterEmptyDeps(t)
}

func TestManifestRewriterAllUnmappedRunner(t *testing.T) {
	TestManifestRewriterAllUnmapped(t)
}

func TestDetectSourceLangRunner(t *testing.T) {
	TestDetectSourceLang(t)
}

func TestFormatDepLineRunner(t *testing.T) {
	TestFormatDepLine(t)
}

func TestQuoteFeaturesRunner(t *testing.T) {
	TestQuoteFeatures(t)
}

func TestManifestRewriterRoundTripWithCanonicalizeRunner(t *testing.T) {
	TestManifestRewriterRoundTripWithCanonicalize(t)
}
