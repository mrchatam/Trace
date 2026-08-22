# E01 — ab-incident-tracker

**Objective:** Clean A/B after experiments wipe — test Trace + Phase 23 enforcement on a **new domain** (on-call incident ops), with Multitask-safe protocol.

| Arm | Trace | Enforcement |
|-----|-------|-------------|
| B0 | No | — |
| G1 | Yes | strict config + cursor rules + cursor-hook |

## Why incident tracker (not CMS)

Prior CMS experiments caused **path reuse** and **pattern-matching** (agents copied old CMS trees). Incident ops uses **severity, assignment, status lifecycles, responder timelines, public status page** — same planning depth, zero CMS overlap.

## Feature ask (vague)

- Incident CRUD + severity + status lifecycle
- Roles: reporter, responder, admin
- Assignment + comment/timeline
- Search/filter
- REST API + auth
- Responder/admin dashboard
- Public read-only status page
- Activity/audit log

## Seed (prefix `e0100000-…`)

| UUID suffix | Entity |
|-------------|--------|
| `...0001` | Goal — Ship incident tracker |
| `...0010` | Task — Design architecture (**IN_PROGRESS** at seed) |
| `...0020` | Task — REST API and auth |
| `...0030` | Task — Responder/admin dashboard |
| `...0040` | Task — Public status page |
| `...0050` | Task — End-to-end verification |
| `...0a1–0a3` | Assumptions (stack, public read-only, ≥3 severities) |

Module: `incidentops` (empty `cmd/incidentd/main.go`).

## Operator one-liner

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-incident-tracker
./prepare.sh          # first time (both)
./prepare.sh G1       # after B0 — does not wipe B0
```

G1 session env:

```bash
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0100000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1
```

## Workspaces (open in Cursor **before** paste)

| Arm | Path |
|-----|------|
| B0 | `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/B0` |
| G1 | `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1` |

## Prompts

| File | Use |
|------|-----|
| [prompts/PROMPT-B0.md](prompts/PROMPT-B0.md) | B0 first message |
| [prompts/PROMPT-G1-ENFORCE.md](prompts/PROMPT-G1-ENFORCE.md) | G1 first message |
| [prompts/SUBAGENT-DELEGATION.md](prompts/SUBAGENT-DELEGATION.md) | G1 worker packet (any Trace task UUID, not a 4-agent cap) |

## Docs

[PROTOCOL.md](PROTOCOL.md) · [MULTITASK.md](MULTITASK.md) · [RUBRIC.md](RUBRIC.md) · [ENFORCEMENT.md](ENFORCEMENT.md)

## Score

```bash
./score.sh B0 --test
./score.sh G1 --test --gate
./run-enforcement-demo.sh   # mechanics only
```
