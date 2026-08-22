package compiler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

// ContextOptions controls TaskContext / ExpandContext budgets.
type ContextOptions struct {
	TokenBudget     int
	MaxItems        int
	MaxLayer        int // progressive layer ceiling 1..3; default 1 (L0–L1)
	IncludeMarkdown bool
	IncludeWhy      bool
	Query           string
}

func (o ContextOptions) withDefaults() ContextOptions {
	if o.TokenBudget <= 0 {
		o.TokenBudget = DefaultTokenBudget
	}
	if o.MaxItems <= 0 {
		o.MaxItems = DefaultMaxItems
	}
	if o.MaxItems > DefaultMaxItems {
		o.MaxItems = DefaultMaxItems
	}
	if o.MaxLayer <= 0 {
		o.MaxLayer = 1
	}
	if o.MaxLayer > 3 {
		o.MaxLayer = 3
	}
	return o
}

// Compiler builds Layer 0–3 packets from store (+ optional retrieval engine).
type Compiler struct {
	store *store.Store
	retr  Retriever
}

// Retriever is the retrieval surface used by Compiler (Expand/Search/Why).
// *retrieval.Engine implements it; tests may stub Why for fail-closed checks.
type Retriever interface {
	Expand(ctx context.Context, seeds []retrieval.Hit, depth int) ([]retrieval.Hit, error)
	Search(ctx context.Context, query string, opts retrieval.SearchOptions) ([]retrieval.Hit, error)
	SearchGraphLabels(ctx context.Context, intent retrieval.Intent, opts retrieval.SearchOptions) ([]retrieval.Hit, error)
	Why(ctx context.Context, entityType, entityID string) (retrieval.WhyResult, error)
}

// New constructs a Compiler. st must be non-nil.
func New(st *store.Store) *Compiler {
	if st == nil {
		panic("compiler: New: store is nil")
	}
	return &Compiler{store: st, retr: retrieval.New(st)}
}

// WithRetrieval replaces/sets the retrieval engine (may share store).
func (c *Compiler) WithRetrieval(e Retriever) *Compiler {
	if e != nil {
		c.retr = e
	}
	return c
}

// TaskContext builds a Layer 0–1 packet (default max_layer=1) with expand depth 1.
func (c *Compiler) TaskContext(ctx context.Context, taskID string, opts ContextOptions) (Packet, error) {
	return c.compileAtDepth(ctx, taskID, 1, opts)
}

// ExpandContext builds a packet with explicit expand depth in 1..2.
func (c *Compiler) ExpandContext(ctx context.Context, taskID string, depth int, opts ContextOptions) (Packet, error) {
	if depth < 1 || depth > 2 {
		return Packet{}, fmt.Errorf("compiler: ExpandContext: depth must be 1..2, got %d", depth)
	}
	return c.compileAtDepth(ctx, taskID, depth, opts)
}

func (c *Compiler) compileAtDepth(ctx context.Context, taskID string, depth int, opts ContextOptions) (Packet, error) {
	opts = opts.withDefaults()
	if taskID == "" {
		return Packet{}, fmt.Errorf("compiler: taskID required")
	}

	task, err := c.store.GetTask(taskID)
	if err != nil {
		return Packet{}, err
	}

	var items []Item
	// Layer 0: task core
	d0 := 0
	items = append(items, Item{
		EntityType: "task",
		EntityID:   task.ID,
		Title:      task.Title,
		Excerpt:    excerptBody(task.Body),
		ReasonCode: retrieval.ReasonDirectTaskScope,
		Distance:   &d0,
		Trust:      TrustUntrustedData,
		Layer:      0,
		Provenance: &Provenance{Status: task.Status, SourceType: task.SourceType, Confidence: task.Confidence},
	})
	// work_state as system-labeled companion via excerpt note on a system item
	items = append(items, Item{
		EntityType: "task_state",
		EntityID:   task.ID,
		Title:      "work_state",
		Excerpt:    task.WorkState,
		ReasonCode: retrieval.ReasonDirectTaskScope,
		Distance:   &d0,
		Trust:      TrustSystem,
		Layer:      0,
	})

	if task.GoalID != nil && *task.GoalID != "" {
		g, err := c.store.GetGoal(*task.GoalID)
		if err == nil {
			items = append(items, Item{
				EntityType: "goal",
				EntityID:   g.ID,
				Title:      g.Title,
				Excerpt:    excerptBody(g.Body),
				ReasonCode: retrieval.ReasonGoalHasTask,
				Distance:   intPtr(0),
				Trust:      TrustUntrustedData,
				Layer:      0,
				Provenance: &Provenance{Status: g.Status, SourceType: g.SourceType, Confidence: g.Confidence},
			})
		}
	}

	seed := retrieval.Hit{
		EntityType: "task",
		EntityID:   task.ID,
		Title:      task.Title,
		Excerpt:    excerptBody(task.Body),
		ReasonCode: retrieval.ReasonDirectTaskScope,
		Score:      1.0,
		Distance:   0,
	}

	expanded, err := c.retr.Expand(ctx, []retrieval.Hit{seed}, depth)
	if err != nil {
		return Packet{}, err
	}

	intentIn := retrieval.IntentInput{
		TaskTitle: task.Title,
		TaskBody:  task.Body,
		Query:     strings.TrimSpace(opts.Query),
	}
	taskIntent := retrieval.ExtractIntent(intentIn)
	searchOpts := retrieval.SearchOptions{Limit: 16, Intent: &intentIn}

	// FTS on task title tokens (bounded) for assumptions etc.
	// Search errors are Expand-only: do not abort the packet (DF-87).
	fts, err := c.retr.Search(ctx, task.Title, searchOpts)
	if err != nil {
		fts = nil
	}

	candidates := append([]retrieval.Hit{}, expanded...)
	candidates = append(candidates, fts...)

	// Optional agent query FTS (G1): merge after title FTS, before file-seed expand.
	// Fail-open like title FTS (DF-87): query search error → skip query hits.
	if q := strings.TrimSpace(opts.Query); q != "" {
		qfts, err := c.retr.Search(ctx, q, searchOpts)
		if err == nil {
			candidates = append(candidates, qfts...)
		}
	}

	labels, err := c.retr.SearchGraphLabels(ctx, taskIntent, searchOpts)
	if err != nil {
		labels = nil // DF-87 fail-open
	}
	candidates = retrieval.MergeConceptHits(candidates, labels)

	// DF-65: Expand file-typed hits from Expand∪FTS so import hops / edge_provenance
	// reach the context packet (reuse S01 resolve inside Retriever.Expand).
	var fileSeeds []retrieval.Hit
	fileSeen := map[string]struct{}{}
	for _, h := range candidates {
		if h.EntityType != "file" || h.EntityID == "" {
			continue
		}
		if _, ok := fileSeen[h.EntityID]; ok {
			continue
		}
		fileSeen[h.EntityID] = struct{}{}
		fileSeeds = append(fileSeeds, h)
	}
	if len(fileSeeds) > 0 {
		fileExp, err := c.retr.Expand(ctx, fileSeeds, 1)
		if err != nil {
			return Packet{}, err
		}
		candidates = append(candidates, fileExp...)
	}

	maxLayer := opts.MaxLayer
	if extra, err := c.enrichLayerCandidates(ctx, candidates, maxLayer); err != nil {
		return Packet{}, err
	} else if len(extra) > 0 {
		candidates = append(candidates, extra...)
	}
	sortHits(candidates)

	seen := map[string]struct{}{
		"task\x00" + task.ID:       {},
		"task_state\x00" + task.ID: {},
	}
	if task.GoalID != nil {
		seen["goal\x00"+*task.GoalID] = struct{}{}
	}

	// DF-63: items_total = layer0 + unique admissible hits in the full candidate
	// list (not truncated by MaxCandidateHits).
	layer0Len := len(items)
	admitUniverse := 0
	admitSeen := map[string]struct{}{}
	for _, h := range candidates {
		k := layerAdmitKey(h, task, maxLayer)
		if k == "" {
			continue
		}
		if _, ok := admitSeen[k]; ok {
			continue
		}
		admitSeen[k] = struct{}{}
		admitUniverse++
	}
	itemsTotal := layer0Len + admitUniverse

	var layerExtra []Item
	candidatesCapped := false
	for i, h := range candidates {
		if h.EntityType == "task" && h.EntityID == task.ID {
			continue
		}
		if task.GoalID != nil && h.EntityType == "goal" && h.EntityID == *task.GoalID {
			continue
		}
		layer, ok := admitLayer(h, task, maxLayer)
		if !ok || layer == 0 {
			continue
		}
		k := h.EntityType + "\x00" + h.EntityID
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}

		trust := TrustUntrustedData
		dist := h.Distance
		sc := h.Score
		layerExtra = append(layerExtra, Item{
			EntityType:     h.EntityType,
			EntityID:       h.EntityID,
			Title:          h.Title,
			Excerpt:        h.Excerpt,
			ReasonCode:     h.ReasonCode,
			Distance:       &dist,
			Score:          &sc,
			Trust:          trust,
			EdgeProvenance: h.EdgeProvenance,
			Layer:          layer,
		})
		if len(layerExtra)+len(items) >= MaxCandidateHits {
			for _, rest := range candidates[i+1:] {
				if layerAdmitKey(rest, task, maxLayer) == "" {
					continue
				}
				rk := rest.EntityType + "\x00" + rest.EntityID
				if _, ok := seen[rk]; ok {
					continue
				}
				candidatesCapped = true
				break
			}
			break
		}
	}
	items = append(items, layerExtra...)

	sortItemsForBudget(items)
	kept, tokensEst, truncated := trimToBudget(items, opts.TokenBudget, opts.MaxItems)
	itemsKept := len(kept)
	if itemsKept < itemsTotal || candidatesCapped {
		truncated = true
	}
	maxLayerIncluded := 0
	for _, it := range kept {
		if it.Layer > maxLayerIncluded {
			maxLayerIncluded = it.Layer
		}
	}

	pkt := Packet{
		SchemaVersion: SchemaVersion,
		Layer:         maxLayerIncluded,
		TaskID:        taskID,
		GeneratedAt:   time.Now().UTC(),
		Budget: Budget{
			TokenLimit:       opts.TokenBudget,
			TokensEst:        tokensEst,
			MaxItems:         opts.MaxItems,
			ItemsTotal:       itemsTotal,
			ItemsKept:        itemsKept,
			CandidatesCapped: candidatesCapped,
			Truncated:        truncated,
		},
		// DF-62: honesty over pre-trim file items (not kept alone).
		IndexHonesty:     buildIndexHonesty(c.store, items),
		GraphSyncHonesty: buildGraphSyncHonesty(c.store),
		Items:            kept,
	}
	if kw := taskIntent.SummaryKeywords(); kw != "" {
		pkt.IntentSummary = &IntentSummary{
			Keywords: kw,
			Scope:    taskIntent.Scope,
			Source:   taskIntent.Source,
		}
	}

	if opts.IncludeWhy {
		why, err := c.retr.Why(ctx, "task", taskID)
		if err != nil {
			return Packet{}, err
		}
		for _, s := range why.Steps {
			pkt.WhyTrace = append(pkt.WhyTrace, WhyTraceStep{
				EntityType:     s.EntityType,
				EntityID:       s.EntityID,
				ReasonCode:     s.ReasonCode,
				Title:          s.Title,
				EdgeProvenance: s.EdgeProvenance,
			})
		}
	}

	attachTaskCapabilities(&pkt, c.store, taskID)
	attachTaskImpact(ctx, &pkt, c.store, taskID)
	attachEvidenceSections(&pkt, c.store, taskID)

	if opts.IncludeMarkdown {
		pkt.SetMarkdownCache(RenderMarkdown(pkt))
	}
	return pkt, nil
}

// attachTaskImpact sets packet.impact from domain ImpactReport (after item trim).
func attachTaskImpact(ctx context.Context, pkt *Packet, st *store.Store, taskID string) {
	summaries, err := domain.New(st).ImpactSummariesForTask(ctx, taskID)
	if err != nil || len(summaries) == 0 {
		return
	}
	pkt.Impact = summaries
}

// attachTaskCapabilities sets required_capabilities + missing_capabilities only
// (never dumps the full catalog). Missing = required with status != AVAILABLE
// or absent capability row.
func attachTaskCapabilities(pkt *Packet, st *store.Store, taskID string) {
	required, err := st.ListCapabilitiesRequiredByTaskID(taskID)
	if err != nil || len(required) == 0 {
		return
	}
	pkt.RequiredCapabilities = make([]CapabilityRef, 0, len(required))
	for _, cap := range required {
		pkt.RequiredCapabilities = append(pkt.RequiredCapabilities, capabilityRefFromStore(cap))
		if cap.Status != store.CapabilityStatusAvailable {
			pkt.MissingCapabilities = append(pkt.MissingCapabilities, capabilityRefFromStore(cap))
		}
	}
}

func capabilityRefFromStore(c store.Capability) CapabilityRef {
	return CapabilityRef{
		ID:     c.ID,
		Kind:   c.Kind,
		Slug:   c.Slug,
		Title:  c.Title,
		Status: c.Status,
	}
}

func excerptBody(body string) string {
	if len(body) <= 240 {
		return body
	}
	return body[:240] + "…"
}

func (c *Compiler) taskIntentInput(taskID, query string) retrieval.IntentInput {
	in := retrieval.IntentInput{Query: strings.TrimSpace(query)}
	if c == nil || c.store == nil || strings.TrimSpace(taskID) == "" {
		return in
	}
	task, err := c.store.GetTask(taskID)
	if err != nil {
		return in
	}
	in.TaskTitle = task.Title
	in.TaskBody = task.Body
	return in
}

func intPtr(v int) *int { return &v }

func eligibleEntityType(entityType string) bool {
	switch entityType {
	case "decision", "assumption", "discovery", "plan_change", "claim", "evidence", "file", "symbol", "goal", "task", "commit":
		return true
	default:
		return false
	}
}

// classifyNaturalLayer assigns the deepest natural layer for a hit (before max_layer cap).
func classifyNaturalLayer(h retrieval.Hit, task store.Task) int {
	if h.ReasonCode == retrieval.ReasonHistoricalVCS {
		return 3
	}
	if h.Distance >= 2 && (h.EntityType == "decision" || h.EntityType == "evidence") {
		return 3
	}
	if h.Distance >= 3 && h.ReasonCode == retrieval.ReasonGraphNeighbor {
		return 3
	}
	if h.ReasonCode == retrieval.ReasonGraphNeighbor && h.Distance >= 2 {
		return 2
	}
	if h.ReasonCode == retrieval.ReasonDiscoveryCausesPlanChg || h.ReasonCode == retrieval.ReasonRecentEvent {
		return 2
	}
	if h.EntityType == "discovery" {
		return 2
	}
	if h.EntityType == "task" && h.EntityID != task.ID {
		return 2
	}
	return 1
}

// admitLayer returns the layer for a candidate hit, capped by maxLayer.
// When natural layer exceeds maxLayer, eligible hits downgrade to L1 (never auto-promote).
func admitLayer(h retrieval.Hit, task store.Task, maxLayer int) (layer int, ok bool) {
	if h.EntityType == "task" && h.EntityID == task.ID {
		return 0, false
	}
	if task.GoalID != nil && h.EntityType == "goal" && h.EntityID == *task.GoalID {
		return 0, false
	}
	if !eligibleEntityType(h.EntityType) {
		return 0, false
	}
	natural := classifyNaturalLayer(h, task)
	if natural > maxLayer {
		if maxLayer >= 1 && natural >= 1 && eligibleEntityType(h.EntityType) {
			return 1, true
		}
		return 0, false
	}
	return natural, true
}

// layerAdmitKey returns a non-empty key if the hit is eligible for admission at maxLayer.
func layerAdmitKey(h retrieval.Hit, task store.Task, maxLayer int) string {
	if _, ok := admitLayer(h, task, maxLayer); !ok {
		return ""
	}
	return h.EntityType + "\x00" + h.EntityID
}

type impactWalker interface {
	ImpactWalk(ctx context.Context, seeds []retrieval.ImpactSeed, depth int) (*retrieval.ImpactWalkResult, error)
}

type historicalEnricher interface {
	HistoricalFileHits(ctx context.Context, paths []string) ([]retrieval.Hit, error)
}

func (c *Compiler) enrichLayerCandidates(ctx context.Context, candidates []retrieval.Hit, maxLayer int) ([]retrieval.Hit, error) {
	if maxLayer < 2 || c.retr == nil {
		return nil, nil
	}
	var out []retrieval.Hit

	if w, ok := c.retr.(impactWalker); ok {
		seeds := impactSeedsFromCandidates(candidates)
		if len(seeds) > 0 {
			if maxLayer >= 2 {
				res, err := w.ImpactWalk(ctx, seeds, 1)
				if err != nil {
					return nil, err
				}
				out = append(out, retrieval.BlastHitsToLayer(res, 1)...)
			}
			if maxLayer >= 3 {
				res, err := w.ImpactWalk(ctx, seeds, 2)
				if err != nil {
					return nil, err
				}
				out = append(out, retrieval.BlastHitsToLayer(res, 2)...)
			}
		}
	}

	if maxLayer >= 3 {
		if he, ok := c.retr.(historicalEnricher); ok {
			hist, err := he.HistoricalFileHits(ctx, filePathsFromHits(candidates))
			if err != nil {
				return nil, err
			}
			out = append(out, hist...)
		}
	}
	return out, nil
}

func impactSeedsFromCandidates(candidates []retrieval.Hit) []retrieval.ImpactSeed {
	seen := map[string]struct{}{}
	var seeds []retrieval.ImpactSeed
	for _, h := range candidates {
		if h.EntityType != "file" && h.EntityType != "symbol" {
			continue
		}
		k := h.EntityType + "\x00" + h.EntityID
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		seeds = append(seeds, retrieval.ImpactSeed{EntityType: h.EntityType, EntityID: h.EntityID})
	}
	return seeds
}

func filePathsFromHits(candidates []retrieval.Hit) []string {
	seen := map[string]struct{}{}
	var paths []string
	for _, h := range candidates {
		p := h.Path
		if p == "" && h.EntityType == "file" {
			p = h.Title
		}
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		paths = append(paths, p)
	}
	return paths
}
