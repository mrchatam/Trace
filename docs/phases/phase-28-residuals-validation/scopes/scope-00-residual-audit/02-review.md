# P28-S00-02 — Residual audit review

## Metadata
- id: P28-S00-02
- todo_ids: [P28-S00-02]
- role: reviewer
- skills: [code-review-and-quality, investigator]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

Independent review of `RESIDUAL-AUDIT.md` — confirm R1–R8 disposition and INT-01..11 mapping before S01 planning. **Fresh subagent** — do not reuse S00-01 session.

## References

- [01-residual-audit.md](01-residual-audit.md) — exit criteria + AUDIT template
- [00-PLANNER.md](00-PLANNER.md) — locked defaults
- [RESIDUAL-AUDIT.md](RESIDUAL-AUDIT.md) (produced by S00-01)
- [Phase 27 VERIFY-NOTES](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md)
- [Phase 28 README](../../README.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Verify checklist

### Deliverable presence

- [ ] `RESIDUAL-AUDIT.md` exists in this scope folder
- [ ] All template sections present: Executive summary, R1–R8 table, INT matrix, dogfood/harness/product, test coverage, closure plan, defers, risks
- [ ] Each of R1–R8 has explicit disposition

### Residual register accuracy

- [ ] R1 (P25-3b / Session-B) cites Phase 27 VERIFY-NOTES L59 “Directed without Session-B” + `PROMPT-G1-DIRECTED-GAP.md` exists
- [ ] R2/R3 (hook allow path) cites `enforcement.go` L106–108 (empty `TRACE_TASK_ID` → allow) — not failClosed
- [ ] R4 (BLOCKING dup) cites `seed.go` L158–170 store-backed BLOCKING check duplicating orphan link in `seed_export_honesty.go`
- [ ] R5 (P25-4) cites `score.sh` L218 skip attestation path
- [ ] R6 maps to INTERVENTION-MATRIX §3 FM rows still open after INT-01..11
- [ ] R7 cites concrete test files vs dogfood-only gaps (loop/deliberation/install/cmd tests)
- [ ] R8 (INT-11 drift) cites `cursorhook.go` `HookDriftNote` — manual only, no automated drift test

### INT matrix completeness

- [ ] INT-01..11 each have implementation + test coverage columns
- [ ] Closed INTs from P26/P27 not falsely marked missing
- [ ] Partial/residual INTs have file targets for S01–S04

### S01–S04 seeds

- [ ] Closure plan rows are actionable (file path + gap + scope)
- [ ] S02 seeds do not require prepare wipe
- [ ] S03 seeds reference strict enforce + missing TRACE_TASK_ID

### Live spot-checks (reviewer runs)

From repo root, confirm audit claims against live code:

```bash
cd /home/ali/Desktop/Trace

# Hook allow without task ID (R2/R3 — expect L106-108 allow branch)
grep -n 'TRACE_TASK_ID\|permission.*allow' internal/install/enforcement.go | head -20

# P25-3a/3b labels + P25-4 skip (R5 — expect L197-218)
grep -n 'P25-3\|P25-4' experiments/ab-p25-gap-pass-validation/score.sh

# Session-B prompt exists (R1)
test -f experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md

# Honesty / BLOCKING dup (R4 — store check seed.go L158-170 + doc orphan seed_export_honesty.go)
grep -n 'BLOCKING\|orphan\|graph honesty' internal/domain/seed_export_honesty.go cmd/trace/seed.go

# Saturation threshold (INT-02 — no internal/loop/saturation.go; use deliberation)
grep -n 'SaturationEmptyThreshold' internal/deliberation/types.go

# Hook drift doc (R8 / INT-11)
grep -n 'HookDriftNote\|drift' internal/install/cursorhook.go

# Promotion helper (INT-01)
grep -n 'PromoteBlocking\|BLOCKING' internal/domain/promote.go | head -10
```

**Path sanity:** If audit cites `internal/loop/saturation.go`, flag as **blocker** — live code uses `internal/deliberation/types.go`.

### Process

- [ ] No product code changed in S00
- [ ] No Session-B dogfood run claimed in S00

## Findings severity

| Level | Action |
|-------|--------|
| blocker | Spawn S00-02a implement + S00-02b review OR fix audit inline if trivial doc-only |
| high | Same as blocker |
| medium | Spawn unless trivial |
| low / nit | Note only |

## Spawn policy

- **HIGH/blocker** finding → insert `P28-S00-02a` (implement audit fix) + `P28-S00-02b` (review) immediately below this row
- **No HIGH/blocker** → APPROVE → next runnable **P28-S01-00**

## Exit criteria

- [ ] Verdict: APPROVE or spawn with pending follow-up
- [ ] Confidence: **high** (or **medium** with explicit residual risks listed)
- [ ] Board row P28-S00-02 status + notes updated

## Todo updates

Status + notes on **P28-S00-02** only.

## Next

**P28-S01-00** (unless spawn pending)
