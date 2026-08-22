package analyzers

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"
)

func TestSupportedLanguagesMatchesAdapters(t *testing.T) {
	got := SupportedLanguages()
	want := make([]string, len(builtinAdapters))
	for i, a := range builtinAdapters {
		want[i] = a.ID()
	}
	sort.Strings(want)
	if len(got) != 5 {
		t.Fatalf("len=%d want 5: %v", len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SupportedLanguages()[%d]=%q want %q; full got=%v want=%v", i, got[i], want[i], got, want)
		}
	}
}

func TestIndexUnsupportedExtMessage(t *testing.T) {
	st := openTemp(t)
	err := IndexFile(context.Background(), st, "main.rs", []byte("fn main() {}"), IndexOptions{})
	if err == nil {
		t.Fatal("expected error for .rs")
	}
	var skip *SkipError
	if !errors.As(err, &skip) {
		t.Fatalf("want SkipError, got %T %v", err, err)
	}
	msg := skip.Reason
	for _, sub := range []string{"unsupported", "INDEX_LANG_POLICY", "Tier-2"} {
		if !strings.Contains(msg, sub) {
			t.Fatalf("reason %q missing %q", msg, sub)
		}
	}
}
