# DESIGN-LOCKS — Phase 34

**Human-promoted 2026-08-21.** Clarified by P34-00 (same date) — light notes only; do not reopen L1–L4.

| Lock | Value |
|------|-------|
| Predecessor | Phase 33 closed — do not rewrite; spawn forward |
| **L1 — Consumer layout** | In a **user project**, Trace may create/use **only** `.trace/` (db, lock, token, …). **Forbidden** as required consumer layout: `web/`, project-root `trace.db`, copying SPA into the repo. |
| **L2 — GUI asset source** | Default: serve the **real Explore SPA from the Trace binary** (`go:embed` of built `web/dist` at Trace build/release time). Disk `<project>/web/dist` is **dev-only** (Trace checkout / contributor DX), never the consumer story. Stub is unacceptable as the shipped default when a full dist was built into the release. |
| **L3 — Auto port** | `trace gui` / `serve`: if default bind is in use, **automatically** try the next free loopback port (or equivalent), print the chosen URL, open browser to that URL. Keep `--addr` for explicit pin. Overturn P32 “no auto-port” defer for the happy path. |
| **L4 — One process = one project** | Still one store root per process (`-C`/cwd). Multi-project = multiple processes on distinct ports (auto). |
| Law 19 / loopback | Unchanged; no public defaults |

## Clarifications (P34-00)

| Topic | Clarification |
|-------|----------------|
| L2 Trace-checkout disk | Preferring disk `web/dist` when present is allowed **only** as contributor DX in the Trace product repo (or explicit `--static-dir`). Arbitrary `-C` consumer roots must **not** require or invent `web/`. |
| L2 stub | Stub / placeholder remains last-resort for empty embed (dev mistake). Release VERIFY must fail if shipped embed is still stub while a full SPA was supposed to be built in. |
| L2 reject | Do **not** copy SPA into consumer `.trace/` as the primary asset path (`.trace/` = store/lock/token). |
| L3 supersession | P32-PORT / P33 RESEARCH **rejected** UA-style auto-increment. **L3 overturns** that reject for **default** bind happy path only. Explicit `--addr` remains **strict**: fail if busy (friendly message may still mention auto behavior for default). |
| L3 scope | Prefer shared bind logic for `gui` and `serve` (same local HTTP). Loopback default unchanged. |
| L4 | Multi-project = N processes × N ports — not one daemon serving many roots. |

## Success sketch

From any Trace-initialized app repo (no `web/`):

```bash
trace gui          # first → :7432 + real SPA
trace gui          # second project → :7433 (or free) + real SPA + opens correct URL
```

Consumer tree still only gains/uses `.trace/`.
