# P29-S07-02 — DR-HANDOFF Phase 29 close

## Metadata
- id: P29-S07-02
- todo_ids: [P29-S07-02]
- role: verify
- skills: [documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Independent **fresh-session** review of S07-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 29 DR-HANDOFF** with explicit successor (**never TBD**). Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 29 complete when this row is `done`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S07-00 locks
- [01-verify.md](01-verify.md) — locked verify floor
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S07-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-29.md](../../../../TODO/phase-29.md)
- [docs/TODO/phase-30.md](../../../../TODO/phase-30.md) — already scaffolded
- [AGENTS.md](../../../../../../AGENTS.md)
- [CLOUD-APPENDIX.md](../../CLOUD-APPENDIX.md) — design-only; **not** Phase 30
- [phase-30 00-PHASE-PLANNER](../../../phase-30-stray-trace-db-hygiene/00-PHASE-PLANNER.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S07-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-07-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p29-s07-01-verify/evidence/` |
| Phase handoff | `DR-HANDOFF.md` |
| Quickstart | `docs/gui-quickstart.md` |
| Cloud notes | `CLOUD-APPENDIX.md` |
| Phase 30 board | `docs/TODO/phase-30.md` |
| Phase 30 planner | `docs/phases/phase-30-stray-trace-db-hygiene/00-PHASE-PLANNER.md` |

## Locked DR-HANDOFF close policy (FINAL — S07-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S07-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S07-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S07-01 `done`; verify blocks 1–7 green per VERIFY-NOTES + independent spot-check |
| Default successor | **Phase 30** — stray root `trace.db` hygiene ([`docs/TODO/phase-30.md`](../../../../TODO/phase-30.md)); first runnable **P30-00** |
| Cloud / hosted SaaS | **Not** Phase 30 — separate product/repo; CLOUD-APPENDIX remains design-only |
| Regression path | Spawn `P29-S07-02a` implement + `02b` review; **do not** close Phase 29 |
| Must not | Leave `Successor decision: TBD`; claim cloud = Phase 30; rewrite S00–S06 `done` history; ship product in this row; start implementing Phase 30 |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 29 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S07-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| Verify floor green (build/tests/API/GUI/security/docs); residuals only as listed | **Phase 30** | Close DR-HANDOFF; point TODO/AGENTS to Phase 30 / **P30-00**; do not implement P30 here |
| Security lock regression / build or httpapi FAIL / SPA placeholder with dist present | **Do not close** — spawn repair | Keep OPEN; insert 02a/02b |
| VERIFY-NOTES missing blocks or evidence dir absent | **Do not close** — spawn repair or send back S07-01 | Keep OPEN |
| VERIFY PASS but human names different successor | Document human theme in Notes | Scaffold only if human promotes **and** Phase 30 still recorded as superseded/deferred |
| Human wants hosted SaaS next | **Not Phase 30** | Close with successor = separate product/repo (or keep OPEN if human forbids close); Phase 30 still the queued Trace core theme unless human cancels it |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`** (still not TBD).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-29-http-api-browser-gui/scopes/scope-07-verify/VERIFY-NOTES.md
test -d experiments/runs/*-p29-s07-01-verify/evidence || ls experiments/runs/ | grep p29-s07-01
go test ./internal/httpapi/... -count=1
go test ./cmd/trace/ -run Serve -count=1
# Spot security (serve briefly if needed):
#   refuse 0.0.0.0 w/o --allow-remote; curl /v1/health; confirm no CORS *; CSP on /; /rpc 404
test -f docs/gui-quickstart.md
test -f docs/TODO/phase-30.md
test -f docs/phases/phase-30-stray-trace-db-hygiene/00-PHASE-PLANNER.md
```

Confirm VERIFY-NOTES lists the four residuals and does not treat cloud as Phase 30.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00 Research (`RESEARCH.md`)
- [ ] S01 ADR + OpenAPI
- [ ] S02 UX IA
- [ ] S03 HTTP API + review
- [ ] S04 GUI MVP + review
- [ ] S05 GUI rich + review
- [ ] S06 Production hardening + review
- [ ] S07 VERIFY + successor documented (**never TBD**; default **Phase 30**)

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| listTasks no paging | Intentional project-local; future only if library grows paging |
| `--static-dir` root-only refuse | Documented footgun; do not advertise `.trace` as static-dir |
| `POST /v1/auth/token` loopback mint | Loopback-trust; bearer required off-loopback |
| `localStorage` `trace.gui.token` | Acceptable loopback XSS surface |
| Cloud / multi-tenant | Separate product — CLOUD-APPENDIX design-only |

### DR-HANDOFF.md update template (on APPROVE — default Phase 30)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **Phase 30** — stray root `trace.db` hygiene |
| Phase 29 outcome | Opt-in `trace serve` + browser GUI; Law 19 adapters; OpenAPI local+cloud-ready contract; S06 production locks verified |
| Verify | Cite VERIFY-NOTES + evidence dir; security checklist green |
| Residuals (non-blocking) | listTasks paging; static-dir root-only; auth/token loopback; localStorage token |
| Cloud | Not Phase 30 — separate product/repo (CLOUD-APPENDIX) |
| Forward | First runnable: **P30-00** @ docs/TODO/phase-30.md |
```

### TODO.md / AGENTS.md updates (on APPROVE)

1. **`docs/TODO.md` orchestrator paste:** Active phase → **Phase 30**; next runnable → **P30-00**; note Phase 29 closed.
2. **Phase boards table:** Phase 29 status `done`; Phase 30 status active/queued→runnable.
3. **`AGENTS.md` Current focus:** Phase 29 complete; Phase 30 next (stray root `trace.db`); do **not** describe cloud SaaS as Phase 30.
4. Confirm Phase 30 scaffold already present (README, `00-PHASE-PLANNER`, scope stubs, board rows 527+). **Do not** invent a Phase 31 cloud board here.

### On FAIL / repair spawn

Insert immediately below this row:

| Order | ID | Role |
|------:|----|------|
| 526a | P29-S07-02a | implement repair |
| 526b | P29-S07-02b | review repair |

Keep DR-HANDOFF **OPEN**. Do not mark Phase 29 done.

## Exit criteria

- [ ] Independent spot-check done; VERIFY-NOTES credible
- [ ] DR-HANDOFF closed with explicit successor (**Phase 30** or OPEN+repairs / human override)
- [ ] TODO.md / AGENTS.md reflect close
- [ ] No TBD successor field; no “cloud = Phase 30” confusion
- [ ] Phase 30 scaffold confirmed present; P30 work **not** started in this row
- [ ] Board Notes on **P29-S07-02** only

## Todo updates

Status + notes on **P29-S07-02**.

## Next

Phase 30 **P30-00** (after green close) — or repair spawn if VERIFY failed
