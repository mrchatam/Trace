# P29-S06-01 — Production hardening implementer

## Metadata
- id: P29-S06-01
- todo_ids: [P29-S06-01]
- role: implementer
- skills: [security-and-hardening, shipping-and-launch, documentation-and-adrs]
- verification: mixed
- hooks: []

## Objective

Make local GUI+API **production-ready** for opt-in `trace serve`: harden security defaults, packaging, user docs, AGENTS carve-out; document cloud extension path (**design only**).

**Do not** start P29-S06-02. **Do not** rewrite S04/S05 feature set. **Do not** ship multi-tenant SaaS / OAuth / public-internet defaults.

## References

- [00-PLANNER.md](00-PLANNER.md) — **final locked defaults**
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)
- [RESEARCH.md](../scope-00-research/RESEARCH.md) § Law carve-out draft (paste into AGENTS / project-rules)
- Live: `internal/httpapi/`, `cmd/trace/serve.go`, `web/`
- Peer hardening pattern (read-only inspiration): codebase-memory `security-ui.sh` / embed — copy **checklist mindset**, not their `/rpc` transport

## Session start

Follow agent-loop-protocol Session start. Human locks and S01–S05 locks are **settled** — do not re-grill bind/auth/CORS deny/`*`, Law 19, or GUI IA.

## Locked defaults

| Item | Value |
|------|-------|
| Deploy | Local-first production use only |
| Default listen | `127.0.0.1:7432` (unchanged) |
| Remote | `--allow-remote` + bearer required (unchanged) |
| CORS default | Deny / omit allow headers; **never** `Access-Control-Allow-Origin: *` |
| CORS Vite DX | Optional `--cors-origin <exact URL>` (e.g. `http://127.0.0.1:5173`). Reflect **only** that exact Origin; reject others; no `*`, no prefix/regex |
| CSP | On static/SPA responses: at least `default-src 'self'; base-uri 'self'; frame-ancestors 'none'; object-src 'none'`. Allow `'unsafe-inline'` for style **only if** current SPA needs it; prefer tightening. `connect-src 'self'` (and the configured `--cors-origin` host only if that flag is set and docs say so — default same-origin) |
| StaticDir default | `<root>/web/dist` |
| StaticDir refuse | If resolved `--static-dir` **equals** project root → refuse to serve (exit ≠ 0) with message naming the footgun (`.trace/` exposure). Help text + quickstart must warn |
| Packaging | **Supported path:** build `web` → `web/dist`, then `trace serve`. **Embed:** `go:embed` of `web/dist` used only when disk `index.html` is missing; when disk dist exists, disk wins |
| Seed strict | Keep HTTP 501 for `strict`/`task_id` on export; document in quickstart — CLI for honesty/gate export |
| Errors | Extend `mapDomainErr`: validation-shaped errors (incl. `*domain.ErrValidation` and messages containing `must be UUID`) → **400** `VALIDATION_ERROR` |
| Promote (low) | Discoveries promote path: if optional `createTransition` returns deny/4xx envelope, show message (task+link success may still stand) — no silent swallow |
| Docs paths | Prefer `docs/gui-quickstart.md` (or `docs/serve-quickstart.md`) + short README pointer; apply RESEARCH carve-out to `AGENTS.md` Hard boundaries + `docs/rules/project-rules.md` settled-stack **new row** (keep historical P0 surface row) |
| Cloud | New appendix file or ADR section — design only; no code for tenants/OAuth/TLS termination product |
| Graph / lists | No new unbounded graph endpoint. Confirm `GET /v1/graph` still requires `center`+`max_nodes`. Where list handlers lack caps (e.g. `listTasks`), add sensible `limit` default+max **or** document intentional “project-local full list” with a hard ceiling — prefer ceiling ≤ 10k with 400 when exceeded only if library supports paging; otherwise document + leave library as-is and cap only search/changes/regressions already capped |

## Work breakdown (ordered)

### A — Security audit + code

1. **Bind/token regression:** Ensure existing tests still PASS (`RefuseRemote`, bearer off-loopback, loopback optional token). Add/adjust only if gaps found.
2. **CORS exact-origin (optional flag):** Wire `--cors-origin` on `trace serve` → `httpapi.Options`. When set and request `Origin` equals it exactly → set `Access-Control-Allow-Origin` to that value (+ `Vary: Origin`); never `*`. Preflight may echo allow-methods/headers needed for `/v1` + Authorization. When unset → current deny behavior.
3. **CSP on static:** In `static.go` (or middleware for non-`/v1` GET), set CSP headers before writing files. Do not weaken API JSON handlers with HTML CSP incorrectly.
4. **`--static-dir` footgun:** Refuse StaticDir == project root (after Abs). Document in `serve` help.
5. **`mapDomainErr`:** Map UUID validation failures to `VALIDATION_ERROR`. Add httptest covering invalid loop `task_id` / UUID string → 400 not 500.
6. **Promote honesty (low):** In `web/src/screens/Discoveries.tsx` (or shared toast), surface transition deny; keep promote success for create+link.

### B — Packaging

7. **Two-artifact docs:** Quickstart steps: `cd web && npm ci && npm run build` → `go build -o bin/trace ./cmd/trace` → `./bin/trace serve` → open `http://127.0.0.1:7432`.
8. **Embed fallback (required attempt):** Add embed of `web/dist` (build-tag or always-on with committed placeholder/`index` stub only if needed so `go test` works without full SPA). Runtime: prefer disk `StaticDir`; else embedded FS; else existing placeholder. Document release: populate `web/dist` before release `go build` when shipping single-binary with UI.
9. If embed proves blocked by empty-dist CI, ship **documented two-artifact only** and note “embed deferred” in Notes with reason — do not invent a second SPA root.

### C — Docs + law text

10. **Quickstart** under `docs/` covering: loopback default, `--allow-remote`+token, `--static-dir` footgun, Vite proxy vs `--cors-origin`, seed 501 honesty, graph budgets.
11. **AGENTS.md:** Apply RESEARCH carve-out bullets (opt-in HTTP, Law 19, still-forbidden, cloud path, FR-P28-X1 historical note). Align “Current focus” next-runnable with board after this phase progresses (do not claim S05-01 still next).
12. **project-rules.md:** Add settled-stack row: `Surface (post–Phase 29) | Library + CLI + MCP (stdio) + **opt-in** trace serve HTTP/GUI (loopback default)` — keep P0 row intact.
13. **Cloud appendix (design only):** e.g. `docs/phases/phase-29-http-api-browser-gui/CLOUD-APPENDIX.md` or ADR subsection — authn/z, tenancy header reserved, TLS at edge, same OpenAPI; **explicit non-goals** for Phase 29.

### D — Perf / Law 6–7 check

14. Grep handlers: no full-graph dump; graph requires budget. Confirm search/changes limits. Note any unbounded list in Notes for S07 if intentionally deferred.

## Test / evidence (Notes must cite)

```bash
cd /home/ali/Desktop/Trace
go test ./internal/httpapi/...
go test ./cmd/trace/ -run Serve
# CORS * still forbidden; with --cors-origin exact reflect works; wrong Origin gets no *
# mapDomainErr: invalid UUID / rl… style → 400 VALIDATION_ERROR
cd web && npm run build
# optional: npm run test:e2e if SPA promote honesty touched
./bin/trace serve --help   # documents cors-origin / static-dir footgun text if added
```

Security checklist (paste evidence into board Notes):

- [ ] Default bind loopback
- [ ] Remote gated + bearer
- [ ] No CORS `*`
- [ ] Exact-origin reflect only when flagged
- [ ] CSP on static
- [ ] StaticDir ≠ project root enforced
- [ ] Tokens not logged
- [ ] No `/rpc` MCP browser transport

## Out of scope

- Hosted multi-tenant product, OAuth, billing, always-on daemon
- Rewriting OpenAPI wave matrix / S05 FEATURE-MATRIX features
- Promoting HTTP seed `strict` without shared library helper
- Changing default bind away from loopback

## Exit criteria

- [ ] Security defaults audited; tests cover bind/CORS*/mapDomainErr UUID→VALIDATION where applicable
- [ ] Packaging path usable from clean clone (quickstart); embed fallback **or** explicit two-artifact-only note
- [ ] AGENTS + project-rules carve-out applied from RESEARCH
- [ ] Cloud appendix present; no hosted deploy claimed
- [ ] `--static-dir` footgun documented (+ refuse root)
- [ ] Board Notes on **P29-S06-01** only with evidence commands

## Todo updates

Status + notes on **P29-S06-01** only.

## Next

**P29-S06-02**
