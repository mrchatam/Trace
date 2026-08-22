# Scope 01 — board map

**S01 plan** — docs-only `PLAN.md`. Serial: **P34-S01-00 → P34-S01-01**. No separate S01 review row (S01-00 gate + S02-00 may spot-check PLAN). Do **not** start S02 until S01-01 done. **No product code** on either S01 row.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 596 | P34-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken `01-plan.md` + this file |
| 597 | P34-S01-01 | [01-plan.md](01-plan.md) | Implementer | Author `PLAN.md` only |

## Inputs (verified for S01-00)

| Input | Status |
|-------|--------|
| S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md) | **PASS** (P34-S00-02) — embed **(A)**; StaticDir opportunistic; auto-port UA-incr shared |
| [`DESIGN-LOCKS.md`](../../DESIGN-LOCKS.md) | L1–L4 locked; L3 supersedes P32/P33 no-auto-port for default only |
| Live baseline | stub-only `embeddist/`; disk→embed→placeholder; no Trace-root Makefile/`.github`; `DefaultAddr` `127.0.0.1:7432`; serve pre-listen print; `flag` `--addr` default string |

## Locked leans (S01-00 → S01-01 PLAN)

1. **Embed:** `web` `npm ci && npm run build` → sync into `internal/httpapi/embeddist` (keep short README) → `go build`; primary entrypoint `scripts/embed-gui.sh` + `//go:generate` (optional `make embed-gui`); S05 stub-fail via missing `#root`/`/assets/` or stub phrase.
2. **StaticDir:** keep default candidate `<root>/web/dist`; no Trace-module probe; refuse == root; consumer needs no `web/`.
3. **Auto-port:** shared `gui`+`serve`; ports `7432`–`7441` (max 10); host `127.0.0.1`; `flag.Changed` strict; serve print post-bind.
4. **Tests:** matrix seeds T1–T9 in `01-plan.md` (S02/S03/S04/S05).
5. **Docs:** RESEARCH audit table owners — do not invent a thinner list.

## Outputs for later scopes

| Downstream | Consumes from PLAN |
|------------|-------------------|
| S02 | Embed pipeline cmds + StaticDir table + SPA-vs-stub markers + T1–T3 |
| S03 | Auto-port / `--addr` table + shared helper + T4–T7 |
| S04 | Docs touch list + T8 |
| S05 | VERIFY floor T9 (T1 + concurrent port + docs) |

## Out of this scope

- Implementing embed or auto-port; creating Makefile/scripts (S02); rewriting quickstart bodies (S04); VERIFY evidence (S05).
- Authoring `PLAN.md` on **P34-S01-00** (that is **P34-S01-01** only).
