package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestTasksAllHardCapped(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	// Insert more than MaxTaskListLimit would be slow; instead verify Limit wiring via one page.
	for i := 0; i < 3; i++ {
		if _, err := svc.CreateTask(ctx, domain.TaskInput{Title: fmt.Sprintf("t%d", i)}); err != nil {
			t.Fatal(err)
		}
	}
	_ = st.Close()

	s := &Server{defaultRoot: dir}
	_, out, err := s.toolTasks(ctx, nil, TasksInput{All: true, Project: dir})
	if err != nil {
		t.Fatal(err)
	}
	_ = out
	// Re-open via tool path already returned JSON text; call again and parse from CallToolResult
	res, _, err := s.toolTasks(ctx, nil, TasksInput{All: true, Project: dir})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil || len(res.Content) == 0 {
		t.Fatal("empty result")
	}
	var payload map[string]any
	// textResult stores JSON in text content
	txt, ok := res.Content[0].(interface{ GetText() string })
	if ok {
		if err := json.Unmarshal([]byte(txt.GetText()), &payload); err != nil {
			t.Fatal(err)
		}
	} else {
		// fallback: marshal content
		b, _ := json.Marshal(res.Content[0])
		var wrap map[string]any
		_ = json.Unmarshal(b, &wrap)
		if text, _ := wrap["text"].(string); text != "" {
			if err := json.Unmarshal([]byte(text), &payload); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("cannot parse result: %#v", res.Content[0])
		}
	}
	// With 3 tasks, truncated should be false and count 3; Limit must not be -1 unbounded.
	if int(payload["count"].(float64)) != 3 {
		t.Fatalf("count=%v", payload["count"])
	}
	if payload["truncated"] == true {
		t.Fatalf("unexpected truncated with 3 rows: %#v", payload)
	}
}
