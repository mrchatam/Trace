# P01-S03-02 — Scope review notes (Agent X0 harness)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Reviewer:** independent (fresh context; claims checked against repo + re-run tests)

## Claims vs evidence

| Claim (P01-S03-01 Notes / 01 prompt) | Evidence |
|--------------------------------------|----------|
| Package `evals/x0` with `schema.json` v1 | `evals/x0/{doc.go,schema.json,x0_test.go}`; schema `schema_version` const `1`, `experiment` const `"X0"` |
| `TestX0DryRunMetricsB0AndG1` | `evals/x0/x0_test.go` — temp-copy fixture, init/seed(abs)/index per condition |
| Temp `metrics-b0.json` + `metrics-g1.json` with `dry_run: true` | Written under `t.TempDir()`; schema-validated; asserts `dry_run == true` |
| B0: no `why`/`context`; stub reads fixture files | `stubB0Agent` reads `src/greeter.ts`; B0 path never calls why/context; `tools_used` = `["read_file"]` and assert excludes why/context |
| G1: live `why task` + `context --format json` on Task `2222…` | `mustTrace(..., "why", "task", TaskID)` + `mustTrace(..., "context", TaskID, "--format", "json")`; non-empty stdout; exit 0 |
| Conditions labeled `B0` / `G1` | Asserted on validated JSON docs |
| Absolute seed path | `filepath.Abs(.../seed/gt.json)` → `seed import` arg + metrics `seed` field |
| Not Gate C / no “G1 beats B0” | Explicit comment on test + `doc.go`; no scoring/comparison exit criteria |
| Honesty not merged into X0 | `evals/honesty` untouched; dry-run `task_family: "understanding"` only |
| `evals/p0x` untouched / 7/7 path | Package files unchanged this review; `CGO_ENABLED=1 go test ./evals/p0x/...` PASS |
| G19: no library import of `cmd/trace` | Shells out via `go build -o … ./cmd/trace`; no Go import of `cmd/trace` |
| No MCP / daemon / HTTP / embeddings | Eval harness only; no product surface creep |
| Tests green | Independent: `CGO_ENABLED=1 go test ./evals/x0/...` + `./evals/p0x/... ./evals/honesty/...` + `./...` PASS |

## Checklist (02-scope-review)

### Conditions & CLI — PASS
- B0 metrics exist; harness does not call `why`/`context` for B0
- G1 metrics exist; harness calls `why task <UUID>` and `context <UUID> --format json` (exit 0, non-empty)
- JSON `condition` correctly `B0` / `G1`

### Schema & dry-run bar — PASS
- Committed `evals/x0/schema.json` (schema_version 1); tests validate emits
- Temp metrics-b0/g1 with `dry_run: true`
- No Gate C scoring / superiority claim as exit criteria
- H5 stays in `evals/honesty`

### Corpus / hygiene — PASS
- `fixtures/x0` + abs seed; Task UUID `22222222-…`
- No committed `.trace/` or live metrics under fixtures/evals
- `evals/p0x` still green

### Laws / regression — PASS
- No MCP requirement creep (DR-AGENT)
- No daemon/HTTP/embeddings
- G19 honored
- x0 + p0x + honesty + `./...` CGO=1 PASS

### Cross-scope — PASS
- S04/S05 prompts already lock “X0 without MCP”, temp metrics paths, and VERIFY commands — no thicken needed
- Did not rewrite S01/S02 `done` history

## Findings

| Sev | Finding | Disposition |
|-----|---------|-------------|
| medium | `github.com/santhosh-tekuri/jsonschema/v5` used directly by `evals/x0` but listed `// indirect` in `go.mod` | **Fixed inline** — promoted to direct `require` via `go get` + `go mod tidy` |
| nit | `SCOPE-TODOS.md` implement/review checkboxes lagged | **Fixed inline** — boxes synced |
| nit | `workEnv.MetricsDir` set in `setupWork` but unused (test uses its own temp dir) | Accept residual |

**Blocker/high:** none  
**Spawns:** none

## Residuals (none material)

- `toolsContainTraceCLI` uses substring match (`why`/`context`); sufficient for current tool strings; tighten later if needed.
- Schema enum still allows `task_family: "honesty"` for future Phase 02 runs; dry-run does not emit it.

## Next

**P01-S04-00** (MCP adapter scope planner)
