# P33-S05-00 — Scope planner (polish + docs)

## Metadata
- id: P33-S05-00
- todo_ids: [P33-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Finalize S05: flip user-facing primary launch story to **`trace gui`**, update `gui-quickstart` / README / AGENTS as needed, and clear residual bugs from S02–S04. Keep `serve` documented as secondary/scripting.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [docs/gui-quickstart.md](../../../../gui-quickstart.md)
- [AGENTS.md](../../../../../AGENTS.md)
- README / CONTRIBUTING as relevant
- S02 REVIEW (addr-in-use dual-word) · S04 REVIEW (no retokenize; EmptyState; optional canvas shot)

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults

| Item | Value |
|------|-------|
| Primary story | `trace gui` (+ PATH install via `go install …/cmd/trace@…`) |
| Secondary | `trace serve` / `./bin/trace serve` for scripting / no-browser / CI |
| Do not | Claim hosted SaaS; change loopback defaults; invent auto-port; conflate PATH with `trace install` |
| Craft (S04 PASS) | Forest-moss + kind chroma strip already landed — **do not** reinvent palette / tokens in polish |
| Product code | Docs-first; residual UI/CLI copy only as listed — **no** compose/route/budget/token rewrite |

## Must answer (handoff to 01) — LOCKED

1. **Doc file list (required):**
   - [`docs/gui-quickstart.md`](../../../../gui-quickstart.md) — rewrite **lead** (title + first walkthrough) to `trace gui`; keep PATH §; demote two-artifact `./bin/trace serve` to secondary; multi-project / Related mention `gui` (serve OK as scripting twin).
   - [`README.md`](../../../../../README.md) — `gui-quickstart` blurb currently says opt-in `trace serve` → flip to `trace gui` (+ serve secondary OK).
   - [`web/README.md`](../../../../../web/README.md) — still “Operator SPA for `trace serve`” + serve-only examples → primary `trace gui`.
   - [`AGENTS.md`](../../../../../AGENTS.md) — Current focus still says “Quickstart still documents `serve` until P33 docs flip”; update after flip; refresh orchestrator **Next runnable** if still stale.
   - Optional: [`docs/TODO.md`](../../../../TODO.md) orchestrator paste `Next runnable` (still cites old row) — align if touching AGENTS.

2. **Residual bug list:**
   | Residual | Live state | S05 disposition |
   |----------|------------|-----------------|
   | EmptyState Tasks CTA | `Graph.tsx` no-seeds already has **inline** `Link` “Open Tasks” inside `EmptyState`; footer `Tasks · Overview` remains | **Verify** inline CTA reads primary; optional visual polish only — do **not** invent new empty UX. If already adequate → Notes “closed / N/A” |
   | Explore canvas screenshot | S04 evidence `explore-{light,dark}.png` is **list-heavy** | **Optional:** one canvas-forward shot under `docs/` or scope evidence + link from quickstart “What you see”; defer to VERIFY Notes if no media pipeline |
   | Addr-in-use dual-word | `httpapi.FormatAddrInUseMessage` still `serve:` / `trace serve` only | **Optional low:** wording `gui`\|`serve` (prefix + hint); update `addr_in_use_test` + quickstart stderr block to match |
   | Craft literacy one-liner | Quickstart craft cues omit kind chroma strip | **Optional:** one sentence — kind literacy = chroma strip + labels (not color-only) |

3. **AGENTS.md / TODO index:** Flip Current-focus “serve until docs flip” sentence; set Next to **P33-S05-02** (after implement) or **P33-S06-00** only after review — implementer updates AGENTS on docs flip; do not rewrite done board history.

4. **Craft residual split:**
   - **Docs-only:** primary launch flip; craft literacy one-liner; optional canvas shot in docs.
   - **CSS/UI-only (tiny):** EmptyState CTA polish if needed — **no** `tokens.css` retokenize / no DESIGN rewrite.
   - **CLI copy-only:** optional `FormatAddrInUseMessage` dual-word.
   - **Out:** retokenize; overview compose; routes; budgets; canvas keyboard residual (accepted S03).

## Planner gate

- [x] `01-implement.md` + `02-review.md` thickened
- [x] `SCOPE-TODOS.md` accurate

## Exit criteria

- [x] Implementer locked; next **P33-S05-01**

## Todo updates

Status + notes on **P33-S05-00** only.

## Next

`P33-S05-01`
