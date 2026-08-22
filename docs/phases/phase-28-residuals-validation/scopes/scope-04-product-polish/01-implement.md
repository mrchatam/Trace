# P28-S04-01 — Product polish implementer

## Metadata
- id: P28-S04-01
- todo_ids: [P28-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Close **R4** (BLOCKING orphan duplicate honesty message) and **R5** (P25-4 attestation automation). Keep scope small: remove the redundant store pass, wire arm-matched env attestation in the harness, add a regression test. **Do not** implement R6 (thin-but-strict stderr hint — deferred).

**No** daemon/HTTP. **Do not** reopen R1–R3/R8.

## References

- [00-PLANNER.md](00-PLANNER.md) — locked defaults (this scope)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R4/R5 seeds
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Live anchors:
  - `cmd/trace/seed.go` `collectExportGraphHonestyViolations` — document pass + store BLOCKING loop (~L132–172)
  - `internal/domain/seed_export_honesty.go` L43–48 — orphan discovery messages
  - `experiments/ab-p25-gap-pass-validation/score.sh` L218 — `skip "P25-4…"`
  - `experiments/ab-p25-gap-pass-validation/PROTOCOL.md` L74 — manual attestation

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P28-S04-00 — do not re-debate)

| Item | Value |
|------|-------|
| **R4 strategy** | **Remove** the store-backed BLOCKING orphan loop in `collectExportGraphHonestyViolations` (`cmd/trace/seed.go`, currently ~L142–171 after `ListDiscoveries`) |
| **R4 sole source** | `domain.CollectSeedDocumentHonestyViolations` — keep existing generic orphan messages |
| **R4 must not** | Add `Severity` to portable `SeedEntity`; keep dual message paths; invent merge/dedupe maps while leaving both loops |
| **R4 test** | BLOCKING orphan discovery → **exactly one** honesty violation mentioning that discovery ID (prefer `cmd/trace` `--strict` stderr count; domain unit OK if it proves single-source after store removal) |
| **R5 mechanism** | Env vars (not RESULTS.md parser) |
| **R5 build** | `P25_ATTEST_BUILD=Y` **and** `--arm build` → **pass** P25-4 |
| **R5 directed** | `P25_ATTEST_DIRECTED=Y` **and** `--arm directed` → **pass** P25-4 |
| **R5 unset** | Env unset for current arm → keep **skip** (backward compatible) |
| **R5 wrong-arm** | Ignore mismatched env (e.g. `P25_ATTEST_DIRECTED=Y` while scoring `--arm build` does **not** pass build attestation) |
| **R5 docs** | Update `PROTOCOL.md` P25-4: set env **before** `./score.sh`; RESULTS.md template remains human narrative; optional one-liner in `RUBRIC.md` |
| **R6** | **Defer** — do not add thin-but-strict-pass stderr hint |
| **Out of scope** | Daemon/HTTP; portable graph severity on discoveries; reopening R1–R3/R8; RESULTS.md auto-parser |

```text
collectExportGraphHonestyViolations
  → CollectSeedDocumentHonestyViolations(doc) only
  → (removed) store ListDiscoveries BLOCKING re-check

score.sh --p25
  → arm=build  + P25_ATTEST_BUILD=Y     → pass P25-4
  → arm=directed + P25_ATTEST_DIRECTED=Y → pass P25-4
  → else                                 → skip P25-4
```

## Preflight

Run from repo root before editing:

```bash
cd /home/ali/Desktop/Trace

test -f docs/phases/phase-28-residuals-validation/scopes/scope-04-product-polish/00-PLANNER.md
test -f docs/phases/phase-28-residuals-validation/scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md

# R4 live dup anchors (store loop must still exist until you remove it)
sed -n '132,172p' cmd/trace/seed.go
sed -n '43,48p' internal/domain/seed_export_honesty.go
grep -n 'BLOCKING discovery\|CollectSeedDocumentHonestyViolations\|ListDiscoveries' cmd/trace/seed.go

# R5 skip anchor
sed -n '196,220p' experiments/ab-p25-gap-pass-validation/score.sh
grep -n 'P25-4\|attestation' experiments/ab-p25-gap-pass-validation/PROTOCOL.md

# Baseline
GOPROXY=direct go test ./internal/domain/ -count=1 -run 'SeedDocumentHonesty'
GOPROXY=direct go test ./cmd/trace/ -count=1 -run 'SeedExport|Enforce|ThinGraph'
```

Abort if `collectExportGraphHonestyViolations` or score.sh P25-4 skip line is missing.

## Files to touch

| File | Change |
|------|--------|
| `cmd/trace/seed.go` | In `collectExportGraphHonestyViolations`: keep the document-honesty loop; **delete** store `ListDiscoveries` + BLOCKING orphan append (~L142–171). Function may no longer need store list if unused — drop unused imports if any. |
| `cmd/trace/*_test.go` and/or `internal/domain/seed_export_honesty_test.go` | Regression: BLOCKING orphan → **one** violation for that ID (not two). Prefer integration under `cmd/trace` that creates a BLOCKING discovery without `discovery_mentions_task`, runs export `--strict`, asserts stderr mentions the ID once. |
| `experiments/ab-p25-gap-pass-validation/score.sh` | Replace unconditional `skip "P25-4…"` with arm-matched env check (`P25_ATTEST_BUILD` / `P25_ATTEST_DIRECTED`). |
| `experiments/ab-p25-gap-pass-validation/PROTOCOL.md` | Document env attestation before score; keep RESULTS.md narrative template. |
| `experiments/ab-p25-gap-pass-validation/RUBRIC.md` | Optional one-liner pointing at env vars. |

Do **not** change `internal/domain/seed_export.go` `SeedEntity` shape. Do **not** implement R6 hints.

## Minimal todos

1. **Preflight** — run bash block; record pass in board Notes.
2. **R4 remove store loop** — `collectExportGraphHonestyViolations` returns only document honesty violations.
3. **R4 regression test** — BLOCKING orphan → exactly one honesty line for that discovery ID; INFO orphan still one line (unchanged).
4. **R5 score.sh** — arm-matched `P25_ATTEST_*=Y` → `pass`; unset → `skip`; wrong-arm ignored.
5. **R5 docs** — PROTOCOL.md (required); RUBRIC.md one-liner (optional).
6. **Regression** — `GOPROXY=direct go test ./internal/... ./cmd/trace/... -count=1` PASS.
7. **Board** — **P28-S04-01** status + notes only (R4/R5 evidence; R6 not done).

## Test commands

```bash
cd /home/ali/Desktop/Trace
GOPROXY=direct go test ./internal/domain/ -count=1 -run 'SeedDocumentHonesty'
GOPROXY=direct go test ./cmd/trace/ -count=1 -run 'SeedExport|Enforce|Honesty|ThinGraph|Orphan|Blocking'
GOPROXY=direct go test ./internal/... ./cmd/trace/... -count=1

# Optional harness smoke (no full E02 run required):
# P25_ATTEST_DIRECTED=Y ./experiments/ab-p25-gap-pass-validation/score.sh G1 --p25 --arm directed
# expect a pass line for P25-4 when graph/rules already satisfy other checks (or at least no unconditional skip)
```

## Exit criteria

- [ ] Store-backed BLOCKING orphan loop **gone** from `collectExportGraphHonestyViolations`
- [ ] Orphan discoveries reported once via `CollectSeedDocumentHonestyViolations` only
- [ ] Regression test: BLOCKING orphan → exactly one violation for that ID
- [ ] `score.sh`: `P25_ATTEST_BUILD=Y` / `P25_ATTEST_DIRECTED=Y` arm-matched → pass P25-4; unset → skip; wrong-arm ignored
- [ ] `PROTOCOL.md` documents env attestation
- [ ] R6 **not** implemented
- [ ] `GOPROXY=direct go test ./internal/... ./cmd/trace/... -count=1` PASS
- [ ] Board **P28-S04-01** Notes cite evidence; next **P28-S04-02**

## Todo updates

Status + notes on **P28-S04-01** only.

## Next

`P28-S04-02`
