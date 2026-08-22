package retrieval

import (
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// Engine is the hybrid retrieval API bound to a store (never opens its own DB).
type Engine struct {
	store *store.Store
	vcs   vcs.Repository // optional
}

// New constructs an Engine. st must be non-nil and already opened.
func New(st *store.Store) *Engine {
	if st == nil {
		panic("retrieval: New: store is nil")
	}
	return &Engine{store: st}
}

// WithVCS attaches an optional VCS repository for temporal enrich (refs only).
func (e *Engine) WithVCS(repo vcs.Repository) *Engine {
	e.vcs = repo
	return e
}
