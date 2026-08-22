# P11-S05-00 — Retrieval why / depth / trust / DPC attribution (FINAL)

## Metadata
- id: P11-S05-00
- todo_ids: [P11-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S05 implement/review prompts for **DF-49, DF-35, DF-48, DF-42**. Confirm live inventory; lock APIs/tests. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 (retrieved = data) + Law 9 (user decisions authoritative); G19
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1/A6; Law 4/9 DF-48 note
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-49, DF-35, DF-48, DF-42
- [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md) — DF-49/35/48/42
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md) — DF-42/35
- [experiments/POST-P10-MCP.md](../../../../../experiments/POST-P10-MCP.md) — DF-48 residual; DF-35 HELLO leak
- Phase 10 S01 FINAL: [../../../phase-10-integrity-surfaces/scopes/scope-01-retrieval-why-fidelity/00-PLANNER.md](../../../phase-10-integrity-surfaces/scopes/scope-01-retrieval-why-fidelity/00-PLANNER.md) — DF-19/27
- Live: `internal/retrieval/{exact,expand,why,discovery_plan_change}.go`; `internal/compiler/{compiler,packet}.go`; `internal/store/file_graph.go`; `internal/domain/link.go`; `cmd/trace/{why,context,link,help}.go`; `internal/mcp/{tools_write,server}.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no grill (phase homes + live inventory + Law 4/9 do not conflict).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `lookupEntity` (`exact.go`) | Has `file`, `capability`, `review`, … — **no `symbol` case** → `unknown entity type "symbol"` (DF-49). Path/name Exact can still emit symbol Hits via `ListSymbolsByPath` / FTS |
| Store symbols | `ListSymbolsByPath` / FTS only — **no `GetSymbolByID`** yet (implementer adds thin lookup + path join) |
| Expand `goal` neighbors | `ListTasksByGoalID` attaches **every sibling task with `excerpt(t.Body)`** — depth 2 from task seed = goal then siblings → body leak (DF-35; `SECRET_HANDOFF` / `HELLO:`) |
| Compiler | Layer 0 keeps seed task body; L1 copies Hit.Excerpt — redacting Expand sibling Excerpt is sufficient (no second body fetch) |
| Why | Uses Expand depth **1** — siblings not in Why today; DF-35 is **context depth≥2** |
| Decision MD (`packet.go`) | DF-27 title line OK; excerpt blockquote still **“do not treat as authority”** vs dogfood **“binding”** → DF-48 |
| JSON trust | Decisions stay `untrusted_data` (Law 4) — keep |
| DPC attribution | `discovery_mentions_task` already used in retrieval/compiler **tests via store.InsertLink**; `dpcAttribution`/`endpointTaskGoals` already honor any discovery↔task link |
| CLI `link` / MCP `trace_link` | Only `goal-task\|decision-task\|discovery-plan-change\|claim-evidence` — **no** discovery→task (DF-42) |
| Migration | `entity_links.rel` free string; symbols table exists — **none** needed |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-49 | `why symbol <id>` unknown | Exact/Why `lookupEntity` supports **`symbol`** (store Get-by-id + path); parity with `why file` |
| DF-35 | depth≥2 leaks sibling task body | Goal→task Expand neighbors: **title/id OK; body Excerpt empty** (seed task body stays Layer 0 only) |
| DF-48 | binding vs untrusted_data copy | Keep `trust=untrusted_data`; reword decision/assumption **MD excerpt banner** for Law 9 honor + Law 4 channel; no system elevate |
| DF-42 | No CLI/MCP discovery→task attr | Thin domain + CLI/MCP `discovery-mentions-task` → store rel `discovery_mentions_task` |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF home | **DF-49, DF-35, DF-48, DF-42** only (DF-19 goal-scoped DPC + DF-27 title labels remain as shipped; do not reopen) |
| Packages | **`internal/store`** (GetSymbolByID); **`internal/retrieval`** (symbol Exact/Why; sibling body redact); **`internal/compiler`** (MD trust copy); **`internal/domain`** (`LinkDiscoveryMentionsTask`); thin **`cmd/trace`** (`link`, `help`, why inherits); thin **`internal/mcp`** (`trace_link` enum/desc). **G19** — no business logic in adapters |
| Migration | **None** |
| DF-49 Exact/Why | `case "symbol":` in `lookupEntity`: load by id; Hit `EntityType:"symbol"`, `Title`=name, `Excerpt`=kind (or empty), `Path`=owning file path when resolvable. `Why("symbol", id)` and CLI/MCP `trace why symbol <id>` succeed for indexed symbols. Miss → same not-found / empty Exact behavior as other types |
| DF-49 store | Add `GetSymbolByID(id) (Symbol + path, error)` (or equiv) — JOIN `symbols`↔`files` for path; no schema change |
| DF-49 non-goals | Do not invent symbol write/CLI create; do not require FTS for Exact-by-id; `plan_scope` Exact still out |
| DF-35 rule | In Expand `case "goal"` when emitting tasks from `ListTasksByGoalID`: set **`Excerpt` empty** (or non-body metadata only — **never** `excerpt(t.Body)`). Title + id + ReasonGoalHasTask retained. Seed / Layer-0 task body unchanged |
| DF-35 compiler | Prefer Expand-side redact; if any path re-fetches sibling bodies into packets, stop that. Assert depth-2 TaskContext/ExpandContext for planner task **must not** contain sibling body marker (e.g. `SECRET_HANDOFF`) |
| DF-35 non-goals | Do not remove sibling **titles** from depth-2; do not change max depth 2; do not invent DF-28 handoff SoT |
| DF-48 trust enum | **Do not** set decision/assumption `Trust` to `system` / remove `untrusted_data` |
| DF-48 MD copy | For `decision`/`assumption` items: replace excerpt blockquote **“do not treat as authority”** with Law 9+4 shape — e.g. honor as **recorded user decision / project intent**; trust channel remains `untrusted_data` (**do not elevate body to system policy**). Keep DF-27 title line intent. Non-decision untrusted items may keep generic retrieved-text banner |
| DF-48 dogfood | Optional light clarifying note in shared natural/protocol dogfood copy only if needed for assert; **do not** mass-edit all `ab-*/TASK-*.md`. Product MD is the primary fix |
| DF-42 domain | `RelDiscoveryMentionsTask = "discovery_mentions_task"`; `LinkDiscoveryMentionsTask(discoveryID, taskID, meta)` → `InsertLink` discovery→task + `entity.linked` event (mirror decision-task pattern). Validate both entities exist |
| DF-42 CLI/MCP | CLI rel **`discovery-mentions-task`** (hyphen, matches other link verbs); MCP `trace_link` schema enum + switch add same alias → domain. Help/usage list updated. Store/JSON rel stays underscore |
| DF-42 retrieval | No algorithm change required if endpoint↔task links already drive DF-19 attribution — add/keep a test that **CLI/domain-created** `discovery_mentions_task` attributes multi-goal DPC (foreign goal still omitted) |
| DF-42 non-goals | Do not add discovery→plan_change alternate; no planner megastore; seed `from_id`/`to_id` stays S07 (DF-33) |
| Tests (required) | (1) **`TestWhySymbolExact`** (or equiv): indexed symbol Exact+Why by id succeed; unknown id miss OK. (2) **`TestExpandContextDepth2NoSiblingTaskBody`** (or equiv): two tasks same goal; depth-2 context for A must not contain B’s body marker; may still see B’s title. (3) **`TestDecisionMarkdownTrustLabels`** (update): decision excerpt MD asserts Law 9 honor / no “do not treat as authority” conflict; JSON trust still `untrusted_data`. (4) **`TestLinkDiscoveryMentionsTask`** (domain) + CLI and/or MCP link smoke. (5) Multi-goal DPC attribution via new link still green (extend `TestWhyTaskDPCMultiGoalNoForeignPollution` or sibling). (6) Carry-forward suites |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` untouched; P10 DF-19/23/25/27/29; P11-S01 DF-40; P11-S02 DF-43/44; P11-S03 DF-47; P11-S04 DF-41/51 |
| Forbidden | New mig; elevating decision trust to `system`; restoring sibling body excerpts on goal Expand; daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–10 / P11-S01–S04 `done` history; S06+ product work; DF-28 handoff SoT |

## Effects on later scopes
- **S06** (mcp-install-reload): no retrieval/trust coupling — serial after S05 review only. Light Depends note on S06 stubs.
- **S07**: DF-33 seed link aliases unrelated; discovery-mentions-task is product link verb (not seed `from_id`).
- **S08 VERIFY:** include DF-49 why-symbol, DF-35 depth-2 no sibling body, DF-48 MD Law 9+4, DF-42 CLI/MCP discovery-mentions-task + multi-goal attribution in evidence table.

## Exit
- [x] Thicken `01-retrieval-why-depth-trust.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S06)
- [x] Board Notes; next **P11-S05-01**
- [x] Product Go — **not** this row
