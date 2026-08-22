package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPeelServeProjectFlagsAfterSubcommand(t *testing.T) {
	dir := t.TempDir()
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, rest, err := peelServeProjectFlags(cwd, []string{"-C", dir})
	if err != nil {
		t.Fatal(err)
	}
	if root != abs {
		t.Fatalf("root=%q want %q", root, abs)
	}
	if len(rest) != 0 {
		t.Fatalf("rest=%v want empty", rest)
	}
}

func TestServeAcceptsProjectFlagAfterSubcommand(t *testing.T) {
	dir := t.TempDir()
	code, _, stderr := runCapture(t, []string{"serve", "-C", dir, "--help"})
	if code != exitOK {
		t.Fatalf("serve -C after subcommand --help: exit=%d stderr=%s", code, stderr)
	}
}
