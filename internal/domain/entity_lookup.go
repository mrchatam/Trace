package domain

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// AddEntityKind is the create-time kind vocabulary shared by trace_add and
// `trace add`. It is the wire form; EntityAddKinds lists the accepted values.
const (
	AddKindGoal       = "goal"
	AddKindTask       = "task"
	AddKindDecision   = "decision"
	AddKindAssumption = "assumption"
	AddKindDiscovery  = "discovery"
	AddKindPlanChange = "plan-change"
	AddKindClaim      = "claim"
	AddKindEvidence   = "evidence"
)

// EntityAddKinds is the closed create-kind vocabulary, in the order the tool
// surfaces it.
var EntityAddKinds = []string{
	AddKindGoal, AddKindTask, AddKindDecision, AddKindAssumption,
	AddKindDiscovery, AddKindPlanChange, AddKindClaim, AddKindEvidence,
}

// NormalizeAddKind validates a create kind against EntityAddKinds.
func NormalizeAddKind(kind string) (string, error) {
	k := strings.ToLower(strings.TrimSpace(kind))
	for _, want := range EntityAddKinds {
		if k == want {
			return k, nil
		}
	}
	return "", &ErrValidation{Msg: "kind must be one of " + strings.Join(EntityAddKinds, "|")}
}

// ValidateEntityID checks that an explicitly supplied entity id is a UUID.
// An empty id is valid and means "allocate one"; a non-empty id that is not a
// UUID is a caller error, because the schema stores ids as opaque text and
// would otherwise store a dangling reference (#151).
func ValidateEntityID(field, id string) (string, error) {
	v := strings.TrimSpace(id)
	if v == "" {
		return "", nil
	}
	if _, err := uuid.Parse(v); err != nil {
		return "", &ErrValidation{Msg: field + " must be a UUID, got " + strconv.Quote(v)}
	}
	return v, nil
}

// entityGetter resolves the store lookup for one create kind.
type entityGetter struct {
	entityType string
	get        func(string) error
}

// addKindGetters maps each create kind to the row it writes, so existence
// probing uses the same vocabulary the write path does.
func (s *Service) addKindGetters() map[string]entityGetter {
	return map[string]entityGetter{
		AddKindGoal:       {EntityGoal, func(id string) error { _, err := s.store.GetGoal(id); return err }},
		AddKindTask:       {EntityTask, func(id string) error { _, err := s.store.GetTask(id); return err }},
		AddKindDecision:   {EntityDecision, func(id string) error { _, err := s.store.GetDecision(id); return err }},
		AddKindAssumption: {EntityAssumption, func(id string) error { _, err := s.store.GetAssumption(id); return err }},
		AddKindDiscovery:  {EntityDiscovery, func(id string) error { _, err := s.store.GetDiscovery(id); return err }},
		AddKindPlanChange: {EntityPlanChange, func(id string) error { _, err := s.store.GetPlanChange(id); return err }},
		AddKindClaim:      {EntityClaim, func(id string) error { _, err := s.store.GetClaim(id); return err }},
		AddKindEvidence:   {EntityEvidence, func(id string) error { _, err := s.store.GetEvidence(id); return err }},
	}
}

// EntityExists reports whether an id is already a row of the given create kind.
// A lookup error other than "not found" is returned so a real store failure is
// never mistaken for a free id.
func (s *Service) EntityExists(ctx context.Context, kind, id string) (bool, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return false, nil
	}
	g, ok := s.addKindGetters()[kind]
	if !ok {
		return false, &ErrValidation{Msg: "unknown kind " + strconv.Quote(kind)}
	}
	err := g.get(id)
	if err == nil {
		return true, nil
	}
	if isNotFoundErr(err) || errors.Is(err, store.ErrNotFound) {
		return false, nil
	}
	return false, err
}
