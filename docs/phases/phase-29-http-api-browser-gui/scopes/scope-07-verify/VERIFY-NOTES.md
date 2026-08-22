# Phase 29 VERIFY notes

- Date: 2026-08-21
- Git SHA: unknown (workspace has no `.git`; binary rebuilt from live tree)
- Evidence: experiments/runs/2026-08-21-p29-s07-01-verify/evidence/
- Verdict: **PASS**

## Blocks

| Block | Result | Evidence |
|-------|--------|----------|
| 0 Evidence dir | PASS | `99-run-metadata.txt` |
| 1 Build (go + web) | PASS | `01-go-build.txt` exit 0; `01-web-build.txt` exit 0; `web/dist/index.html` present |
| 2 Tests (httpapi, Serve, e2e?) | PASS (e2e partial) | `02-httpapi-test.txt` ok; `02-serve-test.txt` ok; e2e: API gates 2/2 PASS; UI 3 fail — Playwright chromium missing in sandbox (`02-e2e.txt`) — non-blocking per floor |
| 3 API smoke (health/tasks/loop) | PASS | `03-health.json` `{"ok":true}`; `03-tasks.json` items; `03-loop-status.json` 400 `VALIDATION_ERROR` without task_id (not 500); `03-version.json` ok/version shape |
| 4 GUI smoke (Overview/Tasks + promote\|seed) | PASS | SPA `id="root"` + `/assets/` (`04-root.html`, `04-spa-grep.txt`); browser Overview + Tasks list; **Seed honesty**: strict export → red `NOT_IMPLEMENTED` (`04-gui-smoke-notes.txt`, `04-gui-tasks.png`) |
| 5 Security S06 locks | PASS | checklist below; artifacts `05-*` |
| 6 Packaging (disk/embed/help) | PASS | disk `web/dist` wins (Block 4 SPA); embed stub **skipped** (disk path verified only); `06-serve-help.txt` documents `--cors-origin`, `--static-dir` footgun, remote+token |
| 7 Docs | PASS | `07-docs-check.txt`; gui-quickstart two-artifact + loopback/remote/static-dir/501; AGENTS carve-out; CLOUD-APPENDIX design-only |

## Security checklist

- [x] Default 127.0.0.1:7432
- [x] 0.0.0.0 refused w/o --allow-remote (exit=2; `05-refuse-remote.log`)
- [x] No CORS * (`05-cors-verdict.txt`)
- [x] CSP on / (`default-src` / `frame-ancestors`; `05-csp-verdict.txt`)
- [x] /rpc 404 envelope (`05-rpc.json` NOT_FOUND MCP message)
- [x] seed strict/task_id → 501 (`05-seed-strict-code.txt` = 501 + NOT_IMPLEMENTED)
- [x] bad UUID → 400 VALIDATION_ERROR (`05-bad-uuid-code.txt` = 400; dogfood `rl…` id)

## Residuals (non-blocking unless regress)

- listTasks paging: intentional project-local; no library paging — scale bound = full project task list (dogfood ~11 items observed)
- static-dir bound: refuses exact project root only (not `.trace/` alone) — operator footgun; help/quickstart warn against root
- auth/token loopback mint: loopback-trust can mint/rotate without prior bearer (tradeoff; unchanged)
- localStorage `trace.gui.token`: local XSS surface; OK for loopback SPA

## Packaging

- disk web/dist wins: yes (built SPA with `id="root"` + hashed `/assets/`)
- embed stub: skipped (reason: disk path verified only; did not move `web/dist` aside)

## DR-HANDOFF

Still **OPEN** — successor owned by P29-S07-02 (default Phase 30 if green; cloud ≠ Phase 30).
