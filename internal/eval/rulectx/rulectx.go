package rulectx

import "context"

// InvariantRules exposes invariant enablement for built-in mechanisms.
type InvariantRules interface {
	IsInvariantEnabled(id string) bool
}

type rulesKey struct{}

// With attaches invariant rules to ctx.
func With(ctx context.Context, rules InvariantRules) context.Context {
	if rules == nil {
		return ctx
	}
	return context.WithValue(ctx, rulesKey{}, rules)
}

// From returns invariant rules previously attached with With.
func From(ctx context.Context) InvariantRules {
	v := ctx.Value(rulesKey{})
	if v == nil {
		return nil
	}
	r, ok := v.(InvariantRules)
	if !ok {
		return nil
	}
	return r
}
