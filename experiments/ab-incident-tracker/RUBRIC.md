# E01 Rubric — incident tracker

Binary checks from `./score.sh`. Human-gradable without LLM judge.

## Product (both arms)

| Check | Pass |
|-------|------|
| Tests | `go test ./...` passes (`--test`) |
| API | REST CRUD + auth + RBAC |
| Workflow | Draft/open → acknowledge → resolve (or equivalent status machine) |
| UI | Admin/responder dashboard + public status page |
| Audit | Activity log visible to admin |

Document manual checks in agent-produced **VERIFY.md**.

## B0 harness

| ID | Pass |
|----|------|
| B0-1 | No `.trace/` |
| B0-2 | `cmd/incidentd` or built binary exists |
| B0-3 | `go test ./...` (with `--test`) |

## G1 harness

| ID | Pass |
|----|------|
| G2 | `trace/graph.json` exported |
| G3 | Graph: ≥1 goal, ≥3 tasks, ≥3 decisions |
| G-seed | All seed task IDs present in export |
| G4 | `TRACE-EVIDENCE*.md` mentions deliberation |
| G5 | ≥5 tasks in Trace store |
| G6 | ≥1 resolved uncertainty in store |
| G7 | Optional: wave-0 evidence before later waves |
| E1–E3 | Cursor rules, hook, `enforce: strict` |
| E4 | With `--gate`: verify task `…0050` gate allows edit (post-run only) |

## Verdict interpretation

| Outcome | Meaning |
|---------|---------|
| G1 harness PASS + B0 PASS | Compare product depth in VERIFY + manual review |
| G1 product PASS, harness FAIL | Trace did not improve planning discipline |
| G1 fail / B0 pass | Investigate Trace confusion or over-blocking |
