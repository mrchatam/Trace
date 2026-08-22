# P28-S04-02 — Product polish independent review

## Metadata
- id: P28-S04-02
- todo_ids: [P28-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review that **R4** and **R5** are closed by repo evidence (not Notes alone): single honesty message source for orphan discoveries; P25-4 env attestation wired and documented. Fresh subagent — do **not** reuse P28-S04-01 session.

## References

- [01-implement.md](01-implement.md) — locked R4/R5
- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R4/R5
- Live: `cmd/trace/seed.go`, `internal/domain/seed_export_honesty.go`, `experiments/ab-p25-gap-pass-validation/score.sh`, `PROTOCOL.md`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked review policy

| Item | Rule |
|------|------|
| R4 | Store-backed BLOCKING orphan loop **must be absent**; document honesty is sole orphan source |
| R4 | No `Severity` field added to portable `SeedEntity` |
| R4 | Regression proves BLOCKING orphan → **exactly one** violation for that ID |
| R5 | `P25_ATTEST_BUILD=Y` / `P25_ATTEST_DIRECTED=Y` arm-matched → pass; unset → skip; wrong-arm ignored |
| R5 | PROTOCOL.md documents env path (RUBRIC one-liner optional) |
| R6 | Must **not** be implemented in this scope |
| Out of scope creep | Daemon/HTTP; RESULTS.md parser; reopening R1–R3/R8 → blocker |
| Close R4/R5 | Only with passing tests + live code/docs matching locks |

## Verify checklist

### R4 — BLOCKING dup removed

- [ ] `collectExportGraphHonestyViolations` no longer calls `ListDiscoveries` for BLOCKING re-check
- [ ] No remaining `"BLOCKING discovery %s missing discovery_mentions_task"` append in that function
- [ ] `CollectSeedDocumentHonestyViolations` still flags orphan discoveries
- [ ] Regression test exists and asserts **one** violation (not two) for a BLOCKING orphan ID
- [ ] `SeedEntity` unchanged (no severity on portable discovery export shape)

### R5 — P25-4 attestation

- [ ] `score.sh` no longer unconditionally `skip "P25-4…"`
- [ ] Build arm + `P25_ATTEST_BUILD=Y` → pass path
- [ ] Directed arm + `P25_ATTEST_DIRECTED=Y` → pass path
- [ ] Unset → skip
- [ ] Wrong-arm env does not falsely pass
- [ ] `PROTOCOL.md` documents env attestation before score

### Process / blast radius

- [ ] Reviewer re-runs `GOPROXY=direct go test ./internal/... ./cmd/trace/... -count=1` PASS
- [ ] R6 thin-but-strict hint **not** added
- [ ] No daemon/HTTP; no R1–R3/R8 rework
- [ ] S05 not silently started

### Live spot-checks (reviewer runs)

```bash
REPO=/home/ali/Desktop/Trace
cd "$REPO"

# R4 — store loop gone
sed -n '132,175p' cmd/trace/seed.go
# Must NOT still contain BLOCKING discovery missing discovery_mentions_task in this function
! grep -n 'BLOCKING discovery' cmd/trace/seed.go || true
grep -n 'CollectSeedDocumentHonestyViolations' cmd/trace/seed.go

# SeedEntity still bare (id/title/body)
grep -A6 'type SeedEntity struct' internal/domain/seed_export.go

# R5 — env attestation
sed -n '196,230p' experiments/ab-p25-gap-pass-validation/score.sh
grep -n 'P25_ATTEST_BUILD\|P25_ATTEST_DIRECTED\|P25-4' experiments/ab-p25-gap-pass-validation/score.sh
grep -n 'P25_ATTEST\|P25-4' experiments/ab-p25-gap-pass-validation/PROTOCOL.md

# Re-verify
GOPROXY=direct go test ./internal/... ./cmd/trace/... -count=1
```

## Findings severity

| Level | Action |
|-------|--------|
| blocker | Store BLOCKING loop still present; dual messages still emit; Severity added to SeedEntity; P25-4 still unconditional skip; wrong-arm env passes; daemon/HTTP |
| high | No regression for single-message; PROTOCOL silent on env; tests only grep without exercising orphan path |
| medium | RUBRIC stale vs PROTOCOL; skip message wording unclear |
| low / nit | Comment polish; optional RUBRIC one-liner missing |

## Spawn policy

- blocker/high → insert `P28-S04-02a` implement + `P28-S04-02b` review immediately below this row with full prompts
- medium → prefer spawn unless trivial inline fix
- Do **not** rewrite `done` P28-S04-00/01 history

## Exit criteria

- [ ] R4 and R5 closed with evidence (or spawn filed)
- [ ] Confidence **high**, or **medium** with explicit residual risks
- [ ] Board **P28-S04-02** Notes; next **P28-S05-00** if APPROVE
- [ ] Optionally update RESIDUAL-AUDIT R4/R5 disposition to closed (forward-only note OK)

## Todo updates

Status + notes on **P28-S04-02** only; spawn rows if needed.

## Next

`P28-S05-00`
