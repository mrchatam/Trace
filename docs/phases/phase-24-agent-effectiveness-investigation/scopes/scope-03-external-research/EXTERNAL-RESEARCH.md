# External research — agent planning / memory comparables

Phase 24 / Scope 03 / P24-S03-01 deliverable. Answers [INVESTIGATION.md](../../INVESTIGATION.md) Q-D using web docs + local `similar projects/` clones. **Investigation only** — no product Go.

**Trace constraints (all transfer rows must respect):** local-first SQLite under `.trace/`; library + CLI primary surface; incremental reindex; **no daemon/HTTP/MCP required on P0-X core path**; progressive context over full-graph dumps; agent claims ≠ evidence.

---

## §1 Executive summary

Peers split into three patterns Trace can borrow without hosted dependencies: **(1) hard plan-before-edit separation** (OpenHands PLAN.md-only planning agent; Aider architect→editor two-stage mode); **(2) progressive graph context instead of seed anchoring** (CodeGraph `codegraph_explore`, Understand-Anything 1-hop context packets, Graphify local AST graphs); **(3) explicit task lifecycle + promotion** (AgentRQ `getTask`/`createTask`/`publishEvent`; OpenHands plan→build handoff file). Trace’s E01 gap is that it combines rich graph primitives with **soft harness enforcement** (Cursor hook allows edits without `TRACE_TASK_ID`) and **product paths that record discoveries without spawning tasks** (`trace add discovery` + `discovery_mentions_task` vs `loop apply` `spawned_tasks[]` only).

Top **anti-patterns for Trace P0-X:** Graphiti/Zep-style HTTP MCP + graph DB daemon as core path; AgentRQ hosted workspace MCP as required memory; full-repo graph dumps on every turn (UA initial `/understand`, graphify without query scoping); autonomous task spawn without human authority (AgentRQ `publishEvent` is powerful but hosted). Literature ([arXiv:2512.13564](https://arxiv.org/abs/2512.13564), [arXiv:2501.13956](https://arxiv.org/abs/2501.13956)) confirms memory must **evolve** (invalidation, episode boundaries, working vs long-term) — Trace’s sticky STOP without deliberation reset violates that dynamic.

---

## §2 Comparable systems table

| ID | Source | Type | approach | transfer to Trace | risk | tags | FM IDs |
|----|--------|------|----------|-------------------|------|------|--------|
| **UA** | [`similar projects/Understand-Anything/`](../../../../similar%20projects/Understand-Anything/) — [README L53–128](../../../../similar%20projects/Understand-Anything/README.md), [`context-builder.ts` L20–79](../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts) | local | Multi-agent `/understand` pipeline builds `.ua/knowledge-graph.json`; incremental re-runs on changed files only. Chat context = SearchEngine top-N + **1-hop edge expand** (default 15 nodes) — not full graph. **No task/backlog layer**; graph is exploration-only. | Phase 25: `trace context` progressive packet (search plan tree + hop-limited links); optional install skill for “graph build pass” before implement; commit exported graph for team onboarding (UA L281). | Initial full-repo scan is token-heavy; no edit gates; multi-agent writes graph JSON not SQLite — import story needed. | product, harness | FM-02, FM-07, FM-09 |
| **CG** | [`similar projects/codegraph/`](../../../../similar%20projects/codegraph/) — [README agent workflow L206–571](../../../../similar%20projects/codegraph/README.md) | local | Local `.codegraph/` SQLite index; **single MCP tool** `codegraph_explore` returns ranked symbols + call paths + blast radius in one call. Benchmark: without index agents burn ~43 tool calls; with index ~1–4 explores. Read-only — no backlog writes. | Consolidate Trace MCP read surface toward one high-signal “explore task/plan neighborhood” tool; keep CLI `trace context`/`trace why` as non-MCP path; document “query graph before edit” in install rules (mirror CG steering). | First `codegraph init` is full index (acceptable local one-time); MCP optional on harness only — must not become P0-X requirement. | product, harness | FM-02, FM-08 |
| **GF** | [`similar projects/graphify/`](../../../../similar%20projects/graphify/) — [README L30–65](../../../../similar%20projects/graphify/README.md), [`worked/example/README.md`](../../../../similar%20projects/graphify/worked/example/README.md) | local | `/graphify` skill: tree-sitter AST graph for code (deterministic, local); exports `graph.json` + `GRAPH_REPORT.md`. Worked example shows code+docs pipeline with explicit EXTRACTED vs INFERRED edges. **No planning loop** — agent queries graph after build. | Trace seed export already ships `graph.json`; borrow **edge provenance labels** and “query don’t grep” skill pattern for gap passes; keep code indexing incremental (GF code path is AST-only, no full LLM rescan). | Non-code assets need LLM pass; `app.graphify.com` hosted path exists — Trace must stay on local CLI/skill pattern only. | product, experiment | FM-02, FM-07 |
| **GT** | [`similar projects/graphiti/`](../../../../similar%20projects/graphiti/) — [README L41–48, L96](../../../../similar%20projects/graphiti/README.md), [`mcp_server/README.md` L18–28, L56–75](../../../../similar%20projects/graphiti/mcp_server/README.md), [arXiv:2501.13956](https://arxiv.org/abs/2501.13956) | local (+ hosted MCP) | **Temporal context graph:** facts have validity windows; **episodes** are provenance units; incremental ingestion without full recompute. MCP adds episodes, search, entity CRUD; default transport **HTTP `:8000/mcp/`** + FalkorDB/Neo4j. | Borrow **episode boundary** vocabulary for deliberation reset (new episode after gap pass); temporal invalidation for superseded decisions/discoveries in export — implement in SQLite, not external graph DB. | **Not transferable:** Docker HTTP MCP + graph DB daemon on core path; conflicts with Trace local-first / no P0 daemon law. | product, protocol | FM-02, FM-10 |
| **CM** | [`similar projects/codebase-memory-mcp/`](../../../../similar%20projects/codebase-memory-mcp/) — [README L17–40, L624–629](../../../../similar%20projects/codebase-memory-mcp/README.md), [arXiv:2603.27277](https://arxiv.org/abs/2603.27277) | local | Native binary; RAM-first index; **15 MCP tools** (search, trace, architecture, impact, **manage_adr** write). 100% local; `install` wires 43 agent surfaces. Read-heavy; ADR CRUD is closest “add structured artifact” but not executable task queue. | Improve MCP tool descriptions: separate **read context** vs **write plan/task** tools with ranked ordering; optional ADR-like `decision` promotion docs in install; benchmark token budget (CM claims ~120× token reduction vs grep loops). | Many tools → discovery paralysis (Trace FM-08 mirror); still MCP-harness optional; index can be large but incremental file watches. | product, harness | FM-08 |
| **AR** | [`similar projects/agentrq/`](../../../../similar%20projects/agentrq/) — [README L17–33, L77–96](../../../../similar%20projects/agentrq/README.md), [`server.go` L264–326](../../../../similar%20projects/agentrq/backend/internal/controller/mcp/server.go) | local repo / **hosted MCP** | Task workspace over **HTTP MCP** (`WORKSPACE_ID.mcp.agentrq.com`). Agent **`getTask`** dequeues next “not started” task; **`updateTaskStatus`** (`ongoing`→`completed`/`blocked`); **`createTask`**; **`publishEvent`** auto-creates subscriber trigger tasks. Instructions require status transition at start/end. | Model **`trace_tasks` pull + `trace_add task` + `trace_transition`** as explicit dequeue/state machine in MCP instructions; optional event→spawn pattern for human-approved backlog expansion (local, not HTTP). | **Hosted HTTP MCP** + OAuth — contrast only; auto-spawn via `publishEvent` needs human gate to match Trace law. | product, protocol | FM-01, FM-10 |
| **SWE** | [SWE-agent](https://github.com/SWE-agent/SWE-agent) — [hello world](https://swe-agent.com/latest/usage/hello_world/), [trajectories](https://github.com/SWE-agent/SWE-agent/blob/main/docs/usage/trajectories.md), [FAQ demos](https://swe-agent.com/latest/faq/) | web | **Issue-in → trajectory-out:** sandbox ACI (bash + editor tools); system prompt encodes **ordered steps** (reproduce → fix → verify); **`submit`** terminates when patch ready; `.traj` JSON logs thought/action/observation/state per step. Demonstration trajectories at run start. **No separate plan graph** — replanning = more LM turns in same issue context. | Harness: embed step template + mandatory **`trace loop gate --for edit`** before product edits; trajectory-style export of gate decisions; demo trajectory in `trace install` for greenfield builds. | Sandbox-bound; soft prompt discipline (edits not blocked by product gate); sticky state — community notes replay fragility ([issue #1068](https://github.com/SWE-agent/SWE-agent/issues/1068)). | harness, protocol | FM-03, FM-07 |
| **OH** | [OpenHands](https://github.com/All-Hands-AI/OpenHands) — [planning agent guide](https://docs.openhands.dev/sdk/guides/agent-custom), [planning preset](https://github.com/OpenHands/software-agent-sdk/blob/main/openhands-tools/openhands/tools/preset/planning.py), [PR #10439](https://github.com/All-Hands-AI/OpenHands/pull/10439) | web | **Hard plan-before-edit:** planning agent = Glob/Grep + **`plan_md` tool writing only `/workspace/PLAN.md`**; code agent has full edits. Plan mode = **new sub-conversation**; **Build Now** handoff via shared **`PLAN.md` file**, not merged chat history ([SDK #2335](https://github.com/OpenHands/software-agent-sdk/issues/2335)). CLI `/plan` `/execute`. | Phase 25 dual-mode: **Plan** = read-only + `trace loop apply`/`plan_changes` only; **Build** = EXECUTE after plan_critiqued; parent orchestrator writes PLAN/SPEC before delegate (matches PROMPT-G1 intent). | Requires sandbox workspace file; sub-conversation UX complexity; plan agent still advisory if user stays in Build mode. | harness, product | FM-03, FM-04 |
| **AID** | [Aider](https://github.com/Aider-AI/aider) — [repo map](https://aider.chat/docs/repomap.html), [architect mode](https://aider.chat/docs/usage/modes.html), [architect/editor split](https://aider.chat/2024/09/26/architect.html) | web | **Repo map** (PageRank-ranked symbols, `--map-tokens` budget) sent every turn. **`architect` mode:** architect model proposes solution → editor model emits diffs (**two inference steps**). **`ask` mode** never edits. Modes swappable mid-session. | Map **architect→editor** to Trace **deliberation packet→EXECUTE edits**; dynamic context budget for `trace context` (rank tasks/plan nodes by relevance); install prompt “use architect equivalent: loop next before edit”. | Mode selection is user/harness opt-in — default `code` mode can skip planning (FM-09 parallel); repo map ≠ task backlog. | harness, product | FM-03, FM-08 |
| **CUR** | [Cursor Rules](https://cursor.com/docs/context/rules), [Cursor Hooks](https://cursor.com/docs/hooks) | web | **Rules:** prompt-level, optionally always-on / glob / agent-selected — **advisory**. **Hooks:** `preToolUse` can **`permission: deny`** on Write/Edit; `beforeShellExecution`/`beforeMCPExecution` block shell/MCP; `subagentStart`/`subagentStop` for Task tool; exit code 2 = deny. Trace E01 hook allows when no `TRACE_TASK_ID` ([`trace-loop-gate.sh`](../../../../../experiments/ab-incident-tracker/runs/G1/.cursor/hooks/trace-loop-gate.sh) L5–7). | **Harness:** parent orchestrator must set `TRACE_TASK_ID`; `preToolUse` Write matcher → `trace loop gate --for edit` with **reason_code in JSON**; `failClosed: true`; Multitask docs in install. Rules alone insufficient (FM-05). | Hook API still evolving (`preToolUse` `ask` not enforced per docs); rules can be long/ignored — hooks required for choke points ([ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md)). | harness, protocol | FM-04, FM-05, FM-09 |
| **LIT** | [arXiv:2512.13564](https://arxiv.org/abs/2512.13564) *Memory in the Age of AI Agents*; [arXiv:2501.13956](https://arxiv.org/abs/2501.13956) Graphiti/Zep temporal KG; [arXiv:2602.05665](https://arxiv.org/html/2602.05665v1) graph memory survey | paper | Surveys unify **memory forms** (token/parametric/latent), **functions** (factual/experiential/working), **dynamics** (formation, evolution, retrieval). Graph memory lifecycle: extract → store → retrieve → **evolve**. Temporal KGs track validity windows — supports replanning when facts superseded. | Use survey vocabulary in S04: Trace discoveries = experiential memory; seed tasks = working memory; need **evolution** step (promote discovery→task, invalidate STOP). Cite when designing deliberation reset + export honesty. | Academic frameworks don’t prescribe CLI; avoid over-fitting to hosted Zep/Graphiti deployment patterns. | protocol, product | FM-09, FM-10 |

---

## §3 Research question answers

### §3.1 Forced replanning patterns

- **OH / AID:** Strongest **forced** replanning — OpenHands planning agent **cannot** edit source (only `PLAN.md`); Aider **`ask` mode** never edits; **`architect` mode** separates planning inference from edit inference ([OpenHands planning guide](https://docs.openhands.dev/sdk/guides/agent-custom), [Aider modes](https://aider.chat/docs/usage/modes.html)).
- **GT / LIT:** Temporal graphs **invalidate** stale facts when new episodes arrive ([Graphiti README L66–69](../../../../similar%20projects/graphiti/README.md), [arXiv:2501.13956](https://arxiv.org/abs/2501.13956)) — replanning is data model operation, not prompt nudge.
- **SWE:** Replanning is **implicit** — same GitHub issue context continues until `submit`; demonstration trajectories steer behavior ([SWE-agent FAQ](https://swe-agent.com/latest/faq/)) — no product gate equivalent to Trace STOP.
- **Trace gap:** P19 saturation + sticky STOP block forward progress without a peer-style **plan artifact handoff** or **episode reset** (S02 §3).

### §3.2 Plan-before-edit harness patterns

- **OH:** Tool allowlist enforcement — planning preset = read-only + single write target ([`planning.py`](https://github.com/OpenHands/software-agent-sdk/blob/main/openhands-tools/openhands/tools/preset/planning.py)).
- **CUR:** `preToolUse` **deny** on Write + structured `user_message`/`agent_message` ([hooks docs](https://cursor.com/docs/hooks)); Trace install uses gate CLI but **allows** parent bypass without task id (FM-04/05).
- **AID:** Repo map + architect stage ensures **context before edit**; still soft unless user picks architect mode.
- **CG / UA:** **Orient-before-edit** via graph query, not policy gate — reduces blind edits but does not block them.

### §3.3 Discovery → task promotion

- **AR:** **`publishEvent`** creates trigger tasks for subscribers; **`createTask`** explicit; agent must **`updateTaskStatus`** on start ([`server.go` L269–326](../../../../similar%20projects/agentrq/backend/internal/controller/mcp/server.go)) — clearest promotion/state machine (hosted).
- **SWE:** Findings stay in **trajectory** until human/issue closure — no auto backlog spawn ([trajectories.md](https://github.com/SWE-agent/SWE-agent/blob/main/docs/usage/trajectories.md)).
- **UA / CG / CM / GF:** **No task promotion** — graph/memory layers only; gap discovery stays in chat or graph nodes.
- **Trace gap (FM-10):** Peers either (a) spawn tasks explicitly (AR) or (b) don’t model tasks at all; Trace records **`discovery` + `discovery_mentions_task`** without **`loop apply` `spawned_tasks[]`** — worst of both (S02 §2 FM-10).

### §3.4 Orchestrator vs worker memory

- **CUR:** **`subagentStart`/`subagentStop`** hooks; Multitask parent can spawn workers — matches E01 Session A orchestrator bypass unless parent holds graph ([hooks docs](https://cursor.com/docs/hooks)).
- **OH:** Planning sub-conversation **does not inherit parent chat history**; plan transfers via **shared file** only ([SDK #2335](https://github.com/OpenHands/software-agent-sdk/issues/2335)).
- **UA:** Multi-agent **pipeline** builds one committed graph artifact; chat uses scoped retrieval — centralized graph, distributed builders ([README L128](../../../../similar%20projects/Understand-Anything/README.md)).
- **PROMPT-G1:** Parent should own Trace graph ([L81–90](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md)); workers gate on `TRACE_TASK_ID` — peer pattern exists but harness doesn’t enforce on parent.

### §3.5 Saturation / budget analogs

- **CG:** Measured **tool-call budget** — without graph, agents exhaust ~43 calls on discovery ([README L206](../../../../similar%20projects/codegraph/README.md)); with graph, 1–4 explores — analog to Trace hop/P19 saturation if “orient” never completes.
- **AID:** **`--map-tokens`** caps repo map size; expands when chat lacks file context ([repo map docs](https://aider.chat/docs/repomap.html)) — explicit token budget knob Trace lacks for plan packets.
- **SWE / OH:** No hop counter; termination via **submit** or agent finish — budget = LM steps in sandbox.
- **Trace FM-03:** `p19_saturated` on first empty `loop apply` + sticky STOP is **stricter than peers** and misfires on greenfield + post-gap verify (S02 §3).

---

## §4 Anti-patterns (Trace law violations)

| Anti-pattern | Example peer | Law violated | Trace stance |
|--------------|--------------|--------------|--------------|
| HTTP MCP / daemon required for core loop | Graphiti MCP `@8000/mcp/` + FalkorDB ([`mcp_server/README.md` L56–75](../../../../similar%20projects/graphiti/mcp_server/README.md)); AgentRQ hosted workspace URL ([README L81](../../../../similar%20projects/agentrq/README.md)) | No daemon/HTTP on P0-X; local-first | MCP optional harness only; library+CLI authoritative |
| Full-graph context dump every turn | UA initial `/understand` whole repo ([README L128–130](../../../../similar%20projects/Understand-Anything/README.md)); unscoped graphify on monorepo | Progressive context; incremental reindex | Scoped `trace context` / hop limits; seed export ≠ prompt dump |
| Full rebuild on every change | Graphiti “clear graph and rebuild indices” MCP op ([`mcp_server/README.md` L22](../../../../similar%20projects/graphiti/mcp_server/README.md)) | Incremental reindex | SQLite upsert + export snapshots |
| Autonomous swarm without task authority | AgentRQ `publishEvent` auto task creation ([`server.go` L323–326](../../../../similar%20projects/agentrq/backend/internal/controller/mcp/server.go)) | User decisions authoritative | Human-approved `trace add task` or `loop apply` spawn |
| Advisory-only enforcement | Cursor Rules without hooks ([Rules docs](https://cursor.com/docs/context/rules)) | P23 product+harness choke points | `trace loop gate` + `preToolUse` deny path |
| Agent stdout as evidence | SWE trajectories without product verifier | Agent claims ≠ evidence | Gate reads SQLite + git; export `--strict` |

---

## §5 S02 residual crosswalk

| S02 residual | Peer pattern | Trace gap | S04 intervention hint |
|--------------|--------------|-----------|------------------------|
| **Sticky STOP / reason-code UX** | CUR hooks return `permission: deny` + `user_message`/`agent_message` ([hooks](https://cursor.com/docs/hooks)); SWE `.traj` stores per-step state | Export shows `p19_saturated`; live gate shows `hop_budget_exceeded` when `Stopped=true` (S02 §3) | **product:** unify persisted vs live `reason_code`; **ux/docs:** explain recovery steps |
| **Deliberation reset after gap pass** | GT **episode** boundaries ([README L78–80](../../../../similar%20projects/graphiti/README.md)); OH **new sub-conversation** + PLAN file handoff ([SDK #2335](https://github.com/OpenHands/software-agent-sdk/issues/2335)); AR **`updateTaskStatus`** cycle | No API to clear `Stopped`/reset `hop_count` (S02 §3) | **product:** gap-pass transition or episode import; **protocol:** directed-gap rubric |
| **Discovery without task promotion** | AR explicit **`createTask`/`publishEvent`**; SWE keeps findings in trajectory only — no false task illusion | 7 discoveries, 0 new tasks; mentions-task links only (FM-10) | **product:** promote BLOCKING discovery → `spawned_tasks`; **harness:** MCP nudge after `trace_add discovery` |
| **Orchestrator vs worker graph** | OH parent/child conversations + shared PLAN; CUR `subagentStart`; UA single committed graph | Multitask parent edits without `TRACE_TASK_ID` (FM-04) | **harness:** parent `preToolUse` gate; **protocol:** PROMPT-G1 enforce parent graph |
| **Plan-before-edit enforcement** | OH read-only planning agent; Aider architect mode; CUR `preToolUse` Write deny | Hook allows without task id; rules advisory (FM-05) | **harness:** failClosed hook; **product:** EXECUTE phase + `plan_critiqued` clarity |

---

## §6 Open gaps / deferred

- **Human product call:** Auto-spawn tasks from discoveries vs always require explicit `loop apply` — AR `publishEvent` shows auto-spawn power but needs local-first equivalent with approval.
- **Hop/P19 calibration:** No peer uses identical saturation signal — need Trace-specific thresholds (possibly relax first empty apply on greenfield); defer ranking to S04.
- **Graph DB vs SQLite temporal model:** Graphiti patterns attractive but infrastructure-heavy — spike “episode table” in SQLite before any external graph engine.
- **Cursor hook stability:** `preToolUse` deny works per docs; community notes API nuances — verify against Trace hook script in Phase 25 harness spike.
- **Optional comparables:** None added — 10 locked IDs + LIT cover scope.

---

## Sources index

| Tier | IDs | Primary evidence |
|------|-----|------------------|
| A local | UA, CG, GF, GT, CM, AR | `similar projects/*` README + mechanism files cited in §2 |
| B web | SWE, OH, AID, CUR | Verified URLs in §2 |
| C literature | LIT | arXiv:2512.13564, arXiv:2501.13956, arXiv:2602.05665 |

Internal baseline (not re-audited): [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md), [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md), [ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md).
