package domain

import "testing"

func TestNormalizeTitle(t *testing.T) {
	got := normalizeTitle("  Implement   Conflict\tDetection  ")
	want := "implement conflict detection"
	if got != want {
		t.Fatalf("normalizeTitle=%q want %q", got, want)
	}
}
