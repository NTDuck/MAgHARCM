package languages

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// CodeElement represents a structural code item extracted into Canonical Code IR.
type CodeElement struct {
	Kind      string `json:"kind"` // "function", "struct", "class", "interface", "trait", "type", "enum", "macro", "impl"
	Name      string `json:"name"`
	Signature string `json:"signature,omitempty"`
	Line      int    `json:"line"`
	EndLine   int    `json:"end_line"`
}

// FileStructureResult represents the complete structural decomposition of a source file.
type FileStructureResult struct {
	FilePath string        `json:"file_path"`
	Language string        `json:"language"`
	Elements []CodeElement `json:"elements"`
	Imports  []string      `json:"imports"`
	RawCode  string        `json:"raw_code,omitempty"`
}

// ExtractFileStructure parses any source file and extracts its canonical Code IR and imports.
func ExtractFileStructure(filePath string) (*FileStructureResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file %s: %w", filePath, err)
	}

	registry := GetRegistry()
	spec, found := registry.FindByPath(filePath)
	if !found {
		// Generic fallback if extension is unknown
		spec = &LanguageSpec{
			Name:          "generic",
			CommentPrefix: "//",
		}
	}

	// 1. Attempt dynamic Tree-Sitter parsing if grammar is available
	loader := GetLoader()
	if tsLang, err := loader.LoadGrammar(spec); err == nil && tsLang != nil {
		if result, err := parseWithTreeSitter(filePath, content, spec, tsLang); err == nil {
			return result, nil
		}
	}

	// 2. Language-aware structural and lexical parser fallback
	result := parseWithLanguageHeuristics(filePath, content, spec)
	return result, nil
}

// parseWithTreeSitter extracts Code IR elements using dynamic Tree-Sitter grammar and spec mappings.
func parseWithTreeSitter(filePath string, src []byte, spec *LanguageSpec, tsLang *tree_sitter.Language) (*FileStructureResult, error) {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(tsLang); err != nil {
		return nil, err
	}

	tree := parser.Parse(src, nil)
	if tree == nil {
		return nil, fmt.Errorf("tree-sitter parse produced nil tree for %s", filePath)
	}
	defer tree.Close()

	var elements []CodeElement
	var imports []string

	root := tree.RootNode()
	if root != nil {
		count := root.NamedChildCount()
		for i := uint(0); i < count; i++ {
			node := root.NamedChild(i)
			if node == nil {
				continue
			}

			nodeKind := node.Kind()
			startRow := int(node.StartPosition().Row) + 1
			endRow := int(node.EndPosition().Row) + 1

			category, recognized := spec.ASTMapping.CategorizeNode(nodeKind)
			if !recognized {
				continue
			}

			if category == CategoryImport {
				imports = append(imports, strings.TrimSpace(node.Utf8Text(src)))
				continue
			}

			name := ""
			if nameNode := node.ChildByFieldName("name"); nameNode != nil {
				name = nameNode.Utf8Text(src)
			}
			if name == "" {
				name = firstLine(node.Utf8Text(src))
				if idx := strings.Index(name, "{"); idx > 0 {
					name = strings.TrimSpace(name[:idx])
				}
				if idx := strings.Index(name, "("); idx > 0 && category == CategoryFunction {
					name = strings.TrimSpace(name[:idx])
				}
			}

			elements = append(elements, CodeElement{
				Kind:      string(category),
				Name:      name,
				Signature: firstLine(node.Utf8Text(src)),
				Line:      startRow,
				EndLine:   endRow,
			})
		}
	}

	return &FileStructureResult{
		FilePath: filePath,
		Language: spec.Name,
		Elements: elements,
		Imports:  imports,
		RawCode:  string(src),
	}, nil
}

// parseWithLanguageHeuristics extracts structural declarations using language spec rules.
func parseWithLanguageHeuristics(filePath string, src []byte, spec *LanguageSpec) *FileStructureResult {
	var elements []CodeElement
	var imports []string

	scanner := bufio.NewScanner(bytes.NewReader(src))
	lineNum := 0

	fnRe := regexp.MustCompile(`(?i)(?:pub\s+)?(?:async\s+)?(?:fn|def|func|function|sub)\s+([a-zA-Z0-9_]+)`)
	cFnRe := regexp.MustCompile(`(?i)^[a-zA-Z0-9_\*]+\s+([a-zA-Z0-9_]+)\s*\([^;]*\)\s*\{?`)
	typeRe := regexp.MustCompile(`(?i)(?:pub\s+)?(?:type|struct|class|interface|trait|enum)\s+([a-zA-Z0-9_]+)`)
	importRe := regexp.MustCompile(`(?i)^(?:import|use|#include|from|require)\s+(.+)`)

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || (spec.CommentPrefix != "" && strings.HasPrefix(line, spec.CommentPrefix)) {
			continue
		}

		if match := importRe.FindStringSubmatch(line); len(match) > 1 {
			imports = append(imports, line)
			continue
		}

		if match := fnRe.FindStringSubmatch(line); len(match) > 1 {
			elements = append(elements, CodeElement{
				Kind:      "function",
				Name:      match[1],
				Signature: line,
				Line:      lineNum,
				EndLine:   lineNum,
			})
			continue
		}

		if match := cFnRe.FindStringSubmatch(line); len(match) > 1 && !strings.HasPrefix(line, "return") && !strings.HasPrefix(line, "if") && !strings.HasPrefix(line, "while") {
			elements = append(elements, CodeElement{
				Kind:      "function",
				Name:      match[1],
				Signature: line,
				Line:      lineNum,
				EndLine:   lineNum,
			})
			continue
		}

		if match := typeRe.FindStringSubmatch(line); len(match) > 1 {
			kind := "type"
			lower := strings.ToLower(line)
			if strings.Contains(lower, "struct") {
				kind = "struct"
			} else if strings.Contains(lower, "class") {
				kind = "class"
			} else if strings.Contains(lower, "interface") {
				kind = "interface"
			} else if strings.Contains(lower, "trait") {
				kind = "trait"
			} else if strings.Contains(lower, "enum") {
				kind = "enum"
			}
			elements = append(elements, CodeElement{
				Kind:      kind,
				Name:      match[1],
				Signature: line,
				Line:      lineNum,
				EndLine:   lineNum,
			})
		}
	}

	return &FileStructureResult{
		FilePath: filePath,
		Language: spec.Name,
		Elements: elements,
		Imports:  imports,
		RawCode:  string(src),
	}
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "\n"); idx > 0 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}
