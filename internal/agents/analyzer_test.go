package agents

import (
	"testing"
)

func TestSelectMigrationStrategy(t *testing.T) {
	// 1. Small project -> BIG_BANG
	strat, _ := SelectMigrationStrategy(2, 200, true, true)
	if strat != "BIG_BANG" {
		t.Errorf("expected BIG_BANG for small repo, got %s", strat)
	}

	// 2. Large project -> PILOT
	strat, _ = SelectMigrationStrategy(60, 15000, true, true)
	if strat != "PILOT" {
		t.Errorf("expected PILOT for large repo, got %s", strat)
	}

	// 3. Untested project -> FROZEN_LEGACY
	strat, _ = SelectMigrationStrategy(10, 2000, false, true)
	if strat != "FROZEN_LEGACY" {
		t.Errorf("expected FROZEN_LEGACY for untested repo, got %s", strat)
	}

	// 4. Standard modular with tests -> PARALLEL_CUTOVER / INCREMENTAL
	strat, _ = SelectMigrationStrategy(15, 3000, true, true)
	if strat != "PARALLEL_CUTOVER" && strat != "INCREMENTAL" {
		t.Errorf("expected PARALLEL_CUTOVER or INCREMENTAL, got %s", strat)
	}
}
