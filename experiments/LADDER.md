# Trace dogfood ladder

Open-ended program. Add a rung when a Trace capability is untested in a **live agent** setting.

## Tracks

| Track | Where | Proves |
|-------|--------|--------|
| Harness | `evals/*` | Deterministic CLI/library behavior |
| Dogfood | `experiments/ab-*` | Cursor agent with vs without Trace |

Status: `done` | `ready` | `next` | `planned` | `deferred`

---

## Tier 0 — Complex real-app + enforcement (current)

| ID | Experiment | Isolates | Status |
|----|------------|----------|--------|
| **E01** | [ab-incident-tracker](ab-incident-tracker/) | Vague **on-call incident tracker**; B0 vs G1 with Phase 23 **strict** + cursor-hook; revealed **two-mode** effectiveness (build vs directed gap) | **done** |
| **E02** | [ab-p25-gap-pass-validation](ab-p25-gap-pass-validation/) | **Equipment checkout desk**; build + directed arms; Phase 28 Session-B + VERIFY | **done** |
| **E03** | [ab-library-hold-desk](ab-library-hold-desk/) | **Library hold / waitlist desk**; full **P25–P28** stack; build then optional Session-B | **done** — P25-3a PASS; Session-B P25-3b PASS + live promotion (2026-08-21) |

**Hypotheses stressed:** H2 progressive planning, H5 honesty/review, P23 enforcement choke points.

**Primary failure modes to watch:**

- B0: codes without structured plan artifacts; thin VERIFY
- G1 without orchestrator Trace: post-hoc graph, G-check failures (decisions, TRACE-EVIDENCE, uncertainties)
- Multitask: wrong workspace / sibling experiment path reuse

---

## Operating rules

1. **One active `ab-*` at a time** unless docs-only.
2. Finish score → update [RESULTS.md](RESULTS.md) → pick next ID from gaps.
3. Prefer [G1-NATURAL.md](G1-NATURAL.md) for G1 product prompts.
4. After each rung: G1 pass / B0 fail → Trace helped; both pass → tighten; G1 fail / B0 pass → investigate confusion.
