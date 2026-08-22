package analyzers

import (
	"sort"

	"github.com/mrchatam/Trace/internal/store"
)

// LanguageAdapterAPIVersion is the contribution-surface contract version.
// Bump only on breaking LanguageAdapter iface or dispatch semantics change.
const LanguageAdapterAPIVersion = 1

// LanguageAdapter is the versioned analyzer contribution surface for one language.
// Extensions must be lowercase with a leading dot (e.g. ".go").
type LanguageAdapter interface {
	ID() string
	Extensions() []string
	Extract(content []byte) ([]store.Symbol, []store.Import, error)
}

// builtinAdapters is the compile-time static table of built-in language adapters.
// There is no public Register mutator and no dynamic .so loader.
var builtinAdapters = []LanguageAdapter{
	jsAdapter{},
	tsAdapter{},
	tsxAdapter{},
	pythonAdapter{},
	goAdapter{},
}

var (
	adaptersByExt = map[string]LanguageAdapter{}
	adaptersByID  = map[string]LanguageAdapter{}
)

func init() {
	for _, a := range builtinAdapters {
		adaptersByID[a.ID()] = a
		for _, ext := range a.Extensions() {
			adaptersByExt[ext] = a
		}
	}
}

func adapterByExt(ext string) (LanguageAdapter, bool) {
	a, ok := adaptersByExt[ext]
	return a, ok
}

func adapterByID(id string) (LanguageAdapter, bool) {
	a, ok := adaptersByID[id]
	return a, ok
}

// SupportedLanguages returns sorted Tier-1 language IDs from the compile-time
// builtinAdapters table (honesty surface for CLI/HTTP index status).
func SupportedLanguages() []string {
	out := make([]string, len(builtinAdapters))
	for i, a := range builtinAdapters {
		out[i] = a.ID()
	}
	sort.Strings(out)
	return out
}

type jsAdapter struct{}

func (jsAdapter) ID() string { return LangJavaScript }
func (jsAdapter) Extensions() []string {
	return []string{".js", ".jsx", ".mjs", ".cjs"}
}
func (jsAdapter) Extract(content []byte) ([]store.Symbol, []store.Import, error) {
	return extractJS(content)
}

type tsAdapter struct{}

func (tsAdapter) ID() string           { return LangTypeScript }
func (tsAdapter) Extensions() []string { return []string{".ts"} }
func (tsAdapter) Extract(content []byte) ([]store.Symbol, []store.Import, error) {
	return extractTS(content, false)
}

type tsxAdapter struct{}

func (tsxAdapter) ID() string           { return LangTSX }
func (tsxAdapter) Extensions() []string { return []string{".tsx"} }
func (tsxAdapter) Extract(content []byte) ([]store.Symbol, []store.Import, error) {
	return extractTS(content, true)
}

type pythonAdapter struct{}

func (pythonAdapter) ID() string           { return LangPython }
func (pythonAdapter) Extensions() []string { return []string{".py"} }
func (pythonAdapter) Extract(content []byte) ([]store.Symbol, []store.Import, error) {
	return extractPython(content)
}

type goAdapter struct{}

func (goAdapter) ID() string           { return LangGo }
func (goAdapter) Extensions() []string { return []string{".go"} }
func (goAdapter) Extract(content []byte) ([]store.Symbol, []store.Import, error) {
	return extractGo(content)
}
