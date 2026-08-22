# P32-S02-01 — Implement API gaps + P32-PORT

## Metadata
- id: P32-S02-01
- todo_ids: [P32-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, api-and-interface-design]
- mcps: []
- verification: automated
- hooks: []

## Objective

Ship **three** deliverables in one row (S02 is **not** API-only):

1. **Library/OpenAPI:** write `NO-GAPS.md` (UX-IA + live: no new library-backed `/v1` ops for the inspector map).
2. **Client glue:** add missing `getImpact` wrapper in `web/src/api/ops.ts` (OpenAPI `/v1/impact` + `handleGetImpact` already exist).
3. **P32-PORT (always):** RESEARCH / OPEN-PORT-MULTI **#1 minimum** — friendly `EADDRINUSE` / in-use stderr + `--addr` examples. **Do not** ship auto free-port (#2) this row.

**No** inventing `/v1/path`, requiring `listChanges`/`listRegressions`, multi-tenant hosting, or changing default bind off loopback.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [00-PLANNER.md](00-PLANNER.md) — Must-answer locks below are **final**
- S01 [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md) — API gap list + P32-PORT note
- S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md) — P32-PORT prefer #1
- Live: `internal/httpapi/` (`server.go` Listen, `handlers_p1.go` `handleGetImpact`, `bind.go` `DefaultAddr`), `cmd/trace/serve.go`, `cmd/trace/serve_test.go`, `api/openapi.yaml` `/v1/impact`, `web/src/api/ops.ts`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Law 19 | Handlers/CLI/UI are **adapters only** — call existing library/domain; no business-logic fork in `web/` or HTTP handlers |
| Library / OpenAPI gaps | **None** for inspector map → author `scopes/scope-02-api-gaps/NO-GAPS.md` with evidence (cite UX-IA gap table + live OpenAPI/handlers). **Do not** add `/v1/path` or new core ops |
| Client glue | **Ship** `getImpact` in `web/src/api/ops.ts` (required for S03 inspector Impact section) |
| `getImpact` shape | Match `getWhy` / `getContext` style — see locked signature below |
| OpenAPI change | **None** for impact (already present). Do **not** invent path; do not require changes/regressions wrappers |
| P32-PORT | **#1 only (min)** — friendly in-use detection + stderr/`--addr` guidance. **Defer #2** auto free-port / `:0` (note in Notes / NO-GAPS or P32-PORT Notes; S05 docs own multi-project narrative polish) |
| Default bind | Keep `httpapi.DefaultAddr` = `127.0.0.1:7432`; fail-on-conflict policy **unchanged** (still exit fail — only message quality improves) |
| Remote bind | Still requires `--allow-remote` (+ token policy as today); no `0.0.0.0` default |
| Defer hatch | `NO-PORT-CHANGE.md` **discouraged** — only with written reason; prefer ship #1 |
| Product scope | Adapter glue + serve UX only — **no** graph-home / inspector UI (that is S03) |

### Must-answer locks (encode in implementation + Notes)

| # | Question | Locked answer |
|---|----------|---------------|
| 1 | Library API gaps? | **`NO-GAPS.md`** — no new library/OpenAPI ops. Evidence: UX-IA API gap list; live `/v1/why\|context\|impact\|graph\|search\|reviews` + handlers. |
| 2 | `getImpact` wrapper? | **Ship** in `ops.ts` (see signature). Response type: loose JSON object (`Record<string, unknown>` or equivalent) — OpenAPI schema is `additionalProperties: true`. |
| 3 | P32-PORT options? | **#1 min only.** Detect bind-in-use; print friendly stderr with `--addr` examples. **#2 deferred.** |
| 4 | Port-in-use test plan? | Automated Go test: occupy `127.0.0.1:<ephemeral>`, second `Listen` / `serve` path yields **non-OK** exit and stderr contains in-use guidance + `--addr` hint (see Role work). |
| 5 | Help / stderr copy? | On in-use: state address already in use; suggest second project via distinct `--addr` (example `127.0.0.1:7433`); keep short. Optionally add one multi-project sentence to `trace serve` help Usage — not a substitute for conflict-path message. |

### Locked `getImpact` client signature

Place near `getWhy` / `getContext` in `web/src/api/ops.ts`:

```typescript
/** getImpact — GET /v1/impact (OpenAPI). Prefer task_id when selection is a task (UX-IA). */
export function getImpact(taskId: string, opt: TokenOpt = {}) {
  return apiFetch<Record<string, unknown>>('/v1/impact', {
    ...opt,
    query: { task_id: taskId },
  })
}
```

**Do not** document/require `decision_id` in the client this row (handler accepts it live; OpenAPI documents `task_id` only — stay on the published contract). S03 will call with task id when selection resolves to a task; omit/honest-empty when non-task.

### Locked P32-PORT stderr intent (wording may vary; must convey)

When listen fails because the address is in use:

```text
serve: address already in use (<addr>)
hint: another process (often trace serve) is bound there.
  For a second project, pick a free port, e.g.:
    trace serve -C /path/to/other --addr 127.0.0.1:7433
```

Classify via portable check (`errors.As` to `*net.OpError` / `EADDRINUSE` / message contains `address already in use`) — prefer a small helper under `internal/httpapi` or `cmd/trace` rather than string-matching only in one ad-hoc place. CLI remains the user-facing printer (`serve.go`).

## Preflight (confirm in Notes; do not change unless implementing)

1. `ops.ts`: has `getWhy` / `getContext`; **no** `getImpact`.
2. `api/openapi.yaml`: `/v1/impact` `operationId: getImpact`; `task_id` query optional in schema; handler requires `task_id` **or** `decision_id`.
3. `httpapi.DefaultAddr` = `127.0.0.1:7432`; `ListenAndServe` = bare `net.Listen`; conflict → `serve: %v` today.
4. UX-IA gap table: library none; client `getImpact`; P32-PORT #1 always.

## Role work

### A — Library story → `NO-GAPS.md`

Author `docs/phases/phase-32-graph-first-gui/scopes/scope-02-api-gaps/NO-GAPS.md`:

- State: no library-backed OpenAPI additions required for Phase 32 inspector map.
- Cite UX-IA API gap table + live paths (`openapi.yaml`, handlers).
- Explicitly list **out**: `/v1/path`, new changes/regressions requirements.
- Note: client `getImpact` + P32-PORT are **separate** deliverables in this same row (S02 is not “API-only”).

### B — Client glue → `getImpact`

1. Add locked `getImpact` to `web/src/api/ops.ts`.
2. Export only — **do not** wire Graph inspector (S03).
3. Sanity: TypeScript still typechecks for `web/` if the repo has a check script; otherwise ensure signature matches `apiFetch` / `TokenOpt` patterns. No need to invent a new vitest suite if none exists.

### C — P32-PORT #1

1. Detect address-in-use on listen failure.
2. Map to friendly stderr from `cmd/trace/serve.go` (and/or shared helper) including `--addr` example(s).
3. Keep default `127.0.0.1:7432` and fail-on-conflict (non-zero exit).
4. Tests (preferred locations: `cmd/trace/serve_test.go` and/or `internal/httpapi`):
   - Bind a listener on `127.0.0.1:0` (or fixed free port), attempt second listen / `run([]string{… "serve", "--addr", …})`, assert exit ≠ OK and stderr includes in-use + `--addr` guidance.
   - Do **not** weaken `TestServeRefuseRemoteCLI`.
5. Optional: one help-text line on multi-project distinct `--addr` (S05 still owns full quickstart docs).

### Out of scope this row

- Graph-home shell, inspector UI, nav reweight (S03+)
- Auto-port / `:0` / port scanning (#2)
- `gui-quickstart.md` multi-project essay (S05)
- OpenAPI regeneration unless a real schema bug blocks `getImpact` (unexpected — escalate in Notes)

## Exit criteria

- [ ] `NO-GAPS.md` present with evidence (library/OpenAPI)
- [ ] `getImpact` exported from `web/src/api/ops.ts` per locked signature
- [ ] P32-PORT #1 shipped: friendly in-use path + `--addr` guidance
- [ ] Automated test(s) cover port-in-use / friendly message (or Notes justify with command transcript if truly untestable — prefer test)
- [ ] Loopback defaults unchanged; no `0.0.0.0` without `--allow-remote`
- [ ] No `/v1/path` invented; no changes/regressions wrappers required
- [ ] Board Notes cite files + commands (`go test …`, any `web` check)

## Minimal todos

- [ ] Write `NO-GAPS.md`
- [ ] Add `getImpact` to `ops.ts`
- [ ] Implement P32-PORT #1 helper + CLI stderr (+ optional help tweak)
- [ ] Add/adjust Go tests for in-use path
- [ ] Run targeted tests; update **P32-S02-01** board Notes

## Todo updates

Status + notes on **P32-S02-01** only.

## Next

`P32-S02-02`
