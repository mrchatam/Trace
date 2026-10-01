# Impact-walk seed follow-ups — agent prompt

Paste-ready prompt for the **three items left over** after `bf524c1`
("resolve file:<path> impact-walk seeds"). That commit is good work and should
be merged; this file is the residue, not a re-litigation of it.

Found by verifying `bf524c1` from outside Trace, from a consumer project
(`mostashar`) that reached the original defect. The reporting project is
`/home/ali/Desktop/mostashar`; its paths are absolute and only used for
reproduction.

**Start here — a correction to the original report.** The first report claimed
the not-found message's advice (`get one with 'trace search <term>'`) could not
work, because `trace search` was thought to return no file entities. **That was
wrong, and this prompt is written from the corrected position.** `fts_docs` does
carry `file` rows — 124 of them in the reporting store:

```
$ sqlite3 /tmp/rp/.trace/trace.db "SELECT entity_type, count(*) FROM fts_docs GROUP BY entity_type;"
symbol|958
file|124
```

So the advice is sound *when the term is a path fragment*:

```
$ trace search 'case-graph/src' --limit 4
[('file','3a83b93b…'), ('file','403ff06c…'), ('file','4c4d9e03…'), ('file','7e185009…')]
```

It silently fails for a *symbol* term, because every `file` row has an empty
`body` and matches on `title` (the path) only. That is R1 below — the message is
not wrong, it is under-specified at exactly the moment the caller is stuck.

**Severity ordering is deliberate.** R1 costs a stuck caller a search round trip
or two. R2 is the same raw-error leak `bf524c1` was written to remove, still
present on the sibling branch. R3 is an ergonomics gap with a correctness edge.
**The hardest rule here is that R2 must not change what a symbol seed resolves
to — only what it says when it fails.** See *Hard rules*.

## How to run

1. Open a fresh agent session in the Trace repo root: `/home/ali/Desktop/Trace`
2. Paste everything under **Prompt (copy below)**.
3. Reproduce each item before fixing it. A defect never observed to fail is
   unproven.
4. Work on a **branch**. Do not commit to `main` without explicit human
   approval.

---

## Prompt (copy below)

```text
You are finishing three small items left over after `bf524c1`, which made
`file:<path>` impact-walk seeds resolve and removed a raw store error from the
file-seed branch. `bf524c1` is correct and is not up for revision — read it
first, and keep its behaviour.

## Repo and conventions (read first)

- Repo root: /home/ali/Desktop/Trace  (Go; module github.com/mrchatam/Trace)
- `AGENTS.md` — stack and hard boundaries
- `docs/rules/agent-loop-protocol.md` — the session-start gate, prompt
  taxonomy, and board authority. Obey it.
- `docs/rules/project-rules.md` — commit/secret rules

### Where the code is

Located by `grep -rn "resolveSeed\|NormalizePath" --include=*.go internal/` —
verified at commit bf524c1. Re-run it; if the file list differs, trust the grep.

| Path                                     | Holds                                              |
| ---------------------------------------- | -------------------------------------------------- |
| `internal/retrieval/impact_walk.go`      | `resolveSeed` (~203–230); the message at ~226      |
| `internal/store/entities.go`             | `NormalizePath` (~330)                             |
| `internal/store/fts.go`                  | `SearchFTS` (~44) — reads `fts_docs`               |
| `internal/retrieval/impact_walk_test.go` | `TestImpactWalkFileSeedByPath` (~362) is the pattern |
| `cmd/trace/cli_test.go`                  | `TestImpactWalkCLIPathSeed` (~1556) is the pattern  |

---

## R1 — the not-found message's advice only works for path-shaped terms

**Symptom.** `resolveSeed` ends its message with:

```
…; get one with `trace search <term>` and use the entity_id
```

`fts_docs` really does index files, so this is **not** false advice. It is
incomplete in the one place it matters. A `file` row has an **empty `body`** and
matches on `title` (the path) only — all 124 rows in the reporting store. So:

```
$ trace search 'case-graph/src' --limit 4
→ four file hits with entity_ids      (works)

$ trace search 'assertNodeInScope' --limit 4
→ ['claim','evidence','symbol']       (no file hit at all)
```

A caller who mistyped a path, then searched the *symbol* they were thinking
about, gets nothing — and concludes the message lied to them.

**Fix.** Make the term explicit and give a fallback that always works:

```
retrieval: ImpactWalk: seed file:<uuid-or-path> not found — %q matched no row in
`files` (by path) nor in `files.id` (by uuid). File ids come from a path-shaped
search, e.g. `trace search "src/store"` (file rows carry the path in `title`,
not `body`); `trace index status` shows what is indexed.
```

If you would rather not name a specific search shape, drop the advice clause
rather than leave a vague one — a wrong shape is worse than none.

## R2 — symbol seeds still leak the raw store error

**Symptom.** `bf524c1` rewrote the message for `file` seeds. The `symbol` branch
still emits a store implementation detail through a user-facing tool:

```
$ trace impact walk --seed symbol:nope --depth 1
impact: retrieval: ImpactWalk: seed symbol:nope: store: symbol id "nope": sql: no rows in result set
```

`resolveSeed` (`impact_walk.go:212`) short-circuits before the path fallback
whenever `typ != "file"`, so a symbol miss keeps `%w`-wrapping `sql.ErrNoRows`.

**Fix.** Give the symbol branch the same treatment — a message that names the
accepted form, the table, and the id space:

```
retrieval: ImpactWalk: seed symbol:<uuid> not found — %q matched no row in
`symbols` (by id); symbol seeds are uuid-only, and a file path is not accepted
here
```

That last clause earns its place: after `bf524c1`, a caller who pastes a path
under `symbol:` gets a confusing miss rather than being told paths belong to the
`file:` prefix.

## R3 — path normalisation is narrower than the schema implies

**Symptom.** The MCP schema now advertises `file:<repo-relative-path>`. Three
shapes a caller will reasonably type behave inconsistently:

```
file:packages\case-graph\src\access.ts      → resolves        (backslash ok)
file:./packages/case-graph/src/access.ts    → resolves        (./ ok)
file:packages//case-graph/src/access.ts     → NOT FOUND      (double slash)
file:/packages/case-graph/src/access.ts     → NOT FOUND      (leading slash)
```

`NormalizePath` (`internal/store/entities.go:330`) only does
`ReplaceAll("\\\\","/")` + `TrimPrefix("./")` — no slash collapsing, no leading
slash. Terminals print absolute paths, so the leading-slash case will be
somebody's first attempt.

**Fix, in preference order:**

(a) **Widen `NormalizePath`** to collapse repeated separators and trim a leading
    `/`. It is the repo's single path canonicaliser, so this is where it belongs
    — but it is used by the **indexer**, not just the walk. Confirm that the
    indexer writes paths through `NormalizePath` on the way *in*, or a widened
    normaliser will make lookup and storage disagree. Check this before changing
    it; if you cannot confirm it quickly, do not do (a).
(b) If (a) is unsafe to do here, **normalise inside `resolveSeed` only**, and say
    so in a comment explaining why the shared normaliser was left alone.
(c) If you do neither, **say so in the schema** — mark the form
    `file:<repo-relative-path, no leading slash>` — and let it be a documented
    limitation rather than a trap.

(a) is right if the indexer round-trip is confirmed; (c) is an acceptable
outcome if it is not. Do not ship (b) and (c) silently together — pick one and
justify it.

## Method: reproduce before you fix

In a scratch store — **never** the consumer's real graph:

```
mkdir -p /tmp/rp/.trace
cp /home/ali/Desktop/mostashar/.trace/trace.db /tmp/rp/.trace/trace.db
```

Then, for each item, capture the exact output **before** and **after** with the
same command. For R2 and R3 that is one command each. For R1, capture both
search shapes — the message is only fixed when *both* tell the caller something
true.

**Never hand-edit a `.trace/trace.db`.** Copy it to `/tmp` to experiment.

## Hard rules

1. **Do not change what a symbol seed resolves to.** R2 is a message change. If
   you find yourself adding path resolution for symbols, that is a different
   proposal — stop and put it in `docs/TODO/` with a rationale.
2. **Do not regress `bf524c1`.** `TestImpactWalkFileSeedByPath`,
   `TestImpactWalkCLIPathSeed` and the MCP path-seed test must pass unmodified.
3. **No behaviour change without a test that failed first.** R1 and R2 are
   message changes, and the test is on the message text. R3 is a resolution
   change, and the test is on all four shapes above.
4. **No behaviour change without operator sign-off.** R3(a) touches the indexer's
   path handling; surface the plan first, per `docs/rules/agent-loop-protocol.md`.
5. **An error message must not promise a workflow that does not exist.** This is
   the whole reason R1 is filed: the previous fix shipped a true sentence that
   is false for a common input. If you are unsure a suggestion works, run it.
6. **Work on a branch.** Do not push to `main` without explicit approval.
7. **A defect you cannot fix is recorded, not worked around.**
8. **`gofmt` is a fail-fast CI gate.** Run it before you commit.

## Verify

```bash
go test ./internal/retrieval/... ./cmd/trace/... ./internal/mcp/... 2>&1 | tail -10
gofmt -l .            # must print nothing
```

Note: `go test ./...` fails on `similar projects/graphify/tests/fixtures` with
`malformed import path … invalid char ' '`. That is **pre-existing** — it
reproduces on `main` and is caused by a directory name containing a space. Do
not fix it here; do not let it look like your regression.

Then re-run each R-item's command and confirm the message is true for the inputs
it claims to cover.

## Out of scope

- **Re-opening `bf524c1`.** Its resolution path, its use of the resolved id in
  `hitKey`, and its single-branch-covers-both-surfaces design are all correct.
- **A `trace files` listing command.** That would be the ideal answer to "how do
  I get a file id", but it is a new surface and a separate decision. R1 only
  requires the message to stop misleading.
- **The consumer's docs.** Its runbooks are no longer blocked — the walk works
  from paths — but four of them still describe the old grammar. That repo owns
  them.

## Report back

Per item: the failing output before, the passing output after, the same command.
For R1, both search shapes. The commit sha, on a branch. Anything you declined,
and why — especially if you skipped R3(a) and why.
```
