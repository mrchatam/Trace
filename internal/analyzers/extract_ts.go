package analyzers

import (
	"unsafe"

	"github.com/mrchatam/Trace/internal/store"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

func extractTS(content []byte, tsx bool) ([]store.Symbol, []store.Import, error) {
	var langPtr unsafe.Pointer
	if tsx {
		langPtr = tree_sitter_typescript.LanguageTSX()
	} else {
		langPtr = tree_sitter_typescript.LanguageTypescript()
	}
	lang := tree_sitter.NewLanguage(langPtr)
	return extractJSTS(lang, content, tsSymbolQuery)
}
