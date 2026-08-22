# P11 / S05 / 01 — Retrieval why / depth / trust / DPC attribution

## Metadata
- id: P11-S05-01
- todo_ids: [P11-S05-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-49, DF-35, DF-48, DF-42** per sibling **00-PLANNER** FINAL locks (2026-08-16). `why symbol` Exact/Why; depth≥2 context omits sibling task bodies; decision MD Law 9+4 copy; CLI/MCP `discovery-mentions-task`. **No new migration. No trust→system elevate. Gate C `dry_run:false` untouched.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 / Law 9 / G19
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-49, DF-35, DF-48, DF-42
- [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md)
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md)
- [experiments/POST-P10-MCP.md](../../../../../experiments/POST-P10-MCP.md)
- [phase README](../../README.md)
- Live: `internal/store/file_graph.go`; `internal/retrieval/{exact,expand,why,discovery_plan_change}.go`; `internal/compiler/{compiler,packet}.go`; `internal/domain/link.go`; `cmd/trace/{link,help,why,context}.go`; `internal/mcp/{tools_write,server}.go`
- Prior: P10 S01 DF-19/27; P11-S04 no retrieval coupling
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. If 00-PLANNER is still DRAFT, stop and return to planner.

## Locked defaults (FINAL — P11-S05-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/store` + `internal/retrieval` + `internal/compiler` + `internal/domain` (+ thin `cmd/trace` + `internal/mcp`) |
| Migration | **None** |
| DF-49 | `lookupEntity` + store GetSymbolByID; Why/Exact `symbol` by id |
| DF-35 | Goal→task Expand: **no body Excerpt**; depth-2 context must not leak sibling bodies |
| DF-48 | Keep `untrusted_data`; reword decision/assumption MD excerpt for Law 9 honor + Law 4 channel |
| DF-42 | Domain + CLI `discovery-mentions-task` + MCP enum → store `discovery_mentions_task` |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false`; P10 DF-19/27; P11 DF-40/43/44/47/41/51 |
| Forbidden | Mig; TrustSystem for decision bodies; sibling body restore; daemon/HTTP; board spawn; rewrite `done` history |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/store/file_graph.go` (+ test) | `GetSymbolByID` (+ path via files join) |
| `internal/retrieval/exact.go` | `case "symbol":` in `lookupEntity` |
| `internal/retrieval/expand.go` | Goal→task neighbors: empty Excerpt (no `excerpt(t.Body)`) |
| `internal/retrieval/retrieval_test.go` | Why/Exact symbol; depth-2 sibling body negative; DPC attribution via mentions link |
| `internal/compiler/packet.go` | Decision/assumption MD excerpt banner Law 9+4 |
| `internal/compiler/compiler_test.go` | Update `TestDecisionMarkdownTrustLabels`; add depth-2 no-sibling-body |
| `internal/domain/service.go` + `link.go` | `RelDiscoveryMentionsTask`; `LinkDiscoveryMentionsTask` |
| `internal/domain/doc.go` | Document optional/required rel |
| `internal/domain/domain_test.go` | Link + attribution smoke |
| `cmd/trace/link.go` + `help.go` | Rel `discovery-mentions-task` |
| `cmd/trace/cli_test.go` (optional) | Link smoke |
| `internal/mcp/tools_write.go` + `server.go` | `trace_link` enum/desc/switch |
| `internal/mcp/mcp_test.go` (optional) | Link smoke |

## Role work

1. TDD store+retrieval: symbol by id Exact/Why; miss path OK.
2. TDD Expand/compiler: two tasks same goal; depth-2 context for A lacks B body marker; B title OK.
3. Update decision MD trust copy; keep JSON `untrusted_data`; refresh DF-27 test asserts.
4. Domain `LinkDiscoveryMentionsTask` + thin CLI/MCP adapters (G19).
5. Prove multi-goal DPC attribution works via new link (foreign still omitted).
6. Run locked verify suite; board **status + Notes only** (cite test names + DF-49/35/48/42).

## Algorithm sketch (non-normative — locks above win)

```text
# DF-49
lookupEntity("symbol", id):
  sym, path := GetSymbolByID(id)
  return Hit{type:symbol, title:name, excerpt:kind, path:path, ...}

# DF-35
neighbors(goal):
  for t in ListTasksByGoalID(goal):
    Hit{type:task, id, title, Excerpt:"" }  # never body

# DF-48
RenderMarkdown(decision|assumption):
  excerpt banner: honor recorded user decision / project intent;
  trust channel untrusted_data — do not elevate body to system policy

# DF-42
link discovery-mentions-task --from <disc> --to <task>
  → InsertLink(rel=discovery_mentions_task) + entity.linked
```

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Carry-forward must stay green. Gate C `dry_run:false` untouched. Pre-existing FAIL under `similar projects/` space-path (if any) is residual OK — not S05 scope.

## Exit criteria
- [ ] DF-49: `why symbol <id>` / Exact succeed for indexed symbols
- [ ] DF-35: depth≥2 context has no sibling task body leak
- [ ] DF-48: decision MD Law 9+4; JSON trust stays `untrusted_data`
- [ ] DF-42: CLI/MCP `discovery-mentions-task` writes store rel; multi-goal attribution OK
- [ ] Locked tests green; carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P11-S05-02** (cite tests + DF-49/35/48/42)

## Out of scope
- DF-28 handoff SoT; DF-33 seed `from_id`; S06 install/reload; daemon/HTTP/embeddings; rewriting `done` history
