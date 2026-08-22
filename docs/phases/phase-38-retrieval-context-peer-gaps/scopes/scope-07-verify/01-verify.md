# P38-S07-01 — VERIFY Phase 38

## Metadata
- id: P38-S07-01
- todo_ids: [P38-S07-01]
- role: implementer
- skills: [qa-lead, research, documentation-and-adrs]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed (artifact checklist + git boundary + spot-read)
- hooks: []

## Objective

Run locked verify blocks **0–6** from [00-PLANNER.md](00-PLANNER.md). Author **`VERIFY-NOTES.md`** with PASS/FAIL per block and artifact manifest. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P38-S07-02**. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S07-00)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — investigation-only law
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S07-02
- S00–S06 artifacts (see Locked artifact manifest)
- [REMEDIATION-PLAN.md](../scope-06-remediation-plan/REMEDIATION-PLAN.md) — §2 G1–G9, §4 rejects, §6 successor
- Pattern: [P24 S05-01 verify](../../../phase-24-agent-effectiveness-investigation/scopes/scope-05-phase-verify/01-verify.md) (investigation archive)
- Pattern: [P37 S03-01 verify](../../../phase-37-p36-residuals/scopes/scope-03-verify/01-verify.md) (block table shape)

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row verifies documentation completeness and archives evidence; it does **not** close DR-HANDOFF, scaffold Phase 39, or change product bodies.

## Locked defaults (FINAL — S07-00)

| Item | Value |
|------|-------|
| Precondition | P38-00 … P38-S06-02 all `done`; S06 review **APPROVE** (high confidence) |
| Product Go / TS / `web/` | **Forbidden** (evidence + notes only) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p38-s07-01-verify/evidence/` |
| Notes artifact | `scopes/scope-07-verify/VERIFY-NOTES.md` (**required**) |
| Archive method | Copy each S00–S06 artifact into evidence dir **and** write `manifest.sha256` |
| DR-HANDOFF | Stays **OPEN** — S07-02 closes + scaffolds Phase 39 |
| Phase 39 scaffold | **Out of scope** — S07-02 owns per agent-loop-protocol Phase handoff |
| P38 git window start | Human promotion **2026-08-22** (P38-00 done) — use `--since=2026-08-22` for Block 0 |

### Locked artifact manifest (must exist — Block 1 spot-read)

| Artifact | Source path (repo root relative) |
|----------|----------------------------------|
| INVESTIGATION-INDEX.md | `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-00-investigation-index/INVESTIGATION-INDEX.md` |
| TRACE-AUDIT.md | `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-01-trace-audit/TRACE-AUDIT.md` |
| PEER-CG.md | `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md` |
| PEER-UA-GF.md | `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md` |
| GAP-REGISTRY.md | `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md` |
| SATURATION-NOTES.md | `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-05-saturation-gate/SATURATION-NOTES.md` |
| REMEDIATION-PLAN.md | `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md` |

### H1–H11 register (Block 1 — must map in VERIFY-NOTES)

| H* | G-ID | Primary artifact |
|----|------|------------------|
| H1 | G-001 | GAP-REGISTRY §2.2 + TRACE-AUDIT |
| H2 | G-002 | GAP-REGISTRY §2.2 + TRACE-AUDIT |
| H3 | G-003 | GAP-REGISTRY §2.2 + TRACE-AUDIT |
| H4 | G-004 (+ G-004a defer, G-004b gap) | GAP-REGISTRY §2.2–§2.3 |
| H5 | G-005 | GAP-REGISTRY §2.2 + PEER-CG |
| H6 | G-006 | GAP-REGISTRY §2.2 + TRACE-AUDIT + PEER-CG + PEER-UA-GF §3 MP |
| H7 | G-007 | GAP-REGISTRY §2.2 + PEER-CG + SATURATION h7-compose-desk-check |
| H8 | G-008 | GAP-REGISTRY §2.2 + PEER-UA-GF |
| H9 | G-009 | GAP-REGISTRY §2.2 + TRACE-AUDIT + PEER-UA-GF §3 MP |
| H10 | G-010 | GAP-REGISTRY §2.2 + TRACE-AUDIT |
| H11 | G-011 | GAP-REGISTRY §2.2 + REMEDIATION-PLAN G4 |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Block 0: any commit in P38 window touches `internal/`, `cmd/`, `web/`, `trace/` (product) outside allowed doc-only paths
- Block 1: any required artifact missing OR any H1–H11 without gap/defer/non-gap verdict in GAP-REGISTRY
- Block 2: SATURATION-NOTES §6 `ready_for_REMEDIATION_PLAN` not `true` OR S05-02 not APPROVE (saturated)
- Block 3: REMEDIATION-PLAN missing G1–G9 rank table OR §4 reject list &lt; 12 OR "implement in P38" language
- Block 4: missing mechanism cites for **CG, UA, GF, or MP** (file:line or `$EV/` paths)
- Block 5: GAP-REGISTRY §3 M-001 moat row absent or marked gap
- Block 6: REMEDIATION-PLAN §6 lacks Phase 39 G1+G3+G4 entry co-wave recommendation
- VERIFY-NOTES missing or evidence dir absent after claimed PASS
- Product code shipped in this row

**Do not fail VERIFY solely for:**

| Topic | Disposition |
|-------|-------------|
| GAP-REGISTRY §6 H7 still "Open" | Forward-only — S05 closed via defer→S06 |
| SATURATION §3 criterion 3 says "CG, UA, GF" without MP | MP covered in GAP-REGISTRY matrix + Block 4 MP gate |
| Phase 39 folder not yet exists | S07-02 scaffold — not S07-01 |
| Optional `$EV/` file named in appendix but claim verified elsewhere | Low — cite alternate evidence |

## Locked verify command floor

Run from Trace repo root. Tee outputs into evidence dir.

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p38-s07-01-verify/evidence"
mkdir -p "$EVID"

{
  echo "verify_id=P38-S07-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P38-S06-02 APPROVE high confidence"
  echo "phase=P38 investigation-only"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists with metadata.

---

### Block 0 — No product Go/TS commits in P38 scope

P38 is **investigate → saturate → plan only**. Product paths must be untouched during the P38 window.

```bash
# Product-path commits since human promotion (expect empty)
git log --oneline --since=2026-08-22 -- internal/ cmd/ trace/ web/ \
  2>&1 | tee "$EVID/00-product-commits-since-promotion.txt"

# Diff stat vs Aug 21 baseline (optional cross-check)
git diff --stat $(git rev-list -1 --before="2026-08-21" HEAD)..HEAD -- internal/ cmd/ web/ trace/ \
  2>&1 | tee "$EVID/00-product-diff-stat.txt" || true

# Allowed paths touched in P38 (expect docs + experiments)
git log --oneline --since=2026-08-22 --name-only -- docs/phases/phase-38* docs/TODO/phase-38.md experiments/runs/ \
  2>&1 | tee "$EVID/00-allowed-p38-paths.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| `00-product-commits-since-promotion.txt` | **Empty** (no lines except header) |
| Allowed paths | Docs under `phase-38*` + `experiments/runs/` present |
| Board history | No P38 implement rows with product code in Notes |

---

### Block 1 — All H1–H11 addressed in GAP-REGISTRY

```bash
P38=docs/phases/phase-38-retrieval-context-peer-gaps

# Artifact existence
for f in \
  scopes/scope-00-investigation-index/INVESTIGATION-INDEX.md \
  scopes/scope-01-trace-audit/TRACE-AUDIT.md \
  scopes/scope-02-codegraph-peer/PEER-CG.md \
  scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md \
  scopes/scope-04-gap-registry/GAP-REGISTRY.md \
  scopes/scope-05-saturation-gate/SATURATION-NOTES.md \
  scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md
do
  test -f "$P38/$f" && echo "OK $f" || echo "MISSING $f"
done | tee "$EVID/01-artifact-existence.txt"

# H* rows in GAP-REGISTRY
grep -E '^\| \*\*G-00[0-9]|^\| \*\*G-01[01]' "$P38/scopes/scope-04-gap-registry/GAP-REGISTRY.md" \
  | tee "$EVID/01-gap-rows.txt"

# Archive + manifest
mkdir -p "$EVID/artifacts"
cp "$P38/scopes/scope-00-investigation-index/INVESTIGATION-INDEX.md" "$EVID/artifacts/"
cp "$P38/scopes/scope-01-trace-audit/TRACE-AUDIT.md" "$EVID/artifacts/"
cp "$P38/scopes/scope-02-codegraph-peer/PEER-CG.md" "$EVID/artifacts/"
cp "$P38/scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md" "$EVID/artifacts/"
cp "$P38/scopes/scope-04-gap-registry/GAP-REGISTRY.md" "$EVID/artifacts/"
cp "$P38/scopes/scope-05-saturation-gate/SATURATION-NOTES.md" "$EVID/artifacts/"
cp "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md" "$EVID/artifacts/"

(
  cd "$P38"
  sha256sum scopes/scope-00-investigation-index/INVESTIGATION-INDEX.md
  sha256sum scopes/scope-01-trace-audit/TRACE-AUDIT.md
  sha256sum scopes/scope-02-codegraph-peer/PEER-CG.md
  sha256sum scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md
  sha256sum scopes/scope-04-gap-registry/GAP-REGISTRY.md
  sha256sum scopes/scope-05-saturation-gate/SATURATION-NOTES.md
  sha256sum scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md
) > "$EVID/manifest.sha256"
```

**Block 1 evidence table (fill in VERIFY-NOTES):**

| Check | Expected | Evidence |
|-------|----------|----------|
| 7 artifacts exist | All OK | `01-artifact-existence.txt` |
| G-001…G-011 rows | 11 gap rows + verdict | `01-gap-rows.txt` + GAP-REGISTRY §2.1 |
| H1–H11 mapping | Each H* → G-ID | GAP-REGISTRY §2.1 table |
| INVESTIGATION-INDEX | §2 register matches | Spot-read §2 H1–H11 |
| Archive | 7 files + manifest | `$EVID/artifacts/`, `manifest.sha256` |

---

### Block 2 — S05 saturation APPROVE on record

```bash
grep -n 'ready_for_REMEDIATION_PLAN\|PROCEED_TO_S06\|APPROVE\|6/6' \
  "$P38/scopes/scope-05-saturation-gate/SATURATION-NOTES.md" \
  | tee "$EVID/02-saturation-gate.txt"

grep -n 'Verdict.*APPROVE\|saturated' \
  docs/TODO/phase-38.md | grep -i '664\|S05-02' \
  | tee "$EVID/02-board-s05-approve.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| SATURATION-NOTES §1 | 6/6 checklist PASS |
| SATURATION-NOTES §6 | `ready_for_REMEDIATION_PLAN: true` |
| SATURATION-NOTES §6 | Recommendation **PROCEED_TO_S06** |
| Board row 664 (P38-S05-02) | `done` + APPROVE (saturated) in Notes |

---

### Block 3 — REMEDIATION-PLAN ranked G* + reject list

```bash
grep -n 'G1\|G2\|G3\|G4\|G5\|G6\|G7\|G8\|G9\|Rank order\|§4 Reject' \
  "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md" \
  | head -40 | tee "$EVID/03-remediation-plan-headers.txt"

grep -c '^| [0-9]' "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md" \
  | tee "$EVID/03-reject-row-count.txt"

grep -i 'implement in P38\|No product code' \
  "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md" \
  | tee "$EVID/03-no-implement-language.txt"

# G-001…G-011 coverage in themes or §4
grep -E 'G-00[0-9]|G-01[01]|G-004a|M-001' \
  "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md" \
  | tee "$EVID/03-gap-coverage-grep.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| §2 summary table | G1–G9 with rank order **G1→G3→G4→G5→G2→G6→G7→G8→G9** |
| §4 reject registry | **≥12** rows (plan has 15) |
| §1 executive | M-001 moat preserved |
| Every G-001…G-011 | In theme or §4 defer/reject |
| No implement in P38 | Explicit in header + T7 self-check |
| Phase 39 co-wave | G1 + G3 + G4 named in §1 or §3 |

Spot-read: REMEDIATION-PLAN §2 G1/G3/G4 peer patterns cite upstream `$EV/` — no duplicate JSON required.

---

### Block 4 — Peer cites CG, UA, GF, **MP**

Mechanism cites must be **file:line or `$EV/` paths** — not README-only.

```bash
# CG mechanism (PEER-CG)
grep -n 'tools.ts\|L1163\|L3193\|handleExplore\|watcher.ts' \
  "$P38/scopes/scope-02-codegraph-peer/PEER-CG.md" \
  | head -15 | tee "$EVID/04-peer-cg-cites.txt"

# UA mechanism (PEER-UA-GF)
grep -n 'context-builder.ts\|search.ts\|L25\|L36' \
  "$P38/scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md" \
  | head -15 | tee "$EVID/04-peer-ua-cites.txt"

# GF mechanism (PEER-UA-GF)
grep -n 'validate.py\|symbol_resolution\|graph.html\|EXTRACTED\|INFERRED' \
  "$P38/scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md" \
  | head -15 | tee "$EVID/04-peer-gf-cites.txt"

# MP mechanism (PEER-UA-GF §3 — LOCK)
grep -n '§3\|Mempalace\|searcher.py\|layers.py\|service.py\|fact_checker.py\|wake_up\|_hybrid_rank' \
  "$P38/scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md" \
  | head -20 | tee "$EVID/04-peer-mp-cites.txt"

# GAP-REGISTRY MP column populated
grep -n '| MP |' "$P38/scopes/scope-04-gap-registry/GAP-REGISTRY.md" \
  | head -5 | tee "$EVID/04-gap-registry-mp-column.txt"

# Spot-check MP evidence files exist (S03 run)
ls experiments/runs/2026-08-22-p38-s03-657/evidence/h4-mp-hybrid-search.md \
   experiments/runs/2026-08-22-p38-s03-657/evidence/h6-mp-mcp-surface.md \
   experiments/runs/2026-08-22-p38-s03-657/evidence/h9-mp-fact-check-contrast.md \
   2>&1 | tee "$EVID/04-mp-evidence-files.txt"
```

**Block 4 peer cite checklist (fill in VERIFY-NOTES):**

| Peer | Required mechanism cites | Primary doc | Spot-check |
|------|-------------------------|-------------|------------|
| **CG** | `codegraph_explore` schema + handler; watcher debounce | PEER-CG §3 | `tools.ts` L1163+, L3193+ |
| **UA** | `buildChatContext(query)` + SearchEngine keys | PEER-UA-GF §2 | `context-builder.ts` L25–79 |
| **GF** | EXTRACTED/INFERRED + graph.html orient | PEER-UA-GF §2 | `symbol_resolution.py`, `graph.html` |
| **MP** | hybrid search, wake_up layers, MCP surface, fact_checker contrast | PEER-UA-GF **§3** | `searcher.py` L276–329, `layers.py` L404–431, `service.py` L60–82, `fact_checker.py` L55–78 |

**Pass:** Each peer row has ≥1 mechanism cite; MP §3 present; GAP-REGISTRY matrix includes **MP** column on gap rows H1/H3/H4/H6/H8/H9/H10.

---

### Block 5 — Moat row M-001

```bash
grep -n 'M-001\|§3 Moat\|non-gap' \
  "$P38/scopes/scope-04-gap-registry/GAP-REGISTRY.md" \
  | tee "$EVID/05-moat-m001.txt"

grep -n 'M-001\|moat' \
  "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md" \
  | head -10 | tee "$EVID/05-remediation-m001.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| GAP-REGISTRY §3 | M-001 **non-gap** with Trace \| CG \| UA \| GF \| MP columns |
| REMEDIATION-PLAN §1 | M-001 preserved; remediation merges into moat |
| Distinct from G-010 | M-001 = moat exists; G-010 = under-promotion |

---

### Block 6 — Successor Phase 39 in DR-HANDOFF prep (notes only)

In VERIFY-NOTES, include **successor recommendation table** for S07-02 (do **not** close DR-HANDOFF or scaffold Phase 39):

| Field | Locked value (from REMEDIATION-PLAN §6) |
|-------|------------------------------------------|
| Successor | **Phase 39 — Context orient & harness** (human promotes) |
| Entry co-wave | **G1 + G3 + G4** |
| Secondary queue | G5, G2 compose-first → Phase 39–40; G2 unified explore → Phase 40+ |
| Idle alternative | `no successor` — if human defers implementation |
| P38 outcome | **Plan only** — no implement rows |

Reference REMEDIATION-PLAN §3 Phase 39 bullets and §5 open questions for human promotion.

---

### Block 7 — WRITE VERIFY-NOTES.md

Create `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-07-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — Phase 38 / S07-01

**Date:** …
**Overall:** PASS | FAIL | PARTIAL
**Git SHA:** …
**Evidence:** experiments/runs/…-p38-s07-01-verify/evidence/

## Precondition cites
- P38-S06-02 APPROVE (high confidence)
- SATURATION-NOTES ready_for_REMEDIATION_PLAN: true
- All S00–S06 scope reviews APPROVE

## Block results
| Block | Check | Result | Evidence file |
|-------|-------|--------|---------------|
| 0 | No product Go/TS in P38 | | 00-product-commits-since-promotion.txt |
| 1 | H1–H11 + 7 artifacts | | 01-artifact-existence.txt, manifest.sha256 |
| 2 | S05 saturation APPROVE | | 02-saturation-gate.txt |
| 3 | REMEDIATION-PLAN G1–G9 + rejects | | 03-remediation-plan-headers.txt |
| 4 | Peer cites CG, UA, GF, MP | | 04-peer-*-cites.txt |
| 5 | Moat M-001 | | 05-moat-m001.txt |
| 6 | Phase 39 successor prep | | (this section) |

## H1–H11 coverage
| H* | G-ID | Verdict | Evidence pointer |
|----|------|---------|------------------|
| H1 | G-001 | gap | … |
| … | … | … | … |

## Peer cite summary
| Peer | PASS | Primary cites |
| CG | | PEER-CG §3 |
| UA | | PEER-UA-GF §2 |
| GF | | PEER-UA-GF §2 |
| MP | | PEER-UA-GF §3 |

## Artifact manifest
| Artifact | SHA256 (from manifest) | Archived |
| … | … | yes |

## Successor recommendation (for S07-02)
- **Default:** Phase 39 — G1+G3+G4 entry co-wave
- **Alternative:** no successor (human idle)
- **Never:** TBD

## DR-HANDOFF
Stays OPEN — P38-S07-02 closes + scaffolds Phase 39

## Next
P38-S07-02
```

## Exit criteria

- [ ] `VERIFY-NOTES.md` with block table 0–6 + H1–H11 map + peer cite summary
- [ ] Evidence dir populated under `experiments/runs/…-p38-s07-01-verify/evidence/`
- [ ] Blocks 0–6 executed in order
- [ ] All 7 artifacts archived with `manifest.sha256`
- [ ] Board Notes on **P38-S07-01** only
- [ ] DR-HANDOFF remains OPEN
- [ ] Next: **P38-S07-02**

## Next

`P38-S07-02`
