# P29-S00-02 — Research review

## Metadata
- id: P29-S00-02
- todo_ids: [P29-S00-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Independent review of S00 `RESEARCH.md` against Phase 29 human locks and live repo. Spawn remediation if blocker/high gaps; otherwise pass with confidence. **Fresh subagent — not the S00-01 author.**

## References

- [01-peer-and-surface-research.md](01-peer-and-surface-research.md)
- [00-PLANNER.md](00-PLANNER.md)
- [RESEARCH.md](RESEARCH.md) (produced by S00-01)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 29 README](../../README.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §19
- [docs/TODO.md](../../../../TODO.md) Later developments

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Preflight (reviewer)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-29-http-api-browser-gui/scopes/scope-00-research/RESEARCH.md
test ! -d internal/httpapi
test ! -d web
! grep -q 'case "serve"' cmd/trace/root.go
```

## Verify checklist

### Completeness
- [ ] `RESEARCH.md` has all template sections (exec summary, peer matrix, surface map, API families, stack, law carve-out, risks)
- [ ] Peer matrix ≥3 rows and covers locked paths from `01`:
  - Understand-Anything `packages/dashboard` (+ skill noted or skimmed)
  - `codebase-memory-mcp/graph-ui`
  - `agentrq/frontend`
- [ ] Trace surface map covers CLI + MCP parity **intent** (not necessarily every admin cmd as P0)
- [ ] API resource families listed for S01 OpenAPI (actionable, not vague)
- [ ] Stack recommendation present with trade-offs; if not TS+Vite+React, evidence bar met
- [ ] Law carve-out draft present: opt-in `trace serve`, loopback default, Law 19, cloud = separate hosted product / same OpenAPI

### Human locks / safety
- [ ] No open-bind / public-internet default proposed as ship target
- [ ] No SaaS/tenancy/OAuth deploy plan as Phase 29 ship target
- [ ] No second SoT / business-logic fork in the browser
- [ ] No “point local MCP at the internet”

### Repo hygiene
- [ ] No premature product code (`internal/httpapi`, `web/`, `case "serve"`)
- [ ] S00-01 Notes claim matches files on disk
- [ ] Open decisions for S01–S02 are listed (not silently assumed locked)

## Findings format

Tag each finding: `blocker` | `high` | `medium` | `low` | `nit`. Cite section/path.

## Spawn policy

| Severity | Action |
|----------|--------|
| blocker / high | Insert **`P29-S00-02a`** (implement fix to RESEARCH) + **`P29-S00-02b`** (review) **immediately below** this board row; write full prompts under `scopes/scope-00-research/` |
| medium | Prefer spawn unless a trivial doc fix (≤ few lines) you can apply inline |
| low / nit | Notes only or inline fix |

Confidence must be **high**, or **medium** with **explicit residual risks** listed — never silent.

Do **not** start S01 architecture until this review (or spawn chain) is `done`.

## Exit criteria

- [ ] Findings severity-tagged; no open blocker/high without pending follow-up row
- [ ] Board Notes cite evidence paths (`RESEARCH.md` sections + preflight)
- [ ] Next runnable **P29-S01-00** (or spawn pair still pending)

## Todo updates

Status + notes on **P29-S00-02**; may spawn upcoming rows + thicken upcoming prompts only.

## Next

**P29-S01-00** (after pass) or **P29-S00-02a** (if spawned)
