# P29-S01-02 — Architecture review

## Metadata
- id: P29-S01-02
- todo_ids: [P29-S01-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Independent review of ADR + OpenAPI against Law 19, human locks, S00 `RESEARCH.md`, and S01-00 locked defaults. Fresh subagent — do not share implementer session.

## References

- [01-adr-and-openapi.md](01-adr-and-openapi.md)
- [00-PLANNER.md](00-PLANNER.md) locked defaults
- [../scope-00-research/RESEARCH.md](../scope-00-research/RESEARCH.md)
- [Phase 29 README](../../README.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §19
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start. Fresh subagent.

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f docs/adr/ADR-HTTP-API-GUI.md
test -f api/openapi.yaml
test ! -d internal/httpapi
test ! -d web
! grep -q 'case "serve"' cmd/trace/root.go
```

## Checklist

### Human locks + S01-00

- [ ] Opt-in `trace serve`; default bind `127.0.0.1`; open bind refused without `--allow-remote` (or equivalent named in ADR)
- [ ] Auth: loopback-trust; bearer required for non-loopback — matches ADR + OpenAPI security notes
- [ ] Package named `internal/httpapi`; prefix `/v1`
- [ ] ADR path `docs/adr/ADR-HTTP-API-GUI.md`; OpenAPI `api/openapi.yaml`
- [ ] CORS: no `*`; default deny / same-origin documented
- [ ] No MCP `/rpc` (or tools/call) as browser transport
- [ ] Static GUI: disk `web/dist` first; embed deferred to S06
- [ ] Cloud: extension hooks only — no tenancy/OAuth implementation claimed

### Law 19 + progressive context

- [ ] Handlers/UI described as adapters → canonical library; no second SoT
- [ ] No full-graph dump default endpoint; `/v1/graph` budgeted (`max_nodes` / center)
- [ ] Seed routes are status/summary oriented (not silent full dump)

### OpenAPI coverage

- [ ] P0 families present: health, version, project, tasks, loop, entities, links, transitions, context, why, search, seed, graph
- [ ] P1 / defer families listed explicitly (reviews, plans, capability, impact, changes, regressions, agents, index, auth token, ops deferred)
- [ ] Shared error schema; consistent `/v1` paths
- [ ] `x-trace-wave` (or equivalent) distinguishes p0/p1/defer — or deferred section is unambiguous
- [ ] OpenAPI is valid YAML and parseable as OpenAPI 3.x (spot-check `openapi:` / `paths:`)

### Process

- [ ] No product `internal/httpapi` / `web/` / `serve` introduced by implementer
- [ ] Paths recorded; confidence **medium+** with residuals listed (never silent)
- [ ] Findings severity tagged: blocker | high | medium | low | nit

## Spawn policy

- **blocker/high:** small inline fix **or** spawn `P29-S01-02a` (implement) + `P29-S01-02b` (review) immediately below this board row; write full prompts under this scope folder.
- **medium:** prefer spawn unless trivial fix.
- **low/nit:** note residual; spawn only if it would block S02/S03.

## Exit criteria

- [ ] No open blocker/high without pending follow-up
- [ ] Confidence medium or high with evidence in Notes
- [ ] Next **P29-S02-00** (unless spawn rows inserted first)

## Todo updates

Status + notes on **P29-S01-02** only (plus spawned rows if any). May thicken **upcoming** S02 prompts if contract gaps affect UX IA — do not rewrite S00 `done` history.

## Next

**P29-S02-00** (or spawned 02a/02b first)
