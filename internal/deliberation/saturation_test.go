package deliberation_test

import (
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
)

func TestSaturationEmptyThreshold(t *testing.T) {
	if deliberation.SaturationEmptyThreshold != 2 {
		t.Fatalf("SaturationEmptyThreshold=%d want 2", deliberation.SaturationEmptyThreshold)
	}
	if deliberation.HopBudget != 12 {
		t.Fatalf("HopBudget=%d want 12", deliberation.HopBudget)
	}
}

func TestNextConsecutiveEmptyApplies(t *testing.T) {
	tests := []struct {
		name                     string
		prev, plan, spawn, discW int
		want                     int
	}{
		{name: "first pure empty", prev: 0, want: 1},
		{name: "second pure empty", prev: 1, want: 2},
		{name: "plan clears", prev: 2, plan: 1, want: 0},
		{name: "spawn clears", prev: 2, spawn: 1, want: 0},
		{name: "discoveries-only no increment", prev: 1, discW: 1, want: 1},
		{name: "discoveries-only from zero", prev: 0, discW: 2, want: 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := deliberation.NextConsecutiveEmptyApplies(tc.prev, tc.plan, tc.spawn, tc.discW)
			if got != tc.want {
				t.Fatalf("got %d want %d", got, tc.want)
			}
		})
	}
}

func TestSaturatedFromCounter(t *testing.T) {
	if deliberation.SaturatedFromCounter(false, 1) {
		t.Fatal("consecutive=1 must not saturate")
	}
	if !deliberation.SaturatedFromCounter(false, 2) {
		t.Fatal("consecutive=2 must saturate")
	}
	if !deliberation.SaturatedFromCounter(true, 0) {
		t.Fatal("max iterations must saturate")
	}
}
