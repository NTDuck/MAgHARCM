// MAgHARCM-tui is the interactive user-facing front-end for MAgHARCM.
//
// Phase 1: collect a translation request via prompts and write a YAML
//
//	config (the same schema cmd/MAgHARCM consumes).
//
// Phase 2: execute the pipeline described by the in-memory YAML, reusing
//
//	internal/runner.
//
// `/clear` returns to Phase 1 without restarting the process. Other slash
// commands mirror the small, conventional set found in agent REPLs
// (agy, claude): /help /show /save /load /samples /status /run /dry-run /quit.
//
// The interactive loop is a Bubble Tea Model/Update/View. The exported
// functions Phase1Step, SetField, Phase1Finalize, HandleSlash, and the
// help string constants are kept as the canonical command surface so
// the test mirror in tests/cmd/MAgHARCM-tui stays in lockstep.
package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"MAgHARCM/internal/config"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/runner"
	"gopkg.in/yaml.v3"
)

const banner = `
================================================================
   MAgHARCM - interactive translation harness
================================================================
   Phase 1: answer a few questions to build a YAML config.
   Phase 2: the harness runs end-to-end and reports the result.
   Slash commands: /help
================================================================
`

// Lip Gloss styles — central colors live here so every screen shares one palette.
var (
	bannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7FB3D5")).
			Border(lipgloss.DoubleBorder()).
			Padding(0, 1)

	promptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F7DC6F"))

	logStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BDC3C7"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#82E0AA")).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")).
			Bold(true)
)

// Phase identifies which screen of the REPL is active.
type Phase int

const (
	PhaseCollect Phase = iota
	PhaseExecute
)

// PhaseLabel returns the human label for a Phase.
func PhaseLabel(p Phase) string {
	if p == PhaseExecute {
		return "Phase 2 · execute"
	}
	return "Phase 1 · collect"
}

// ReplState carries per-session REPL state that persists between commands.
// Log history lives in the logger package's ring buffer (see logger.Snapshot)
// so this struct only holds the debug toggle.
type ReplState struct {
	Debug bool
}

// model is the Bubble Tea Model. Exported fields are not used by tests;
// the REPL is driven through the exported functions (SetField, HandleSlash,
// Phase1Step, Phase1Finalize) so the test mirror stays unchanged.
type model struct {
	cfg      *config.Config
	phase    Phase
	rs       *ReplState
	input    textinput.Model
	viewport viewport.Model
	spinner  spinner.Model
	history  []string
	quitting bool
	width    int
	height   int
	ready    bool
}

func newModel(cfg *config.Config, rs *ReplState) model {
	ti := textinput.New()
	ti.Placeholder = "type ? for help"
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#A3E4D7"))

	return model{
		cfg:     cfg,
		phase:   PhaseCollect,
		rs:      rs,
		input:   ti,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 6
		footerHeight := 3
		vHeight := msg.Height - headerHeight - footerHeight
		if vHeight < 5 {
			vHeight = 5
		}
		if !m.ready {
			m.viewport = viewport.New(msg.Width, vHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = vHeight
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			line := strings.TrimSpace(m.input.Value())
			m.input.SetValue("")
			if line == "" {
				return m, nil
			}
			m.history = append(m.history, line)
			m = m.dispatch(line)
			if m.quitting {
				return m, tea.Quit
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// dispatch routes a single line to slash commands, Phase1 step, or
// the Phase 2 read-only message.
func (m model) dispatch(line string) model {
	if strings.HasPrefix(line, "/") {
		next, cont, err := HandleSlash(line, m.cfg, m.phase, m.rs)
		if err != nil {
			logger.LogError("%v", err)
			return m
		}
		m.phase = next
		if !cont {
			m.quitting = true
		}
		return m
	}
	if m.phase == PhaseCollect {
		if err := Phase1Step(line, m.cfg); err != nil {
			logger.LogError("%v", err)
		}
		return m
	}
	logger.LogStep("Phase 2 is read-only. Use /clear to revisit the YAML, or /help for commands.")
	return m
}

func (m model) View() string {
	if m.quitting {
		return "bye.\n"
	}

	var b strings.Builder

	// Banner.
	renderedBanner, err := glamour.Render(strings.TrimSpace(banner), "dark")
	if err != nil {
		renderedBanner = bannerStyle.Render(strings.TrimSpace(banner))
	}
	b.WriteString(renderedBanner)
	b.WriteString("\n")

	if m.phase == PhaseExecute {
		b.WriteString(promptStyle.Render(fmt.Sprintf("[%s] %s Running pipeline...", PhaseLabel(m.phase), m.spinner.View())))
	} else {
		b.WriteString(promptStyle.Render(fmt.Sprintf("[%s]", PhaseLabel(m.phase))))
	}
	b.WriteString("\n")
	for _, h := range lastNLines(8) {
		b.WriteString(logStyle.Render(h))
		b.WriteString("\n")
	}

	// Input line.
	b.WriteString(promptStyle.Render("> "))
	b.WriteString(m.input.View())
	b.WriteString("\n")

	// Footer hint.
	b.WriteString(helpStyle.Render("type /help, ? for phase 1 fields, done to finalize"))

	return b.String()
}

func lastNLines(n int) []string {
	all := logger.Snapshot()
	if len(all) <= n {
		return all
	}
	return all[len(all)-n:]
}

// RunCLI launches the Bubble Tea program. Tests do NOT call this; they
// drive the exported command surface directly.
func RunCLI() {
	cfg := config.Config{}
	rs := &ReplState{}

	logger.SetOutput(logger.Tee(os.Stdout))
	logger.LogStep("%s", strings.TrimSpace(banner))

	if _, err := tea.NewProgram(newModel(&cfg, rs)).Run(); err != nil {
		logger.LogError("tui: %v", err)
		os.Exit(1)
	}
}

// Phase1Step walks through the YAML fields. The user types "ok" to accept
// the default in brackets, or types a value; an empty line also accepts
// the default. Typing "?" prints help for the current field; "abort"
// returns to the field-picker without changing anything.
func Phase1Step(input string, cfg *config.Config) error {
	switch input {
	case "?", "help":
		logger.LogStep("%s", Phase1Help)
		return nil
	case "done", "next", "finish":
		return Phase1Finalize(cfg)
	case "abort", "back", "cancel":
		logger.LogStep("(aborted current edit; pick a field or type ? for help)")
		return nil
	}

	if eq := strings.SplitN(input, "=", 2); len(eq) == 2 {
		key := strings.TrimSpace(eq[0])
		val := strings.TrimSpace(eq[1])
		if err := SetField(cfg, key, val); err != nil {
			return err
		}
		logger.LogStep("set %s = %s", key, val)
		return nil
	}

	logger.LogStep(`unrecognised. type a field = value, "?" for help, or "done" to finish.`)
	return nil
}

// SetField updates one cfg field by its short name; matches the keys
// printed in Phase1Help.
func SetField(cfg *config.Config, key, val string) error {
	switch strings.ToLower(strings.ReplaceAll(key, "-", "_")) {
	case "source_dir":
		cfg.SourceDir = val
	case "source_language", "source_lang":
		cfg.SourceLang = val
	case "target_dir":
		cfg.TargetDir = val
	case "target_language", "target_lang":
		cfg.TargetLang = val
	case "toolchain":
		cfg.Toolchain = val
	case "reasoning_model", "reasoning":
		cfg.ReasoningModel = val
	case "coding_model", "coding":
		cfg.CodingModel = val
	case "ollama_url", "ollama":
		cfg.OllamaBaseURL = val
	case "max_iterations", "iterations":
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return fmt.Errorf("max_iterations must be a positive integer")
	}
	cfg.MaxIterations = n
	case "timeout_seconds", "timeout":
		n, err := strconv.Atoi(val)
		if err != nil || n <= 0 {
			return fmt.Errorf("timeout_seconds must be a positive integer")
		}
		cfg.Timeout = time.Duration(n) * time.Second
	case "lsp_provider", "lsp":
		cfg.LSPProvider = val
	default:
		return fmt.Errorf("unknown field %q (type ? for the list)", key)
	}
	return nil
}

// Phase1Finalize validates the collected config, writes it to a YAML file
// (the default path is `./magharcm-request.yml`), and returns the next
// Phase. The caller (main) updates its own Phase variable.
func Phase1Finalize(cfg *config.Config) error {
	if cfg.SourceDir == "" || cfg.TargetDir == "" {
		return fmt.Errorf("set source_dir and target_dir before finishing Phase 1")
	}
	if cfg.SourceLang == "" || cfg.TargetLang == "" {
		return fmt.Errorf("set source_language and target_language before finishing Phase 1")
	}
	out := defaultRequestPath()
	if err := WriteYAML(out, cfg); err != nil {
		return fmt.Errorf("write %s: %w", out, err)
	}
	logger.LogStep("Wrote %s", out)
	logger.LogStep("Phase 1 complete. Type /run to execute, /clear to revise, or /show to reprint.")
	return nil
}

func defaultRequestPath() string {
	if v := os.Getenv("MAGHARCM_REQUEST_PATH"); v != "" {
		return v
	}
	return "magharcm-request.yml"
}

// WriteYAML emits the nested translation.* schema that config.ParseYAML
// expects. Keys mirror config.sample1.yml and the live configs.
func WriteYAML(path string, cfg *config.Config) error {
	doc := struct {
		Translation struct {
			Source struct {
				Dir      string `yaml:"dir"`
				Language string `yaml:"language"`
			} `yaml:"source"`
			Target struct {
				Dir       string `yaml:"dir"`
				Language  string `yaml:"language"`
				Toolchain string `yaml:"toolchain"`
			} `yaml:"target"`
			Models struct {
				Reasoning string `yaml:"reasoning"`
				Coding    string `yaml:"coding"`
				OllamaURL string `yaml:"ollama_url"`
			} `yaml:"models"`
			Execution struct {
				MaxIterations  int `yaml:"max_iterations"`
				TimeoutSeconds int `yaml:"timeout_seconds"`
			} `yaml:"execution"`
			LSP struct {
				Provider string `yaml:"provider"`
			} `yaml:"lsp"`
		} `yaml:"translation"`
	}{}
	doc.Translation.Source.Dir = cfg.SourceDir
	doc.Translation.Source.Language = cfg.SourceLang
	doc.Translation.Target.Dir = cfg.TargetDir
	doc.Translation.Target.Language = cfg.TargetLang
	doc.Translation.Target.Toolchain = cfg.Toolchain
	doc.Translation.Models.Reasoning = cfg.ReasoningModel
	doc.Translation.Models.Coding = cfg.CodingModel
	doc.Translation.Models.OllamaURL = cfg.OllamaBaseURL
	doc.Translation.Execution.MaxIterations = cfg.MaxIterations
	doc.Translation.Execution.TimeoutSeconds = int(cfg.Timeout.Seconds())
	doc.Translation.LSP.Provider = cfg.LSPProvider

	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// HandleSlash processes /-prefixed commands. Returns the next Phase, a
// continuation flag, and any error.
func HandleSlash(line string, cfg *config.Config, current Phase, state *ReplState) (Phase, bool, error) {
	args := strings.Fields(line)
	cmd := strings.ToLower(args[0])

	switch cmd {
	case "/help", "/?":
		logger.LogStep("%s", strings.TrimSpace(slashHelp))
		return current, true, nil

	case "/clear":
		*cfg = config.Config{}
		logger.LogStep("Cleared. Back to Phase 1: answer the questions to build a new YAML.")
		return PhaseCollect, true, nil

	case "/show":
		logger.LogStep("%s", RenderConfigTable(cfg))
		return current, true, nil

	case "/save":
		path := defaultRequestPath()
		if len(args) > 1 {
			path = args[1]
		}
		if err := WriteYAML(path, cfg); err != nil {
			return current, true, err
		}
		logger.LogStep("Saved to %s", path)
		return current, true, nil

	case "/load":
		if len(args) < 2 {
			return current, true, fmt.Errorf("usage: /load <path>")
		}
		loaded, err := config.LoadYAML(args[1])
		if err != nil {
			return current, true, fmt.Errorf("load %s: %w", args[1], err)
		}
		*cfg = *loaded
		logger.LogStep("Loaded %s. Type /run to execute or /show to inspect.", args[1])
		return PhaseExecute, true, nil

	case "/samples":
		logger.LogStep("%s", samplesList())
		return current, true, nil

	case "/status":
		logger.LogStep("Phase=%s source=%s (%s) target=%s (%s) toolchain=%s models: reasoning=%s, coding=%s iterations=%d timeout=%ds",
			PhaseLabel(current), cfg.SourceDir, cfg.SourceLang, cfg.TargetDir, cfg.TargetLang, cfg.Toolchain,
			cfg.ReasoningModel, cfg.CodingModel, cfg.MaxIterations, int(cfg.Timeout.Seconds()))
		return current, true, nil

	case "/run":
		if err := runPhase2(*cfg); err != nil {
			logger.LogError("Phase 2 failed: %v", err)
		}
		return PhaseExecute, true, nil

	case "/dry-run":
		logger.LogStep("(dry-run) execute: source=%s target=%s reasoning=%s coding=%s iterations=%d",
			cfg.SourceDir, cfg.TargetDir, cfg.ReasoningModel, cfg.CodingModel, cfg.MaxIterations)
		return PhaseExecute, true, nil

	case "/debug":
		state.Debug = !state.Debug
		if state.Debug {
			logger.LogStep("debug: ON (verbose stderr logging)")
		} else {
			logger.LogStep("debug: OFF")
		}
		return current, true, nil

	case "/logs":
		lines := logger.Snapshot()
		if len(lines) == 0 {
			logger.LogStep("(no log lines captured yet)")
		} else {
			for _, l := range lines {
				logger.LogStep("%s", l)
			}
		}
		return current, true, nil

	case "/quit", "/exit":
		logger.LogStep("bye.")
		return current, false, nil

	default:
		return current, true, fmt.Errorf("unknown command %q (type /help)", cmd)
	}
}

func runPhase2(cfg config.Config) error {
	if cfg.SourceDir == "" || cfg.TargetDir == "" {
		return fmt.Errorf("set source_dir and target_dir first (use /load <path> or finish Phase 1)")
	}
	final, err := runner.Run(context.Background(), &cfg)
	if err != nil {
		return err
	}
	if runner.Success(final) {
		logger.LogStep("Phase 2: success. Type /clear to start a new request, /quit to exit.")
	} else {
		logger.LogStep("Phase 2: incomplete. See validator logs above. Type /clear to revise the YAML.")
	}
	return nil
}


// RenderConfigTable formats the in-memory configuration as an idiomatic Charm table.
func RenderConfigTable(cfg *config.Config) string {
	columns := []table.Column{
		{Title: "Configuration Field", Width: 24},
		{Title: "Assigned Value", Width: 46},
	}
	rows := []table.Row{
		{"source_dir", cfg.SourceDir},
		{"source_language", cfg.SourceLang},
		{"target_dir", cfg.TargetDir},
		{"target_language", cfg.TargetLang},
		{"toolchain", cfg.Toolchain},
		{"reasoning_model", cfg.ReasoningModel},
		{"coding_model", cfg.CodingModel},
		{"ollama_url", cfg.OllamaBaseURL},
		{"max_iterations", fmt.Sprintf("%d", cfg.MaxIterations)},
		{"timeout_seconds", fmt.Sprintf("%d", int(cfg.Timeout.Seconds()))},
		{"lsp_provider", cfg.LSPProvider},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(13),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7FB3D5")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#F7DC6F"))
	t.SetStyles(s)
	return t.View()
}

// globSamples wraps filepath.Glob so samplesList stays declarative.
func globSamples() ([]string, error) {
	matches, err := filepath.Glob("config.sample*.yml")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// samplesList returns the names of every config.sample*.yml file in the
// current working directory. Empty list is fine — MAgHARCM-tui is useful even
// without samples.
func samplesList() string {
	matches, err := globSamples()
	if err != nil || len(matches) == 0 {
		return "(no config.sample*.yml files in this directory; /load <path> accepts any YAML.)"
	}
	return strings.Join(matches, "\n")
}


const Phase1Help = `Phase 1 fields (type a value, or hit enter to accept the default):
  source_dir       directory that holds the source code
  source_language  one of C, Go, Java, Rust
  target_dir       directory to write the translated code into
  target_language  one of Rust, Go
  toolchain        build or test toolchain (cargo, go, maven)
  reasoning_model  Ollama reasoning model tag
  coding_model     Ollama coding model tag
  ollama_url       base URL for the Ollama HTTP API
  max_iterations   repair-loop iterations (1 to 50)
  timeout_seconds  per-run wall-clock budget
  lsp_provider     "native" (tree-sitter) or "abcoder" (MCP)

Type ` + "`field = value`" + ` to set a field. Example: ` + "`source_dir = assets/samples/foo`" + `.
Type ` + "`done`" + ` when you finish.`

const slashHelp = `slash commands (always start with /):
  /help                 show this message
  /show                 print the in-memory YAML as a summary
  /save [path]          write the YAML (default: magharcm-request.yml)
  /load <path>          load a YAML and (if valid) jump to Phase 2
  /samples              list bundled sample configs in this directory
  /status               print Phase + key config fields
  /run                  execute Phase 2 with the current YAML
  /dry-run              print what /run would do, but do not execute
  /debug                toggle verbose debug output
  /logs                 print the in-memory ring buffer of recent log lines
  /clear
  /quit                 exit the REPL

Phase 1 also accepts:
  field = value         set a YAML field directly (see ` + "`?`" + `)
  ?                     list Phase 1 fields
  done                  finish Phase 1 and write the YAML`
