# Charm-stack idiomatic audit — `internal/tui/`

Sprint: 2026-09-26  ·  Subagent: G  ·  ADR: ADR-C-011

## Scope

Single file scanned:

- `internal/tui/tui.go` (613 lines after edits)

There is exactly one file in `internal/tui/`; no sub-packages, no test files
inside the directory. The test mirror lives in `tests/cmd/MAgHARCM-tui/` and
is **not** in scope.

## Charm imports in scope

| Import | Used for | Verdict |
|---|---|---|
| `github.com/charmbracelet/bubbletea` | Model / Update / View loop, `tea.Quit`, `tea.Batch`, `tea.KeyMsg`, `tea.WindowSizeMsg`, `tea.NewProgram` | idiomatic |
| `github.com/charmbracelet/bubbles/spinner` | Phase-Execute running indicator (`spinner.Dot`, styled) | idiomatic |
| `github.com/charmbracelet/bubbles/table` | `RenderConfigTable` — `/show` command output | idiomatic |
| `github.com/charmbracelet/bubbles/textinput` | REPL input line + `textinput.Blink` cursor command | idiomatic |
| `github.com/charmbracelet/bubbles/viewport` | **removed — was constructed but never `.View()`-ed** | dead, deleted |
| `github.com/charmbracelet/glamour` | banner rendering (with lipgloss fallback) | idiomatic |
| `github.com/charmbracelet/lipgloss` | styles + borders (`DoubleBorder`, `NormalBorder`, `Color`) | idiomatic |

## Audit checklist results

### 1. Manual ANSI escape sequences (`\033`, `\x1b`, `\u001b`)

Pattern search across `internal/tui/`: **0 matches.** No raw ESC bytes,
no octal escapes, no unicode-escape escapes, no string concatenation
that would emit a CSI sequence. Confirmed via three regex passes
(`\\033`, `\\x1b`, `\\u001b`, plus `SetForeground` / `SetBackground`).

### 2. `SetForeground` / `SetBackground` from outside `lipgloss`

Pattern search: **0 matches.** The only color setters are
`lipgloss.NewStyle().Foreground(lipgloss.Color(...))` and
`table.DefaultStyles()` followed by `.BorderForeground(...)` /
`.Foreground(...)` — all Charm-provided surfaces.

### 3. Raw cursor-positioning or box-drawing characters written to screen

Pattern search for U+2500..U+257F box-drawing block: **0 matches.**
All borders come from `lipgloss.DoubleBorder()` (banner) and
`lipgloss.NormalBorder()` (table header). No manual `┌─┐`/`│`/`└─┘`
characters in source.

### 4. Reimplemented Charm primitives

Searched for: manual spinner loops, manual table cell layouts,
manual text input (read-line + echo), manual list/paginator.

| Candidate | Status |
|---|---|
| Manual spinner | **none.** `spinner.New()` + `spinner.Dot` + `m.spinner.Tick` in `Init()`. |
| Manual table | **none.** `RenderConfigTable` constructs `[]table.Column` and `[]table.Row` then `table.New(table.WithColumns(...))` with `table.DefaultStyles()`. |
| Manual text input | **none.** `textinput.New()` with placeholder, focus, char limit, width. |
| Manual list/paginator | **none.** The log tail is rendered line-by-line through `logStyle.Render(h)`, which is text styling, not a list primitive. `lastNLines(8)` is a slice window over `logger.Snapshot()`, not a UI widget — it has no selection state, no scroll, no keymap. |

## Violations found and fixed

| # | Location | Pattern | Recommended replacement | Applied? |
|---|---|---|---|---|
| 1 | `tui.go:33` (import) | `github.com/charmbracelet/bubbles/viewport` imported but unused (no `.View()` call site) | remove import | yes |
| 2 | `tui.go:106` (model field) | `viewport viewport.Model` field allocated but never read after `Update` set it | remove field | yes |
| 3 | `tui.go:151-157` (`Update`) | `m.viewport = viewport.New(...)` / width / height assignments were the only writes to a dead field | remove the inner `if !m.ready` branch (only the `m.ready = true` flip remains meaningful) | yes |

Edits were a single conceptual cleanup: drop the unused `bubbles/viewport`
import + field + their lone write site. `Update` retains `m.width` /
`m.height` capture so future Charm primitives can size off them, and
`m.ready` stays for parity with future first-window handling.

**Total violations:** 3 (all from the same dead-import site).
**Total replacements applied:** 3 (one import removal, one field removal,
one branch reduction).

## Logger-style TUI output review

`RunCLI()` (lines 261–272) calls:

```go
logger.SetOutput(logger.Tee(os.Stdout))
logger.LogStep("%s", strings.TrimSpace(banner))
```

This writes the banner to stdout **before** `tea.NewProgram(...).Run()`
takes the terminal. Once the program starts, the TUI exclusively uses
`logger.LogStep` / `logger.LogError` to push events into the ring buffer
that `View()` reads via `logger.Snapshot()`. That is the documented
pattern from `internal/logger/logger.go` (`SetOutput` comment: "Production
code uses stdout; tests and the interactive TUI pass other writers.")

**Verdict:** correct, in-place. Not a violation. The logger module is the
canonical sink; the TUI reads from it for `/logs` and for the in-screen
tail. No migration needed.

## Acceptance verification

- [x] Files scanned: `internal/tui/tui.go`.
- [x] ZERO manual ANSI escapes in `internal/tui/`.
- [x] ZERO `SetForeground`/`SetBackground` outside lipgloss in `internal/tui/`.
- [x] ZERO raw box-drawing characters in `internal/tui/`.
- [x] Every TUI module that has a Charm equivalent uses it
      (spinner, table, textinput, glamour, lipgloss).
- [x] Dead `bubbles/viewport` import + field removed.
- [x] `go vet ./internal/tui/...` clean (run, no output).
- [x] No file changes outside `internal/tui/` + this report.

## Result

- Report path: `local://sprint-2026-09-26-charm-audit.md`
- Violation count: **3** (one dead-import site, three write lines).
- Replacement count: **3** (all at the same dead-import site).
- Net lines removed from `internal/tui/tui.go`: 5
  (1 import line, 1 field line, 3 viewport write lines).
