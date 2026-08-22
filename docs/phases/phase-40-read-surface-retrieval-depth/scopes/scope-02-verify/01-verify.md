# P40-S02-01 — Verify (Phase 40+ gate)

## Metadata
- id: P40-S02-01
- todo_ids: [P40-S02-01]
- role: verifier
- skills: [qa-lead, production-validator, test-driven-development]
- mcps: [Shell, user-trace]
- verification: mixed (test floor + review sign-offs + artifact checklist)

## Objective

Run locked verify blocks **0–6** from [00-PLANNER.md](00-PLANNER.md). Author **`VERIFY-NOTES.md`** with PASS/FAIL per block and evidence manifest. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P40-S02-02**. **No new features.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor (FINAL — S02-00)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S02-02
- [INTAKE.md](../../INTAKE.md)
- [REMEDIATION-PLAN §3 Phase 40+](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- S00–S01 review APPROVE notes on [docs/TODO/phase-40.md](../../../../TODO/phase-40.md)
- Pattern: [P39 S03-01 verify](../../../phase-39-context-orient-harness/scopes/scope-03-verify/01-verify.md)
- Code anchors (live post-S01):
  - G5: `web/src/components/GraphOrientPanel.tsx`, `web/src/lib/orientDismiss.ts`, `web/src/screens/Graph.tsx:465`
  - G5 caps: `web/src/lib/overviewCompose.ts` — SEED_CAP=8, SEED_MAX_NODES=40, UI_CAP=100, EXPAND_MAX_NODES=50, DEPTH=2
  - G2 library: `internal/compiler/explore.go` + `internal/compiler/explore_test.go` (G2-T1–T7)
  - G2 adapters: `internal/mcp/tools_explore.go`, `cmd/trace/explore.go`
  - MCP: `internal/mcp/server.go:228–237`, `internal/mcp/instructions.go`, `internal/mcp/mcp_test.go`

## Session start

Follow agent-loop-protocol Session start. Unattended: run blocks in order; do **not** close DR-HANDOFF or scaffold Phase 41.

## Locked defaults (FINAL — S02-00)

| Item | Value |
|------|-------|
| Precondition | P40-S00-02, P40-S01-02 all `done` + **APPROVE** (high confidence) |
| Product Go / TS / `web/` | **Forbidden** (evidence + notes only) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p40-s02-01-verify/evidence/` |
| Notes artifact | `scopes/scope-02-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S02-02 closes + scaffolds Phase 41+ |
| Phase 41 scaffold | **Out of scope** — S02-02 owns per agent-loop-protocol Phase handoff |
| P40 git window | From P40-S00-01 first product commit through verify run |
| Trace repo root | `/home/ali/Desktop/Trace` |

### Review sign-off precondition (Block 0–1 gate)

| Row | Verdict | Confidence | Notes cite |
|-----|---------|------------|------------|
| P40-S00-02 | **APPROVE** | high | G5-A1–A7; Law 19 GUI adapter; board row 687 |
| P40-S01-02 | **APPROVE** | high | G2-T1–T7; 17 tools; moat preserved; board row 690 |

### G5 accept map (Block 0 — must tick in VERIFY-NOTES)

| ID | Criterion | Primary evidence |
|----|-----------|------------------|
| G5-A1 | First-visit orient panel on `/` Explore | `GraphOrientPanel.tsx` + `data-testid="graph-orient-panel"`; mount `Graph.tsx:465` |
| G5-A2 | Moat-first narrative Tasks→Loop→gate→review | Copy grep `GraphOrientPanel.tsx` |
| G5-A3 | Dismiss persistence | `orientDismiss.ts` key `trace.orient.dismissed`; `orientDismiss.test.ts` 3/3 |
| G5-A4 | Budget/confidence labels on law/truncation/budget UX | `graph-law-banner-confidence`, `graph-truncation-confidence`, `graph-budget-confidence` testids |
| G5-A5 | Law 19 adapter — no new retrieval logic in `web/` | Caps unchanged in `overviewCompose.ts`; no orient in `internal/httpapi/` |
| G5-A6 | Install hook narrative | CONTRIBUTING graph-first GUI `:70`; `bootstrap_hint.go` graph-first line |
| G5-A7 | No graph regression | `overviewCompose.test.ts` 7/7; `npm run build` OK |

### G2 accept map (Block 1 — must tick in VERIFY-NOTES)

| ID | Test / check | Assert summary |
|----|--------------|----------------|
| G2-T1 | `TestExploreTaskRequired` | Empty task_id rejected |
| G2-T2 | `TestExploreTaskMoatPreserved` | Task identity + task_summary in response |
| G2-T3 | `TestExploreQueryMerged` | Query merges via G1 path |
| G2-T4 | `TestExploreCappedHonest` | Truncation flags honest |
| G2-T5 | `TestExploreNoDump` | No unbounded dump |
| G2-T6 | `TestExploreWhyIncluded` | Bounded why/neighborhood |
| G2-T7 | `TestExploreFailOpenSearch` | Search fail-open |
| G2-T1-MCP | `TestMCPExploreTaskRequired` | MCP validation |
| G2-T3-MCP | `TestMCPExploreQueryMerged` | MCP query field |
| G2-INST | `TestServerInstructionsExploreOptional` | explore after moat; 9/17 hygiene |
| Tool count | `TestToolNamesRegistered` | **17** tools; last = `trace_explore` |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Block 0: any G5-A1–A7 not met; S00-02 not APPROVE
- Block 1: any G2-T1–T7 or MCP mirror fails; tool count ≠ 17; S01-02 not APPROVE; stale hygiene still `9/16`
- Block 2: M-001 violation (query-only explore, 1-tool facade, graph-only drift, explore replaces moat lead)
- Block 3: full-graph dump default; cap default inflation; business logic in MCP/GUI adapters (Explore logic must live in `internal/compiler/`)
- Block 4: G6/G7 accidentally implemented without board spawn; secondary queue missing from DR-HANDOFF
- Block 5: Phase 41+ successor TBD in VERIFY-NOTES
- Block 6: P40 entity schema changed but `trace/graph.json` not updated
- VERIFY-NOTES missing or evidence dir absent after claimed PASS
- Product code shipped in this row

**Do not fail VERIFY solely for:**

| Topic | Disposition |
|-------|-------------|
| G6/G7 not implemented | Expected — secondary queue only |
| G8/G9 not started | Phase 41+ per DR-HANDOFF |
| HTTP explore route absent | Residual — MCP+CLI shipped; HTTP mirrors context/search only |
| Phase 41 folder not exists yet | S02-02 scaffold |
| Redundant double `dismissOrient()` on G5 dismiss | Low nit from S00-02 — idempotent |

## Locked verify command floor

Run from Trace repo root. Tee outputs into evidence dir.

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p40-s02-01-verify/evidence"
mkdir -p "$EVID"

{
  echo "verify_id=P40-S02-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=S00-02+S01-02 APPROVE high confidence"
  echo "phase=P40 G5+G2 delivery"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists with metadata.

---

### Block 0 — G5 GUI orient (G5-A1–A7 + S00-02 APPROVE)

```bash
cd /home/ali/Desktop/Trace

# Board sign-off
grep -n 'P40-S00-02.*APPROVE\|687.*APPROVE' docs/TODO/phase-40.md \
  | tee "$EVID/00-board-s00-approve.txt"

# Orient component + mount point
test -f web/src/components/GraphOrientPanel.tsx
grep -n 'graph-orient-panel\|GraphOrientPanel' \
  web/src/components/GraphOrientPanel.tsx web/src/screens/Graph.tsx \
  | tee "$EVID/00-g5-orient-file.txt"

# Moat copy + dismiss key
grep -nE 'Tasks|Loop|gate|review|trace\.orient\.dismissed|localStorage' \
  web/src/components/GraphOrientPanel.tsx web/src/lib/orientDismiss.ts \
  | tee "$EVID/00-g5-copy-grep.txt"

# Confidence testids
grep -n 'graph-law-banner-confidence\|graph-truncation-confidence\|graph-budget-confidence' \
  web/src/screens/Graph.tsx \
  | tee "$EVID/00-g5-confidence-testids.txt"

# Law 19 — no orient logic in httpapi
grep -rn 'orient' internal/httpapi/ 2>/dev/null \
  | tee "$EVID/00-g5-law19-httpapi.txt" || echo "OK: no orient in httpapi" | tee "$EVID/00-g5-law19-httpapi.txt"

# Caps unchanged (presentation constants only)
grep -nE 'SEED_CAP|SEED_MAX_NODES|UI_CAP|EXPAND_MAX_NODES|DEPTH' \
  web/src/lib/overviewCompose.ts \
  | tee "$EVID/00-g5-caps.txt"

# Install narrative
grep -nE 'Graph-first GUI|trace serve|trace gui|Explore' CONTRIBUTING.md \
  | head -15 | tee "$EVID/00-g5-contributing.txt"
grep -n 'graph-first GUI' internal/install/bootstrap_hint.go \
  | tee "$EVID/00-g5-bootstrap-hint.txt"

# Web unit tests (node:test — no vitest in package.json)
(cd web && node --experimental-strip-types --test src/lib/orientDismiss.test.ts 2>&1) \
  | tee "$EVID/00-g5-orient-dismiss-test.txt"
(cd web && node --experimental-strip-types --test src/lib/overviewCompose.test.ts 2>&1) \
  | tee "$EVID/00-g5-overview-compose-test.txt"
(cd web && npm run build 2>&1) | tee "$EVID/00-g5-web-build.txt"
```

**Block 0 evidence table (fill in VERIFY-NOTES):**

| Check | Expected | Evidence file |
|-------|----------|---------------|
| S00-02 APPROVE | Board row 687 `done` + APPROVE high | `00-board-s00-approve.txt` |
| G5-A1 | Panel + testid; mount `:465` | `00-g5-orient-file.txt` |
| G5-A2 | Moat copy Tasks→Loop→gate→review | `00-g5-copy-grep.txt` |
| G5-A3 | Dismiss key + 3/3 tests | `00-g5-orient-dismiss-test.txt` |
| G5-A4 | Confidence testids | `00-g5-confidence-testids.txt` |
| G5-A5 | Caps unchanged; no httpapi orient | `00-g5-caps.txt`, `00-g5-law19-httpapi.txt` |
| G5-A6 | CONTRIBUTING + bootstrap hint | `00-g5-contributing.txt`, `00-g5-bootstrap-hint.txt` |
| G5-A7 | overviewCompose 7/7 + build OK | `00-g5-overview-compose-test.txt`, `00-g5-web-build.txt` |

---

### Block 1 — G2 unified explore (G2-T1–T7 + 17 tools + S01-02 APPROVE)

Re-run G2 acceptance + MCP mirrors. Library lives in **`internal/compiler/explore.go`** (not `retrieval/` — import-cycle move at S01-01).

```bash
cd /home/ali/Desktop/Trace

# G2 acceptance tests (compiler)
go test ./internal/compiler/... -count=1 -v \
  -run 'TestExploreTaskRequired|TestExploreTaskMoatPreserved|TestExploreQueryMerged|TestExploreCappedHonest|TestExploreNoDump|TestExploreWhyIncluded|TestExploreFailOpenSearch' \
  2>&1 | tee "$EVID/01-g2-acceptance.txt"

# MCP mirrors + instructions
go test ./internal/mcp/... -count=1 -v \
  -run 'TestMCPExploreTaskRequired|TestMCPExploreQueryMerged|TestServerInstructionsExploreOptional|TestToolNamesRegistered' \
  2>&1 | tee "$EVID/01-g2-mcp-acceptance.txt"

# Library-first spot-read
grep -n 'func Explore\|ExploreOpts\|DefaultExplore' internal/compiler/explore.go \
  | head -20 | tee "$EVID/01-g2-library-spot-read.txt"

# Thin adapters
wc -l internal/mcp/tools_explore.go cmd/trace/explore.go \
  | tee "$EVID/01-g2-adapter-size.txt"
grep -n 'compiler.Explore' internal/mcp/tools_explore.go cmd/trace/explore.go \
  | tee "$EVID/01-g2-adapter-wiring.txt"

# MCP wiring — 17th tool ReadOnlyHint
grep -c 'AddTool' internal/mcp/server.go | tee "$EVID/01-g2-addtool-count.txt"
grep -n 'trace_explore\|ReadOnlyHint' internal/mcp/server.go internal/mcp/tools_explore.go \
  | tee "$EVID/01-g2-explore-wiring.txt"

# Instructions — explore optional + 9/17 hygiene
grep -n 'trace_explore\|compose-first\|9/17\|stale' internal/mcp/instructions.go \
  | tee "$EVID/01-g2-instructions.txt"

# Stale hygiene in install docs (S01-02 inline fix)
grep -n '9/17' CONTRIBUTING.md internal/install/cursor.go \
  | tee "$EVID/01-g2-stale-hygiene-docs.txt"

# Board sign-off
grep -n 'P40-S01-02.*APPROVE\|690.*APPROVE' docs/TODO/phase-40.md \
  | tee "$EVID/01-board-s01-approve.txt"
```

**Block 1 evidence table:**

| ID | Expected | Evidence |
|----|----------|----------|
| G2-T1–T7 | All PASS | `01-g2-acceptance.txt` |
| G2-T1-MCP / T3-MCP | PASS | `01-g2-mcp-acceptance.txt` |
| G2-INST | explore after moat; 9/17 | `01-g2-mcp-acceptance.txt`, `01-g2-instructions.txt` |
| Tool count | 17; last = trace_explore | `01-g2-mcp-acceptance.txt`, `01-g2-addtool-count.txt` |
| Library | `compiler.Explore` in explore.go | `01-g2-library-spot-read.txt` |
| Adapters | Thin MCP + CLI pass-through | `01-g2-adapter-wiring.txt` |
| S01-02 | APPROVE on board | `01-board-s01-approve.txt` |

---

### Block 2 — M-001 moat preserved

```bash
cd /home/ali/Desktop/Trace

# Explore requires task_id (library + MCP)
grep -n 'task_id\|TaskID\|required' internal/compiler/explore.go internal/mcp/tools_explore.go \
  | tee "$EVID/02-moat-task-required.txt"

# Moat lead in Instructions — explore optional after compose block
grep -n 'trace_tasks\|trace_context\|trace_loop\|trace_review\|trace_explore\|Optional convenience' \
  internal/mcp/instructions.go \
  | tee "$EVID/02-moat-instructions.txt"

# Task summary in explore response (G2-T2)
grep -n 'task_summary\|TaskSummary' internal/compiler/explore.go internal/compiler/explore_test.go \
  | tee "$EVID/02-moat-task-summary.txt"

# Write tools still registered (8 write + 9 read = 17)
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered -v \
  2>&1 | tee "$EVID/02-moat-tool-count.txt"

# No CG 1-tool facade language
grep -rn '1-tool\|one tool only\|codegraph only' internal/ cmd/ CONTRIBUTING.md AGENTS.md 2>/dev/null \
  | tee "$EVID/02-moat-no-facade.txt" || echo "OK: no facade language" | tee "$EVID/02-moat-no-facade.txt"

# DR-HANDOFF M-001 forward note
grep -n 'M-001\|moat\|never replace\|query-only' \
  docs/phases/phase-40-read-surface-retrieval-depth/DR-HANDOFF.md \
  | tee "$EVID/02-moat-dr-handoff.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| task_id required | Empty rejected at library + MCP |
| Moat lead | trace_tasks/context before trace_explore in Instructions |
| Explore optional | “Optional convenience” / “does not replace” language |
| Tool count | 17 — explore is 17th read-only, not facade replacement |
| Dual-stack | CG optional complement — Trace moat primary |

---

### Block 3 — Laws 6–7 caps honest; Law 19 library-first

```bash
cd /home/ali/Desktop/Trace

# Compiler packet defaults unchanged (G1 path)
grep -n 'DefaultTokenBudget\|DefaultMaxItems\|MaxCandidateHits' internal/compiler/packet.go \
  | tee "$EVID/03-law67-packet-caps.txt"

# Explore-specific caps (library constants)
grep -n 'DefaultExploreSearchLimit\|MaxExploreSearchLimit\|DefaultExploreMaxNodes\|DefaultExploreDepth' \
  internal/compiler/explore.go \
  | tee "$EVID/03-law67-explore-caps.txt"

# GUI caps unchanged
grep -n 'SEED_CAP\|UI_CAP\|EXPAND_MAX_NODES' web/src/lib/overviewCompose.ts \
  | tee "$EVID/03-law67-gui-caps.txt"

# No dump API
go test ./internal/compiler/... -count=1 -run 'TestNoDumpAPI|TestExploreNoDump' -v \
  2>&1 | tee "$EVID/03-law67-no-dump.txt"

# Law 19 — Explore logic in compiler; thin adapters
grep -n 'compileAtDepth\|ContextOptions' internal/compiler/compiler.go \
  | head -10 | tee "$EVID/03-law19-compiler-context.txt"
wc -l internal/mcp/tools_explore.go cmd/trace/explore.go \
  web/src/components/GraphOrientPanel.tsx \
  | tee "$EVID/03-law19-adapter-size.txt"

# Reject — no explore business logic in web/
grep -rn 'Explore\|neighborhood\|FTS' web/src/ --include='*.ts' --include='*.tsx' \
  | grep -v overviewCompose | grep -v GraphOrientPanel \
  | tee "$EVID/03-law19-no-web-retrieval.txt" || echo "OK: no web retrieval logic" | tee "$EVID/03-law19-no-web-retrieval.txt"
```

**Pass asserts:**

| Law | Expected |
|-----|----------|
| 6–7 | Packet 4096/32/64; explore limits 32/64/100/depth2; GUI SEED_CAP=8/UI_CAP=100; truncation honest |
| 19 | G2 in `internal/compiler/explore.go`; MCP/CLI/GUI adapters thin |

---

### Block 4 — Secondary queue (G6/G7 not P40 implement)

```bash
cd /home/ali/Desktop/Trace

grep -A25 'Secondary queue' docs/phases/phase-40-read-surface-retrieval-depth/DR-HANDOFF.md \
  | tee "$EVID/04-secondary-queue.txt"

# Scope folders — S00–S02 only (no S03/S04 G6/G7 rows)
ls docs/phases/phase-40-read-surface-retrieval-depth/scopes/ \
  | tee "$EVID/04-scope-folders.txt"

# Board — no G6/G7 implement rows
grep -n 'G6\|G7\|S03\|S04' docs/TODO/phase-40.md \
  | tee "$EVID/04-board-no-g6g7-rows.txt"
```

**Pass:** DR-HANDOFF lists G6/G7 as forward-only; scope folders = `scope-00-gui-graph-orient`, `scope-01-unified-explore`, `scope-02-verify` only.

---

### Block 5 — Successor Phase 41+ (notes only — do not scaffold)

In VERIFY-NOTES, include **successor recommendation table** for S02-02:

| Field | Locked value (from DR-HANDOFF) |
|-------|--------------------------------|
| Successor | **Phase 41+ — Layers & intent** (human promotes P41-00) |
| Entry themes | **G8** Progressive layers L2–L3 + **G9** Intent pipeline |
| Secondary queue | G6, G7 per REMEDIATION-PLAN rank — document in Phase 41 INTAKE |
| P40 outcome | **G5+G2 delivered** — GUI orient + unified trace_explore |
| Idle alternative | `no successor` — if human defers Phase 41 |

```bash
grep -n 'Phase 41\|G8\|G9\|successor\|Secondary queue' \
  docs/phases/phase-40-read-surface-retrieval-depth/DR-HANDOFF.md \
  docs/phases/phase-40-read-surface-retrieval-depth/INTAKE.md \
  | tee "$EVID/05-successor-grep.txt"

grep -n 'Phase 41\|G8\|G9' \
  docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md \
  | head -15 | tee "$EVID/05-successor-remediation-plan.txt"
```

**Pass:** VERIFY-NOTES Block 5 cites **Phase 41+ G8+G9** by name; **never TBD**.

---

### Block 6 — Graph export (if applicable)

```bash
cd /home/ali/Desktop/Trace

# P40 product window — entity schema changes?
git log --oneline --since=2026-08-22 -- \
  internal/store/ internal/graph/ trace/graph.json pkg/domain/ \
  2>&1 | tee "$EVID/06-entity-commits.txt"

if grep -q . "$EVID/06-entity-commits.txt" 2>/dev/null; then
  go build -o /tmp/trace ./cmd/trace
  /tmp/trace seed export -o trace/graph.json 2>&1 | tee "$EVID/06-graph-export.txt"
  git diff --stat trace/graph.json | tee "$EVID/06-graph-diff-stat.txt"
else
  echo "N/A — no entity schema commits in P40 window (expected — G5/G2 adapter-only)" \
    | tee "$EVID/06-graph-export-na.txt"
fi
```

**Pass:** N/A with note if no entity changes (expected per S01-01 Notes); or `trace/graph.json` updated if schema changed.

---

### Full test floor (aggregate)

```bash
cd /home/ali/Desktop/Trace
go test ./internal/retrieval/... ./internal/compiler/... ./internal/mcp/... ./cmd/trace/... -count=1 \
  2>&1 | tee "$EVID/go-test-p40-full.txt"
```

---

### WRITE VERIFY-NOTES.md

Create `docs/phases/phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — Phase 40 / S02-01

**Date:** …
**Overall:** PASS | FAIL | PARTIAL
**Git SHA:** …
**Evidence:** experiments/runs/…-p40-s02-01-verify/evidence/

## Precondition cites
- P40-S00-02 APPROVE (high) — G5 GUI orient
- P40-S01-02 APPROVE (high) — G2 unified explore

## Block results
| Block | Check | Result | Evidence file |
|-------|-------|--------|---------------|
| 0 | G5 G5-A1–A7 + S00 APPROVE | | 00-g5-*.txt |
| 1 | G2 G2-T1–T7 + MCP + 17 tools + S01 APPROVE | | 01-g2-*.txt |
| 2 | M-001 moat | | 02-moat-*.txt |
| 3 | Laws 6–7 / 19 | | 03-law*.txt |
| 4 | G6/G7 forward-only | | 04-*.txt |
| 5 | Phase 41+ successor prep | | 05-successor-*.txt |
| 6 | Graph export | PASS/N/A | 06-graph-*.txt |

## G5 accept map
| ID | Result | Evidence |
| G5-A1 | | |
| … | | |

## G2 accept map
| ID | Result | Evidence |
| G2-T1 | | |
| … | | |

## Successor recommendation (for S02-02)
- **Default:** Phase 41+ — G8 layers L2–L3 + G9 intent pipeline
- **Secondary:** G6, G7 per DR-HANDOFF — Phase 41 INTAKE queue only
- **Never:** TBD

## DR-HANDOFF
Stays OPEN — P40-S02-02 closes + scaffolds Phase 41+

## Next
P40-S02-02
```

## Exit criteria

- [ ] `VERIFY-NOTES.md` with block table 0–6 + G5/G2 accept maps
- [ ] Evidence dir populated under `experiments/runs/…-p40-s02-01-verify/evidence/`
- [ ] Blocks 0–6 executed in order
- [ ] Board Notes on **P40-S02-01** only
- [ ] DR-HANDOFF remains OPEN
- [ ] Next: **P40-S02-02**

## Next

`P40-S02-02`
