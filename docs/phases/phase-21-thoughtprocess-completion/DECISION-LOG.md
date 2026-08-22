# Phase 21 — Who decided what was not implemented (P20)

This log answers: **who deferred each gap, residual, or MVP cut**, and **what Phase 21 does about it**.

Legend: **promote** = implement in Phase 21 · **out** = still excluded (hard boundary) · **keep** = intentional permanent simplification

| ID | Topic | Not in P20 because | Decided by | Evidence | P21 disposition |
|----|-------|-------------------|------------|----------|-----------------|
| D-01 | §16 Experiment objects | MVP cut — bake-offs add cost/complexity before foundation proven | **P20-00 phase planner** + [`COVERAGE.md`](../phase-20-cognitive-deliberation/COVERAGE.md) §Future | `00-PHASE-PLANNER.md` forbidden §16/§18 product; COVERAGE row §16 Future | **promote** thin `experiments` table or experiment-linked outcome rows (S07) — not full multi-agent bake-off engine |
| D-02 | §18 Risk-adaptive verification | MVP cut — needs history + policy engine | **P20-00** + **README §30 Challenge** | COVERAGE §18 Future; README lists risk-adaptive as out-of-scope | **promote** minimal rules from change paths/metadata (S07) — not ML matrix |
| D-03 | EXECUTE/TEST/EVALUATE/REFLECT/REPLAN in SelectNext | MVP table intentionally 8 rows; EXECUTE never auto-selected | **P20-S01-00 scope planner** | `scope-01/00-PLANNER.md` SelectNext FINAL table; note "EXECUTE not selected by MVP table" | **promote** full cycle routing (S03) |
| D-04 | EXPLORE as controller phase | Merged into DecisionAlternative; enum only | **P20-00** + **COVERAGE** entity merge | Option → Reuse DecisionAlternative | **keep** merged unless S03 adds optional EXPLORE when alternatives exist without decision |
| D-05 | Seed export omits P20 tables | P17 export scope; cognition deferred to forward work | **P20-S07-00 verify planner** | `01-verify.md` omit policy; DR-HANDOFF residual | **promote** (S01) |
| D-06 | Retrieval lacks `uncertainty` etc. | S02/S06 shipped loop sections first; retrieval lag | **P20-S06-01 implementer** + verify notes | `deliberation_packet.go` comment; VERIFY-NOTES residual | **promote** (S02) |
| D-07 | FTS skip on P20 upserts | Fail-closed until entity types registered | **P20-S02-01 implementer** | `cognitive.go` "Does not SyncEntityFTS" | **promote** (S02) |
| D-08 | Non-transactional `loop apply` | Inherited P19 pattern; scope not expanded | **P19-S02** + P20 reviews (low residual) | S01–S06 review notes | **promote** (S06) |
| D-09 | Baseline promotion B100→B101 | S04 thin baseline only | **P20-S04-00 planner** | Should not full lifecycle in MVP | **promote** (S04) |
| D-10 | Eval regression blocks promotion | Comparison JSON exists; promotion policy not wired | **P20-S04** gate scope + §31 gap audit | `overall_regression` not in TransitionTask | **promote** (S04) |
| D-11 | `trace why` for P20 nouns | Retrieval pre-P20 entity set | **P20 scope chain** (implicit) | `retrieval/exact.go` no uncertainty | **promote** (S05) |
| D-12 | Historical relationships in loop packet | S05 library only; S06 packet section not added | **P20-S05-00 planner** | §17 Should — links exist, no loop surfacing | **promote** (S05) |
| D-13 | `goal_id` validation on transition | Low residual; time box | **P20-S01-02 reviewer** | review notes | **promote** (S06) |
| D-14 | §31 live multi-hop eval | Fixture-scale sufficient for phase close | **P20-S07-00/01** | VERIFY-NOTES fixture residual | **promote** (S08 verify) |
| D-15 | S06 tests in `internal/loop` | Tests live in `cmd/trace/loop_test.go` by convention | **P20-S07-01 verifier** | vacuous `./internal/loop/...` pass | **promote** doc fix or move tests (S06) |
| D-16 | Trace runs tests (§31 narrative) | Harness-agnostic — agent runs, Trace records | **P20 README** + **COVERAGE** §29F Future execution | "result kinds, not three runners" | **out** — keep; optional test-intent metadata only |
| D-17 | Requirement as own table | Entity merge — Goal body + links | **P20-00** + **COVERAGE** §29B | Requirement → Merge into Goal | **keep** merged |
| D-18 | `plan_confidence` / `requirement_coverage` floats | MVP PolicyInputs boolean/thin | **P20-S01-00** | PolicyInputs struct | **promote** optional derived scores in S03 if cheap |
| D-19 | Hosted MCP / daemon / HTTP | Hard project boundary | **Human** + **`docs/TODO.md` Later developments** + **P20-00** forbidden list | G_PROJECT_LAWS; TODO Later developments | **out** |
| D-20 | Autonomous implementer loop | Out of scope MVP | **P20 README** completion bar | "No autonomous code-edit loop" | **out** |
| D-21 | ML / learned phase policies | §26 anti-complexity | **P20 README §30** + **TRACE_THOUGHTPROCESS §26** | deterministic table locked | **out** |
| D-22 | Graph DB / embeddings | Law 13 + P20 forbidden | **G_PROJECT_LAWS** + **P20-00** | forbidden list | **out** |
| D-23 | Phase 20 board itself | Human forward promotion after P19 | **Human operator** | AGENTS.md "human-promoted"; P19 DR-HANDOFF no successor | N/A — P20 **done**; P21 is **new human promotion** (this board) |
| D-24 | P20 DR-HANDOFF `no successor` | Phase 20 complete; no auto-spawn next phase | **P20-S07-02 reviewer** | DR-HANDOFF.md CLOSED | **superseded** by human promotion of Phase 21 |

## Authority stack (when decisions conflict)

1. **Human operator** — promotes phases, approves hosted product, overrides board
2. **`docs/init/G_PROJECT_LAWS.md`** — non-negotiable architecture law
3. **Phase N-00 planner** — locks MVP vs Future for that phase
4. **`COVERAGE.md` / WORK-MAP.md`** — section disposition within a phase
5. **Scope 00-PLANNER** — FINAL locks for implementers
6. **Reviewers** — may inline-fix; spawn Na/Nb rows; do not rewrite done history

## TRACE_THOUGHTPROCESS.md itself

The prompt document **§26 Avoid Premature Complexity** explicitly asks planners to separate MVP / Phase 2 / Future. P20 planners executed that instruction — they did **not** ignore the doc; they **scoped** it. Phase 21 is the promoted **Phase 2** tranche for cognition completion.
