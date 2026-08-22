# P27-S01-01 — Protocol v2 implementation (INT-08/10)

## Metadata
- id: P27-S01-01
- todo_ids: [P27-S01-01, S01-T01, S01-T02, S01-T03, S01-T04, S01-T05, S01-T06, S01-T07]
- role: implementer
- skills: [incremental-implementation, test-driven-development, shell]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

Implement experiment **protocol v2** under `experiments/ab-p25-gap-pass-validation/`: automate export preflight in `score.sh`, split build-only vs directed-gap scoring arms (P25-3a/3b), document two-session flow in PROTOCOL/RUBRIC, add Session-B prompt — per **locked S01 planner defaults** (P27-S01-00). **No product code** (`cmd/trace/seed.go`, `internal/loop/gate.go`, etc.) — INT-07 belongs to **S02**.

## References

- [00-PLANNER.md](00-PLANNER.md) — planner row (locked defaults source)
- [../scope-00-investigation/AUDIT.md](../scope-00-investigation/AUDIT.md) — S01-T01..T07 seeds, Phase 26 residual
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- Phase 26 evidence:
  - [`experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt)
  - [`p26-export-snippet.json`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P27-S01-00 — do not re-debate)

| Item | Locked value | Rationale |
|------|--------------|-----------|
| **Scope boundary** | `experiments/ab-p25-gap-pass-validation/**` only | INT-07 product honesty = **S02**; no `seed.go` / `gate.go` / `seed_export.go` edits |
| **Export automation** | **`score.sh` preflight-only** | Export from DB when `trace/graph.json` missing; do **not** add post-install export to `prepare.sh` (would stamp stale seed-only graph before agent session) |
| **`--strict` in harness** | **Warn-only** (`trace seed export -o … --strict`, no `--enforce`) | Today's product strict does not fail thin graphs; enforce belongs after S02 honesty rules |
| **P25-3 split** | **P25-3a** build-only (FAIL = expected baseline) + **P25-3b** directed-gap (PASS required) | Preserves Phase 26 verdict matrix; separates install wiring from gap-pass behavior |
| **Arm isolation** | **`score.sh G1 --p25 --arm build\|directed`** | Default `--arm build` when `--p25`; `directed` scores Session-B graph only; clearer than `--session-b` |
| **FM-07 git-sparsity** | **Warn-only** | Compare `exported_at_commit` vs workspace `git rev-parse HEAD`; WARN if behind, do not FAIL G2 |
| **P25-4 attestation** | Operator records in `experiments/RESULTS.md` notes per arm; score prints reminder | Invalidates build arm if gap prompt sent before build score |
| **Backward compat** | `./score.sh G1 --p25` without `--arm` behaves as `--arm build` | Phase 26 repro still valid |

## Hard boundary

```text
ALLOWED:  experiments/ab-p25-gap-pass-validation/{PROTOCOL,RUBRIC,score.sh,prepare.sh,prompts/**}
FORBIDDEN: cmd/trace/seed.go, internal/loop/gate.go, internal/domain/seed_export.go,
           internal/domain/seed_eval_rules_test.go, cmd/trace/enforce_test.go
```

If product changes seem required, **stop** and note in board row — schedule under S02, do not implement here.

## Preflight

Run from repo root before editing:

```bash
cd /home/ali/Desktop/Trace

# Harness (edit targets)
test -f experiments/ab-p25-gap-pass-validation/PROTOCOL.md
test -f experiments/ab-p25-gap-pass-validation/RUBRIC.md
test -x experiments/ab-p25-gap-pass-validation/score.sh
test -x experiments/ab-p25-gap-pass-validation/prepare.sh
test -f experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-BUILD.md

# Trace binary (preflight export needs it)
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
test -x bin/trace

# Phase 26 anchor (read-only)
test -f experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt
```

If any path missing, mark row `blocked` with path in Notes.

## Files to change

| Path | Task IDs | Change summary |
|------|----------|----------------|
| `experiments/ab-p25-gap-pass-validation/score.sh` | T01–T03, T06 | Preflight export, `--strict` warn, FM-07 warn, `--arm build\|directed`, P25-3a/3b branching |
| `experiments/ab-p25-gap-pass-validation/prepare.sh` | T01 | Add **guard message only** (no export): remind operator that score preflight exports if missing |
| `experiments/ab-p25-gap-pass-validation/PROTOCOL.md` | T04, T07 | Mandatory export flow; two-session arms; P25-4 invalidation rules |
| `experiments/ab-p25-gap-pass-validation/RUBRIC.md` | T05, T07 | P25-3a/3b rows; verdict matrix update; expected-failures table names build-only thin graph |
| `experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md` | T07 | **New** Session-B prompt (directed gap pass) |

**Do not create** product Go changes or modify `experiments/RESULTS.md` on this row unless repro requires a scratch note (operator-owned).

## Per-task acceptance criteria

### S01-T01 — Export guard (prepare + score preflight)

**prepare.sh**
- [ ] After G1 install block, print one line: `note: score.sh exports trace/graph.json from DB if missing before scoring`
- [ ] **No** `trace seed export` invocation in `prepare.sh`

**score.sh (G1 arm only)**
- [ ] Before G2/P25 checks: if `$WS/trace/graph.json` missing, run `"$TRACE_BIN" -C "$WS" seed export -o trace/graph.json`
- [ ] If export command fails, `fail "G2 graph export — preflight export failed"` with repro hint
- [ ] If file existed, skip auto-export (respect agent/operator export)

### S01-T02 — `--strict` warn preflight

- [ ] After graph file exists, run `"$TRACE_BIN" -C "$WS" seed export -o trace/graph.json --strict` (capture stderr)
- [ ] Print `WARN strict: …` for each stderr line; **do not** increment FAILS or use `--enforce`
- [ ] Document in score header comment: strict warn is harness hygiene until S02 graph-honesty rules

### S01-T03 — FM-07 git-sparsity (warn-only)

- [ ] Read `exported_at_commit` from graph JSON (python3 or jq one-liner)
- [ ] If workspace has git: compare to `git -C "$WS" rev-parse HEAD`
- [ ] If export SHA ≠ HEAD: `echo "WARN  FM-07 exported_at_commit ($export_sha) behind HEAD ($head_sha) — re-export recommended"`
- [ ] If no git or field absent: skip with `skip "FM-07 git-sparsity — no git or exported_at_commit"`
- [ ] **Never** fail G2 on drift alone

### S01-T04 — PROTOCOL.md two-session + export steps

- [ ] §1 Prepare: step documenting that **post-session** export is automated by `score.sh` preflight if missing; agent may export during session per prompt
- [ ] §3 G1: rename/clarify **Session A (build-only)** — score with `./score.sh G1 --p25 --arm build`
- [ ] §4 (new or expanded): **Session B (directed-gap)** — when to run, score with `./score.sh G1 --p25 --arm directed`
- [ ] §4 Optional confirmation run text upgraded from footnote to first-class Session-B section with cross-link to new prompt

### S01-T05 — RUBRIC.md P25-3a / P25-3b

- [ ] Replace single P25-3 row with:
  - **P25-3a** (build arm): `discoveries≥1 OR decisions≥1` — **expected FAIL** on build-only G1 (documented baseline)
  - **P25-3b** (directed arm): same threshold — **PASS required** for gap-pass validation
- [ ] Update expected-failures table to name `discoveries=0 decisions=0` on build-only explicitly
- [ ] Update verdict matrix: P25-3a FAIL + P25-1/2 PASS → "Install OK; behavior unchanged"; P25-3b PASS + P25-1/2 PASS → "P25-C validated (directed)"

### S01-T06 — `--arm build|directed` in score.sh

- [ ] Parse `--arm build` or `--arm directed` (default `build` when `--p25`)
- [ ] `--arm build`: run P25-3a check; label output `P25-3a`
- [ ] `--arm directed`: run P25-3b check; label output `P25-3b`; require graph from **after** Session-B (same file path; operator responsibility)
- [ ] Update usage: `usage: $0 [B0|G1] [--test] [--p25] [--arm build|directed]`
- [ ] Unknown `--arm` → exit 2 with message

### S01-T07 — Session-B prompt + P25-4 formalization

- [ ] Add `prompts/PROMPT-G1-DIRECTED-GAP.md`: single directed message "Run the mandatory gap pass from your Trace rules" (+ Trace env block matching G1-BUILD)
- [ ] PROTOCOL: gap prompt **before** build score invalidates build arm (operator must `./prepare.sh G1` fresh)
- [ ] RUBRIC P25-4: separate attestation rows for build vs directed arms in RESULTS.md template text
- [ ] score.sh: `skip "P25-4 operator attestation — record in RESULTS.md (arm=…)"` with arm name

## Repro steps (implementer self-check)

Baseline Phase 26 residual (build-only FAIL expected):

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
cd experiments/ab-p25-gap-pass-validation
./prepare.sh G1

# Simulate post-session: remove export if present to test preflight
rm -f runs/G1/trace/graph.json

export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
./score.sh G1 --p25 --arm build 2>&1 | tee /tmp/s01-build-score.txt
# Expect: preflight export runs; P25-3a FAIL if graph thin (discoveries=0 decisions=0)
grep -E 'P25-3a|preflight|WARN|FM-07' /tmp/s01-build-score.txt
```

Directed arm dry run (with existing G1 workspace graph):

```bash
./score.sh G1 --p25 --arm directed 2>&1 | tee /tmp/s01-directed-score.txt
grep -E 'P25-3b|WARN strict' /tmp/s01-directed-score.txt
```

## Test plan (score.sh dry runs)

| # | Command | Expected |
|---|---------|----------|
| 1 | `./score.sh G1 --p25` | Same as `--arm build`; P25-3a label in output |
| 2 | `./score.sh G1 --p25 --arm build` | P25-3a check; strict WARN lines if any; FM-07 WARN if SHA drift |
| 3 | `./score.sh G1 --p25 --arm directed` | P25-3b check only (not 3a) |
| 4 | `rm -f runs/G1/trace/graph.json && ./score.sh G1 --test` | Preflight creates export; G2 PASS |
| 5 | `./score.sh B0 --test` | Unchanged B0 path (no regression) |
| 6 | `./score.sh G1 --p25 --arm foo` | Exit 2 usage error |

No `go test` required unless product files touched (they must not be).

## Role work

1. Implement tasks **T01→T07** in order (score.sh core first, then docs, then new prompt).
2. Keep bash `set -euo pipefail`; match existing `pass`/`fail`/`skip` style.
3. Self-check repro + test plan table.
4. Update board row **P27-S01-01** status + Notes only (no future prompt edits).

## Todo updates

Per board-rights: **P27-S01-01** status + notes only.

## Exit criteria

- [ ] All S01-T01..T07 acceptance checkboxes satisfied
- [ ] Locked defaults table honored (especially: no product code, warn-only strict/FM-07, score preflight-only export)
- [ ] `./score.sh G1 --p25 --arm build` runs on live workspace with labeled P25-3a output
- [ ] `./score.sh G1 --p25 --arm directed` runs with labeled P25-3b output
- [ ] PROTOCOL + RUBRIC document two-session flow and P25-3a/3b semantics
- [ ] `prompts/PROMPT-G1-DIRECTED-GAP.md` exists
- [ ] Board row `done` with evidence paths in Notes

## Minimal todos

- [ ] T01: score.sh preflight export + prepare.sh guard message
- [ ] T02: score.sh `--strict` warn pass
- [ ] T03: score.sh FM-07 warn
- [ ] T04: PROTOCOL.md export + two-session sections
- [ ] T05: RUBRIC.md P25-3a/3b + verdict matrix
- [ ] T06: score.sh `--arm build|directed`
- [ ] T07: PROMPT-G1-DIRECTED-GAP.md + P25-4 attestation text
- [ ] Self-check repro + test plan
- [ ] Board P27-S01-01 done

## Next

`P27-S01-02`
