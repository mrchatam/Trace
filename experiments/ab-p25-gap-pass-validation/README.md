# E02 — P25 gap-pass validation (dogfood)

**Objective:** Test whether **Phase 25 (P25-C)** changes default agent behavior — without a human “run gap analysis” prompt.

## Hypothesis

After `trace install cursor --write`, G1 agents in a **build-only** session will:

1. See **mandatory gap pass** text in installed rules / AGENTS.md (INT-03)
2. Record **≥1 discovery** or **≥1 decision** after initial build (partial Mode B collapse)
3. Parent Multitask orchestrator sets **`TRACE_TASK_ID`** before subagent edits (INT-04)

**Null hypothesis (E01-like):** Build-only G1 still produces seed-only tasks, 0 discoveries until human directs gap work.

## Domain

**Equipment checkout desk** — different from E01 incident tracker to reduce pattern reuse.

- Check out / return assets (laptops, monitors, keys)
- Roles: requester, desk staff, admin
- REST API + simple admin UI + public “availability” page
- Audit log of checkouts

Vague stakeholder ask only; agent picks stack (expect Go + SQLite like prior arms).

## Arms

| Arm | Trace | Install | Prompt |
|-----|-------|---------|--------|
| B0 | No | — | Build only |
| G1 | Yes + strict + hook | Phase 25 `trace install cursor --write` | **Build only** — must NOT say “gap analysis” |

Critical: G1 prompt deliberately **does not** mention gap pass — we test whether **install bundle** (GapPassPrompt) is enough.

## Seed prefix

`e0200000-0000-4000-8000-…`

## Phase 26 decision tree (after E02)

| E02 outcome | Next step |
|-------------|-----------|
| G1 records discoveries/decisions **without** human gap prompt | **P25-C validated** — defer Phase 26 or run E02 repeat for confidence |
| G1 product OK but graph still thin (E01 Session A replay) | **Phase 26 → P25-A** (discovery→task promotion) and/or **P25-B** (hop reset) |
| G1 still ≡ B0 code + thin graph | **Phase 26 → full stack** (P25-A + P25-B + protocol INT-08/10) |
| Orchestrator still bypasses Trace | **Phase 26 → harness hardening** (INT-04 enforcement beyond docs) |

## Quick start

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation
./prepare.sh          # first time
./prepare.sh G1       # after B0 — does not wipe B0
```

Open **`runs/B0`** or **`runs/G1`** in Cursor; paste matching prompt from `prompts/`.

## Docs

[PROTOCOL.md](PROTOCOL.md) · [RUBRIC.md](RUBRIC.md) · [HYPOTHESIS.md](HYPOTHESIS.md)

## Score

```bash
./score.sh B0 --test
./score.sh G1 --test --p25
```
