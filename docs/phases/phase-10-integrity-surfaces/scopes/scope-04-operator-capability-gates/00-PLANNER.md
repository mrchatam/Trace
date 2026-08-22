# P10-S04-00 — Operator + capability gates (FINAL)

## Metadata
- id: P10-S04-00
- todo_ids: [P10-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S04 implement/review prompts for **DF-17, DF-18, DF-24, DF-26, DF-31**. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19; claims ≠ evidence
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A2 operator DONE; A3 reopen invalidates PASS; A4 fail-closed cap gate
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-17/18/24/26/31
- [experiments/ab-operator-gate/results/G1/PROBE-AGENT-DONE.md](../../../../../experiments/ab-operator-gate/results/G1/PROBE-AGENT-DONE.md) — acceptance evidence shape
- S01–S03 inherit (do not re-litigate): why/Exact; nine MCP tools incl. `trace_capability`/`trace_transition`; index GC
- Live: `internal/domain/{task_state,service,capability,review}.go`; `cmd/trace/{transition,capability,seed,help}.go`; `internal/mcp/tools_write.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no user grill required (A2–A4 + live inventory + dogfood probes).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `TransitionOptions` | `Actor`, `Reason`, `EvidenceIDs`, `AllowDoneWithoutReview` only — **no** operator flag; **no** cap override |
| `TransitionTask` → DONE | Linked Review `PASS` **or** hatch; **any** actor (cli/mcp/agent) may promote after PASS (**DF-17**) |
| `findPassReviewID` | Any linked `review_judges_task` with `result=PASS` — survives DONE→PENDING/STALE reopen (**DF-18**) |
| Legal graph | `DONE` → `STALE`\|`PENDING` only (not direct `IN_PROGRESS`) |
| Capability | `MissingCapabilities` exists; **never consulted** in `TransitionTask` (**DF-24**) |
| CLI `--allow-done` | Quiet success JSON; usage note only (**DF-26**) |
| `capability missing` | Empty `--task` → bare usage line; no `trace tasks` hint (**DF-31**) |
| MCP | Nine tools (S02); `trace_transition` mirrors CLI; `as_operator` / `allow_missing_caps` **absent** |
| Honesty Path C | `Actor: "implementer"` + PASS → DONE (no operator flag) — **must update under explicit supersession below** |
| Gate G | Escape = `AllowDoneWithoutReview` without PASS — **keep hatch**; louder UX only |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-17 | Operator owns DONE (not text-only) | Explicit **`AllowOperatorDone`** flag (CLI `--as-operator`, MCP `as_operator`); **do not** trust `Actor` string |
| DF-18 | Sticky PASS after reopen | Leaving **DONE** invalidates linked task PASS reviews (→ `UNCERTAIN`); new PASS required |
| DF-24 | Missing caps gate transitions | Fail-closed: any transition blocked when `MissingCapabilities` nonempty unless override |
| DF-26 | `--allow-done` footgun | Keep hatch for Gate G; **loud** CLI stderr + MCP `warning` on use |
| DF-31 | `capability missing` UX | Clear usage + hint to `trace tasks` / require `task_id` on MCP |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Packages | **`internal/domain`** (gates in `TransitionTask` + helpers); thin **`cmd/trace`** (`transition`, `capability`, `seed`, `help`); thin **`internal/mcp`** (`trace_transition` / `trace_capability` params only). **G19** — no business logic in adapters |
| Migration | **None** (prefer). Invalidate via existing `SetReviewResult` → `UNCERTAIN`. Do **not** add `011_*` unless a blocker forces it (then Notes + reviewer) |
| Check order in `TransitionTask` | (1) actor+reason non-empty (2) legal edge (3) **DF-24** missing-cap gate (4) if →DONE: hatch **or** (PASS + operator flag) |
| DF-17 API | `TransitionOptions.AllowOperatorDone bool`. CLI `--as-operator`. MCP `as_operator`. Seed JSON `as_operator` (bool; mirror `allow_done`). **Actor string is never authorization.** |
| DF-17 DONE rule | →`DONE` allowed iff `AllowDoneWithoutReview` **OR** (`findPassReviewID` nonempty **AND** `AllowOperatorDone`). EvidenceIDs alone still insufficient. |
| DF-17 reject reason | Must mention Review PASS / `--as-operator` (or `AllowOperatorDone`) / `--allow-done` (hatch) as appropriate — keep assertable strings for honesty/domain tests |
| DF-17 honesty supersession | Paths A/B unchanged (reject). **Path C:** after PASS, DONE must set `AllowOperatorDone: true` (actor may stay `"implementer"` or become `"operator"` — flag is what matters). Gate **semantics** (PASS authorizes DONE; FAIL/EvidenceIDs do not) preserved; **who may promote** is now explicit operator flag |
| DF-17 Gate G | Escape path still uses `AllowDoneWithoutReview: true` **without** requiring `AllowOperatorDone` (hatch bypasses both review and operator) |
| DF-17 non-goals | No authn/OIDC; no allowlist of actor names; no daemon; no graph edge `DONE→IN_PROGRESS` this scope |
| DF-18 invalidate | On any transition **from** `DONE` → allowed next (`PENDING` or `STALE`): for each linked `review_judges_task` review with `result=PASS`, set result **`UNCERTAIN`** (reuse `SetReviewResult` or internal helper) with reason like `invalidated on reopen (DONE→…)`; actor = transition `opts.Actor` (or `"system"` if helper). Append normal `review.result` event. Then apply work_state change. |
| DF-18 DONE auth | Subsequent →DONE needs a **new** PASS (old invalidated reviews must not satisfy `findPassReviewID`) |
| DF-18 non-goals | Do not delete reviews/links; do not invent `INVALIDATED` result enum; scope-only `review_judges_scope` PASS does not authorize task DONE (unchanged) |
| DF-24 gate | Before mutating state: if `MissingCapabilities(taskID)` returns nonempty → `ErrInvalidTransition` (or typed validation) unless `AllowMissingCapabilities` |
| DF-24 override | `TransitionOptions.AllowMissingCapabilities bool`; CLI `--allow-missing-caps`; MCP `allow_missing_caps` |
| DF-24 scope | **All** transitions (not only →DONE). Empty requirements ⇒ no block. UNKNOWN/UNAVAILABLE already in `MissingCapabilities` |
| DF-24 MCP/CLI | Same domain path (G19); MCP `trace_transition` must pass both new flags |
| DF-26 louder | On **successful** transition with `AllowDoneWithoutReview`: CLI **stderr** WARNING (must include `allow-done` / escape / bypass wording); MCP success JSON includes `"warning":"<non-empty>"`. Refresh `help.go` + MCP tool description. Do **not** remove hatch |
| DF-26 non-goals | No double-confirm flag; no env var required; seed `allow_done` may stay quiet or share warning — prefer warning on CLI/MCP interactive paths |
| DF-31 CLI | Missing `--task`: usage + hint e.g. `task id required; list tasks: trace tasks` (exit usage). Do not return empty/mysterious store errors |
| DF-31 MCP | `trace_capability` action `missing` without `task`/`task_id`: clear error naming required param + suggest `trace_tasks` |
| Tests (required) | Domain: operator DONE deny without flag; PASS+flag OK; hatch OK; reopen invalidates PASS then deny DONE until new PASS; missing cap blocks transition; override OK. CLI/MCP flag wiring + DF-26 warning. Honesty A/B/C (+ Gate G) green after Path C supersession. Capability missing usage test optional but preferred |
| ab-operator-gate | Probe shape: agent/MCP DONE after PASS **without** `as_operator` must fail; acceptance evidence for review Notes |
| Carry-forward | honesty A/B/C + Gate G; Gate E/F; ablation; Gate H; compat; p0x; x0; S01 why/Exact; S02 nine tools; S03 GC tests; Gate C `dry_run:false` intact |
| Forbidden | New mig (default); daemon/HTTP/embeddings; trusting Actor for DONE; weakening honesty A/B or removing Gate G hatch; plan/impact/index MCP; rewriting Phase 00–09 / S01–S03 `done` history; Mode-B Gate C pack rewrite |

## Effects on later scopes
- **S05 VERIFY:** must spot-check DF-17/18/24/26/31 + honesty Path C operator flag + Gate G hatch + nine MCP tools still present; light Depends note on S05 stubs.
- **Experiments:** ab-operator-gate PROBE becomes product-true after implement (not board-blocking).

## Exit
- [x] Thicken `01-operator-capability-gates.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light S05 Depends note
- [x] Board Notes; next **P10-S04-01**
- [x] Product Go — **not** this row
