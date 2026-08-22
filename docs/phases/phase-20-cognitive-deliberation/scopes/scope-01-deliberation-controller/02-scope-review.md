# P20-S01-02 — Review deliberation controller

## Metadata
- id: P20-S01-02
- todo_ids: [P20-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent fresh-session review: policies are deterministic and table-tested; EXECUTE cannot proceed on blocking uncertainty; hop budget fail-closed; P19 loop tests still green; no raw CoT storage; events auditable.

## Keeper tests (must re-run)

```bash
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
```

## Review checklist

- [ ] SelectNext is pure/deterministic — same inputs → same phase
- [ ] Table tests cover every row locked in S01-00
- [ ] `deliberation.transition` payload includes inputs + chosen_phase + reason
- [ ] Hop budget fail-closed — no silent infinite ORIENT↔INVESTIGATE
- [ ] No modification to P19 schema contracts in S01
- [ ] No source blobs / no raw CoT in persistence
- [ ] Law 19: CLI not required for S01 unless explicitly scoped

## Exit criteria

- blocker/high fixed or spawned forward (implement+review pair immediately below this row)
- confidence medium or high with evidence
- residuals listed explicitly if medium
