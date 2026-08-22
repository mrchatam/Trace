# P39-S03-01 — VERIFY Phase 39

## Metadata
- id: P39-S03-01
- todo_ids: [P39-S03-01]
- role: verify
- skills: [qa-lead, test-driven-development, shipping-and-launch]
- mcps: [Shell, user-trace]
- verification: mixed (test floor + review sign-offs + artifact checklist)
- hooks: []

## Objective

Run locked verify blocks **0–6** from [00-PLANNER.md](00-PLANNER.md). Author **`VERIFY-NOTES.md`** with PASS/FAIL per block and evidence manifest. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P39-S03-02**. **No new features.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor (FINAL — S03-00)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S03-02
- [INTAKE.md](../../INTAKE.md)
- [REMEDIATION-PLAN §3 Phase 39](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- S00–S02 review APPROVE notes on [docs/TODO/phase-39.md](../../../../TODO/phase-39.md)
- Pattern: [P37 S03-01 verify](../../../phase-37-p36-residuals/scopes/scope-03-verify/01-verify.md)
- Pattern: [P38 S07-01 verify](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-07-verify/01-verify.md)
- Code anchors: `internal/compiler/compiler.go`, `internal/mcp/instructions.go`, `internal/mcp/server.go`, `CONTRIBUTING.md`, `AGENTS.md`

## Session start

Follow agent-loop-protocol Session start. Unattended: run blocks in order; do **not** close DR-HANDOFF or scaffold Phase 40.

## Locked defaults (FINAL — S03-00)

| Item | Value |
|------|-------|
| Precondition | P39-S00-02, P39-S01-02, P39-S02-02 all `done` + **APPROVE** (high confidence) |
| Product Go / TS / `web/` | **Forbidden** (evidence + notes only) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p39-s03-01-verify/evidence/` |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S03-02 closes + scaffolds Phase 40+ |
| Phase 40 scaffold | **Out of scope** — S03-02 owns per agent-loop-protocol Phase handoff |
| P39 git window | From P39-S00-01 first product commit through verify run |
| Trace repo root | `/home/ali/Desktop/Trace` |

### Review sign-off precondition (Block 0–2 gate)

| Row | Verdict | Confidence | Notes cite |
|-----|---------|------------|------------|
| P39-S00-02 | **APPROVE** | high | G1 T1–T6 + T1-MCP; 16 tools; query merge at `compiler.go:158–165` |
| P39-S01-02 | **APPROVE** | high | G3-A1–A6; `ServerInstructions()` wired; moat + compose + 9/16 |
| P39-S02-02 | **APPROVE** | high | G4-D1–D8 doc-only; H11; S01 coordination PASS |

### G1 accept map (Block 0 — must tick in VERIFY-NOTES)

| ID | Test name | Assert summary |
|----|-----------|----------------|
| T1 | `TestG1QueryHitMerged` | Query token merges hit; empty Query excludes |
| T2 | `TestG1TaskMoatPreserved` | Layer 0 task + task_state with Query; empty task_id rejected |
| T3 | `TestG1TitleFTSStillRunsWithQuery` | Title + query channels both contribute |
| T4 | `TestG1QueryExpandDedupe` | Expand + query same entity → one item |
| T5 | `TestG1QueryCapHonesty` | Many matches → `Truncated` / cap honest |
| T6 | `TestG1QuerySearchFailOpen` | Query search error → packet still valid (DF-87) |
| T1-MCP | `TestMCPContextQueryMerged` | MCP `query` field merges hit |

### G3 accept map (Block 1 — must tick in VERIFY-NOTES)

| ID | Criterion | Primary evidence |
|----|-----------|------------------|
| G3-A1 | Non-empty MCP Instructions via go-sdk | `TestServerInstructionsNonEmpty` |
| G3-A2 | Moat lead: tasks → context(+query) → loop → review → plan | `TestServerInstructionsMoatLead` |
| G3-A3 | Compose-first read recipe | `TestServerInstructionsComposeRecipe` |
| G3-A4 | `trace_version` + 9/16 stale hygiene | `TestServerInstructionsStaleHygiene` |
| G3-A5 | CONTRIBUTING moat-first + reload | Doc grep `:68–72` |
| G3-A6 | 16 tools unchanged | `TestToolNamesRegistered` |

### G4 accept map (Block 2 — must tick in VERIFY-NOTES)

| ID | Requirement | Primary doc |
|----|-------------|-------------|
| G4-D1 | Complementary dual-stack title | `CONTRIBUTING.md` dual-stack heading |
| G4-D2 | `.trace/` vs `.codegraph/` storage | Same section |
| G4-D3 | When to use Trace (moat) | Same + cross-link S01 |
| G4-D4 | When to use Codegraph (optional) | Same |
| G4-D5 | Law 19 adapter boundaries | Same |
| G4-D6 | Setup recipe (independent paths) | Same |
| G4-D7 | Not shipping rejects | Same |
| G4-D8 | PEER-CG / PEER-FIXTURES links | Same |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Block 0: any G1 T1–T6 or T1-MCP test fails; S00-02 not APPROVE
- Block 1: any G3-A1–A6 not met; MCP tool count ≠ 16; S01-02 not APPROVE
- Block 2: G4 product code in P39 S02 diff; G4-D1–D8 incomplete; S02-02 not APPROVE
- Block 3: M-001 violation (query-only orient, 1-tool facade, gate bypass, new explore tool)
- Block 4: full-graph dump default added; compiler caps changed without justification; business logic in MCP/CLI adapters
- Block 5: DR-HANDOFF forward note missing G5/G2 secondary queue; Phase 40+ successor still TBD in VERIFY-NOTES
- Block 6: P39 entity schema changed but `trace/graph.json` not updated
- VERIFY-NOTES missing or evidence dir absent after claimed PASS
- Product code shipped in this row

**Do not fail VERIFY solely for:**

| Topic | Disposition |
|-------|-------------|
| G5 GUI orient not implemented | Forward queue — Phase 40+ |
| G2 unified `trace_explore` absent | Forward queue — Phase 40+ after G1 + law spike |
| Phase 40 folder not yet exists | S03-02 scaffold — not S03-01 |
| `instructions.go:25` "Phase 39 S02" stub | Optional P40 doc hygiene (S02-02 nit) |
| G6/G7/G8/G9 not started | Secondary queue per DR-HANDOFF |

## Locked verify command floor

Run from Trace repo root. Tee outputs into evidence dir.

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p39-s03-01-verify/evidence"
mkdir -p "$EVID"

{
  echo "verify_id=P39-S03-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=S00-02+S01-02+S02-02 APPROVE high confidence"
  echo "phase=P39 G1+G3+G4 delivery"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists with metadata.

---

### Block 0 — G1 context merge (T1–T6 + S00-02 APPROVE)

Re-run G1 acceptance + regression subset.

```bash
cd /home/ali/Desktop/Trace

# G1 acceptance tests
go test ./internal/compiler/... -count=1 -v \
  -run 'TestG1QueryHitMerged|TestG1TaskMoatPreserved|TestG1TitleFTSStillRunsWithQuery|TestG1QueryExpandDedupe|TestG1QueryCapHonesty|TestG1QuerySearchFailOpen' \
  2>&1 | tee "$EVID/00-g1-acceptance.txt"

# MCP query mirror
go test ./internal/mcp/... -count=1 -v -run TestMCPContextQueryMerged \
  2>&1 | tee "$EVID/00-g1-mcp-query.txt"

# Regression subset (must stay green)
go test ./internal/compiler/... -count=1 -v \
  -run 'TestTaskContextContinuesWhenSearchErrors|TestCandidateCapSetsTruncated|TestItemCapNeverExceeded|TestBudgetLoudTotals|TestNoDumpAPI' \
  2>&1 | tee "$EVID/00-g1-regression.txt"

# Board sign-off
grep -n 'P39-S00-02.*APPROVE\|674.*APPROVE' docs/TODO/phase-39.md \
  | tee "$EVID/00-board-s00-approve.txt"

# Spot-read merge point + caps
grep -n 'Query\|Search(ctx, q\|ReasonFTSMatch\|DefaultTokenBudget\|DefaultMaxItems\|MaxCandidateHits' \
  internal/compiler/compiler.go internal/compiler/packet.go \
  | head -30 | tee "$EVID/00-g1-spot-read.txt"

# CLI + MCP wiring
grep -n 'query\|Query' cmd/trace/context.go internal/mcp/tools_context.go \
  | tee "$EVID/00-g1-adapter-wiring.txt"
```

**Block 0 evidence table (fill in VERIFY-NOTES):**

| Check | Expected | Evidence file |
|-------|----------|---------------|
| S00-02 APPROVE | Board row 674 `done` + APPROVE high | `00-board-s00-approve.txt` |
| T1–T6 | All PASS | `00-g1-acceptance.txt` |
| T1-MCP | PASS | `00-g1-mcp-query.txt` |
| Regression | PASS (DF-87, caps, no dump) | `00-g1-regression.txt` |
| Merge point | Query merge after title FTS, before file-seed | `00-g1-spot-read.txt` |
| Adapters | `--query` CLI; MCP `query` optional; task_id required | `00-g1-adapter-wiring.txt` |

---

### Block 1 — G3 harness orient (G3-A1–A6 + 16 tools + S01-02 APPROVE)

```bash
cd /home/ali/Desktop/Trace

# G3 acceptance tests
go test ./internal/mcp/... -count=1 -v \
  -run 'TestServerInstructionsNonEmpty|TestServerInstructionsMoatLead|TestServerInstructionsComposeRecipe|TestServerInstructionsStaleHygiene|TestToolNamesRegistered' \
  2>&1 | tee "$EVID/01-g3-acceptance.txt"

# Instructions wiring
grep -n 'ServerInstructions\|ServerOptions\|Instructions:' \
  internal/mcp/server.go internal/mcp/instructions.go \
  | tee "$EVID/01-g3-instructions-wiring.txt"

# Moat + compose + hygiene string targets
grep -n 'trace_tasks\|trace_context\|trace_loop\|trace_review\|trace_plan\|trace_version\|trace_search\|9/16\|stale\|compose\|Codegraph' \
  internal/mcp/instructions.go \
  | tee "$EVID/01-g3-instructions-content.txt"

# CONTRIBUTING moat-first + reload (S01 region)
grep -n 'moat\|reload\|trace_version\|9/16\|dual-stack' CONTRIBUTING.md \
  | head -20 | tee "$EVID/01-g3-contributing.txt"

# 16 tools — no 17th AddTool
grep -c 'AddTool' internal/mcp/server.go | tee "$EVID/01-g3-addtool-count.txt"
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered -v \
  2>&1 | tee "$EVID/01-g3-tool-count.txt"

# Board sign-off
grep -n 'P39-S01-02.*APPROVE\|677.*APPROVE' docs/TODO/phase-39.md \
  | tee "$EVID/01-board-s01-approve.txt"

# Reject checks — no explore tool, no tool reduction
grep -rn 'trace_explore' internal/mcp/ cmd/trace/ 2>/dev/null \
  | tee "$EVID/01-g3-no-explore.txt" || echo "OK: no trace_explore" | tee "$EVID/01-g3-no-explore.txt"
```

**Block 1 evidence table:**

| ID | Expected | Evidence |
|----|----------|----------|
| G3-A1 | Instructions non-empty | `01-g3-acceptance.txt` |
| G3-A2 | Moat tool order in Instructions | `01-g3-instructions-content.txt` |
| G3-A3 | Compose-first recipe | same |
| G3-A4 | trace_version + 9/16 | same |
| G3-A5 | CONTRIBUTING expanded | `01-g3-contributing.txt` |
| G3-A6 | 16 tools | `01-g3-tool-count.txt` |
| S01-02 | APPROVE on board | `01-board-s01-approve.txt` |

---

### Block 2 — G4 dual-stack docs (G4-D1–D8 docs-only + S02-02 APPROVE)

```bash
cd /home/ali/Desktop/Trace

# Doc-only boundary for S02 scope
git log --oneline --since=2026-08-22 --name-only -- CONTRIBUTING.md AGENTS.md \
  2>&1 | tee "$EVID/02-g4-doc-commits.txt"

git diff --stat $(git rev-list -1 --before="2026-08-22" HEAD)..HEAD -- CONTRIBUTING.md AGENTS.md \
  2>&1 | tee "$EVID/02-g4-doc-diff-stat.txt"

# Confirm no S02 product paths touched
git log --oneline --since=2026-08-22 -- internal/ cmd/ web/ trace/ \
  2>&1 | tee "$EVID/02-g4-product-commits-since-s02.txt"

# G4 checklist greps
grep -n 'Trace + Codegraph\|optional dual-stack\|complement' CONTRIBUTING.md \
  | tee "$EVID/02-g4-d1-title.txt"
grep -n '\.trace/\|\.codegraph/' CONTRIBUTING.md \
  | tee "$EVID/02-g4-d2-storage.txt"
grep -n 'Law 19\|adapter\|does not index' CONTRIBUTING.md \
  | tee "$EVID/02-g4-d5-law19.txt"
grep -n 'Not shipping\|no default dual-index\|bundled' CONTRIBUTING.md \
  | tee "$EVID/02-g4-d7-rejects.txt"
grep -n 'PEER-CG\|PEER-FIXTURES' CONTRIBUTING.md \
  | tee "$EVID/02-g4-d8-links.txt"
grep -n 'Optional Codegraph\|complement' AGENTS.md \
  | tee "$EVID/02-g4-agents-subsection.txt"

# Link resolve (relative paths exist)
{
  test -f docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md \
    && echo "OK PEER-CG" || echo "MISSING PEER-CG"
  test -f docs/phases/phase-38-retrieval-context-peer-gaps/PEER-FIXTURES.md \
    && echo "OK PEER-FIXTURES" || echo "MISSING PEER-FIXTURES"
} | tee "$EVID/02-g4-link-resolve.txt"

# Board sign-off
grep -n 'P39-S02-02.*APPROVE\|680.*APPROVE' docs/TODO/phase-39.md \
  | tee "$EVID/02-board-s02-approve.txt"
```

**Block 2 pass asserts:**

| Check | Expected |
|-------|----------|
| S02 diff | `CONTRIBUTING.md` + `AGENTS.md` only (no `internal/` in S02 implement) |
| G4-D1–D8 | Each grep file non-empty or spot-read PASS |
| S01 coordination | Moat-first `:68–72` preserved; dual-stack anchor linked |
| S02-02 APPROVE | Board row 680 |

Copy diff stat → `$EVID/doc-diff-stat.txt` (alias for planner reference).

---

### Block 3 — M-001 moat preserved

```bash
cd /home/ali/Desktop/Trace

# Task loop + gate tools present
grep -n 'trace_loop\|trace_tasks\|trace_review\|trace_transition\|action=gate' \
  internal/mcp/server.go internal/mcp/tools_loop.go internal/mcp/instructions.go \
  | tee "$EVID/03-moat-tools.txt"

# No query-only orient path (task_id still required)
grep -n 'task_id\|TaskID\|required' internal/mcp/tools_context.go \
  | tee "$EVID/03-moat-task-required.txt"

# No CG 1-tool facade language in product
grep -rn '1-tool\|one tool only\|codegraph only' internal/ cmd/ CONTRIBUTING.md AGENTS.md 2>/dev/null \
  | tee "$EVID/03-moat-no-facade.txt" || echo "OK: no facade language" | tee "$EVID/03-moat-no-facade.txt"

# DR-HANDOFF M-001 forward note
grep -n 'M-001\|moat\|never replace' docs/phases/phase-39-context-orient-harness/DR-HANDOFF.md \
  | tee "$EVID/03-moat-dr-handoff.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| Task loop tools | `trace_tasks`, `trace_loop`, `trace_review`, gate path documented |
| Query additive | Optional query on context — not query-only packet |
| Tool count | Still 16 — no `trace_explore` |
| Dual-stack | CG optional complement — Trace moat primary |

---

### Block 4 — Laws 6–7 caps honest; Law 19 library-first

```bash
cd /home/ali/Desktop/Trace

# Default caps unchanged
grep -n 'DefaultTokenBudget\|DefaultMaxItems\|MaxCandidateHits' internal/compiler/packet.go \
  | tee "$EVID/04-law67-caps.txt"

# No dump API
go test ./internal/compiler/... -count=1 -run TestNoDumpAPI -v \
  2>&1 | tee "$EVID/04-law67-no-dump.txt"

# Law 19 — merge in compiler; thin CLI/MCP
grep -n 'compileAtDepth\|ContextOptions' internal/compiler/compiler.go \
  | head -15 | tee "$EVID/04-law19-compiler.txt"
wc -l cmd/trace/context.go internal/mcp/tools_context.go \
  | tee "$EVID/04-law19-adapter-size.txt"

# HTTP/GUI unchanged for orient (G5 deferred)
git log --oneline --since=2026-08-22 -- web/ internal/httpapi/ \
  2>&1 | tee "$EVID/04-law19-no-gui-orient.txt"
```

**Pass asserts:**

| Law | Expected |
|-----|----------|
| 6–7 | Caps 4096/32/64 defaults; truncation honest; no full-graph dump |
| 19 | G1 logic in `internal/compiler/`; CLI/MCP pass-through only |

---

### Block 5 — Successor documented (Phase 40+ G5/G2 — notes only)

In VERIFY-NOTES, include **successor recommendation table** for S03-02 (do **not** close DR-HANDOFF or scaffold Phase 40):

| Field | Locked value (from DR-HANDOFF forward note) |
|-------|---------------------------------------------|
| Successor | **Phase 40+ — Read surface & retrieval depth** (human promotes P40-00) |
| Entry themes | **G5** GUI graph orient start + **G2** unified `trace_explore` |
| Secondary queue | G6, G7 per REMEDIATION-PLAN rank; G8/G9 Phase 41+ |
| P39 outcome | **G1+G3+G4 delivered** — compose-first shipped in S01 Instructions |
| Idle alternative | `no successor` — if human defers Phase 40 |

```bash
grep -n 'Phase 40\|G5\|G2\|Secondary queue' \
  docs/phases/phase-39-context-orient-harness/DR-HANDOFF.md \
  | tee "$EVID/05-successor-dr-handoff.txt"

grep -n 'Phase 40' \
  docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md \
  | head -10 | tee "$EVID/05-successor-remediation-plan.txt"
```

**Pass:** VERIFY-NOTES Block 5 cites Phase 40+ by name; G5 + G2 listed; **never TBD**.

---

### Block 6 — Graph export (if applicable)

```bash
cd /home/ali/Desktop/Trace

# P39 product window — entity schema changes?
git log --oneline --since=2026-08-22 -- \
  internal/store/ internal/graph/ trace/graph.json pkg/ \
  2>&1 | tee "$EVID/06-entity-commits.txt"

# If commits present, export required
if grep -q . "$EVID/06-entity-commits.txt" 2>/dev/null; then
  go build -o /tmp/trace ./cmd/trace
  /tmp/trace seed export -o trace/graph.json 2>&1 | tee "$EVID/06-graph-export.txt"
  git diff --stat trace/graph.json | tee "$EVID/06-graph-diff-stat.txt"
else
  echo "N/A — no entity schema commits in P39 window" | tee "$EVID/06-graph-export-na.txt"
fi
```

**Pass:** N/A with note if no entity changes (expected — G1/G3 adapter-only); or `trace/graph.json` updated if schema changed.

---

### Block 7 — Full test floor (optional aggregate)

```bash
cd /home/ali/Desktop/Trace
go test ./internal/compiler/... ./internal/mcp/... ./cmd/trace/... -count=1 \
  2>&1 | tee "$EVID/go-test-p39-full.txt"
```

---

### Block 8 — WRITE VERIFY-NOTES.md

Create `docs/phases/phase-39-context-orient-harness/scopes/scope-03-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — Phase 39 / S03-01

**Date:** …
**Overall:** PASS | FAIL | PARTIAL
**Git SHA:** …
**Evidence:** experiments/runs/…-p39-s03-01-verify/evidence/

## Precondition cites
- P39-S00-02 APPROVE (high) — G1
- P39-S01-02 APPROVE (high) — G3
- P39-S02-02 APPROVE (high) — G4 doc-only

## Block results
| Block | Check | Result | Evidence file |
|-------|-------|--------|---------------|
| 0 | G1 T1–T6 + T1-MCP + S00 APPROVE | | 00-g1-*.txt |
| 1 | G3 G3-A1–A6 + 16 tools + S01 APPROVE | | 01-g3-*.txt |
| 2 | G4 G4-D1–D8 docs-only + S02 APPROVE | | 02-g4-*.txt |
| 3 | M-001 moat | | 03-moat-*.txt |
| 4 | Laws 6–7 / 19 | | 04-law*.txt |
| 5 | Phase 40+ successor prep | | 05-successor-*.txt |
| 6 | Graph export | PASS/N/A | 06-graph-*.txt |

## G1 accept map
| ID | Result | Evidence |
| T1 | | |
| … | | |

## G3 accept map
| ID | Result | Evidence |
| G3-A1 | | |
| … | | |

## G4 accept map
| ID | Result | Evidence |
| G4-D1 | | |
| … | | |

## Successor recommendation (for S03-02)
- **Default:** Phase 40+ — G5 GUI orient + G2 unified explore
- **Secondary:** G6, G7 per DR-HANDOFF
- **Never:** TBD

## DR-HANDOFF
Stays OPEN — P39-S03-02 closes + scaffolds Phase 40+

## Next
P39-S03-02
```

## Exit criteria

- [ ] `VERIFY-NOTES.md` with block table 0–6 + G1/G3/G4 accept maps
- [ ] Evidence dir populated under `experiments/runs/…-p39-s03-01-verify/evidence/`
- [ ] Blocks 0–6 executed in order
- [ ] Board Notes on **P39-S03-01** only
- [ ] DR-HANDOFF remains OPEN
- [ ] Next: **P39-S03-02**

## Next

`P39-S03-02`
