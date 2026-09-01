// magtui is the interactive user-facing front-end for MAgHARCM.
//
// Phase 1: collect a translation request via prompts and write a YAML
//
//	config (the same schema cmd/MAgHARCM consumes).
//
// Phase 2: execute the pipeline described by the in-memory YAML, reusing
//
//	internal/runner.
//
// `/clear` returns to phase 1 without restarting the process. Other slash
// commands mirror the small, conventional set found in agent REPLs
// (agy, claude): /help /show /save /load /samples /status /run /dry-run /quit.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"MAgHARCM/internal/config"
	"MAgHARCM/internal/runner"
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

func main() {
	cfg := *config.Defaults()
	phase := phaseCollect

	fmt.Print(banner)
	for {
		fmt.Printf("\n[%s] > ", phaseLabel(phase))
		line, err := readLine()
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return
			}
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/") {
			next, cont, err := handleSlash(line, &cfg, phase)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				continue
			}
			if !cont {
				return
			}
			phase = next
			continue
		}

		if phase == phaseCollect {
			if err := phase1Step(line, &cfg); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			continue
		}

		fmt.Println("Phase 2 is read-only. Use /clear to revisit the YAML, or /help for commands.")
	}
}

type phase int

const (
	phaseCollect phase = iota
	phaseExecute
)

func phaseLabel(p phase) string {
	if p == phaseExecute {
		return "phase 2 · execute"
	}
	return "phase 1 · collect"
}

// readLine returns one line from stdin (stripped of trailing newline).
// readLine is a thin wrapper around the shared package-level scanner so
// we don't allocate a fresh *bufio.Scanner on every prompt — each
// Scanner buffers ahead, and creating one per call drops buffered bytes.
func readLine() (string, error) {
	if !stdinScanner.Scan() {
		if err := stdinScanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return stdinScanner.Text(), nil
}

var stdinScanner = func() *bufio.Scanner {
	s := bufio.NewScanner(os.Stdin)
	s.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return s
}()
// phase1Step walks through the YAML fields. The user types "ok" to accept
// the default in brackets, or types a value; an empty line also accepts
// the default. Typing "?" prints help for the current field; "abort"
// returns to the field-picker without changing anything.
func phase1Step(input string, cfg *config.Config) error {
	switch input {
	case "?", "help":
		fmt.Println(phase1Help)
		return nil
	case "done", "next", "finish":
		return phase1Finalize(cfg)
	case "abort", "back", "cancel":
		fmt.Println("(aborted current edit; pick a field or type ? for help)")
		return nil
	}

	if eq := strings.SplitN(input, "=", 2); len(eq) == 2 {
		key := strings.TrimSpace(eq[0])
		val := strings.TrimSpace(eq[1])
		if err := setField(cfg, key, val); err != nil {
			return err
		}
		fmt.Printf("  -> %s = %s\n", key, val)
		return nil
	}

	fmt.Println(`unrecognised. type a field = value, "?" for help, or "done" to finish.`)
	return nil
}

// setField updates one cfg field by its short name; matches the keys
// printed in phase1Help.
func setField(cfg *config.Config, key, val string) error {
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

// phase1Finalize validates the collected config, writes it to a YAML file
// (the default path is `./magharcm-request.yml`), and returns the next
// phase. The caller (main) updates its own phase variable.
func phase1Finalize(cfg *config.Config) error {
	if cfg.SourceDir == "" || cfg.TargetDir == "" {
		return fmt.Errorf("source_dir and target_dir are required before finishing phase 1")
	}
	if cfg.SourceLang == "" || cfg.TargetLang == "" {
		return fmt.Errorf("source_language and target_language are required")
	}
	out := defaultRequestPath()
	if err := writeYAML(out, cfg); err != nil {
		return fmt.Errorf("write %s: %w", out, err)
	}
	fmt.Printf("Wrote %s\n", out)
	fmt.Println("Phase 1 complete. Type /run to execute, /clear to revise, or /show to reprint.")
	return nil
}

func defaultRequestPath() string {
	if v := os.Getenv("MAGHARCM_REQUEST_PATH"); v != "" {
		return v
	}
	return "magharcm-request.yml"
}

// writeYAML emits the nested translation.* schema that config.ParseYAML
// expects. Keys mirror config.sample1.yml and the live configs.
func writeYAML(path string, cfg *config.Config) error {
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

// handleSlash processes /-prefixed commands. Returns the next phase, a
// continuation flag, and any error.
func handleSlash(line string, cfg *config.Config, current phase) (phase, bool, error) {
	args := strings.Fields(line)
	cmd := strings.ToLower(args[0])

	switch cmd {
	case "/help", "/?":
		fmt.Println(strings.TrimSpace(slashHelp))
		return current, true, nil

	case "/clear":
		*cfg = *config.Defaults()
		fmt.Println("Cleared. Back to phase 1 — answer questions to build a new YAML.")
		return phaseCollect, true, nil

	case "/show":
		fmt.Println(dumpYAML(cfg))
		return current, true, nil

	case "/save":
		path := defaultRequestPath()
		if len(args) > 1 {
			path = args[1]
		}
		if err := writeYAML(path, cfg); err != nil {
			return current, true, err
		}
		fmt.Printf("Saved to %s\n", path)
		return current, true, nil

	case "/load":
		if len(args) < 2 {
			return current, true, fmt.Errorf("usage: /load <path>")
		}
		loaded, err := config.LoadYAML(args[1])
		if err != nil {
			return current, true, fmt.Errorf("load %s: %w", args[1], err)
		}
		if loaded.SourceDir == "" || loaded.TargetDir == "" {
			return current, true, fmt.Errorf("refusing: %s is missing source_dir or target_dir", args[1])
		}
		*cfg = *loaded
		fmt.Printf("Loaded %s. Type /run to execute or /show to inspect.\n", args[1])
		return phaseExecute, true, nil

	case "/samples":
		fmt.Println(samplesList())
		return current, true, nil

	case "/status":
		fmt.Printf("phase=%s\nsource=%s (%s)\ntarget=%s (%s)\ntoolchain=%s\nmodels: reasoning=%s, coding=%s\niterations=%d timeout=%ds\n",
			phaseLabel(current), cfg.SourceDir, cfg.SourceLang, cfg.TargetDir, cfg.TargetLang, cfg.Toolchain,
			cfg.ReasoningModel, cfg.CodingModel, cfg.MaxIterations, int(cfg.Timeout.Seconds()))
		return current, true, nil

	case "/run":
		if err := runPhase2(*cfg); err != nil {
			fmt.Fprintf(os.Stderr, "phase 2 failed: %v\n", err)
		}
		return phaseExecute, true, nil

	case "/dry-run":
		fmt.Printf("(dry-run) would execute with: source=%s target=%s reasoning=%s coding=%s iterations=%d\n",
			cfg.SourceDir, cfg.TargetDir, cfg.ReasoningModel, cfg.CodingModel, cfg.MaxIterations)
		return phaseExecute, true, nil

	case "/quit", "/exit":
		fmt.Println("bye.")
		return current, false, nil

	default:
		return current, true, fmt.Errorf("unknown command %q (type /help)", cmd)
	}
}

func runPhase2(cfg config.Config) error {
	if cfg.SourceDir == "" || cfg.TargetDir == "" {
		return fmt.Errorf("refusing: source_dir and target_dir are required (use /load <path> or finish phase 1)")
	}
	final, err := runner.Run(context.Background(), &cfg)
	if err != nil {
		return err
	}
	if runner.Success(final) {
		fmt.Println("phase 2: success. Type /clear to start a new request, /quit to exit.")
	} else {
		fmt.Println("phase 2: incomplete. See validator logs above. Type /clear to revise the YAML.")
	}
	return nil
}

// dumpYAML is a human-readable summary, not the canonical YAML.
func dumpYAML(cfg *config.Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "source_dir        %s\n", cfg.SourceDir)
	fmt.Fprintf(&b, "source_language   %s\n", cfg.SourceLang)
	fmt.Fprintf(&b, "target_dir        %s\n", cfg.TargetDir)
	fmt.Fprintf(&b, "target_language   %s\n", cfg.TargetLang)
	fmt.Fprintf(&b, "toolchain         %s\n", cfg.Toolchain)
	fmt.Fprintf(&b, "reasoning_model   %s\n", cfg.ReasoningModel)
	fmt.Fprintf(&b, "coding_model      %s\n", cfg.CodingModel)
	fmt.Fprintf(&b, "ollama_url        %s\n", cfg.OllamaBaseURL)
	fmt.Fprintf(&b, "max_iterations    %d\n", cfg.MaxIterations)
	fmt.Fprintf(&b, "timeout_seconds   %d\n", int(cfg.Timeout.Seconds()))
	fmt.Fprintf(&b, "lsp_provider      %s\n", cfg.LSPProvider)
	return b.String()
}

// samplesList enumerates the bundled config.sample*.yml files in the
// current working directory. Empty list is fine — magtui is useful even
// without samples.
func samplesList() string {
	matches, err := filepath.Glob("config.sample*.yml")
	if err != nil || len(matches) == 0 {
		return "(no config.sample*.yml files in this directory; /load <path> accepts any YAML.)"
	}
	return strings.Join(matches, "\n")
}

const phase1Help = `phase 1 fields (type a value, or hit enter for the default):
  source_dir       directory holding the source code
  source_language  e.g. C, Go, Java, Rust
  target_dir       directory to write translated code into
  target_language  e.g. Rust, Go
  toolchain        build/test toolchain (cargo, go, maven)
  reasoning_model  Ollama reasoning model tag
  coding_model     Ollama coding model tag
  ollama_url       base URL for the Ollama HTTP API
  max_iterations   repair-loop iterations (1-50)
  timeout_seconds  per-run wall-clock budget
  lsp_provider     "native" (tree-sitter) or "abcoder" (MCP)

type ` + "`field = value`" + ` to jump to a specific field, e.g.
` + "`source_dir = assets/samples/foo`" + `. type ` + "`done`" + `
when finished.`

const slashHelp = `slash commands (always start with /):
  /help                 show this message
  /show                 print the in-memory YAML as a summary
  /save [path]          write the YAML (default: magharcm-request.yml)
  /load <path>          load a YAML and (if valid) jump to phase 2
  /samples              list bundled sample configs in this directory
  /status               print phase + key config fields
  /run                  execute phase 2 with the current YAML
  /dry-run              print what /run would do, but don't execute
  /clear                drop the in-memory YAML and return to phase 1
  /quit                 exit the REPL

phase 1 also accepts:
  field = value         set a YAML field directly (see ` + "`?`" + `)
  ?                     list phase 1 fields
  done                  finish phase 1 and write the YAML`
