package install

import (
	"fmt"
	"sort"
	"strings"
)

var registry = []Target{
	cursorTarget{},
	claudeTarget{},
	cursorHookTarget{},
	gitHookTarget{},
}

// Registry returns registered install targets in stable id order.
func Registry() []Target {
	out := make([]Target, len(registry))
	copy(out, registry)
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID() < out[j].ID()
	})
	return out
}

// ListTargets returns DetectInfo for every registered target (no writes).
func ListTargets(opts InstallOpts) []DetectInfo {
	targets := Registry()
	out := make([]DetectInfo, 0, len(targets))
	for _, t := range targets {
		d := t.Detect(opts)
		out = append(out, DetectInfo{
			ID:       t.ID(),
			Tier:     t.Tier(),
			Detected: d.Detected,
			Reason:   d.Reason,
		})
	}
	return out
}

// Lookup returns a target by id (case-sensitive slug).
func Lookup(id string) (Target, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("install: target id is required")
	}
	for _, t := range registry {
		if t.ID() == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("install: unknown target %q", id)
}
