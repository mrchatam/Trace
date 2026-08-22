package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/mrchatam/Trace/internal/analyzers"
)

func TestIndexStatusSupportedLanguages(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	statusJSON := captureStdout(t, func() int {
		return run([]string{"-C", dir, "index", "status"})
	})
	var status indexStatusJSON
	if err := json.Unmarshal([]byte(statusJSON), &status); err != nil {
		t.Fatalf("status json: %v\n%s", err, statusJSON)
	}
	want := analyzers.SupportedLanguages()
	if !reflect.DeepEqual(status.SupportedLanguages, want) {
		t.Fatalf("supported_languages=%v want %v", status.SupportedLanguages, want)
	}
	if len(status.SupportedLanguages) != 5 {
		t.Fatalf("len=%d want 5", len(status.SupportedLanguages))
	}
}
