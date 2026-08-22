# P20-S03-00 — Change + effects planner

## Metadata
- id: P20-S03-00
- todo_ids: [P20-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Change as a first-class object with Git SHA/path refs and expected vs actual effects. **No blobs. No product Go this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [COVERAGE.md](../../COVERAGE.md) merge table (Change / Metric·Observation)
- Law 1 (Git canonical — OIDs + paths, never source blobs)
- Live: `internal/store/schema/001_init.sql` (`files.git_oid`, no file body), `002_vcs_index.sql` (`vcs_commits` + `vcs_commit_paths` — refs only), `016_cognitive_artifacts.sql` (`hypotheses`, `decision_reconsiderations.trigger` includes `contradicted_effect`)
- Live VCS: `internal/vcs.Repository.ShowFile` / `Changes` / `Head`; `internal/gitcli` implements CLI-only; domain **imports `vcs` iface only** (never `gitcli`)
- Live S02: `internal/domain/cognitive.go` (`CreateHypothesis`, `RecordDecisionReconsideration`); do **not** fork Discovery as a hypothesis stand-in
- Laws 1, 2, 5, 8, 9, 11, 16, 19

## Doc map
§10, 14, 29I, 29L

## Live inventory (2026-08-18)

| Surface | Live API | S03 action |
|---------|----------|------------|
| `files.git_oid` | Optional OID on file metadata; **no** content/body/blob column (`TestNoSourceContentColumns`) | **Do not** store change content on `files`. Change paths are Git paths, not `files.id` FKs (new file may not be indexed yet). |
| `vcs_commits` / `vcs_commit_paths` | Thin commit+path index after `Refresh`; no patch bodies | **Do not** dual-write from Change APIs. Change is semantic; VCS index stays Git-owned. |
| `vcs.Repository.ShowFile` | Bytes from `git show rev:path` | **Read path** for change content. Inject iface; tests use `vcs.Fake`. |
| Hypothesis | `CreateHypothesis` + `hypothesis_supported_by` → evidence | Contradicted effect **may link** Hypothesis. Not a second Discovery. |
| Discovery | `CreateDiscovery` severity INFO\|PLAN_AFFECTING\|BLOCKING | Contradicted effect **may** emit PLAN_AFFECTING (Law 8). Not the hypothesis. |
| Decision reconsideration | Child table; trigger `contradicted_effect` already CHECK'd | On contradicted effect, if `change_implements_decision` exists: FIRED reconsideration. |
| Regression / Baseline / test results | **missing** (S04/S05) | **Out.** No `tests` / `verification_runs` / `baseline` / `score_*` columns on `changes`. |
| Seed export | No `changes`/`effects` keys | **Out of S03** (residual S07). |
| Next migration | max **016** | S03-01 adds **`017_changes_effects.sql`** |

## Paths storage decision (FINAL)

| Option | Verdict |
|--------|---------|
| JSON array column on `changes` (`[{path, symbol_id?}]`) | **Rejected** |
| Child table `change_paths` | **Locked** |

Rejected JSON because: (1) `vcs_commit_paths` and S02 `decision_reconsiderations` are already child tables — queryable, CHECK-able, no extra `content` key can appear in a blob; (2) S05/S06 need “changes touching path X” and bounded `recent_changes[]` without parsing JSON; (3) Law 1 fail-closed is a missing column, not a convention on map keys. Implementer must **not** add `paths_json` / `files_json` / `diff` / `patch`.

## FINAL locked defaults (S03-01 must not re-debate)

| Item | Value |
|------|-------|
| Migration | **`017_changes_effects.sql`** — additive; do not rewrite 001–016; **no ALTER** on `files` / `vcs_*` / `hypotheses` / `discoveries` / `decisions` |
| Compat ceiling | **17** after S03-01 (forbid `018+`); bump `evals/compat`, `production_hardening_test`, `deliberation_test` EmbedExpected, `TestOpenCreatesDBAndMigratesIdempotent` versions. This planner row does not bump (no product Go). |
| New tables | `changes`, `change_paths`, `effects` only |
| Forbidden tables / columns | Finding, Regression, Baseline, Experiment; any `content`/`blob`/`patch`/`diff`/`file_body`/`source_text`; JSON path arrays |
| Git | Store commit OID + paths only. Diff/file bytes via `vcs.Repository.ShowFile` at read time. **Never** insert into `vcs_commits` / `vcs_commit_paths` from this API. |
| History | Law 11 — no `DeleteChange*`. SUPERSEDED is explicit. |
| Raw CoT | Forbidden |
| CLI / MCP / loop / SelectNext | **Library-only** — S06 owns apply keys `changes`/`effects` |
| FTS / seed JSON | Do **not** extend this row |
| Comparison | Agent-supplied enum; library validates. **Not** string-equality auto-score. Law 2: comparison ≠ evidence (evidence is optional links). |
| S05 | Contradiction is an **input**. Do not create Regression. Do not set `attribution=caused`. |
| S04 | Test/verify/eval/baseline stay off `changes`. |

### `changes`

```sql
CREATE TABLE IF NOT EXISTS changes (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    git_commit TEXT NOT NULL DEFAULT '',
    parent_change_id TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    reason TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'RECORDED', 'COMPARED', 'SUPERSEDED')),
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_verified_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_changes_task_id ON changes(task_id);
CREATE INDEX IF NOT EXISTS idx_changes_git_commit ON changes(git_commit);
CREATE INDEX IF NOT EXISTS idx_changes_parent ON changes(parent_change_id);
```

| Field | Lock |
|-------|------|
| `task_id` | **Required** originating task (doc §10). Fail closed if empty. No SQL FK. No duplicate `change_originates_from_task` link. |
| `git_commit` | Git OID reference only. Empty allowed on `OPEN`. Non-empty must match `^[0-9a-fA-F]{7,64}$`, store lowercase. Need not exist in `vcs_commits`. |
| `parent_change_id` | Optional; if set, `GetChange` must succeed. Creating a child does **not** auto-SUPERSEDE the parent. |
| `actor` | Opaque string (agent/human id). Not a user table. |
| `reason` | Why it changed. Max **8192** bytes (fail closed). Semantic text, not a diff. |
| `status` | Lifecycle (not provenance ACTIVE/STALE). |

**Status transitions**

| From | To | When |
|------|----|------|
| (create) | `OPEN` | `git_commit` empty |
| (create) | `RECORDED` | `git_commit` non-empty on create |
| `OPEN` | `RECORDED` | `RecordChangeCommit` with valid SHA |
| `RECORDED` | `COMPARED` | Every expected effect row for this change has non-empty `comparison` |
| `OPEN`/`RECORDED`/`COMPARED` | `SUPERSEDED` | `SupersedeChange` (reason required) |
| `SUPERSEDED` | — | Terminal |

`RecordChangeCommit` on already-`RECORDED`/`COMPARED` may update SHA only if currently empty (fail closed if replacing a different non-empty SHA — new child change instead).

Out of this table (doc §10 extras → other scopes): `tests`, `verification_runs`, `baseline`, `score_before`, `score_after`, `risks` (Uncertainty `kind=risk`).

### `change_paths` (not JSON)

```sql
CREATE TABLE IF NOT EXISTS change_paths (
    change_id TEXT NOT NULL,
    path TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    symbol_id TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (change_id, path)
);
CREATE INDEX IF NOT EXISTS idx_change_paths_path ON change_paths(path);
```

| Field | Lock |
|-------|------|
| `path` | Repo-relative; domain calls `store.NormalizePath`. Empty fail closed. **No file contents.** |
| `status` | Optional Git letter (`A`/`M`/`D`/`R`/…) or `''`. Free TEXT — do not invent a second VCS enum. |
| `symbol_id` | Optional `symbols.id` snapshot ref. Empty default. **No** symbol source text. |
| Cardinality | Create requires **≥1** path. |

### `effects` (one row per change × dimension)

```sql
CREATE TABLE IF NOT EXISTS effects (
    id TEXT PRIMARY KEY,
    change_id TEXT NOT NULL,
    dimension TEXT NOT NULL,
    expected TEXT NOT NULL DEFAULT '',
    actual TEXT NOT NULL DEFAULT '',
    comparison TEXT NOT NULL DEFAULT ''
        CHECK (comparison IN ('', 'supported', 'partially_supported', 'contradicted')),
    confidence REAL NOT NULL DEFAULT 0,
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE(change_id, dimension)
);
CREATE INDEX IF NOT EXISTS idx_effects_change_id ON effects(change_id);
CREATE INDEX IF NOT EXISTS idx_effects_comparison ON effects(comparison);
```

| Field | Lock |
|-------|------|
| `dimension` | Trimmed non-empty, max **64** bytes. Free string (`latency`, `correctness`, …). Not an enum. |
| `expected` | Required non-empty on insert. Max **8192** bytes. Semantic, not a metric dump or file body. |
| `actual` | Empty until `RecordActualEffect`. Same max. |
| `comparison` | `''` until actual recorded. Recording actual **requires** `supported` \| `partially_supported` \| `contradicted`. Unknown values **fail closed**. |

**RecordActualEffect rules:** expected row for `(change_id, dimension)` **must already exist**. Do not invent expected from actual. Comparison is caller-supplied; do not auto-derive by string match.

### Contradiction downstream (S05 consumes; S03 does not own regression)

When `comparison=contradicted`:

| Action | Lock |
|--------|------|
| Event | Always `effect.contradicted` on `entity_type=change` (payload includes `effect_id`, `dimension`) plus `effect.compared` |
| Discovery | **Optional** `EmitDiscovery` → `CreateDiscovery` severity **`PLAN_AFFECTING`** + rel `discovery_from_contradicted_effect` (discovery → effect). Title required if emitting. |
| Hypothesis | **Optional** `HypothesisID` **or** `CreateHypothesis` (not both). Link `hypothesis_explains_effect` (hypothesis → effect). New hypothesis is `OPEN` via existing `CreateHypothesis`. **Do not** insert a Discovery as the hypothesis. |
| Decision | For each `change_implements_decision` link: `RecordDecisionReconsideration` trigger=`contradicted_effect`, status=`FIRED`, `related_type=effect`, `related_id=<effect id>` |
| Uncertainty | **Do not** auto-spawn (S05 Regression → Question). |
| Regression | **Do not** create (table does not exist; S05). |
| Replan / SelectNext | **Do not** auto-hop or auto-replan (Laws 9, 16). S06 may later pass flags into PolicyInputs. |

`supported` / `partially_supported` do **not** spawn Discovery/Hypothesis/reconsideration.

### Rels (new)

| Rel | From → To | Purpose |
|-----|-----------|---------|
| `change_implements_decision` | change → decision | Optional at create; drives FIRED `contradicted_effect` reconsider |
| `effect_supported_by` | effect → evidence | Optional; Law 2 evidence for the comparison |
| `hypothesis_explains_effect` | hypothesis → effect | Contradiction → hypothesis (S05) |
| `discovery_from_contradicted_effect` | discovery → effect | Optional Law 8 spawn |

### Domain APIs

```text
CreateChange(ctx, ChangeInput) (Change, error)
RecordChangeCommit(ctx, changeID, gitCommit string) (Change, error)
SupersedeChange(ctx, changeID, reason string) (Change, error)
RecordExpectedEffect(ctx, changeID, ExpectedEffectInput) (Effect, error)
RecordActualEffect(ctx, changeID, RecordActualEffectInput) (Effect, *Discovery, error)
GetChange / ListChangesByTaskID / ListChangePaths / ListEffectsByChangeID
ResolveChangePath(ctx, changeID, path string, repo vcs.Repository) ([]byte, error)
```

```text
ChangeInput:
  TaskID          string   // required
  GitCommit       string   // empty → OPEN; non-empty → RECORDED
  ParentChangeID  string
  Actor, Reason   string
  Paths           []ChangePathInput  // ≥1 required
  Expected        []ExpectedEffectInput // optional at create
  DecisionID      string   // optional → change_implements_decision
  SourceType      string
  Confidence      float64

ChangePathInput: Path (required), Status, SymbolID

ExpectedEffectInput: Dimension, Expected (both required), Confidence, SourceType
  comparison must be empty on this path

RecordActualEffectInput:
  Dimension, Actual, Comparison   // all required; Comparison in the 3-value enum
  Confidence                      float64
  EvidenceIDs                     []string  // optional effect_supported_by
  EmitDiscovery                   bool
  DiscoveryTitle, DiscoveryBody   string    // Title required if EmitDiscovery
  HypothesisID                    string    // XOR CreateHypothesis
  CreateHypothesis                bool
  HypothesisTitle, HypothesisBody string    // Title required if CreateHypothesis
```

**ResolveChangePath:** `GetChange` + path must exist on `change_paths`; `git_commit` non-empty; `repo != nil`; return `repo.ShowFile(ctx, git_commit, NormalizePath(path))`. Fail closed otherwise. **Never** persist returned bytes.

Empty `git_commit` → fail closed (cannot resolve). Unknown path on the change → fail closed.

Events: `entity.created` (change); `change.recorded` when SHA first set after create; `effect.compared`; `effect.contradicted`; `entity.created` for optional Discovery/Hypothesis via existing helpers.

### Files (S03-01)

| Path | Role |
|------|------|
| `internal/store/schema/017_changes_effects.sql` | Tables + CHECKs + indexes |
| `internal/store/changes.go` | Upsert/Get/List paths+effects |
| `internal/store/changes_test.go` | Store round-trip + no-blob pragma on new tables |
| `internal/domain/changes.go` | Create/record/resolve APIs + rel constants |
| `internal/domain/changes_test.go` | Named domain tests below |
| `internal/domain/service.go` | `EntityChange`, `EntityEffect`, rels, event names |
| Compat / embed tests listed above | Ceiling **17** |

Do **not** edit `internal/loop`, `cmd/trace`, MCP, `internal/deliberation/select.go`, `internal/gitcli` (callers pass `vcs.Repository`).

### Named tests (minimum — exact names)

1. `TestCreateChangeWithGitSHAAndPathsNoBlob`
2. `TestChangeSchemaHasNoBlobOrPatchColumns`
3. `TestCreateChangeRequiresTaskIDAndPath`
4. `TestRecordExpectedThenActualSupported`
5. `TestRecordActualRequiresExpectedDimension`
6. `TestUnknownEffectComparisonFailClosed`
7. `TestRecordActualContradictedLinksHypothesisWithoutDiscoveryFork`
8. `TestRecordActualContradictedOptionalPlanAffectingDiscovery`
9. `TestContradictedEffectFiresDecisionReconsideration`
10. `TestContradictedEffectDoesNotCreateRegressionOrAutoReplan`
11. `TestParentChangeChain`
12. `TestResolveChangePathViaGitNotSQLite`
13. `TestResolveChangePathFailsClosedWithoutCommit`
14. `TestOversizedEffectTextFailClosed`

## Merge table

COVERAGE.md Change row **thickened** to name `change_paths` + `effects`. Metric/Observation still merge into effect comparison (no Metric table). No gap found that needs a new noun.

## Later scopes (upcoming notes only)

- **S04:** do not add tests/baseline/score columns on `changes`.
- **S05:** contradicted comparison + optional `hypothesis_explains_effect` are inputs; never auto-`caused`.
- **S06:** apply keys `changes`/`effects`; `recent_changes[]` SHA+path+comparison only (no file bytes in packets); ceiling **17** after S03-01.
- **S07:** seed export of `changes`/`effects`/`change_paths` is an explicit residual; VERIFY ceiling = live embed max (17 after S03-01).

## Planner work

1. [x] Pick paths storage: **`change_paths` table** (JSON rejected).
2. [x] Lock `changes` + `effects` schema, comparison enum, contradiction hooks.
3. [x] Lock named tests for 01; compat ceiling **17** after 01.
4. [x] Thicken `01-change-effects.md` + `02-scope-review.md`.
5. [x] Note later scopes; no product Go.

## Exit criteria

- [x] 01/02 thickened with named tests
- [x] Git duplication forbidden in locks (ShowFile read path; no vcs_* dual-write)
- [x] Paths storage picked (`change_paths`)
- [x] No product Go

## Next

Orchestrator: **P20-S03-01** after this row is `done`.
