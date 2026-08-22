package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// ImpactSeedRef is a structural walk seed (file|symbol only).
type ImpactSeedRef struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// ImpactEntityRef is one blast node (keys only — no paths or source text).
type ImpactEntityRef struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// ImpactWalkSnapshot is the domain-facing impact walk outcome for predict/compare.
type ImpactWalkSnapshot struct {
	Seeds         []ImpactSeedRef   `json:"seeds"`
	Blast         []ImpactEntityRef `json:"blast"`
	AffectedTests []ImpactEntityRef `json:"affected_tests"`
	BlastTotal    int               `json:"blast_total"`
	BlastKept     int               `json:"blast_kept"`
	Truncated     bool              `json:"truncated"`
	Depth         int               `json:"depth"`
}

// ImpactWalker runs a graph impact walk (implemented by retrieval; injected to avoid import cycle).
type ImpactWalker interface {
	ImpactWalk(ctx context.Context, seeds []ImpactSeedRef, depth int) (*ImpactWalkSnapshot, error)
}

// PredictedImpactPayload is the stored predict JSON (no source blobs).
type PredictedImpactPayload struct {
	Seeds            []ImpactSeedRef `json:"seeds"`
	BlastKeys        []string        `json:"blast_keys"`
	AffectedTestKeys []string        `json:"affected_test_keys"`
	Depth            int             `json:"depth"`
	BlastTotal       int             `json:"blast_total"`
	BlastKept        int             `json:"blast_kept"`
	Truncated        bool            `json:"truncated"`
}

// ImpactCompareDelta is the persisted compare result (sorted keys).
type ImpactCompareDelta struct {
	Matched    []string `json:"matched"`
	Unexpected []string `json:"unexpected"`
	Missed     []string `json:"missed"`
}

// ImpactCompareResult is returned from CompareActualImpact.
type ImpactCompareResult struct {
	ChangeID   string             `json:"change_id"`
	Delta      ImpactCompareDelta `json:"delta"`
	ComparedAt string             `json:"compared_at"`
}

// EntityKey formats a file|symbol id as a blast comparison key.
func EntityKey(entityType, entityID string) string {
	typ := strings.ToLower(strings.TrimSpace(entityType))
	return typ + ":" + strings.TrimSpace(entityID)
}

func blastKeysFromSnapshot(snap *ImpactWalkSnapshot) []string {
	keys := make([]string, 0, len(snap.Blast))
	for _, h := range snap.Blast {
		keys = append(keys, EntityKey(h.EntityType, h.EntityID))
	}
	sort.Strings(keys)
	return keys
}

func affectedTestKeysFromSnapshot(snap *ImpactWalkSnapshot) []string {
	keys := make([]string, 0, len(snap.AffectedTests))
	for _, h := range snap.AffectedTests {
		keys = append(keys, EntityKey(h.EntityType, h.EntityID))
	}
	sort.Strings(keys)
	return keys
}

func predictedKeysFromPayload(p PredictedImpactPayload) []string {
	seen := map[string]struct{}{}
	var keys []string
	for _, k := range p.BlastKeys {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		keys = append(keys, k)
	}
	for _, k := range p.AffectedTestKeys {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func actualKeysFromSnapshot(snap *ImpactWalkSnapshot) []string {
	seen := map[string]struct{}{}
	var keys []string
	for _, k := range blastKeysFromSnapshot(snap) {
		seen[k] = struct{}{}
		keys = append(keys, k)
	}
	for _, k := range affectedTestKeysFromSnapshot(snap) {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func diffKeySets(predicted, actual []string) ImpactCompareDelta {
	pSet := map[string]struct{}{}
	for _, k := range predicted {
		pSet[k] = struct{}{}
	}
	aSet := map[string]struct{}{}
	for _, k := range actual {
		aSet[k] = struct{}{}
	}
	var matched, unexpected, missed []string
	for k := range pSet {
		if _, ok := aSet[k]; ok {
			matched = append(matched, k)
		} else {
			missed = append(missed, k)
		}
	}
	for k := range aSet {
		if _, ok := pSet[k]; !ok {
			unexpected = append(unexpected, k)
		}
	}
	sort.Strings(matched)
	sort.Strings(unexpected)
	sort.Strings(missed)
	if matched == nil {
		matched = []string{}
	}
	if unexpected == nil {
		unexpected = []string{}
	}
	if missed == nil {
		missed = []string{}
	}
	return ImpactCompareDelta{
		Matched:    matched,
		Unexpected: unexpected,
		Missed:     missed,
	}
}

func (s *Service) requireImpactWalker() (ImpactWalker, error) {
	if s.impactWalker == nil {
		return nil, &ErrValidation{Msg: "impact walker not configured"}
	}
	return s.impactWalker, nil
}

// SetImpactWalker wires retrieval (or tests) for predict/compare walks.
func (s *Service) SetImpactWalker(w ImpactWalker) {
	s.impactWalker = w
}

// RecordPredictedImpact serializes walk keys and upserts impact_predictions.
func (s *Service) RecordPredictedImpact(ctx context.Context, changeID string, walkResult *ImpactWalkSnapshot) (store.ImpactPrediction, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return store.ImpactPrediction{}, &ErrValidation{Msg: "change_id is required"}
	}
	if walkResult == nil {
		return store.ImpactPrediction{}, &ErrValidation{Msg: "walk result is required"}
	}
	if _, err := s.store.GetChange(changeID); err != nil {
		return store.ImpactPrediction{}, err
	}

	payload := PredictedImpactPayload{
		Seeds:            walkResult.Seeds,
		BlastKeys:        blastKeysFromSnapshot(walkResult),
		AffectedTestKeys: affectedTestKeysFromSnapshot(walkResult),
		Depth:            walkResult.Depth,
		BlastTotal:       walkResult.BlastTotal,
		BlastKept:        walkResult.BlastKept,
		Truncated:        walkResult.Truncated,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return store.ImpactPrediction{}, fmt.Errorf("domain: marshal predicted impact: %w", err)
	}
	return s.store.UpsertImpactPrediction(store.ImpactPrediction{
		ChangeID:      changeID,
		PredictedJSON: string(raw),
		Depth:         walkResult.Depth,
	})
}

func (s *Service) seedsFromChangePaths(paths []store.ChangePath) ([]ImpactSeedRef, error) {
	if len(paths) == 0 {
		return nil, &ErrValidation{Msg: "change has no paths"}
	}
	seen := map[string]struct{}{}
	var seeds []ImpactSeedRef
	for _, p := range paths {
		symID := strings.TrimSpace(p.SymbolID)
		if symID != "" {
			k := EntityKey("symbol", symID)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			seeds = append(seeds, ImpactSeedRef{EntityType: "symbol", EntityID: symID})
			continue
		}
		path := store.NormalizePath(strings.TrimSpace(p.Path))
		if path == "" {
			return nil, &ErrValidation{Msg: "change path is empty"}
		}
		f, err := s.store.GetFileByPath(path)
		if err != nil {
			return nil, &ErrValidation{Msg: fmt.Sprintf("path %q is not indexed", path)}
		}
		k := EntityKey("file", f.ID)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		seeds = append(seeds, ImpactSeedRef{EntityType: "file", EntityID: f.ID})
	}
	if len(seeds) == 0 {
		return nil, &ErrValidation{Msg: "no seeds from change paths"}
	}
	return seeds, nil
}

// PredictImpactForChange walks change paths, runs ImpactWalk, and stores the snapshot.
func (s *Service) PredictImpactForChange(ctx context.Context, changeID string, depth int) (store.ImpactPrediction, error) {
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return store.ImpactPrediction{}, &ErrValidation{Msg: "change_id is required"}
	}
	if depth < 1 || depth > 2 {
		return store.ImpactPrediction{}, &ErrValidation{Msg: "depth must be 1 or 2"}
	}
	walker, err := s.requireImpactWalker()
	if err != nil {
		return store.ImpactPrediction{}, err
	}
	paths, err := s.store.ListChangePaths(changeID)
	if err != nil {
		return store.ImpactPrediction{}, err
	}
	seeds, err := s.seedsFromChangePaths(paths)
	if err != nil {
		return store.ImpactPrediction{}, err
	}
	res, err := walker.ImpactWalk(ctx, seeds, depth)
	if err != nil {
		return store.ImpactPrediction{}, err
	}
	return s.RecordPredictedImpact(ctx, changeID, res)
}

// CompareActualImpact re-walks stored seeds+depth and diffs blast key sets.
func (s *Service) CompareActualImpact(ctx context.Context, changeID string) (ImpactCompareResult, error) {
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return ImpactCompareResult{}, &ErrValidation{Msg: "change_id is required"}
	}
	walker, err := s.requireImpactWalker()
	if err != nil {
		return ImpactCompareResult{}, err
	}
	row, err := s.store.GetImpactPrediction(changeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ImpactCompareResult{}, &ErrValidation{Msg: "no impact prediction for change"}
		}
		return ImpactCompareResult{}, err
	}
	var predicted PredictedImpactPayload
	if err := json.Unmarshal([]byte(row.PredictedJSON), &predicted); err != nil {
		return ImpactCompareResult{}, fmt.Errorf("domain: parse predicted_json: %w", err)
	}
	if len(predicted.Seeds) == 0 {
		return ImpactCompareResult{}, &ErrValidation{Msg: "stored prediction has no seeds"}
	}
	depth := row.Depth
	if depth < 1 {
		depth = predicted.Depth
	}
	if depth < 1 || depth > 2 {
		return ImpactCompareResult{}, &ErrValidation{Msg: "stored depth invalid"}
	}

	actual, err := walker.ImpactWalk(ctx, predicted.Seeds, depth)
	if err != nil {
		return ImpactCompareResult{}, err
	}
	delta := diffKeySets(predictedKeysFromPayload(predicted), actualKeysFromSnapshot(actual))
	raw, err := json.Marshal(delta)
	if err != nil {
		return ImpactCompareResult{}, fmt.Errorf("domain: marshal compare_json: %w", err)
	}
	updated, err := s.store.UpdateImpactPredictionCompare(changeID, string(raw), "")
	if err != nil {
		return ImpactCompareResult{}, err
	}
	return ImpactCompareResult{
		ChangeID:   changeID,
		Delta:      delta,
		ComparedAt: updated.ComparedAt,
	}, nil
}
