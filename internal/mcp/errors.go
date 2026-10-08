package mcp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// linkRelShape is the from/to entity-type expectation for each trace_link rel.
// It is what turns a bare "no rows" into a message that tells the caller which
// parameter is wrong (#155).
var linkRelShape = map[string]struct {
	from, to string
}{
	"goal-task":               {from: "goal", to: "task"},
	"decision-task":           {from: "decision", to: "task"},
	"discovery-plan-change":   {from: "discovery", to: "plan_change"},
	"discovery-mentions-task": {from: "discovery", to: "task"},
	"claim-evidence":          {from: "claim", to: "evidence"},
	"scope-member":            {from: "any linkable entity", to: "scope"},
	"scope_member":            {from: "any linkable entity", to: "scope"},
	"api-contract":            {from: "task", to: "task"},
	"api_contract":            {from: "task", to: "task"},
	"implements":              {from: "task or decision", to: "task"},
	"blocks":                  {from: "task", to: "task"},
}

// wrapToolErr is the single MCP error-mapping point. It renders a missing row
// with the entity type that was actually queried and never leaks driver text;
// every other error passes through untouched.
func wrapToolErr(tool string, err error) error {
	if err == nil {
		return nil
	}
	var nf *store.NotFoundError
	if errors.As(err, &nf) {
		return fmt.Errorf("%s: no %s %q", tool, nf.Entity, nf.ID)
	}
	return fmt.Errorf("%s: %w", tool, err)
}

// wrapLinkErr renders a trace_link failure with the rel's own vocabulary, so a
// wrong-type from/to is named instead of guessed from the called function.
func wrapLinkErr(rel, fromID, toID string, err error) error {
	if err == nil {
		return nil
	}
	var nf *store.NotFoundError
	if !errors.As(err, &nf) {
		return wrapToolErr("trace_link", err)
	}
	shape, known := linkRelShape[strings.ToLower(strings.TrimSpace(rel))]
	if !known {
		return fmt.Errorf("trace_link: no %s %q", nf.Entity, nf.ID)
	}
	// Name the parameter the caller actually got wrong, not the function's
	// fixed lookup order.
	role := "from"
	if nf.ID != strings.TrimSpace(fromID) && nf.ID == strings.TrimSpace(toID) {
		role = "to"
	}
	return fmt.Errorf(
		"trace_link: rel %s: no %s %q for %s (rel %s expects from=%s to=%s)",
		rel, nf.Entity, nf.ID, role, rel, shape.from, shape.to)
}
