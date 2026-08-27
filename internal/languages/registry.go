package languages

import (
	"path/filepath"
	"strings"
	"sync"
)

// Registry manages registered language specifications.
type Registry struct {
	languages map[string]*LanguageSpec
	extMap    map[string]*LanguageSpec
	mu        sync.RWMutex
}

var defaultRegistry *Registry
var registryOnce sync.Once

// GetRegistry returns the global language registry initialized with built-in language definitions.
func GetRegistry() *Registry {
	registryOnce.Do(func() {
		defaultRegistry = NewRegistry()
		defaultRegistry.registerBuiltins()
	})
	return defaultRegistry
}

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		languages: make(map[string]*LanguageSpec),
		extMap:    make(map[string]*LanguageSpec),
	}
}

// Register adds or updates a language specification in the registry.
func (r *Registry) Register(spec LanguageSpec) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s := &spec
	r.languages[strings.ToLower(spec.Name)] = s
	for _, alias := range spec.Aliases {
		r.languages[strings.ToLower(alias)] = s
	}
	for _, ext := range spec.Extensions {
		cleanExt := strings.ToLower(ext)
		if !strings.HasPrefix(cleanExt, ".") {
			cleanExt = "." + cleanExt
		}
		r.extMap[cleanExt] = s
	}
}

// FindByExtension locates the language spec matching a file extension.
func (r *Registry) FindByExtension(ext string) (*LanguageSpec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cleanExt := strings.ToLower(ext)
	if !strings.HasPrefix(cleanExt, ".") {
		cleanExt = "." + cleanExt
	}
	spec, ok := r.extMap[cleanExt]
	return spec, ok
}

// FindByPath locates the language spec matching a file path.
func (r *Registry) FindByPath(path string) (*LanguageSpec, bool) {
	return r.FindByExtension(filepath.Ext(path))
}

// FindByName locates the language spec by language name or alias.
func (r *Registry) FindByName(name string) (*LanguageSpec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	spec, ok := r.languages[strings.ToLower(strings.TrimSpace(name))]
	return spec, ok
}

// registerBuiltins populates default language specifications for all major programming languages.
func (r *Registry) registerBuiltins() {
	// 1. Rust
	r.Register(LanguageSpec{
		Name:             "rust",
		Aliases:          []string{"rs"},
		Extensions:       []string{".rs"},
		TreeSitterSymbol: "tree_sitter_rust",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Functions:  []string{"function_item"},
			Structs:    []string{"struct_item"},
			Enums:      []string{"enum_item"},
			Traits:     []string{"trait_item"},
			Impls:      []string{"impl_item"},
			Types:      []string{"type_item"},
			Imports:    []string{"use_declaration"},
			Macros:     []string{"macro_definition", "macro_invocation"},
			Variables:  []string{"static_item", "const_item"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "cargo",
				BuildFiles:      []string{"Cargo.toml"},
				BuildCommand:    []string{"cargo", "check", "--tests"},
				TestCommand:     []string{"cargo", "test", "--lib", "--tests", "--", "--nocapture"},
				SourceDirectory: "src",
				TestDirectory:   "tests",
			},
		},
		DefaultSkeleton: map[string]string{
			"Cargo.toml": "[package]\nname = \"{{.ProjectName}}\"\nversion = \"0.1.0\"\nedition = \"2021\"\n\n[dependencies]\n",
			"src/lib.rs": "// Library interface\n",
		},
	})

	// 2. C
	r.Register(LanguageSpec{
		Name:             "c",
		Aliases:          []string{"h"},
		Extensions:       []string{".c", ".h"},
		TreeSitterSymbol: "tree_sitter_c",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Functions: []string{"function_definition"},
			Structs:   []string{"struct_specifier"},
			Enums:     []string{"enum_specifier"},
			Types:     []string{"type_definition", "declaration"},
			Imports:   []string{"preproc_include"},
			Macros:    []string{"preproc_def", "preproc_function_def"},
			Variables: []string{"declaration"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "make",
				BuildFiles:      []string{"Makefile", "makefile"},
				BuildCommand:    []string{"make"},
				TestCommand:     []string{"make", "test"},
				SourceDirectory: ".",
				TestDirectory:   "tests",
			},
			{
				Name:            "cmake",
				BuildFiles:      []string{"CMakeLists.txt"},
				BuildCommand:    []string{"cmake", "--build", "build"},
				TestCommand:     []string{"ctest", "--test-dir", "build"},
				SourceDirectory: "src",
				TestDirectory:   "tests",
			},
		},
		DefaultSkeleton: map[string]string{
			"Makefile": "all:\n\t@echo \"Build\"\ntest:\n\t@echo \"Test\"\n",
		},
	})

	// 3. C++
	r.Register(LanguageSpec{
		Name:             "cpp",
		Aliases:          []string{"c++", "cc", "cxx", "hpp", "hh"},
		Extensions:       []string{".cpp", ".cc", ".cxx", ".hpp", ".hh"},
		TreeSitterSymbol: "tree_sitter_cpp",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Functions:  []string{"function_definition"},
			Methods:    []string{"field_declaration"},
			Classes:    []string{"class_specifier"},
			Structs:    []string{"struct_specifier"},
			Enums:      []string{"enum_specifier"},
			Interfaces: []string{"class_specifier"},
			Types:      []string{"type_definition", "declaration", "alias_declaration"},
			Imports:    []string{"preproc_include"},
			Macros:     []string{"preproc_def", "preproc_function_def"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "cmake",
				BuildFiles:      []string{"CMakeLists.txt"},
				BuildCommand:    []string{"cmake", "--build", "build"},
				TestCommand:     []string{"ctest", "--test-dir", "build"},
				SourceDirectory: "src",
				TestDirectory:   "tests",
			},
			{
				Name:            "make",
				BuildFiles:      []string{"Makefile"},
				BuildCommand:    []string{"make"},
				TestCommand:     []string{"make", "test"},
				SourceDirectory: "src",
				TestDirectory:   "tests",
			},
		},
	})

	// 4. Go
	r.Register(LanguageSpec{
		Name:             "go",
		Aliases:          []string{"golang"},
		Extensions:       []string{".go"},
		TreeSitterSymbol: "tree_sitter_go",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Functions:  []string{"function_declaration"},
			Methods:    []string{"method_declaration"},
			Structs:    []string{"type_declaration", "type_spec"},
			Interfaces: []string{"interface_type"},
			Types:      []string{"type_declaration", "type_alias"},
			Imports:    []string{"import_declaration"},
			Variables:  []string{"var_declaration", "const_declaration"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "go",
				BuildFiles:      []string{"go.mod"},
				BuildCommand:    []string{"go", "build", "./..."},
				TestCommand:     []string{"go", "test", "-v", "./..."},
				SourceDirectory: ".",
				TestDirectory:   ".",
			},
		},
		DefaultSkeleton: map[string]string{
			"go.mod": "module {{.ProjectName}}\n\ngo 1.22\n",
			"lib.go": "package {{.ProjectName}}\n",
		},
	})

	// 5. Java
	r.Register(LanguageSpec{
		Name:             "java",
		Extensions:       []string{".java"},
		TreeSitterSymbol: "tree_sitter_java",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Classes:    []string{"class_declaration"},
			Interfaces: []string{"interface_declaration"},
			Enums:      []string{"enum_declaration"},
			Methods:    []string{"method_declaration"},
			Imports:    []string{"import_declaration"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "maven",
				BuildFiles:      []string{"pom.xml"},
				BuildCommand:    []string{"mvn", "compile"},
				TestCommand:     []string{"mvn", "test"},
				SourceDirectory: "src/main/java",
				TestDirectory:   "src/test/java",
			},
			{
				Name:            "gradle",
				BuildFiles:      []string{"build.gradle", "build.gradle.kts"},
				BuildCommand:    []string{"gradle", "build"},
				TestCommand:     []string{"gradle", "test"},
				SourceDirectory: "src/main/java",
				TestDirectory:   "src/test/java",
			},
		},
	})

	// 6. Python
	r.Register(LanguageSpec{
		Name:             "python",
		Aliases:          []string{"py"},
		Extensions:       []string{".py", ".pyi"},
		TreeSitterSymbol: "tree_sitter_python",
		CommentPrefix:    "#",
		ASTMapping: ASTNodeKindMap{
			Functions: []string{"function_definition"},
			Classes:   []string{"class_definition"},
			Imports:   []string{"import_statement", "import_from_statement"},
			Variables: []string{"assignment"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "pytest",
				BuildFiles:      []string{"pyproject.toml", "setup.py", "requirements.txt"},
				BuildCommand:    []string{"python", "-m", "py_compile"},
				TestCommand:     []string{"pytest", "-v"},
				SourceDirectory: "src",
				TestDirectory:   "tests",
			},
		},
	})

	// 7. TypeScript
	r.Register(LanguageSpec{
		Name:             "typescript",
		Aliases:          []string{"ts", "tsx"},
		Extensions:       []string{".ts", ".tsx"},
		TreeSitterSymbol: "tree_sitter_typescript",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Functions:  []string{"function_declaration", "arrow_function"},
			Methods:    []string{"method_definition"},
			Classes:    []string{"class_declaration"},
			Interfaces: []string{"interface_declaration"},
			Types:      []string{"type_alias_declaration"},
			Enums:      []string{"enum_declaration"},
			Imports:    []string{"import_statement"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "npm",
				BuildFiles:      []string{"package.json"},
				BuildCommand:    []string{"npm", "run", "build"},
				TestCommand:     []string{"npm", "test"},
				SourceDirectory: "src",
				TestDirectory:   "test",
			},
		},
	})

	// 8. JavaScript
	r.Register(LanguageSpec{
		Name:             "javascript",
		Aliases:          []string{"js", "jsx", "mjs", "cjs"},
		Extensions:       []string{".js", ".jsx", ".mjs", ".cjs"},
		TreeSitterSymbol: "tree_sitter_javascript",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Functions: []string{"function_declaration", "arrow_function"},
			Methods:   []string{"method_definition"},
			Classes:   []string{"class_declaration"},
			Imports:   []string{"import_statement"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:            "npm",
				BuildFiles:      []string{"package.json"},
				BuildCommand:    []string{"npm", "run", "build"},
				TestCommand:     []string{"npm", "test"},
				SourceDirectory: "src",
				TestDirectory:   "test",
			},
		},
	})

	// 9. Kotlin
	r.Register(LanguageSpec{
		Name:             "kotlin",
		Aliases:          []string{"kt", "kts"},
		Extensions:       []string{".kt", ".kts"},
		TreeSitterSymbol: "tree_sitter_kotlin",
		CommentPrefix:    "//",
		ASTMapping: ASTNodeKindMap{
			Functions:  []string{"function_declaration"},
			Classes:    []string{"class_declaration"},
			Interfaces: []string{"class_declaration"},
			Imports:    []string{"import_header"},
		},
		Toolchains: []ToolchainSpec{
			{
				Name:         "gradle",
				BuildFiles:   []string{"build.gradle.kts", "build.gradle"},
				BuildCommand: []string{"gradle", "build"},
				TestCommand:  []string{"gradle", "test"},
			},
		},
	})
}
