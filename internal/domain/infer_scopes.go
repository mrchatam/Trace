package domain

import (
	"context"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Inference provenance (Phase 44 S03). Distinct from USER_ASSERTED / IMPORTED.
const (
	SourceTypeInferred     = "INFERRED"
	SourceTypeImported     = "IMPORTED"
	InferredLinkConfidence = 0.4
	inferRulePath          = "R-PATH"
	inferRuleTitle         = "R-TITLE"
	inferRulePlan          = "R-PLAN"
)

// InferOptions controls InferScopes.
type InferOptions struct {
	DryRun bool
}

// InferCandidate is a proposed scope_member edge (dry-run or inserted).
type InferCandidate struct {
	FromType   string  `json:"from_type"`
	FromID     string  `json:"from_id"`
	ToID       string  `json:"to_id"` // scope id
	Rel        string  `json:"rel"`
	Rule       string  `json:"rule"`
	Reason     string  `json:"reason,omitempty"`
	Confidence float64 `json:"confidence"`
}

// InferSkip records why an entity did not receive an INFERRED membership.
type InferSkip struct {
	FromType string `json:"from_type"`
	FromID   string `json:"from_id"`
	Reason   string `json:"reason"` // conflict | explicit | no_candidate
}

// InferReport summarizes an inference pass.
type InferReport struct {
	Inserted        int              `json:"inserted"`
	SkippedExisting int              `json:"skipped_existing"`
	SkippedConflict int              `json:"skipped_conflict"`
	SkippedExplicit int              `json:"skipped_explicit"`
	Candidates      []InferCandidate `json:"candidates,omitempty"`
	Skips           []InferSkip      `json:"skips,omitempty"`
}

type inferCand struct {
	scopeID string
	rule    string
	reason  string
}

type entityKey struct {
	typ string
	id  string
}

// InferScopes runs MVP scope-membership rules (R-PATH, R-TITLE, R-PLAN).
// Writes only via InsertLinkOrIgnore with source_type=INFERRED and confidence 0.4.
// Never deletes or updates existing links. Fail-closed on multi-scope conflicts;
// any USER_ASSERTED/IMPORTED scope_member blocks further INFERRED membership for that entity.
func (s *Service) InferScopes(ctx context.Context, opts InferOptions) (InferReport, error) {
	_ = ctx
	var report InferReport

	scopes, err := s.store.ListScopes()
	if err != nil {
		return report, err
	}
	if len(scopes) == 0 {
		return report, nil
	}

	existing, err := s.store.ListLinksByRel(RelScopeMember)
	if err != nil {
		return report, err
	}
	explicitMember := map[string]bool{} // entity id → has USER_ASSERTED/IMPORTED membership
	for _, l := range existing {
		if l.ToType != EntityScope {
			continue
		}
		if l.SourceType == DefaultSourceType || l.SourceType == SourceTypeImported {
			explicitMember[l.FromID] = true
		}
	}

	bag := map[entityKey][]inferCand{}

	if err := s.collectPathCandidates(bag, scopes); err != nil {
		return report, err
	}
	if err := s.collectTitleCandidates(bag, scopes); err != nil {
		return report, err
	}
	if err := s.collectPlanCandidates(bag, scopes); err != nil {
		return report, err
	}

	keys := make([]entityKey, 0, len(bag))
	for k := range bag {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].typ != keys[j].typ {
			return keys[i].typ < keys[j].typ
		}
		return keys[i].id < keys[j].id
	})

	for _, k := range keys {
		cands := bag[k]
		if explicitMember[k.id] {
			report.SkippedExplicit++
			report.Skips = append(report.Skips, InferSkip{
				FromType: k.typ, FromID: k.id, Reason: "explicit",
			})
			continue
		}
		scopeIDs := uniqueScopeIDs(cands)
		if len(scopeIDs) == 0 {
			continue
		}
		if len(scopeIDs) > 1 {
			report.SkippedConflict++
			report.Skips = append(report.Skips, InferSkip{
				FromType: k.typ, FromID: k.id, Reason: "conflict",
			})
			continue
		}
		scopeID := scopeIDs[0]
		rule, reason := pickRuleReason(cands, scopeID)
		cand := InferCandidate{
			FromType:   k.typ,
			FromID:     k.id,
			ToID:       scopeID,
			Rel:        RelScopeMember,
			Rule:       rule,
			Reason:     reason,
			Confidence: InferredLinkConfidence,
		}
		report.Candidates = append(report.Candidates, cand)
		if opts.DryRun {
			continue
		}
		inserted, _, err := s.store.InsertLinkOrIgnore(store.EntityLink{
			FromType:   k.typ,
			FromID:     k.id,
			Rel:        RelScopeMember,
			ToType:     EntityScope,
			ToID:       scopeID,
			SourceType: SourceTypeInferred,
			Confidence: InferredLinkConfidence,
		})
		if err != nil {
			return report, err
		}
		if inserted {
			report.Inserted++
		} else {
			report.SkippedExisting++
		}
	}
	return report, nil
}

func uniqueScopeIDs(cands []inferCand) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, c := range cands {
		if c.scopeID == "" {
			continue
		}
		if _, ok := seen[c.scopeID]; ok {
			continue
		}
		seen[c.scopeID] = struct{}{}
		out = append(out, c.scopeID)
	}
	sort.Strings(out)
	return out
}

func pickRuleReason(cands []inferCand, scopeID string) (rule, reason string) {
	for _, c := range cands {
		if c.scopeID == scopeID {
			return c.rule, c.reason
		}
	}
	return "", ""
}

func addCand(bag map[entityKey][]inferCand, typ, id, scopeID, rule, reason string) {
	if typ == "" || id == "" || scopeID == "" {
		return
	}
	k := entityKey{typ: typ, id: id}
	bag[k] = append(bag[k], inferCand{scopeID: scopeID, rule: rule, reason: reason})
}

func (s *Service) collectPathCandidates(bag map[entityKey][]inferCand, scopes []store.Scope) error {
	tasks, err := s.store.ListTasks()
	if err != nil {
		return err
	}
	for _, t := range tasks {
		paths, err := s.taskChangePaths(t.ID)
		if err != nil {
			return err
		}
		paths = append(paths, pathHintsFromBody(t.Body)...)
		for _, p := range paths {
			sid, reason := resolvePathToScope(p, scopes)
			if sid == "" {
				continue
			}
			addCand(bag, EntityTask, t.ID, sid, inferRulePath, reason)
		}
	}
	return nil
}

func pathHintsFromBody(body string) []string {
	var out []string
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		line = strings.Trim(line, "`\"'")
		n := store.NormalizePath(line)
		if n == "" {
			continue
		}
		if strings.HasPrefix(n, "web/") || strings.HasPrefix(n, "internal/") || strings.HasPrefix(n, "cmd/") {
			out = append(out, n)
		}
	}
	return out
}

// resolvePathToScope maps a repo-relative path to an existing graph scope.
// Signal: change paths (and body path hints). Prefer feature refine (auth/billing +
// matching slug) over layer (frontend/backend) when both exist and unique.
func resolvePathToScope(path string, scopes []store.Scope) (scopeID, reason string) {
	path = store.NormalizePath(strings.TrimSpace(path))
	if path == "" {
		return "", ""
	}
	lower := strings.ToLower(path)
	layer := ""
	switch {
	case strings.HasPrefix(lower, "web/"), strings.HasPrefix(lower, "web/src/"):
		layer = "frontend"
	case strings.HasPrefix(lower, "internal/"), strings.HasPrefix(lower, "cmd/"):
		layer = "backend"
	default:
		return "", ""
	}

	feature := ""
	switch {
	case pathHasToken(lower, "auth"), pathHasToken(lower, "login"), pathHasToken(lower, "session"):
		feature = "auth"
	case pathHasToken(lower, "billing"), pathHasToken(lower, "payment"):
		feature = "billing"
	}

	if feature != "" {
		if sid := findFeatureScope(scopes, feature, layer); sid != "" {
			return sid, "path feature " + feature + " (" + layer + ")"
		}
	}
	if sid := findLayerScope(scopes, layer); sid != "" {
		return sid, "path layer " + layer
	}
	return "", ""
}

func pathHasToken(path, token string) bool {
	return strings.Contains(path, "/"+token+"/") ||
		strings.Contains(path, "/"+token+".") ||
		strings.HasSuffix(path, "/"+token) ||
		strings.HasPrefix(path, token+"/")
}

func findFeatureScope(scopes []store.Scope, feature, layer string) string {
	var preferred, soft []string
	for _, sc := range scopes {
		slug := strings.ToLower(sc.Slug)
		title := strings.ToLower(sc.Title)
		if !strings.Contains(slug, feature) && !strings.Contains(title, feature) {
			continue
		}
		soft = append(soft, sc.ID)
		switch layer {
		case "frontend":
			if strings.HasSuffix(slug, "-fe") || strings.Contains(slug, "front") ||
				strings.Contains(title, "front") || strings.Contains(title, "frontend") {
				preferred = append(preferred, sc.ID)
			}
		case "backend":
			if strings.HasSuffix(slug, "-be") || strings.Contains(slug, "back") ||
				strings.Contains(title, "back") || strings.Contains(title, "backend") {
				preferred = append(preferred, sc.ID)
			}
		}
	}
	preferred = uniqueStrings(preferred)
	if len(preferred) == 1 {
		return preferred[0]
	}
	soft = uniqueStrings(soft)
	if len(soft) == 1 {
		return soft[0]
	}
	return ""
}

func findLayerScope(scopes []store.Scope, layer string) string {
	var matches []string
	for _, sc := range scopes {
		slug := strings.ToLower(sc.Slug)
		title := strings.ToLower(sc.Title)
		switch layer {
		case "frontend":
			if slug == "frontend" || slug == "fe" || strings.Contains(slug, "frontend") ||
				strings.Contains(title, "frontend") || title == "fe" {
				matches = append(matches, sc.ID)
			}
		case "backend":
			if slug == "backend" || slug == "be" || strings.Contains(slug, "backend") ||
				strings.Contains(title, "backend") || title == "be" {
				matches = append(matches, sc.ID)
			}
		}
	}
	matches = uniqueStrings(matches)
	if len(matches) == 1 {
		return matches[0]
	}
	// Prefer kind=layer when multiple.
	var layered []string
	for _, sc := range scopes {
		for _, id := range matches {
			if sc.ID == id && sc.Kind == store.ScopeKindLayer {
				layered = append(layered, id)
			}
		}
	}
	layered = uniqueStrings(layered)
	if len(layered) == 1 {
		return layered[0]
	}
	return ""
}

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func (s *Service) collectTitleCandidates(bag map[entityKey][]inferCand, scopes []store.Scope) error {
	tasks, err := s.store.ListTasks()
	if err != nil {
		return err
	}
	for _, t := range tasks {
		sid, reason := resolveTitleToScope(t.Title, scopes)
		if sid != "" {
			addCand(bag, EntityTask, t.ID, sid, inferRuleTitle, reason)
		}
	}
	decisions, err := s.store.ListDecisions()
	if err != nil {
		return err
	}
	for _, d := range decisions {
		sid, reason := resolveTitleToScope(d.Title, scopes)
		if sid != "" {
			addCand(bag, EntityDecision, d.ID, sid, inferRuleTitle, reason)
		}
	}
	discoveries, err := s.store.ListDiscoveries()
	if err != nil {
		return err
	}
	for _, d := range discoveries {
		sid, reason := resolveTitleToScope(d.Title, scopes)
		if sid != "" {
			addCand(bag, EntityDiscovery, d.ID, sid, inferRuleTitle, reason)
		}
	}
	return nil
}

func resolveTitleToScope(title string, scopes []store.Scope) (scopeID, reason string) {
	lower := strings.ToLower(title)
	tokens := tokenize(lower)
	type hit struct {
		id     string
		score  int
		reason string
	}
	var hits []hit

	authTokens := map[string]bool{"auth": true, "login": true, "session": true}
	billingTokens := map[string]bool{"billing": true, "payment": true}
	feTokens := map[string]bool{"gui": true, "overview": true, "frontend": true, "fe": true}

	hasAuth, hasBilling, hasFE := false, false, false
	for _, tok := range tokens {
		if authTokens[tok] {
			hasAuth = true
		}
		if billingTokens[tok] {
			hasBilling = true
		}
		if feTokens[tok] {
			hasFE = true
		}
	}

	if hasAuth {
		for _, sc := range scopes {
			slug := strings.ToLower(sc.Slug)
			ttl := strings.ToLower(sc.Title)
			if strings.Contains(slug, "auth") || strings.Contains(ttl, "auth") {
				hits = append(hits, hit{id: sc.ID, score: 2, reason: "title auth"})
			}
		}
	}
	if hasBilling {
		for _, sc := range scopes {
			slug := strings.ToLower(sc.Slug)
			ttl := strings.ToLower(sc.Title)
			if strings.Contains(slug, "billing") || strings.Contains(ttl, "billing") ||
				strings.Contains(slug, "payment") || strings.Contains(ttl, "payment") {
				hits = append(hits, hit{id: sc.ID, score: 2, reason: "title billing"})
			}
		}
	}
	if hasFE {
		if sid := findLayerScope(scopes, "frontend"); sid != "" {
			hits = append(hits, hit{id: sid, score: 1, reason: "title frontend"})
		}
		for _, sc := range scopes {
			slug := strings.ToLower(sc.Slug)
			ttl := strings.ToLower(sc.Title)
			if strings.Contains(slug, "front") || strings.Contains(ttl, "front") ||
				strings.Contains(slug, "gui") || strings.Contains(ttl, "gui") {
				hits = append(hits, hit{id: sc.ID, score: 1, reason: "title frontend feature"})
			}
		}
	}

	if len(hits) == 0 {
		return "", ""
	}
	best := 0
	for _, h := range hits {
		if h.score > best {
			best = h.score
		}
	}
	var bestIDs []string
	var bestReason string
	seen := map[string]struct{}{}
	for _, h := range hits {
		if h.score != best {
			continue
		}
		if _, ok := seen[h.id]; ok {
			continue
		}
		seen[h.id] = struct{}{}
		bestIDs = append(bestIDs, h.id)
		bestReason = h.reason
	}
	if len(bestIDs) != 1 {
		return "", "" // fail-closed within title rule
	}
	return bestIDs[0], bestReason
}

func tokenize(s string) []string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	parts := strings.Fields(b.String())
	return parts
}

func (s *Service) collectPlanCandidates(bag map[entityKey][]inferCand, scopes []store.Scope) error {
	goals, err := s.store.ListGoals()
	if err != nil {
		return err
	}
	for _, g := range goals {
		planScopes, err := s.store.ListPlanScopesByGoal(g.ID)
		if err != nil {
			return err
		}
		tasks, err := s.store.ListTasksByGoalID(g.ID)
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			continue
		}
		for _, ps := range planScopes {
			sid, reason := alignPlanScopeToGraph(ps, scopes)
			if sid == "" {
				continue
			}
			for _, t := range tasks {
				addCand(bag, EntityTask, t.ID, sid, inferRulePlan, reason)
			}
		}
	}
	return nil
}

func alignPlanScopeToGraph(ps store.PlanScope, scopes []store.Scope) (scopeID, reason string) {
	text := strings.ToLower(strings.TrimSpace(ps.Title))
	if body := strings.TrimSpace(ps.Body); body != "" {
		first := body
		if i := strings.IndexByte(body, '\n'); i >= 0 {
			first = body[:i]
		}
		text = text + " " + strings.ToLower(strings.TrimSpace(first))
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ""
	}
	tokens := tokenize(text)

	type scored struct {
		id    string
		score int
	}
	var scores []scored
	for _, sc := range scopes {
		slug := strings.ToLower(sc.Slug)
		title := strings.ToLower(sc.Title)
		score := 0
		if slug != "" && strings.Contains(text, slug) {
			score += 3
		}
		if title != "" && strings.Contains(text, title) {
			score += 3
		}
		for _, tok := range tokens {
			if len(tok) < 3 {
				continue
			}
			if strings.Contains(slug, tok) || strings.Contains(title, tok) {
				score++
			}
		}
		if score > 0 {
			scores = append(scores, scored{id: sc.ID, score: score})
		}
	}
	if len(scores) == 0 {
		return "", ""
	}
	best := 0
	for _, s := range scores {
		if s.score > best {
			best = s.score
		}
	}
	var winners []string
	for _, s := range scores {
		if s.score == best {
			winners = append(winners, s.id)
		}
	}
	winners = uniqueStrings(winners)
	if len(winners) != 1 {
		return "", "" // fail-closed: zero or ≥2 equally good
	}
	return winners[0], "plan_scope soft-align"
}
