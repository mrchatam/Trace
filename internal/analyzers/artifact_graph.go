package analyzers

import (
	"strings"
	"unicode"

	"github.com/mrchatam/Trace/internal/store"
)

func computeArtifactEdges(st *store.Store, p string, content []byte, lang string) ([]store.CodeEdge, error) {
	f, err := st.GetFileByPath(p)
	if err != nil {
		return nil, err
	}
	syms, err := st.ListSymbolsByPath(p)
	if err != nil {
		return nil, err
	}
	var jstsKeys map[string]bool
	switch lang {
	case LangJavaScript, LangTypeScript, LangTSX:
		jstsKeys, err = jstsExportedKeys(lang, content)
		if err != nil {
			return nil, err
		}
	}
	var out []store.CodeEdge
	for _, sym := range syms {
		id := sym.ID
		out = append(out, store.CodeEdge{
			FromFileID: f.ID,
			ToFileID:   f.ID,
			ToSymbolID: &id,
			Rel:        store.RelContainsModule,
			Provenance: store.ImportProvenanceExtracted,
		})
		if !symbolExportsAPI(lang, sym, jstsKeys) {
			continue
		}
		out = append(out, store.CodeEdge{
			FromFileID: f.ID,
			ToFileID:   f.ID,
			ToSymbolID: &id,
			Rel:        store.RelExportsAPI,
			Provenance: store.ImportProvenanceExtracted,
		})
	}
	return out, nil
}

func symbolExportsAPI(lang string, sym store.Symbol, jstsKeys map[string]bool) bool {
	switch lang {
	case LangGo:
		return goSymbolExported(sym)
	case LangPython:
		return pythonSymbolExported(sym)
	case LangJavaScript, LangTypeScript, LangTSX:
		if jstsKeys == nil {
			return false
		}
		return jstsKeys[exportKey(sym.Kind, sym.Name, sym.StartLine)]
	default:
		return false
	}
}

func goSymbolExported(sym store.Symbol) bool {
	switch sym.Kind {
	case "function", "method", "type":
	default:
		return false
	}
	for _, r := range sym.Name {
		return unicode.IsUpper(r)
	}
	return false
}

func pythonSymbolExported(sym store.Symbol) bool {
	switch sym.Kind {
	case "function", "method", "class":
	default:
		return false
	}
	return sym.Name != "" && !strings.HasPrefix(sym.Name, "_")
}
