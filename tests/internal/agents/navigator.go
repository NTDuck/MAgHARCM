package agents_test

import (
	"MAgHARCM/internal/agents"
	"context"
	"errors"
	"reflect"
	"testing"

	"MAgHARCM/internal/tools"
)

// mockLSPProvider is a programmable LSPProvider used to drive the agents.Navigator
// through deterministic test paths (success, partial failure, full failure).
type mockLSPProvider struct {
	defOut    *tools.DefinitionOutput
	defErr    error
	refOut    *tools.ReferencesOutput
	refErr    error
	hoverOut  *tools.HoverOutput
	hoverErr  error
	defCalls  int
	refCalls  int
	hoverCalls int
}

func (m *mockLSPProvider) Name() string { return "mock" }

func (m *mockLSPProvider) GetDefinition(ctx context.Context, in *tools.DefinitionInput) (*tools.DefinitionOutput, error) {
	m.defCalls++
	return m.defOut, m.defErr
}

func (m *mockLSPProvider) GetDiagnostics(ctx context.Context, in *tools.DiagnosticsInput) (*tools.DiagnosticsOutput, error) {
	return nil, nil
}

func (m *mockLSPProvider) EditFile(ctx context.Context, in *tools.EditFileInput) (*tools.EditFileOutput, error) {
	return nil, nil
}

func (m *mockLSPProvider) GetHover(ctx context.Context, in *tools.HoverInput) (*tools.HoverOutput, error) {
	m.hoverCalls++
	return m.hoverOut, m.hoverErr
}

func (m *mockLSPProvider) GetReferences(ctx context.Context, in *tools.ReferencesInput) (*tools.ReferencesOutput, error) {
	m.refCalls++
	return m.refOut, m.refErr
}

func (m *mockLSPProvider) RenameSymbol(ctx context.Context, in *tools.RenameSymbolInput) (*tools.RenameSymbolOutput, error) {
	return nil, nil
}

// TestNavigatorLookupSymbolNoProvider: passing nil provider yields the
// sentinel agents.ErrNoLSPProvider and no sub-calls are made.
func TestNavigatorLookupSymbolNoProvider(t *testing.T) {
	n := agents.NewNavigator(nil)
	res := n.LookupSymbol(context.Background(), "foo", "/tmp/foo.c")

	if res.Symbol != "foo" {
		t.Fatalf("expected Symbol=%q, got %q", "foo", res.Symbol)
	}
	if !errors.Is(res.Error, agents.ErrNoLSPProvider) {
		t.Fatalf("expected agents.ErrNoLSPProvider, got %v", res.Error)
	}
	if res.Definition != nil || res.References != nil || res.Hover != nil {
		t.Fatalf("expected all sub-results nil, got def=%v refs=%v hover=%v",
			res.Definition, res.References, res.Hover)
	}
}

// TestNavigatorLookupSymbolWithMockProvider: a fully-successful mock provider
// should populate all three sub-results and leave Error=nil.
func TestNavigatorLookupSymbolWithMockProvider(t *testing.T) {
	mock := &mockLSPProvider{
		defOut: &tools.DefinitionOutput{
			Symbol: "foo",
			Definitions: []tools.DefinitionLocation{
				{FilePath: "/tmp/foo.c", Line: 42, Signature: "int foo(int x)", Snippet: "int foo(int x) {"},
			},
		},
		refOut: &tools.ReferencesOutput{
			Symbol:     "foo",
			Total:      2,
			References: []tools.ReferenceMatch{{FilePath: "a.c", Line: 1, LineText: "foo"}, {FilePath: "b.c", Line: 9, LineText: "foo"}},
		},
		hoverOut: &tools.HoverOutput{Symbol: "foo", Signature: "int foo(int x)", Doc: "Returns foo", Found: true},
	}
	n := agents.NewNavigator(mock)
	res := n.LookupSymbol(context.Background(), "foo", "/tmp/foo.c")

	if res.Error != nil {
		t.Fatalf("expected nil Error, got %v", res.Error)
	}
	if res.Definition == nil || !reflect.DeepEqual(res.Definition.Definitions, mock.defOut.Definitions) {
		t.Fatalf("Definition mismatch: got %v", res.Definition)
	}
	if res.References == nil || len(res.References.References) != 2 {
		t.Fatalf("References mismatch: got %+v", res.References)
	}
	if res.Hover == nil || res.Hover.Signature != "int foo(int x)" {
		t.Fatalf("Hover mismatch: got %+v", res.Hover)
	}
	if mock.defCalls != 1 || mock.refCalls != 1 || mock.hoverCalls != 1 {
		t.Fatalf("expected 1 call per sub-tool, got def=%d ref=%d hover=%d",
			mock.defCalls, mock.refCalls, mock.hoverCalls)
	}
}

// TestNavigatorLookupSymbolPartialFailure: one sub-call fails, the others
// succeed; the failing call's output must be nil and the overall Error must
// still be set (with the other fields populated).
func TestNavigatorLookupSymbolPartialFailure(t *testing.T) {
	refBoom := errors.New("references exploded")
	mock := &mockLSPProvider{
		defOut: &tools.DefinitionOutput{
			Symbol:      "foo",
			Definitions: []tools.DefinitionLocation{{FilePath: "/tmp/foo.c", Line: 1, Signature: "foo"}},
		},
		refErr:   refBoom,
		hoverOut: &tools.HoverOutput{Symbol: "foo", Signature: "foo()", Found: true},
	}
	n := agents.NewNavigator(mock)
	res := n.LookupSymbol(context.Background(), "foo", "/tmp/foo.c")

	if res.Error == nil {
		t.Fatalf("expected non-nil Error, got nil")
	}
	if !errors.Is(res.Error, refBoom) {
		t.Fatalf("expected Error=%v, got %v", refBoom, res.Error)
	}
	if res.Definition == nil {
		t.Fatalf("expected partial Definition result despite references failure")
	}
	if res.References != nil {
		t.Fatalf("expected References nil on failure, got %+v", res.References)
	}
	if res.Hover == nil {
		t.Fatalf("expected partial Hover result despite references failure")
	}
}

// TestNavigatorLookupSymbolAllFailures: every sub-call fails; Error must be
// the first error reported and all sub-results must be nil.
func TestNavigatorLookupSymbolAllFailures(t *testing.T) {
	mock := &mockLSPProvider{
		defErr:   errors.New("def failed"),
		refErr:   errors.New("ref failed"),
		hoverErr: errors.New("hover failed"),
	}
	n := agents.NewNavigator(mock)
	res := n.LookupSymbol(context.Background(), "foo", "/tmp/foo.c")

	if res.Error == nil {
		t.Fatalf("expected non-nil Error when all sub-calls fail")
	}
	if res.Definition != nil || res.References != nil || res.Hover != nil {
		t.Fatalf("expected all sub-results nil, got def=%v refs=%v hover=%v",
			res.Definition, res.References, res.Hover)
	}
}

// TestNavigatorLookupSymbols: batch lookup returns a map keyed by symbol name
// with one resolution per symbol.
func TestNavigatorLookupSymbols(t *testing.T) {
	mock := &mockLSPProvider{
		defOut: &tools.DefinitionOutput{Symbol: "x", Definitions: []tools.DefinitionLocation{{FilePath: "f.c", Line: 1}}},
		refOut: &tools.ReferencesOutput{Symbol: "x", Total: 0},
		hoverOut: &tools.HoverOutput{Symbol: "x", Found: true},
	}
	n := agents.NewNavigator(mock)
	out := n.LookupSymbols(context.Background(), []string{"alpha", "beta", "gamma"}, "/tmp")

	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	for _, sym := range []string{"alpha", "beta", "gamma"} {
		res, ok := out[sym]
		if !ok {
			t.Fatalf("expected key %q in result map", sym)
		}
		if res.Symbol != sym {
			t.Fatalf("expected Symbol=%q, got %q", sym, res.Symbol)
		}
		if res.Error != nil {
			t.Fatalf("expected nil Error for %q, got %v", sym, res.Error)
		}
		if res.Definition == nil {
			t.Fatalf("expected Definition populated for %q", sym)
		}
	}
	if mock.defCalls != 3 || mock.refCalls != 3 || mock.hoverCalls != 3 {
		t.Fatalf("expected 3 calls per sub-tool, got def=%d ref=%d hover=%d",
			mock.defCalls, mock.refCalls, mock.hoverCalls)
	}
}

// TestNavigatorNewNavigatorNil: agents.NewNavigator(nil) is a valid construction;
// agents.LookupSymbol must produce the no-provider sentinel.
func TestNavigatorNewNavigatorNil(t *testing.T) {
	n := agents.NewNavigator(nil)
	if n == nil {
		t.Fatal("agents.NewNavigator(nil) returned nil")
	}
	if n.Provider != nil {
		t.Fatalf("expected Provider nil, got %T", n.Provider)
	}
}

// TestRefCount: nil, empty, and populated references produce the expected counts.
func TestRefCount(t *testing.T) {
	if got := agents.RefCount(nil); got != 0 {
		t.Fatalf("agents.RefCount(nil) = %d, want 0", got)
	}
	if got := agents.RefCount(&tools.ReferencesOutput{}); got != 0 {
		t.Fatalf("agents.RefCount(empty) = %d, want 0", got)
	}
	populated := &tools.ReferencesOutput{
		References: []tools.ReferenceMatch{{Line: 1}, {Line: 2}, {Line: 3}, {Line: 4}, {Line: 5}},
	}
	if got := agents.RefCount(populated); got != 5 {
		t.Fatalf("agents.RefCount(populated) = %d, want 5", got)
	}
}

// TestProjectDirOrDot: file path → its directory; empty → ".".
func TestProjectDirOrDot(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "."},
		{"/tmp/foo.c", "/tmp"},
		{"relative/path.c", "relative"},
	}
	for _, tc := range cases {
		if got := agents.ProjectDirOrDot(tc.in); got != tc.want {
			t.Errorf("agents.ProjectDirOrDot(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
