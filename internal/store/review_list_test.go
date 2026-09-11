package store

import (
	"fmt"
	"testing"
)

func TestListReviewsFilteredLimitAndCursor(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	task, err := s.UpsertTask(Task{Title: "linked", WorkState: WorkStatePending})
	if err != nil {
		t.Fatal(err)
	}
	var linkedIDs []string
	for i := 0; i < 5; i++ {
		r, err := s.UpsertReview(Review{Title: fmt.Sprintf("r%d", i), Body: fmt.Sprintf("body-%d", i)})
		if err != nil {
			t.Fatal(err)
		}
		if i%2 == 0 {
			if _, err := s.InsertLink(EntityLink{
				FromType: "review", FromID: r.ID, Rel: "review_judges_task",
				ToType: "task", ToID: task.ID,
			}); err != nil {
				t.Fatal(err)
			}
			linkedIDs = append(linkedIDs, r.ID)
		}
	}

	page1, err := s.ListReviewsFiltered(ReviewListFilter{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1.Reviews) != 2 || !page1.Truncated || page1.NextCursor == "" {
		t.Fatalf("page1=%+v", page1)
	}
	if page1.Reviews[0].Body != "" {
		t.Fatalf("list should omit body, got %q", page1.Reviews[0].Body)
	}

	page2, err := s.ListReviewsFiltered(ReviewListFilter{Limit: 2, Cursor: page1.NextCursor})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.Reviews) != 2 {
		t.Fatalf("page2 len=%d", len(page2.Reviews))
	}
	if page1.Reviews[0].ID == page2.Reviews[0].ID {
		t.Fatal("cursor did not advance")
	}

	byTask, err := s.ListReviewsFiltered(ReviewListFilter{TaskID: task.ID, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(byTask.Reviews) != len(linkedIDs) {
		t.Fatalf("by task count=%d want %d", len(byTask.Reviews), len(linkedIDs))
	}

	full, err := s.ListReviewsFiltered(ReviewListFilter{Limit: -1, IncludeBody: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(full.Reviews) != 5 || full.Truncated {
		t.Fatalf("full=%+v", full)
	}
	if full.Reviews[0].Body == "" {
		t.Fatal("IncludeBody should return body")
	}
}
