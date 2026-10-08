# S06 — Guards — SCOPE-TODOS

Test-only scope; no issues of its own. Protects #149–#162 from recurrence.

| ID | Item | Status |
|----|------|--------|
| S06-0 | P45-S06-00 thicken 01/02 | **done** 2026-10-04 |
| S06-1 | Key-shape sweep (21 calls) + protocol conformance suite | **done** 2026-10-04 (P45-S06-01) |
| S06-2 | Verification record (self-review) | **done** 2026-10-04 **PASS** (P45-S06-02) |

## Acceptance

- [x] Zero uppercase object keys across 21 real tool responses (incl. scope-member / api-contract / regression paths)
- [x] Guard self-keeper proves it can fail
- [x] Real JSON-RPC: negotiation, locked tool set, ping, unknown method, unknown params (`IsError`), out-of-enum (`IsError`), concurrent readers
- [x] `jsonschema-go` promoted to a direct dependency via `go mod tidy`