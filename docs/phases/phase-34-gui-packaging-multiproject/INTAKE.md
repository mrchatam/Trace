# INTAKE — Phase 34 packaging / multi-project GUI

**Human 2026-08-21** after Phase 33 dogfood.

## What happened

1. Trace repo already running `trace gui` on `127.0.0.1:7432`.
2. In another project (`feet seller telegram app`): `trace gui` → **address already in use** (manual `--addr` hint only; no auto free-port).
3. Opening `http://127.0.0.1:7432` showed the **embedded stub** (“Build the browser GUI…”), not the Phase 33 Explore SPA — because that bind was the first process, and/or a second attempt would look for **`<other-project>/web/dist`**, which does not exist.

## Root causes (light)

| Issue | Cause |
|-------|--------|
| Stub UI in consumer repo | Default `StaticDir` = `<projectRoot>/web/dist`. Consumer projects have **no** `web/`. Missing disk dist → **embed stub**, not Trace’s real SPA. |
| Port conflict | P32-PORT shipped friendly error only; **auto free-port deferred**. Multi-project requires manual `--addr`. |
| Directory hazard | Documented two-artifact path encourages `web/` beside app source — wrong for **consumer** repos. Trace product UI must not live in the user’s project tree. |

## Human locks (desired)

1. **Automatic multi-project / port** — `trace gui` in a second repo must Just Work (pick a free port, print + open the right URL). No mandatory manual `--addr` for the happy path.
2. **GUI assets from Trace product**, not from the served project — real SPA embedded (or beside the installed binary), never require consumer `web/`.
3. **Consumer footprint = `.trace/` only** for Trace state/data. No `web/`, no root `trace.db`, no Trace SPA dirs in the user’s repo.

## Out of scope reminders

- Hosted SaaS
- Putting SPA under consumer `.trace/` as a copy of `web/dist` (prefer binary embed; `.trace/` is for store/lock/token)
