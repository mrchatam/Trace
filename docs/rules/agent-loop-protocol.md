# Trace agent loop protocol

Every prompt under `docs/phases/**` obeys this protocol. Shared rules live here; prompts reference this file instead of duplicating.

Trace inherits the Perfect Planner *shape* from AI DJ with **Trace-specific policy** below.

---

## Session start (mandatory — all roles)

Agents begin in **Agent mode** (full tools), not Plan mode.

```text
1. Start in Agent mode
2. Read the linked prompt + board row + referenced laws/init docs
3. If anything material is unclear:
     - Ask the user (prefer numbered / multiple-choice)
     - Use grilling / brainstorming / other skills when the task warrants
     - Do not guess past a decision that affects architecture or exit criteria
4. When questions are resolved OR none were needed:
     - Switch to Plan mode
     - Produce a short plan against live repo state
5. After plan is accepted (or unattended run with prior prompt approval):
     - Return to Agent mode and execute the role loop
```

Planners, implementers, and reviewers all follow this gate. Skipping clarification when blocked on ambiguity is a defect.

---

## Prompt taxonomy

| File | Role | Product code? | Board edits |
|------|------|---------------|-------------|
| `00-PHASE-PLANNER.md` | Phase scaffold: refine scopes, order, gaps; light locks | **No** | Status/notes; may thicken *upcoming* scope folders |
| `00-PLANNER.md` | Scope planner: finalize sibling prompts + locked defaults | **No** | Status/notes; thicken *this scope’s upcoming* prompts |
| `01-*.md` / `0N-implement*.md` | Implement one bounded deliverable | **Yes** | **Status + notes only** |
| `02-*-review.md` | Independent review; small fixes or spawn | Small fixes OK | Status/notes; **spawn / thicken upcoming** |
| `03+` / `Na`/`Nb` | Spawned implement/review pairs | Per role | Per role |
| `VERIFY.md` / phase review | Gate: aggregate evidence (e.g. P0-X 7/7) | No features | Status/notes; may spawn remediation |
| Human-gated rows | Human supplies evidence | N/A | Human / orchestrator marks done |

---

## Board authority

Source of order: [`docs/TODO.md`](../TODO.md) (index). Row tables live in [`docs/TODO/phase-NN.md`](../TODO/).

1. Run **top to bottom** on the **active phase board**. Never start a higher-order pending row while a lower-order row is `pending` / `in_progress` / `failed`.
2. After every subagent finishes, **re-read** the active phase board (`docs/TODO/phase-NN.md`) before launching the next (reviewers may have inserted `Na`/`Nb` rows).
3. One **fresh subagent per row**. Review rows must not share the implementer’s session.
4. Do not start Phase N+1 until Phase N is complete (all rows `done` / `skipped` with reason) and phase gate passed.

### Status vocabulary

| Status | Meaning |
|--------|---------|
| `pending` | Not started |
| `in_progress` | Running (orchestrator sets) |
| `done` | Exit criteria met with evidence in Notes |
| `failed` | Attempted; reason in Notes |
| `blocked` | External / human / dependency |
| `skipped` | Explicitly out of scope; reason required |

---

## Forward-only immutability (best-effort hard rule)

Treat completed work as **history**:

- A row marked `done` and its prompt file are **immutable**. Do not rewrite the prompt body or pretend the work never happened.
- Fixes go **forward**: insert `Pxx-Syy-02a` (implement) / `02b` (review) **immediately below** the review that discovered the gap (or next free suffix). Prefer `a/b/c` suffixes over renumbering the entire board.
- Never edit an earlier phase’s delivered code “as if Phase N+1 can rewrite Phase N history.” If behavior must change, schedule new work in the **current or future** phase that explicitly supersedes prior behavior (document the supersession).

Absolute enforcement is impossible with LLMs — reviewers and orchestrators must **check** for backward edits and reject them.

---

## Implementer vs reviewer board rights

| Action | Implementer (`01`, `03a`, …) | Reviewer (`02`, `03b`, …) | Planner |
|--------|------------------------------|---------------------------|---------|
| Set own row `in_progress` / `done` / `failed` / `blocked` | Yes | Yes | Yes |
| Add Notes on own row | Yes | Yes | Yes |
| Rewrite future prompt substance / locked defaults | **No** | **Yes** (upcoming only) | **Yes** (upcoming only) |
| Insert new board rows + prompt files | **No** | **Yes** (spawn) | Yes (during plan, before implement wave) |
| Edit a `done` prompt or erase history | **No** | **No** | **No** |

Implementers who discover a structural gap: record it in **Notes** and continue if possible; otherwise `blocked` / `failed` with reason. The **next review row** (or phase planner if before reviews exist) owns real backlog changes.

---

## Planner duties

### Phase planner (`00-PHASE-PLANNER`)

- Inspect `docs/init/*`, design docs, current repo, prior phase outcomes.
- Improve phase README, scope list, run order, and **light** locks.
- Ensure each scope folder has at least stub `00-PLANNER` / `01` / `02` prompts.
- Consider cross-scope blast radius; only change **upcoming** artifacts.
- Do **not** deep-plan every task — scope planners finalize.

### Scope planner (`00-PLANNER`)

- Replan against **live** repo.
- Lock defaults so implementers do not re-debate stack/paths.
- Thicken `01-*` enough to run alone (objective, files, exit criteria, skills/MCPs).
- Note effects on later scopes; update **upcoming** prompts only.
- No product code.

---

## Implementer loop

After session-start clarification + Plan mode:

```text
LOOP:
  1. Implement next incomplete item in this prompt’s Minimal todos
  2. Self-check against Exit criteria (not a substitute for independent review)
  3. Fix gaps
  UNTIL exit criteria met OR failed/blocked with reason
→ Update own board row status + Notes only
```

If module boundaries change: **Note** the blast radius on the board; do not silently rewrite later prompts (reviewer/planner does that).

---

## Reviewer loop (quality-first)

Fresh subagent. After session-start + Plan mode:

```text
1. Compare claimed deliverables (prompt + Notes) to repo evidence
2. Findings by severity: blocker | high | medium | low | nit
3. blocker/high: small inline fix OR spawn implement+review pair (full prompts)
4. medium: prefer spawn unless trivial
5. Re-verify
6. UNTIL no open blocker/high without a pending follow-up
   AND Confidence medium or high with evidence
```

Reviewers **must** know how to write prompts: spawned files use this protocol (metadata, objective, locked defaults, exit criteria, skills/MCPs, session-start gate).

Insert spawns on the board **immediately below** the review row. Typical pattern:

```text
… 05 implement done
… 06 review  → spawns 06a implement, 06b review
… 06a …
… 06b …  → may spawn 06c / 06d …
```

Loop until a review exits with **high** confidence (or **medium** with explicit residual risks listed — never silent).

---

## Phase handoff (mandatory)

When a phase is known to have a successor (roadmap / init plan / phase README says so), **the closing phase owns the next-phase scaffold** — not a later ad-hoc session.

Before marking the phase VERIFY / final phase-review row `done`:

1. Ensure the **next** phase folder exists with at least:
   - `README.md` (goal, scope list, in/out)
   - `00-PHASE-PLANNER.md` (runnable)
   - Per-scope stubs: `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
2. Register next-phase board rows in `docs/TODO/phase-NN.md` and add an index link in `docs/TODO.md`, with the phase planner as the **first pending** row after this phase’s last `done` row.
3. Do **not** leave “stub README only / blocked until someone notices.”
4. If there is intentionally **no** next phase, say so in VERIFY Notes (`no successor`).

Deep tasking still belongs to the **next** phase’s `00-PHASE-PLANNER` and scope planners when that phase runs. The prior phase only delivers a **runnable handoff**, not finished implement prompts.

Forward-only: do not rewrite the prior phase’s `done` history; create/update **upcoming** next-phase files and board rows from the VERIFY/review row (reviewer/planner rights).

- Prepare instructions + artifacts
- Mark row `blocked` awaiting human evidence file under `docs/verification/` when required
- **Must not** mark human-gated criteria `done` from model self-claim

---

## Required prompt skeleton

```markdown
# <Phase> / <Scope> / <NN-title>

## Metadata
- id: PXX-SYY-NN
- todo_ids: […]
- role: planner | implementer | reviewer | verify | human
- skills: […]
- mcps: […]
- agents: […]          # optional subagent types
- verification: automated | human | mixed
- hooks: […]           # if any

## Objective
…

## References
- docs/rules/agent-loop-protocol.md (this file)
- docs/rules/project-rules.md
- docs/rules/skills-map.md
- docs/init/* as relevant
- …

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (if any)
| Item | Value |

## Preflight / Plan
…

## Role work
…

## Todo updates
Per board-rights table (implementers: status+notes only).

## Exit criteria
- [ ] …

## Minimal todos
- [ ] …
```

---

## Orchestrator snippet (user / parent agent)

```text
Use subagents to complete Phase NN — <name>
@docs/TODO.md
@docs/TODO/phase-NN.md

- One agent per todo row (fresh context)
- Do not start a row until every higher-order row is done or skipped
- After each subagent: re-read docs/TODO/phase-NN.md (spawns may have appeared)
- Follow docs/rules/agent-loop-protocol.md and project-rules.md
```

---

## Post-bootstrap critique path (Phase 37 R11)

After `trace plan bootstrap` or MCP `trace_plan action=bootstrap`, seed plan critique via **`trace loop apply`** with a `plan_changes` envelope — not a separate MCP critique-seed tool. Canonical Block 0 pattern: `TestGreenfield_MCPPlanBootstrap_EditGatePasses` in `internal/mcp/mcp_test.go`. VERIFY cross-ref: `docs/phases/phase-37-p36-residuals/scopes/scope-03-verify/01-verify.md` blocks 0 and 3.

Overview plan-gap surface (Phase 37 R8): GUI consumes `GET /v1/loop/status` `advisories[]` and gate violations only (Law 19 — no planner logic in `web/`).
