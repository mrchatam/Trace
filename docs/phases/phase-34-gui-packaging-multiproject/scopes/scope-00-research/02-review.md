# P34-S00-02 — Research review

## Metadata
- id: P34-S00-02
- todo_ids: [P34-S00-02]
- role: reviewer
- skills: [code-review-and-quality, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of `RESEARCH.md` against DESIGN-LOCKS, INTAKE, and live repo facts. Fix trivial RESEARCH gaps inline; thicken **upcoming** S01 prompts if structural. **No product code.** Do not rewrite Phases 00–33 done history.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1–L4 + clarifications
- [INTAKE.md](../../INTAKE.md)
- [01-research.md](01-research.md) — template + rejects + must-answer
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Artifact: [RESEARCH.md](RESEARCH.md)
- Supersession cites: [P33 RESEARCH](../../../phase-33-gui-craft-hook-launch/scopes/scope-00-research/RESEARCH.md), [P32 RESEARCH](../../../phase-32-graph-first-gui/scopes/scope-00-research/RESEARCH.md)
- Spot-check live: `internal/httpapi/static.go`, `embeddist/`, `addr_in_use.go`, `cmd/trace/local_http.go`, `cmd/trace/help.go`, `docs/gui-quickstart.md`
- Peer (optional spot): UA `viewer.mjs` `listen` + `portExplicit`

## Session start

Follow agent-loop-protocol Session start. Fresh context — do not share the S00-01 session. Unattended: proceed from locks + artifact.

## Checklist

### Locks (DESIGN-LOCKS + INTAKE)

- [ ] L1–L4 honored (research informs; does not reopen)
- [ ] L1: consumer footprint = `.trace/` only; no required `web/`
- [ ] L2: real SPA from Trace product/binary recommended; disk = Trace-checkout / `--static-dir` DX only
- [ ] L2: stub last-resort; release VERIFY must fail if stub shipped when full SPA intended (called out or compatible with handoff)
- [ ] L2 reject: SPA under consumer `.trace/` as primary **not** recommended
- [ ] L3: default auto free **loopback** port; print + open URL story present
- [ ] L3: `--addr` remains **strict** fail-if-busy
- [ ] L3 supersession: P33/P32 reject/defer of UA auto-port **documented as overturned for default bind only**
- [ ] L4: one process = one project; multi-project = N processes × N ports
- [ ] Law 19 / loopback defaults noted; no public bind default; no always-on daemon; no SaaS

### Artifact quality

- [ ] All RESEARCH template headings present and non-empty (incl. **L3 supersession**)
- [ ] Live baseline cites real paths and matches spot-check (or documented drift)
- [ ] Embed options ranked + **one** recommendation (prefer go:embed pipeline unless blocked)
- [ ] StaticDir policy actionable for S01/S02 (consumer vs Trace-checkout vs `--static-dir`)
- [ ] Auto-port algorithm actionable for S01/S03 (start port, attempts/range, host, `--addr` detection, gui/serve share)
- [ ] Docs/layout audit table non-empty with S02/S03/S04 owners
- [ ] Rejected alternatives cover rejects **1–5** from 01-research.md
- [ ] Peer cite accurate: UA increment + `portExplicit` + loopback (borrow under L3 for default only)

### Forward board rights

- [ ] PASS → next **P34-S01-00**; optionally thicken S01 stubs if plan-shaped gaps (pipeline cmds, flag semantics, test seeds)
- [ ] blocker/high → inline fix RESEARCH or spawn `P34-S00-02a` / `02b` before PASS
- [ ] Do **not** start S02 product work from this review

## Common failure modes (fail or high)

| Symptom | Likely severity |
|---------|-----------------|
| Recommends keeping no-auto-port on default | blocker (vs L3) |
| Auto-hop on explicit `--addr` | blocker |
| Consumer must have `web/` or two-artifact primary | blocker |
| SPA primary path = copy into consumer `.trace/` | blocker |
| Missing L3 supersession / claims P33 still forbids auto-port | high |
| Embed lean = install-sidecar without proving embed blocked | high (conflicts L2 lean) |
| `:0` as primary without UX justification | medium |
| Docs audit empty or only Phase README | medium |
| gui auto-port but serve left fail-only with no justification | medium |

## Findings format

`blocker` | `high` | `medium` | `low` | `nit` — every blocker/high needs fix or spawn before PASS.

## Exit criteria

- [ ] Confidence medium or high with evidence in Notes
- [ ] No open blocker/high without pending follow-up
- [ ] Board Notes: PASS/FAIL + short evidence bullets
- [ ] Next **P34-S01-00** on PASS

## Todo updates

Status + notes on **P34-S00-02** only.

## Next

`P34-S01-00`
