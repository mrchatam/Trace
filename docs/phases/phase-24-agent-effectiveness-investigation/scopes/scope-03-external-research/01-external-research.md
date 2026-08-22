# P24-S03-01 — External research investigation

## Metadata
- id: P24-S03-01
- todo_ids: [P24-S03-01]
- role: implementer (investigator / researcher)
- skills: [research, documentation-and-adrs, writing-for-agents]
- mcps: [Read, Glob, Grep, Shell, WebSearch, WebFetch]
- agents: [researcher, tech-researcher]
- verification: manual (URLs + local repo paths cited)
- hooks: none

## Objective

Answer **INVESTIGATION.md Q-D** (external comparables): how do other agent-memory / planning systems force replanning, expand backlog on gap discovery, and gate edits until plan is ready?

Write **`EXTERNAL-RESEARCH.md`** synthesizing **web research** + **`similar projects/`** local scan. Update [FINDINGS.md](../../FINDINGS.md) external-comparables section (brief bullets + link — do not duplicate full tables).

**Investigation only** — no product Go commits. Do not recommend hosted MCP/daemon on Trace P0-X core path.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — local-first, no full-rebuild, agent claims ≠ evidence
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md) — Q-D, FM taxonomy, intervention categories
- [FINDINGS.md](../../FINDINGS.md) — external section owner
- S03-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S01 evidence: [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md) — §3 FM matrix, two-mode model
- S02 product mechanisms: [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) — §2 FM table, §3 residuals, §4 cross-cutting

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md): Agent mode → clarify if comparable tier or FM mapping unclear → Plan mode → execute. **Use S01/S02 as internal baseline** — external research explains what peers do differently, not re-audit Trace code.

## Locked defaults

| Item | Value |
|------|-------|
| Output file | `scopes/scope-03-external-research/EXTERNAL-RESEARCH.md` |
| Minimum comparables | **≥6** rows in §2 table (≥3 must have **verified URL**; ≥2 from **`similar projects/`**) |
| Per-comparable columns | approach \| transfer to Trace \| risk |
| Tag taxonomy | `product` \| `harness` \| `protocol` \| `experiment` \| `ux/docs` (per INVESTIGATION intervention categories) |
| FM linkage | Each comparable maps to **≥1** FM-ID (FM-01..FM-10) or investigation question A–D |
| Product Go | **Forbidden** |
| Local repo root | `/home/ali/Desktop/Trace/similar projects/` |
| Trace constraints (cite in §1 + each risk row) | Local-first; `.trace/` gitignored SQLite; **no daemon/HTTP on P0-X**; incremental reindex; library+CLI primary surface |

## Locked comparables list (mandatory coverage)

Implementer must produce a row for **each** tier below. May add **≤2** optional rows if evidence warrants (mark `optional` in Notes).

### Tier A — Local repo (`similar projects/`)

| ID | Path / project | Why locked | FM focus |
|----|----------------|------------|----------|
| **UA** | `similar projects/Understand-Anything/` | Multi-agent graph build + context packet for Cursor/Claude; closest “graph → agent context” peer | FM-02, FM-07, FM-09 |
| **CG** | `similar projects/codegraph/` | Local semantic code graph + MCP `codegraph_explore`; surgical context before edit | FM-02, FM-08 |
| **GF** | `similar projects/graphify/` | Agent-oriented codebase graph extraction (worked examples) | FM-02, FM-07 |
| **GT** | `similar projects/graphiti/` | Temporal context graph; episode ingestion; MCP server (contrast Trace local-first) | FM-02, FM-10 |
| **CM** | `similar projects/codebase-memory-mcp/` | MCP-native codebase memory layer | FM-08 |
| **AR** | `similar projects/agentrq/` | Task workspace + MCP task pull/update (hosted contrast) | FM-01, FM-10 |

**Local read hints (start paths):**

| ID | Start paths |
|----|-------------|
| UA | `README.md`, `understand-anything-plugin/src/context-builder.ts`, `understand-anything-plugin/skills/understand/` |
| CG | `README.md` (agent workflow section), MCP tool descriptions in repo |
| GF | `README.md`, `worked/*/README.md` |
| GT | `README.md`, `mcp_server/README.md`, arXiv link in README |
| CM | `README.md`, MCP server entry |
| AR | `README.md`, backend MCP tool definitions |

### Tier B — External harness / IDE (web — verify URLs)

| ID | System | Why locked | FM focus |
|----|--------|------------|----------|
| **SWE** | [SWE-agent](https://github.com/SWE-agent/SWE-agent) | Issue→plan→edit loop; agent-computer interface; trajectory discipline | FM-03, FM-07 |
| **OH** | [OpenHands](https://github.com/All-Hands-AI/OpenHands) | Software-agent sandbox; planning vs action separation | FM-03, FM-04 |
| **AID** | [Aider](https://github.com/Aider-AI/aider) | Repo map / architect mode before edits; chat-driven plan expansion | FM-03, FM-08 |
| **CUR** | Cursor docs: [Rules](https://cursor.com/docs/context/rules), [Hooks](https://cursor.com/docs/agent/hooks), Multitask/subagent patterns | Orchestrator vs worker; hook enforcement (Trace E01 harness) | FM-04, FM-05, FM-09 |

### Tier C — Literature / surveys (≥1 required; verify URL)

| ID | Theme | Why locked | FM focus |
|----|-------|------------|----------|
| **LIT** | Agent memory + plan revision surveys | Gap detection / replanning vocabulary for S04 | FM-09, FM-10 |

**Suggested seeds (verify + cite primary):**

- Graphiti paper: [arXiv:2501.13956](https://arxiv.org/abs/2501.13956) (temporal KG — may overlap GT row; use LIT for *patterns*, GT for *product*)
- Plan-and-solve / ReAct lineage — search for recent SE-agent planning surveys (2024–2026)
- “Discovery → task promotion” or “open-world task decomposition” in coding agents

Do **not** fabricate URLs. If a source is paywalled, note in EXTERNAL-RESEARCH §4 and cite abstract/author preprint only.

## S02 review residuals (must address in research)

Forwarded from P24-S02-02; each needs **≥1 comparable pattern** or explicit “no peer pattern found” in EXTERNAL-RESEARCH §3:

| Residual | Research question | Compare against |
|----------|-------------------|-----------------|
| **Sticky STOP / reason-code UX** | Do peers surface *why* edit is blocked (saturation vs budget vs plan missing)? | SWE, OH, AID, CUR hooks |
| **Deliberation reset after gap pass** | Do peers reset planning state after fix pass or require new task/issue? | SWE issue loop, AR task states, GT episode boundaries |
| **Discovery without task promotion** | Do peers auto-spawn backlog items from findings? | AR, SWE, UA pipeline stages |
| **Orchestrator vs worker graph** | Do peers centralize memory/plan on parent? | UA multi-agent, CUR Multitask, OH delegation |
| **Plan-before-edit enforcement** | Hard gates vs advisory rules? | SWE, AID architect, CUR hooks, Trace install contrast |

## Research questions → FM mapping (must answer in §3)

Map each INVESTIGATION Q-D bullet + selected A/B harness questions to comparables:

| Research question | FM IDs | Deliverable section |
|-------------------|--------|---------------------|
| Q-D: How do others force replanning? | FM-03, FM-09, FM-10 | §3.1 |
| Q-D: Harness patterns requiring plan expansion before edit? | FM-03, FM-05, FM-08 | §3.2 |
| Q-A: Why no `trace add` after discovery? (peer task promotion) | FM-01, FM-08, FM-10 | §3.3 |
| Q-A: Multitask / orchestrator bypass | FM-04, FM-09 | §3.4 |
| Q-B: Saturation / hop budget calibration (peer analogs) | FM-03 | §3.5 |
| Q-B: Discovery_mentions_task without new tasks | FM-10 | §3.3 |

## Required reads (mandatory)

Read before writing comparable rows. Record path or URL in every claim.

### Phase 24 handoff (internal)

| # | Path | Extract |
|---|------|---------|
| 1 | [INVESTIGATION.md](../../INVESTIGATION.md) | Q-D; FM taxonomy; intervention categories; two-mode table |
| 2 | [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md) | §2 two-mode; §3 FM confirmed statuses |
| 3 | [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) | §2 FM mechanisms; §3 residuals; §4 apply-vs-add split |
| 4 | [FINDINGS.md](../../FINDINGS.md) | Do not contradict S01/S02; add external bullets only |
| 5 | [project-rules.md](../../../../rules/project-rules.md) | Trace law summary for §1 constraints |
| 6 | [ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md) | What Trace install already enforces (contrast peers) |

### Local comparables (≥2 projects — skim README + one mechanism file each)

| # | Path | Extract |
|---|------|---------|
| 7 | `similar projects/Understand-Anything/README.md` | Graph build pipeline; agent integration |
| 8 | `similar projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts` | How context is assembled for agents |
| 9 | `similar projects/codegraph/README.md` | Local graph; agent workflow; MCP surface |
| 10 | `similar projects/graphify/README.md` | Graph extraction approach |
| 11 | `similar projects/graphiti/README.md` + `mcp_server/README.md` | Temporal episodes; MCP (hosted caution) |
| 12 | `similar projects/agentrq/README.md` | Task MCP; agent pull model |

### E01 harness (Trace-specific contrast)

| # | Path | Extract |
|---|------|---------|
| 13 | [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) L75–90 | Harness recommends task add; parent owns graph |
| 14 | [SUBAGENT-DELEGATION.md](../../../../../experiments/ab-incident-tracker/prompts/SUBAGENT-DELEGATION.md) | Worker packet; TRACE_TASK_ID |
| 15 | `experiments/.../runs/G1/.cursor/hooks/trace-loop-gate.sh` | Hook allow path without task id |

## Investigation tasks (by theme)

1. **Plan expansion gates (FM-03, FM-05)** — For SWE/OH/AID/CUR: what blocks file edit until plan/issue/architect step complete? Hard tool denial vs soft prompt?
2. **Task/backlog promotion (FM-01, FM-10)** — For AR/SWE/UA: when findings appear, is a new task/issue auto-created or linked-only?
3. **Graph/memory honesty (FM-02, FM-07)** — For UA/CG/GF/GT: what gets written to graph vs ephemeral chat? Export/portability?
4. **Tool surface / discoverability (FM-08)** — For CG/CM/CUR MCP: how are “add work item” tools described vs read-only context?
5. **Orchestrator pattern (FM-04, FM-09)** — UA multi-agent vs CUR Multitask vs OH: who owns plan state?
6. **Reset after fix pass (S02 residual)** — Peers: new issue, new session, state machine transition, or sticky “done”?
7. **Anti-patterns (§4)** — Hosted-only memory, full-graph dump, daemon-required core path — label explicitly **not transferable** to Trace P0.

## Search themes (web supplement)

Use WebSearch/WebFetch; prefer official docs + GitHub README over blog summaries.

- `SWE-agent planning before edit trajectory`
- `OpenHands agent planning action separation`
- `Aider architect mode repo map before edit`
- `Cursor hooks agent gate before edit`
- `agent memory knowledge graph replanning software engineering`
- `temporal knowledge graph agent task decomposition`

## Deliverable

### `EXTERNAL-RESEARCH.md` — locked template

Create in this scope folder.

#### §1 Executive summary

3–5 sentences: top transferable patterns + top anti-patterns. State Trace constraints (local-first, no P0 daemon) up front.

#### §2 Comparable systems table (required)

| ID | Source | Type | approach | transfer to Trace | risk | tags | FM IDs |
|----|--------|------|----------|-------------------|------|------|--------|
| UA | `similar projects/Understand-Anything/` | local | … | … | … | product,harness | FM-02,… |
| … | … | local\|web\|paper | … | … | … | … | … |

**Column rules:**

- **approach** — how they handle plan expansion, memory graph, edit gates (concrete mechanism, not marketing)
- **transfer to Trace** — actionable delta for Phase 25 themes (library, CLI, install, protocol) respecting Trace laws
- **risk** — adoption cost, hosted dependency, conflicts with incremental reindex / no daemon
- **tags** — one or more intervention categories
- **FM IDs** — comma-separated FM-01..FM-10

#### §3 Research question answers (required)

##### §3.1 Forced replanning patterns
##### §3.2 Plan-before-edit harness patterns
##### §3.3 Discovery → task promotion
##### §3.4 Orchestrator vs worker memory
##### §3.5 Saturation / budget analogs

Each subsection: 2–4 bullets with **comparable ID citations** (not generic claims).

#### §4 Anti-patterns (required)

Table or bullets: hosted MCP/daemon on core path, full-graph context dump, autonomous swarm without human task authority — map to Trace laws violated.

#### §5 S02 residual crosswalk (required)

| S02 residual | Peer pattern | Trace gap | S04 intervention hint |
|--------------|--------------|-----------|------------------------|
| Sticky STOP UX | … | … | product / harness |
| Deliberation reset | … | … | … |
| … | … | … | … |

#### §6 Open gaps / deferred

Items needing human product call or Phase 25 spike — do **not** rank here (S04 owns INTERVENTION-MATRIX).

### FINDINGS.md update

Add or refresh **External comparables** subsection under investigation answers:

- 3–5 bullets summarizing §2 top transfers
- Link to `EXTERNAL-RESEARCH.md`
- Mark row `draft` until S03-02 review approves

## Trace law compliance notes (implementer must respect)

When writing **transfer** and **risk** columns:

1. **No daemon/HTTP/MCP-on-P0-X** as *required* core path — MCP as optional harness OK to *describe* in peers; Trace recommendation must stay library+CLI first ([project-rules.md](../../../../rules/project-rules.md)).
2. **No full-graph dumps** — prefer progressive context (UA/CG patterns OK if scoped retrieval).
3. **Incremental reindex** — reject peer patterns that rebuild entire graph on every edit.
4. **Agent claims ≠ evidence** — cite URLs, README paths, or file:line from local clones.
5. **User decisions authoritative** — peer “auto-spawn task” patterns note where human approval is still required.

## Exit criteria

- [ ] `EXTERNAL-RESEARCH.md` exists with §1–§6 per template
- [ ] **≥6** comparables in §2; **≥3** verified external URLs; **≥2** from `similar projects/`
- [ ] Every locked Tier A/B ID (UA, CG, GF, GT, CM, AR, SWE, OH, AID, CUR) has a §2 row OR explicit “skipped” with reason in §6
- [ ] **≥1** Tier C literature row with verified URL
- [ ] Each §2 row has approach + transfer + risk + tags + FM IDs
- [ ] §3 answers all research questions in mapping table
- [ ] §4 anti-patterns explicitly flag hosted/daemon/full-dump peers
- [ ] §5 addresses all S02 forwarded residuals
- [ ] FINDINGS external section updated (brief + link)
- [ ] No product Go in diff
- [ ] Board row P24-S03-01 set `done` with evidence paths in Notes

## Minimal todos

- [ ] Read required internal handoff (POSTMORTEM, CODEBASE-AUDIT, INVESTIGATION Q-D)
- [ ] Skim locked Tier A projects (README + one mechanism file each)
- [ ] Web research Tier B harness docs (SWE, OH, AID, Cursor) — capture URLs
- [ ] Find ≥1 Tier C paper/survey with verified URL
- [ ] Draft §2 table (all locked IDs)
- [ ] Write §3–§5 mapped to FMs and S02 residuals
- [ ] Write §4 anti-patterns vs Trace laws
- [ ] Update FINDINGS external section
- [ ] Self-check exit criteria; set board row done

## Next

**P24-S03-02**
