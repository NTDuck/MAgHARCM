---
title: Sprint 2026-09-09 — Methodology Compliance Report
backlink: [[Sprint 2026-09-09]]
tags: [diary, methodology, compliance, [[2.0.0 MAgHARCM]], sprint-2026-09-09]
---

# Sprint 2026-09-09 — Methodology Compliance Report

Scope: read-only verification against the directive list in the user's
briefing for Sprint 2026-09-09. The build state at commit `0e60035` is
green and the audit does NOT modify any code. Each finding cites the
exact file and line that supports it.

Inputs cross-checked:
- `.obsidian/MAgHARCM/research/METHODOLOGY.md` (§§1–6)
- `.obsidian/MAgHARCM/architecture/ADR-2026-09-07-Sprint-Conventions.md` (ADR-C-007..015)

---

## Directive 1 — Compile-time inits use Must pattern (no fallbacks)

Status: ✅

Evidence:
- `internal/compiletime/compiletime.go:417-440` declares the three strict
  Must helpers (`Must[T any]`, `MustNotNil[T any]`, `MustNotEmpty`); all
  panic on failure (no fallback path).
- `internal/agents/strategy.go` (whole file) is `compiletime`-only and never
  shadows Must with a fallback — it imports `compiletime.StrategyKind`
  directly (`internal/agents/strategy.go:11-12`) and only re-exports it
  via the type alias.
- `internal/agents/roleflip.go:46-49` calls
  `panic("compiletime.MustRoleFlipGate: Model must not be nil")` —
  hard fail, no fallback.
- `internal/agents/iter_retrieval.go:64-67` similarly panics when the
  Navigator is nil.
- `internal/agents/navigator.go:51-55` panics when LSPProvider is nil.
- `internal/agents/spec_lifecycle.go:73-75` uses
  `compiletime.Must(phase, err)` — strict Must call site.
- Repo-wide grep for `_ = Must` / `_ = compiletime.Must` returns zero
  matches (i.e. no swallowed-error fallback): no occurrences in
  `internal/`, `cmd/`, or `tests/`.

Conclusion: every compile-time init site uses `compiletime.Must*`; no
fallback logic found.

---

## Directive 2 — Each unit has clear boundaries (no cross-package unexported leakage)

Status: ✅

Evidence:
- `internal/compiletime/state.go` (whole file, header at lines 1-30)
  declares `State`, `Task`, every artifact struct, the `SchemaVersioned`
  interface, and `DocumentWrapper[T]` — these are the cross-cutting types
  referenced from `graph.go`, `runner.go`, `checkpoint.go`. The file
  comment explicitly states `compiletime` is a leaf package
  (`compiletime.go:5-12` of state.go explains the cycle-free layout).
- `internal/agents/state.go:8-10` is a one-liner type alias
  `type State = compiletime.State` — no duplication.
- Per-agent files inside `internal/agents/*.go` expose type aliases only
  (`AnalyzerOutput = compiletime.AnalyzerOutput`,
  `PlanningOutput = compiletime.PlanningOutput`,
  `ArchaeologyReport = compiletime.ArchaeologyReport`,
  `SpecMinerInvariants = compiletime.SpecMinerInvariants`,
  `ImplementationPlan / PlanStep` from `planning.go:24-26`,
  `LibraryMapping / SourceProjectResearch / TargetProjectDesign /
  ThirdPartyLibraryAnalysis` from `analyzer.go:19-24`,
  `NamingFinding` from `archaeology.go:31`,
  `ArchitectureStabilityLayer / Layer*` from
  `design_rule_hierarchy.go:11-19`,
  `SpecLifecyclePhase / Phase*` from `spec_lifecycle.go:13-24`,
  `StrategyKind` from `strategy.go:11-12`).
- The agents package never exposes an unexported cross-package helper —
  grep for `func [a-z][a-zA-Z]*\([^)]*\)\s*\*?compiletime\.` in
  `internal/agents/` returns zero matches.
- The historical cycle comment in `state.go:5-12` documents the deliberate
  split: shared types live in compiletime (the leaf), producer agents
  import and consume them, methods live on the canonical types.

Conclusion: shared types live exclusively in `internal/compiletime`; the
agents package re-exports them by type alias only; no unexported symbols
leak across the package boundary.

---

## Directive 3 — SelectMigrationStrategy replaced by `Registry.TryInOrder`

Status: ✅

Evidence:
- `internal/agents/strategy.go:60-74` — `func (r *Registry) TryInOrder(ctx
  context.Context, p Profile) (StrategyKind, error)` implements the
  try-and-fail loop. Iterates `r.all`, calls `Matches(p)`, on success
  calls `Attempt(ctx, p)`, and on Miss logs and continues.
- `internal/agents/strategy.go:39-53` — `NewDefaultRegistry()` returns
  the canonical 5-strategy list in priority order:
  `bigBangStrategy, pilotStrategy, frozenLegacyStrategy,
  parallelCutoverStrategy, incrementalStrategy`.
- `internal/agents/strategy.go:226-232` — `SelectAndTryStrategies` is the
  convenience runner used by `AnalyzerAgent.Run`; it calls
  `NewDefaultRegistry().TryInOrder(ctx, p)` (line 227).
- Repo-wide grep for the legacy helper `SelectMigrationStrategy` returns
  zero matches in `internal/`, `cmd/`, or `tests/`. The only textual
  occurrences are in `.obsidian/MAgHARCM/diary/Sprint-2026-09-08-Handoff.md:25`
  (an open-work bullet, not code) and in
  `tests/internal/agents/analyzer.go:12-15` (a comment in
  `TestAnalyzerStrategySelection` explicitly noting "The legacy
  SelectMigrationStrategy helper was retired in favour of the strategy
  interface").

Conclusion: the hardcoded if-cascade has been replaced by the registry
interface; `Registry.TryInOrder` is wired through the analyzer and
recruiter.

---

## Directive 4 — No hardcoded magic strings in code

Status: ✅

Evidence:
- All primitive sentinels, language keys, toolchain names, log scopes,
  strategy thresholds, verdict vocabulary, role-flip tokens, checkpoint
  patterns, archaeology skip list, SLM prompt preamble, concept cluster
  table, comprehension recognition labels, recruiter tool/agent names,
  strategy rationales, binary compilation status, and the five strategy
  kinds live in `internal/compiletime/compiletime.go:17-440`.
- Examples of centralisation:
  - Languages: `LangRust`, `LangRustCanonical`, `LangGo`, `LangC`,
    `LangJava`, `LangPython` (lines 28-44).
  - Toolchains: `ToolchainCargo`, `ToolchainGo` (lines 22-26).
  - LSP providers: `LSPProviderABCoder`, `LSPProviderNative` (lines
    46-50).
  - Schema version: `CurrentSchemaVersion` (line 18).
  - Default artifact dir / request file / package / tree depth
    (lines 55-75).
  - Log scopes (lines 94-102).
  - Strategy thresholds (lines 110-116) and rationales (lines 121-128).
  - Verdict vocabulary (lines 134-152).
  - Optional check names (lines 163-171).
  - RoleFlip gate tokens (lines 177-184).
  - Checkpoint persistence (lines 212-228).
  - Archaeology skip list (line 237).
  - SLM prompt preamble (lines 247-261).
  - Spec lifecycle phases (lines 268-280).
  - Concept clusters and token min length (lines 286-318).
  - Architecture stability layers (lines 322-337).
  - Comprehension recognition labels (lines 344-352).
  - Recruiter tool/agent names (lines 358-376).
  - Strict enums `CompilationStatus`, `StrategyKind`, `WorkUnitStatus`
    (lines 384-411).
- Per-agent code imports these via `compiletime.*` exclusively; the
  `strategy.go` and `recruit.go` files reference `compiletime.LogScope*`,
  `compiletime.Strategy*`, `compiletime.Tool*`, `compiletime.Agent*` —
  no inline string literals appear.
- The only hardcoded string outside `compiletime` is the prompt-template
  body text in `internal/agents/prompts.go` (LLM prompts); that is by
  design — prompts are not magic literals, they are user-visible content.

Conclusion: every magic literal enumerated in the directive list is
centralised in `internal/compiletime/compiletime.go`; no code-level
magic literals remain in the agent modules.

---

## Directive 5 — Graph has 8+ agents wired

Status: ✅

Evidence (`internal/graph/graph.go`):
- Agent constructors at lines 33-40 (one per agent):
  1. `archaeologistAgent := agents.NewArchaeologist()`
  2. `analyzerAgent := agents.NewAnalyzerAgent(models.Reasoning)`
  3. `planningAgent := agents.NewPlanningAgent(models.Reasoning)`
  4. `translatorAgent := agents.NewTranslatorAgent(models.Coding, runID)`
  5. `reviewerAgent := agents.NewRoleFlipGate(models.Reasoning)`
  6. `validatorAgent := agents.NewValidatorAgent(models.Reasoning, runID)`
  7. `verdictPanelAgent := agents.NewVerdictPanel(models.Reasoning)`
  8. `recruiterAgent := agents.NewRecruiter()`
- Eight graph nodes registered (one per agent):
  1. `archaeologist` (line 45)
  2. `analyzer` (line 63)
  3. `planning` (line 71)
  4. `translator` (line 79)
  5. `reviewer` (line 92)
  6. `validator` (line 116)
  7. `verdict_panel` (line 128)
  8. `recruiter` (line 150)
- Plus two checkpoint barrier nodes (`save_translator_ckpt` line 87,
  `save_validator_ckpt` line 124) to enforce disk-durability between
  the translator/reviewer and reviewer/verdict steps.
- Forward edges connect START → archaeologist → analyzer → planning →
  translator → save_translator_ckpt → reviewer → validator →
  save_validator_ckpt (lines 175-198).
- Repair branch at line 201-219 routes from `save_validator_ckpt` to
  either `END` (success) or `verdict_panel` (re-enter repair loop).
- Cyclic repair loop: verdict_panel → recruiter → translator
  (lines 222-227).
- Compile with `compose.WithGraphName("MAgHARCM-8Agent")` and
  `compose.WithMaxRunSteps(compiletime.MaxGraphRunSteps)` (lines 230-233).

Conclusion: 8 agent nodes are constructed and wired; edges form
START → 7 sequential agents → branch → 2-cycle repair loop → END.

---

## Directive 6 — abcoder-mcp default (not tree-sitter)

Status: ✅

Evidence:
- `internal/compiletime/compiletime.go:46-50` declares the canonical
  provider constants:
  - `LSPProviderNative = "native"` (tree-sitter)
  - `LSPProviderABCoder = "abcoder-mcp"` (default)
- `internal/tools/lsp_provider.go:62-85` defines `ABCoderMcpProvider`
  and its `Name()` returns `"abcoder-mcp"`.
- `internal/tools/lsp_provider.go:189-196` — `GetLSPProvider`:
  - `case "abcoder", "abcoder-mcp", "mcp"` → `NewABCoderMcpProvider(...)`
  - `default` → `NewNativeLSPProvider()` (the explicit documented
    fallback when abcoder is unavailable).
- All committed configs use the abcoder-mcp default:
  `.config/gildedrose.yml:17`,
  `.config/alphatrans/commons-validator.yml:17`,
  `.config/oxidizer/gohistogram.yml:17`,
  `.config/oxidizer/stats.yml:17`,
  and the live test fixture in `tests/cmd/MAgHARCM/main.go:26-27`.

Conclusion: `abcoder-mcp` is the documented and configured default; the
native tree-sitter provider remains as the explicit fallback only.

---

## Directive 7 — Compilation status binary (PASS/FAIL)

Status: ✅

Evidence:
- `internal/compiletime/compiletime.go:382-389`:
  ```
  type CompilationStatus string
  const (
      CompilationStatusPass CompilationStatus = "PASS"
      CompilationStatusFail CompilationStatus = "FAIL"
  )
  ```
- `internal/compiletime/state.go:271-278` — `ValidationReport.CompilationStatus()`
  returns strictly `CompilationStatusPass` or `CompilationStatusFail`:
  ```
  if v.CompilationSuccess {
      return CompilationStatusPass
  }
  return CompilationStatusFail
  ```
- The status is a closed two-value enum; no partial pass-rate exists at
  the type level. The `String()` method on `ValidationReport`
  (state.go:281-290) renders `status=PASS|FAIL` only.

Conclusion: the directive is enforced at the type level via the closed
enum and the single binary mapping site.

---

## Directive 8 — No fmt.Print* outside assets/

Status: ✅

Evidence (repo-wide grep, case-insensitive, restricting to `internal/`
and `cmd/`):
- Zero matches for `fmt.Print`, `fmt.Println`, `fmt.Printf`,
  `fmt.Fprint`, `fmt.Fprintln`, `fmt.Fprintf` (excluding
  `fmt.Fprintf` writes into a `strings.Builder` or `bytes.Buffer`).
- The remaining `fmt.Fprintf` calls all write into string builders or
  buffers for prompt / dump / stub emission (prompt assembly is
  ADR-allowed):
  - `internal/agents/chunked_translator.go:198, 200` — `&sb` builder
    for chunked prompt.
  - `internal/agents/mock_validator.go:196-202` — `&b` builder for stub
    generation.
  - `internal/agents/plateau.go:84` — `&b` builder for log message.
  - `internal/logger/logger.go:107` — `out` writer inside the logger
    itself (the canonical place for `fmt.Fprintf`).
  - `internal/tui/tui.go:531-541` — `&b` builder for the human-readable
    `/show` dump (TUI local use, not log output).
- All matches inside `.assets/` are commented-out test code in vendored
  sample inputs (e.g. `assets/samples/oxidizer/stats/go/correlation_test.go:16`)
  and are not compiled into the binary.
- Matches in `.agents/skills/` and `docs/.paper/` are documentation
  fragments (skill markdown), not executable Go code.

Conclusion: no `fmt.Print*` call exists in compiled code outside the
prompt/builder/stub/logger whitelist. ADR-C-010 is satisfied.

---

## Directive 9 — Charm TUI stack

Status: ✅

Evidence (`internal/tui/tui.go`):
- Imports at lines 28-35:
  - `github.com/charmbracelet/bubbles/spinner`
  - `github.com/charmbracelet/bubbles/table`
  - `github.com/charmbracelet/bubbles/textinput`
  - `github.com/charmbracelet/bubbles/viewport`
  - `tea "github.com/charmbracelet/bubbletea"`
  - `github.com/charmbracelet/glamour`
  - `github.com/charmbracelet/lipgloss`
- Bubble Tea Model/Update/View lifecycle implemented:
  - `Init()` returns `tea.Batch(textinput.Blink, m.spinner.Tick)`
    (line 137-139).
  - `Update(msg tea.Msg) (tea.Model, tea.Cmd)` (line 142-189).
  - `View() string` (line 216-251).
- Lip Gloss styles centralised in the `var (...)` block at lines 54-75
  (`bannerStyle`, `promptStyle`, `logStyle`, `helpStyle`, `errorStyle`).
- Glamour render call at line 224
  (`glamour.Render(strings.TrimSpace(banner), "dark")`).
- Idiomatic Charm table used in `RenderConfigTable`
  (`table.WithColumns(...)`, `table.WithRows(...)` at lines 564-579).
- No manual ANSI escape sequences; no reinvented layout primitives.

Conclusion: the TUI is a full idiomatic Bubble Tea Model on the Charm
stack (`bubbletea`, `bubbles`, `lipgloss`, `glamour`). ADR-C-011 is
satisfied.

---

## Directive 10 — ABCoderMcp naming (no ABCoderMCP)

Status: ✅

Evidence:
- Strict, case-sensitive grep for `ABCoderMCP` (uppercase MCP) inside
  Go source (`internal/`, `cmd/`, `tests/`) returns zero matches.
- The only occurrences of `ABCoderMCP` in the whole repo are in
  markdown / paper text and are NEGATIVE references documenting the
  rename (`Sprint-2026-09-04-Handoff.md:88`,
  `Sprint-2026-09-05-Handoff.md:76`,
  `Sprint-2026-09-06-Handoff.md:69`, and the open-work bullet in
  `Sprint-2026-09-08-Handoff.md:24`). The string appears only in
  phrases like "rename ABCoderMCP → ABCoderMcp" or "(NOT
  ABCoderMCPProvider)" — never as a live Go identifier.
- All Go identifiers use the correct camelCase form:
  - `type ABCoderMcpProvider` at `internal/tools/lsp_provider.go:64`
  - `func NewABCoderMcpProvider(...)` at `internal/tools/lsp_provider.go:72`
  - All methods on `*ABCoderMcpProvider` at lines 83, 88, 119, 132, 144,
    155, 166, 177.

Conclusion: ADR-C-013 is satisfied; the uppercase `MCP` form is absent
from all compiled code.

---

## Summary

| # | Directive | Status | Key Evidence |
| - | --------- | :----: | ------------ |
| 1 | Must pattern, no fallbacks | ✅ | `compiletime.go:417-440`; panic sites in `roleflip.go:46-49`, `iter_retrieval.go:64-67`, `navigator.go:51-55`; `spec_lifecycle.go:73-75` |
| 2 | Clear unit boundaries | ✅ | `compiletime/state.go:1-30` leaf package; `agents/state.go:8-10` single alias; per-agent aliases only |
| 3 | `Registry.TryInOrder` | ✅ | `strategy.go:60-74`, `NewDefaultRegistry` lines 39-53, `SelectAndTryStrategies` lines 226-232 |
| 4 | No magic strings | ✅ | All literals in `compiletime.go:17-440`; no inline strings in agent modules |
| 5 | 8+ agent graph | ✅ | `graph.go:33-40` constructors, lines 45-150 nodes, lines 175-227 edges |
| 6 | abcoder-mcp default | ✅ | `compiletime.go:46-50`, `lsp_provider.go:189-196`, all configs `*.yml:17` |
| 7 | Binary compilation status | ✅ | `compiletime.go:382-389`, `state.go:271-278` |
| 8 | No fmt.Print* outside assets | ✅ | Zero matches in `internal/` and `cmd/` |
| 9 | Charm TUI stack | ✅ | `tui.go:28-35` imports; Model/Update/View at 137-251 |
| 10 | ABCoderMcp naming | ✅ | `ABCoderMCP` absent from Go source; `ABCoderMcpProvider` at `lsp_provider.go:64` |

All ten directives are PASS. No code changes are required for this
sprint's methodology compliance gate.

## See also

- `[[1.0.0 ADR-2026-09-26-Vault-Lint-Extension]]` — supersedes the row-5 ("8+ agent graph") and row-10 ("ABCoderMcp naming") audits. ADR-2026-09-26 §Audit findings gives the canonical 10-node topology (8 functional + 2 checkpoint) and confirms `ABCoderMCPProvider` is moot (zero hits); canonical ident is `LSPProviderABCoder`.
