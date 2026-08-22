# Phase 18 — Context FTS + clone honesty (thin)

**Status:** **complete** (2026-08-18) — S01–S03 + S04 VERIFY green + S05 rebuild; DR-HANDOFF **`no successor`**. Human-scheduled **queue** after Phase 17 (P17 historical DR-HANDOFF remains `no successor`). Phase folder evidence: [`scopes/scope-04-phase-verify/VERIFY-NOTES.md`](scopes/scope-04-phase-verify/VERIFY-NOTES.md) + [`scopes/scope-05-rebuild-binaries/REVIEW-NOTES.md`](scopes/scope-05-rebuild-binaries/REVIEW-NOTES.md). Next runnable: **none**.

## Why this phase exists

D40 A/B (`experiments/runs/2026-08-17-ab-compare/`) hit two **Trace** defects/gaps that block honest clone+context use:

1. **`trace context` dies** on task titles containing `/` (`GET /notes`, `GET /notes/search`) — FTS5 MATCH syntax, **DF-87**.
2. Clone import lands tasks **PENDING** while live G1/G2 were DONE/SKIPPED — seed omits reviews/`work_state`. **P17 already excluded those surfaces.** This phase **does not reverse that**. It **documents expected clone PENDING** (**DF-88 wontfix** + honesty docs).

A third hunt candidate: G2 Go indexer appeared to miss `Search` / `SearchCursor`. Live `extract_go.go` already captures `method_declaration`. S03 adds a **handler-shaped golden** so that miss is either a regression fix or proven operator/`why`-by-UUID confusion.

## Investigation (2026-08-17, P18-00)

| Claim | Still true? | Evidence |
|-------|-------------|----------|
| MATCH already parameterized | **Yes** | `internal/store/fts.go` `WHERE fts_docs MATCH ?` |
| `/` left in MATCH expression | **Yes** | `sanitizeFTSQuery` replacer omits `/`; `strings.Fields("GET /notes")` → `GET`, `/notes`; FTS5 `GET AND /notes` → `syntax error near "/"` |
| Context uses title as FTS query | **Yes** | `internal/compiler/compiler.go` `c.retr.Search(ctx, task.Title, …)` — error **aborts the packet** |
| Why still works | **Yes** | `trace why` is graph/Exact by UUID — does not MATCH the title |
| P17 export exclude reviews/`work_state` | **Yes** | DF-84-FORWARD + `TestSeedExportOmitsDeniedSurfaces`; clone `transitions: 0` |
| Go methods extracted today | **Yes** (golden + G2 files) | `TestIndexFileGoGolden` `method:Run`; live `IndexFile` of D40 G2 `notes.go`/`memory.go` emits `method:Search`/`method:SearchCursor` (S03-00) |
| CANCELLED work_state | **No such state** | `legalTransitions` PENDING → SKIPPED (not CANCELLED). G2 used SKIPPED correctly |
| App HTTP 400/429 | **Not Trace** | notes-api experiment contract |
| MCP stdio EOF / rsync hang | **Harness** | `run-isolated.sh` stdin close; `run-verify.sh` rsync; product tests PASS |

**Root cause lock (DF-87):** do **not** treat this as “add SQL parameters” (already `?`). Do **not** lock “escape `/` only.” Sanitize **punctuation as token separators** (class includes `/`) and **quote remaining FTS5 tokens**, keep MATCH bound.

**S01-00 FINAL:** token charset = Unicode letter/number (`L*`/`N*`, matching unicode61); every other rune (incl. `/`) is a separator; quote tokens; `AND`-join. `GET /notes` → `"GET" AND "notes"`. TaskContext Search error → Expand-only (`fts = nil`), do not abort the packet. Named: `TestSanitizeFTSQueryPunctuationClass`, `TestSearchFTSSlashInQuery`, `TestTaskContextSlashTitle`, `TestTaskContextContinuesWhenSearchErrors`. Keepers: `TestFTSFindsEntityTitleAndPathSymbol`, `TestIncludeWhyFailClosed`. MATCH is already `?` — not an SQL-params fix.

**S02-00 FINAL:** docs/help/comments only — keep exclude. Named `TestHelpCloneTasksImportPending` (`pending` + `import` + `omits reviews` + `transitions` + `work_state` on `trace help`). Keepers `TestSeedExportOmitsDeniedSurfaces`, `TestHelpSeedExportPath`. CONTRIBUTING bullet 7 + README clone-PENDING sentence. No `--include-reviews`.

**S03-00 FINAL:** named `TestIndexFileGoHandlerMethods` on `internal/analyzers/testdata/handler_methods.go` — exact `kind:name` `method:Search`, `method:SearchCursor`, `type:Memory`, `type:Notes`. Keeper `TestIndexFileGoGolden`. Live CGO1 `IndexFile` of D40 G2 `notes.go`/`memory.go` already emits those methods (A/B miss is not a query hole). Golden-only if green; fix `goSymbolQuery` **only if** named test red. **CGO_ENABLED=1**. S05 still after VERIFY.

**S04-00 FINAL:** named DF-87/88/89 + keepers + carry-forward. **Two-clone not required** (no shell recipe; no dedicated `-run`; do not implement `TestPortableGraphTwoCloneWhyContextPlan`). DF-88 re-prove is **document-only** (`TestHelpCloneTasksImportPending` + omit/path keepers). CGO: DF-87 CGO0 authoritative + CGO1 corroboration; DF-88 `cmd/trace` and DF-89 analyzers **CGO1**. S04 **starts** DR-HANDOFF Notes; does **not** close. Stale binaries **non-fail**. S05 still **after** S04-02.

**S05 CGO:** `bin/trace` CGO1; `bin/trace-mcp` CGO0 (MCP has no analyzers). S05-01 aligned README/`help.go` MCP lines to CGO0.

**S05-00 FINAL:** rebuild **required** after VERIFY even if pre-rebuild `-h` already lists 10 (binaries mtime 2026-08-17 17:32 is pre-S01). Preferred `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` (sandbox 403 `segmentio/encoding` class — retry ≠ product defect). Catalog: `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version`. Optional DF-87 temp-dir context on `GET /notes/search`: **skip = non-fail** (Notes required); **run red = FAIL**. Optional two-line MCP CGO0 docs. Cursor MCP lag = DF-22/37 non-fail. **S05-02** closes DR-HANDOFF (`no successor`); 00/01 do not close. Not research S05.

## Disposition matrix (P18-00 FINAL)

| ID / residual | Disposition | One-line rationale | Home |
|---------------|-------------|--------------------|------|
| **DF-87** | **fix** | Context unusable on HTTP-route titles; MATCH syntax not SQL injection | **S01** |
| **DF-88** | **wontfix** + **document** | Keep P17 exclude; clone PENDING is expected. Do not explode portable-graph identity | **S02** ([DF-88-DECISION.md](DF-88-DECISION.md)) |
| **DF-89** | **fix** (golden / extract if red) | Handler methods `Search`/`SearchCursor` must appear as `method:*` | **S03** |
| DF-86 git-hook | **defer** | Not an A/B product error; stays P17 pack residual | — |
| DF-67 symbol-entity staleness | **defer** | Already deferred; not D40 | — |
| DF-22 / DF-37 MCP reload | **defer** | Ops residual; tip already shipped | — |
| DF-44 `--as-operator` ≠ identity | **wontfix** (closed) | Residual class; not reopened | — |
| PENDING→CANCELLED | **wontfix** | No CANCELLED state; SKIPPED is the legal abandon edge | — |
| App 400 unknown cursor / 429 rate limit | **wontfix** | notes-api contract, not Trace | — |
| MCP stdio harness EOF | **harness-only** | Product MCP tests PASS | — |
| verify `rsync` hang | **harness-only** | One-line: keep `cp`; do not board | — |
| `context` vs `why` “feature gap” | **depend** | After DF-87, context lexical+graph works; why stays UUID graph. No why-by-name CLI this phase | S01 |
| Hosted MCP / HTTP / OAuth | **out** | TODO Later developments | — |
| Research S05 / `plan simulate` / D21+ / ranks 7+ | **out** | Not this board (≠ phase **S05** binary rebuild) | — |

## Scope order (locked)

| Scope | Focus | DFs |
|-------|--------|-----|
| S00 / phase planner | Inventory + disposition + spawn | **done** (P18-00) |
| S01 | FTS MATCH sanitizer + context packet on slash titles + Search error does not abort packet | DF-87 |
| S02 | Clone-PENDING honesty docs + help; keep omit tests as fail bar | DF-88 |
| S03 | Go method golden (`Search` / `SearchCursor` shape) | DF-89 |
| S04 | Phase VERIFY (product evidence) | named S01–S03 |
| S05 | Rebuild `bin/trace` + `bin/trace-mcp` | ops (stale-binary SKIP) |

## Out of scope unless promoted

- Hosted/authenticated MCP; daemon/HTTP on Trace core
- Default or CONDITIONAL export of reviews / transitions / `work_state`
- `trace install git-hook` (DF-86)
- Why-by-name CLI; embeddings; full-rebuild indexer
- Adding `CANCELLED` to the work-state machine
- Rewriting Phase 00–17 `done` history; claiming P17 DR-HANDOFF was wrong
- Kitchen-sink harness (rsync, stdio EOF scripts)

## Assumptions

1. Human cut is **thin dogfood residuals** from D40, not a new product surface.
2. Portable-graph **identity stays UUIDs**; git SHA remains evidence (DF-85). Clone process history stays local.
3. VERIFY default DR-HANDOFF = **`no successor`**. S04-01 starts Notes; S04-02 re-verifies product only; **S05-02 closes**. Do **not** auto-scaffold Phase 19.
4. P17 two-clone recipe stays green; S01 must not break slash-free fixture titles.
5. After product VERIFY, **rebuild** workspace CLI+MCP so experiments/Cursor cannot SKIP on stale `bin/trace-mcp`.

## Completion bar (VERIFY + S05)

Named S01–S03 tests green + carry-forward gates. Context on a task titled `GET /notes/search` returns a packet (not FTS5 syntax error). Seed export still omits reviews/`work_state`. Go golden includes handler-shaped methods. **S05:** `CGO_ENABLED=1` `bin/trace` and `CGO_ENABLED=0` `bin/trace-mcp`; `trace-mcp -h` lists **10** tools including `trace_impact`. **Met** (P18-S05-02 **APPROVE**).

## Parallel track (not board-blocking)

Optional dogfood under `experiments/`; feed new DF-* **forward** only (next free after this phase: **DF-90**).
