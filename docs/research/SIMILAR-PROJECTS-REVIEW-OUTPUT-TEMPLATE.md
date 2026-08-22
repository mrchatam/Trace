# Similar-projects review — output template

Copy this file to a dated findings doc (e.g. `SIMILAR-PROJECTS-REVIEW-YYYY-MM-DD.md`) and fill every section. Do not leave empty headings. See [`SIMILAR-PROJECTS-REVIEW-PROMPT.md`](SIMILAR-PROJECTS-REVIEW-PROMPT.md) for the reviewing-agent prompt.

**Meta**

| Field | Value |
|-------|-------|
| Review date | |
| Reviewer / agent | |
| Trace commit / branch | |
| Peers path | `similar projects/` |

---

## A. Inventory

| Dir | One-line what it appears to be |
|-----|--------------------------------|
| `agentrq` | |
| `codebase-memory-mcp` | |
| `codegraph` | |
| `graphify` | |
| `graphiti` | |
| _(add rows if more clones appear)_ | |

---

## B. Technique findings

### `agentrq`

- 

### `codebase-memory-mcp`

- 

### `codegraph`

- 

### `graphify`

- 

### `graphiti`

- 

---

## C. Candidate steals (ranked)

| Rank | Technique | Source | Trace surface (H1–H7 / named) | Fit vs laws | Effort | Risk | Lane | Notes |
|-----:|-----------|--------|-------------------------------|-------------|--------|------|------|-------|
| 1 | | | | accept / adapt / reject | S/M/L | low/med/high | product / dogfood-only / reject | |
| 2 | | | | | | | | |
| 3 | | | | | | | | |

---

## D. Explicit rejects

| Idea | Why reject (law / hard boundary / wrong-product) |
|------|--------------------------------------------------|
| | |
| | |

---

## E. FUTURE PHASE (run when ready; do not assume Phase 11)

> Paste onto `docs/TODO.md` only when the human schedules it. Phase number is TBD — use `Pxx` placeholders below.

### Suggested name + slug

- **Name:**
- **Slug:** `phase-xx-…`

### Problem



### Non-goals

- 
- 

### Ordered scopes

| Scope | One-line |
|-------|----------|
| S01 | |
| S02 | |
| S03 | |
| S04 VERIFY | |

### Board-ready TODO stub

```markdown
## Phase XX — <name> (FUTURE; schedule when ready — do not assume Phase 11)

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 1 | Pxx-00 | pending | phases/phase-xx-<slug>/00-PHASE-PLANNER.md | phase planner |
| 2 | Pxx-S01-00 | pending | phases/phase-xx-<slug>/scopes/scope-01-…/00-PLANNER.md | |
| 3 | Pxx-S01-01 | pending | phases/phase-xx-<slug>/scopes/scope-01-…/01-….md | |
| 4 | Pxx-S01-02 | pending | phases/phase-xx-<slug>/scopes/scope-01-…/02-scope-review.md | |
| 5 | Pxx-S02-00 | pending | … | |
| … | … | pending | … | |
| N | Pxx-S0N-02 | pending | …/02-scope-review.md | VERIFY review + DR-HANDOFF |
```

---

## F. Open questions for the human

1. 
2. 
