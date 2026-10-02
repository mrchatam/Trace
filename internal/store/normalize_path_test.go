package store

import "testing"

// TestNormalizePathCanonicalization locks the canonical form of stored file
// paths: repo-relative, forward slashes, single separators, no leading slash.
// Every writer and every lookup runs through this function, so storage and
// lookup must agree by construction (R3).
func TestNormalizePathCanonicalization(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a.go", "a.go"},
		{"./a.go", "a.go"},
		{"pkg/mod.go", "pkg/mod.go"},
		{"./pkg/mod.go", "pkg/mod.go"},
		{"pkg\\mod.go", "pkg/mod.go"},
		{"/pkg/mod.go", "pkg/mod.go"},
		{"/a.go", "a.go"},
		{"pkg//mod.go", "pkg/mod.go"},
		{"pkg///mod.go", "pkg/mod.go"},
		{"//a.go", "a.go"},
		{"a//b//c.go", "a/b/c.go"},
		{"", ""},
		{"/", ""},
		{"//", ""},
	}
	for _, c := range cases {
		if got := NormalizePath(c.in); got != c.want {
			t.Errorf("NormalizePath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestUpsertFileCanonicalizesPath locks the write side of the round-trip:
// whatever the caller passes, the stored path is the canonical form.
func TestUpsertFileCanonicalizesPath(t *testing.T) {
	s, _ := openTempStore(t)
	for _, in := range []string{"/a.go", ".//a.go", "pkg//mod.go"} {
		want := NormalizePath(in)
		f, err := s.UpsertFile(in, "h:"+in, nil)
		if err != nil {
			t.Fatalf("UpsertFile(%q): %v", in, err)
		}
		if f.Path != want {
			t.Fatalf("UpsertFile(%q) stored path %q, want %q", in, f.Path, want)
		}
	}
}
