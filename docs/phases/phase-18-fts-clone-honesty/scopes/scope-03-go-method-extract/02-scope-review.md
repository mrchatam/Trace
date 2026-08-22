# P18 / S03 / 02 — Go method extract review

## Metadata
- id: P18-S03-02
- todo_ids: [P18-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated
- hooks: []

## Objective
Independent review of **DF-89** golden vs sibling [00-PLANNER.md](00-PLANNER.md) **FINAL**. Write `REVIEW-NOTES.md` (landed `func Test*` names for S04). Next **P18-S04-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT** (it is **FINAL**). Fresh session ≠ implementer.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-go-method-extract.md](01-go-method-extract.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [docs/TODO.md](../../../../TODO.md)
- S02 (no coupling; do not reverse exclude): [../scope-02-clone-pending-honesty/REVIEW-NOTES.md](../scope-02-clone-pending-honesty/REVIEW-NOTES.md)

## Session start
Follow agent-loop-protocol. Re-run verify; do not trust Notes. Do not re-open S01 FTS or S02 exclude. Do not edit S05 board rows. Do not add why-by-name.

## Checklist

| # | Check | How |
|---|--------|-----|
| 1 | Named golden exact list is `method:Search`, `method:SearchCursor`, `type:Memory`, `type:Notes` | Read `TestIndexFileGoHandlerMethods` |
| 2 | Fixture is 00 locked handler shape (`package testdata`; `func (n *Notes) Search(w http.ResponseWriter, r *http.Request)`; `func (m *Memory) SearchCursor(...)`; types Notes+Memory only) | Read `testdata/handler_methods.go` vs 00 locked source |
| 3 | Fixture path is `internal/analyzers/testdata/handler_methods.go`; virtual index path `pkg/handler_methods.go` | Test + testdata |
| 4 | Existing `TestIndexFileGoGolden` still PASS; `sample.go` / keeper `wantSym` untouched | Re-run + read |
| 5 | `extract_go.go` unchanged **unless** named test was red; if changed, only query/extract (file-local incremental stays) | Diff |
| 6 | No why-by-name CLI; no FTS/seed/MCP; no full-rebuild indexer; DF-88 exclude not reversed | Diff / grep |
| 7 | CGO analyzer tests green on locked verify | `CGO_ENABLED=1` re-run |
| 8 | S05 rows still pending after VERIFY; no binary rebuild in this scope | Board P18-S05-00/01/02 |

Reject-if (blocker/high): named test missing the two methods; fixture is not handler-shaped; keeper mutated; `sample.go` used as the DF-89 vehicle; extract drive-by while golden was already green; full-rebuild indexer; why-by-name; reversing DF-88.

## Locked verify (re-run)

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoHandlerMethods|TestIndexFileGoGolden'
```

Do **not** use CGO0 as this scope’s bar. Do **not** use a bare `-run 'TestIndexFileGo'`.

## REVIEW-NOTES.md (required)

Record: verdict; checklist; landed `func Test*` names (must include `TestIndexFileGoHandlerMethods`); keeper `TestIndexFileGoGolden`; CGO1 evidence; whether `extract_go.go` was touched; residuals (CGO0 analyzers carry-forward is non-fail). S04 imports names from this file.

## Exit criteria
- [ ] REVIEW-NOTES.md; confidence high or medium with residuals listed
- [ ] Landed test names recorded for S04 import
- [ ] No open blocker/high without a pending spawn
- [ ] Board Notes; next **P18-S04-00** unless spawn (S05 still after S04)

## Minimal todos
- [ ] Verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
