# Phase 22 — Locks and supersessions

This log answers: **what Phase 22 implements, what it still forbids, and which prior deferrals it supersedes.**

Legend: **promote** = implement in Phase 22 · **out** = still excluded (hard boundary) · **keep** = intentional permanent simplification · **supersede** = prior phase lock is replaced for upcoming P22 work only (do not rewrite done P17–P21 prompts)

| ID | Topic | Prior state | Decided by | P22 disposition |
|----|-------|-------------|------------|-----------------|
| D-22-01 | 43 unchecked checklist bullets | Audit of `docs/CAPABILITIES_CHECKLIST.md` after P21 close | **Human operator** (phase promotion) + this planner | **promote** all 43 — no post-MVP / later / optional / deferred leftovers |
| D-22-02 | Git-hook DF-86 | P17 CONDITIONAL/deferred; P18/P21 VERIFY non-fail if absent | **P17-00** + DF-84-FORWARD | **promote** `trace install git-hook` (S02) — still must **not** wrap `git commit`; still not a daemon |
| D-22-03 | Autonomous / embedded test runner | P20/P21 D-16 **out** — harness records via `loop apply` | **P20 README** + P21 DECISION-LOG D-16 | **supersede** for **explicit** invoke: `trace test run` may exec the project’s test command and **must** record `outcome_results`. **Keep out:** background daemon watcher, hosted CI, silent auto-run without agent/hook trigger |
| D-22-04 | `BuildPolicyInputs` cycle flags | SelectNext 14-row table exists; live `BuildPolicyInputs` never sets ExecutePending/TestPending/EvaluationPending/ReflectPending | **P21-S03-01** stub + P21-S08-02 residual | **promote** (S03-01) — 14-phase cycle must be store-driven |
| D-22-05 | MCP catalog frozen at 10 tools | P21 VERIFY required unchanged catalog | **P21 S05/S07/S08 locks** | **supersede** — stdio MCP may add tools that mirror new CLI (search, loop, history, **trace_agents**). Hosted MCP still **out**. Update `TestToolNamesRegistered` when tools land |
| D-22-06 | Code graph in seed JSON | P17 omits index; clones `index` after import | **P17 seed deny list** | **keep** omit — rebuild via `trace index`. C04 closes via hook+honesty+index, not by exporting symbols |
| D-22-07 | Change capture opt-in `loop apply` only | P20 CreateChange via apply writes | **P20-S03/S06** | **promote** VCS-promoted capture (S02-05) so “every meaningful change” is not leftover |
| D-22-08 | Regression attribution | Create always `correlated`; `caused` fail-closed without evidence | **P20-S05** Law 5 | **keep** fail-closed; **promote** evidence-backed `caused` when confirmed hypothesis + evidence exist (already API; S04-03 must actually associate the **change** and make it queryable) |
| D-22-09 | Predicted vs actual impact | Impact walk exists; no stored prediction snapshot compared after | **P14/P16 impact** | **promote** (S04-01) |
| D-22-10 | State comparison API | VCS commits indexed; no compare-two-states library/CLI | **P0 index** | **promote** (S02-07) |
| D-22-11 | Pattern learning / similar-change | Not built | — | **promote** deterministic aggregation (S06) — **out** ML |
| D-22-12 | Project-specific eval rules | Scoring criteria exist; no project rule file / mechanism interface | **P20-S04** outcomes | **promote** (S07) additive `Evaluator` contract + committed `trace/eval-rules.json` |
| D-22-13 | Multi-agent conflict/redundant work | Shared state exists; no overlap detector | **P11 coordination** | **promote** advisory overlap (S08-03) — does not take locks across agents |
| D-22-14 | Hosted MCP / daemon / HTTP | Hard project boundary | **Human** + `docs/TODO.md` Later developments | **out** |
| D-22-15 | Full-rebuild-on-any-change indexer | Forbidden architecture | **G_PROJECT_LAWS 12** | **out** — hook/index stay file-local |
| D-22-16 | Wrapping `git commit` | DF-86 lock | **DF-84-FORWARD** | **out** — hook is post-commit and/or pre-push only |
| D-22-17 | Requirement as own table | Merged into Goal | **P20 D-17** | **keep** merged |
| D-22-18 | Experiments bake-off engine | Thin table in P21 | **P21 D-01** | **keep** thin — S07 evaluators are not a multi-agent bake-off runner |
| D-22-19 | Seed export of new P22 tables | P21 added P20 cognition keys | **P21 D-05** | **promote** additive seed keys for knowledge + eval-rules pointer + improvements when those tables exist (S06/S07). Code graph still omitted (D-22-06) |
| D-22-20 | DONE/Review PASS policy | Operator + Review PASS | **P10/P11** | **keep** — verification cycle is deliberation/status, not a silent rewrite of `TransitionTask` DONE |
| D-22-21 | Compat ceiling | P21 = **21** (`021_experiments`); forbid 022+ | **P21-S07** | **promote** sequential migs 022–026 (see WORK-MAP). Each owning implement row bumps `evals/compat` ceiling |
| D-22-22 | Phase 21 DR-HANDOFF `no successor` | P21 complete | **P21-S08-02** | **superseded** by human promotion of Phase 22 — do not rewrite P21 close notes |
| D-22-23 | Scope count | Human sketch ~5–8; S09 added post-scaffold | **This planner** + **Human** | **9 scopes** (S01–S09); multiple implement pairs per scope; VERIFY after S09 |
| D-22-24 | Review residuals | Prior phases listed non-blocking leftovers at VERIFY | **Human intent P22** | **supersede** — reviewers **must** spawn `Na`/`Nb` for unmet capabilities; VERIFY **must not** close with leftover checklist `[ ]` unless an in-phase spawn is still runnable |
| D-22-25 | Harness agent / subagent delegation | Not in original P22 scaffold | **Human operator** (post-scaffold) | **promote** (S09) — Trace **recommends** harness agents and fresh subagents for independent review; Trace **never** spawns or runs agents |
| D-22-26 | Agent catalog kind | Capabilities had SKILL/RULE/MCP/TOOL/HOOK only | **S09 planner** | **promote** add **`AGENT`** kind + `harness_agents` table (mig **027**); map agents → required capability slugs |
| D-22-27 | Bundled default agents | No install-time agent profiles | **Human intent** | **promote** committed `trace/agents/default.json` + `trace install agents`; idempotent upsert into `.trace/` |
| D-22-28 | Phase/task routing | SelectNext chooses phase only | **S09** | **promote** deterministic `RecommendAgents` — e.g. CRITIQUE→code-reviewer, perf keywords→performance-reviewer; loop next `harness_recommendations[]` |
| D-22-29 | Subagent availability gate | agent-loop-protocol fresh-subagent is docs-only | **S09** | **promote** when `harness:subagent` capability is AVAILABLE, recommend `use_subagent: true` + prompt stub; when UNAVAILABLE, inline review with missing-cap honesty |
| D-22-30 | Future trace host agent registry | Later developments hosted Trace | **Human intent** | **promote** schema extension fields (`registry_source`, `external_url`) + `trace/agents/README.md`; **out** network fetch/sync in P22 |

## Approaches considered (brainstorm — required before scaffold)

| Approach | Trade-off | Choice |
|----------|-----------|--------|
| 17 one-section scopes | Clean mapping; planner explosion; file thrash on shared MCP/CLI | **Reject** |
| 5 mega-scopes | Fewer boards; implementer sessions too large | **Reject** |
| 8 scopes × several implement/review pairs | Matches human sketch; session-sized rows; shared packages sequenced | **Accept** |
| Filesystem watcher daemon for C04/C25 | True continuous sync | **Reject** — daemon/HTTP out; git-hook + honesty + `trace index` |
| Export full code graph in `trace/graph.json` | Clone has symbols without index | **Reject** — Law 1 / deny-list; clone runs `index`; not a leftover |
| Change ML model for patterns | Richer “tends to help/hurt” | **Reject** — Law 13; deterministic counts + structured knowledge rows |

## Authority stack (when decisions conflict)

1. **Human operator** — promotes phases, overrides board
2. **`docs/init/G_PROJECT_LAWS.md`** — non-negotiable architecture law
3. **This DECISION-LOG** for Phase 22 locks (including supersessions of P17–P21 deferrals)
4. **Phase 22 README coverage matrix** — capability ownership
5. **Scope `00-PLANNER`** — FINAL locks for implementers
6. **Reviewers** — spawn Na/Nb; do not rewrite done history; do not close with residuals

## Schema / compat (locked sequence)

| Mig | Owner | Compat after | Forbids |
|-----|-------|--------------|---------|
| `022_code_relationships.sql` | S01-01 | **22** | 023+ until S02 |
| `023_graph_sync.sql` | S02-03 | **23** | 024+ until S04 (S02-05/07 reuse 023 or existing `changes`/`vcs_*`) |
| *(S03 — no SQL unless planner finds a hole; reuse `outcome_results`)* | S03 | stays **23** | — |
| `024_impact_compare.sql` | S04-01 | **24** | 025+ until S06 (S04-05 improvements may use 024) |
| *(S05 — query only)* | S05 | stays **24** | — |
| `025_engineering_knowledge.sql` | S06-03 | **25** | 026+ until S07 |
| `026_eval_rules.sql` | S07-01 | **26** | 027+ |
| *(S08 — no SQL; conflict is computed)* | S08 | stays **26** until S09 | — |
| `027_harness_agents.sql` | S09-01 | **27** | 028+ until post-P22 |

If a scope planner proves a table is unnecessary, they may drop **that scope’s** mig and keep the previous ceiling — they must not steal a later number. They **must** still close owned capabilities.

## Spawn / residual policy (locked)

- Implementers: board **status + notes only**.
- Reviewers: if any owned checklist bullet is unmet, spawn `P22-Sxx-0Na` implement + `0Nb` review **immediately below**. Never mark review `done` with “later” leftovers.
- VERIFY: fail the gate if any `docs/CAPABILITIES_CHECKLIST.md` item is still `[ ]` unless a spawned **in-phase** remediation row exists and is still runnable.
- DR-HANDOFF successor **`no successor`** only if the checklist is fully `[x]` (or only hard-boundary **out** items remain, which are already `[x]` or non-goals). Do **not** hand leftover capabilities to Phase 23.
