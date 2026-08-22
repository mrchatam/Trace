# P32-S02-02 — API gaps + P32-PORT review

## Metadata
- id: P32-S02-02
- todo_ids: [P32-S02-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening, api-and-interface-design]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S02-01 deliverables. **P32-PORT is a hard checklist item** even when library API is `NO-GAPS.md`. Confirm client `getImpact` glue, Law 19 adapters-only, loopback safety, and that S02 did **not** invent `/v1/path` or ship deferred #2 auto-port as if it were required. Thicken upcoming S03 only if blast radius changed. Fresh context — do not share implementer session.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [00-PLANNER.md](00-PLANNER.md) / [01-implement.md](01-implement.md) locked defaults
- S01 [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md)
- S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md) — P32-PORT #1
- Expected artifacts: [`NO-GAPS.md`](NO-GAPS.md), `web/src/api/ops.ts`, `cmd/trace/serve.go` (+ helper/tests), board Notes on P32-S02-01

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Checklist

### Library / OpenAPI
- [ ] `NO-GAPS.md` exists and cites UX-IA + live evidence (or gaps shipped are library-backed only — unexpected per planner)
- [ ] **No** new `/v1/path` (or other invented core ops)
- [ ] **No** required `listChanges` / `listRegressions` client/API work without IA proof
- [ ] OpenAPI unchanged for impact **or** any change is justified and still Law 19

### Client glue
- [ ] `getImpact` present in `web/src/api/ops.ts`
- [ ] Signature matches planner lock: `getImpact(taskId: string, opt?: TokenOpt)` → `GET /v1/impact?task_id=…`
- [ ] Style consistent with `getWhy` / `getContext` (`apiFetch`, `TokenOpt`)
- [ ] **No** Graph/inspector UI wiring claimed as S02 (belongs S03) — glue-only is correct

### P32-PORT (hard — always)
- [ ] **#1 addressed:** friendly in-use / `EADDRINUSE` (or equivalent) detection
- [ ] Stderr (and/or help) guides user to distinct `--addr` with concrete example(s)
- [ ] Default remains `127.0.0.1:7432`; still fail-on-conflict (improved message ≠ auto-bind)
- [ ] Automated test evidence for conflict path **or** explicit residual with severity if missing
- [ ] **#2 auto-port not silently shipped** as the sole story (defer OK; must not replace #1)
- [ ] `NO-PORT-CHANGE.md` absent **or** has written reason (discouraged — treat as high if used without strong justification)

### Security / Law 19 / bind policy
- [ ] Loopback defaults safe — no `0.0.0.0` / public bind without `--allow-remote`
- [ ] Law 19: no SQL / domain business fork in `web/` or new handler logic beyond adapter error formatting
- [ ] Token / `--allow-remote` policy unchanged in spirit
- [ ] Help/stderr usable for multi-project local users (not opaque `serve: %v` only)

### Process
- [ ] Board Notes on P32-S02-01 cite files + commands
- [ ] Confidence medium or high; blocker/high → inline fix or spawn `02a`/`02b` below this row
- [ ] Upcoming S03 thickened only if S02 changed handoff (e.g. `getImpact` shape differs)

## Findings protocol

Per agent-loop-protocol reviewer loop: severity blocker | high | medium | low | nit. Blocker/high → small inline fix **or** spawn implement+review pair immediately below this row. Medium: prefer spawn unless trivial. Re-verify until no open blocker/high without pending follow-up.

## Exit criteria

- [ ] Checklist complete with evidence in Notes
- [ ] No open blocker/high without follow-up row
- [ ] Confidence **medium** or **high** (residuals listed if medium)
- [ ] P32-PORT explicitly ticked
- [ ] Next: **P32-S03-00**

## Todo updates

Status + notes on **P32-S02-02**; may thicken upcoming prompts only; may insert spawn rows immediately below if needed.

## Next

`P32-S03-00`
