# P17-S04-00 — Phase VERIFY / portable-graph-git (FINAL)

## Metadata
- id: P17-S04-00
- todo_ids: [P17-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 17 VERIFY evidence: **S01–S03 named DF regressions** + **two-clone git-JSON recipe** (no shared `.trace/`; init + import + index + why/context + plan readable) + carry-forward gates. Decide **DR-HANDOFF** = **`no successor`** (default). **DF-86 hook absence is non-fail.** **No product Go.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [DF-84-FORWARD.md](../../DF-84-FORWARD.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 16 S06 VERIFY [`../../../phase-16-assert-root-and-surfaces/scopes/scope-06-phase-verify/00-PLANNER.md`](../../../phase-16-assert-root-and-surfaces/scopes/scope-06-phase-verify/00-PLANNER.md)
- S01–S03 REVIEW-NOTES (all **APPROVE high**): [S01](../scope-01-seed-export/REVIEW-NOTES.md), [S02](../scope-02-commit-convention/REVIEW-NOTES.md), [S03](../scope-03-idempotent-import/REVIEW-NOTES.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Depends-on: S01–S03 APPROVE (landed). Default handoff remains **`no successor`**. Unattended: no Plan-mode switch; no product Go.

## Depends-on S03 FINAL (landed — do not re-lock)

Sibling [../scope-03-idempotent-import/00-PLANNER.md](../scope-03-idempotent-import/00-PLANNER.md) **FINAL** delivers:

- Idempotent **`seed import`**: entity UUID upsert (insert-only events), duplicate-link no-op, plan-tree upsert, findings/alternatives upsert
- Named **`TestSeedImportIdempotent`**, **`TestSeedImportDuplicateLinksNoOp`**, **`TestSeedImportSameIdLastWins`**, **`TestSeedImportPlanTreeIdempotent`**
- CONTRIBUTING merge **union-by-id** (entities + plan arrays); last-import-wins upsert after human git resolve

S04 VERIFY assumes re-import and two-clone merge behavior from S03; does not re-implement upsert.

## Depends-on (S01–S03 — landed; live named tests confirmed 2026-08-17)

| Scope | Board | Named tests imported (live `func Test*` exists) |
|-------|-------|--------------------------------------------------|
| S01 DF-80/84/85 | **APPROVE high** (P17-S01-02) | `TestSeedExportRoundTrip` (plan tree + findings); `TestSeedExportOmitsDeniedSurfaces`; `TestSeedExportWritesExportedAtCommit`; P16 keepers `TestSeedImportAndWhy` / `TestSeedImportDiscoveryMentionsTask` / `TestSeedImportImpactFindings` / `TestSeedImportFromIDAliases` / `TestSeedImportRelativePathAgainstC` / `TestSeedImportMissingEndpointsMessage` |
| S02 DF-82/85 | **APPROVE high** (P17-S02-02) | `TestHelpSeedExportPath`; `TestHelpHandoffSoT`; `TestAsOperatorFlagIdentityDocs`; S01 `TestSeedExport*` carry-forward |
| S03 DF-81/83/84 | **APPROVE high** (P17-S03-02) | `TestSeedImportIdempotent`; `TestSeedImportDuplicateLinksNoOp`; `TestSeedImportSameIdLastWins`; `TestSeedImportPlanTreeIdempotent`; S01/S02 keepers green |

## Live residuals → DR-HANDOFF decision (2026-08-17)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 17 | DF-80/81/82/83/84/85 | Closed by S01–S03 APPROVE — VERIFY must **re-prove named tests** |
| Explicit residual OK into VERIFY | Encryption wontfix; reviews omitted from default export; **DF-86 hook absent**; CGO=0 `cmd/trace` tree-sitter FAIL; no dedicated `work_state` re-import named test | **Do not fail VERIFY** for these |
| Goals sequence #2–#4 | Research S05 / `plan simulate` / D21+ | Stay off-board — **not** auto-boarded |
| Hosted MCP / HTTP / OAuth | TODO Later developments (separate repo) | **Not** a VERIFY successor |
| Product bar | `./cmd\|internal\|evals` with CGO1 | Prefer product pkgs over full-module `./...` when graphify space FAIL present |

**DR-HANDOFF = `no successor`.** S04 Notes must not promote research S05 / plan simulate / D21+ / hosted server. Phase 16 historical `no successor` left intact as history.

## Planner work
1. [x] Import named `-run` filters from S01–S03 REVIEW-NOTES
2. [x] Lock verify command set + two-clone recipe + evidence path
3. [x] Thicken 01-verify / 02-scope-review / SCOPE-TODOS
4. [x] Stamp DR-HANDOFF intent; mark this prompt **FINAL**; next **P17-S04-01**

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Phase 17 portable-graph-git closeout — S01–S03 named DFs + **two-clone recipe** — **not** a new planted eval gate |
| Path lock | **`trace/graph.json`** committed or exportable at VERIFY time |
| Two-clone recipe | **Fail bar:** two independent temp dirs/worktrees from the **same git tree** containing `trace/graph.json`; **neither** reads the original `.trace/`; each: `init` → `seed import trace/graph.json` → `index` → `why` on a seeded decision + `context` on a seeded task + plan hierarchy readable (`plan show`). **Offline** — no account, no HTTP, no MCP server |
| Named test (two-clone) | **`TestPortableGraphTwoCloneWhyContextPlan`** in `cmd/trace` (CGO=1). **If absent at S04-01 preflight:** S04-01 implements as VERIFY remediation (allowed product Go on verify row) — do **not** weaken bar by skipping |
| DF-86 | **`trace install git-hook` absent is non-fail.** VERIFY records grep evidence; hook **must not** wrap `git commit` if ever added |
| Residuals OK | Encryption wontfix; reviews omitted; DF-86 hook absent; CGO=0 `cmd/trace` FAIL; S03 `work_state` preservation SQL-only — **do not fail VERIFY** |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat (ceiling as of P16 close); p0x; x0; Gate C `dry_run:false`; product `./cmd\|internal\|evals` |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` — **product pkgs PASS**; R3 graphify space FAIL on full `./...` OK; R4 CGO0 analyzers / `cmd/trace` CGO0 FAIL OK |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; **except** two-clone named test if absent + spawn remediation if fail |
| Evidence artifact | **`scopes/scope-04-phase-verify/VERIFY-NOTES.md`** (verdict + evidence table + clone recipe + law checks + handoff). Gate C inspect only: `docs/verification/gate-c-x0/` |
| Spawn | On fail: `P17-S04-01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — **S04-01 starts** Notes; **S04-02 owns completion**. Do **not** auto-board research S05 / plan simulate / D21+ / hosted server |
| Forbidden | Product features on VERIFY (except two-clone test if absent); claiming P16 DFs; rewriting P16 history; new MCP seed tool; pointing `trace-mcp` at internet |

### Locked verify command set (FINAL)

Per-scope named `-run` lines match S01–S03 REVIEW-NOTES (independent DF re-proof). **CGO0 `./cmd/trace/...` is R4 — do not use as fail bar.** `GOMODCACHE`+`GOPROXY=off` on full product bar.

```bash
# --- S01 DF-80/84/85 export + P16 import keepers (CGO=1 authoritative for cmd/trace) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit|TestSeedImportAndWhy|TestSeedImportDiscoveryMentionsTask|TestSeedImportImpactFindings|TestSeedImportFromIDAliases|TestSeedImportRelativePathAgainstC|TestSeedImportMissingEndpointsMessage'

# --- S02 DF-82/85 help/path + actor≠auth keepers ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpSeedExportPath|TestHelpHandoffSoT|TestAsOperatorFlagIdentityDocs|TestSeedExport'

# --- S03 DF-81/83/84 idempotent import + S01/S02 carry-forward ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportIdempotent|TestSeedImportDuplicateLinksNoOp|TestSeedImportSameIdLastWins|TestSeedImportPlanTreeIdempotent|TestSeedExportRoundTrip|TestHelpSeedExportPath'

# --- Two-clone git-JSON recipe (fail bar) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run TestPortableGraphTwoCloneWhyContextPlan

# --- Honesty: Paths A/B/C + Gate G (CGO-free) ---
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# --- Gate E / F / capability ablation ---
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# --- Gate H + compat (compat covers mig ceiling as of P16 close) ---
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# --- P0-X + X0 ---
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# --- Product regression bar ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Optional (strong evidence, **not** substitutes for package PASS):

```bash
# DF-86: grep — no `install git-hook` / `git-hook` implementation in product Go (non-fail if absent)
# .gitignore: `.trace/` only; `trace/graph.json` not ignored
# Gate C artifact inspect (jq/grep OK): docs/verification/gate-c-x0/ metrics dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# P16 catalog 10 + trace_impact still green (carry-forward spot-check via compat/product bar)
# Do NOT fail for R3 graphify space FAIL or R4 CGO0 analyzers / CGO0 cmd/trace FAIL
# Do NOT fail for encryption wontfix / reviews omitted / DF-86 hook absent
```

### Locked two-clone shell recipe (manual evidence — S04-01 documents in VERIFY-NOTES)

Use when named test is green **or** as independent corroboration. **Offline; no HTTP; no MCP server.**

```bash
# Prerequisite: repo root has trace/graph.json (export first if missing: trace seed export -o trace/graph.json)
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
GRAPH="$REPO_ROOT/trace/graph.json"
test -f "$GRAPH"

run_clone() {
  local dest="$1"
  mkdir -p "$dest"
  rsync -a --exclude='.trace' --exclude='.git' "$REPO_ROOT/" "$dest/"
  ( cd "$dest" && \
    trace init && \
    trace seed import trace/graph.json && \
    trace index && \
    trace plan show && \
    trace why decision "<seeded-decision-id>" && \
    trace context "<seeded-task-id>" )
}

CLONE_A="$(mktemp -d)"
CLONE_B="$(mktemp -d)"
run_clone "$CLONE_A"
run_clone "$CLONE_B"
# Record: neither clone used REPO_ROOT/.trace/; each has its own .trace/
```

Seeded decision/task ids: from `trace/graph.json` in repo or from round-trip fixture ids documented in VERIFY-NOTES. Named test **`TestPortableGraphTwoCloneWhyContextPlan`** is the **primary** fail bar; shell recipe is **secondary corroboration**.

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + two-clone recipe + DR-HANDOFF locked (`no successor`)
- [x] Evidence path = `VERIFY-NOTES.md` in this scope folder
- [x] SCOPE-TODOS + board Notes; next `P17-S04-01`
- [x] Product Go — **not** this row

## Out of scope
- Running VERIFY (S04-01)
- Product Go / new MCP tools / daemon / DF-86 hook (unless S04-01 remediation)
- Rewriting Phase 00–16 `done` history
- Auto-boarding research S05 / plan simulate / D21+ / hosted MCP
- Claiming P16 DFs or rewriting P16 DR-HANDOFF

## Next
**P17-S04-01** (independent VERIFY run).
