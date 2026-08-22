# P27-S00-02 — Investigation review

## Metadata
- id: P27-S00-02
- todo_ids: [P27-S00-02]
- role: reviewer
- skills: [code-review-and-quality, investigator]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

Independent review of `AUDIT.md` — confirm INT-07/08/10 mapping and Phase 26 P25-3 residual linkage before S01 planning. **Fresh subagent** — do not reuse S00-01 session.

## References

- [01-investigation.md](01-investigation.md) — exit criteria + AUDIT template
- [00-PLANNER.md](00-PLANNER.md) — locked defaults
- [AUDIT.md](AUDIT.md) (produced by S00-01)
- [Phase 26 VERIFY-NOTES](../../../phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md)
- [Phase 24 CODEBASE-AUDIT](../../../phase-24-agent-effectiveness-investigation/scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Verify checklist

### Deliverable presence

- [ ] `AUDIT.md` exists in this scope folder
- [ ] All template sections present: Executive summary, Phase 26 residual, INT mapping table, P24 delta, Risks/open decisions
- [ ] Each of INT-07, INT-08, INT-10 has ≥1 row in mapping table

### Phase 26 residual accuracy

- [ ] P25-3 failure mode: `discoveries=0 decisions=0` on build-only G1
- [ ] P25-1/2 PASS cited; explicitly **not** framed as P25-C install regression
- [ ] Manual export step documented (`prepare.sh` does not export; operator runs `seed export` before score)
- [ ] Evidence paths match live files under `experiments/runs/2026-08-20-p26-s05-01-verify/evidence/`

### Live spot-checks (reviewer runs)

From repo root, confirm AUDIT claims against live code:

```bash
# P25-3 check location
grep -n "P25-3" experiments/ab-p25-gap-pass-validation/score.sh

# --strict / GateForExport wiring
grep -n "strict\|GateForExport\|collectExportViolations" cmd/trace/seed.go internal/loop/gate.go

# prepare.sh: import yes, export no
grep -n "seed import\|seed export" experiments/ab-p25-gap-pass-validation/prepare.sh

# Two-session separation in protocol
grep -n "build-only\|Session-B\|gap analysis" experiments/ab-p25-gap-pass-validation/PROTOCOL.md
```

- [ ] At least 3 AUDIT path claims verified with grep/read (not trust implementer alone)
- [ ] Stale P24 line refs flagged or updated in AUDIT

### Scope hygiene

- [ ] No product code changed on S00-01 (`git diff` — only `AUDIT.md` expected)
- [ ] S01/S02 task seeds name **files + gaps** (e.g. `score.sh` missing import-before-gate), not INT themes only
- [ ] Threshold numbers listed as **options**, not locked decisions

### Confidence rubric

| Level | When to use |
|-------|-------------|
| **high** | All checklist items pass; spot-checks match AUDIT; no open high findings |
| **medium** | Minor gaps (stale line ref, missing optional risk row) with explicit residual risks in Notes |
| **low** | Missing INT, wrong P25-3 interpretation, or unverified paths — **do not APPROVE** |

## Spawn policy

| Severity | Action |
|----------|--------|
| **blocker** | Missing INT section, wrong P25-3 interpretation, or AUDIT absent → spawn `P27-S00-02a` (AUDIT fix) + `P27-S00-02b` (re-review) **immediately below this row** |
| **high** | Actionable mapping wrong on spot-check → inline AUDIT fix if ≤5 lines; else spawn 02a/02b |
| **medium** | Missing P24 delta row, weak task seeds → spawn unless trivial one-line AUDIT fix |
| **low / nit** | Note only in review Notes |

Spawned prompts must follow agent-loop-protocol skeleton (metadata, objective, exit criteria). Insert board rows in `docs/TODO/phase-27.md` directly under order 455.

## Findings template (use in Notes)

```text
| Severity | Finding | Evidence | Disposition |
|----------|---------|----------|-------------|
| blocker  | …       | path:Lnn | spawn 02a |
| high     | …       | …        | inline fix / spawn |
```

## Exit criteria

- [ ] AUDIT verified against live repo spot-checks (commands above)
- [ ] No open blocker/high without pending follow-up row (02a/02b)
- [ ] Confidence **medium** or **high** with evidence cited in Notes
- [ ] Verdict recorded: `APPROVE` or `REQUEST_CHANGES`
- [ ] Next runnable: **P27-S01-00** (only after APPROVE)

## Verdict

`APPROVE` | `REQUEST_CHANGES`

**APPROVE** only when INT-07/08/10 mapping is actionable for S01/S02 planners and P25-3 residual is correctly anchored to Phase 26 evidence.

## Todo updates

Status + notes on **P27-S00-02** only.

## Next

`P27-S01-00`
