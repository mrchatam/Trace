# E03 — Full-stack regression (library hold desk)

**Objective:** End-to-end dogfood of the **Phase 25–28** Trace stack on a **new domain** (no E02 pattern reuse).

## Choices locked

| # | Choice |
|---|--------|
| 1a | Full stack regression on a new domain |
| 2a | Library hold / waitlist desk |
| 3a | B0 + G1 **build-only**; Session-B only if build is thin |
| 4a | Experiment only — promote Phase 29 **only if E03 fails** |

## Hypothesis

After Phases 25–28 (gap pass, promotion, reset, honesty, Option A hook, FM nudges), a **build-only** G1 session on a fresh domain will:

1. Ship a working app with tests
2. Record **≥1 discovery OR ≥1 decision** without a human “gap analysis” prompt (P25-3a PASS)
3. Keep install wiring green (P25-1/2)
4. Fail-closed on thin `--strict --enforce` export when graph is dishonest

**Null:** Build-only still thin (E01 Session A / early E02 replay) → run Session-B; if still thin or install broken → consider Phase 29.

## Domain

**Library hold / waitlist desk** — patrons place holds on titles; staff fulfill; public page shows waitlist position.

- Titles + copies (available / on hold / checked out)
- Hold request + cancel + fulfill lifecycle
- Roles: patron, desk staff, admin
- REST API + desk UI + public waitlist / availability page
- Audit log of hold events

Vague stakeholder ask; agent picks stack (expect Go + SQLite).

## Arms

| Arm | Trace | Prompt |
|-----|-------|--------|
| B0 | No | Build only |
| G1 Session A | Yes + strict + hook (post-P28 install) | **Build only** — no directed gap prompt |
| G1 Session B | Same workspace | **Only if** Session A P25-3a FAIL |

## Seed prefix

`e0300000-0000-4000-8000-…`

## Phase 29 gate

| E03 outcome | Next |
|-------------|------|
| P25-1/2 PASS + P25-3a PASS (build rich) | **Stack validated** — no Phase 29 |
| P25-1/2 PASS + P25-3a FAIL, Session-B P25-3b PASS | Install OK; default gap still weak — optional Phase 29 harness |
| P25-1/2 FAIL or honesty broken | Fix product/install first — Phase 29 repair |
| Promotion never used + thin forever | Phase 29 for D1 decision or stronger nudges |

## Quick start

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace

cd experiments/ab-library-hold-desk
./prepare.sh
# Open runs/B0 → prompts/PROMPT-B0.md
# After B0:
./prepare.sh G1
# Open runs/G1 → prompts/PROMPT-G1-BUILD.md
P25_ATTEST_BUILD=Y ./score.sh G1 --test --p25 --arm build
```

Docs: [PROTOCOL.md](PROTOCOL.md) · [RUBRIC.md](RUBRIC.md) · [HYPOTHESIS.md](HYPOTHESIS.md)
