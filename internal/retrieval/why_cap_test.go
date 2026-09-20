package retrieval_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
)

// TestWhyStepsCappedAndTruncated covers #106: dense expand neighborhoods must
// hard-cap Why steps and set truncated=true honestly.
func TestWhyStepsCappedAndTruncated(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "dense goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "dense task", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < retrieval.MaxWhySteps+20; i++ {
		d, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: fmt.Sprintf("disc-%02d", i)})
		if err != nil {
			t.Fatal(err)
		}
		if err := svc.LinkDiscoveryMentionsTask(ctx, d.ID, task.ID, domain.LinkMeta{}); err != nil {
			t.Fatal(err)
		}
	}

	why, err := eng.Why(ctx, "task", task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(why.Steps) > retrieval.MaxWhySteps {
		t.Fatalf("steps=%d want <= %d", len(why.Steps), retrieval.MaxWhySteps)
	}
	if len(why.Steps) != retrieval.MaxWhySteps {
		t.Fatalf("steps=%d want exactly %d (cap filled)", len(why.Steps), retrieval.MaxWhySteps)
	}
	if !why.Truncated {
		t.Fatal("expected truncated=true when expand exceeds MaxWhySteps")
	}
	if why.Steps[0].ReasonCode != retrieval.ReasonExactID {
		t.Fatalf("seed step missing: %+v", why.Steps[0])
	}
}

// TestWhyMinimalUntruncated covers #106: small graphs must not claim truncation.
func TestWhyMinimalUntruncated(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()
	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "tiny"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "tiny task", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}
	why, err := eng.Why(ctx, "task", task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if why.Truncated {
		t.Fatalf("minimal why must not be truncated: steps=%d", len(why.Steps))
	}
	if len(why.Steps) == 0 {
		t.Fatal("want at least seed step")
	}
}
