# P11-S02-00 — Review PASS+FAIL / operator identity (FINAL)

## Metadata
- id: P11-S02-00
- todo_ids: [P11-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S02 implement/review prompts for **DF-43, DF-44**. Confirm live inventory; lock APIs/tests. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19; claims ≠ evidence
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A5 DF-44 posture choice
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-43, DF-44
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md) — adv8 PASS+FAIL; adv1 flag≠identity
- Phase 10 S04 FINAL: [../../../phase-10-integrity-surfaces/scopes/scope-04-operator-capability-gates/00-PLANNER.md](../../../phase-10-integrity-surfaces/scopes/scope-04-operator-capability-gates/00-PLANNER.md) — DF-17 operator flag; Gate G hatch; Actor≠auth
- Live: `internal/domain/task_state.go` (`findPassReviewID`, `TransitionTask`); `cmd/trace/{transition,help}.go`; `internal/mcp/{server,tools_write}.go`; `evals/honesty/honesty_test.go` Path B/C
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no grill (phase A5 + residual-hardening A1/A3 + live inventory + honesty Path C conflict resolved by explicit supersession).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `findPassReviewID` | First linked `review_judges_task` with `result=PASS` — **ignores sibling FAIL** → DF-43 |
| `TransitionTask` →DONE | Hatch **or** (`findPassReviewID` nonempty ∧ `AllowOperatorDone`); no FAIL scan |
| Honesty Path C | FAIL review left linked + still FAIL, then sibling PASS + `AllowOperatorDone` → DONE **succeeds** (encodes the DF-43 hole) |
| DF-17 | `--as-operator` / MCP `as_operator` freestanding; `Actor` never auth — **holds**; residual freeness → DF-44 |
| Help / MCP | Partial Actor≠auth wording already present; does **not** say flag≠verified identity |
| Gate G | `--allow-done` / `AllowDoneWithoutReview` bypasses PASS + operator — **retain** |
| Migration | None needed for review results enum or links |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-43 | Sibling FAIL must block DONE even if another PASS exists | Before authorizing via PASS: any linked `review_judges_task` with current `result=FAIL` → reject; UNCERTAIN/empty do not block |
| DF-44 | Flag≠identity residual | Keep **conscious-flag** design (no OAuth/identity); close via explicit help/MCP schema honesty (+ optional usage note) |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF home | **DF-43, DF-44** only (DF-17/18/24/26 remain as shipped; DF-51 hatch-vs-caps → S04 — do not expand here) |
| Packages | **`internal/domain`** (`TransitionTask` / `findPassReviewID` or sibling helper); thin **`cmd/trace`** (`transition`, `help`); thin **`internal/mcp`** (`trace_transition` desc + `as_operator` schema). **G19** — no business logic in adapters |
| Migration | **None** |
| DF-43 rule | →DONE (non-hatch) requires: **no** linked `review_judges_task` review with current `result=FAIL`, **and** nonempty PASS id (existing `findPassReviewID`), **and** `AllowOperatorDone`. FAIL blocks even when a sibling PASS exists. **Not** latest-wins-by-timestamp |
| DF-43 active set | Only current result on linked reviews. **UNCERTAIN** and unset/empty do **not** block. Scope-only `review_judges_scope` FAIL does not authorize or block task DONE (unchanged) |
| DF-43 reject reason | Assertable string mentioning linked/sibling **FAIL** (and existing PASS / `--as-operator` / `--allow-done` guidance as appropriate) |
| DF-43 hatch | `AllowDoneWithoutReview` **bypasses** FAIL check as well as PASS+operator (Gate G retained) |
| DF-43 recovery | Clearing a blocking FAIL is **explicit** (`SetReviewResult` → `UNCERTAIN` or other non-FAIL) — do **not** auto-invalidate FAIL when a new PASS is created |
| DF-43 honesty supersession | **Path B** unchanged (FAIL alone rejects). **Path C:** after FAIL, remediation must **supersede** the FAIL (e.g. set prior FAIL → `UNCERTAIN` with reason like `superseded by later review`) **before** PASS + `AllowOperatorDone` → DONE. Path C must **not** leave an active FAIL linked while promoting. Optional assert: FAIL→UNCERTAIN then PASS, DONE OK |
| DF-44 posture | **Keep conscious flag** (phase A5 pick). `--as-operator` / `as_operator` remains freestanding; **Actor string never authorization** (DF-17 holds) |
| DF-44 close | Help + transition usage note + MCP `trace_transition` description / `as_operator` jsonschema must state clearly: flag is a **conscious claim**, **not** verified operator identity / not authn. Prefer assertable phrases (`flag≠identity` or `not verified identity` / `conscious` + `Actor` ≠ auth). Optional one-line usage note OK; **do not** require a WARNING on every successful `--as-operator` DONE (avoid colliding with DF-26 hatch WARNING) |
| DF-44 non-goals | No OAuth/OIDC/tokens; no actor-name allowlist; no env secret gate; do not remove freestanding flag; do not trust `Actor=="operator"` |
| Check order in `TransitionTask` | Keep existing: actor+reason → legal edge → DF-24 caps → →DONE: hatch **or** (**no linked FAIL** ∧ PASS ∧ operator). Implement FAIL scan in domain only |
| Tests (required) | (1) **`TestSiblingFailBlocksDone`** (or equiv): linked FAIL+PASS + `AllowOperatorDone` → reject; reason mentions FAIL. (2) PASS alone + flag → OK. (3) PASS+UNCERTAIN + flag → OK. (4) Hatch with FAIL present → OK (Gate G). (5) Honesty A/B/C green after Path C supersession. (6) DF-44: help and/or MCP schema/desc contain identity/conscious wording; keep `TestOperatorDoneRequiresFlag` (Actor≠auth) |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` untouched; P10 DF-17/18/24/26; P11-S01 DF-40 |
| Forbidden | New mig; OAuth/real identity; latest-wins-only without FAIL block; auto-clear FAIL on PASS; removing Gate G hatch; trusting Actor; daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–10 / P11-S01 `done` history; Mode-B Gate C pack rewrite; S03+ product work; DF-51 scope creep |

## Effects on later scopes
- **S03** (store lock): no review-gate coupling — serial after S02 review only. Light Depends note already on S03 stubs.
- **S04** (DF-51): hatch WARNING may later mention missing-caps; this scope must not change hatch↔caps independence.
- **S08 VERIFY:** include DF-43 sibling FAIL+PASS reject + DF-44 conscious-flag copy + honesty Path C supersession in evidence table.

## Exit
- [x] Thicken `01-review-pass-fail-operator.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S03 already Depends-on S02-02)
- [x] Board Notes; next **P11-S02-01**
- [x] Product Go — **not** this row
