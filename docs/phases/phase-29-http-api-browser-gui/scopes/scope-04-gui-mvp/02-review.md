# P29-S04-02 — GUI MVP review

## Metadata
- id: P29-S04-02
- todo_ids: [P29-S04-02]
- role: reviewer
- skills: [code-review-and-quality, accessibility-auditing]
- mcps: [cursor-ide-browser]
- verification: mixed
- hooks: []

## Objective

Independent review of the S04 MVP SPA vs UX-IA P0, Law 19, ADR static serve path, and implementer Notes. Fresh subagent. Small inline fixes OK; structural gaps spawn `P29-S04-02a` / `02b`.

## Session start

Follow agent-loop-protocol Session start. Do not share the implementer session. Do not start S05.

## References

- [00-PLANNER.md](00-PLANNER.md) locked defaults
- [01-implement.md](01-implement.md) layout, routes, op map
- [UX-IA.md](../scope-02-ux-ia/UX-IA.md) `gui_ship: S04`
- [ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)
- [api/openapi.yaml](../../../../../api/openapi.yaml)
- Implementer board Notes on **P29-S04-01**

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f web/dist/index.html || test -f web/package.json
test -d internal/httpapi
grep -q 'case "serve"' cmd/trace/root.go
# Prefer: cd web && npm run build
```

If `web/` missing or build broken with no Notes evidence, **fail** or spawn 02a — do not rewrite the whole SPA in review unless a small fix restores exit criteria.

## Checklist

### Ship matrix (S04 only)

- [ ] Project gate: blocked empty when `store_ready: false` / unreachable
- [ ] Overview: goals summary (search/tasks), active task, loop status + gate strip
- [ ] Tasks: list + detail + `TRACE_TASK_ID` copy; light transition uses envelope-only on deny
- [ ] Loop: read-only status/gate — **no** next/apply/reset controls calling those APIs
- [ ] Graph: stub only — budgeted `getGraph` (center+max_nodes) and/or S05 placeholder; **no** xyflow/Three; truncated honesty when `truncated`
- [ ] Discoveries: read list/search + detail — **no** create/promote
- [ ] Seed: status + honesty copy — **no** export/import CTAs
- [ ] Settings: theme + token paste + version/project display
- [ ] **No** reviews UI; **no** agents nav/screen

### Law 19 / contract

- [ ] All data via `/v1` only (no SQLite/IndexedDB SoT, no editing `graph.json` as SoT)
- [ ] Search consumes `{items:[...]}` — not CLI `hits`
- [ ] Error UI maps ADR envelope `error.code` / `message`
- [ ] Static assets from `web/dist` via `trace serve` (placeholder gone when dist present)

### Build + smoke

- [ ] `npm run build` (or equiv) evidence in Notes or re-run PASS
- [ ] Smoke: serve → UI → tasks visible (or honest empty)
- [ ] Vite proxy documented for DX; no production CORS `*`

### A11y floor (UX-IA §7)

- [ ] Primary nav focusable with visible focus
- [ ] Icon-only controls have accessible names
- [ ] Status/gate strips are text (not icon-only severity)
- [ ] Prefer `prefers-reduced-motion` respected for non-essential motion

## Findings policy

| Severity | Action |
|----------|--------|
| blocker / high | Inline small fix **or** spawn `P29-S04-02a` (implement) + `02b` (review) immediately below this row |
| medium | Prefer spawn unless trivial one-liner |
| low / nit | Note residuals; thicken upcoming **S05-00** only if needed |

## Exit criteria

- [ ] No open blocker/high without a pending spawn
- [ ] Confidence **medium+** with evidence (build + smoke + checklist)
- [ ] Board Notes; next runnable **P29-S05-00** (after this row `done`)

## Todo updates

Status + notes on **P29-S04-02** only (plus spawn rows if inserted).

## Next

**P29-S05-00**
