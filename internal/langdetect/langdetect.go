// Package langdetect maps file paths to Tier-1 language IDs with no
// dependencies — not even the store. Language detection used to live only in
// internal/analyzers, so consumers that needed nothing but "is this an
// indexable extension?" (indexwalk's tree walk, domain's path classification)
// pulled the whole tree-sitter analyzer graph into their build. The table here
// is the single source of truth for extension→language; the analyzers package
// delegates to it and a test pins the adapter registry to this table, so a
// contributor adding a language cannot let the two drift apart.
package langdetect

import (
	"path/filepath"
	"sort"
	"strings"
)

// Language identifiers persisted on files.language.
const (
	JavaScript = "javascript"
	TypeScript = "typescript"
	TSX        = "tsx"
	Python     = "python"
	Go         = "go"
)

// langByExt is the Tier-1 extension table. Keys must be lowercase with a
// leading dot; Detect lowercases the path extension before lookup.
var langByExt = map[string]string{
	".js":  JavaScript,
	".jsx": JavaScript,
	".mjs": JavaScript,
	".cjs": JavaScript,
	".ts":  TypeScript,
	".tsx": TSX,
	".py":  Python,
	".go":  Go,
}

// Detect maps a path extension to a supported language id. ok is false for
// unsupported extensions.
func Detect(path string) (lang string, ok bool) {
	lang, ok = langByExt[strings.ToLower(filepath.Ext(path))]
	return lang, ok
}

// AllExts returns every extension in the table, sorted. Test-only surface for
// pinning the analyzers adapter registry to this table.
func AllExts() []string {
	out := make([]string, 0, len(langByExt))
	for ext := range langByExt {
		out = append(out, ext)
	}
	sort.Strings(out)
	return out
}
