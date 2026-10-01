# Impact-walk seed resolution — agent prompt

Paste-ready prompt to fix **two defects in Trace's impact-walk surface**
(`trace impact walk`, `trace_impact action=walk`), found from outside Trace by a
consumer project (`mostashar`) whose phase runbooks make the walk a **required
close step**. Both are Trace bugs.

Everything this prompt needs is in this file. The reporting project is
`/home/ali/Desktop/mostashar`; paths into it are absolute and only used for
reproduction.

**Severity ordering is deliberate.** D1 is a one-line error message. D2 is the
reason the walk could not be run at all by a caller holding a file path — and it
is a few lines, because the capability already exists and is simply unreachable.
Neither is a correctness bug: the walk does not return a wrong blast radius when
it refuses. **The hardest rule in this prompt is that a path seed must resolve to
the same node as its UUID — never to a broader or guessed match.** See *Hard
rules*.

## How to run

1. Open a fresh agent session in the Trace repo root: `/home/ali/Desktop/Trace`
2. Paste everything under **Prompt (copy below)**.
3. Reproduce each defect before fixing it. A defect never observed to fail is
   unproven — that rule is the reason this file exists.
4. Work on a **branch**. Do not commit to `main` without explicit human
   approval.

---

## Prompt (copy below)

```text
You are fixing two defects in Trace's impact-walk surface (`trace impact walk`,
and the `trace_impact action=walk` MCP tool). They were found by a consumer
project, not by Trace's own dogfooding, which is why they survived: the only
form the tests exercise is the UUID form.

## Repo and conventions (read first)

- Repo root: /home/ali/Desktop/Trace  (Go; module github.com/mrchatam/Trace)
- `AGENTS.md` — stack and hard boundaries
- `docs/rules/agent-loop-protocol.md` — the session-start gate, prompt
  taxonomy, and board authority. Obey it.
- `docs/rules/project-rules.md` — commit/secret rules
- `docs/ARCHITECTURE.md`, `docs/REVIEW_AND_VERIFICATION.md` — the retrieval
  surfaces these defects sit in

### Where the code is

Located by `grep -rn "ImpactWalk" --include=*.go .` — verified to reproduce at
commit d18523b. Re-run it; if the file list differs, trust the grep over this
table.

| Path                                     | Holds                                                |
| ---------------------------------------- | ---------------------------------------------------- |
| `internal/retrieval/impact_walk.go`      | the walk; the seed loop is at ~74-104                |
| `internal/retrieval/exact.go`            | `lookupEntity` (`file` case ~198)                    |
| `internal/retrieval/exact.go`            | `q.Path` resolution at ~45-61 — **what you want**    |
| `internal/store/file_graph.go`           | `GetFileByPath` (~188), `NormalizePath`              |
| `internal/mcp/tools_impact.go`           | MCP seed parsing (~263-276)                          |
| `cmd/trace/impact.go`                    | CLI seed parsing (~281-315) — a **duplicate** of it  |
| `internal/retrieval/impact_walk_test.go` | the walk's canonical tests                           |
| `cmd/trace/cli_test.go`                  | `TestImpactWalkCLI` (~1477) — the UUID-form pattern  |

---

## D1 — a path-shaped file seed dies on a raw store error

**Symptom.** A caller who has a file *path* — which is what every consumer's
runbook names — writes `file:<path>` and gets a storage-layer error that names
neither the accepted format nor the table:

```
impact: retrieval: ImpactWalk: seed file:packages/case-graph/src/access.ts: store: file id "packages/case-graph/src/access.ts": sql: no rows in result set
```

**Ground truth.** `internal/retrieval/impact_walk.go:82`:

```go
h, err := e.lookupEntity(typ, s.EntityID, ReasonExactID, 0, 1.0)
if err != nil {
    return nil, fmt.Errorf("retrieval: ImpactWalk: seed %s:%s: %w", typ, s.EntityID, err)
}
```

`lookupEntity`'s `"file"` case (`internal/retrieval/exact.go:198`) calls
`store.GetFileByID(id)`. `ImpactWalk` has **no path branch at all**, and the
validation above it only checks that the string has a `file:`/`symbol:` prefix
and a non-empty remainder — so a path passes validation and dies two layers
down on `%w`-wrapped `sql.ErrNoRows`.

This reproduces on **both** surfaces: the CLI prints
`impact: retrieval: ImpactWalk: …` and the MCP tool prints
`trace_impact: retrieval: ImpactWalk: …` with the same body.

Note what the message does *not* say, in order of cost to the caller: not that a
UUID was expected, not that `files.path` is a valid alternative, and not which
table was searched. `sql: no rows in result set` is a store implementation
detail leaking through a user-facing tool.

**Fix.** At minimum, resolve the seed. The resolution is nine lines away and
**already exists** in this repo: `internal/retrieval/exact.go:45-61` resolves
`q.Path` through `store.NormalizePath` + `store.GetFileByPath` and builds the
same `Hit{EntityType: "file", EntityID: f.ID, …}` that `lookupEntity` would have
returned. Prefer resolving inside `ImpactWalk` so the fix covers the CLI and the
MCP tool at once, rather than duplicating the branch into both parsers.

Either way, **name the accepted form in the error**. A caller who guessed wrong
should get:

```
retrieval: ImpactWalk: seed file:<uuid-or-path> not found — "…" matched no row in
`files` (by path) nor in `files.id` (by uuid); get one with `trace search <term>`
and use the entity_id
```

The last clause is only true if you also do D2 — otherwise do not print it.

---

## D2 — nothing surfaces a file UUID, so the documented form is unreachable

**Symptom.** The documented seed form is `file:<uuid>`. A consumer holding a
path has **no supported way to obtain that uuid**: `trace search` indexes
entities, not files (`internal/retrieval/search.go` has no file branch), and no
MCP tool or CLI command accepts a path and returns a file id. `Engine.Exact` —
the one function that resolves a path (`exact.go:45-61`) — has **no caller
outside its own tests**:

```
$ grep -rn "Exact(" --include=*.go internal/ cmd/ | grep -v _test
internal/retrieval/exact.go:16:func (e *Engine) Exact(ctx context.Context, q ExactQuery) ([]Hit, error) {
```

So the walk's required input is producible only by reading `.trace/trace.db`
directly. In the reporting session the only escape was to seed with **symbol**
ids taken from `trace_context` output and hope one lived in the file of
interest — which walked two unrelated test files. That is luck, not an API.

**Fix, in preference order:**

(a) **Resolve path seeds** (D1). This makes the surface usable and is the real
    fix. `files.path` is unique (single-row `QueryRow`), so there is no
    ambiguity question to settle.
(b) If you prefer paths stay out of the seed grammar, then **return the file
    uuid wherever a blast already names a file** — the walk's result rows carry
    `path` and `entity_id`, so a caller can at least harvest ids from a first
    walk. Say so in the tool schema.

(a) is the better fix and is small. (b) is the fallback if you decide paths
should not be a seed form — in which case D1's error must say so explicitly,
because the current one implies a path was legitimate.

---

## Method: reproduce before you fix

In a scratch store — **never** the consumer's real graph:

1. `mkdir -p /tmp/probe/.trace && cp /home/ali/Desktop/mostashar/.trace/trace.db /tmp/probe/.trace/trace.db`
   (read the copy only; never write the original — it holds real case memory)
2. Prove the file is **in** the store, or you are reproducing the wrong thing.
   This is the trap, and it is easy to hit: a probe pointed at a project root
   whose store does not index that path fails identically for the right reason
   and the wrong one.
   ```
   $ sqlite3 /tmp/probe/.trace/trace.db "SELECT id FROM files WHERE path='packages/case-graph/src/access.ts';"
   ed622487-7d8d-406d-bea4-a66d386c4992
   ```
3. Run the A/B on that same row:
   ```
   $ trace -C /tmp/probe impact walk --seed file:ed622487-7d8d-406d-bea4-a66d386c4992 --depth 1
   → blast returned (works today)

   $ trace -C /tmp/probe impact walk --seed file:packages/case-graph/src/access.ts --depth 1
   impact: retrieval: ImpactWalk: seed file:packages/case-graph/src/access.ts: store: file id "packages/case-graph/src/access.ts": sql: no rows in result set
   ```
   **A fix is proven only when the second command returns the same blast as the
   first.** A *different* blast is a new defect, not a fix.
4. Write a Go test that **fails** on current `main`
5. Fix, re-run, confirm the test passes
6. Capture the passing output with the same command

**Never hand-edit a `.trace/trace.db`.** Copy it to `/tmp` to experiment.

## Hard rules

1. **A path seed must resolve to exactly the row its path names.** No prefix
   matching, no basename fallback, no "closest" match, and never a silently
   empty result standing in for a resolution. An unresolved path is an error,
   not an empty walk — an empty blast reads as "nothing depends on this file",
   which is the one answer an impact tool must never give by accident. Fuzzy
   matching, if you want it, is a human decision with a written rationale.
2. **Run `store.NormalizePath` on every path you accept**, exactly as
   `exact.go:46` does, so its handling of `./`, leading slashes and separators
   stays in one place.
3. **A UUID seed must behave exactly as today** — same blast, same ordering,
   same `hop`/`hop_risk` accounting. `internal/retrieval/impact_walk_test.go`
   and `TestImpactWalkCLI` are the regression net; they must pass unmodified.
4. **Fix both surfaces or neither.** `internal/mcp/tools_impact.go:263-276` and
   `cmd/trace/impact.go:281-315` are independent copies of the same seed
   parsing. Resolving inside `ImpactWalk` gets both for free; duplicating the
   branch into each parser lets them drift. If you do duplicate, say why in the
   commit and add a test per surface.
5. **No behaviour change without a test that failed first.** For a resolution
   fix the test is the **A/B pair** in *Method* — same store, same file, uuid
   vs path, identical blast.
6. **No behaviour change without operator sign-off.** Surface the plan before
   implementing, per `docs/rules/agent-loop-protocol.md`.
7. **Work on a branch.** Do not push to `main` without explicit approval.
8. **A defect you cannot fix is recorded, not worked around** — write it up
   with the exact failing output and a reason, in the commit message or a
   `docs/TODO/` row. A silent workaround is the failure mode here.
9. **`gofmt` is a fail-fast CI gate.** Run it before you commit.

## Verify

```bash
go test ./... 2>&1 | tail -20
gofmt -l .            # must print nothing
```

Then re-run the A/B pair from **Method** against the fixed build and confirm the
path seed returns the **same** blast as the uuid seed.

## Out of scope

- **Changing what a seed means for `symbol:`.** `GetSymbolByID` takes a uuid and
  there is no symbol-path resolution to borrow; do not add one here.
- **A `files` listing command.** D2's fix is path seeds, or ids carried on the
  blast; a new `trace files` surface is a separate decision.
- **The consumer's phase docs.** Four of its runbooks instruct
  `trace_impact action=walk` seeded with file **paths**, which this prompt makes
  unnecessary but does not invalidate. That repo owns its own docs; do not edit
  them from here. It is mentioned only because it is how the defect was reached:
  a runbook a phase agent is *required* to follow, describing a call the tool
  cannot perform.

## Report back

Per defect: the failing output before, the passing output after, the same
command — and for D1 specifically, the **A/B blast comparison** proving the two
seed forms agree. The commit sha, on a branch. Anything you declined to fix, and
why.
```
