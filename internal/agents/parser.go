package agents

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Matches headers like "### File: src/lib.rs", "**FILE:** `src/lib.rs`", "File: `src/lib.rs`", "Path: src/main.rs"
	fileHeaderRegex = regexp.MustCompile(`(?i)^(?:[#*>\s-]*)\b(?:FILE|PATH)\b[:\s*]*[` + "`" + `"'({[]*([a-zA-Z0-9_./\\-]+(?:\.[a-zA-Z0-9_]+)?)[` + "`" + `"'\s)}\]]*`)

	// Matches code fence annotations like "```rust src/lib.rs", "```rust:src/lib.rs", "```rust file="src/lib.rs""
	fenceAnnotationRegex = regexp.MustCompile(`^` + "```" + `[a-zA-Z0-9_-]*(?:[:\s]+(?:file=)?|\s+)[` + "`" + `"'({[]*([a-zA-Z0-9_./\\-]+(?:\.[a-zA-Z0-9_]+)?)[` + "`" + `"'\s)}\]]*`)

	// Matches comment markers like "// File: src/lib.rs", "/* FILE: src/lib.rs */", "# File: src/lib.rs"
	commentFileRegex = regexp.MustCompile(`^(?:\/{2,}|#|\/\*|--)\s*\b(?:FILE|PATH)\b[:\s*]*[` + "`" + `"'({[]*([a-zA-Z0-9_./\\-]+(?:\.[a-zA-Z0-9_]+)?)[` + "`" + `"'\s)}\]]*`)
)

// cleanFilePath sanitizes and normalizes a relative file path extracted from model output.
func cleanFilePath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, "`'\"*:#()[];{}<>")
	p = strings.TrimSpace(p)
	if idx := strings.IndexAny(p, " \t\r\n("); idx != -1 {
		p = p[:idx]
	}
	p = strings.Trim(p, "`'\"*:#()[];{}<>")
	p = strings.ReplaceAll(p, "\\", "/")
	return filepath.Clean(p)
}

func isValidFilePath(p string) bool {
	if p == "" || p == "." || p == "/" || p == "..." {
		return false
	}
	return strings.Contains(p, "/") || strings.Contains(p, ".")
}

// parseAllFileMarkers parses all file blocks and their contents from model markdown output.
func parseAllFileMarkers(text string) map[string]string {
	files := make(map[string]string)
	lines := strings.Split(text, "\n")

	var currentFile string
	var currentLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 1. Check for code fence annotation: ```rust src/lib.rs
		if strings.HasPrefix(trimmed, "```") {
			if match := fenceAnnotationRegex.FindStringSubmatch(trimmed); len(match) > 1 {
				filePath := cleanFilePath(match[1])
				if isValidFilePath(filePath) {
					if currentFile != "" {
						files[currentFile] = extractCleanCode(currentLines)
						currentLines = nil
					}
					currentFile = filePath
					// Store generic opening fence so extractCleanCode correctly captures fenced lines
					currentLines = append(currentLines, "```")
					continue
				}
			}
		}

		// 2. Check for explicit file header marker: FILE: src/lib.rs
		if match := fileHeaderRegex.FindStringSubmatch(trimmed); len(match) > 1 {
			filePath := cleanFilePath(match[1])
			if isValidFilePath(filePath) {
				if currentFile != "" {
					files[currentFile] = extractCleanCode(currentLines)
					currentLines = nil
				}
				currentFile = filePath
				continue
			}
		}

		// 3. Check for inline comment marker: // FILE: src/lib.rs
		if match := commentFileRegex.FindStringSubmatch(trimmed); len(match) > 1 {
			filePath := cleanFilePath(match[1])
			if isValidFilePath(filePath) {
				if currentFile != "" {
					files[currentFile] = extractCleanCode(currentLines)
					currentLines = nil
				}
				currentFile = filePath
				continue
			}
		}

		if currentFile != "" {
			currentLines = append(currentLines, line)
		}
	}

	if currentFile != "" {
		files[currentFile] = extractCleanCode(currentLines)
	}

	return files
}

func extractCleanCode(lines []string) string {
	var codeLines []string
	inFence := false
	sawFence := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			if !inFence {
				inFence = true
				sawFence = true
				continue
			} else {
				inFence = false
				continue
			}
		}
		if inFence {
			codeLines = append(codeLines, line)
		} else if !sawFence {
			codeLines = append(codeLines, line)
		}
	}

	if sawFence && len(codeLines) > 0 {
		return strings.TrimSpace(strings.Join(codeLines, "\n"))
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func parseFileBlocks(text, sectionHeader string) map[string]string {
	sectionText := text
	if sectionHeader != "" {
		idx := strings.Index(text, sectionHeader)
		if idx != -1 {
			sectionText = text[idx+len(sectionHeader):]
		}
	}
	return parseAllFileMarkers(sectionText)
}
