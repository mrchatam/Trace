# P29-S07-01 — VERIFY implementer

## Metadata
- id: P29-S07-01
- todo_ids: [P29-S07-01]
- role: implementer
- skills: [grinding-until-pass, webapp-testing]
- mcps: [cursor-ide-browser]
- verification: mixed
- hooks: []

## Objective

Run Phase 29 verify floor after S00–S06. Capture evidence, write **`VERIFY-NOTES.md`** with per-block PASS/FAIL, packaging check, and S06 residuals. Keep **`DR-HANDOFF.md` OPEN** — S07-02 owns successor close. **No product code.** Do **not** start Phase 30.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — verify floor + residuals (FINAL)
- [ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)
- [api/openapi.yaml](../../../../../../api/openapi.yaml)
- [docs/gui-quickstart.md](../../../../../gui-quickstart.md)
- [CLOUD-APPENDIX.md](../../CLOUD-APPENDIX.md) — design-only
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S07-02
- [FEATURE-MATRIX.md](../scope-05-gui-rich/FEATURE-MATRIX.md) — M01–M07 done baseline
- Live: `internal/httpapi/`, `cmd/trace/serve.go`, `web/`, `internal/httpapi/embeddist/`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or run P30-*.

## Locked defaults (FINAL — S07-00)

| Item | Value |
|------|-------|
| Precondition | P29-S00-00 … P29-S06-02 all `done` |
| Product Go / SPA | **Forbidden** (evidence + notes only) |
| Binary | Rebuild `bin/trace` from repo HEAD before smoke |
| Default bind | `127.0.0.1:7432` |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p29-s07-01-verify/` (tee under `evidence/` subdir) |
| Notes artifact | `scopes/scope-07-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S07-02 closes |
| Successor | **Out of scope** — S07-02 only (Phase 30 already queued; cloud ≠ Phase 30) |
| E2E | Optional `cd web && npm run test:e2e` — prefer run; note skip reason if blocked |
| Embed | Confirm disk `web/dist` wins; embed stub when disk missing — do **not** rewrite `embeddist/` for everyday verify |

### Fail vs residual (locked)

**Fail VERIFY for:** build FAIL; `go test ./internal/httpapi/...` FAIL; `go test ./cmd/trace/ -run Serve` FAIL; `/v1/health` or `/v1/tasks` smoke FAIL; security lock regression (open bind default, CORS `*`, missing CSP on `/`, `/rpc` not 404 envelope, seed `strict`/`task_id` not 501, bad UUID → not 400 `VALIDATION_ERROR`); missing/broken `docs/gui-quickstart.md` or AGENTS carve-out contradiction; SPA placeholder when `web/dist` present.

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES):

| Residual | Disposition |
|----------|-------------|
| `listTasks` no library paging | Intentional project-local; note scale bound |
| `--static-dir` refuses exact project root only (not `.trace/` alone) | Operator footgun; do not document unsafe dirs |
| `POST /v1/auth/token` on loopback-trust can mint/rotate without prior bearer | Loopback-trust tradeoff |
| GUI bearer in `localStorage` key `trace.gui.token` | Local XSS surface; OK for loopback SPA |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir.

### Block 0 — Evidence dir

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p29-s07-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P29-S07-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
} > "$EVID/99-run-metadata.txt"
```

### Block 1 — Build (required)

```bash
go build -o bin/trace ./cmd/trace 2>&1 | tee "$EVID/01-go-build.txt"
cd web && npm run build 2>&1 | tee "../$EVID/01-web-build.txt"
cd ..
test -f web/dist/index.html
```

**Pass:** both builds exit 0; `web/dist/index.html` exists (two-artifact path).

### Block 2 — Tests (required + optional e2e)

```bash
go test ./internal/httpapi/... -count=1 2>&1 | tee "$EVID/02-httpapi-test.txt"
go test ./cmd/trace/ -run Serve -count=1 2>&1 | tee "$EVID/02-serve-test.txt"
# Optional but preferred:
cd web && npm run test:e2e 2>&1 | tee "../$EVID/02-e2e.txt" || echo "e2e_skipped_or_fail" | tee "../$EVID/02-e2e.txt"
cd ..
```

**Pass:** httpapi + Serve tests exit 0. E2E: PASS preferred; if skip/fail, document in VERIFY-NOTES (do not alone fail verify unless promote/export regression is clear and reproducible).

### Block 3 — API smoke (required)

```bash
# Start serve in background; kill after curls
./bin/trace serve --addr 127.0.0.1:7432 >"$EVID/03-serve.log" 2>&1 &
SERVE_PID=$!
sleep 1
curl -sS "http://127.0.0.1:7432/v1/health" | tee "$EVID/03-health.json"
curl -sS "http://127.0.0.1:7432/v1/tasks" | tee "$EVID/03-tasks.json"
curl -sS "http://127.0.0.1:7432/v1/loop/status" | tee "$EVID/03-loop-status.json" || true
kill $SERVE_PID 2>/dev/null || true
wait $SERVE_PID 2>/dev/null || true
```

**Pass:** health OK (200 + ok/version shape); tasks returns JSON items (dogfood tasks OK); loop status returns JSON envelope (may be error body for non-UUID dogfood ids — must not be 500 `INTERNAL_ERROR` for invalid UUID; expect 400 `VALIDATION_ERROR` when probing bad UUID — see Block 5).

### Block 4 — GUI smoke (required)

```bash
./bin/trace serve --addr 127.0.0.1:7432 >"$EVID/04-serve-gui.log" 2>&1 &
SERVE_PID=$!
sleep 1
curl -sS -D "$EVID/04-root-headers.txt" -o "$EVID/04-root.html" "http://127.0.0.1:7432/"
# Expect SPA (id="root" or bundled assets), not plain placeholder when web/dist present
grep -E 'id="root"|/assets/' "$EVID/04-root.html"
kill $SERVE_PID 2>/dev/null || true
wait $SERVE_PID 2>/dev/null || true
```

Browser (cursor-ide-browser or equivalent): open `http://127.0.0.1:7432/` → Overview + Tasks visible; exercise **one** of: Discoveries promote path **or** Seed honesty path (export/import UI honesty / 501 messaging). Screenshot or note path under `$EVID/`.

**Pass:** SPA served; Overview/Tasks usable; promote **or** seed honesty path demonstrated.

### Block 5 — Security locks from S06 (required)

```bash
# Default refuse open bind
./bin/trace serve --addr 0.0.0.0:7432 >"$EVID/05-refuse-remote.log" 2>&1; echo "exit=$?" | tee -a "$EVID/05-refuse-remote.log"
# Expect non-zero without --allow-remote

./bin/trace serve --addr 127.0.0.1:7432 >"$EVID/05-sec-serve.log" 2>&1 &
SERVE_PID=$!
sleep 1
# CORS never *
curl -sS -D "$EVID/05-cors-headers.txt" -o /dev/null -H "Origin: https://evil.example" "http://127.0.0.1:7432/v1/health"
grep -i 'access-control-allow-origin: \*' "$EVID/05-cors-headers.txt" && echo "FAIL cors star" || echo "PASS no cors star" | tee "$EVID/05-cors-verdict.txt"
# CSP on /
curl -sS -D "$EVID/05-csp-headers.txt" -o /dev/null "http://127.0.0.1:7432/"
grep -i 'content-security-policy' "$EVID/05-csp-headers.txt" | tee "$EVID/05-csp-verdict.txt"
# /rpc → 404 envelope
curl -sS -D "$EVID/05-rpc-headers.txt" -o "$EVID/05-rpc.json" "http://127.0.0.1:7432/rpc"
# seed strict/task_id → 501
curl -sS -o "$EVID/05-seed-strict.json" -w "%{http_code}" -X POST "http://127.0.0.1:7432/v1/seed/export" \
  -H 'Content-Type: application/json' \
  -d '{"path":"trace/graph.json","strict":true}' | tee "$EVID/05-seed-strict-code.txt"
# bad UUID → 400 VALIDATION_ERROR (loop status or similar)
curl -sS -o "$EVID/05-bad-uuid.json" -w "%{http_code}" \
  "http://127.0.0.1:7432/v1/loop/status?task_id=rl010000-0000-4000-8000-000000000010" | tee "$EVID/05-bad-uuid-code.txt"
kill $SERVE_PID 2>/dev/null || true
wait $SERVE_PID 2>/dev/null || true
```

**Pass checklist (all required):**

- [ ] Default bind loopback (`127.0.0.1:7432`)
- [ ] `0.0.0.0` without `--allow-remote` refused (exit ≠ 0)
- [ ] CORS never `Access-Control-Allow-Origin: *`
- [ ] CSP present on `/` (at least `default-src` / `frame-ancestors` family)
- [ ] `/rpc` → 404 JSON envelope (not MCP transport)
- [ ] Seed export `strict`/`task_id` → 501 `NOT_IMPLEMENTED` (or documented 501 envelope)
- [ ] Bad/non-UUID loop id → 400 `VALIDATION_ERROR` (not 500 `INTERNAL_ERROR`)

### Block 6 — Packaging (required)

| Check | How |
|-------|-----|
| Disk wins | With `web/dist/index.html` present, `/` serves built SPA (Block 4) |
| Embed fallback | Optional: temporarily move/rename `web/dist` aside, restart serve, confirm embed stub or placeholder (not crash); restore dist after. If skip, note “disk path verified only” |
| Help text | `./bin/trace serve --help` mentions `--cors-origin`, `--static-dir` footgun, remote+token |

```bash
./bin/trace serve --help 2>&1 | tee "$EVID/06-serve-help.txt"
```

### Block 7 — Docs (required)

Confirm files exist and are coherent (spot-read; no rewrite required unless contradiction blocks ship):

| Doc | Expect |
|-----|--------|
| [`docs/gui-quickstart.md`](../../../../../gui-quickstart.md) | Two-artifact path; loopback default; remote+token; static-dir footgun; seed 501 honesty |
| [`AGENTS.md`](../../../../../../AGENTS.md) Hard boundaries | Opt-in `trace serve`; Law 19 adapters; still-forbidden open defaults; cloud via OpenAPI later |
| [`CLOUD-APPENDIX.md`](../../CLOUD-APPENDIX.md) | Design-only; no SaaS ship claim |

```bash
test -f docs/gui-quickstart.md
test -f docs/phases/phase-29-http-api-browser-gui/CLOUD-APPENDIX.md
grep -n 'trace serve\|127.0.0.1' AGENTS.md | head
```

## VERIFY-NOTES.md template (required)

Write to `scopes/scope-07-verify/VERIFY-NOTES.md`:

```markdown
# Phase 29 VERIFY notes

- Date: YYYY-MM-DD
- Git SHA: …
- Evidence: experiments/runs/YYYY-MM-DD-p29-s07-01-verify/evidence/
- Verdict: PASS | FAIL

## Blocks

| Block | Result | Evidence |
|-------|--------|----------|
| 1 Build (go + web) | PASS/FAIL | … |
| 2 Tests (httpapi, Serve, e2e?) | PASS/FAIL/SKIP | … |
| 3 API smoke (health/tasks/loop) | PASS/FAIL | … |
| 4 GUI smoke (Overview/Tasks + promote\|seed) | PASS/FAIL | … |
| 5 Security S06 locks | PASS/FAIL | checklist ticks |
| 6 Packaging (disk/embed/help) | PASS/FAIL | … |
| 7 Docs | PASS/FAIL | … |

## Security checklist

- [ ] Default 127.0.0.1:7432
- [ ] 0.0.0.0 refused w/o --allow-remote
- [ ] No CORS *
- [ ] CSP on /
- [ ] /rpc 404 envelope
- [ ] seed strict/task_id → 501
- [ ] bad UUID → 400 VALIDATION_ERROR

## Residuals (non-blocking unless regress)

- listTasks paging: …
- static-dir bound: …
- auth/token loopback mint: …
- localStorage `trace.gui.token`: …

## Packaging

- disk web/dist wins: …
- embed stub: verified | skipped (reason)

## DR-HANDOFF

Still **OPEN** — successor owned by P29-S07-02 (default Phase 30 if green; cloud ≠ Phase 30).
```

## Out of scope

- Closing `DR-HANDOFF.md` or updating TODO/AGENTS phase focus for “Phase 29 done”
- Starting Phase 30 / inventing a cloud phase
- Product code fixes (if blocker: mark `failed`/`blocked` + note; S07-02 or repair spawn)

## Exit criteria

- [ ] VERIFY-NOTES.md complete (security + docs blocks explicit)
- [ ] Evidence dir referenced with Block 0–7 artifacts
- [ ] Residuals recorded (listTasks, static-dir, auth/token loopback, localStorage token)
- [ ] DR-HANDOFF still OPEN
- [ ] Board Notes on **P29-S07-01** only with evidence path + verdict

## Todo updates

Status + notes on **P29-S07-01** only.

## Next

**P29-S07-02**
