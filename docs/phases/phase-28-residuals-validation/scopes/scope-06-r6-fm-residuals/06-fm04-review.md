# P28-S06-06 — FM-04 / FR-P28-03 reviewer

## Metadata
- id: P28-S06-06
- todo_ids: [P28-S06-06]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening]
- verification: mixed
- hooks: []

## Objective

Fresh-session review of P28-S06-05 (FM-04). Confirm parent/worker Trace gap closed or honestly documented; Option A intact; Option B not sneak-shipped.

## Checklist

- [ ] Parent cannot silently offload graph while editing without task (or limits documented)
- [ ] Option A strict+empty deny still green
- [ ] No Option B; no daemon/HTTP
- [ ] Tests/docs evidence adequate

## Exit criteria

Verdict in Notes; next **P28-S06-07** on APPROVE.

## Next

`P28-S06-07`
