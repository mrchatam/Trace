# P11 / S02 / 01 — Review PASS+FAIL / operator identity

## Metadata
- id: P11-S02-01
- todo_ids: [P11-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-43, DF-44** per sibling **00-PLANNER** FINAL locks (2026-08-16). Sibling linked FAIL must block →DONE even when another PASS exists; `--as-operator` stays a conscious flag with explicit flag≠identity docs (no real auth). **No new migration. Gate G hatch retained.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19; claims ≠ evidence
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-43, DF-44
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md)
- [phase README](../../README.md)
- Live: `internal/domain/task_state.go`; `cmd/trace/{transition,help}.go`; `internal/mcp/{server,tools_write}.go`; `evals/honesty/honesty_test.go`
- Prior: P10 S04 DF-17/18/24/26 — keep green
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. If 00-PLANNER is still DRAFT, stop and return to planner.

## Locked defaults (FINAL — P11-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/domain` (+ thin `cmd/trace` + `internal/mcp` copy only) |
| Migration | **None** |
| DF-43 | Non-hatch →DONE needs **no linked FAIL** ∧ PASS ∧ `AllowOperatorDone`; UNCERTAIN/empty do not block; not latest-wins |
| DF-43 hatch | `AllowDoneWithoutReview` bypasses FAIL as well |
| DF-43 recovery | Explicit `SetReviewResult` clear of FAIL (e.g. →UNCERTAIN); no auto-clear on new PASS |
| Honesty Path C | Supersede: clear prior FAIL before PASS+operator DONE |
| DF-44 | Keep freestanding `--as-operator` / `as_operator`; document flag≠identity / conscious claim; Actor≠auth |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false`; P10 operator gates; P11 DF-40 |
| Forbidden | OAuth/identity; auto-clear FAIL; remove hatch; trust Actor; mig; daemon/HTTP; board spawn |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/domain/task_state.go` | FAIL scan before PASS authorize (helper OK); reject reason mentions FAIL; hatch skips FAIL |
| `internal/domain/*_test.go` | `TestSiblingFailBlocksDone` (+ PASS-alone / PASS+UNCERTAIN / hatch+FAIL cases as locked) |
| `evals/honesty/honesty_test.go` | Path C: supersede FAIL→UNCERTAIN (or equiv) before PASS+DONE; A/B unchanged |
| `cmd/trace/help.go`, `transition.go` | DF-44 conscious-flag / flag≠identity wording in help + usage note |
| `internal/mcp/server.go`, `tools_write.go` | DF-44: `trace_transition` description + `as_operator` schema honesty |
| CLI/MCP tests (optional but preferred) | Assert copy phrases; keep Actor≠auth / `--as-operator` wiring green |

## Role work

1. TDD domain: FAIL+PASS + `AllowOperatorDone` rejects; PASS alone OK; PASS+UNCERTAIN OK; hatch with FAIL OK.
2. Wire FAIL gate in `TransitionTask` (domain only; G19).
3. Update honesty Path C supersession; keep Path A/B + Gate G.
4. DF-44: thicken help/MCP schema/desc with flag≠identity / conscious claim (no auth mechanism).
5. Run locked verify suite; board **status + Notes only** (cite test names + DF-43/44).

## Algorithm sketch (non-normative — locks above win)

```text
if to == DONE && !AllowDoneWithoutReview:
  if any linked review_judges_task has result==FAIL:
    reject (mention FAIL)
  passID = findPassReviewID(task)
  if passID == "": reject (need PASS)
  if !AllowOperatorDone: reject (need --as-operator)
# hatch path: skip FAIL + PASS + operator checks
```

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: sibling FAIL blocks; Path C after supersede works; hatch still DONE without PASS; Actor≠auth still holds; help/MCP mention flag≠identity.

## Exit criteria
- [ ] `TestSiblingFailBlocksDone` (or equiv) green — DF-43
- [ ] Honesty A/B/C + Gate G green after Path C supersession
- [ ] DF-44 help/MCP copy asserts conscious flag≠identity; no OAuth; freestanding flag retained
- [ ] P10 DF-17/18/24/26 behaviors still green
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P11-S02-02**

## Out of scope
- S03+ Phase 11 product work (DF-47, DF-51, …)
- DF-45 review get/list/show (S07)
- Real operator identity / authn
- Rewriting Mode-B Gate C packs / Phase 00–10 history

## Todo updates
Implementer: **status + notes only**. Record test names + DF-43/44 evidence. No spawning; no rewriting upcoming prompts.

## Minimal todos
- [ ] Domain FAIL+PASS reject test + implementation
- [ ] Honesty Path C supersession
- [ ] DF-44 help/MCP wording
- [ ] Locked verify suite + board Notes
