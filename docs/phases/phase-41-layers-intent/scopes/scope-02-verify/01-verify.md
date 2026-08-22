# P41-S02-01 — Verify (G8 + G9)

## Metadata
- id: P41-S02-01
- todo_ids: [P41-S02-01]
- role: verifier
- skills: [qa-lead, production-validator, test-driven-development]
- mcps: [Shell, user-trace]
- verification: mixed (test floor + review sign-offs + artifact checklist)

## Objective

Run locked verify blocks **0–6** from [00-PLANNER.md](00-PLANNER.md). Author **`VERIFY-NOTES.md`** with PASS/FAIL per block and evidence manifest. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P41-S02-02**. **No new features.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor (FINAL — S02-00)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S02-02
- [INTAKE.md](../../INTAKE.md)
- [REMEDIATION-PLAN §3 Phase 41+](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- S00–S01 review APPROVE notes on [docs/TODO/phase-41.md](../../../../TODO/phase-41.md)
- Pattern: [P40 S02-01 verify](../../../phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/01-verify.md)
- Code anchors (live post-S01 — from S00-02 / S01-02 reviews):
  - G8: `internal/compiler/compiler.go:34–38` (`MaxLayer` default 1); `layer_enrich.go`; `retrieval/layer_enrich.go`; `budget.go:29–67` trim L0→L3; `retrieval/doc.go:13–16` reason codes
  - G8 adapters: `cmd/trace/context.go` (`--max-layer`); `internal/mcp/tools_context.go` (`max_layer`); `internal/mcp/mcp_test.go` `TestMCPContextMaxLayer2`
  - G8 tests: `internal/compiler/compiler_test.go` G8-L1–L7 block (~1423–1650); `TestNoDumpAPI:813`
  - G9: `internal/retrieval/intent.go` (`ExtractIntent`, caps 32/16/512/256); `search.go:13–44`; `compiler.go:155–179`; `explore.go:115–117`
  - G9 doc: `docs/RETRIEVAL_AND_CONTEXT.md` §3 intent shipped + DR-NOSSEM semantic defer
  - G9 tests: `internal/retrieval/intent_test.go` G9-I1–I6

## Session start

Follow agent-loop-protocol Session start. Unattended: run blocks in order; do **not** close DR-HANDOFF or scaffold Phase 42.

## Locked defaults (FINAL — S02-00)

| Item | Value |
|------|-------|
| Precondition | P41-S00-02, P41-S01-02 all `done` + **APPROVE** (high confidence) |
| Product Go / TS / `web/` | **Forbidden** (evidence + notes only) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p41-s02-01-verify/evidence/` |
| Notes artifact | `scopes/scope-02-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S02-02 closes + scaffolds Phase 42+ |
| Phase 42 scaffold | **Out of scope** — S02-02 owns per agent-loop-protocol Phase handoff |
| P41 git window | From P41-S00-01 first product commit through verify run |
| Trace repo root | `/home/ali/Desktop/Trace` |

### Review sign-off precondition (Block 0–1 gate)

| Row | Verdict | Confidence | Notes cite |
|-----|---------|------------|------------|
| P41-S00-02 | **APPROVE** | high | G8-L1–L7; opt-in L2/L3; board row **697** |
| P41-S01-02 | **APPROVE** | high | G9-I1–I6; §3 revised; board row **700** |

### G8 accept map (Block 0 — must tick in VERIFY-NOTES)

| ID | Criterion | Primary evidence |
|----|-----------|------------------|
| G8-L1 | Default context → layer max 1 | `TestContextDefaultLayer1` (`compiler_test.go:1423`); `TestTaskContextAndBudgets:61` regression |
| G8-L2 | `max_layer=2` → layer-2 items + reason_codes | `TestContextMaxLayer2` (`:1443`); MCP `TestMCPContextMaxLayer2` (`mcp_test.go:168`) |
| G8-L3 | `max_layer=3` → L3 when graph supports | `TestContextMaxLayer3` (`:1474`); `historical_vcs` via VCS |
| G8-L4 | L2/L3 budget caps | `TestContextLayerBudgetCap` (`:1520`) |
| G8-L5 | Trim priority L0→L3 | `TestContextLayerTrimPriority` (`:1542`) |
| G8-L6 | depth ≠ layer | `TestContextDepthIndependentOfLayer` (`:1594`) |
| G8-L7 | No dump; MaxCandidateHits honored | `TestContextNoDump` (`:1611`); `TestNoDumpAPI` (`:813`) |

### G9 accept map (Block 1 — must tick in VERIFY-NOTES)

| ID | Test / check | Assert summary |
|----|--------------|----------------|
| G9-I1 | `TestExtractIntentFromTask` | Keywords from task fields |
| G9-I2 | `TestExtractIntentEntityHints` | UUID/path/symbol hints |
| G9-I3 | `TestExtractIntentQueryMerge` | Query tokens merged |
| G9-I4 | `TestSearchUsesIntent` | Search enriched by intent |
| G9-I5 | `TestIntentNoSemantic` | No vector/semantic channel |
| G9-I6 | `TestExtractIntentDeterministic` | Deterministic output |
| G9-DOC | §3 doc check | Intent shipped + DR-NOSSEM semantic defer (`RETRIEVAL_AND_CONTEXT.md:46–70`) |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Block 0: any G8-L1–L7 not met; S00-02 not APPROVE; default context auto-loads L2/L3
- Block 1: G9 ship path incomplete without doc supersede; S01-02 not APPROVE; semantic channel introduced
- Block 2: M-001 violation (query-only layer/intent, dump replaces moat)
- Block 3: full-graph dump default; cap inflation; business logic in MCP adapters only
- Block 4: G6/G7 accidentally implemented without board spawn; secondary queue missing from DR-HANDOFF
- Block 5: Phase 42+ successor TBD in VERIFY-NOTES
- Block 6: P41 entity schema changed but `trace/graph.json` not updated
- VERIFY-NOTES missing or evidence dir absent after claimed PASS
- Product code shipped in this row

**Do not fail VERIFY solely for:**

| Topic | Disposition |
|-------|-------------|
| HTTP `max_layer` mirror absent | Residual — CLI+MCP sufficient (S00-02 low) |
| G6/G7 not implemented | Expected — secondary queue |
| G-004a vector | Permanent defer (DR-NOSSEM) |
| `IntentSummary` JSON-only (not Markdown render) | Residual — S01-02 low |
| Search multi-OR vs `FTSQuery()` doc path | Residual — behavior OK, doc drift (S01-02 low) |
| Trim comment vs layer-only sort nit | Residual — S00-02 nit |
| `TaskContext` godoc still "L0–L1" | Residual — S00-02 nit |

## Locked verify command floor

Run from Trace repo root. Tee outputs into evidence dir.

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p41-s02-01-verify/evidence"
mkdir -p "$EVID"

{
  echo "verify_id=P41-S02-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=S00-02+S01-02 APPROVE high confidence"
  echo "phase=P41 G8+G9 delivery"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists with metadata.

---

### Block 0 — G8 progressive layers (G8-L1–L7 + S00-02 APPROVE)

```bash
cd /home/ali/Desktop/Trace

# Board sign-off (row 697)
grep -n 'P41-S00-02.*APPROVE\|697.*APPROVE' docs/TODO/phase-41.md \
  | tee "$EVID/00-board-s00-approve.txt"

# G8 acceptance tests
go test ./internal/compiler/... -count=1 -v \
  -run 'TestContextDefaultLayer1|TestContextMaxLayer2|TestContextMaxLayer3|TestContextLayerBudgetCap|TestContextLayerTrimPriority|TestContextDepthIndependentOfLayer|TestContextNoDump|TestNoDumpAPI' \
  2>&1 | tee "$EVID/00-g8-acceptance.txt"

# MCP max_layer mirror (G8-L2)
go test ./internal/mcp/... -count=1 -v -run 'TestMCPContextMaxLayer2' \
  2>&1 | tee "$EVID/00-g8-mcp-maxlayer.txt"

# Library-first — MaxLayer default + admission
grep -n 'MaxLayer\|classifyNaturalLayer\|enrichLayerCandidates' \
  internal/compiler/compiler.go internal/compiler/layer_enrich.go \
  | head -30 | tee "$EVID/00-g8-library-spot-read.txt"

# Reason codes documented
grep -n 'graph_neighbor\|recent_event\|historical_vcs' \
  internal/retrieval/doc.go internal/retrieval/layer_enrich.go \
  | tee "$EVID/00-g8-reason-codes.txt"

# Trim priority L0→L3
grep -n 'layer.*trim\|L0\|L1\|L2\|L3' internal/compiler/budget.go \
  | tee "$EVID/00-g8-trim-priority.txt"

# Thin adapters — CLI + MCP pass-through only
grep -n 'max.layer\|MaxLayer\|max_layer' \
  cmd/trace/context.go internal/mcp/tools_context.go \
  | tee "$EVID/00-g8-adapter-wiring.txt"
wc -l cmd/trace/context.go internal/mcp/tools_context.go \
  | tee "$EVID/00-g8-adapter-size.txt"

# Default caps unchanged (4096/32/64)
grep -n 'DefaultTokenBudget\|DefaultMaxItems\|MaxCandidateHits' internal/compiler/packet.go \
  | tee "$EVID/00-g8-packet-caps.txt"
```

**Block 0 evidence table (fill in VERIFY-NOTES):**

| Check | Expected | Evidence file |
|-------|----------|---------------|
| S00-02 APPROVE | Board row 697 `done` + APPROVE high | `00-board-s00-approve.txt` |
| G8-L1 | Default layer≤1 | `00-g8-acceptance.txt` |
| G8-L2 | L2 items + reason_codes; MCP mirror | `00-g8-acceptance.txt`, `00-g8-mcp-maxlayer.txt` |
| G8-L3 | L3 when graph supports | `00-g8-acceptance.txt` |
| G8-L4 | Budget caps on L2/L3 | `00-g8-acceptance.txt` |
| G8-L5 | Trim L0→L3 | `00-g8-acceptance.txt`, `00-g8-trim-priority.txt` |
| G8-L6 | depth ≠ layer | `00-g8-acceptance.txt` |
| G8-L7 | No dump | `00-g8-acceptance.txt` |
| Law 19 | Logic in compiler; thin adapters | `00-g8-library-spot-read.txt`, `00-g8-adapter-wiring.txt` |

---

### Block 1 — G9 intent pipeline (G9-I1–I6 + §3 doc + S01-02 APPROVE)

```bash
cd /home/ali/Desktop/Trace

# G9 acceptance tests
go test ./internal/retrieval/... -count=1 -v \
  -run 'TestExtractIntentFromTask|TestExtractIntentEntityHints|TestExtractIntentQueryMerge|TestSearchUsesIntent|TestIntentNoSemantic|TestExtractIntentDeterministic' \
  2>&1 | tee "$EVID/01-g9-acceptance.txt"

# Library spot-read — intent.go caps + Search wiring
grep -n 'ExtractIntent\|IntentInput\|FTSQuery\|maxKeywords\|maxEntityHints' \
  internal/retrieval/intent.go \
  | tee "$EVID/01-g9-library-spot-read.txt"
grep -n 'Intent\|ExtractIntent\|intentSearchTerms' internal/retrieval/search.go \
  | tee "$EVID/01-g9-search-wiring.txt"

# Compiler + explore wiring (task moat path)
grep -n 'IntentInput\|ExtractIntent\|IntentSummary' \
  internal/compiler/compiler.go internal/compiler/explore.go internal/compiler/packet.go \
  | tee "$EVID/01-g9-compiler-wiring.txt"

# §3 doc — intent shipped + DR-NOSSEM
grep -n 'Shipped\|DR-NOSSEM\|semantic\|intent extraction' \
  docs/RETRIEVAL_AND_CONTEXT.md \
  | head -20 | tee "$EVID/01-g9-doc-section3.txt"

# No semantic/embedding imports
grep -rn 'embedding\|vector\|semantic_match' internal/retrieval/ internal/compiler/ \
  --include='*.go' \
  | tee "$EVID/01-g9-no-semantic-imports.txt" || echo "OK: no semantic imports" | tee "$EVID/01-g9-no-semantic-imports.txt"

# G1 regression (boundary preserved)
go test ./internal/compiler/... -count=1 -v -run 'TestG1' \
  2>&1 | tee "$EVID/01-g9-g1-regression.txt"

# Board sign-off (row 700)
grep -n 'P41-S01-02.*APPROVE\|700.*APPROVE' docs/TODO/phase-41.md \
  | tee "$EVID/01-board-s01-approve.txt"
```

**Block 1 evidence table:**

| ID | Expected | Evidence |
|----|----------|----------|
| G9-I1–I6 | All PASS | `01-g9-acceptance.txt` |
| G9-DOC | §3 intent shipped + DR-NOSSEM | `01-g9-doc-section3.txt` |
| Library | `ExtractIntent` in retrieval/ | `01-g9-library-spot-read.txt` |
| Compiler wire | title+query IntentInput | `01-g9-compiler-wiring.txt` |
| G1 boundary | `TestG1*` green | `01-g9-g1-regression.txt` |
| S01-02 | APPROVE on board | `01-board-s01-approve.txt` |

---

### Block 2 — M-001 moat preserved

```bash
cd /home/ali/Desktop/Trace

# Default context unchanged (G8-L1 moat)
go test ./internal/compiler/... -count=1 -v -run 'TestContextDefaultLayer1|TestTaskContextAndBudgets' \
  2>&1 | tee "$EVID/02-moat-default-layer.txt"

# Task UUID required on compile/explore paths
grep -n 'task_id\|TaskID\|required' \
  internal/compiler/compiler.go internal/compiler/explore.go internal/mcp/tools_context.go \
  | tee "$EVID/02-moat-task-required.txt"

# Moat lead in Instructions — context before optional explore
grep -n 'trace_tasks\|trace_context\|trace_loop\|trace_review\|trace_explore\|Optional convenience' \
  internal/mcp/instructions.go \
  | tee "$EVID/02-moat-instructions.txt"

# Intent requires task on compile path (not query-only bypass)
grep -n 'taskIntentInput\|IntentInput{.*TaskTitle' \
  internal/compiler/compiler.go internal/compiler/explore.go \
  | tee "$EVID/02-moat-intent-task-bound.txt"

# Layer/intent merge into packet — not standalone dump
grep -n 'IntentSummary\|Items\|Layer' internal/compiler/packet.go \
  | head -15 | tee "$EVID/02-moat-packet-merge.txt"

# DR-HANDOFF M-001 forward note
grep -n 'M-001\|moat\|never replace\|query-only' \
  docs/phases/phase-41-layers-intent/DR-HANDOFF.md \
  | tee "$EVID/02-moat-dr-handoff.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| Default layer≤1 | G8-L1 green — moat path unchanged |
| task_id required | compile/explore reject empty task |
| Intent task-bound | `IntentInput` from task fields on compile path |
| Layer merge | L2/L3 items in packet.Items — not dump API |
| Moat lead | trace_tasks/context before trace_explore |

---

### Block 3 — Laws 6–7 caps honest; Law 19 library-first

```bash
cd /home/ali/Desktop/Trace

# Packet defaults unchanged
grep -n 'DefaultTokenBudget\|DefaultMaxItems\|MaxCandidateHits' internal/compiler/packet.go \
  | tee "$EVID/03-law67-packet-caps.txt"

# Search limits unchanged (32 default, 64 hard)
grep -n 'DefaultSearchLimit\|MaxSearchLimit\|Limit' internal/retrieval/search.go internal/retrieval/types.go \
  | tee "$EVID/03-law67-search-caps.txt"

# Intent output caps (32/16/512/256 per S01-02)
grep -n 'maxKeywords\|maxEntityHints\|maxFTSQueryLen\|SummaryKeywords' \
  internal/retrieval/intent.go internal/compiler/packet.go \
  | tee "$EVID/03-law67-intent-caps.txt"

# No dump API
go test ./internal/compiler/... -count=1 -run 'TestNoDumpAPI|TestContextNoDump' -v \
  2>&1 | tee "$EVID/03-law67-no-dump.txt"

# Law 19 — admission in compiler; intent in retrieval; thin adapters
grep -n 'MaxLayer\|ContextOptions' internal/compiler/compiler.go \
  | head -10 | tee "$EVID/03-law19-compiler.txt"
wc -l internal/retrieval/intent.go cmd/trace/context.go internal/mcp/tools_context.go \
  | tee "$EVID/03-law19-adapter-size.txt"

# Reject — no intent/layer logic duplicated in web/
grep -rn 'ExtractIntent\|MaxLayer\|max_layer' web/ --include='*.ts' --include='*.tsx' 2>/dev/null \
  | tee "$EVID/03-law19-no-web-logic.txt" || echo "OK: no web retrieval logic" | tee "$EVID/03-law19-no-web-logic.txt"
```

**Pass asserts:**

| Law | Expected |
|-----|----------|
| 6–7 | Packet 4096/32/64; search 32/64; intent caps bounded; truncation honest |
| 19 | G8 in `internal/compiler/`; G9 in `internal/retrieval/`; CLI/MCP adapters thin |

---

### Block 4 — Secondary queue (G6/G7 not P41 implement)

```bash
cd /home/ali/Desktop/Trace

grep -A20 'Secondary queue' docs/phases/phase-41-layers-intent/DR-HANDOFF.md \
  | tee "$EVID/04-secondary-queue.txt"

# Scope folders — S00–S02 only (no S03/S04 G6/G7 rows)
ls docs/phases/phase-41-layers-intent/scopes/ \
  | tee "$EVID/04-scope-folders.txt"

# Board — no G6/G7 implement rows
grep -n 'G6\|G7\|S03\|S04' docs/TODO/phase-41.md \
  | tee "$EVID/04-board-no-g6g7-rows.txt"
```

**Pass:** DR-HANDOFF lists G6/G7 as Phase 42+ forward-only; scope folders = `scope-00-progressive-layers`, `scope-01-intent-pipeline`, `scope-02-verify` only.

---

### Block 5 — Successor Phase 42+ (notes only — do not scaffold)

In VERIFY-NOTES, include **successor recommendation table** for S02-02:

| Field | Locked value (from DR-HANDOFF) |
|-------|--------------------------------|
| Successor | **Phase 42+ — G6/G7 secondary queue** (human promotes P42-00) |
| Entry themes | **G6** Non-semantic concept retrieval + **G7** Index freshness & langs |
| P41 outcome | **G8+G9 delivered** — opt-in L2/L3 + rule-based intent |
| Idle alternative | `no successor` — if human defers Phase 42 |

```bash
grep -n 'Phase 42\|G6\|G7\|successor\|Secondary queue' \
  docs/phases/phase-41-layers-intent/DR-HANDOFF.md \
  docs/phases/phase-41-layers-intent/INTAKE.md \
  | tee "$EVID/05-successor-grep.txt"

grep -n 'G6\|G7\|G-004b\|G-005' \
  docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md \
  | head -15 | tee "$EVID/05-successor-remediation-plan.txt"
```

**Pass:** VERIFY-NOTES Block 5 cites **Phase 42+ G6/G7** by name; **never TBD**.

---

### Block 6 — Graph export (if applicable)

```bash
cd /home/ali/Desktop/Trace

# P41 product window — entity schema changes?
git log --oneline --since=2026-08-22 -- \
  internal/store/ internal/graph/ trace/graph.json pkg/domain/ \
  2>&1 | tee "$EVID/06-entity-commits.txt"

if grep -q . "$EVID/06-entity-commits.txt" 2>/dev/null; then
  go build -o /tmp/trace ./cmd/trace
  /tmp/trace seed export -o trace/graph.json 2>&1 | tee "$EVID/06-graph-export.txt"
  git diff --stat trace/graph.json | tee "$EVID/06-graph-diff-stat.txt"
else
  echo "N/A — no entity schema commits in P41 window (expected — G8/G9 retrieval/compiler only)" \
    | tee "$EVID/06-graph-export-na.txt"
fi
```

**Pass:** N/A with note if no entity changes (expected per S00-01/S01-01 Notes); or `trace/graph.json` updated if schema changed.

---

### Full test floor (aggregate)

```bash
cd /home/ali/Desktop/Trace
go test ./internal/compiler/... ./internal/retrieval/... ./internal/mcp/... ./cmd/trace/... -count=1 \
  2>&1 | tee "$EVID/go-test-p41-full.txt"
```

---

### WRITE VERIFY-NOTES.md

Create `docs/phases/phase-41-layers-intent/scopes/scope-02-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — Phase 41 / S02-01

**Date:** …
**Overall:** PASS | FAIL | PARTIAL
**Git SHA:** …
**Evidence:** experiments/runs/…-p41-s02-01-verify/evidence/

## Precondition cites
- P41-S00-02 APPROVE (high) — G8 progressive layers L2–L3
- P41-S01-02 APPROVE (high) — G9 rule-based intent pipeline

## Block results
| Block | Check | Result | Evidence file |
|-------|-------|--------|---------------|
| 0 | G8 G8-L1–L7 + S00 APPROVE | | 00-g8-*.txt |
| 1 | G9 G9-I1–I6 + §3 + S01 APPROVE | | 01-g9-*.txt |
| 2 | M-001 moat | | 02-moat-*.txt |
| 3 | Laws 6–7 / 19 | | 03-law*.txt |
| 4 | G6/G7 forward-only | | 04-*.txt |
| 5 | Phase 42+ successor prep | | 05-successor-*.txt |
| 6 | Graph export | PASS/N/A | 06-graph-*.txt |

## G8 accept map
| ID | Result | Evidence |
| G8-L1 | | |
| … | | |

## G9 accept map
| ID | Result | Evidence |
| G9-I1 | | |
| … | | |

## Successor recommendation (for S02-02)
- **Default:** Phase 42+ — G6 non-semantic concept retrieval + G7 index freshness
- **Secondary:** per REMEDIATION-PLAN rank
- **Never:** TBD

## Residuals (non-blocking)
- HTTP max_layer absent (S00-02)
- IntentSummary JSON-only (S01-02)
- Search multi-OR vs FTSQuery doc (S01-02)

## DR-HANDOFF
Stays OPEN — P41-S02-02 closes + scaffolds Phase 42+

## Next
P41-S02-02
```

## Exit criteria

- [ ] `VERIFY-NOTES.md` with block table 0–6 + G8/G9 accept maps
- [ ] Evidence dir populated under `experiments/runs/…-p41-s02-01-verify/evidence/`
- [ ] Blocks 0–6 executed in order
- [ ] Board Notes on **P41-S02-01** only
- [ ] DR-HANDOFF remains OPEN
- [ ] Next: **P41-S02-02**

## Next

`P41-S02-02`
