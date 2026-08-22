# P00-S05-02 — Scope review notes (2026-08-15)

Independent review of S05 against `01-causal.md` + TODO Notes for `P00-S05-01`. Fresh session; claims verified in-repo.

## Plan (executed)

1. Diff claims vs `internal/domain` + store mig `003` + causal Upsert/link APIs
2. Re-run `CGO_ENABLED=0` domain+store (+ vcs/gitcli) and full `./...` with CGO for analyzers
3. Severity-tag findings; fix/spawn only for blocker/high
4. Write these notes; mark board + SCOPE-TODOS; light thicken upcoming S06/S07 Depends if helpful

## Verdict

**APPROVE** — no blocker/high. Confidence: **high**. Spawns: **none**. Next board row: **P00-S06-00**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| `domain.New(*store.Store)` only; no `store.Open` / second DB | Pass (`service.go`; imports: store only among product pkgs) |
| Migration `003_causal_domain.sql`: `tasks.work_state` + `entity_links` + indexes | Pass (embedded via `schema/*.sql`; store test reads column) |
| Provenance on Create* (G5): default `USER_ASSERTED`, `ACTIVE`, empty title rejected | Pass (`create.go` + `TestCreateRoundtripProvenance`) |
| `work_state` ≠ provenance `status` | Pass (separate column + Task fields; defaults PENDING / ACTIVE) |
| Goal→Task via `goal_id` (not duplicate `entity_links`) | Pass (`LinkGoalTask` + `TestGoalTaskLinkViaGoalID`) |
| Decision→Task `decision_affects_task`; Discovery→PlanChange `discovery_causes_plan_change` | Pass (link tests) |
| Transition graph + illegal fails closed; `task.transition` payload actor/from/to/reason/`evidence_ids` | Pass (`task_state.go` + `TestTransitionLegalAndIllegal`) |
| DONE policy stub (flag or EvidenceIDs) | Pass (`TestDonePolicyStub`) |
| MarkStale → provenance STALE + event | Pass (`TestMarkStale`; event type `entity.stale`) |
| Claim/Evidence stubs + `claim_has_evidence`; no promotion | Pass (`claim_stub.go` + smoke test) |
| Events on create/link/transition; append-only store Events | Pass |
| No source BLOBs; no MCP/daemon/HTTP; no new CLI | Pass (`cmd/trace` still help/version only) |
| Domain does not import `analyzers` / `gitcli` | Pass (`go list` imports) |
| `CGO_ENABLED=0` domain+store | Pass (re-run below) |
| Cross-scope S06/S07 notes still accurate | Pass (Depends already name domain Create*/Link*/`work_state`; lightly thickened) |

## Re-verification commands (2026-08-15)

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... -count=1   # ok
CGO_ENABLED=0 go test ./internal/vcs/... ./internal/gitcli/... -count=1     # ok
CGO_ENABLED=1 go test ./... -count=1                                       # ok
```

Domain tests: CreateRoundtripProvenance, GoalTaskLinkViaGoalID, DecisionAffectsTaskLink, DiscoveryCausesPlanChangeLink, TransitionLegalAndIllegal, DonePolicyStub, MarkStale, ClaimEvidenceStub.

## Findings

### Blocker / high

None.

### Medium (residual — no spawn)

- **Non-atomic mutate-then-event:** `TransitionTask` / links / creates Upsert or Insert then `AppendEvent`. If the event write fails, durable state can advance without the matching event. Acceptable for P0 thin events; transactional wrap is later harden (same class as S04 IndexFile residual).
- **CreateTask accepts caller `WorkState`:** DONE policy and the legal graph apply only on `TransitionTask`. A caller can create a task already `DONE` (or an illegal vocabulary string) and skip the stub. Locked API exposes `TaskInput.WorkState` with default PENDING — intentional for seed/fixture, but S07/S08 should prefer transitions (or empty WorkState) for honest provenance of DONE.

### Low / nit

- `MarkStale` emits `entity.stale` (string literal), not one of the three locked create/link/transition type constants — fine per MarkStale wording; consider a package const for why/FTS consumers.
- Several APIs ignore `ctx` (`_ = ctx`) — fine until cancellation is wired.
- `LinkMeta` on `LinkGoalTask` is defaulted then unused for persistence (Goal→Task lives on `goal_id` only) — matches lock; `_ = meta` is honest.

## Spawns

None.

## Residual risks

- No Review/VerifiedFact / Claim promotion (explicitly Phase 01).
- Decision workflow enums (`PROPOSED`/`ACCEPTED`/`REJECTED`) deferred; Decisions use provenance ACTIVE.
- Duplicate `entity_links` UNIQUE fails closed at store — domain does not soft-dedupe.
- `EvidenceIDs` on DONE are recorded in the event payload only (not validated as existing evidence rows).

## Forward edits this review

- `SCOPE-TODOS.md` — mark S05-01/02 done
- `docs/TODO.md` — `P00-S05-02` → done + notes
- `scopes/scope-06-retrieval-context/01-retrieval.md` — light API note: ListLinks*, event types, MarkStale
- `scopes/scope-07-cli/01-cli.md` — light note: prefer TransitionTask for DONE; Link* helpers named
