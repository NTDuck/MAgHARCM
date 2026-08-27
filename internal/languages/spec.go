package languages

// NodeCategory classifies semantic AST node roles into canonical Code IR categories.
type NodeCategory string

const (
	CategoryFunction  NodeCategory = "function"
	CategoryMethod    NodeCategory = "method"
	CategoryStruct    NodeCategory = "struct"
	CategoryClass     NodeCategory = "class"
	CategoryInterface NodeCategory = "interface"
	CategoryTrait     NodeCategory = "trait"
	CategoryType      NodeCategory = "type"
	CategoryEnum      NodeCategory = "enum"
	CategoryImport    NodeCategory = "import"
	CategoryImpl      NodeCategory = "impl"
	CategoryMacro     NodeCategory = "macro"
	CategoryVariable  NodeCategory = "variable"
)

// ASTNodeKindMap defines mapping rules from language-specific Tree-Sitter node kinds to canonical categories.
type ASTNodeKindMap struct {
	Functions  []string `json:"functions" yaml:"functions"`
	Methods    []string `json:"methods" yaml:"methods"`
	Structs    []string `json:"structs" yaml:"structs"`
	Classes    []string `json:"classes" yaml:"classes"`
	Interfaces []string `json:"interfaces" yaml:"interfaces"`
	Traits     []string `json:"traits" yaml:"traits"`
	Types      []string `json:"types" yaml:"types"`
	Enums      []string `json:"enums" yaml:"enums"`
	Imports    []string `json:"imports" yaml:"imports"`
	Impls      []string `json:"impls" yaml:"impls"`
	Macros     []string `json:"macros" yaml:"macros"`
	Variables  []string `json:"variables" yaml:"variables"`
}

// ToolchainSpec defines how to build and test projects in this language.
type ToolchainSpec struct {
	Name            string   `json:"name" yaml:"name"`
	BuildFiles      []string `json:"build_files" yaml:"build_files"`           // e.g. ["Cargo.toml"], ["go.mod"], ["Makefile"]
	BuildCommand    []string `json:"build_command" yaml:"build_command"`       // e.g. ["cargo", "check", "--tests"]
	TestCommand     []string `json:"test_command" yaml:"test_command"`         // e.g. ["cargo", "test", "--lib", "--tests", "--", "--nocapture"]
	SourceDirectory string   `json:"source_directory" yaml:"source_directory"` // e.g. "src", "."
	TestDirectory   string   `json:"test_directory" yaml:"test_directory"`     // e.g. "tests", "."
}

// LanguageSpec defines everything required to parse, analyze, generate skeletons, build, and test a language.
type LanguageSpec struct {
	Name             string            `json:"name" yaml:"name"`
	Aliases          []string          `json:"aliases" yaml:"aliases"`
	Extensions       []string          `json:"extensions" yaml:"extensions"`
	TreeSitterSymbol string            `json:"tree_sitter_symbol" yaml:"tree_sitter_symbol"` // e.g. "tree_sitter_rust"
	GrammarLibraries []string          `json:"grammar_libraries" yaml:"grammar_libraries"`   // e.g. ["libtree-sitter-rust.so", "tree-sitter-rust.so"]
	CommentPrefix    string            `json:"comment_prefix" yaml:"comment_prefix"`
	ASTMapping       ASTNodeKindMap    `json:"ast_mapping" yaml:"ast_mapping"`
	Toolchains       []ToolchainSpec   `json:"toolchains" yaml:"toolchains"`
	DefaultSkeleton  map[string]string `json:"default_skeleton,omitempty" yaml:"default_skeleton,omitempty"`
}

// CategorizeNode maps a Tree-Sitter node type to a canonical Code IR category.
func (m *ASTNodeKindMap) CategorizeNode(nodeType string) (NodeCategory, bool) {
	for _, k := range m.Functions {
		if k == nodeType {
			return CategoryFunction, true
		}
	}
	for _, k := range m.Methods {
		if k == nodeType {
			return CategoryMethod, true
		}
	}
	for _, k := range m.Structs {
		if k == nodeType {
			return CategoryStruct, true
		}
	}
	for _, k := range m.Classes {
		if k == nodeType {
			return CategoryClass, true
		}
	}
	for _, k := range m.Interfaces {
		if k == nodeType {
			return CategoryInterface, true
		}
	}
	for _, k := range m.Traits {
		if k == nodeType {
			return CategoryTrait, true
		}
	}
	for _, k := range m.Types {
		if k == nodeType {
			return CategoryType, true
		}
	}
	for _, k := range m.Enums {
		if k == nodeType {
			return CategoryEnum, true
		}
	}
	for _, k := range m.Imports {
		if k == nodeType {
			return CategoryImport, true
		}
	}
	for _, k := range m.Impls {
		if k == nodeType {
			return CategoryImpl, true
		}
	}
	for _, k := range m.Macros {
		if k == nodeType {
			return CategoryMacro, true
		}
	}
	for _, k := range m.Variables {
		if k == nodeType {
			return CategoryVariable, true
		}
	}
	return "", false
}
