# P38-S05-01 — Saturation assessment

## Metadata
- id: P38-S05-01
- todo_ids: [P38-S05-01]
- role: implementer
- skills: [analyst, research, doubt-driven-development]
- mcps: [user-trace]
- verification: mixed
- hooks: none

## Objective

Author **`SATURATION-NOTES.md`**. Decide: **exit investigation** (`PROCEED_TO_S06`) or **spawn more S01–S04 rows** (`SPAWN`). This is the **only authorized exit** from investigation loops. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — § Saturation exit criteria (authoritative checklist)
- [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md) — §2–§7 (APPROVED input)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §4.4, §5 spawn rules
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner locks from P38-S05-00
- **Upstream summaries (read-only — do not re-audit):**
  - [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md)
  - [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md)
  - [PEER-UA-GF.md](../scope-03-ua-graphify-peer/PEER-UA-GF.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P38-S05-00 — do not re-debate)

| Item | Value |
|------|-------|
| Output path | `scopes/scope-05-saturation-gate/SATURATION-NOTES.md` |
| Product edits | **Forbidden** (Go/TS/web) |
| Method | Checklist synthesis of APPROVED GAP-REGISTRY + DESIGN-LOCKS — **not** full S01–S03 re-audit |
| Evidence root | `experiments/runs/YYYY-MM-DD-p38-s05-663/evidence/` (desk-check + saturation notes only) |
| H7 compose-equivalence | **Defer-with-trigger** — run T4 desk-check; **do not spawn** S01/S02 unless T4 reveals uncovered structural dimension |
| Verdict vocabulary | `PROCEED_TO_S06` \| `SPAWN` (with row sketches) |
| Blocks S06 | Until S05-02 **APPROVE (saturated)** |

## Must answer (embed in SATURATION-NOTES.md)

1. **Checklist** — pass/fail per DESIGN-LOCKS § Saturation exit criteria (all six boxes) with evidence pointers.
2. **Confidence** — `high` \| `medium` \| `low` + rationale: would another S01–S03 row discover **new** gaps or duplicate?
3. **Rejected spawns** — duplicate / out-of-scope / law-conflict ideas with one-line reason each.
4. **`ready_for_REMEDIATION_PLAN`** — explicit boolean (`true` only if recommendation is `PROCEED_TO_S06`).

---

## Investigation todos (run in order)

### T0 — Preflight + evidence folder

```bash
EV=experiments/runs/$(date +%Y-%m-%d)-p38-s05-663/evidence
mkdir -p "$EV"
```

- Confirm GAP-REGISTRY exists and S04-02 APPROVED (board row 661).
- Confirm S01–S03 artifacts APPROVED (rows 652, 655, 658).
- Record date + row id in every new evidence file header.

### T1 — DESIGN-LOCKS checklist walk (item-by-item)

**Read:** DESIGN-LOCKS § Saturation exit criteria (six bullets).

For **each** criterion, record in `$EV/t1-saturation-checklist.md`:

| # | Criterion | Pass/Fail | Evidence pointer (link, don't duplicate) |
|---|-----------|-----------|------------------------------------------|
| 1 | H1–H11 verified / rejected / deferred | | GAP-REGISTRY §2 rows G-001…G-011; §5 defer/reject |
| 2 | Live Trace command per major gap | | GAP-REGISTRY §7 → S01 `$EV/` paths |
| 3 | Peer mechanism cites CG, UA, GF | | GAP-REGISTRY §2 dual-side; S02/S03 mechanism files |
| 4 | Moat row (Trace strengths) | | GAP-REGISTRY §3 M-001 |
| 5 | Spawn empty or deferred with trigger | | §6 + T4/T5 defer closures |
| 6 | High confidence — further rows duplicate | | T6 reject list |

**Rule:** Link existing `$EV/` from S01–S04. Re-run live CLI only if spot-checking a **contested** claim (unlikely post-S04 APPROVE).

### T2 — Live Trace coverage confirmation (link audit)

**Do not re-run full S01 audit.** Confirm GAP-REGISTRY §7 index maps major gaps to S01 live evidence:

| Gap cluster | S01 evidence file(s) | Live command class |
|-------------|----------------------|-------------------|
| G-001, G-002, G-003 | h1-*, h2-*, h3-* | `trace context`, compiler path |
| G-005, G-006, G-009, G-010 | h5-*, h6-*, h9-*, h10-* | `trace index`, MCP inventory, doc grep |
| G-008 partial | h8-gui-partial | GUI route read (+ S03 resolved) |

**Verdict target:** `$EV/t2-live-coverage-confirm.md` — table: G-ID \| S01 evidence link \| sufficient Y/N.

**Spot-check (optional, one command max):** e.g. `go test ./internal/mcp/ -run TestToolNamesRegistered -count=1` — only if checklist item 2 contested.

### T3 — Peer mechanism cite confirmation (link audit)

Confirm **CG, UA, GF** have mechanism cites (not README-only) via GAP-REGISTRY §2:

| Peer | Mechanism evidence | Scope |
|------|-------------------|-------|
| CG | `h7-explore-mechanism.md`, `tools.ts` cites | S02 |
| UA | `context-builder.ts`, `h1-ua-partial.md` | S03 |
| GF | `symbol_resolution.py`, `h4-gf-extracted-inferred.md` | S03 |
| MP (bonus) | `searcher.py`, `layers.py` | S03 |

**Verdict target:** `$EV/t3-peer-mechanism-confirm.md`.

### T4 — H7 compose-equivalence desk-check (planner lock)

**Purpose:** Close GAP-REGISTRY §6 open trigger **without spawn** unless desk-check finds uncovered dimension.

**Read (do not re-run peer stack):**

- S02: [`h7-explore-gap.md`](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h7-explore-gap.md)
- S02: [`h7-explore-mechanism.md`](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h7-explore-mechanism.md)
- S01: TRACE-AUDIT H1 partial, H2; `h1-trace-partial.md`

**Compare dimensions** (table in evidence file):

| Dimension | CG `codegraph_explore` | Trace compose (`trace_search` + `trace_why` + `trace_context` [+ `trace_impact`]) |
|-----------|------------------------|-------------------------------------------------------------------------------------|
| Single capped call | | |
| Query-driven orient | | |
| Verbatim source grouped by file | | |
| Call path among symbols | | |
| Blast radius summary | | |
| Task-scoped packet merge | | |
| Tool count / discovery cost | | |

**Expected outcome:** **Not equivalent** — document structural deltas (already evidenced). If equivalent on all dimensions → **stop** and note SPAWN need (unlikely).

**Verdict target:** `$EV/h7-compose-desk-check.md`

**Defer trigger (record in SATURATION-NOTES §4):**

- Owner: **S06 REMEDIATION-PLAN**
- Trigger: G-007 remediation sketch must compare unified `trace_explore` vs compose-first UX; optional Phase 39 pre-implement live spike
- **Not** a P38 investigation spawn

### T5 — Remaining “what if X?” triage

Brainstorm residual investigation angles; classify each:

| Candidate | Class | Rationale |
|-----------|-------|-----------|
| H7 live MCP side-by-side | **reject spawn** | T4 desk-check + S02 mechanism sufficient |
| Optional S01-01a symbol richness | **defer** | GAP-REGISTRY §6; S06 optional |
| Cursor 9/16 MCP tools | **defer → S06** | Harness hygiene |
| H12+ peer slice | **reject** | MP mapped in S03 |
| Re-litigate H8 GUI screenshot | **reject** | S03 supported; T5 skip documented |
| DR-NOSSEM vector channel | **reject** | G-004a law defer |
| Full CG live explore on Trace | **reject** | No `.codegraph/`; out of scope |

**Verdict target:** `$EV/t5-residual-triage.md`

### T6 — Duplication confidence + reject list

Answer explicitly:

> Would another S01–S03 investigate row discover **new** evidence-backed gaps, or mostly duplicate APPROVED artifacts?

Document **≥8 rejected spawn ideas** (see planner 00-PLANNER § Rejected spawn ideas as seed).

**Verdict target:** `$EV/t6-duplication-confidence.md`

**Confidence target:** **high** (acceptable: **medium** only if T4 leaves genuine structural unknown — then recommendation may be SPAWN).

### T7 — Author SATURATION-NOTES.md

**Path:** `scopes/scope-05-saturation-gate/SATURATION-NOTES.md`

**Required sections:**

```markdown
## §1 Checklist (DESIGN-LOCKS saturation exit criteria)
Pass/fail per criterion 1–6 + evidence links

## §2 Confidence statement
high | medium | low + rationale (duplication vs discovery)

## §3 Spawn list
Empty = good for exit. If non-empty: row sketches + owner scope.

## §4 Deferred investigations (trigger + owner)
Include H7 compose-equivalence defer → S06 trigger

## §5 Rejected duplicative investigations
Table: idea | reject reason (duplicate | out of scope | law conflict)

## §6 Recommendation
PROCEED_TO_S06 | SPAWN (with sketches)

## §7 ready_for_REMEDIATION_PLAN
Explicit boolean + one-line gate note ("S05-02 APPROVE required")
```

**Link** `$EV/` files; link GAP-REGISTRY §2–§7; **no** REMEDIATION-PLAN content (ranked G1–Gn is S06).

### T8 — Self-check before done

- [ ] All six DESIGN-LOCKS boxes addressed in §1
- [ ] H7 defer documented in §4 with S06 owner (not silent)
- [ ] §5 reject list ≥8 items
- [ ] §6 recommendation consistent with §1 (any FAIL → SPAWN, not S06)
- [ ] §7 boolean matches §6
- [ ] No product code diff
- [ ] No ranked remediation themes

---

## Hard stop

If recommendation is **SPAWN**:

1. **Do not start S06.**
2. Document spawn row sketches in §3 (scope, hypothesis, owner S01–S04).
3. Board reviewer (S05-02) inserts rows per INVESTIGATION-INDEX §5.
4. Re-enter S05 after spawn cycle completes.

## Exit criteria

- [ ] `SATURATION-NOTES.md` §§1–7 complete
- [ ] Evidence under `experiments/runs/…-p38-s05-663/evidence/` (minimum: T1, T4, T6)
- [ ] Board row P38-S05-01 → `done` with Notes (recommendation + confidence)
- [ ] No product code

## Next

`P38-S05-02`
