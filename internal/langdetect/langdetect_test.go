package langdetect

import "testing"

func TestDetectKnownExtensions(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"app.js", JavaScript},
		{"app.jsx", JavaScript},
		{"app.mjs", JavaScript},
		{"app.cjs", JavaScript},
		{"src/index.ts", TypeScript},
		{"src/App.tsx", TSX},
		{"main.py", Python},
		{"cmd/trace/main.go", Go},
	}
	for _, c := range cases {
		got, ok := Detect(c.path)
		if !ok || got != c.want {
			t.Fatalf("Detect(%q) = (%q,%v), want (%q,true)", c.path, got, ok, c.want)
		}
	}
}

func TestDetectCaseInsensitiveExtension(t *testing.T) {
	got, ok := Detect("README.GO")
	if !ok || got != Go {
		t.Fatalf("Detect(README.GO) = (%q,%v), want (%q,true)", got, ok, Go)
	}
}

func TestDetectUnsupported(t *testing.T) {
	// Note: Detect(".js") is (JavaScript, true) — filepath.Ext treats a bare
	// ".js" as extension-only. Pre-existing semantics; don't change silently.
	for _, p := range []string{
		"foo.txt", "foo.rs", "foo", "foo.", "dir.js/", "a.min.js.bak",
	} {
		if lang, ok := Detect(p); ok {
			t.Fatalf("Detect(%q) = (%q,true), want unsupported", p, lang)
		}
	}
}

func TestAllExtsMatchesTable(t *testing.T) {
	exts := AllExts()
	if len(exts) != len(langByExt) {
		t.Fatalf("AllExts len = %d, want %d", len(exts), len(langByExt))
	}
	for _, ext := range exts {
		if _, ok := langByExt[ext]; !ok {
			t.Fatalf("AllExts returned %q not in table", ext)
		}
		if lang := langByExt[ext]; lang == "" {
			t.Fatalf("extension %q maps to empty language", ext)
		}
	}
}
