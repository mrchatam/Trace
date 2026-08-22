# P28-S06-03 — FM-02 / FR-P28-02 implementer

## Metadata
- id: P28-S06-03
- todo_ids: [P28-S06-03]
- role: implementer
- skills: [incremental-implementation, context-engineering]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

**FR-P28-02 / FM-02:** Close pre-export skip — agents write discoveries/decisions **before** `seed export --strict --enforce`, not only at gate time. Prefer stronger gap-pass / write-before-edit harness nudges; optional early thin-graph warn (not a silent weaken of enforce).

## References

- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — FM-02
- INT-07 honesty: `internal/domain/seed_export_honesty.go`; gap-pass: `internal/install/gappass.go`
- Harness: `experiments/ab-p25-gap-pass-validation/` PROTOCOL / score.sh
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start.

## Acceptance hint

Directed or build arm: consecutive session shows disc/dec writes **before** export; `--strict --enforce` still blocks thin export (regression green).

## Preflight

```bash
cd /home/ali/Desktop/Trace
grep -n 'GapPass\|write-before\|thin' internal/install/gappass.go experiments/ab-p25-gap-pass-validation/PROTOCOL.md | head
GOPROXY=direct go test ./cmd/trace/ -count=1 -run 'SeedExport|Enforce|Thin'
```

## Suggested work

1. Strengthen harness/product nudge for write-before-export.
2. Optional early warn on thin graph (warn-only; enforce remains fail-closed).
3. Evidence: notes + regression that thin export still fails enforce.
4. Deliver `FM02-NOTES.md` if behavior is dogfood/harness-only.

## Out of scope

- Removing `--strict --enforce`; FM-01/04+; daemon/HTTP

## Exit criteria

- [ ] Acceptance hint met; enforce regression green
- [ ] Next runnable **P28-S06-04**

## Todo updates

Status + notes on **P28-S06-03** only.

## Next

`P28-S06-04`
