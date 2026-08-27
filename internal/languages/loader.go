package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	"MAgHARCM/internal/logger"
)

// LanguageLoader dynamically loads Tree-Sitter grammars from shared libraries at runtime via purego.
type LanguageLoader struct {
	searchPaths []string
	cache       map[string]*tree_sitter.Language
	mu          sync.RWMutex
}

var globalLoader *LanguageLoader
var loaderOnce sync.Once

// GetLoader returns the singleton LanguageLoader instance.
func GetLoader() *LanguageLoader {
	loaderOnce.Do(func() {
		homeDir, _ := os.UserHomeDir()
		paths := []string{
			".",
			"./lib",
			"./grammars",
			filepath.Join(homeDir, ".local", "lib", "tree-sitter"),
			filepath.Join(homeDir, ".config", "tree-sitter", "bin"),
			"/usr/local/lib",
			"/usr/lib",
			"/usr/lib64",
		}
		if envDir := os.Getenv("TREE_SITTER_DIR"); envDir != "" {
			paths = append([]string{envDir}, paths...)
		}
		if ldPath := os.Getenv("LD_LIBRARY_PATH"); ldPath != "" {
			paths = append(paths, filepath.SplitList(ldPath)...)
		}

		globalLoader = &LanguageLoader{
			searchPaths: paths,
			cache:       make(map[string]*tree_sitter.Language),
		}
	})
	return globalLoader
}

// AddSearchPath adds an additional directory to search for grammar shared libraries.
func (l *LanguageLoader) AddSearchPath(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.searchPaths = append([]string{path}, l.searchPaths...)
}

// LoadGrammar attempts to dynamically locate and load the Tree-Sitter grammar for a language at runtime.
func (l *LanguageLoader) LoadGrammar(spec *LanguageSpec) (*tree_sitter.Language, error) {
	l.mu.RLock()
	if cached, ok := l.cache[spec.Name]; ok {
		l.mu.RUnlock()
		return cached, nil
	}
	l.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	// Double check under write lock
	if cached, ok := l.cache[spec.Name]; ok {
		return cached, nil
	}

	symbolName := spec.TreeSitterSymbol
	if symbolName == "" {
		symbolName = fmt.Sprintf("tree_sitter_%s", spec.Name)
	}

	// Candidates for shared library filenames based on OS
	var libNames []string
	switch runtime.GOOS {
	case "darwin":
		libNames = []string{
			fmt.Sprintf("libtree-sitter-%s.dylib", spec.Name),
			fmt.Sprintf("tree-sitter-%s.dylib", spec.Name),
			fmt.Sprintf("libtree-sitter-%s.so", spec.Name),
			fmt.Sprintf("tree_sitter_%s.dylib", spec.Name),
		}
	case "windows":
		libNames = []string{
			fmt.Sprintf("tree-sitter-%s.dll", spec.Name),
			fmt.Sprintf("libtree-sitter-%s.dll", spec.Name),
			fmt.Sprintf("tree_sitter_%s.dll", spec.Name),
		}
	default: // linux, bsd
		libNames = []string{
			fmt.Sprintf("libtree-sitter-%s.so", spec.Name),
			fmt.Sprintf("tree-sitter-%s.so", spec.Name),
			fmt.Sprintf("libtree-sitter-%s.so.0", spec.Name),
			fmt.Sprintf("tree_sitter_%s.so", spec.Name),
		}
	}
	for _, custom := range spec.GrammarLibraries {
		libNames = append([]string{custom}, libNames...)
	}

	// 1. Try loading from dynamic shared object libraries via purego.Dlopen
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
							if tsLang != nil {
								l.cache[spec.Name] = tsLang
								logger.LogStep("Loaded dynamic Tree-Sitter grammar for `%s` from `%s`", spec.Name, fullPath)
								return tsLang, nil
							}
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
			if tsLang != nil {
				l.cache[spec.Name] = tsLang
				return tsLang, nil
			}
		}
	}

	return nil, fmt.Errorf("dynamic tree-sitter grammar for '%s' (symbol %s) not found in search paths", spec.Name, symbolName)
}
