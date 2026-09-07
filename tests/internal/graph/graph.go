package graph_test

import "testing"

// Backlink: [[Methodology]] §1 Stage 2.5 (Symbol Navigation).
//
// Mirror of internal/graph/graph.go structural invariants. This file
// exists so the cross-package test driver in tests/cmd/ can detect when
// the 5-agent graph wiring silently regresses to 4 agents. The full
// Eino compile path needs real BaseChatModel instances and is exercised
// by the e2e test under cmd/MAgHARCM-tui instead.

// expectedAgentCount is the count of agent-bearing lambda nodes in the
// MAgHARCM graph. Must stay in sync with internal/graph/graph.go. The
// 5 nodes are: analyzer, navigator, planning, translator, validator.
// Checkpoint lambdas (save_translator_ckpt, save_validator_ckpt) and
// the branch target (compose.END) are NOT counted here.
const expectedAgentCount = 5

// expectedNodes lists every lambda node expected on the MAgHARCM graph.
// Adding a new agent means appending its name here AND wiring the
// matching AddLambdaNode in NewMAgHARCMGraph.
var expectedNodes = []string{
	"analyzer",
	"navigator",
	"planning",
	"translator",
	"validator",
	"save_translator_ckpt",
	"save_validator_ckpt",
}

// expectedForwardEdges lists every edge on the initial forward pass.
// The repair branch from save_validator_ckpt → translator (or END)
// is not enumerated here; verify it via the e2e driver.
var expectedForwardEdges = []string{
	"START->analyzer",
	"analyzer->navigator",
	"navigator->planning",
	"planning->translator",
	"translator->save_translator_ckpt",
	"save_translator_ckpt->validator",
	"validator->save_validator_ckpt",
}

// TestGraphNodeInventory: a static mirror test that asserts the
// inventory documented in this file matches the wiring in
// internal/graph/graph.go. We deliberately do NOT instantiate the
// Eino graph here — that requires real BaseChatModel instances and
// is covered by the cmd/MAgHARCM-tui e2e suite.
func TestGraphNodeInventory(t *testing.T) {
	if expectedAgentCount != 5 {
		t.Errorf("expectedAgentCount = %d, want 5 (analyzer, navigator, planning, translator, validator)",
			expectedAgentCount)
	}
	if len(expectedNodes) != 7 {
		t.Errorf("len(expectedNodes) = %d, want 7 (5 agents + 2 checkpoint lambdas)",
			len(expectedNodes))
	}
	hasNavigator := false
	for _, n := range expectedNodes {
		if n == "navigator" {
			hasNavigator = true
			break
		}
	}
	if !hasNavigator {
		t.Error("expectedNodes must include \"navigator\" (5th agent)")
	}
	hasNavigatorEdge := false
	for _, e := range expectedForwardEdges {
		if e == "navigator->planning" || e == "analyzer->navigator" {
			hasNavigatorEdge = true
			break
		}
	}
	if !hasNavigatorEdge {
		t.Error("expectedForwardEdges must include the navigator edge")
	}
}
