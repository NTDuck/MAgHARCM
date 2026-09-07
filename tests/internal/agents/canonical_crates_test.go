package agents_test

import "testing"

func TestCanonicalizeCargoTomlRunner(t *testing.T) {
	TestCanonicalizeCargoToml(t)
}

func TestCrateCanonicalHintsRunner(t *testing.T) {
	TestCrateCanonicalHints(t)
}

func TestRenameCrateKeyRunner(t *testing.T) {
	TestRenameCrateKey(t)
}
