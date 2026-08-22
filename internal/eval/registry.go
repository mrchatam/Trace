package eval

import (
	"sort"
	"sync"
)

// Built-in mechanism ids (locked in S07 planner).
const (
	MechanismStoredTest             = "stored_test"
	MechanismStoredVerification     = "stored_verification"
	MechanismStoredEvaluation       = "stored_evaluation"
	MechanismArchitecturalInvariant = "architectural_invariant"
)

var (
	defaultRegistry     *Registry
	defaultRegistryOnce sync.Once
)

// Registry holds registered evaluation mechanisms.
type Registry struct {
	mu   sync.RWMutex
	byID map[string]Mechanism
}

// DefaultRegistry returns the process-wide mechanism registry.
func DefaultRegistry() *Registry {
	defaultRegistryOnce.Do(func() {
		defaultRegistry = &Registry{byID: map[string]Mechanism{}}
	})
	return defaultRegistry
}

// Register adds a mechanism to the default registry (init-time or tests).
func Register(m Mechanism) {
	DefaultRegistry().register(m)
}

func (r *Registry) register(m Mechanism) {
	if m == nil {
		return
	}
	id := m.ID()
	if id == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.byID == nil {
		r.byID = map[string]Mechanism{}
	}
	r.byID[id] = m
}

// ListMechanismIDs returns registered ids in stable sorted order.
func (r *Registry) ListMechanismIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.byID))
	for id := range r.byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func (r *Registry) mechanism(id string) (Mechanism, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.byID[id]
	return m, ok
}
