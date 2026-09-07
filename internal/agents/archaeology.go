// Backlink: [[Primitives]] §NEW-PRIM-14 (Software-Archaeology Stage).
//
// Derived from: pp-besm dev.to "Software Archaeology", AgentPatterns.ai
// Legacy Code Archaeology, Rajlich, Müller, Foltz, Baldwin & Clark, Kazman
// & Cai. Emits an compiletime.ArchaeologyReport: boundaries, time-capsule reproduction
// commands, git churn hotspots, legacy naming findings, co-occurrence
// concept map.
package agents

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"MAgHARCM/internal/compiletime"
 	"MAgHARCM/internal/logger"
)

// Type aliases (cycle-free Locality of Behaviour).
// Producer files declare the algorithms; canonical artifact types
// live in internal/compiletime/state.go. Aliases let method receivers
// reference the type name without package qualification.
type ArchaeologyReport = compiletime.ArchaeologyReport
type NamingFinding = compiletime.NamingFinding

const (
	churnHotspotLimit         = 25
	archaeologyGitLogTimeout   = 30 * time.Second
)


// compiletime.NamingFinding records a legacy-identifier forensic finding produced by
type Archaeologist struct{}

func NewArchaeologist() *Archaeologist { return &Archaeologist{} }

func (a *Archaeologist) Investigate(ctx context.Context, sourceDir string) (compiletime.ArchaeologyReport, error) {
	logger.LogAgent("Archaeology", "starting pass for %s", sourceDir)
	if err := ctx.Err(); err != nil {
		return compiletime.ArchaeologyReport{}, err
	}
	b, _ := a.ExtractBoundaries(ctx, sourceDir)
	t, _ := a.BuildTimeCapsule(ctx, sourceDir)
	c, _ := a.FindChurnHotspots(ctx, sourceDir)
	n, _ := a.ForensicNaming(ctx, sourceDir)
	m, _ := a.MapConcepts(ctx, sourceDir)
	logger.LogStep("archaeology: %d/%d/%d/%d/%d", len(b), len(t), len(c), len(n), len(m))
	return compiletime.ArchaeologyReport{BoundaryMap: b, TimeCapsuleCommands: t, ChurnHotspots: c, NamingForensics: n, ConceptMap: m}, nil
}

// ExtractBoundaries reports module, file and function boundaries.
func (a *Archaeologist) ExtractBoundaries(ctx context.Context, sourceDir string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root := cleanRoot(sourceDir)
	var out []string
	for _, m := range []string{"go.mod", "Cargo.toml", "package.json", "pom.xml", "build.gradle", "Makefile", "CMakeLists.txt"} {
		if _, err := os.Stat(filepath.Join(root, m)); err == nil {
			out = append(out, "module: "+filepath.Join(root, m))
		}
	}
	walkSources(ctx, root, func(path string, lines []string) {
		out = append(out, "file: "+path)
		if strings.HasSuffix(path, ".go") {
			for i, line := range lines {
				if m := regexp.MustCompile(`^func\s+([A-Z][A-Za-z0-9]*)\s*\(`).FindStringSubmatch(line); m != nil {
					rel, _ := filepath.Rel(root, path)
					out = append(out, "func: "+rel+":"+strconv.Itoa(i+1)+" "+m[1])
				}
			}
		}
	})
	return out, nil
}

// BuildTimeCapsule emits shell commands that reproduce the original build /
// test environment.
func (a *Archaeologist) BuildTimeCapsule(ctx context.Context, sourceDir string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root, q := cleanRoot(sourceDir), shellQuote
	// manifest → (label, build, test) reproduction lines; test may be empty.
	caps := map[string][3]string{
		"Makefile":     {"Makefile", "make -C " + q(root) + " build || true", "make -C " + q(root) + " test  || true"},
		"Cargo.toml":   {"Rust", "(cd " + q(root) + " && cargo build)", "(cd " + q(root) + " && cargo test)"},
		"go.mod":       {"Go", "(cd " + q(root) + " && go build ./...)", "(cd " + q(root) + " && go test ./...)"},
		"package.json": {"Node", "(cd " + q(root) + " && npm ci || npm install)", "(cd " + q(root) + " && npm test --silent)"},
		"pom.xml":      {"Maven", "(cd " + q(root) + " && mvn -B -q test)", ""},
		"Dockerfile":   {"Container", "docker build -t archaeology-target " + q(root), ""},
	}
	var cmds []string
	for m, c := range caps {
		if _, err := os.Stat(filepath.Join(root, m)); err == nil {
			cmds = append(cmds, "# "+c[0], c[1])
			if c[2] != "" {
				cmds = append(cmds, c[2])
			}
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "*.tf")); len(matches) > 0 {
		cmds = append(cmds, "# Terraform",
			"(cd "+q(root)+" && terraform init -input=false)",
			"(cd "+q(root)+" && terraform plan -input=false)")
	}
	if _, err := os.Stat(filepath.Join(root, "README.md")); err == nil {
		cmds = append(cmds, "# README",
			"head -n 80 "+q(filepath.Join(root, "README.md")))
	}
	return cmds, nil
}

// FindChurnHotspots returns top-N churned files from git history. No git
// → empty slice. Uses exec.CommandContext (same pattern as tools/exec.go).
func (a *Archaeologist) FindChurnHotspots(ctx context.Context, sourceDir string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root := cleanRoot(sourceDir)
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		return []string{}, nil
	}
	logCtx, cancel := context.WithTimeout(ctx, archaeologyGitLogTimeout)
	defer cancel()
	cmd := exec.CommandContext(logCtx, "git", "log", "--pretty=format:", "--name-only")
	cmd.Dir = root
	raw, err := cmd.Output()
	if err != nil {
		logger.LogWarning("archaeology: git log failed: %v", err)
		return []string{}, nil
	}
	counts := map[string]int{}
	for sc := bufio.NewScanner(strings.NewReader(string(raw))); sc.Scan(); {
		if f := strings.TrimSpace(sc.Text()); f != "" {
			counts[f]++
		}
	}
	ranked := make([]string, 0, len(counts))
	for f, c := range counts {
		ranked = append(ranked, f+"\t"+strconv.Itoa(c))
	}
	sort.Slice(ranked, func(i, j int) bool { return churnCount(ranked[i]) > churnCount(ranked[j]) })
	if len(ranked) > churnHotspotLimit {
		ranked = ranked[:churnHotspotLimit]
	}
	return ranked, nil
}

// ForensicNaming scans source files for legacy identifier encodings:
// Hungarian prefixes, leading / trailing underscores, ALL_CAPS, m_ members.
func (a *Archaeologist) ForensicNaming(ctx context.Context, sourceDir string) ([]compiletime.NamingFinding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root := cleanRoot(sourceDir)
	patterns := []struct{ style string; re *regexp.Regexp }{
		{"hungarian", regexp.MustCompile(`\b(?:str|i|b|f|d|n|p|o|c|s)([A-Z][A-Za-z0-9]*)\b`)},
		{"leading_underscore", regexp.MustCompile(`\b_[A-Za-z][A-Za-z0-9]*\b`)},
		{"trailing_underscore", regexp.MustCompile(`\b[A-Za-z][A-Za-z0-9]*_\b`)},
		{"all_caps", regexp.MustCompile(`\b[A-Z]{3,}\b`)},
		{"mfc_member", regexp.MustCompile(`\bm_[A-Za-z][A-Za-z0-9]*\b`)},
	}
	var findings []compiletime.NamingFinding
	walkSources(ctx, root, func(path string, lines []string) {
		for i, line := range lines {
			for _, p := range patterns {
				if loc := p.re.FindStringIndex(line); loc != nil {
					tok := line[loc[0]:loc[1]]
					findings = append(findings, compiletime.NamingFinding{
						Style: p.style, File: path, Line: i + 1,
						Token: tok, Suggestion: suggestGoName(p.style, tok),
					})
				}
			}
		}
	})
	return findings, nil
}

// MapConcepts clusters files via union-find on shared top-level identifiers.
func (a *Archaeologist) MapConcepts(ctx context.Context, sourceDir string) (map[string][]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root := cleanRoot(sourceDir)
	identRE := regexp.MustCompile(`\b[A-Z][A-Za-z0-9]{2,}\b|\bfunc\s+([A-Za-z][A-Za-z0-9]*)\b|\btype\s+([A-Za-z][A-Za-z0-9]*)\b`)
	idPerFile, parent, fileID := map[string]map[string]int{}, map[string]string{}, map[string]string{}
	walkSources(ctx, root, func(path string, lines []string) {
		ids := map[string]int{}
		for _, line := range lines {
		for _, m := range identRE.FindAllStringSubmatch(line, -1) {
			switch {
			case len(m) > 1 && m[1] != "":
				ids[m[1]]++
			case len(m) > 2 && m[2] != "":
				ids[m[2]]++
			default:
				ids[m[0]]++
			}
		}
		}
		if len(ids) > 0 {
			idPerFile[path] = ids
		}
	})
	var find func(string) string
	find = func(x string) string {
		if _, ok := parent[x]; !ok {
			parent[x] = x
		}
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	for f, ids := range idPerFile {
		for id := range ids {
			if other, ok := fileID[id]; ok {
				if ra, rb := find(f), find(other); ra != rb {
					parent[ra] = rb
				}
			} else {
				fileID[id] = f
			}
		}
	}
	clusters := map[string][]string{}
	for f := range idPerFile {
		clusters[find(f)] = append(clusters[find(f)], f)
	}
	out := map[string][]string{}
	for rep, files := range clusters {
		best, bestN := "", -1
		tally := map[string]int{}
		for _, f := range files {
			for id, c := range idPerFile[f] {
				tally[id] += c
				if tally[id] > bestN {
					best, bestN = id, tally[id]
				}
			}
		}
		if best == "" {
			best = filepath.Base(rep)
		}
		sort.Strings(files)
		out["concept:"+strings.ToLower(best)] = files
	}
	return out, nil
}

func cleanRoot(dir string) string {
	if r := filepath.Clean(dir); r != "" {
		return r
	}
	return "."
}

func isSourceFile(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".go", ".rs", ".py", ".js", ".ts", ".tsx", ".jsx", ".java", ".c", ".cc", ".cpp", ".h", ".hpp":
		return true
	}
	return false
}

// walkSources invokes fn(path, lines) for every source file under root.
func walkSources(ctx context.Context, root string, fn func(path string, lines []string)) {
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.IsDir() {
				for _, skip := range compiletime.ArchaeologySkipDirs {
					if d.Name() == skip {
						return filepath.SkipDir
					}
				}
			}
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if !isSourceFile(d.Name()) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		var lines []string
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		f.Close()
		fn(path, lines)
		return nil
	})
}

func suggestGoName(style, tok string) string {
	switch style {
	case "hungarian":
		return tok[1:]
	case "leading_underscore":
		return tok[1:]
	case "trailing_underscore":
		return strings.TrimSuffix(tok, "_")
	case "all_caps":
		return strings.ToLower(tok)
	case "mfc_member":
		return tok[2:]
	}
	return tok
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func churnCount(line string) int {
	i := strings.LastIndexByte(line, '\t')
	if i < 0 {
		return 0
	}
	n, err := strconv.Atoi(line[i+1:])
	if err != nil {
		return 0
	}
	return n
}
