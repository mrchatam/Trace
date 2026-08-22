# P34-S05-01 — VERIFY

## Metadata
- id: P34-S05-01
- todo_ids: [P34-S05-01]
- role: verify
- skills: [planning-and-task-breakdown, grinding-until-pass]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Run Phase 34 VERIFY floor after S00–S04 (**PLAN T9**). Aggregate prior PASS cites + live re-checks into **`VERIFY-NOTES.md`** (+ evidence dir). Confirm locks **L1–L4** + docs consumer story. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P34-S05-02**. **No product code.** Do **not** start S05-02 or invent a successor phase.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S05-00)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1–L4
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S05-02
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — T1–T11; **T9** = this floor
- Prior PASS artifacts (cite in notes):
  - S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md) + [`02-review`](../scope-00-research/02-review.md)
  - S01 [`PLAN.md`](../scope-01-plan/PLAN.md) + board Notes
  - S02 board Notes (P34-S02-01/02) — T1–T3 embed
  - S03 board Notes (P34-S03-01/02) — T4–T7/T11 auto-port
  - S04 board Notes (P34-S04-01/02) — T8 docs
- Live anchors: `cmd/trace/` (`gui`, `serve`, `local_http.go`, `help.go`), `internal/httpapi/` (`static.go`, `embeddist/`, `auto_port.go`, `addr_in_use.go`), `scripts/embed-gui.sh`, `docs/gui-quickstart.md`, `web/README.md`, `AGENTS.md`, root `README.md`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product bodies.

## Locked defaults (FINAL — S05-00)

| Item | Value |
|------|-------|
| Precondition | P34-00 … P34-S04-02 all `done`; S02/S03/S04 reviews **PASS** |
| Product / Go / TS / docs changes | **Forbidden** (evidence + notes only). Failures → spawn remediation from this row or leave FAIL for S05-02 to spawn |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p34-s05-01-verify/evidence/` |
| Notes artifact | `scopes/scope-05-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S05-02 closes |
| Successor | **Out of scope** — S05-02 only (lean default **no successor**) |
| VERIFY = PLAN **T9** | Combine **T1** (real SPA markers) + **T4/T5** (concurrent auto-port) + **T8** (docs) + **T10** stub-fail if embed still stub when full SPA intended |
| Port range | Default auto-hop `127.0.0.1:7432`–`7441` (`MaxAutoPortAttempts` = 10); `--addr` pin-strict (`fs.Visit` / AddrExplicit) |
| Static order | disk-if-`index.html` → embed → `placeholderHTML`; **no** SPA copy under consumer `.trace/` |
| Contributor DX | Trace-checkout `web/` / disk `web/dist` OK if labeled — not a consumer requirement |

### L1–L4 + docs acceptance map (must tick in VERIFY-NOTES)

| Lock | Acceptance (spot-check + prior PASS cite) |
|------|-------------------------------------------|
| **L1 — Consumer layout** | Consumer-like temp (`.trace/` init, **no** `web/`) uses only `.trace/` Trace artifacts; no required project `web/` |
| **L2 — GUI asset source** | `GET /` (or static index) from that temp is **real SPA**: `id="root"` and/or `/assets/` module; **not** phrase `Embedded GUI stub`; shipped `embeddist/index.html` not stub-only (**T10**) |
| **L3 — Auto port** | Default busy → hop next free in range; second concurrent default `gui`/`serve` → **distinct** port; printed/opened URL matches bound addr (**T4/T5**); free default stays `:7432` (**T11** cite OK) |
| **L4 — One process = one root** | Each process one store root (`-C`/cwd); multi-project = N processes × N ports (spot-check concurrent test / docs) |
| **Docs (T8)** | Quickstart / help / `web/README` / AGENTS / embeddist: embed + auto-port; **no** consumer two-artifact primary; **no** “no auto free-port” for default |

### Aggregate evidence map (S00–S04 — cite in VERIFY-NOTES)

| Scope | Must cite / re-check |
|-------|----------------------|
| S00 | RESEARCH leans: embed=A; UA-incr auto-port; L3 supersedes P33 reject |
| S01 | PLAN T1–T11 + board S02→S05; StaticDir opportunistic; flag.Changed/`Visit` intent |
| S02 | Real SPA in `embeddist`; T1–T3 tests; README consumer `.trace/` only; placeholder tone |
| S03 | Shared hop in httpapi; T4–T7/T11 + T5 concurrent CLI; T6 explicit busy no hop; help/usage auto-port |
| S04 | T8 greps PASS; `gui-quickstart` primary embed + `7432`–`7441` + `--addr` pin; contributor `web/README` |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Stub / placeholder SPA on consumer path when release embed should be full (`Embedded GUI stub` in shipped embed, or `GET /` lacks `#root`/`/assets/` markers)
- No auto-port on default busy (second concurrent still fails or same port / wrong URL)
- Docs still require consumer `web/` or teach two-artifact / “no auto free-port” as primary default story
- Public bind default (`0.0.0.0` / non-loopback as default)
- SPA copied into consumer `.trace/` as primary asset path
- Focused Go tests for embed/auto-port FAIL (blocks below)
- Law 19 fork / always-on daemon / OS `:0` as Trace default claimed as shipped

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES):

| Residual | Disposition |
|----------|-------------|
| Cosmetic help / wording nits | Accept if T8 greps clean and positive story present |
| Contributor Trace-checkout `web/` DX still documented | Accept if labeled **contributor**, not consumer |
| Default StaticDir path string still `<root>/web/dist` | Accept — resolution is disk→embed→placeholder; consumers rarely need `--static-dir` |
| Optional CI workflow for embed-gui not invented | Out of phase (PLAN deferred) |
| Explore UI / craft redesign | Phase 33 closed — out of phase |
| Hosted SaaS / brew/deb | Out of scope |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name. Prefer `-p 1` on packages that bind `7432`–`7441` to avoid parallel port races.

### Block 0 — Evidence dir + preflight

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p34-s05-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P34-S05-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P34-S04-02 PASS; S00–S04 done; PLAN T9 floor"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists; metadata cites S04-02 PASS + S00–S04 complete.

### Block 1 — Embed / static (T1–T3 + T10 stub-fail)

```bash
go test ./internal/httpapi/ -run 'TestStaticCSPAndEmbedFallback|TestStaticDiskWinsOverEmbed|TestStaticDirRefuseProjectRoot' -count=1 -p 1 \
  2>&1 | tee "$EVID/01-static-embed.txt"
{
  echo "=== embeddist index markers ==="
  test -f internal/httpapi/embeddist/index.html && echo "index.html present"
  grep -n 'id="root"\|/assets/\|Embedded GUI stub' internal/httpapi/embeddist/index.html | head -40
  echo "=== stub phrase must be ABSENT from embeddist index ==="
  if grep -q 'Embedded GUI stub' internal/httpapi/embeddist/index.html; then
    echo "FAIL: stub phrase in shipped embed"
  else
    echo "PASS: no Embedded GUI stub in embeddist/index.html"
  fi
} 2>&1 | tee "$EVID/01b-embeddist-markers.txt"
```

**Pass:** focused static tests exit 0; `embeddist/index.html` has `#root` and/or `/assets/`; **no** `Embedded GUI stub`.

### Block 2 — Auto-port Go (T4/T7/T11 + explicit)

```bash
go test ./internal/httpapi/ -run 'TestListenAutoPort_' -count=1 -p 1 \
  2>&1 | tee "$EVID/02-httpapi-auto-port.txt"
```

**Pass:** exit 0 — free stays 7432 (T11); busy hops next (T4); explicit busy no hop; exhausted range fails with range/`--addr` (T7).

### Block 3 — Concurrent CLI + pin-strict (T5/T6)

```bash
go test ./cmd/trace/ -run 'TestGuiServeConcurrentDefaultDistinctPorts|TestGuiExplicitDefaultAddrBusyNoHop|TestServeExplicitDefaultAddrBusyNoHop' -count=1 -p 1 \
  2>&1 | tee "$EVID/03-cmd-concurrent-pin.txt"
```

**Pass:** exit 0 — concurrent defaults distinct ports + open URL matches; explicit `--addr`==DefaultAddr busy fails (no hop).

### Block 4 — Broader httpapi + cmd/trace smoke (required)

```bash
go test ./internal/httpapi/ ./cmd/trace/ -count=1 -p 1 \
  2>&1 | tee "$EVID/04-go-packages.txt"
go run ./cmd/trace gui --help 2>&1 | tee "$EVID/04b-gui-help.txt"
go run ./cmd/trace serve --help 2>&1 | tee "$EVID/04c-serve-help.txt"
```

**Pass:** package tests exit 0; help mentions auto-port / `7432`–`7441` (or equivalent) and `--addr` pin; loopback default intact.

### Block 5 — Docs consumer story (T8) spot-check

```bash
{
  echo "=== gui-quickstart head ==="
  head -50 docs/gui-quickstart.md
  echo "=== forbidden phrases (must be empty / non-primary) ==="
  grep -nE 'no auto free-port|does not auto-pick|Pick a free port yourself|embedded stub|two-artifact' \
    docs/gui-quickstart.md web/README.md AGENTS.md README.md \
    cmd/trace/help.go cmd/trace/local_http.go \
    internal/httpapi/embeddist/README.md internal/httpapi/addr_in_use.go 2>/dev/null || true
  echo "=== positive markers ==="
  grep -nE '7432|7441|embed|\.trace/|auto|pin|--addr' docs/gui-quickstart.md | head -60
  echo "=== contributor label on web/README ==="
  head -30 web/README.md
} 2>&1 | tee "$EVID/05-docs-t8.txt"
```

**Pass checklist (all required in VERIFY-NOTES):**

1. Quickstart primary = `trace gui` + binary embed + `.trace/` only.
2. Multi-project auto-hop `7432`–`7441` without required `--addr` for default.
3. `--addr` pin-strict documented.
4. Greppable consumer surfaces free of forbidden “no auto free-port” / consumer two-artifact primary / shipped-stub teaching.
5. `web/README` contributor-labeled (if it mentions `web/` build).

### Block 6 — Live consumer-temp optional smoke (recommended)

If environment allows (built `trace` on PATH or `go run`), prove consumer-like temp without `web/`:

```bash
TMP=$(mktemp -d)
cleanup() { rm -rf "$TMP"; }
trap cleanup EXIT
go run ./cmd/trace -C "$TMP" init 2>&1 | tee "$EVID/06a-init.txt"
# Prefer a short-lived serve/gui --no-open; capture GET / markers via curl if listen succeeds.
# Record exact commands + ports in VERIFY-NOTES. If flaky CI/env, cite Block 1–3 tests as primary and mark this WAIVED with reason.
```

**Pass:** live SPA markers from temp **or** explicit waive with reason + Blocks 1–3 green (automated T1/T5 cover the floor).

### Block 7 — Residuals + aggregate (required in notes)

1. Cite S00–S04 PASS board Notes / artifacts (paths above).
2. Tick L1–L4 + Docs in VERIFY-NOTES.
3. List non-blocking residuals from the locked table.
4. Overall PASS only if blocks 0–5 green (+ block 6 live or waived) and no fail criteria tripped.

## Do not

- Change product Go/TS/CSS/docs in this row (spawn or FAIL instead)
- Close DR-HANDOFF (that is **P34-S05-02**)
- Claim/scaffold successor beyond citing lean **no successor** (S05-02 only)
- Reopen S02 embed / S03 hop / S04 docs as “fix” without spawn
- Treat contributor `web/` DX or StaticDir path-string as VERIFY FAIL
- Fake PASS without evidence dir + notes

## VERIFY-NOTES.md template (required)

Write `docs/phases/phase-34-gui-packaging-multiproject/scopes/scope-05-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — P34-S05-01

**Date:** YYYY-MM-DD
**Git SHA:** …
**Overall:** PASS | FAIL
**Evidence:** experiments/runs/YYYY-MM-DD-p34-s05-01-verify/evidence/
**Precondition:** P34-S04-02 PASS; S00–S04 done; PLAN T9

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS/FAIL | |
| 1 Static embed T1–T3 + T10 | PASS/FAIL | |
| 2 httpapi auto-port T4/T7/T11 | PASS/FAIL | |
| 3 Concurrent + pin T5/T6 | PASS/FAIL | |
| 4 Packages + help | PASS/FAIL | |
| 5 Docs T8 | PASS/FAIL | |
| 6 Live consumer-temp | LIVE path… \| WAIVED reason… | |
| 7 Residuals + aggregate | listed | |

## L1–L4 + Docs
- [ ] L1 consumer `.trace/` only — …
- [ ] L2 real SPA from binary (not stub) — …
- [ ] L3 auto-port concurrent + correct URL — …
- [ ] L4 one process = one root — …
- [ ] Docs T8 embed + auto-port — …

## Aggregate (S00–S04)
- S00 RESEARCH: …
- S01 PLAN: …
- S02 embed: … (T1–T3)
- S03 auto-port: … (T4–T7/T11)
- S04 docs: … (T8)

## Residuals (non-blocking)
- …

## Failures (if any)
- …

## DR-HANDOFF
remains OPEN — close owner **P34-S05-02**

## Next
P34-S05-02
```

## Todo updates

Status + notes on **P34-S05-01** only. Do not mark S05-02 done.

## Exit criteria

- [ ] `VERIFY-NOTES.md` with overall PASS/FAIL and blocks 0–7
- [ ] Evidence dir under `experiments/runs/…-p34-s05-01-verify/evidence/`
- [ ] L1–L4 + Docs explicitly ticked
- [ ] PLAN T9 / T10 stub-fail addressed in notes
- [ ] Board Notes summarize results + evidence path
- [ ] DR-HANDOFF still **OPEN**
- [ ] Next: **P34-S05-02**

## Next

`P34-S05-02`
