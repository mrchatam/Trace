package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestCapabilityListMissingSnakeCase(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	decl := captureStdout(t, func() int {
		return run([]string{"-C", dir, "capability", "declare",
			"--kind", "MCP", "--slug", "mcp:demo", "--title", "Demo", "--status", "UNAVAILABLE"})
	})
	var declOut map[string]any
	if err := json.Unmarshal([]byte(decl), &declOut); err != nil {
		t.Fatalf("declare json: %v\n%s", err, decl)
	}
	capID, _ := declOut["id"].(string)
	if capID == "" {
		t.Fatalf("declare missing id: %s", decl)
	}

	listOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "capability", "list"})
	})
	assertCapabilitySnakeCaseRows(t, listOut, "capabilities")
	if strings.Contains(listOut, `"ID"`) || strings.Contains(listOut, `"Kind"`) {
		t.Fatalf("list must not emit PascalCase keys: %s", listOut)
	}

	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "T"})
	})
	var taskMap map[string]any
	if err := json.Unmarshal([]byte(taskOut), &taskMap); err != nil {
		t.Fatalf("add task: %v\n%s", err, taskOut)
	}
	taskID, _ := taskMap["id"].(string)
	if taskID == "" {
		t.Fatalf("task id: %s", taskOut)
	}

	_ = captureStdout(t, func() int {
		return run([]string{"-C", dir, "capability", "require", "--task", taskID, "--capability", "mcp:demo"})
	})

	missOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "capability", "missing", "--task", taskID})
	})
	assertCapabilitySnakeCaseRows(t, missOut, "missing")
	if strings.Contains(missOut, `"ID"`) || strings.Contains(missOut, `"Slug"`) {
		t.Fatalf("missing must not emit PascalCase keys: %s", missOut)
	}

	// Absolute path sanity for -C (seed/init already used temp).
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	_ = captureStdout(t, func() int {
		return run([]string{"-C", abs, "capability", "list", "--kind", "MCP"})
	})
}

func assertCapabilitySnakeCaseRows(t *testing.T, raw, field string) {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		t.Fatalf("%s envelope: %v\n%s", field, err, raw)
	}
	if env["ok"] != true {
		t.Fatalf("%s ok: %s", field, raw)
	}
	arr, ok := env[field].([]any)
	if !ok || len(arr) < 1 {
		t.Fatalf("%s want non-empty array: %s", field, raw)
	}
	row, ok := arr[0].(map[string]any)
	if !ok {
		t.Fatalf("%s row type: %T", field, arr[0])
	}
	for _, k := range []string{"id", "kind", "slug", "title", "status"} {
		if _, ok := row[k]; !ok {
			t.Fatalf("%s missing key %q in %v", field, k, row)
		}
	}
	for k := range row {
		switch k {
		case "id", "kind", "slug", "title", "status":
		default:
			t.Fatalf("%s unexpected key %q (body/timestamps must be omitted)", field, k)
		}
	}
}
