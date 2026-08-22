# P32-S02-00 — Scope planner (API gaps + P32-PORT)

## Metadata
- id: P32-S02-00
- todo_ids: [P32-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design]
- mcps: []
- verification: automated
- hooks: []

## Objective

Finalize S02 implement/review prompts for:

1. Library-backed OpenAPI additions required by `UX-IA.md` (if none → `NO-GAPS.md`).
2. **P32-PORT** — **always required**, even when API is `NO-GAPS.md`. Lock chosen serve/UX measure(s) from S00 RESEARCH + OPEN-PORT-MULTI (minimum: friendly `EADDRINUSE` / in-use message + `--addr` guidance).

**No product code in this planner row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- `../scope-01-ux-ia/UX-IA.md`
- `../scope-00-research/RESEARCH.md`
- Live: `internal/httpapi/`, `cmd/trace/serve.go`, `api/openapi.yaml`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Law 19 | Handlers/CLI call library only — no business-logic fork |
| API gaps (library) | Per UX-IA: **none** for inspector map → expect `NO-GAPS.md` for library/OpenAPI |
| Client glue | UX-IA requires **`getImpact` `ops.ts` wrapper** (OpenAPI `/v1/impact` exists) — ship in S02-01 as adapter glue, not a new core op |
| Do not invent | `/v1/path`, requiring `listChanges`/`listRegressions` without new IA proof |
| P32-PORT | **Always ship** in S02-01 — RESEARCH prefer **#1** (friendly `EADDRINUSE` + `--addr` examples); #2 auto-port optional/defer |
| Bind defaults | Keep loopback `127.0.0.1:7432`; no `0.0.0.0` without `--allow-remote` |
| Defer hatch | `NO-PORT-CHANGE.md` only with written reason (discouraged) |

## Must answer (handoff to 01) — LOCKED 2026-08-21

1. **Library API = `NO-GAPS.md`?** Yes — UX-IA + live OpenAPI/handlers (why/context/impact/graph/search/reviews) confirm no new library ops. S02-01 authors `NO-GAPS.md` with evidence.
2. **`getImpact` wrapper shape + tests?** Ship `getImpact(taskId: string, opt?: TokenOpt)` → `GET /v1/impact?task_id=` matching `getWhy`/`getContext`; response `Record<string, unknown>`. No vitest suite required if none exists; TypeScript pattern-match + S02-02 review. Go tests own P32-PORT.
3. **P32-PORT option(s)?** **#1 min only** (friendly in-use + `--addr` examples). **#2 deferred** (not this row).
4. **Test plan for port-in-use?** Go test: occupy addr → second serve/listen fails with stderr containing in-use guidance + `--addr` hint (`cmd/trace/serve_test.go` and/or `internal/httpapi`).
5. **Help/stderr copy?** Conflict path must name in-use + suggest `trace serve -C <other> --addr 127.0.0.1:7433` (wording flexible). Optional one-line help multi-project note; S05 owns full quickstart docs.

## Planner gate

- [x] `01-implement.md` thick enough for fresh subagent (API + P32-PORT)
- [x] `02-review.md` checklist includes P32-PORT always
- [x] `SCOPE-TODOS.md` states S02 is not “API-only”

## Exit criteria

- [x] Implementer locked; next **P32-S02-01**

## Todo updates

Status + notes on **P32-S02-00** only.

## Next

`P32-S02-01`
