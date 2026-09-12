package store

import (
	"fmt"
	"testing"
)

func TestListTasksFilteredLimitAndCursor(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	for i := 0; i < 5; i++ {
		ws := WorkStatePending
		if i%2 == 0 {
			ws = WorkStateDone
		}
		if _, err := s.UpsertTask(Task{Title: fmt.Sprintf("t%d", i), WorkState: ws}); err != nil {
			t.Fatal(err)
		}
	}

	page1, err := s.ListTasksFiltered(TaskListFilter{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1.Tasks) != 2 || !page1.Truncated || page1.NextCursor == "" {
		t.Fatalf("page1=%+v", page1)
	}
	// list rows omit body
	if page1.Tasks[0].Body != "" {
		t.Fatalf("list should omit body, got %q", page1.Tasks[0].Body)
	}

	page2, err := s.ListTasksFiltered(TaskListFilter{Limit: 2, Cursor: page1.NextCursor})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.Tasks) != 2 {
		t.Fatalf("page2 len=%d", len(page2.Tasks))
	}
	if page1.Tasks[0].ID == page2.Tasks[0].ID {
		t.Fatal("cursor did not advance")
	}

	done, err := s.ListTasksFiltered(TaskListFilter{WorkStates: []string{WorkStateDone}, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(done.Tasks) != 3 { // i=0,2,4
		t.Fatalf("done count=%d", len(done.Tasks))
	}
	for _, task := range done.Tasks {
		if task.WorkState != WorkStateDone {
			t.Fatalf("want DONE got %s", task.WorkState)
		}
	}

	all, err := s.ListTasksFiltered(TaskListFilter{Limit: -1, IncludeBody: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Tasks) != 5 || all.Truncated {
		t.Fatalf("all=%+v", all)
	}
}
