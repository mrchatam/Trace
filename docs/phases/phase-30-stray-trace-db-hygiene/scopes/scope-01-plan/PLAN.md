# PLAN — Phase 30 stray root trace.db hygiene

**Date:** 2026-08-21 · **Author:** P30-S01-01 · **Depends:** INVESTIGATION.md + P30-S00-02 PASS

## 1. Verdict taken

From [`../scope-00-investigate/INVESTIGATION.md`](../scope-00-investigate/INVESTIGATION.md):

> **agent hygiene** — Trace never creates or opens `<projectRoot>/trace.db`; every executable store join is `<root>/.trace/trace.db`, and a disposable repro shows a **0-byte** root stub from `python3` `sqlite3.connect('trace.db')` after `trace init`, while CLI continues to use the live `.trace/` store.

**Also recorded (S00 / P30-S00-02):** INTAKE **confirmed**; **no Trace dual-store bug**; **no store-path change**. Canonical store remains `<projectRoot>/.trace/trace.db` via `store.Open` / `OpenExisting`. S00 overturn path for path redesign is **closed**. This plan has **no “path fix” task**.

## 2. Goals / non-goals

### Goals

- Operator + agent clarity: only `.trace/trace.db` is authoritative; root `trace.db` is never the Trace store.
- Ship **T1–T4 only**: docs, warn-once-on-open (stderr, non-fatal), root `/trace.db` gitignore, tests.

### Non-goals (explicit)

- Store-path redesign / “path fix” / changing `dbPath` join
- Silent or default delete of root `trace.db`
- Optional documented delete flag (deferred **future-only**; **not** implemented in S02)
- Treating root stub as a second store (open, migrate, or sync)
- HTTP / daemon / GUI scope creep (Phase 29 owns `trace serve`)
- New phase for a Trace dual-store bug (**none found**)
- Product Go outside the locked T2 warn + T4 tests surfaces

## 3. Tasks

| ID | Change | Files (likely) | Acceptance | Deps |
|----|--------|----------------|------------|------|
| T1 | Docs: state canonical path `.trace/trace.db`; never open/create project-root `trace.db` (agents use CLI/MCP) | `AGENTS.md` (stack one-liner); `docs/rules/project-rules.md` (Store row); `CONTRIBUTING.md` (`.trace/` local store section); optional one-liner in `cmd/trace/help.go` (init/store clarity only). **Skip** Phase 29 HTTP docs. **Skip** `internal/install/*` rule bodies unless a single shared agent-facing string already documents store location — prefer AGENTS + project-rules as durable agent surfaces (INVESTIGATION install gap closed by those, not by rewriting enforcement hooks). | Readers see: live store is `.trace/trace.db` only; root `trace.db` is not Trace; no dual-store claim. | — |
| T2 | Warn once per successful open path when `<absRoot>/trace.db` exists as a **regular file**; open still uses `.trace/trace.db` | `internal/store/open.go` — single choke in `openStore` (covers `Open`, `OpenExisting` → `openStore`, and `Restore` rebind). Optional package-level `io.Writer` (default `os.Stderr`) for tests. **No** join change. | Trigger after Abs root; once per `openStore` call that sees stray; stderr (or test writer); non-fatal; never deletes/renames/opens root file; `DBPath()` still under `.trace/`. | — |
| T3 | Root-only gitignore `/trace.db`; keep `.trace/` ignored | This repo `.gitignore` (add `/trace.db` beside existing `.trace/`); consumer scaffold `fixtures/x0/.gitignore` (same pattern); CONTRIBUTING (or T1) mentions consumers should ignore root `/trace.db` without un-ignoring `.trace/`. | Pattern is `/trace.db` (leading slash); does not target `.trace/trace.db`; `.trace/` remains ignored. | Parallel with T1 if file-disjoint; else after T1 prose for CONTRIBUTING |
| T4 | Tests for warn + path invariants | `internal/store/*_test.go` (new cases beside existing open tests) | Named cases below; `go test ./internal/...` PASS. | T2 |

## 4. Warn design (T2)

| Item | Decision |
|------|----------|
| **Choke** | `openStore` in `internal/store/open.go`, after `filepath.Abs` succeeds and before (or immediately after) `.trace/` mkdir — **one** place so CLI, MCP (`OpenExisting` → `openStore`), and `trace serve` (`store.Open`) all warn without per-adapter duplication. |
| **Trigger** | `os.Stat(filepath.Join(absRoot, "trace.db"))` succeeds and mode is regular file. Directories / missing / permission errors on the stub check → **no warn** (do not fail open). |
| **Frequency** | Once per `openStore` invocation that sees the stray. Not a persistent on-disk flag. Acceptable if a process opens twice and warns twice. |
| **Channel** | stderr only (or injectable `io.Writer` for tests). Never stdout / JSON MCP payloads. |
| **Severity** | Non-fatal. Open **must still succeed** (subject to existing lock/auth/db errors). Never fail-closed on stray presence. |
| **Must not** | Delete, rename, migrate, or `sql.Open` the root file; change `filepath.Join(absRoot, ".trace", "trace.db")`; treat stub as second store. |
| **Draft message** | `trace: warning: project-root trace.db exists but is not the Trace store; using .trace/trace.db. Do not open or create a root trace.db (agents: use CLI/MCP).` |

## 5. Implementation order for S02

1. **T2** — add stray check + warn writer hook in `openStore` (no path join change).
2. **T4** — tests first against T2 API (stub present / absent / stub untouched).
3. **T3** — `.gitignore` + `fixtures/x0/.gitignore` (`/trace.db`).
4. **T1** — AGENTS / project-rules / CONTRIBUTING (+ optional help one-liner).

**Parallel-safe:** T3 gitignore files are disjoint from T2/T4 Go; T1 docs are disjoint from T2/T4 except CONTRIBUTING ignore sentence (coordinate with T3). Prefer T2+T4 before docs so acceptance is executable.

Green bar before S02-01 done: `go test ./internal/...` PASS.

## 6. Test plan

**Command:** `go test ./internal/...`

| Case | Expect |
|------|--------|
| Stub present + `Open` | Warn observed (stderr or test writer); `DBPath()` / open target is `.trace/trace.db`; open succeeds |
| Stub present + `OpenExisting` (after init) | Same warn + path + success |
| No stub | No warn; open succeeds |
| After warn | Live DB under `.trace/`; root stub still present; size unchanged |

Do not assert delete, rename, or open of the root file. Do not require HTTP/GUI tests.

## 7. Risks / residuals

Low-risk hygiene only.

| Residual | Disposition |
|----------|-------------|
| Documented optional delete of root stub | **Future-only** note — not S02 |
| Agents that ignore docs and still `sqlite3.connect('trace.db')` | Mitigated by T2 warn + T3 gitignore; cannot prevent all agent mistakes |
| Warn on every open in long-lived `serve` | Acceptable once-per-open; not a persistent suppress flag in Phase 30 |
| Install bundled enforcement strings omit store path | Closed via T1 AGENTS/project-rules; do not overload Phase 23 enforcement hook text |

## 8. S02 handoff

**P30-S02-00** must re-read this `PLAN.md` as SoT, re-assert S00 **agent hygiene** + INTAKE confirmed + **no store-path change**, lock implementer defaults to **T1–T4** exactly as above (warn: once-per-open stderr non-fatal in `openStore`; gitignore `/trace.db`; no silent delete; no HTTP creep), thicken `01-implement.md` / `02-review.md` against these file lists and acceptance rows, then hand to **P30-S02-01**. Do not re-open S00 investigation.
