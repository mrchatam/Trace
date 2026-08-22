# Scope 05 — board map (polish)

Production polish + docs (quickstart / multi-project ports / optional screenshots). Serial: **S05-00 → S05-01 → S05-02**. **No** new explorer features.

## Locked for implement (from P32-S05-00)

| Lock | Value |
|------|-------|
| Port story | P32-PORT **#1** as shipped: fail-on-conflict; friendly in-use; distinct `--addr` per project |
| Not shipped | Auto free-port / `:0` (**#2 deferred**) — do not document as available |
| Docs primary | `docs/gui-quickstart.md` multi-project / ports section **required** |
| Docs secondary | `web/README.md` if single-port-only |
| OPEN #3/#4 | Docs (optional helper not required) |
| UI hero | Explore `/` canvas-first + inspector — not Overview ops dual-card |
| Craft cite (screenshots) | `.graph-shell` / `--graph-canvas-height` / denser `PacketView` / calm nodes / settle+node+chrome motions + `prefers-reduced-motion` |
| Residuals | S04 chrome box-shadow nit optional; keyboard-via-list defer OK |
| Out | New depth/IA/API; Three.js; public bind defaults; auto-port |

## Board rows

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 564 | P32-S05-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner |
| 565 | P32-S05-01 | [01-implement.md](01-implement.md) | Implementer |
| 566 | P32-S05-02 | [02-review.md](02-review.md) | Reviewer |

## Implementer checklist (mirror 01)

- [x] `gui-quickstart.md` multi-project `--addr` + fail-on-conflict honesty
- [x] Explorer orientation (graph home, not ops hero)
- [x] `web/README.md` ports/DX if needed
- [x] Residuals fixed or deferred with reason
- [x] Optional screenshots of Explore craft — skipped (prose craft cues; no docs/assets)
- [x] Tests if code touched; Board Notes with paths — docs-only (no test gate)

## Reviewer focus (mirror 02)

Port docs accuracy vs live `FormatAddrInUseMessage` / help; no auto-port claim; Law 19 / bind defaults; residual closure; UI docs = explorer not Overview.
