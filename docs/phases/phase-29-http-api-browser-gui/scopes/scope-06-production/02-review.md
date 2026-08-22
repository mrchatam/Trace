# P29-S06-02 — Production review

## Metadata
- id: P29-S06-02
- todo_ids: [P29-S06-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening]
- verification: mixed
- hooks: []

## Objective

Independent review of production hardening vs loopback defaults, packaging, docs, AGENTS carve-out, and **cloud-appendix-only** rule. Fresh subagent. Small inline fixes OK; structural gaps spawn `P29-S06-02a` / `02b`.

## Session start

Follow agent-loop-protocol Session start. Do not share the implementer session. Do not start S07 until this row exits.

## References

- [00-PLANNER.md](00-PLANNER.md) locked defaults
- [01-implement.md](01-implement.md) work breakdown + security checklist
- [`docs/adr/ADR-HTTP-API-GUI.md`](../../../../../adr/ADR-HTTP-API-GUI.md)
- [`RESEARCH.md`](../scope-00-research/RESEARCH.md) carve-out draft
- Implementer board Notes on **P29-S06-01**

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -d internal/httpapi && test -d web
grep -q 'case "serve"' cmd/trace/root.go
go test ./internal/httpapi/...
# quickstart + carve-out present:
test -f docs/gui-quickstart.md -o -f docs/serve-quickstart.md
grep -q 'trace serve' AGENTS.md
```

If hardening/docs missing with no Notes reason, **fail** or spawn 02a — do not re-implement the whole scope in review.

## Checklist

### Security (human locks)

- [ ] Default listen remains loopback (`127.0.0.1:7432` or documented equivalent)
- [ ] Non-loopback / `0.0.0.0` refused without `--allow-remote`
- [ ] Bearer required off-loopback; 401 `UNAUTHORIZED`
- [ ] CORS: **no** `Access-Control-Allow-Origin: *` in default or flagged paths
- [ ] If `--cors-origin` exists: exact Origin only; wrong Origin gets no ACAO; never wildcard
- [ ] CSP present on static/SPA responses (`default-src` / `frame-ancestors` sensible)
- [ ] `--static-dir` = project root refused **or** equivalent hard guard + docs
- [ ] Tokens not logged; error bodies do not leak secrets
- [ ] No MCP `/rpc` browser transport

### Residuals from S04/S05

- [ ] `mapDomainErr`: UUID / “must be UUID” → **400** `VALIDATION_ERROR` (not 500 `INTERNAL_ERROR`) — test or live evidence
- [ ] Promote `createTransition` deny: honesty surfaced **or** explicit low residual noted for S07 (not silent forever without Notes)
- [ ] HTTP seed `strict`/`task_id` still **501** (CLI honesty retained); documented

### Packaging / docs

- [ ] Quickstart works as written (clone → build SPA → serve → browser) or embed fallback documented
- [ ] Two-artifact vs embed behavior matches 00-PLANNER locks
- [ ] AGENTS.md carve-out matches RESEARCH intent (opt-in, loopback, Law 19, still-forbidden, cloud separate)
- [ ] project-rules settled-stack has post–Phase 29 surface row; P0 historical row intact
- [ ] Cloud appendix **design only** — no accidental SaaS/OAuth/hosted deploy code

### Law 6–7 / 19

- [ ] No unbounded graph dump endpoint; `center`+`max_nodes` still required
- [ ] Handlers still library-only (no SQL in `internal/httpapi`)
- [ ] List/search limits not regressively removed

### Tests

- [ ] `go test ./internal/httpapi/...` PASS
- [ ] `go test ./cmd/trace/ -run Serve` PASS (or equivalent serve tests)
- [ ] If SPA changed: `cd web && npm run build` PASS; e2e if promote path touched

## Findings

Classify: blocker | high | medium | low | nit. Cite files.

## Spawn policy

- blocker/high: inline fix if small; else insert **P29-S06-02a** (implement) + **P29-S06-02b** (review) immediately below this row
- medium: prefer spawn unless trivial
- Do not rewrite P29-S06-01 prompt; do not edit done S00–S05 history
- Upcoming only: thicken **S07** prompts if VERIFY needs new evidence hooks

## Exit criteria

- [ ] No open blocker/high without a pending follow-up row
- [ ] Confidence **medium** or **high** with evidence (commands, greps, doc paths)
- [ ] Board Notes on **P29-S06-02**; next **P29-S07-00**

## Next

**P29-S07-00**
