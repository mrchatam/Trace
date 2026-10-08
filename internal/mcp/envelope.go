package mcp

// listEnvelope is the one list-shaped response shape every MCP list action
// uses, so a caller never has to guess whether a bare array or an object came
// back and always learns whether the list was truncated (#157).
//
// Field order in the emitted JSON follows the order added here.
type listEnvelope struct {
	OK         bool   `json:"ok"`
	Items      any    `json:"items"`
	Count      int    `json:"count"`
	Truncated  bool   `json:"truncated"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	Hint       string `json:"hint,omitempty"`
}

// newListEnvelope builds the canonical envelope. items must be a non-nil
// slice so an empty result marshals as [] rather than null.
func newListEnvelope(items any, count int, truncated bool) listEnvelope {
	return listEnvelope{OK: true, Items: items, Count: count, Truncated: truncated}
}

// withLimit records the effective (clamped) limit actually applied.
func (e listEnvelope) withLimit(limit int) listEnvelope {
	e.Limit = limit
	return e
}

// withCursor records the opaque cursor for the next page.
func (e listEnvelope) withCursor(cursor string) listEnvelope {
	e.NextCursor = cursor
	return e
}

// withHint records an actionable note, e.g. an empty-catalog remedy.
func (e listEnvelope) withHint(hint string) listEnvelope {
	e.Hint = hint
	return e
}
