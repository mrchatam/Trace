# VERIFY-NOTES — P34-S05-01

**Date:** 2026-08-21
**Git SHA:** unknown (workspace has no `.git` in this environment)
**Overall:** PASS
**Evidence:** `experiments/runs/2026-08-21-p34-s05-01-verify/evidence/`
**Precondition:** P34-S04-02 PASS; S00–S04 done; PLAN T9 floor

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS | `00-run-metadata.txt`; EVID dir created |
| 1 Static embed T1–T3 + T10 | PASS | Focused static tests ok; `embeddist/index.html` has `#root` + `/assets/`; **no** `Embedded GUI stub` (`01-static-embed.txt`, `01b-embeddist-markers.txt`) |
| 2 httpapi auto-port T4/T7/T11 | PASS | `TestListenAutoPort_*` ok (`02-httpapi-auto-port.txt`) |
| 3 Concurrent + pin T5/T6 | PASS | Concurrent defaults + explicit DefaultAddr no-hop ok (`03-cmd-concurrent-pin.txt`) |
| 4 Packages + help | PASS | `go test ./internal/httpapi/ ./cmd/trace/ -p 1` ok; `gui`/`serve` help mention auto-port `7432`–`7441` + `--addr` pin; loopback default (`04-*.txt`) |
| 5 Docs T8 | PASS | Quickstart primary = `trace gui` + binary embed + `.trace/` only; multi-project hop; pin-strict; forbidden greps empty; `web/README` contributor-labeled (`05-docs-t8.txt`) |
| 6 Live consumer-temp | LIVE PASS | `trace init` → `.trace/` only (no `web/`); `serve` GET `/` = real SPA (`#root` + `/assets/`); concurrent second root → `:7433`; no SPA under `.trace/` (`06*.txt`) |
| 7 Residuals + aggregate | listed | No fail criteria tripped |

## L1–L4 + Docs

- [x] **L1** consumer `.trace/` only — live temp after `init` has `.trace/` (db+lock), no project `web/`; GET `/` served from binary embed (Block 6)
- [x] **L2** real SPA from binary (not stub) — shipped `embeddist/index.html` + live GET `/` have `id="root"` and `/assets/` module; T10 stub phrase **absent** (Blocks 1 + 6)
- [x] **L3** auto-port concurrent + correct URL — Go T4/T5/T11 + live `:7432` then `:7433` for second process (Blocks 2–3 + 6)
- [x] **L4** one process = one root — concurrent smoke used two `-C` temps × two ports; docs/help state one root per process (Blocks 5–6)
- [x] **Docs T8** embed + auto-port — quickstart / help / `web/README` / embeddist; no consumer two-artifact primary; no “no auto free-port” (Block 5)

## Aggregate (S00–S04)

- **S00 RESEARCH:** lean embed=A; UA-incr auto-port; L3 supersedes P33 reject — `scopes/scope-00-research/RESEARCH.md` + board P34-S00-02 PASS
- **S01 PLAN:** T1–T11 + board S02→S05; StaticDir opportunistic; Visit/`flag.Changed` intent — `scopes/scope-01-plan/PLAN.md` + P34-S01-01
- **S02 embed:** real SPA in `embeddist`; T1–T3; README consumer `.trace/` only — board P34-S02-01/02 PASS
- **S03 auto-port:** shared hop in httpapi; T4–T7/T11 + T5 concurrent CLI; T6 explicit busy no hop — board P34-S03-01/02 PASS
- **S04 docs:** T8 greps PASS; `gui-quickstart` primary embed + `7432`–`7441` + `--addr` pin; contributor `web/README` — board P34-S04-01/02 PASS

## Residuals (non-blocking)

- Cosmetic help / wording nits — T8 greps clean; positive story present
- Contributor Trace-checkout `web/` DX still documented — labeled **contributor** in `web/README.md`
- Default StaticDir path string still `<root>/web/dist` — resolution disk→embed→placeholder; consumers rarely need `--static-dir`
- Optional CI workflow for `embed-gui` not invented — out of phase (PLAN deferred)
- Explore UI / craft redesign — Phase 33 closed; out of phase
- Hosted SaaS / brew/deb — out of scope
- Workspace has no `.git` → VERIFY SHA recorded as `unknown` (does not affect product gates)

## Failures (if any)

- None

## DR-HANDOFF

remains **OPEN** — close owner **P34-S05-02**

## Next

P34-S05-02
