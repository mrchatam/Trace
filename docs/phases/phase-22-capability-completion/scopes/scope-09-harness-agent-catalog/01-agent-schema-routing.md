# P22-S09-01 — Implement: harness agent schema + routing library

## Metadata
- id: P22-S09-01
- todo_ids: [P22-S09-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, domain-modeling]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Add the **harness agent catalog** data model and deterministic **routing library**. Trace stores agents it understands and maps them to required environment capabilities (skills, MCPs, hooks, tools). **E03** foundation; **E02** routing core. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) — D-22-25…D-22-30
- [WORK-MAP.md](../../WORK-MAP.md) — W-37, W-38
- Live: `internal/store/schema/010_capability_surface.sql`, `internal/store/capability.go`, `internal/domain/capability.go`, `internal/store/eval_rules.go` (recent mig pattern)

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| `capabilities` table (SKILL/RULE/MCP/TOOL/HOOK only) | `AGENT` kind |
| Compat ceiling **26**; **26** SQL files | mig **027** / `harness_agents` tables |
| `TestUpsertCapabilityGetAndReject` rejects `Kind: "AGENT"` | `internal/agents/` package |
| `trace/agents/README.md` (doc stub) | `trace/agents/default.json` |
| MCP catalog **14** tools | `trace_agents`, `harness_recommendations[]` |
| Deliberation phases: CRITIQUE, VERIFY, INVESTIGATE, ORIENT, … | Agent routing library |

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **`027_harness_agents.sql`** only; compat ceiling **27**; forbid **028+** until a later scope owns it |
| AGENT kind | Extend `store.CapabilityKindAgent` + `domain.CapabilityKindAgent`; `NormalizeCapabilityKind` accepts **`AGENT`** |
| Agent slug convention | `agent:<name>` (e.g. `agent:code-reviewer`) — stored in `harness_agents.slug`, not duplicated in `capabilities` |
| Trace role | **Recommend only** — library returns ranked suggestions; **no** Task tool, subprocess, HTTP, ML |
| Seed export | Additive keys `harness_agents[]` + nested `requirements[]` on `SeedDocument` (follow `engineering_knowledge` pattern) |

## Schema (mig 027 — implement exactly)

```sql
-- Migration v27: harness agent catalog (Phase 22 S09).
-- Additive only; do not rewrite 001–026.

CREATE TABLE IF NOT EXISTS harness_agents (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,              -- agent:<name>
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    subagent_type TEXT NOT NULL DEFAULT '', -- harness profile string e.g. performance-reviewer
    deliberation_phases TEXT NOT NULL DEFAULT '[]',  -- JSON array of phase strings
    task_keywords TEXT NOT NULL DEFAULT '[]',        -- JSON array of lowercase keywords
    recommend_subagent INTEGER NOT NULL DEFAULT 0,   -- 1 when profile suits fresh-subagent review
    registry_source TEXT NOT NULL DEFAULT 'bundled', -- bundled|user|host
    registry_version TEXT NOT NULL DEFAULT '',
    external_url TEXT,                               -- nullable; no fetch in P22
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_harness_agents_slug ON harness_agents(slug);
CREATE INDEX IF NOT EXISTS idx_harness_agents_subagent_type ON harness_agents(subagent_type);

CREATE TABLE IF NOT EXISTS harness_agent_requirements (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    required_capability_slug TEXT NOT NULL,  -- skill:…, mcp:…, hook:…, tool:…
    created_at TEXT NOT NULL,
    UNIQUE(agent_id, required_capability_slug)
);
CREATE INDEX IF NOT EXISTS idx_harness_agent_requirements_agent
    ON harness_agent_requirements(agent_id);
CREATE INDEX IF NOT EXISTS idx_harness_agent_requirements_slug
    ON harness_agent_requirements(required_capability_slug);
```

Optional `harness_agent_routing` table: **defer** — phase/keyword columns on `harness_agents` suffice for P22.

## Domain + store API

1. **`internal/domain/harness_agents.go`**: `HarnessAgent`, `HarnessAgentRequirement`, input types, validation (slug prefix `agent:`, JSON arrays valid).
2. **`internal/store/harness_agents.go`**: CRUD — `UpsertHarnessAgent`, `GetHarnessAgentBySlug`, `ListHarnessAgents`, `ListHarnessAgentRequirements`, `DeleteHarnessAgentRequirementsForAgent` (for idempotent upsert).
3. Wire domain service methods on `domain.Service` mirroring eval_rules / knowledge patterns.

## Routing library (`internal/agents/routing.go`)

```go
type RecommendInput struct {
    Phase         string            // deliberation.Phase string e.g. CRITIQUE
    TaskTitle     string
    TaskTags      []string
    GoalKeywords  []string
    HarnessCaps   map[string]string // slug → status (AVAILABLE|UNAVAILABLE|UNKNOWN)
}

type Recommendation struct {
    AgentSlug            string   `json:"agent_slug"`
    SubagentType         string   `json:"subagent_type"`
    Reason               string   `json:"reason"`
    Confidence           string   `json:"confidence"` // high|medium|low
    UseSubagent          bool     `json:"use_subagent"`
    PromptStub           string   `json:"prompt_stub,omitempty"`
    MissingCapabilities  []string `json:"missing_capabilities,omitempty"`
}

func RecommendAgents(ctx context.Context, st *store.Store, in RecommendInput) ([]Recommendation, error)
```

### Deterministic routing rules (locked)

| Signal | Preferred agent(s) | Notes |
|--------|-------------------|-------|
| Phase **CRITIQUE** | `agent:code-reviewer`, then `agent:nested-reviewer` | Rank code-reviewer first |
| Phase **VERIFY** + perf keywords in title/tags/goal | `agent:performance-reviewer` | Keywords: perf, performance, latency, benchmark, slow, memory |
| Phase **CRITIQUE/VERIFY** + security keywords | `agent:security-reviewer` | auth, injection, owasp, secret, xss, csrf |
| Phase **INVESTIGATE** or **ORIENT** | `agent:explore` | Investigation / codebase exploration |
| No match | `agent:generalPurpose` | Fallback; confidence **low** |
| `recommend_subagent=true` on matched agent **and** `harness:subagent` **AVAILABLE** in HarnessCaps | set `use_subagent: true` | Prompt stub: "Fresh subagent for independent review — not the implementer session." |
| Required capability slug **not** AVAILABLE | append to `missing_capabilities[]`; do **not** drop recommendation | Honest partial match |

- Same input → same output order and scores (`TestRoutingDeterministic`).
- Cap results at **4** entries (S09-05 consumes; implement library with `maxResults` param default 4).
- Load agents from DB; empty catalog → empty slice (no error).

## Requirements

1. Create mig **027** + bump `evals/compat/compat_test.go` ceiling to **27** (mirror S07-01 → 026 pattern).
2. Add **AGENT** kind in store + domain; update `TestUpsertCapabilityGetAndReject` — AGENT must **succeed** (adjust test: unknown kinds still fail; AGENT is known).
3. Implement store/domain CRUD for harness agents + requirements.
4. Implement `RecommendAgents` with table rules above (string match on lowercased title/tags/keywords).
5. Seed export/import: additive `harness_agents` on `SeedDocument`; round-trip in `TestSeedExportRoundTrip` or dedicated keeper.
6. **Do not** add CLI, MCP, install loader, or loop packet wiring — S09-03/05/07 own those.
7. Grep guard: no `Task(`, no subprocess agent runner, no HTTP client for registry fetch.

## Touch files

- `internal/store/schema/027_harness_agents.sql` (new)
- `internal/store/harness_agents.go` (new)
- `internal/store/harness_agents_test.go` (new)
- `internal/domain/capability.go` — AGENT kind + `CapabilityKindAgent`
- `internal/domain/capability_test.go` — AGENT acceptance
- `internal/domain/harness_agents.go` (new)
- `internal/agents/routing.go`, `routing_test.go` (new)
- `internal/domain/seed_export.go` — additive export/import keys
- `evals/compat/compat_test.go` — ceiling **27**

## Named tests

| Test | Proves |
|------|--------|
| `TestHarnessAgentCatalogMigrate027` | mig 027 idempotent; 27 files |
| `TestCapabilityKindAgent` | AGENT kind normalizes + upserts |
| `TestRecommendAgentForPhaseCritique` | CRITIQUE prefers code-reviewer profile |
| `TestRecommendPerformanceReviewerForPerfTask` | perf keywords → performance-reviewer |
| `TestRoutingDeterministic` | same input → same output |
| `TestHarnessAgentSeedExportRoundTrip` | portable graph includes agents (new or extend seed keeper) |

```bash
go test ./internal/agents/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestHarnessAgent|TestCapabilityKindAgent|TestRecommend|TestRoutingDeterministic'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 27
```

## Exit criteria

- [ ] **E03** schema + store/domain CRUD exist
- [ ] **E02** routing library returns deterministic ranked recommendations
- [ ] Named tests PASS; compat **27**; exactly **27** SQL files
- [ ] No spawn/runner/HTTP in new code
- [ ] Board Notes: test summary + E03/E02 partial closure

## Minimal todos

- [ ] Mig 027 + store + domain CRUD
- [ ] AGENT kind + test update
- [ ] Routing library with phase + keyword matching
- [ ] Seed export additive keys
- [ ] Compat bump + board notes

## Residual risks (carry to S09-02)

- JSON array validation on `deliberation_phases` / `task_keywords` — invalid JSON fail-closed vs default `[]`
- Requirement slug references nonexistent capabilities — recommend anyway with `missing_capabilities`
- Seed import merge vs replace semantics for user-edited agents
- `harness:subagent` HOOK not yet declared — routing must not assume AVAILABLE until S09-03
