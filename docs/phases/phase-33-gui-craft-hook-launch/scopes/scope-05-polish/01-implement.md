# P33-S05-01 — Polish + docs

## Metadata
- id: P33-S05-01
- todo_ids: [P33-S05-01]
- role: implementer
- skills: [incremental-implementation, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Update docs so primary launch is **`trace gui`** (PATH assumed). Fix residual bugs from S02–S04 Notes. Keep `serve` as secondary. **Do not** retokenize (S04 PASS).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — locked defaults + Must-answer
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — Theme C docs demote
- [`../scope-02-gui-launch/REVIEW.md`](../scope-02-gui-launch/REVIEW.md) — addr-in-use dual-word residual
- [`../scope-04-color-craft/REVIEW.md`](../scope-04-color-craft/REVIEW.md) — no retokenize; EmptyState; optional canvas shot
- Live docs: `docs/gui-quickstart.md`, `README.md`, `web/README.md`, `AGENTS.md`
- Live code (residuals only): `web/src/screens/Graph.tsx` (EmptyState), `internal/httpapi/addr_in_use.go` (+ test)

## Session start

Follow agent-loop-protocol Session start. Prefer planner locks; do not reopen S01–S04 craft/compose. Do **not** start P33-S05-02.

## Locked defaults (planner — do not re-debate)

| Item | Value |
|------|-------|
| Primary | `trace gui` (PATH via `CGO_ENABLED=1 go install github.com/mrchatam/Trace/cmd/trace@latest`) |
| Secondary | `./bin/trace serve` / `trace serve` for no-browser / scripting / CI |
| Already shipped (S02) | Help Build note + quickstart §Install CLI on PATH — **keep**; rewrite the **lead** walkthrough off `./bin/trace serve` |
| PATH ≠ install | Never teach `trace install` as PATH |
| Loopback / SaaS | Default `127.0.0.1:7432`; no hosted SaaS claims; no auto-port |
| Craft | **Do not** edit `tokens.css` kind/state palette or chroma-strip rules |
| Out | Compose/budgets/routes; Three.js; canvas arrow-roving; brew/deb packaging |

## Doc file checklist (required)

| File | Change |
|------|--------|
| `docs/gui-quickstart.md` | Title + lead path = `trace gui` (opens browser → Explore `/`). Keep PATH §. Demote two-artifact `go build` + `./bin/trace serve` to **secondary** (scripting / `--no-open` twin = `serve`). Multi-project examples: prefer `trace gui` (+ `--no-open` if documenting headless); keep `serve` as alt. Update stderr sample if addr-in-use copy changes. Related: `trace gui --help` (+ serve OK). Craft cues: optional one-liner kind = chroma strip + labels |
| `README.md` | Repo docs bullet for gui-quickstart → primary `trace gui` |
| `web/README.md` | Lead “Operator SPA for `trace gui`”; examples primary `trace gui`; note `serve` for headless/DX |
| `AGENTS.md` | Remove “Quickstart still documents serve until P33 docs flip”; Current focus reflects docs primary = `trace gui`; refresh orchestrator Next if stale |

Optional: `docs/TODO.md` orchestrator paste Next runnable — align if editing AGENTS.

## Residual work (ordered)

1. **Docs flip** (required) — table above.
2. **EmptyState CTA** — Live: no-seeds `EmptyState` already wraps inline `Link` “Open Tasks” (`Graph.tsx`). **Verify** it is the primary CTA; polish only if footer `Tasks · Overview` still reads as the only/main CTA. Do not redesign empty/error copy substance (S03).
3. **Addr-in-use** (optional low) — `FormatAddrInUseMessage`: dual-word so `gui` conflicts don’t say only `serve:` / “often trace serve”. Update `addr_in_use_test.go` + quickstart sample block to match. Keep fail-on-conflict + `--addr` example.
4. **Canvas shot** (optional) — If quickstart shows GUI visuals, add **one** Explore **canvas** screenshot (S04 PNGs are list-focused) under e.g. `docs/phases/…/scope-05-polish/evidence/` or reuse S04 with a new canvas crop; link from “What you see”. If awkward → **defer** explicitly in Notes for VERIFY.
5. **Craft literacy** (optional) — One quickstart sentence: kind literacy = left chroma strip + text labels (not color-only).

## Exit criteria

- [ ] Quickstart **leads** with `trace gui` (PATH assumed); `serve` demoted not deleted
- [ ] README + `web/README` + AGENTS aligned (no serve-as-primary user story)
- [ ] Residuals addressed **or** explicitly deferred to VERIFY Notes (EmptyState / canvas shot / addr-in-use)
- [ ] No palette/token rewrite undoing S04
- [ ] Board Notes with doc paths (+ code paths if residuals touched)
- [ ] Verify: re-read quickstart H1 + first code fence = `trace gui`; `go test ./internal/httpapi/ -run FormatAddrInUse` if message changed

## Minimal todos

- [ ] Docs flip (`gui-quickstart`, README, `web/README`, AGENTS)
- [ ] Residual fixes or defer Notes
- [ ] Board Notes

## Todo updates

Status + notes on **P33-S05-01** only.

## Next

`P33-S05-02`
