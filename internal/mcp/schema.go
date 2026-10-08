package mcp

import (
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

// The SDK infers a tool's input schema from the Go input struct and reads the
// `jsonschema` struct tag as the property *description* only — a tag that
// begins with `WORD=` is rejected outright. Enums and numeric ranges therefore
// cannot come from struct tags; they must be supplied as an explicit
// Tool.InputSchema (#158).
//
// schemaFor infers the struct schema exactly as the SDK would, then applies the
// property constraints declared here, so descriptions keep coming from the
// tags while validation keywords come from code.

// propConstraint declares enum and/or numeric bounds for one input property.
type propConstraint struct {
	enum    []string
	minimum *float64
	maximum *float64
}

func enumOf(values ...string) propConstraint {
	return propConstraint{enum: values}
}

func rangeOf(min, max float64) propConstraint {
	return propConstraint{minimum: &min, maximum: &max}
}

func (c propConstraint) and(other propConstraint) propConstraint {
	if other.enum != nil {
		c.enum = other.enum
	}
	if other.minimum != nil {
		c.minimum = other.minimum
	}
	if other.maximum != nil {
		c.maximum = other.maximum
	}
	return c
}

// schemaFor builds the input schema for In with the given property constraints.
// It panics on a missing property or an invalid schema, because both are
// programming errors that must surface at server construction, not at call time.
func schemaFor[In any](constraints map[string]propConstraint) *jsonschema.Schema {
	schema, err := jsonschema.For[In](nil)
	if err != nil {
		panic(fmt.Sprintf("mcp: schema inference for %v: %v", reflect.TypeFor[In](), err))
	}
	if schema == nil {
		panic(fmt.Sprintf("mcp: nil schema for %v", reflect.TypeFor[In]()))
	}
	for name, c := range constraints {
		prop, ok := schema.Properties[name]
		if !ok || prop == nil {
			panic(fmt.Sprintf("mcp: no property %q on %v", name, reflect.TypeFor[In]()))
		}
		for _, v := range c.enum {
			prop.Enum = append(prop.Enum, v)
		}
		prop.Minimum = c.minimum
		prop.Maximum = c.maximum
	}
	if _, err := schema.Resolve(&jsonschema.ResolveOptions{ValidateDefaults: true}); err != nil {
		panic(fmt.Sprintf("mcp: schema for %v does not resolve: %v", reflect.TypeFor[In](), err))
	}
	return schema
}

// Shared vocabularies. Tool descriptions are generated from these so a
// description and its schema enum cannot drift apart (#159).

var (
	capabilityActions = []string{"declare", "list", "require", "unrequire", "missing"}
	capabilityKinds   = []string{"SKILL", "RULE", "MCP", "TOOL", "HOOK", "AGENT"}
	capabilityStatus  = []string{"AVAILABLE", "UNAVAILABLE", "UNKNOWN"}
	reviewActions     = []string{"create", "set"}
	reviewResults     = []string{"PASS", "FAIL", "UNCERTAIN"}
	changeActions     = []string{"list", "show", "compare"}
	regressionAction  = []string{"list"}
	agentActions      = []string{"list", "describe", "recommend"}
	planActions       = []string{"create-coarse", "set-current", "deep", "show", "bootstrap"}
	impactActions     = []string{"finding", "alternative", "report", "walk"}
	impactOps         = []string{"add", "list", "recommend"}
	loopActions       = []string{"next", "apply", "status", "gate"}
	loopGateFors      = []string{"orient", "edit", "execute", "done", "export"}
	linkRels          = []string{
		"goal-task", "decision-task", "discovery-plan-change", "discovery-mentions-task",
		"claim-evidence", "scope-member", "api-contract", "implements", "blocks",
	}
	contextFormats = []string{"json", "markdown", "both"}
	matchers       = []string{"and", "ops"}
	workStates     = []string{
		"PENDING", "IN_PROGRESS", "AWAITING_REVIEW", "BLOCKED", "FAILED", "DONE", "STALE", "SKIPPED",
	}
	addKinds = []string{
		"discovery", "task", "goal", "decision", "assumption", "plan-change", "claim", "evidence",
	}
	provenanceStatus = []string{"ACTIVE", "STALE", "SUPERSEDED"}
)

// join renders a vocabulary for use in a description or a schema tag.
func join(values []string) string {
	out := ""
	for i, v := range values {
		if i > 0 {
			out += "|"
		}
		out += v
	}
	return out
}
