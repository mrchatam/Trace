# P29-S07-00 — Scope planner (VERIFY)

## Metadata
- id: P29-S07-00
- todo_ids: [P29-S07-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective

Lock VERIFY + DR-HANDOFF prompts and evidence layout for Phase 29 close.

## Session start

Follow agent-loop-protocol Session start.

## Verify floor

| Block | Check |
|-------|-------|
| Build | `go build -o bin/trace ./cmd/trace`; `cd web && npm run build` (two-artifact path) |
| Tests | `go test ./internal/httpapi/...`; `go test ./cmd/trace/ -run Serve`; optional `cd web && npm run test:e2e` |
| API smoke | `trace serve` → `/v1/health` + `/v1/tasks` + loop status |
| GUI smoke | Browser open; Overview/Tasks + promote or seed honesty path |
| Security (S06 locks) | Default `127.0.0.1:7432`; `0.0.0.0` without `--allow-remote` refused; CORS never `*`; CSP on `/`; `/rpc` → 404 envelope; seed `strict`/`task_id` → 501; bad UUID loop id → 400 `VALIDATION_ERROR` |
| Docs | [`docs/gui-quickstart.md`](../../../../../gui-quickstart.md) + AGENTS carve-out + `CLOUD-APPENDIX.md` design-only |

Evidence dir: `experiments/runs/YYYY-MM-DD-p29-s07-01-verify/`

### Residuals to record in VERIFY-NOTES (not fail gates unless regress)

- **listTasks** has no library paging (project-local intentional) — note scale bound
- **`--static-dir`** only refuses exact project root (not `.trace/` alone) — operator footgun; do not document unsafe dirs
- **`POST /v1/auth/token`** on loopback-trust can mint/rotate without prior bearer (loopback-trust tradeoff)
- GUI bearer in `localStorage` (`trace.gui.token`) — local XSS surface; OK for loopback SPA

## Successor defaults

| Outcome | Successor |
|---------|-----------|
| Production checklist green | Close Phase 29; **Phase 30** already queued (stray root `trace.db` hygiene) — runnable only after **P29-S07-02** |
| API/GUI gaps | repair spawn; keep DR-HANDOFF OPEN |
| Cloud / hosted SaaS | **separate product/repo** — not Phase 30; see CLOUD-APPENDIX (design only) |

## Exit criteria

- [x] `01-verify.md` + `02-dr-handoff.md` runnable
- [x] Next **P29-S07-01**

## Next

P29-S07-01 → P29-S07-02
