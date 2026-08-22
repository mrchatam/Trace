package domain

import (
	"context"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/mrchatam/Trace/internal/store"
)

const (
	WorkConflictReasonPathOverlap     = "path_overlap"
	WorkConflictReasonSeedOverlap     = "seed_overlap"
	WorkConflictReasonTitleRedundancy = "title_redundancy"

	minTitleSimilarRunes = 8
)

var activeWorkStates = map[string]struct{}{
	store.WorkStatePending:        {},
	store.WorkStateInProgress:     {},
	store.WorkStateAwaitingReview: {},
	store.WorkStateBlocked:        {},
}

// WorkConflict is an advisory overlap between two active tasks.
type WorkConflict struct {
	TaskA       string          `json:"task_a"`
	TaskB       string          `json:"task_b"`
	Reason      string          `json:"reason"`
	Paths       []string        `json:"paths,omitempty"`
	SeedOverlap []ImpactSeedRef `json:"seed_overlap,omitempty"`
}

// DetectWorkConflictsOpts controls advisory conflict detection.
type DetectWorkConflictsOpts struct {
	TaskID string // optional: only conflicts involving this task
	Limit  int    // max entries; 0 = unlimited
}

type taskConflictProfile struct {
	id       string
	title    string
	goalID   string
	paths    []string
	seedKeys map[string]ImpactSeedRef
}

// DetectWorkConflicts returns bounded, stably sorted advisory conflicts between active tasks.
func (s *Service) DetectWorkConflicts(ctx context.Context, opts DetectWorkConflictsOpts) ([]WorkConflict, error) {
	_ = ctx
	tasks, err := s.store.ListTasks()
	if err != nil {
		return nil, err
	}
	filterID := strings.TrimSpace(opts.TaskID)

	var active []store.Task
	for _, t := range tasks {
		if _, ok := activeWorkStates[t.WorkState]; !ok {
			continue
		}
		if filterID != "" && t.ID != filterID {
			// keep all active tasks for pairwise scan; filter applied after detection
		}
		active = append(active, t)
	}

	profiles := make([]taskConflictProfile, 0, len(active))
	for _, t := range active {
		paths, err := s.taskChangePaths(t.ID)
		if err != nil {
			return nil, err
		}
		seeds := s.seedsFromPathsBestEffort(paths)
		seedKeys := make(map[string]ImpactSeedRef, len(seeds))
		for _, seed := range seeds {
			seedKeys[EntityKey(seed.EntityType, seed.EntityID)] = seed
		}
		goalID := ""
		if t.GoalID != nil {
			goalID = *t.GoalID
		}
		profiles = append(profiles, taskConflictProfile{
			id:       t.ID,
			title:    t.Title,
			goalID:   goalID,
			paths:    paths,
			seedKeys: seedKeys,
		})
	}

	var conflicts []WorkConflict
	for i := 0; i < len(profiles); i++ {
		for j := i + 1; j < len(profiles); j++ {
			a, b := profiles[i], profiles[j]
			if c, ok := detectPairConflict(a, b); ok {
				if filterID != "" && c.TaskA != filterID && c.TaskB != filterID {
					continue
				}
				conflicts = append(conflicts, c)
			}
		}
	}

	sortWorkConflicts(conflicts)
	if opts.Limit > 0 && len(conflicts) > opts.Limit {
		conflicts = conflicts[:opts.Limit]
	}
	if conflicts == nil {
		conflicts = []WorkConflict{}
	}
	return conflicts, nil
}

func detectPairConflict(a, b taskConflictProfile) (WorkConflict, bool) {
	taskA, taskB := a.id, b.id
	if taskA > taskB {
		taskA, taskB = taskB, taskA
		a, b = b, a
	}

	overlappingPaths := intersectingPaths(a.paths, b.paths)
	seedOverlap := intersectingSeeds(a.seedKeys, b.seedKeys)
	titleRedundant := titlesRedundant(a, b)

	if len(overlappingPaths) == 0 && len(seedOverlap) == 0 && !titleRedundant {
		return WorkConflict{}, false
	}

	reasons := make([]string, 0, 3)
	if len(overlappingPaths) > 0 {
		reasons = append(reasons, WorkConflictReasonPathOverlap)
	}
	if len(seedOverlap) > 0 {
		reasons = append(reasons, WorkConflictReasonSeedOverlap)
	}
	if titleRedundant {
		reasons = append(reasons, WorkConflictReasonTitleRedundancy)
	}

	c := WorkConflict{
		TaskA:  taskA,
		TaskB:  taskB,
		Reason: strings.Join(reasons, ";"),
	}
	if len(overlappingPaths) > 0 {
		c.Paths = overlappingPaths
	}
	if len(seedOverlap) > 0 {
		c.SeedOverlap = seedOverlap
	}
	return c, true
}

func (s *Service) taskChangePaths(taskID string) ([]string, error) {
	changes, err := s.store.ListChangesByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	hasOpen := false
	for _, c := range changes {
		if c.Status == store.ChangeStatusSuperseded {
			continue
		}
		if c.Status == store.ChangeStatusOpen {
			hasOpen = true
			break
		}
	}

	seen := map[string]struct{}{}
	var paths []string
	for _, c := range changes {
		if c.Status == store.ChangeStatusSuperseded {
			continue
		}
		if hasOpen && c.Status != store.ChangeStatusOpen {
			continue
		}
		cpaths, err := s.store.ListChangePaths(c.ID)
		if err != nil {
			return nil, err
		}
		for _, p := range cpaths {
			path := store.NormalizePath(strings.TrimSpace(p.Path))
			if path == "" {
				continue
			}
			if _, ok := seen[path]; ok {
				continue
			}
			seen[path] = struct{}{}
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func (s *Service) seedsFromPathsBestEffort(paths []string) []ImpactSeedRef {
	seen := map[string]struct{}{}
	var seeds []ImpactSeedRef
	for _, path := range paths {
		batch, err := s.seedsFromChangePaths([]store.ChangePath{{Path: path}})
		if err != nil {
			continue
		}
		for _, seed := range batch {
			k := EntityKey(seed.EntityType, seed.EntityID)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			seeds = append(seeds, seed)
		}
	}
	sort.Slice(seeds, func(i, j int) bool {
		ki := EntityKey(seeds[i].EntityType, seeds[i].EntityID)
		kj := EntityKey(seeds[j].EntityType, seeds[j].EntityID)
		return ki < kj
	})
	return seeds
}

func intersectingPaths(aPaths, bPaths []string) []string {
	var out []string
	seen := map[string]struct{}{}
	for _, pa := range aPaths {
		for _, pb := range bPaths {
			if !pathsOverlap(pa, pb) {
				continue
			}
			key := pa + "\x00" + pb
			if pa > pb {
				key = pb + "\x00" + pa
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, pa)
			if pb != pa {
				out = append(out, pb)
			}
		}
	}
	sort.Strings(out)
	if out == nil {
		return []string{}
	}
	return out
}

func pathsOverlap(a, b string) bool {
	a = store.NormalizePath(a)
	b = store.NormalizePath(b)
	if a == b {
		return true
	}
	if strings.HasPrefix(b, a+"/") {
		return true
	}
	if strings.HasPrefix(a, b+"/") {
		return true
	}
	return false
}

func intersectingSeeds(aKeys, bKeys map[string]ImpactSeedRef) []ImpactSeedRef {
	var keys []string
	for k := range aKeys {
		if _, ok := bKeys[k]; ok {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	out := make([]ImpactSeedRef, 0, len(keys))
	for _, k := range keys {
		out = append(out, aKeys[k])
	}
	return out
}

func titlesRedundant(a, b taskConflictProfile) bool {
	if a.goalID == "" || b.goalID == "" || a.goalID != b.goalID {
		return false
	}
	return titlesSimilar(a.title, b.title)
}

// normalizeTitle lowercases, trims, and collapses whitespace for title comparison.
func normalizeTitle(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return ""
	}
	var b strings.Builder
	space := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !space && b.Len() > 0 {
				b.WriteByte(' ')
				space = true
			}
			continue
		}
		b.WriteRune(r)
		space = false
	}
	return strings.TrimSpace(b.String())
}

func titlesSimilar(a, b string) bool {
	na := normalizeTitle(a)
	nb := normalizeTitle(b)
	if na == "" || nb == "" {
		return false
	}
	if na == nb {
		return true
	}
	return mutualSubstring(na, nb, minTitleSimilarRunes)
}

func mutualSubstring(a, b string, minRunes int) bool {
	if containsSubstringMinRunes(a, b, minRunes) || containsSubstringMinRunes(b, a, minRunes) {
		return true
	}
	runesA := []rune(a)
	for i := 0; i <= len(runesA)-minRunes; i++ {
		for j := i + minRunes; j <= len(runesA); j++ {
			sub := string(runesA[i:j])
			if strings.Contains(b, sub) {
				return true
			}
		}
	}
	return false
}

func containsSubstringMinRunes(haystack, needle string, minRunes int) bool {
	if utf8.RuneCountInString(needle) < minRunes {
		return false
	}
	return strings.Contains(haystack, needle)
}

func sortWorkConflicts(conflicts []WorkConflict) {
	sort.Slice(conflicts, func(i, j int) bool {
		if conflicts[i].TaskA != conflicts[j].TaskA {
			return conflicts[i].TaskA < conflicts[j].TaskA
		}
		return conflicts[i].TaskB < conflicts[j].TaskB
	})
}
