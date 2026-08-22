# P40-S00-01 — Implement (G5 GUI graph orient)

## Metadata
- id: P40-S00-01
- todo_ids: [P40-S00-01]
- role: implementer
- skills: [frontend-ui-engineering, incremental-implementation, accessibility-reviewer]
- mcps: [user-trace]
- verification: mixed

## Objective

Implement **G5**: graph-first onboarding UX on the existing `/` Explore route — orient panel, moat-first narrative, confidence-labeled budget UX, install hook docs (G-008). Law 19: thin GUI adapter over canonical HTTP API.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 19
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [REMEDIATION-PLAN G5](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-008](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Live anchors (P40-S00-00 verified 2026-08-22):
  - `web/src/App.tsx:21–22` — `/` → Graph
  - `web/src/layout/Nav.tsx:4` — nav **Explore** (orient copy uses Explore; `<h1>` stays “Graph”)
  - `web/src/screens/Graph.tsx` — mount orient after `<h1>` / before `.page-lead` (`:456–461`)
    - Law 6–7 banner `:470–473`; truncation `:605–611`; budget line `:612–617`
  - `web/src/lib/overviewCompose.ts` — seed/cap constants (do not fork graph logic into orient)
  - `web/src/lib/overviewCompose.test.ts` — node:test pattern (no vitest in `web/package.json`)
  - `web/src/styles/app.css`, `web/src/styles/tokens.css` — styling
  - `CONTRIBUTING.md:64–72` — moat-first tone; add graph-first GUI subsection here or adjacent
  - `internal/install/bootstrap_hint.go` — optional one-line `trace serve` → `/` pointer
  - `internal/httpapi/handlers_retrieval.go` — graph/search handlers (no change unless library adds orient endpoint)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-008 |
| Verdict | **Accept** — enhance existing graph route |
| Orient panel | Dismissible; persist dismiss in `localStorage` key e.g. `trace.orient.dismissed` |
| Copy lead | Task loop primary: Tasks → Loop → gate → review; graph is **orient entry**, not moat replacement |
| Budget honesty | Reuse/enhance existing truncation banner + budget line — add plain-language confidence label |
| Install narrative | Doc-only touch: `CONTRIBUTING.md` and/or install hint — `trace serve` / `trace gui` graph-first path |
| Must preserve | Seed compose, ReactFlow canvas, Inspector, search/manual center, empty states |
| Must not | New business rules in TS; cap increases; dump API; `trace_explore`; compiler/MCP changes |
| Tests | node:test strip-types if helper extracted; manual/browser verify for panel UX |

## Touch-list (adapter → docs)

| Step | File | Action |
|------|------|--------|
| 1 | `web/src/components/GraphOrientPanel.tsx` (new) | Orient panel: moat-first copy, Laws 6–7 budget explainer, links to `/tasks` + `/loop`, dismiss control |
| 2 | `web/src/screens/Graph.tsx` | Mount orient after `<h1>` / before `.page-lead`; wire dismiss + `localStorage`; `data-testid="graph-orient-panel"` |
| 3 | `web/src/styles/app.css` | Orient panel styles (use existing tokens) |
| 4 | `web/src/components/GraphOrientPanel.test.ts` (optional) | node:test strip-types — dismiss persistence / render when not dismissed |
| 5 | `CONTRIBUTING.md` | Short **Graph-first GUI** subsection: `trace serve` → `/` Explore orient; Law 19 adapter note |
| 6 | `internal/install/bootstrap_hint.go` or `agents.go` (optional) | One-line graph-first pointer if fits existing hint pattern |

**Explicit non-touch:**

- `internal/httpapi/` — unless library adds orient helper first (default: **no new route**)
- `internal/compiler/`, `internal/mcp/` — S01 scope
- `web/src/lib/overviewCompose.ts` — graph algorithm unchanged
- Cap constants (`UI_CAP`, `SEED_MAX_NODES`, etc.)

## Implementation order

```text
1. GraphOrientPanel component + copy (moat-first, budget honesty)
2. Wire into Graph.tsx with dismiss localStorage
3. Styles + a11y (focus, aria-label on dismiss)
4. CONTRIBUTING graph-first install hook narrative
5. Optional node:test; browser smoke on `trace serve`
6. Self-check G5-A1–A6 before marking done
```

## Acceptance criteria (must pass)

| ID | Criterion | Assert |
|----|-----------|--------|
| G5-A1 | First-visit orient panel visible on `/` | Panel renders when `localStorage` dismiss flag absent |
| G5-A2 | Moat-first narrative | Copy mentions task loop + gate path before graph-only workflow |
| G5-A3 | Dismiss persistence | Dismiss hides panel; reload respects flag |
| G5-A4 | Budget/confidence honesty | Enhance `:470–473` banner and/or `:605–611` truncation + `:612–617` budget line with plain-language confidence labels |
| G5-A5 | Law 19 adapter | No new graph retrieval logic in `web/` — uses existing `ops.ts` APIs only |
| G5-A6 | Install hook narrative | CONTRIBUTING (or install hint) documents graph-first GUI entry after `trace serve` |
| G5-A7 | No regression | Existing graph load, seed compose, expand, Inspector still work |

## Regression checks

```bash
cd web && node --experimental-strip-types --test src/lib/overviewCompose.test.ts
cd web && npm run build
```

- Existing `overviewCompose.test.ts` stays green (7/7)
- Manual: empty project, partial seed fail, truncation banner still shown at `:605–611`

## Role work

1. Build orient panel as pure presentation — no duplicate of `overviewCompose` algorithms.
2. Integrate without breaking Graph loading states (`overviewLoading`, `noSeeds`, etc.).
3. Write install hook narrative — align with Phase 39 moat-first CONTRIBUTING tone (`:68–72` area).
4. Self-check G5-A1–A7; record evidence in board Notes.

## Exit criteria

- [ ] G5-A1–A7 met with evidence in Notes
- [ ] No Law 6–7 / 19 regression
- [ ] Board row → `done` with files touched + verify command

## Next

`P40-S00-02`
