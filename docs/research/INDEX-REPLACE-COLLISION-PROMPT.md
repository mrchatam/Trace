# Index replace collision & the missing test fallback — agent prompt

Paste-ready prompt to fix **three defects in Trace's code-graph indexer and test
selector**, found from outside Trace by a consumer project (`mostashar`) that
indexes its whole tree on every phase.

Everything this prompt needs is in this file. The reporting project is
`/home/ali/Desktop/mostashar`; paths into it are absolute and only used for
reproduction.

**Lead with the number: on the reporting repository, 20 of 145 source files are
silently absent from the graph, spanning 21 commits and 7 phases.** Every one of
those phases ran, `pnpm check` was green, and nothing said otherwise. The
indexer has been failing since commit `2118d92b` and every phase since has
recorded its work into a graph that quietly did not contain it.

**Read D1 first — it changes what the other two mean.** D1 is **not a new bug**.
It is a regression that defeats a fix this repo shipped and closed on the board
in Phase 22 (row 323, `P22-S01-04`). That row says, in its own words:

> **high** inline-fixed: unique-collide on SET NULL **confirmed** — two
> named-import `validates` from one test (`import { alpha, beta }`) to one lib
> made `IndexFile` of the lib fail (`idx_code_edges_unique`). `ReplaceFileSymbols`
> was DELETE-all then INSERT. Fix: upsert-first (stable deterministic ids keep
> incoming FKs) + collapse leftover incoming before leftover DELETE.

The three regression tests that row names **all pass on current `main`** (two in
`internal/store`, one in `internal/analyzers`), and a real repository still
fails with the identical constraint. So the question is not
"write a fix" — it is **"what invariant do those three tests encode that the
real case violates?"** A fourth test that reproduces the real case is the
deliverable; the fix follows from it.

D2 is a missing safety net that is invisible in a Go repo. D3 is
discoverability. Neither is urgent next to D1.

## How to run

1. Open a fresh agent session in the Trace repo root: `/home/ali/Desktop/Trace`
2. Paste everything under **Prompt (copy below)**
3. Reproduce D1 before fixing it. A defect never observed to fail is unproven —
   that rule is the reason this file exists
4. Work on a **branch**. Do not commit to `main` without explicit human approval

---

## Prompt (copy below)

```text
You are fixing three defects in Trace's code-graph indexer (`trace index`) and
its test-target selector (`trace test run`). They were found by a consumer
project, not by Trace's own dogfooding, which is why they survived.

## Repo and conventions (read first)

- Repo root: /home/ali/Desktop/Trace  (Go; module github.com/mrchatam/Trace)
- `AGENTS.md` — stack and hard boundaries
- `docs/rules/agent-loop-protocol.md` — the session-start gate and board
  authority. Obey it.
- `docs/rules/project-rules.md` — commit/secret rules
- `docs/ARCHITECTURE.md` — the code-graph surfaces these defects sit in
- `docs/TODO/phase-22.md` row 323 (`P22-S01-04`) — **read this before you start
  on D1.** It records the same constraint collision as already diagnosed and
  fixed. Your job is to find out why that fix does not hold.

### Where the code is

Verified to reproduce at commit `0032a20`. Re-run the greps; if a file list
differs, trust the grep over this table.

| Path                                  | Holds                                              |
| ------------------------------------- | -------------------------------------------------- |
| `internal/analyzers/index.go:105,108` | the per-file order: `ReplaceFileEdges` → `ReplaceFileSymbols` |
| `internal/store/file_graph.go:329`    | `ReplaceFileSymbols` — upsert-first, then delete leftovers |
| `internal/store/file_graph.go:370,379` | the two `DELETE FROM symbols` statements that trigger the FK |
| `internal/store/file_graph.go:423`    | `collapseEdgesTargetingSymbols` — the attempted mitigation |
| `internal/store/file_graph.go:652`    | `ReplaceFileEdges` (outgoing only)                   |
| `internal/store/schema/022_code_relationships.sql:18-25` | `idx_code_edges_unique`        |
| `internal/testrun/select.go:234,241`  | `packageFallbackTargets` — **the D2 defect**         |
| `internal/testrun/run.go:73`          | the `no relevant tests selected` error               |
| `internal/store/file_graph_test.go:229,278` | two of the three tests the board row names — **read these first** |
| `internal/analyzers/analyzers_test.go:674` | the third (`TestValidatesMultipleNamedImportsSurviveTargetReindex`) |

---

## D1 — `trace index` aborts the whole run, and a shipped fix does not hold

**Symptom.** On a real consumer repository, a whole-tree index fails
deterministically:

```
$ trace -C <repo> index
index: analyzers: replace symbols: store: delete leftover symbols:
constraint failed: UNIQUE constraint failed: index 'idx_code_edges_unique' (2067)
$ echo $?
2
```

**`--force` does not help** — it re-indexes from scratch and hits the same
constraint, so there is no re-run that clears it and no repair command:

```
$ trace -C <repo> index --force
index: analyzers: replace symbols: ... idx_code_edges_unique (2067)
$ echo $?
2
```

**Ground truth — the index silently truncates, and has for 21 commits.** On a
clean copy of the reporting repo, the store is stuck at commit `2118d92b`:

```
$ sqlite3 .trace/trace.db "select count(*) from files;"   -> 127
```

Measured against git, **20 of 145 tracked source files are absent**:

```
source files in git : 145
indexed in store    : 125
MISSING             : 20
```

The missing 20 include P14's `packages/deterministic/src/toolkit.ts` and
`services/agent-runtime/src/tool-step.ts` — not just the files added by the most
recent phase. The two files that fail to index are
`packages/agents/src/index.ts` (last changed in **P14**) and
`services/update-subsystem/src/index.ts` (changed in the most recent phase).
Neither is exotic; both are barrel files that a test file imports five and two
symbols from respectively.

Nothing reports this. `index status` correctly says `"stale": true`, and the
index command exits 2 — but every downstream command sees a **smaller world**
and returns an empty answer rather than an error. That is the real severity:
not "indexing is broken" but **"a graph that has been quietly losing commits for
21 commits, across 7 phases, and reporting success the whole time."** A consumer
that trusts its graph is reasoning about code that is not there.

**Reproduction, verified end to end:**

```bash
# 1. Build the current binary
go build -o /tmp/probe/trace ./cmd/trace

# 2. Copy the reporting repo — NEVER index the original; it holds real case memory
cp -a /home/ali/Desktop/mostashar /tmp/ms-probe

# 3. Confirm the failure
/tmp/probe/trace -C /tmp/ms-probe index ; echo "EXIT=$?"     # -> 2, constraint message

# 4. Confirm --force is not an escape
/tmp/probe/trace -C /tmp/ms-probe index --force ; echo "EXIT=$?"

# 5. Confirm the silent truncation, measured against git
python3 - <<'EOF'
import subprocess, sqlite3
tracked = subprocess.run(['git','ls-files'],capture_output=True,text=True).stdout.split()
c=sqlite3.connect('.trace/trace.db')
have={r[0] for r in c.execute("select path from files")}
src=[p for p in tracked if p.endswith(('.ts','.go','.py','.tsx','.js'))]
missing=[p for p in src if p not in have]
print(f"in git {len(src)} / indexed {len(src)-len(missing)} / MISSING {len(missing)}")
EOF
# -> in git 145 / indexed 125 / MISSING 20

# 6. The two files that trigger it, found by per-file bisection
for f in $(git -C /tmp/ms-probe ls-files '*.ts'); do
  /tmp/probe/trace -C /tmp/ms-probe index "$f" 2>&1 | grep -q idx_code_edges_unique \
    && echo "TRIGGER: $f"
done
# -> packages/agents/src/index.ts
# -> services/update-subsystem/src/index.ts
```

**The trigger is not "a file that was edited."** One of the two was last
touched in **P14**, two phases before the change that surfaced this. What both
share is a shape in the *graph*, not in the file:

| File                                        | incoming `validates` from one test file |
| ------------------------------------------- | --------------------------------------- |
| `packages/agents/src/index.ts`              | **5** (from `agents.test.ts`)           |
| `services/update-subsystem/src/index.ts`    | **2** (from `update-subsystem.test.ts`) |

Six other files have *more* incoming `validates` than either of these —
`packages/tier-gate/src/index.ts` has 8, `packages/billing/src/index.ts` 7,
`services/law-mcp/src/tools.ts` 7 — and **none of them collides.** So the edge
count alone is not the trigger.

The working hypothesis, which you should confirm or refute rather than adopt:
**collision requires ≥2 incoming same-`rel` edges into a file AND a re-index
that leaves leftover symbols to delete.** The six non-colliding files have stable
deterministic symbol ids, so nothing is leftover and the `DELETE` never fires.
The two that collide have symbol sets that changed since the store was last
built (`content_hash` in the store does not match the file on disk). Confirm
this before you rely on it — it is a hypothesis derived from a reading, not
something anyone has measured.

**Note what this means for severity.** The defect was latent through P14 and
surfaced only when a later change happened to disturb a second file. Nothing
about P14's edit was unusual. A consumer cannot distinguish "indexed and
correct" from "indexed 125 of 145", and a fix that only handles the observed
pair will leave the next one latent.

**The mechanism, as read from the source.** `code_edges.to_symbol_id` is
`REFERENCES symbols(id) ON DELETE SET NULL` (migration 022, line 10), and
`idx_code_edges_unique` (lines 18-25) is:

```sql
CREATE UNIQUE INDEX IF NOT EXISTS idx_code_edges_unique
    ON code_edges(
        from_file_id,
        IFNULL(from_symbol_id, ''),
        to_file_id,
        IFNULL(to_symbol_id, ''),
        rel
    );
```

So two edges that differ **only** in `to_symbol_id` are distinct today, and
become identical the instant one of their targets is deleted — `IFNULL` maps the
NULLed side to `''`. Any file with **two or more edges of the same `rel` into
the same target file** collides as soon as one of those targets is deleted.

`ReplaceFileSymbols` deletes leftovers at lines 370 and 379, and
`collapseEdgesTargetingSymbols` (line 423) is the mitigation: it de-duplicates
edges whose `to_symbol_id` is in the leftover set, grouping by
`(from_file_id, IFNULL(from_symbol_id,''), to_file_id, rel)`.

**What is not yet pinned, and is the first thing to establish:** why the
mitigation misses the real case. The dedup only considers edges whose
`to_symbol_id` is in `leftover`. A collision between an edge that collapses to
NULL and an edge that was **already** file-level (`to_symbol_id` NULL from a
prior delete), or a group that spans the `keep`/`leftover` boundary, is not
obviously covered. Do not assume — the code comments at `index.go:102-105` and
`file_graph.go:338-341` both assert this path is handled, and it is not.

**Deliverable, in this order:**

1. **A minimal reproduction as a Go test in `internal/store`**, reduced from the
   real case until it is small enough to read. Name the invariant it violates.
2. **A failing test proving the three existing tests do not cover it.** The
   three span two packages, so name both — scoping to one silently skips a test
   and the "all green" is then partly a lie:

   ```bash
   go test ./internal/store/... ./internal/analyzers/... -count=1 \
     -run 'TestValidatesMultipleNamedImportsSurviveTargetReindex|TestReplaceFileSymbolsKeepsIncomingValidatesOnStableIDs|TestReplaceFileSymbolsCollapsesIncomingWhenSymbolDropped' -v
   ```

   On unmodified `main` that prints three `--- PASS` lines. Your new test must
   fail while they pass. Verified at `0032a20`:
   `ok internal/store 0.062s`, `ok internal/analyzers 0.052s`.
3. **Then the fix.**

**Constraints on the fix:**

- **Do not drop or weaken `idx_code_edges_unique`.** It is load-bearing: it is
  what makes "two named imports of the same symbol" idempotent. A fix that
  relaxes the index is the wrong fix; find the real ordering bug.
- **The whole-tree walk must not abort on one bad file.** Even with the
  collision fixed, an indexer that loses 7 files and every edge touching them
  because of one record is a **fail-closed-shaped, fail-silent-shaped** design
  question in its own right. Consider whether the walk should collect and
  report per-file failures, and land that behaviour explicitly. If you decide
  the current abort is correct, say so with a reason — but a consumer cannot
  distinguish "indexed and correct" from "indexed 127 of 134 files".
- Prefer a transactional, order-independent guarantee over a patch for the
  observed case: the symptom depends on the order the walker happens to visit
  files, which is exactly the kind of thing that reappears.

---

## D2 — the test-target fallback is Go-only, so a TypeScript repo has none

**Symptom.** On a TypeScript repository, `test run` fails closed with:

```
$ trace -C <repo> test run --task <uuid>
test run: no relevant tests selected
$ echo $?
2
```

**Ground truth.** `SelectTestTargets` (`internal/testrun/select.go:29`) tries
`targetsFromImpactAndValidates` and, when that yields nothing, falls back to
`packageFallbackTargets` (line 234). The fallback skips every non-Go path:

```go
for _, p := range paths {
    if !strings.HasSuffix(p, ".go") {
        continue
    }
    ...
}
```

In a Go repo the fallback rescues the run. In a TypeScript repo it returns an
empty slice and the run fails with no further explanation — so the two
languages get **opposite** outcomes from the same code path, and the message
names neither.

**This is independent of D1 — proved, not assumed.** On a store where the
workaround above has recovered the index to **145 of 145 files, 0 missing**,
`test run` still fails:

```
$ sqlite3 … # MISSING 0
$ trace -C <repo> tests verifying --file packages/deterministic/src/index.ts
{"count":4,"ok":true,"tests":[…]}                 # the graph knows these tests
$ trace -C <repo> test run --task <uuid>
test run: no relevant tests selected                # and refuses to use them
```

The graph has the edges; the selector does not use them. Supplying `--paths`
makes the same command succeed, which is what makes this a fallback problem
rather than an indexing one. Fixing D1 alone will **not** fix D2, and a
verification that only re-runs `trace index` will not notice.

Failing closed is correct here. Failing closed while naming **neither the reason
nor the remedy** is D2: a caller cannot tell "the graph has no `validates` edges
for these files" from "the index is stale" from "this repo has no tests for that
package" — and, as above, cannot tell it is the last of those when the first
has already been ruled out.

**Fix, all three are legitimate — pick and justify in the commit:**

(a) **Make the message name the reason.** Say which stage produced zero targets
    and what would change it ("no `validates` edges and no package fallback for
    non-Go paths"). Cheapest fix, no behaviour change.
(b) **Give non-Go languages a package fallback** — for a JS/TS package, the
    sibling `test/` directory is the same idea `goPackageArg` implements.
(c) **Distinguish "no tests exist" from "the graph does not know about them"**,
    e.g. by checking whether the seed path is present in the store at all, and
    saying so. This one is the most valuable and the most work.

(a) alone is a real fix. If you skip (b) or (c), say why.

---

## D3 — the index error names a constraint, not the file that caused it

**Symptom.** The D1 error names an index and a table but not the path, and
there is no escape hatch for the files that trigger it:

```
index: analyzers: replace symbols: store: delete leftover symbols: ...
```

A caller debugging this must bisect the whole tree by hand, one file at a
time, to find the trigger. The reporting session did exactly that.

**Fix.** Name the path in the error. `IndexFile` already knows it — thread it
into the wrapped error at `internal/analyzers/index.go:108`, where
`ReplaceFileSymbols` is called and its error wrapped:

```go
if err := stx.ReplaceFileSymbols(path, symbols); err != nil {
    return fmt.Errorf("analyzers: replace symbols: %w", err)
}
```

A caller currently has to bisect the entire tree one file at a time to find the
trigger, and on a 145-file repo that is 145 invocations. The path is known at
the wrap site; naming it costs one `%s`.

**There is no per-file escape hatch — check this before you promise one.**
It is tempting to conclude that `trace index <one-file>` is a workaround,
because it works for most files. It does not work for the ones that matter:

```bash
trace -C <repo> index packages/agents/src/index.ts
# -> same constraint failure
trace -C <repo> index services/update-subsystem/src/index.ts
# -> same constraint failure
```

The only thing that works is indexing every file **except** the triggers, which
a consumer cannot discover and would never think to try:

```bash
# verified on the reporting repo: 145 tracked source files
for f in $(git ls-files '*.ts'); do
  case "$f" in
    packages/agents/src/index.ts|services/update-subsystem/src/index.ts) continue ;;
  esac
  trace -C <repo> index "$f"
done
# -> in git 145 / indexed 145 / MISSING 0
```

That recovers 20 files and leaves the two triggers permanently stale. Land it
only as a diagnostic note for whoever hits D1 next, and only alongside the D1
fix — documenting it as *the* answer would be worse than the current error,
because it looks like a solution and silently freezes two files forever.

---

## Method: reproduce before you fix

For each defect, in a scratch store — **never** the consumer's real graph:

1. `cp -a /home/ali/Desktop/mostashar /tmp/ms-probe` (copy; never index the
   original — it holds real case memory)
2. Point a locally-built `trace` at the copy and run the failing command
3. Capture the exact failing output
4. Write a Go test that **fails** on current `main`
5. Fix; re-run; confirm the test passes
6. Capture the passing output with the same command

**Never hand-edit a `.trace/trace.db`.** Copy it to `/tmp` to experiment.

**Find what *writes* a constraint violation; do not infer it from a comment.**
The comments at `index.go:102-105` and `file_graph.go:338-341` both assert this
path is already handled, and it is not. A prior session read
`plan_critiqued`, assumed it meant "a critique was recorded", and shipped a
diagnosis that a later session had to formally supersede. Read the code that
assigns the field, and believe the failing test over the comment.

## Hard rules

1. **Never weaken a uniqueness constraint to make a test pass.** `idx_code_edges_unique`
   is the thing preventing duplicate edges. If you cannot fix D1 without touching
   it, D1 is not fixed.
2. **No behaviour change without a test that failed first.** For D1 that test is
   the deliverable, not a follow-up.
3. **Do not "fix" D1 by making the indexer skip the file.** That converts a loud
   failure into a silently missing file, which is the worse defect. If you
   believe a per-file error report is the right shape instead of a hard abort,
   that is a design decision — write the rationale.
4. **No behaviour change without operator sign-off.** Surface the D1 plan to the
   human before implementing; it touches an indexer invariant.
5. **Work on a branch.** Do not push to `main` without explicit approval.
6. **A defect you cannot fix is recorded, not worked around** — exact failing
   output and a reason, in the commit message or a `docs/TODO/` row. A silent
   workaround is the failure mode here.
7. **`gofmt` is a fail-fast CI gate.** Run it before you commit.
8. **Do not edit the consumer repo.** Every path into `/home/ali/Desktop/mostashar`
   is read-only reproduction material. If you find a defect *there*, write it up
   in your report.

## Verify

```bash
go test ./... 2>&1 | tail -20
gofmt -l .            # must print nothing
```

Then, per defect, the before/after pair from **Method**. For D1, the
before/after is the whole-tree `trace index` on the copied consumer repo, and
the `select count(*) from files` check that proves the truncation is gone.

## Out of scope

A **fourth** finding is **not Trace's** and must not be fixed from here. The
consumer's coding agent issued four `insert_line` edits; each was written
**partially** — the tail of the payload was silently dropped, with no error and
a success response, leaving a file that parsed as neither the old nor the new
content. That is a defect in the agent harness, not in Trace or in the consumer
repo, and it is recorded in the consumer's phase log. It is mentioned only
because it cost this investigation three rewrites, and because a silently
truncated write is the single most dangerous failure mode in this whole
ecosystem: it is not a syntax error, it is a paragraph that quietly stops
existing, and the tool reports success.

**Do not touch it from here.**

## Report back

Per defect: the failing output before, the passing output after, the same
command. The commit sha, on a branch. For D1 specifically: the minimal
reproduction, the invariant it violates, and which of the three existing tests
came closest to covering it — that answer is worth more than the fix. Anything
you declined to fix, and why. If a fix changes behaviour for Go consumers, say
which previously-recorded results were produced under the old behaviour and
whether they still hold.
```
