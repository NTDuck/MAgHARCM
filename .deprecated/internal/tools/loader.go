package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)

var (
	grammarMu    sync.RWMutex
	grammarCache = make(map[string]*tree_sitter.Language)
)

// LanguageLoader dynamically loads tree-sitter grammars at runtime using purego.
type LanguageLoader struct {
	searchPaths []string
}

// NewLanguageLoader creates a loader with configurable search paths.
func NewLanguageLoader(searchPaths ...string) *LanguageLoader {
	paths := append([]string{}, searchPaths...)
	if envDir := os.Getenv("TREE_SITTER_LIB_DIR"); envDir != "" {
		paths = append(paths, envDir)
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths,
			filepath.Join(home, ".tree-sitter", "bin"),
			filepath.Join(home, ".local", "lib"),
			filepath.Join(home, ".config", "tree-sitter", "bin"),
		)
	}
	paths = append(paths, ".", "lib", "/usr/local/lib", "/usr/lib")
	return &LanguageLoader{searchPaths: paths}
}

// LoadLanguage dynamically loads a tree-sitter Language by language name (e.g. "c", "rust", "go", "cpp")
// using purego runtime symbol resolution.
func (l *LanguageLoader) LoadLanguage(lang string) (*tree_sitter.Language, error) {
	grammarMu.RLock()
	if cached, ok := grammarCache[lang]; ok {
		grammarMu.RUnlock()
		return cached, nil
	}
	grammarMu.RUnlock()

	symbolName := fmt.Sprintf("tree_sitter_%s", lang)

	// 1. Try loading from dynamic shared object libraries via purego.Dlopen
	var libNames []string
	switch runtime.GOOS {
	case "darwin":
		libNames = []string{
			fmt.Sprintf("libtree-sitter-%s.dylib", lang),
			fmt.Sprintf("tree-sitter-%s.dylib", lang),
			fmt.Sprintf("libtree-sitter-%s.so", lang),
		}
	case "windows":
		libNames = []string{
			fmt.Sprintf("tree-sitter-%s.dll", lang),
			fmt.Sprintf("libtree-sitter-%s.dll", lang),
		}
	default: // linux, bsd
		libNames = []string{
			fmt.Sprintf("libtree-sitter-%s.so", lang),
			fmt.Sprintf("tree-sitter-%s.so", lang),
			fmt.Sprintf("libtree-sitter-%s.so.0", lang),
		}
	}

	for _, dir := range l.searchPaths {
		for _, libName := range libNames {
			fullPath := filepath.Join(dir, libName)
			if _, err := os.Stat(fullPath); err == nil {
				handle, err := purego.Dlopen(fullPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
				if err == nil {
					sym, err := purego.Dlsym(handle, symbolName)
					if err == nil && sym != 0 {
						var fn func() unsafe.Pointer
						purego.RegisterFunc(&fn, sym)
						ptr := fn()
						if ptr != nil {
							tsLang := tree_sitter.NewLanguage(ptr)
							grammarMu.Lock()
							grammarCache[lang] = tsLang
							grammarMu.Unlock()
							return tsLang, nil
						}
					}
				}
			}
		}
	}

	// 2. Try in-process symbol resolution via purego Dlsym from RTLD_DEFAULT
	if sym, err := purego.Dlsym(0, symbolName); err == nil && sym != 0 {
		var fn func() unsafe.Pointer
		purego.RegisterFunc(&fn, sym)
		ptr := fn()
		if ptr != nil {
			tsLang := tree_sitter.NewLanguage(ptr)
			grammarMu.Lock()
			grammarCache[lang] = tsLang
			grammarMu.Unlock()
			return tsLang, nil
		}
	}

	// 3. Fallback to statically linked binding pointers if dynamic libraries are not on disk
	var ptr unsafe.Pointer
	switch lang {
	case "c":
		ptr = tree_sitter_c.Language()
	case "rust":
		ptr = tree_sitter_rust.Language()
	case "go":
		ptr = tree_sitter_go.Language()
	case "cpp":
		ptr = tree_sitter_cpp.Language()
	}

	if ptr != nil {
		tsLang := tree_sitter.NewLanguage(ptr)
		grammarMu.Lock()
		grammarCache[lang] = tsLang
		grammarMu.Unlock()
		return tsLang, nil
	}

	return nil, fmt.Errorf("tree-sitter language grammar not found for: %s", lang)
}

// DefaultLanguageLoader is the package-level runtime grammar loader.
var DefaultLanguageLoader = NewLanguageLoader()
