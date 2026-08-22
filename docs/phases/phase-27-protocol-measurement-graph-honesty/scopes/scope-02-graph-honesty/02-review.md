# P27-S02-02 — Graph honesty review

## Metadata
- id: P27-S02-02
- todo_ids: [P27-S02-02]
- role: reviewer
- skills: [code-review-and-quality, investigator, silent-failure-hunter]
- mcps: [user-codegraph, user-trace]
- verification: mixed
- hooks: []

## Objective

Independent review of **INT-07** graph honesty implementation against **locked S02 planner defaults** (P27-S02-00). Confirm product-only scope, thin-graph `--strict --enforce` behavior, and gate parity before S03 VERIFY. **Fresh subagent** — do not reuse S02-01 session.

## References

- [01-implement.md](01-implement.md) — task acceptance criteria + test commands
- [00-PLANNER.md](00-PLANNER.md)
- [../scope-00-investigation/AUDIT.md](../scope-00-investigation/AUDIT.md) — S02 seeds, P26 root cause
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Phase 26 baseline: [`p26-export-snippet.json`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Verify checklist — locked defaults

### Scope boundary (harness deferred to S03)

- [ ] Diff touches **only** product paths listed in implement hard boundary (`cmd/trace/`, `internal/domain/`, optional `cmd/trace/testdata/`)
- [ ] **No** edits to `experiments/ab-p25-gap-pass-validation/**`
- [ ] **No** daemon/HTTP/MCP feature work

### Thin-graph strict rule (A)+(B) hybrid

- [ ] `collectExportGraphHonestyViolations` (or equivalent) exists and runs on `--strict`
- [ ] `discoveries=0 && decisions=0` produces document-level violation with stable message substring
- [ ] When `discoveries>0`, orphan discovery (no `discovery_mentions_task`) produces violation
- [ ] When `decisions>0`, orphan decision (no `decision_affects_task`) produces violation
- [ ] **Not** driven by eval-rules.json body evaluation

### Export vs done gate split = (A) seed.go honesty only

- [ ] `internal/loop/gate.go`: `GateForExport` still routes to `evaluateDone` (L60–61 unchanged in behavior)
- [ ] **No** new `evaluateExport` function unless spawn-approved
- [ ] Thin P26 graph blocked by honesty layer, not solely by changing STOP/ORIENT gate logic

### BLOCKING uncertainty rule = (A) store-backed promotion

- [ ] BLOCKING severity read from **store** (not export `SeedEntity`)
- [ ] BLOCKING discovery without `discovery_mentions_task` link → violation
- [ ] No false requirement for uncertainty row when promotion link present

### Done-task gate skip = (A) keep skip

- [ ] `collectExportViolations` still skips done/skipped/stale tasks (`seed.go` ~L121–124)
- [ ] No `--strict`-only full gate sweep added in gate.go

### Strict vs enforce behavior

- [ ] `--strict` alone: violations on stderr, exit 0, file written (when other checks pass)
- [ ] `--strict --enforce`: violations → `exitGateBlocked`, no output file
- [ ] Matches existing verification-debt enforce pattern in `enforce_test.go`

### Harness interaction = S03 VERIFY owns enforce upgrade

- [ ] `score.sh` unchanged (no `--enforce` added on S02-01 row)
- [ ] Implementer did not flip harness to fail G2 on strict

### Test fixtures = P26 baseline

- [ ] P26 export snippet (or equivalent) used in test/fixture
- [ ] Thin-graph test: `--strict --enforce` → blocked
- [ ] `TestSeedExportStrictCleanAllowsWrite` passes (fixture has ≥1 decision/discovery + link)
- [ ] `TestEvaluateGate_Export_SameAsDone` passes **unchanged**

## Live spot-check commands (reviewer runs)

From repo root:

```bash
cd /home/ali/Desktop/Trace

# Scope hygiene — no experiments harness
git diff --name-only
git diff --name-only | grep -v '^cmd/trace/' | grep -v '^internal/' | grep -v '^docs/phases/phase-27' || true

# Honesty collector wired in strict path
grep -n 'GraphHonesty\|graph honesty\|collectExportGraph' cmd/trace/seed.go

# Gate parity unchanged
grep -n 'GateForExport\|evaluateDone\|evaluateExport' internal/loop/gate.go

# Gate skip preserved
grep -n 'WorkStateDone\|WorkStateSkipped\|WorkStateStale' cmd/trace/seed.go

# BLOCKING store check
grep -n 'BLOCKING\|SeverityBlocking\|discovery_mentions_task' cmd/trace/seed.go internal/domain/*honesty* 2>/dev/null || true

# Tests
grep -n 'p26-export\|thin.graph\|discoveries=0\|graph honesty' cmd/trace/enforce_test.go internal/domain/*honesty* 2>/dev/null || true

# Harness untouched
git diff --name-only | grep 'experiments/ab-p25-gap-pass-validation' && echo UNEXPECTED || echo 'OK: no harness diff'
```

### Functional tests

```bash
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace

go test ./internal/... -count=1
go test ./cmd/trace/... -count=1 -run 'SeedExport|Strict|Enforce|Honesty'

# Gate parity
go test ./internal/loop/... -count=1 -run TestEvaluateGate_Export_SameAsDone
```

- [ ] All test commands PASS
- [ ] At least 4 grep spot-checks match implementer claims
- [ ] Thin-graph enforce test exists and passes

## Confidence rubric

| Level | When to use |
|-------|-------------|
| **high** | All locked defaults verified; tests + grep match; P26 thin graph blocked on enforce; no harness diff |
| **medium** | Minor test/fixture nits; core honesty rules correct; residual risks listed in Notes |
| **low** | Gate split without spawn, harness edited, eval-rules engine added, or thin graph still passes enforce — **do not APPROVE** |

## Spawn policy

| Severity | Action |
|----------|--------|
| **blocker** | Thin graph passes `--strict --enforce`; gate.go split without approval; experiments/ touched; eval-rules body engine → spawn `P27-S02-02a` + `P27-S02-02b` **immediately below this row** |
| **high** | Missing link rules or clean fixture broken → inline fix if ≤20 lines; else spawn 02a/02b |
| **medium** | BLOCKING check incomplete, message strings unstable → spawn unless trivial |
| **low / nit** | Note only in review Notes |

Insert spawned rows in `docs/TODO/phase-27.md` directly under order 461.

## Findings template (use in Notes)

```text
| Severity | Finding | Evidence | Disposition |
|----------|---------|----------|-------------|
| blocker  | …       | path:Lnn | spawn 02a |
| high     | …       | …        | inline / spawn |
```

## Exit criteria

- [ ] All locked-default checklist items verified (commands + diff)
- [ ] No open blocker/high without pending 02a/02b row
- [ ] Confidence **medium** or **high** with evidence in Notes
- [ ] Verdict: `APPROVE` or `REQUEST_CHANGES`
- [ ] Next runnable: **P27-S03-00** (only after APPROVE)

## Verdict

`APPROVE` | `REQUEST_CHANGES`

**APPROVE** only when INT-07 product rules match P27-S02-00 locks and P26-equivalent thin export fails `--strict --enforce` while clean full-cycle export passes.

## Todo updates

Status + notes on **P27-S02-02** only.

## Next

`P27-S03-00`
