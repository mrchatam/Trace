# P28-S02-01 — Session-B run and score (P25-3b dogfood)

## Metadata
- id: P28-S02-01
- todo_ids: [P28-S02-01]
- role: implementer
- skills: [incremental-implementation]
- mcps: []
- agents: []
- verification: mixed
- hooks: []

## Objective

Run **live Session-B** on E02 G1 (directed gap pass) **without wiping Session-A**. Score **`--arm directed`**. Prove **P25-3b** (`discoveries ≥ 1` OR `decisions ≥ 1`) with honesty-clean export. Do **not** treat this as Trace product work (`internal/`, `cmd/trace/`).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md)
- [PROMPT-G1-DIRECTED-GAP.md](../../../../../experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md)
- [PROTOCOL.md](../../../../../experiments/ab-p25-gap-pass-validation/PROTOCOL.md) §4
- [RUBRIC.md](../../../../../experiments/ab-p25-gap-pass-validation/RUBRIC.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R1
- [TEST-MATRIX.md](../scope-01-integration-tests/TEST-MATRIX.md) — M-16 is **not** P25-3b PASS

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Execution mode (locked)

| Mode | When | Who |
|------|------|-----|
| **Default: coding subagent (agent-executable)** | G1 workspace present (live 2026-08-20: **yes**) | P28-S02-01 implementer runs the gap **in** `runs/G1` via `TRACE_BIN -C "$G1"` |
| Human paste (optional equivalent) | Operator prefers Cursor folder = G1 | Open `runs/G1`, paste `PROMPT-G1-DIRECTED-GAP.md`, then same export/score/notes |
| **Blocked (human restore)** | G1 missing / no `.trace/` | Do **not** `./prepare.sh`. Write `SESSION-B-BLOCKED.md` and mark board `blocked` |

A subagent whose Cursor cwd is the **Trace repo** is valid. Do **not** wait for a human to open G1 as the workspace. Isolate with `-C "$G1"` on every `trace` invocation.

**Do not** use MCP `user-trace` against the Trace product repo for this dogfood — that store is not G1 (`trace_context` on the product root has no G1 task). CLI `-C "$G1"` only.

## Locked defaults

| Item | Value |
|------|-------|
| G1 | `/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1` |
| Harness | `/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation` |
| `TRACE_BIN` | `/home/ali/Desktop/Trace/bin/trace` |
| `TRACE_TASK_ID` | `e0200000-0000-4000-8000-000000000010` |
| `TRACE_PROJECT_ROOT` | same as G1 |
| Goal id | `e0200000-0000-4000-8000-000000000001` |
| Directed prompt | `experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md` |
| Score | `./score.sh G1 --p25 --arm directed --test` from harness dir |
| Session-A RESULTS | `experiments/RESULTS.md` E02 row 2026-08-20 P27 verify (thin P25-3a FAIL expected) |
| Evidence dir | this folder (`scope-02-session-b-dogfood/`) |
| Gap CLI | **`trace loop status`** — `trace gap` **does not exist** (install text still says `trace gap` OR `loop status`) |
| Honesty | every discovery → `discovery-mentions-task`; every decision → `decision-task` |
| Product edits | **only** under `runs/G1` (checkout-desk dogfood). Never Trace `internal/` / `cmd/trace/` |
| Out | `./prepare.sh`, `--arm build` after mutation, hook deny (S03), honesty-dup/attestation product (S04), `apply_promotion_test.go` |

## Arm isolation (do not overwrite Session-A as directed)

1. **Before any G1 graph writes:** snapshot `runs/G1/trace/graph.json` → `SESSION-A-GRAPH-SNAPSHOT.json` in this folder.
2. Live Session-A baseline (planner 2026-08-20): goals=1 tasks=5 **discoveries=0 decisions=0**.
3. After Session-B the **on-disk** `graph.json` will be rich. That is expected. **Never** re-run `./score.sh G1 --p25 --arm build` (default `--arm` is build) — that would score the **directed** graph as P25-3a and destroy the thin baseline.
4. Session-A score stays in P27 VERIFY-NOTES + RESULTS E02 P27 row + this snapshot. Session-B is a **new** RESULTS row **`E02-SB`**.
5. P27 already ran directed score on the **thin** graph (P25-3b FAIL OK). Session-B **replaces** that directed measurement only.

## If G1 is missing

```text
Mark P28-S02-01 blocked.
Write SESSION-B-BLOCKED.md in this folder:
  - what is missing (dir, .trace, graph.json, TRACE_BIN)
  - do not run ./prepare.sh (wipes Session-A)
  - human restore: re-run PROTOCOL Session A (prepare G1 + PROMPT-G1-BUILD + score --arm build) THEN re-queue S02-01
```

## Preflight (abort → blocked, not invent G1)

```bash
REPO=/home/ali/Desktop/Trace
HARNESS="$REPO/experiments/ab-p25-gap-pass-validation"
G1="$HARNESS/runs/G1"
SCOPE="$REPO/docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood"

test -d "$G1"
test -d "$G1/.trace"
test -f "$G1/trace/graph.json"
test -f "$HARNESS/prompts/PROMPT-G1-DIRECTED-GAP.md"
test -x "$REPO/bin/trace" || (cd "$REPO" && CGO_ENABLED=1 go build -o bin/trace ./cmd/trace)
grep -q -- '--arm' "$HARNESS/score.sh"
grep -q 'P25-3b' "$HARNESS/score.sh"
grep -q 'E02' "$REPO/experiments/RESULTS.md"
# R1 still the S02 target
grep -q 'R1' "$REPO/docs/phases/phase-28-residuals-validation/scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md"

# Snapshot Session-A graph BEFORE mutation
cp "$G1/trace/graph.json" "$SCOPE/SESSION-A-GRAPH-SNAPSHOT.json"
```

Record snapshot counts in notes (expect disc=0 dec=0). If snapshot already has disc/dec ≥ 1, do **not** wipe; document “partial Session-B already present” and continue to honesty-clean export + directed score.

## Role work — directed gap (mandatory)

Follow installed G1 rules (`.cursor/rules/trace-enforcement.mdc`) **and** `PROMPT-G1-DIRECTED-GAP.md`. Minimum that closes P25-3b:

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0200000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1
G1="$TRACE_PROJECT_ROOT"

# Surface gaps (installed “trace gap” alias — use loop status)
"$TRACE_BIN" -C "$G1" loop status --task "$TRACE_TASK_ID"
"$TRACE_BIN" -C "$G1" loop gate --task "$TRACE_TASK_ID" --for edit
```

Live 2026-08-20: status is **edit blocked** (`plan_missing` / recommended PLAN). That is a real gap, not a reason to wipe G1.

### Required graph richness (P25-3b)

Add **at least one** discovery **or** decision that is honest (linked). Example (adapt titles to actual G1 findings — do not invent empty placeholders):

```bash
# Discovery (print id from command output)
"$TRACE_BIN" -C "$G1" add discovery --title "…" --body "…"
"$TRACE_BIN" -C "$G1" link discovery-mentions-task --from <discovery-id> --to "$TRACE_TASK_ID"

# And/or decision
"$TRACE_BIN" -C "$G1" add decision --title "…" --body "…"
"$TRACE_BIN" -C "$G1" link decision-task --from <decision-id> --to "$TRACE_TASK_ID"
```

Use `trace add --help` / `trace link --help` if flags differ; do not guess past a CLI error.

### Optional (not required for P25-3b)

- Unstick PLAN: `trace plan create-coarse` (or loop apply) then continue gap pass.
- Product edits **in G1 only** after `loop gate --for edit` allows (or after plan exists).
- BLOCKING promotion: `trace add task --from-discovery <id>` or `loop apply` with `spawned_tasks[].discovery_id`.
- `loop reset` if sticky STOP blocks useful work — record in notes.
- Do **not** force `transition DONE` if gate still fails (`hop_budget_exceeded` / plan) — expected residual.

Do **not** edit `runs/B0`.

## Export then score

```bash
"$TRACE_BIN" -C "$G1" seed export -o trace/graph.json --strict --enforce
# must exit 0. If honesty FAIL: add missing links; do not disable --enforce.

cd /home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation
./score.sh G1 --p25 --arm directed --test 2>&1 | tee \
  /home/ali/Desktop/Trace/docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood/SESSION-B-SCORE.txt
```

Expect **P25-1/2 PASS**, **P25-3b PASS** if disc/dec ≥ 1, **G2 honesty PASS** if links present. `--test` runs G1 `go test ./...` (checkout desk). Overall `VERDICT: PASS` is desired; if G-opt tests fail, record and still keep P25-3b line — reviewer decides spawn.

**Forbidden:** `./score.sh G1 --p25` (defaults to **build**).

## Evidence files (this folder)

### `SESSION-B-NOTES.md` (required)

```markdown
# Session-B notes — E02 G1 directed gap

**Date:** YYYY-MM-DD
**Row:** P28-S02-01
**Mode:** agent-executable | human-paste
**TRACE_TASK_ID:** e0200000-0000-4000-8000-000000000010

## Arm isolation
- Snapshot: SESSION-A-GRAPH-SNAPSHOT.json — discoveries=… decisions=…
- Did NOT run ./prepare.sh
- Did NOT re-score --arm build after mutation

## Gap pass
- loop status summary (phase, violations)
- entities added (ids) + links
- optional promotion / G1 product edits (paths)

## Score
- Command: ./score.sh G1 --p25 --arm directed --test
- Transcript: SESSION-B-SCORE.txt
- P25-1 / P25-2 / P25-3b / G2: PASS/FAIL
- graph counts after export: disc=… dec=… tasks=…

## P25-4 attestation (manual until S04)
- Directed gap prompt / this row executed: Y
- No human/agent directed-gap before Session-A score: Y (P27 Session-A already scored)
- No build-only prompt sent during Session-B: Y/N

## git (G1 workspace)
- \`git -C runs/G1 log -1 --oneline\` and HEAD SHA
```

### `experiments/RESULTS.md` — **append** a new row (do not edit P27 E02 cells)

```text
| YYYY-MM-DD | E02-SB | ab-p25-gap-pass-validation | PASS/FAIL | Session-B directed gap; P25-3b …; snapshot SESSION-A-GRAPH-SNAPSHOT.json; score SESSION-B-SCORE.txt; notes …/SESSION-B-NOTES.md |
```

Optional attestation subsection under the table (RUBRIC template) is allowed **in addition to** the row.

## Do not

- Re-run `PROMPT-G1-BUILD.md`
- `./prepare.sh` / `./prepare.sh G1`
- Implement hook deny, honesty-dup, P25-4 harness parser
- Claim M-16 grep as P25-3b
- Mark R1 closed on the board (reviewer S02-02)

## Todo updates

Status + notes on **P28-S02-01** only.

## Exit criteria

- [ ] Preflight PASS or `blocked` with `SESSION-B-BLOCKED.md`
- [ ] `SESSION-A-GRAPH-SNAPSHOT.json` captured before mutation
- [ ] Directed gap executed against G1 (`loop status` + ≥1 linked discovery or decision)
- [ ] `seed export -o trace/graph.json --strict --enforce` exit 0
- [ ] `SESSION-B-SCORE.txt` from `./score.sh G1 --p25 --arm directed --test`
- [ ] `SESSION-B-NOTES.md` with P25-3b line + P25-4 attestation
- [ ] `experiments/RESULTS.md` has a **new** `E02-SB` row
- [ ] No `--arm build` score after graph mutation
- [ ] No Trace product code changes

## Minimal todos

1. Preflight + snapshot.
2. Directed gap in G1 (CLI `-C`).
3. Export `--strict --enforce`.
4. Directed score + tee.
5. Notes + RESULTS append.
6. Board P28-S02-01 done/blocked/failed.

## Next

`P28-S02-02`
