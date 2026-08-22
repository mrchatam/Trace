# Phase 33 — GUI craft, Explore hook, `trace gui`

Close post–Phase 32 satisfaction gaps: **deliberate color/craft**, **Graphify-class Explore hook**, and practical **`trace gui`** + PATH install (not `./bin/trace serve` as the primary story).

**Complete** (2026-08-21). DR-HANDOFF **CLOSED**; successor **no successor**. Design SoT: [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) · [`INTAKE.md`](INTAKE.md). Verify: [`scopes/scope-06-verify/VERIFY-NOTES.md`](scopes/scope-06-verify/VERIFY-NOTES.md) · handoff review: [`scopes/scope-06-verify/REVIEW-NOTES.md`](scopes/scope-06-verify/REVIEW-NOTES.md).

## Human locks (do not reopen)

| Lock | Value |
|------|-------|
| Theme A | Color system + craft pass — skills **impeccable**, **ui-ux-pro-max**, **frontend-design** on S01/S03/S04 implement+review |
| Theme B | Explore (`/`) = interactive **project overview graph** hook (pan/zoom/click); Laws **6–7** budgeted / progressive — not unbounded dump |
| Theme C | **`trace gui`** (cwd/`-C`) = serve + print URL + **open browser**; PATH install documented (`go install` / symlink / make) |
| Stack | Keep `web/` React+Vite+`@xyflow/react` 2D; **Law 19** adapters only |
| Serve | Reuse `trace serve` + P32-PORT friendly bind errors; loopback default unchanged |
| Out | Hosted SaaS; always-on daemon; Three.js-first; silent delete of root `trace.db`; treating `trace install` (agents/MCP) as PATH install |

Full table: [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md).

## Live repo baseline (P33-00, 2026-08-21)

| Area | State |
|------|-------|
| Explore home | `App.tsx` index → `Graph.tsx`; nav label **Explore**; `/graph` redirects to `/` |
| Hook gap | Graph is still **center-first neighborhood** (`getGraph` + `DEFAULT_MAX=50` / `UI_CAP=100`) — empty until user picks a center. Not yet a Graphify-like **project overview** canvas on open |
| Craft | Calm gray-leaning CSS; kind/state color encoding weak vs INTAKE “bland” complaint |
| CLI launch | `trace serve` only — **no** `gui` subcommand; no auto-open browser |
| PATH | `trace install *` = agents/MCP/hooks — **not** “put `trace` on PATH”. Docs lead with `./bin/trace serve` ([`gui-quickstart.md`](../../gui-quickstart.md)) |
| Port story | P32-PORT shipped (friendly in-use + `--addr`); keep for `gui` reuse |
| Peers (local) | `similar projects/graphify/`, `similar projects/Understand-Anything/` (if present) |

## Scope index (serial)

```
S00 research (peers + Laws 6–7 overview + launch/PATH)
  → S01 design tokens + Explore-as-graph UX
  → S02 CLI `trace gui` + install/PATH (+ open browser)
  → S03 Explore interactive project graph hook
  → S04 colorize / craft full shell
  → S05 polish + docs (`trace gui` primary)
  → S06 VERIFY + DR-HANDOFF
```

| Scope | Title | Primary artifact | Board |
|-------|-------|------------------|-------|
| S00 | Peer + overview/PATH research | `RESEARCH.md` | P33-S00-00…02 |
| S01 | Color tokens + Explore IA | `DESIGN.md` / `UX-IA.md` | P33-S01-00…02 |
| S02 | `trace gui` + PATH install | CLI + tests + install docs/scripts | P33-S02-00…02 |
| S03 | Explore project-graph hook | Evolve `web/` Explore per S01 | P33-S03-00…02 |
| S04 | Colorize / craft shell | Tokens + craft across shell | P33-S04-00…02 |
| S05 | Polish + docs | `gui-quickstart` / README / AGENTS primary = `trace gui` | P33-S05-00…02 |
| S06 | VERIFY + handoff | `VERIFY-NOTES.md`; close DR-HANDOFF | P33-S06-00…02 |

**Ownership notes**

- **Overview-graph budget model:** S00 recommends → S01 locks IA → S03 ships UI (and only then any proven `/v1` gap — prefer reuse).
- **Launch:** S00 recommends CLI/PATH shape → **S02 ships** `trace gui` + PATH story → S05 docs flip primary → S06 VERIFY.
- **Color:** S01 tokens → S03 may use kind/state colors on nodes → **S04** owns full-shell colorize/bolder (do not mark S03 “craft done”).
- **Skills:** Mandatory on S01/S03/S04 implement **and** review prompts (list in metadata).

## In scope

- Opinionated color system + craft pass (not gray monochrome)
- Explore as interactive project overview graph (budgeted / progressive under Laws 6–7)
- `trace gui` + browser open + documented PATH install
- Docs / quickstart primary story flip away from `./bin/trace serve`

## Out of scope

- Hosted multi-tenant SaaS / OAuth / billing
- Always-on daemon; `0.0.0.0` as default bind
- Second SQLite or business-logic fork in `web/`
- Pointing local `trace-mcp` at the public internet
- 3D / Three.js as default canvas
- Rewriting Phases 00–32 `done` history
- Silent auto-delete of root `trace.db`
- Replacing `trace serve` (keep for scripting; demote in user-facing docs)

## Cross-scope blast radius

| If … | Then … |
|------|--------|
| S00 finds need for new overview API | Flag in RESEARCH; S01/S03 planners decide reuse vs thin OpenAPI gap — **no** unbounded dump endpoint |
| S02 PATH install conflicts with `trace install` naming | Document distinction; prefer `go install` / make symlink — do not overload agent-install |
| S03 lands before S04 | S03 may apply token placeholders; S04 owns shell-wide colorize |
| Docs still say `./bin/trace serve` after S02 | S05 owns primary-story flip; S02 may land minimal `--help` + install note |

## Handoff

[`DR-HANDOFF.md`](DR-HANDOFF.md) — **OPEN** until `P33-S06-02`. Successor lean: **no successor** unless residuals need a thin follow-on.
