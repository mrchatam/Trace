# P32-S05-02 — Polish review

## Metadata
- id: P32-S05-02
- todo_ids: [P32-S05-02]
- role: reviewer
- skills: [code-review-and-quality, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S05 polish + docs vs [`01-implement.md`](01-implement.md) locks, [`OPEN-PORT-MULTI.md`](../../OPEN-PORT-MULTI.md), and live S02 serve behavior. **Block** if multi-project port docs are missing, wrong (e.g. claim auto-port), or contradict fail-on-conflict / loopback defaults. **Block** if “polish” silently added explorer features or regresses Laws 6–7/19. Fresh subagent — do not share implementer session.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [01-implement.md](01-implement.md) — doc + residual locks
- Live: `docs/gui-quickstart.md`, `web/README.md`, `internal/httpapi/addr_in_use.go`, `cmd/trace/serve.go`, `cmd/trace/help.go`
- S04 craft (for UI doc honesty): `.graph-shell`, `--graph-canvas-height`, `PacketView`, motions + reduced-motion

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Checklist

### P32-PORT / docs (hard)

- [ ] `gui-quickstart.md` has an explicit multi-project / port section (or equivalent headed content)
- [ ] Documents default `127.0.0.1:7432` and **distinct `--addr`** for a second project (example `127.0.0.1:7433` or equivalent)
- [ ] Matches S02 #1: fail-on-conflict + friendly in-use guidance — **does not** claim auto free-port / `:0` (#2 deferred)
- [ ] Example command shape includes `-C` / root + `--addr` (one serve = one project root)
- [ ] Security defaults still loopback-first; remote still explicit (`--allow-remote` + token) — no public-bind softening in docs
- [ ] `web/README.md` not left implying only one global port forever without a multi-project hint (or Notes say “confirmed OK”)
- [ ] OPEN-PORT #3/#4 narrative addressed via docs (helper script **not** required)

### Explorer / craft honesty (if UI mentioned or screenshots)

- [ ] Hero story is Explore `/` graph-home (canvas-first + inspector), **not** Overview CRUD / Phase 29 ops dual-card
- [ ] Any craft callouts match live S04 (shell / PacketView / calm nodes / reduced-motion) — no invented Three.js or new brand

### Residuals + scope discipline

- [ ] S04 residuals (chrome box-shadow nit; keyboard-via-list) closed **or** explicitly deferred with reason
- [ ] No new inspector sections, routes, API clients, or `/v1/path`
- [ ] Depth/IA invariants untouched (order, select≠expand, budgets 50/100, `getImpact`)
- [ ] Law 19: no docs/code implying browser business-logic SoT or always-on daemon

### Verify evidence

- [ ] Re-read live `FormatAddrInUseMessage` / serve help vs docs wording
- [ ] If implementer touched product CSS/TS: re-run `cd web && npm run build` + `npm run test:e2e -- e2e/s03-depth.spec.ts e2e/s05-gates.spec.ts` and cite PASS
- [ ] Docs-only change: confirm no unintended product diffs beyond stated polish

## Findings

Severity: blocker | high | medium | low | nit.

- blocker/high: inline fix **or** spawn `P32-S05-02a` / `02b` immediately below this row
- medium: prefer spawn unless trivial
- Do **not** rewrite `done` S02–S04 history
- Gaps for VERIFY (S06) → note forward (P32-PORT tick still required on VERIFY)

## Exit criteria

- [ ] No open blocker/high without pending follow-up spawn
- [ ] Confidence medium or high with evidence (port docs vs live serve strings)
- [ ] If PASS: lightly thicken S06 prompts only if VERIFY needs port-doc evidence path (optional)
- [ ] Next: **P32-S06-00** (do not start VERIFY implement)

## Todo updates

Status + notes on **P32-S05-02**.

## Next

`P32-S06-00`
