package retrieval

import (
	"context"
	"fmt"
	"sort"

	"github.com/mrchatam/Trace/internal/store"
)

const (
	// MaxImpactBlast is the loud blast list cap (aligned with MaxCandidateHits).
	MaxImpactBlast     = 64
	maxImpactDepth     = 2
	defaultImpactDepth = 2
)

// ImpactSeed is a structural walk seed (file|symbol only).
type ImpactSeed struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// BlastHit is one non-seed node in an impact walk result.
type BlastHit struct {
	EntityType     string  `json:"entity_type"`
	EntityID       string  `json:"entity_id"`
	Title          string  `json:"title,omitempty"`
	Path           string  `json:"path,omitempty"`
	Hop            int     `json:"hop"`
	HopRisk        float64 `json:"hop_risk"`
	EdgeProvenance string  `json:"edge_provenance,omitempty"`
}

// ImpactWalkResult is the loud multi-seed BFS blast (seeds excluded).
type ImpactWalkResult struct {
	Seeds         []ImpactSeed `json:"seeds"`
	Blast         []BlastHit   `json:"blast"`
	AffectedTests []BlastHit   `json:"affected_tests"`
	BlastTotal    int          `json:"blast_total"`
	BlastKept     int          `json:"blast_kept"`
	Truncated     bool         `json:"truncated"`
	Depth         int          `json:"depth"`
}

type impactFrontier struct {
	entityType       string
	entityID         string
	title            string
	path             string
	hop              int
	edgeProvenance   string
	symbolKind       string
	allowContainsOut bool // file→symbols OK; false after contains-UP (no sibling climb)
}

// ImpactWalk runs one multi-seed BFS over file|symbol seeds with incoming-import
// deps and contains asymmetry. Seeds are hop 0 and excluded from blast.
// Depth must be 1..2; callers that omit depth at the CLI use defaultImpactDepth (2).
func (e *Engine) ImpactWalk(ctx context.Context, seeds []ImpactSeed, depth int) (*ImpactWalkResult, error) {
	_ = ctx
	if len(seeds) == 0 {
		return nil, fmt.Errorf("retrieval: ImpactWalk: at least one seed required")
	}
	if depth < 1 || depth > maxImpactDepth {
		return nil, fmt.Errorf("retrieval: ImpactWalk: depth must be 1..%d, got %d", maxImpactDepth, depth)
	}

	seedKeys := map[string]struct{}{}
	normalizedSeeds := make([]ImpactSeed, 0, len(seeds))
	var frontier []impactFrontier
	seen := map[string]*impactFrontier{}

	for _, s := range seeds {
		typ := NormalizeEntityType(s.EntityType)
		if typ != "file" && typ != "symbol" {
			return nil, fmt.Errorf("retrieval: ImpactWalk: seed entity_type must be file|symbol, got %q", s.EntityType)
		}
		if s.EntityID == "" {
			return nil, fmt.Errorf("retrieval: ImpactWalk: seed entity_id required")
		}
		h, err := e.lookupEntity(typ, s.EntityID, ReasonExactID, 0, 1.0)
		if err != nil {
			return nil, fmt.Errorf("retrieval: ImpactWalk: seed %s:%s: %w", typ, s.EntityID, err)
		}
		k := hitKey(typ, s.EntityID)
		seedKeys[k] = struct{}{}
		normalizedSeeds = append(normalizedSeeds, ImpactSeed{EntityType: typ, EntityID: s.EntityID})

		n := &impactFrontier{
			entityType:       typ,
			entityID:         s.EntityID,
			title:            h.Title,
			path:             h.Path,
			hop:              0,
			allowContainsOut: typ == "file", // file seeds may expand contains-OUT
		}
		if _, ok := seen[k]; !ok {
			seen[k] = n
			frontier = append(frontier, *n)
		}
	}

	for d := 1; d <= depth; d++ {
		var next []impactFrontier
		for _, fi := range frontier {
			neighbors, err := e.impactNeighbors(fi)
			if err != nil {
				return nil, err
			}
			for _, n := range neighbors {
				n.hop = d
				k := hitKey(n.entityType, n.entityID)
				if existing, ok := seen[k]; ok {
					upgraded := false
					if n.hop < existing.hop {
						existing.hop = n.hop
						existing.edgeProvenance = n.edgeProvenance
						existing.title = n.title
						existing.path = n.path
						upgraded = true
					}
					if n.allowContainsOut && !existing.allowContainsOut {
						existing.allowContainsOut = true
						upgraded = true
					}
					if n.symbolKind == "test" && existing.symbolKind != "test" {
						existing.symbolKind = "test"
						upgraded = true
					}
					if upgraded && existing.hop == d {
						next = append(next, *existing)
					}
					continue
				}
				cp := n
				seen[k] = &cp
				next = append(next, n)
			}
		}
		frontier = next
	}

	var blast []BlastHit
	for k, n := range seen {
		if _, isSeed := seedKeys[k]; isSeed {
			continue
		}
		blast = append(blast, BlastHit{
			EntityType:     n.entityType,
			EntityID:       n.entityID,
			Title:          n.title,
			Path:           n.path,
			Hop:            n.hop,
			HopRisk:        float64(n.hop),
			EdgeProvenance: n.edgeProvenance,
		})
	}
	sort.Slice(blast, func(i, j int) bool {
		if blast[i].Hop != blast[j].Hop {
			return blast[i].Hop < blast[j].Hop
		}
		if blast[i].EntityType != blast[j].EntityType {
			return blast[i].EntityType < blast[j].EntityType
		}
		return blast[i].EntityID < blast[j].EntityID
	})

	total := len(blast)
	truncated := total > MaxImpactBlast
	kept := blast
	if truncated {
		kept = blast[:MaxImpactBlast]
	}

	var affected []BlastHit
	for _, h := range kept {
		k := hitKey(h.EntityType, h.EntityID)
		if n, ok := seen[k]; ok && n.symbolKind == "test" {
			affected = append(affected, h)
		}
	}
	if affected == nil {
		affected = []BlastHit{}
	}

	return &ImpactWalkResult{
		Seeds:         normalizedSeeds,
		Blast:         kept,
		AffectedTests: affected,
		BlastTotal:    total,
		BlastKept:     len(kept),
		Truncated:     truncated,
		Depth:         depth,
	}, nil
}

// DefaultImpactDepth returns the CLI default depth (2).
func DefaultImpactDepth() int { return defaultImpactDepth }

func (e *Engine) impactNeighbors(fi impactFrontier) ([]impactFrontier, error) {
	switch fi.entityType {
	case "file":
		return e.impactNeighborsFile(fi)
	case "symbol":
		return e.impactNeighborsSymbol(fi)
	default:
		return nil, nil
	}
}

func (e *Engine) impactNeighborsFile(fi impactFrontier) ([]impactFrontier, error) {
	path := fi.path
	if path == "" {
		f, err := e.store.GetFileByID(fi.entityID)
		if err != nil {
			return nil, err
		}
		path = f.Path
	}

	var out []impactFrontier

	// Contains-OUT: file → symbols (only when allowed — not after contains-UP).
	if fi.allowContainsOut {
		syms, err := e.store.ListSymbolsByPath(path)
		if err != nil {
			return nil, err
		}
		for _, sym := range syms {
			out = append(out, impactFrontier{
				entityType:       "symbol",
				entityID:         sym.ID,
				title:            sym.Name,
				path:             path,
				allowContainsOut: false,
			})
		}
	}

	// Incoming imports only: files whose resolved import targets this path.
	importers, err := e.listIncomingImporters(path)
	if err != nil {
		return nil, err
	}
	for _, imp := range importers {
		out = append(out, impactFrontier{
			entityType:       "file",
			entityID:         imp.fileID,
			title:            imp.path,
			path:             imp.path,
			edgeProvenance:   imp.provenance,
			allowContainsOut: true,
		})
	}

	validates, err := e.listReverseValidatesForFile(fi.entityID, fi.allowContainsOut, path)
	if err != nil {
		return nil, err
	}
	out = append(out, validates...)
	return out, nil
}

func (e *Engine) impactNeighborsSymbol(fi impactFrontier) ([]impactFrontier, error) {
	path := fi.path
	if path == "" {
		_, p, err := e.store.GetSymbolByID(fi.entityID)
		if err != nil {
			if isNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		path = p
	}
	f, err := e.store.GetFileByPath(path)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []impactFrontier
	// Contains-UP into blast so incoming deps of the file can be walked;
	// allowContainsOut=false prevents sibling symbols via contains climb.
	out = append(out, impactFrontier{
		entityType:       "file",
		entityID:         f.ID,
		title:            f.Path,
		path:             f.Path,
		allowContainsOut: false,
	})

	validates, err := e.store.ListValidatesForSymbol(fi.entityID)
	if err != nil {
		return nil, err
	}
	testNeighbors, err := e.validatesEdgesToTestNeighbors(validates)
	if err != nil {
		return nil, err
	}
	out = append(out, testNeighbors...)
	return out, nil
}

type incomingImporter struct {
	fileID     string
	path       string
	provenance string
}

func (e *Engine) listReverseValidatesForFile(fileID string, allowContainsOut bool, path string) ([]impactFrontier, error) {
	edges, err := e.store.ListValidatesForFile(fileID)
	if err != nil {
		return nil, err
	}
	var out []impactFrontier
	if allowContainsOut {
		syms, err := e.store.ListSymbolsByPath(path)
		if err != nil {
			return nil, err
		}
		symbolTargets := map[string]struct{}{}
		for _, sym := range syms {
			symbolTargets[sym.ID] = struct{}{}
		}
		var fileLevel []store.CodeEdge
		for _, edge := range edges {
			if edge.ToSymbolID != nil {
				if _, ok := symbolTargets[*edge.ToSymbolID]; ok {
					continue // covered by ListValidatesForSymbol below
				}
			}
			fileLevel = append(fileLevel, edge)
		}
		testNeighbors, err := e.validatesEdgesToTestNeighbors(fileLevel)
		if err != nil {
			return nil, err
		}
		out = append(out, testNeighbors...)
		for _, sym := range syms {
			symEdges, err := e.store.ListValidatesForSymbol(sym.ID)
			if err != nil {
				return nil, err
			}
			testNeighbors, err := e.validatesEdgesToTestNeighbors(symEdges)
			if err != nil {
				return nil, err
			}
			out = append(out, testNeighbors...)
		}
		return out, nil
	}
	// After contains-UP: only file-level validates (no contained-symbol targets).
	var fileLevel []store.CodeEdge
	for _, edge := range edges {
		if edge.ToSymbolID == nil {
			fileLevel = append(fileLevel, edge)
		}
	}
	return e.validatesEdgesToTestNeighbors(fileLevel)
}

func (e *Engine) validatesEdgesToTestNeighbors(edges []store.CodeEdge) ([]impactFrontier, error) {
	var out []impactFrontier
	for _, edge := range edges {
		if edge.FromSymbolID == nil {
			continue
		}
		sym, symPath, err := e.store.GetSymbolByID(*edge.FromSymbolID)
		if err != nil {
			if isNotFound(err) {
				continue
			}
			return nil, err
		}
		if sym.Kind != "test" {
			continue
		}
		out = append(out, impactFrontier{
			entityType:       "symbol",
			entityID:         sym.ID,
			title:            sym.Name,
			path:             symPath,
			edgeProvenance:   edge.Provenance,
			symbolKind:       "test",
			allowContainsOut: false,
		})
	}
	return out, nil
}

func (e *Engine) listIncomingImporters(targetPath string) ([]incomingImporter, error) {
	edges, err := e.store.ListImportEdges()
	if err != nil {
		return nil, err
	}
	var out []incomingImporter
	seen := map[string]struct{}{}
	for _, edge := range edges {
		resolved, err := e.resolveImportedFile(edge.ImporterPath, edge.ImportedPath)
		if err != nil {
			if isNotFound(err) {
				continue
			}
			return nil, err
		}
		if resolved.Path != targetPath {
			continue
		}
		if _, ok := seen[edge.ImporterPath]; ok {
			continue
		}
		impFile, err := e.store.GetFileByPath(edge.ImporterPath)
		if err != nil {
			if isNotFound(err) {
				continue
			}
			return nil, err
		}
		seen[edge.ImporterPath] = struct{}{}
		out = append(out, incomingImporter{
			fileID:     impFile.ID,
			path:       impFile.Path,
			provenance: edge.Provenance,
		})
	}
	return out, nil
}
