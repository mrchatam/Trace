# P11-S05-02 — REVIEW-NOTES (Retrieval why / depth / trust / DPC attribution)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | **DF-49** — Exact/Why `symbol` by id; miss OK; no mig | **Pass** — `GetSymbolByID` JOIN path; `lookupEntity` `case "symbol"`; `TestWhySymbolExact` + `TestGetSymbolByID`; schema still `001`–`010` only |
| 2 | **DF-35** — depth≥2 omits sibling task bodies (titles OK); Expand goal→task no body Excerpt | **Pass** — Expand `Excerpt: ""` (no `excerpt(t.Body)`); `TestExpandDepth2NoSiblingTaskBody` + `TestExpandContextDepth2NoSiblingTaskBody` (`SECRET_HANDOFF` absent; sibling title present) |
| 3 | **DF-48** — decision/assumption MD Law 9 honor + Law 4 channel; JSON `trust` stays `untrusted_data`; no TrustSystem elevate | **Pass** — `packet.go` honor/intent banners; `TestDecisionMarkdownTrustLabels` asserts `untrusted_data`, no “do not treat as authority” on decision path, no system elevate |
| 4 | **DF-42** — CLI/MCP `discovery-mentions-task` → store `discovery_mentions_task`; multi-goal DPC via link; G19 thin adapters | **Pass** — domain `LinkDiscoveryMentionsTask`; CLI/MCP switch→domain only; `TestLinkDiscoveryMentionsTask` + `TestLinkDiscoveryMentionsTaskCLI`; MCP parity in write-tool smoke; `TestWhyTaskDPCMultiGoalNoForeignPollution` uses domain link |
| 5 | DF-19 goal-scope + DF-27 title intent retained; no forbidden architecture | **Pass** — multi-goal DPC foreign filter still green; DF-27 title “recorded user decision” retained; no mig/daemon/HTTP/embeddings/full-rebuild |
| 6 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01–S04 + Gate C `dry_run:false` | **Pass** — locked CGO0/CGO1 + product suites green; Gate C metrics still `dry_run:false` |
| 7 | Board Notes accurate; planner row had no product Go | **Pass** — P11-S05-00 Notes claim no product Go; P11-S05-01 Notes match live APIs/tests |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | Broader TaskContext budget test still OR-accepts legacy “do not treat as authority” phrasing | Residual OK — dedicated DF-48 test is strict; product decision path uses Law 9+4 copy |
| low | Full-module `go test ./...` may still fail setup under `similar projects/graphify` space path | Pre-existing non-product; product pkgs PASS |

## Residuals (explicit)

1. Product packages (`./cmd/... ./internal/... ./evals/...`) PASS; research `similar projects/` space-path setup fail remains residual OK.
2. DOGFOOD still lists DF-49/35/48/42 as **scheduled** → S05 (status flip deferred to phase VERIFY / findings closeout).

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS (product)
Named: TestWhySymbolExact / TestGetSymbolByID / TestExpandDepth2NoSiblingTaskBody / TestExpandContextDepth2NoSiblingTaskBody / TestDecisionMarkdownTrustLabels / TestLinkDiscoveryMentionsTask / TestLinkDiscoveryMentionsTaskCLI / TestWhyTaskDPCMultiGoalNoForeignPollution → PASS
Gate C: docs/verification/gate-c-x0/metrics-{b0,g1}.json dry_run:false intact
```

## Next

**P11-S06-00**
