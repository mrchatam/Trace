# P27-S01-02 — Protocol v2 review

## Metadata
- id: P27-S01-02
- todo_ids: [P27-S01-02]
- role: reviewer
- skills: [code-review-and-quality, investigator]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

Independent review of INT-08/10 protocol v2 implementation against **locked S01 planner defaults** (P27-S01-00). Confirm harness-only scope, two-session rubric, and `score.sh` preflight behavior before S02 graph-honesty work. **Fresh subagent** — do not reuse S01-01 session.

## References

- [01-implement.md](01-implement.md) — task acceptance criteria + test plan
- [00-PLANNER.md](00-PLANNER.md)
- [../scope-00-investigation/AUDIT.md](../scope-00-investigation/AUDIT.md)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Phase 26 evidence: [`p26-s05-score.txt`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Verify checklist — locked defaults

### Scope boundary (INT-07 deferred to S02)

- [ ] `git diff` touches **only** `experiments/ab-p25-gap-pass-validation/**` (and scope docs if any)
- [ ] **No** changes to `cmd/trace/seed.go`, `internal/loop/gate.go`, `internal/domain/seed_export.go`
- [ ] Implementer did not add `--enforce` strict failures that depend on unimplemented S02 honesty rules

### Export automation = score.sh preflight-only

- [ ] `prepare.sh` has **no** `seed export` invocation
- [ ] `prepare.sh` prints guard/note about score preflight export
- [ ] `score.sh` exports from DB when `trace/graph.json` missing (G1 path)
- [ ] Existing export file is not overwritten blindly when present (preflight only when missing)

### `--strict` = warn-only

- [ ] `score.sh` runs `seed export … --strict` without `--enforce`
- [ ] Strict stderr surfaced as WARN; does not increment FAIL count
- [ ] No G2 failure tied to strict output alone

### P25-3a / P25-3b split

- [ ] `RUBRIC.md` defines P25-3a (build, FAIL expected) and P25-3b (directed, PASS required)
- [ ] Expected-failures / verdict matrix mention build-only thin graph explicitly
- [ ] `score.sh` labels checks `P25-3a` vs `P25-3b` (not generic P25-3 only)

### Arm isolation = `--arm build|directed`

- [ ] `./score.sh G1 --p25` defaults to build arm (backward compatible)
- [ ] `--arm directed` scores P25-3b without mixing build-only verdict
- [ ] Usage documents `--arm build|directed`
- [ ] Invalid arm exits non-zero

### FM-07 git-sparsity = warn-only

- [ ] Reads `exported_at_commit` from graph JSON
- [ ] Compares to workspace HEAD when git available
- [ ] Emits WARN on drift; does not FAIL G2
- [ ] Skips gracefully when no git / no field

### Session-B deliverables

- [ ] `prompts/PROMPT-G1-DIRECTED-GAP.md` exists with Trace env + directed gap instruction
- [ ] `PROTOCOL.md` documents Session A vs Session B score commands
- [ ] P25-4 attestation documented per arm; build arm invalidation if gap prompt early

## Live spot-check commands (reviewer runs)

From repo root:

```bash
cd /home/ali/Desktop/Trace

# Scope hygiene — only experiments harness
git diff --name-only
git diff --name-only | grep -v '^experiments/ab-p25-gap-pass-validation/' | grep -v '^docs/phases/phase-27' || true

# prepare: no export, has guard note
grep -n 'seed export\|preflight\|score.sh' experiments/ab-p25-gap-pass-validation/prepare.sh

# score: preflight, strict warn, arm flag, P25-3a/3b
grep -n 'preflight\|--strict\|--arm\|P25-3a\|P25-3b\|exported_at_commit\|FM-07\|--enforce' \
  experiments/ab-p25-gap-pass-validation/score.sh

# rubric split
grep -n 'P25-3a\|P25-3b\|build-only\|directed' experiments/ab-p25-gap-pass-validation/RUBRIC.md

# protocol two-session
grep -n 'Session\|--arm\|directed\|build-only\|PROMPT-G1-DIRECTED' \
  experiments/ab-p25-gap-pass-validation/PROTOCOL.md

# Session-B prompt exists
test -f experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md
```

### Functional dry runs

```bash
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
cd experiments/ab-p25-gap-pass-validation
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace

# Backward compat
./score.sh G1 --p25 2>&1 | tee /tmp/p27-s01-review-build.txt | grep -E 'P25-3a|VERDICT'

# Directed arm
./score.sh G1 --p25 --arm directed 2>&1 | tee /tmp/p27-s01-review-directed.txt | grep -E 'P25-3b|VERDICT'

# Preflight export (destructive to graph file — use prepared G1 only)
# rm -f runs/G1/trace/graph.json && ./score.sh G1 --test 2>&1 | grep -E 'G2|export'
```

- [ ] At least 3 grep spot-checks match implementer claims
- [ ] One live `./score.sh G1 --p25 --arm build` produces P25-3a-labeled output
- [ ] B0 scoring path still works: `./score.sh B0` exits 0 or expected FAIL without arm errors

## Confidence rubric

| Level | When to use |
|-------|-------------|
| **high** | All locked defaults verified; spot-checks + dry run match; no open blocker/high |
| **medium** | Minor doc nits (line refs, wording); core behavior correct; residual risks listed in Notes |
| **low** | Product code touched, enforce strict added, prepare exports, or arms conflated — **do not APPROVE** |

## Spawn policy

| Severity | Action |
|----------|--------|
| **blocker** | Product files changed; `--enforce` strict fails G2; prepare.sh exports; arms not isolated → spawn `P27-S01-02a` (fix) + `P27-S01-02b` (re-review) **immediately below this row** |
| **high** | P25-3a/3b missing or score defaults wrong → inline fix if ≤10 lines bash/markdown; else spawn 02a/02b |
| **medium** | PROTOCOL/RUBRIC drift, missing FM-07 skip path → spawn unless trivial fix |
| **low / nit** | Note only in review Notes |

Insert spawned rows in `docs/TODO/phase-27.md` directly under order 458. Spawned prompts follow agent-loop-protocol skeleton.

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
- [ ] Next runnable: **P27-S02-00** (only after APPROVE)

## Verdict

`APPROVE` | `REQUEST_CHANGES`

**APPROVE** only when INT-08/10 harness changes match P27-S01-00 locks and Phase 26 build-only P25-3a FAIL remains reproducible as documented baseline.

## Todo updates

Status + notes on **P27-S01-02** only.

## Next

`P27-S02-00`
