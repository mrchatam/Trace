package analyzers

import (
	"path/filepath"
	"strings"
)

// Language identifiers persisted on files.language.
const (
	LangJavaScript = "javascript"
	LangTypeScript = "typescript"
	LangTSX        = "tsx"
	LangPython     = "python"
	LangGo         = "go"
)

// DetectLanguage maps a path extension to a supported language id via the
// compile-time builtin LanguageAdapter table. ok is false for unsupported extensions.
func DetectLanguage(path string) (lang string, ok bool) {
	ext := strings.ToLower(filepath.Ext(path))
	a, ok := adapterByExt(ext)
	if !ok {
		return "", false
	}
	return a.ID(), true
}
