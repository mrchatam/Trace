# P42-S02-01 — VERIFY (G6 + G7)

## Metadata
- id: P42-S02-01
- todo_ids: [P42-S02-01]
- role: verifier
- skills: [qa-lead, production-validator, test-driven-development]
- mcps: [Shell, user-trace]
- verification: mixed (test floor + review sign-offs + artifact checklist)

## Objective

Run locked verify blocks **0–6** from [00-PLANNER.md](00-PLANNER.md). Author **`VERIFY-NOTES.md`** with PASS/FAIL per block and evidence manifest. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P42-S02-02**. **No new features.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor (FINAL — S02-00)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S02-02
- [INTAKE.md](../../INTAKE.md)
- [REMEDIATION-PLAN §2 G6/G7](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- S00–S01 review APPROVE notes on [docs/TODO/phase-42.md](../../../../TODO/phase-42.md)
- Pattern: [P41 S02-01 verify](../../../phase-41-layers-intent/scopes/scope-02-verify/01-verify.md)
- Code anchors (live post-S01 — from S00-02 / S01-02 reviews):
  - G6: `internal/retrieval/concept.go:8–15` (concept types); `types.go:14` (`ReasonGraphLabelMatch`); `concept.go:22–38` (SearchGraphLabels cap≤64)
  - G6 merge: `internal/compiler/compiler.go:183–187`; `internal/compiler/explore.go:126–130` (fail-open DF-87)
  - G6 law: `scopes/scope-00-non-semantic-concept/LAW-REVIEW-NOTES.md` PASS
  - G6 doc: `docs/RETRIEVAL_AND_CONTEXT.md:34–36`, `:48–59` (G6 shipped + DR-NOSSEM defer)
  - G6 tests: `internal/retrieval/concept_test.go` G6-C1–C7; `internal/compiler/compiler_test.go:1662` (`TestContextIncludesGraphLabels`)
  - G7: `docs/INDEX_LANG_POLICY.md`; `internal/analyzers/language_adapter.go:55–` (`SupportedLanguages()`)
  - G7 CLI: `cmd/trace/index_status.go:20`; `cmd/trace/index_watch.go` (foreground watch)
  - G7 HTTP: `internal/httpapi/handlers_p1.go:356` (`supported_languages`)
  - G7 hook: `internal/install/githook.go` (primary freshness)
  - G7 tests: `cmd/trace/index_status_test.go` G7-F1; `language_adapter_test.go` G7-F2–F3; `index_watch_test.go` G7-F4–F5; `httpapi_test.go:710` G7-F6

## Session start

Follow agent-loop-protocol Session start. Unattended: run blocks in order; do **not** close DR-HANDOFF or scaffold Phase 43.

## Locked defaults (FINAL — S02-00)

| Item | Value |
|------|-------|
| Precondition | P42-S00-02, P42-S01-02 all `done` + **APPROVE** (high confidence) |
| Product Go / TS / `web/` | **Forbidden** (evidence + notes only) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p42-s02-01-verify/evidence/` |
| Notes artifact | `scopes/scope-02-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S02-02 closes |
| Phase 43 scaffold | **Out of scope** — S02-02 owns per agent-loop-protocol Phase handoff |
| P42 git window | From P42-S00-01 first product commit through verify run |
| Trace repo root | `/home/ali/Desktop/Trace` |

### Review sign-off precondition (Block 0–1 gate)

| Row | Verdict | Confidence | Notes cite |
|-----|---------|------------|------------|
| P42-S00-02 | **APPROVE** | high | G6-C1–C7; LAW-REVIEW PASS; board row **707** |
| P42-S01-02 | **APPROVE** | high | G7-F1–F6; board row **710** |

### G6 accept map (Block 0 — must tick in VERIFY-NOTES)

| ID | Criterion | Primary evidence |
|----|-----------|------------------|
| G6-C1 | Concept body match → `graph_label_match` | `TestSearchGraphLabelsDiscovery` (`concept_test.go:24`) |
| G6-C2 | Entity filter excludes file/symbol | `TestSearchGraphLabelsEntityFilter` (`:65`) |
| G6-C3 | Compile merge | `TestContextIncludesGraphLabels` (`compiler_test.go:1662`) |
| G6-C4 | Cap ≤64 | `TestSearchGraphLabelsCap` (`concept_test.go:101`) |
| G6-C5 | No semantic/vector | `TestSearchGraphLabelsNoSemantic` (`:131`) |
| G6-C6 | Deterministic | `TestSearchGraphLabelsDeterministic` (`:172`) |
| G6-C7 | Fail-open | `TestSearchGraphLabelsFailOpen` (`:225`) |
| LAW | DR-NOSSEM desk-check | `LAW-REVIEW-NOTES.md` PASS |

### G7 accept map (Block 1 — must tick in VERIFY-NOTES)

| ID | Criterion | Primary evidence |
|----|-----------|------------------|
| G7-F1 | Status JSON langs | `TestIndexStatusSupportedLanguages` (`index_status_test.go:11`) |
| G7-F2 | Adapter parity | `TestSupportedLanguagesMatchesAdapters` (`language_adapter_test.go:11`) |
| G7-F3 | Unsupported ext UX | `TestIndexUnsupportedExtMessage` (`:28`) |
| G7-F4 | Watch debounce | `TestIndexWatchDebounced` (`index_watch_test.go:13`) |
| G7-F5 | Foreground exit | `TestIndexWatchForegroundExit` (`:56`) |
| G7-F6 | HTTP mirror | `TestHTTPIndexStatusLanguages` (`httpapi_test.go:710`) or N/A residual |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Block 0: any G6-C1–C7 not met; S00-02 not APPROVE; LAW-REVIEW not PASS; vector channel introduced
- Block 1: G7 policy/watch incomplete; S01-02 not APPROVE; always-on daemon introduced
- Block 2: M-001 violation (query-only concept path, dump replaces moat)
- Block 3: full-graph dump; cap inflation; business logic in adapters only
- Block 4: G-004a vector shipped; semantic_match in concept path
- Block 5: Phase 43+ successor TBD in VERIFY-NOTES
- Block 6: P42 entity schema changed but `trace/graph.json` not updated
- VERIFY-NOTES missing or evidence dir absent after claimed PASS
- Product code shipped in this row

**Do not fail VERIFY solely for:**

| Topic | Disposition |
|-------|-------------|
| HTTP G7-F6 absent | Residual — CLI sufficient |
| Tier-2 langs not shipped | Expected — policy defer |
| G-004a vector | Permanent defer — DR-NOSSEM |
| Phase 43+ not scaffolded | S02-02 owns — expected at S02-01 |
| explore graph-label merge gated on non-empty searchQ | Residual — S00-02 low |
| watch indexOne HEAD-first in git repos | Residual — S01-02 low |

## Locked verify command floor

Run from Trace repo root. Tee outputs into evidence dir.

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p42-s02-01-verify/evidence"
mkdir -p "$EVID"

{
  echo "verify_id=P42-S02-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=S00-02+S01-02 APPROVE high confidence"
  echo "phase=P42 G6+G7 delivery"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists with metadata.

---

### Block 0 — G6 graph-label concept retrieval (G6-C1–C7 + LAW-REVIEW + S00-02 APPROVE)

```bash
cd /home/ali/Desktop/Trace

# Board sign-off (row 707)
grep -n 'P42-S00-02.*APPROVE\|707.*APPROVE' docs/TODO/phase-42.md \
  | tee "$EVID/00-board-s00-approve.txt"

# LAW-REVIEW desk-check PASS
grep -n 'Verdict.*PASS\|DR-NOSSEM' \
  docs/phases/phase-42-concept-index/scopes/scope-00-non-semantic-concept/LAW-REVIEW-NOTES.md \
  | tee "$EVID/00-law-review-pass.txt"

# G6 acceptance tests
go test ./internal/retrieval/... -count=1 -v \
  -run 'TestSearchGraphLabelsDiscovery|TestSearchGraphLabelsEntityFilter|TestSearchGraphLabelsCap|TestSearchGraphLabelsNoSemantic|TestSearchGraphLabelsDeterministic|TestSearchGraphLabelsFailOpen|TestMergeConceptHitsDedupe' \
  2>&1 | tee "$EVID/00-g6-acceptance.txt"

# Compile merge (G6-C3)
go test ./internal/compiler/... -count=1 -v \
  -run 'TestContextIncludesGraphLabels|TestNoDumpAPI' \
  2>&1 | tee "$EVID/00-g6-compile-merge.txt"

# Library spot-read — concept types + reason code + cap
grep -n 'conceptEntityTypes\|ReasonGraphLabelMatch\|SearchGraphLabels\|MergeConceptHits\|limit > 64' \
  internal/retrieval/concept.go internal/retrieval/types.go \
  | tee "$EVID/00-g6-library-spot-read.txt"

# Compiler + explore wiring (fail-open DF-87)
grep -n 'SearchGraphLabels\|MergeConceptHits\|labels = nil' \
  internal/compiler/compiler.go internal/compiler/explore.go \
  | tee "$EVID/00-g6-compiler-wiring.txt"

# §2 doc — G6 shipped + DR-NOSSEM semantic defer
grep -n 'graph_label_match\|G6\|DR-NOSSEM\|Phase 42' \
  docs/RETRIEVAL_AND_CONTEXT.md \
  | head -25 | tee "$EVID/00-g6-doc-section2.txt"

# No vector/semantic in concept path
grep -rn 'semantic_match\|embedding\|vector' internal/retrieval/concept.go \
  | tee "$EVID/00-g6-no-vector.txt" || echo "OK: no vector in concept.go" | tee "$EVID/00-g6-no-vector.txt"
```

**Block 0 evidence table (fill in VERIFY-NOTES):**

| Check | Expected | Evidence file |
|-------|----------|---------------|
| S00-02 APPROVE | Board row 707 `done` + APPROVE high | `00-board-s00-approve.txt` |
| LAW-REVIEW | PASS + DR-NOSSEM | `00-law-review-pass.txt` |
| G6-C1 | Discovery → graph_label_match | `00-g6-acceptance.txt` |
| G6-C2 | Entity filter | `00-g6-acceptance.txt` |
| G6-C3 | Compile merge | `00-g6-compile-merge.txt` |
| G6-C4 | Cap ≤64 | `00-g6-acceptance.txt`, `00-g6-library-spot-read.txt` |
| G6-C5 | No semantic | `00-g6-acceptance.txt`, `00-g6-no-vector.txt` |
| G6-C6 | Deterministic | `00-g6-acceptance.txt` |
| G6-C7 | Fail-open | `00-g6-acceptance.txt` |
| Law 19 | Logic in retrieval/; compile merge only | `00-g6-library-spot-read.txt`, `00-g6-compiler-wiring.txt` |

---

### Block 1 — G7 index freshness & langs (G7-F1–F6 + INDEX_LANG_POLICY + S01-02 APPROVE)

```bash
cd /home/ali/Desktop/Trace

# Board sign-off (row 710)
grep -n 'P42-S01-02.*APPROVE\|710.*APPROVE' docs/TODO/phase-42.md \
  | tee "$EVID/01-board-s01-approve.txt"

# Policy doc exists + tier table
test -f docs/INDEX_LANG_POLICY.md
grep -n 'Tier-1\|Tier-2\|git-hook\|watch' docs/INDEX_LANG_POLICY.md \
  | head -20 | tee "$EVID/01-g7-policy-doc.txt"

# G7 acceptance tests
CGO_ENABLED=1 go test ./internal/analyzers/... ./cmd/trace/... ./internal/httpapi/... -count=1 -v \
  -run 'TestIndexStatusSupportedLanguages|TestSupportedLanguagesMatchesAdapters|TestIndexUnsupportedExtMessage|TestIndexWatchDebounced|TestIndexWatchForegroundExit|TestHTTPIndexStatusLanguages' \
  2>&1 | tee "$EVID/01-g7-acceptance.txt"

# SupportedLanguages + status JSON
grep -n 'SupportedLanguages\|supported_languages' \
  internal/analyzers/language_adapter.go cmd/trace/index_status.go internal/httpapi/handlers_p1.go \
  | tee "$EVID/01-g7-lang-wiring.txt"

# 5 tier-1 langs frozen
grep -n 'go\|javascript\|typescript\|tsx\|python' \
  internal/analyzers/detect.go internal/analyzers/language_adapter.go \
  | tee "$EVID/01-g7-tier1-langs.txt"

# Git-hook primary (no always-on daemon)
grep -n 'git.hook\|GitHook\|install git-hook' \
  internal/install/githook.go cmd/trace/install.go docs/INDEX_LANG_POLICY.md \
  | tee "$EVID/01-g7-githook-primary.txt"

# Foreground watch — fsnotify, debounce, SIGINT exit
grep -n 'fsnotify\|debounce\|SIGINT\|daemon' \
  cmd/trace/index_watch.go docs/INDEX_LANG_POLICY.md \
  | tee "$EVID/01-g7-watch-foreground.txt"

# Unsupported ext cites policy
grep -n 'INDEX_LANG_POLICY\|unsupported' \
  internal/analyzers/errors.go \
  | tee "$EVID/01-g7-unsupported-ext.txt"
```

**Block 1 evidence table:**

| ID | Expected | Evidence |
|----|----------|----------|
| G7-F1–F6 | All PASS (F6 or N/A residual) | `01-g7-acceptance.txt` |
| INDEX_LANG_POLICY | Tier table + hook primary + watch optional | `01-g7-policy-doc.txt` |
| SupportedLanguages | Exported; sorted; 5 langs | `01-g7-lang-wiring.txt`, `01-g7-tier1-langs.txt` |
| Git-hook primary | No daemon default | `01-g7-githook-primary.txt`, `01-g7-watch-foreground.txt` |
| S01-02 | APPROVE on board | `01-board-s01-approve.txt` |

---

### Block 2 — M-001 moat preserved

```bash
cd /home/ali/Desktop/Trace

# Task context moat — compile requires task; concept merges into packet
go test ./internal/compiler/... -count=1 -v \
  -run 'TestTaskContextAndBudgets|TestContextDefaultLayer1|TestContextIncludesGraphLabels' \
  2>&1 | tee "$EVID/02-moat-task-context.txt"

# Task UUID required on compile/explore paths
grep -n 'task_id\|TaskID\|required' \
  internal/compiler/compiler.go internal/compiler/explore.go internal/mcp/tools_context.go \
  | tee "$EVID/02-moat-task-required.txt"

# Moat lead in Instructions — tasks before explore
grep -n 'trace_tasks\|trace_context\|trace_loop\|trace_review\|trace_explore\|Optional convenience' \
  internal/mcp/instructions.go \
  | tee "$EVID/02-moat-instructions.txt"

# Concept merge into candidates — not query-only bypass
grep -n 'SearchGraphLabels\|MergeConceptHits\|taskIntent' \
  internal/compiler/compiler.go internal/compiler/explore.go \
  | tee "$EVID/02-moat-concept-merge.txt"

# DR-HANDOFF M-001 forward note
grep -n 'M-001\|moat\|never replace\|query-only' \
  docs/phases/phase-42-concept-index/DR-HANDOFF.md docs/phases/phase-42-concept-index/INTAKE.md \
  | tee "$EVID/02-moat-dr-handoff.txt"

# Index honesty (G7) — no silent stale claims
go test ./internal/analyzers/... ./cmd/trace/... -count=1 -v \
  -run 'TestIndexHonesty|TestGraphSyncHonesty' 2>/dev/null \
  | tee "$EVID/02-moat-index-honesty.txt" || echo "N/A if test names differ" | tee "$EVID/02-moat-index-honesty.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| Task loop primary | compile/explore require task_id |
| Concept merge | G6 hits merge into packet candidates — not standalone dump |
| Index honesty | status JSON honest; unsupported ext errors cite policy |
| Moat lead | trace_tasks/context before trace_explore |

---

### Block 3 — Laws 6–7 caps honest; Law 19 library-first

```bash
cd /home/ali/Desktop/Trace

# Packet defaults unchanged (4096/32/64)
grep -n 'DefaultTokenBudget\|DefaultMaxItems\|MaxCandidateHits' internal/compiler/packet.go \
  | tee "$EVID/03-law67-packet-caps.txt"

# Search + concept limits (32 default, 64 hard)
grep -n 'DefaultSearchLimit\|MaxSearchLimit\|limit > 64\|limit = 16' \
  internal/retrieval/search.go internal/retrieval/concept.go internal/retrieval/types.go \
  | tee "$EVID/03-law67-search-caps.txt"

# No dump API
go test ./internal/retrieval/... ./internal/compiler/... -count=1 -v \
  -run 'TestNoDump|TestNoDumpAPI|TestContextNoDump' \
  2>&1 | tee "$EVID/03-law67-no-dump.txt"

# Law 19 — G6 in retrieval/; G7 in analyzers/cmd; HTTP thin mirror
grep -n 'SearchGraphLabels\|SupportedLanguages' \
  internal/retrieval/concept.go internal/analyzers/language_adapter.go \
  | tee "$EVID/03-law19-library.txt"
wc -l cmd/trace/index_status.go cmd/trace/index_watch.go internal/httpapi/handlers_p1.go \
  | tee "$EVID/03-law19-adapter-size.txt"

# Reject — no concept/index logic duplicated in web/
grep -rn 'SearchGraphLabels\|SupportedLanguages\|index watch' web/ --include='*.ts' --include='*.tsx' 2>/dev/null \
  | tee "$EVID/03-law19-no-web-logic.txt" || echo "OK: no web index/retrieval logic" | tee "$EVID/03-law19-no-web-logic.txt"
```

**Pass asserts:**

| Law | Expected |
|-----|----------|
| 6–7 | Packet 4096/32/64; search/concept 32/64; truncation honest |
| 19 | G6 in `internal/retrieval/`; G7 in analyzers/cmd; HTTP handlers thin mirror |

---

### Block 4 — G-004a vector absent; DR-NOSSEM honored

```bash
cd /home/ali/Desktop/Trace

# DR-NOSSEM — no semantic_match reason code shipped
grep -rn 'semantic_match' internal/retrieval/ internal/compiler/ --include='*.go' \
  | tee "$EVID/04-dr-nossem-semantic-match.txt" || echo "OK: no semantic_match" | tee "$EVID/04-dr-nossem-semantic-match.txt"

# No embedding/vector imports in retrieval
grep -rn 'embedding\|vector\|openai\|cohere' internal/retrieval/ --include='*.go' \
  | tee "$EVID/04-dr-nossem-no-vector.txt" || echo "OK: no vector imports" | tee "$EVID/04-dr-nossem-no-vector.txt"

# INTAKE + DR-HANDOFF reject preserved
grep -n 'G-004a\|vector\|DR-NOSSEM\|Rejects preserved' \
  docs/phases/phase-42-concept-index/INTAKE.md docs/phases/phase-42-concept-index/DR-HANDOFF.md \
  | tee "$EVID/04-dr-nossem-intake.txt"

# REMEDIATION-PLAN — G6 non-semantic only
grep -n 'G-004a\|G-004b\|vector\|graph-label' \
  docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md \
  | head -15 | tee "$EVID/04-dr-nossem-remediation-plan.txt"

# Scope folders — S00–S02 only (no accidental vector scope)
ls docs/phases/phase-42-concept-index/scopes/ \
  | tee "$EVID/04-scope-folders.txt"
```

**Pass:** No `semantic_match` in shipped code; G-004a permanently deferred; scope folders = `scope-00-non-semantic-concept`, `scope-01-index-freshness-langs`, `scope-02-verify` only.

---

### Block 5 — Successor **`no successor`** (notes only — do not scaffold Phase 43)

In VERIFY-NOTES, include **successor recommendation table** for S02-02:

| Field | Locked value (from DR-HANDOFF) |
|-------|--------------------------------|
| Successor | **`no successor`** — G1–G9 remediation complete after P42 |
| Optional | Phase 43+ residuals (HTTP index, Tier-2 lang) — human promotion only |
| P42 outcome | **G6+G7 delivered** — graph-label channel + index/lang policy |
| Idle alternative | Same as default when G6+G7 ship |

```bash
grep -n 'no successor\|Phase 43\|G1–G9\|remediation complete' \
  docs/phases/phase-42-concept-index/DR-HANDOFF.md \
  docs/phases/phase-42-concept-index/INTAKE.md \
  | tee "$EVID/05-successor-grep.txt"

grep -n 'G6\|G7\|G-004b\|G-005\|complete' \
  docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md \
  | head -20 | tee "$EVID/05-successor-remediation-plan.txt"
```

**Pass:** VERIFY-NOTES Block 5 cites **`no successor`** by name; **never TBD**.

---

### Block 6 — Graph export (if applicable)

```bash
cd /home/ali/Desktop/Trace

# P42 product window — entity schema changes?
git log --oneline --since=2026-08-22 -- \
  internal/store/ internal/graph/ trace/graph.json pkg/domain/ \
  2>&1 | tee "$EVID/06-entity-commits.txt"

if grep -q . "$EVID/06-entity-commits.txt" 2>/dev/null; then
  go build -o /tmp/trace ./cmd/trace
  /tmp/trace seed export -o trace/graph.json 2>&1 | tee "$EVID/06-graph-export.txt"
  git diff --stat trace/graph.json | tee "$EVID/06-graph-diff-stat.txt"
else
  echo "N/A — no entity schema commits in P42 window (expected — G6 reason_code + G7 policy only)" \
    | tee "$EVID/06-graph-export-na.txt"
fi
```

**Pass:** N/A with note if no entity changes (expected per S00-01/S01-01 Notes); or `trace/graph.json` updated if schema changed.

---

### Full test floor (aggregate)

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go test ./internal/retrieval/... ./internal/compiler/... ./internal/analyzers/... ./cmd/trace/... ./internal/httpapi/... -count=1 \
  2>&1 | tee "$EVID/go-test-p42-full.txt"
```

---

### WRITE VERIFY-NOTES.md

Create `docs/phases/phase-42-concept-index/scopes/scope-02-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — Phase 42 / S02-01

**Date:** …
**Overall:** PASS | FAIL | PARTIAL
**Git SHA:** …
**Evidence:** experiments/runs/…-p42-s02-01-verify/evidence/

## Precondition cites
- P42-S00-02 APPROVE (high) — G6 graph-label concept retrieval
- P42-S01-02 APPROVE (high) — G7 index freshness & langs

## Block results
| Block | Check | Result | Evidence file |
|-------|-------|--------|---------------|
| 0 | G6 G6-C1–C7 + LAW + S00 APPROVE | | 00-g6-*.txt, 00-law-*.txt |
| 1 | G7 G7-F1–F6 + policy + S01 APPROVE | | 01-g7-*.txt |
| 2 | M-001 moat | | 02-moat-*.txt |
| 3 | Laws 6–7 / 19 | | 03-law*.txt |
| 4 | G-004a absent; DR-NOSSEM | | 04-dr-nossem-*.txt |
| 5 | **`no successor`** named | | 05-successor-*.txt |
| 6 | Graph export | PASS/N/A | 06-graph-*.txt |

## G6 accept map
| ID | Result | Evidence |
| G6-C1 | | |
| G6-C2 | | |
| G6-C3 | | |
| G6-C4 | | |
| G6-C5 | | |
| G6-C6 | | |
| G6-C7 | | |
| LAW | | |

## G7 accept map
| ID | Result | Evidence |
| G7-F1 | | |
| G7-F2 | | |
| G7-F3 | | |
| G7-F4 | | |
| G7-F5 | | |
| G7-F6 | | PASS or N/A residual |

## Successor recommendation (for S02-02)
- **Default:** **`no successor`** — G1–G9 remediation complete
- **Optional:** Phase 43+ residuals (HTTP POST /v1/index, first Tier-2 lang) — human promotion only
- **Never:** TBD

## Residuals (non-blocking)
- explore graph-label merge gated on non-empty searchQ (S00-02 low)
- watch indexOne HEAD-first in git repos (S01-02 low)
- fsnotify indirect in go.mod (S01-02 low)
- HTTP G7-F6 absent if CLI sufficient

## DR-HANDOFF
Stays OPEN — P42-S02-02 closes

## Next
P42-S02-02
```

## Role work

1. Run blocks 0–6 in order; record PASS/FAIL per block in VERIFY-NOTES.md.
2. Cite board row Notes for S00-02 / S01-02 APPROVE.
3. Name successor: **`no successor`** (default) unless human deferred with reason.
4. Leave DR-HANDOFF **OPEN**.

## Exit criteria

- [ ] `VERIFY-NOTES.md` with block table 0–6 + G6/G7 accept maps
- [ ] Evidence dir populated under `experiments/runs/…-p42-s02-01-verify/evidence/`
- [ ] Blocks 0–6 executed in order
- [ ] Successor named (**`no successor`** — never TBD)
- [ ] Board Notes on **P42-S02-01** only
- [ ] DR-HANDOFF remains OPEN
- [ ] Next: **P42-S02-02**

## Next

`P42-S02-02`
