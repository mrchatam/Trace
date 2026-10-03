package analyzers

import (
	"github.com/mrchatam/Trace/internal/langdetect"
)

// Language identifiers persisted on files.language.
const (
	LangJavaScript = langdetect.JavaScript
	LangTypeScript = langdetect.TypeScript
	LangTSX        = langdetect.TSX
	LangPython     = langdetect.Python
	LangGo         = langdetect.Go
)

// DetectLanguage maps a path extension to a supported language id via the
// compile-time builtin LanguageAdapter table. ok is false for unsupported extensions.
//
// Detection itself is delegated to internal/langdetect (pure, CGO-free) so
// callers that only need extension classification do not link the tree-sitter
// graph; the builtinAdapters registry remains the source of truth for
// extraction. TestBuiltinLanguageAdaptersContributionPath pins both directions:
// every adapter extension is detectable, and every table extension has an adapter.
func DetectLanguage(path string) (lang string, ok bool) {
	return langdetect.Detect(path)
}
