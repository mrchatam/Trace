# Loop gates & envelope discoverability — agent prompt

Paste-ready prompt to fix **four defects in Trace's own deliberation loop**, found
from outside Trace by a consumer project (`mostashar`) that drives the loop on
every phase. Three are Trace bugs; the fourth is a runbook defect in the consumer
and is listed here only for context — **do not fix it from this repo**.

Everything this prompt needs is in this file. The reporting project is
`/home/ali/Desktop/mostashar`; paths into it are absolute and only used for
reproduction.

**Severity ordering is deliberate.** D1 and D2 cost whole sessions to
diagnose. D3 is a one-call-per-session tax. None of them is a correctness
bug — the gates still refuse what they must. **The single hardest rule in this
prompt is that none of these fixes may widen a gate.** See *Hard rules*.

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
You are fixing four defects in Trace's deliberation loop (`trace loop`,
`trace_loop` MCP tool). They were found by a consumer project, not by Trace's
own dogfooding, which is why they survived.

All four are in the loop's *ergonomics and discoverability*, not its
correctness: the gates still refuse what they must refuse. Fix them without
changing what any gate accepts.

## Repo and conventions (read first)

- Repo root: /home/ali/Desktop/Trace  (Go; module github.com/mrchatam/Trace)
- `AGENTS.md` — stack and hard boundaries
- `docs/rules/agent-loop-protocol.md` — the session-start gate, prompt
  taxonomy, and board authority. Obey it.
- `docs/rules/project-rules.md` — commit/secret rules
- `docs/ARCHITECTURE.md`, `docs/REVIEW_AND_VERIFICATION.md` — the deliberation
  and evaluation surfaces these defects sit in

### Where the code is

Located by `grep -rln 'plan_uncritiqued\|execute_pending\|structured array'
--include=*.go .` — verified to reproduce at commit 1346d37. Re-run it; if the
file list differs, trust the grep over this table.

| Path                              | Holds                                          |
| --------------------------------- | ---------------------------------------------- |
| `internal/loop/apply.go`          | the envelope and its validation                |
| `internal/loop/gate_test.go`      | the `edit` / `done` gate preconditions          |
| `internal/loop/policy_test.go`    | `execute_pending` / `test_pending` transitions  |
| `internal/deliberation/select.go` | phase selection, `plan_uncritiqued`             |
| `internal/deliberation/types.go`  | the `ReasonCode` vocabulary                     |
| `internal/domain/regressions.go`  | `CreateReflection` — raises the D1 error (~733) |
| `internal/mcp/mcp_test.go`        | `TestGreenfield_MCPPlanBootstrap_EditGatePasses` — the canonical D3 pattern |
| `cmd/trace/loop_test.go`          | the CLI surface these defects surface through   |

---

## D1 — the reflection envelope rejects every key a caller can guess

**Symptom.** `loop apply` with `writes.reflections[]` fails:

```
domain: reflection requires at least one structured array
```

The message names no key. A caller must guess. One consumer session burned a
whole session on this and tried **25** different key names, all rejected
(`findings`, `weaknesses`, `issues`, `risks`, `questions`, `objections`,
`gaps`, `learnings`, `lessons`, `observations`, `improvements`, `next_actions`,
`follow_ups`, `actions`, `takeaways`, `changes`, `notes`, `evidence`,
`hypotheses`, `dimensions`, `affected_tests`, `attribution`, `comparison`,
`claims`, `assumptions`, `alternatives` — as bare string arrays and as
`{"title","detail"}` objects).

**Ground truth.** The error is raised in
`internal/domain/regressions.go` (`CreateReflection`, ~line 733):

```go
if len(assumpIDs) == 0 && len(deps) == 0 && len(tests) == 0 {
    return store.Reflection{}, &ErrValidation{Msg: "reflection requires at least one structured array"}
}
```

`assumpIDs`, `deps`, `tests` are built from
`InvalidatedAssumptionIDs`, `NewDependencies`, and `UsefulTests`. So the three
accepted keys are **`invalidated_assumptions`**, **`new_dependencies`**, and
**`useful_tests`** — and each is an array of **strings**, not objects. (The
`{kind, ref}` shape the caller eventually used is accepted at the envelope
layer and projected to these string fields; check the input type before
assuming.) A working item also needs `id` and `task_id`.

There is a real design choice here, and it is yours to make, not mine to
prescribe: an "essay-only" reflection (just `summary` / `broaden_tests_note`)
**deliberately fails closed** — see the comment on `CreateReflection`. Naming
the keys must not weaken that. A caller who supplies only prose still gets an
error; the error just tells them which key to reach for.

**Fix.** The error must name all three, and say they are arrays of strings.

```
domain: reflection requires at least one structured array — use
`invalidated_assumptions`, `new_dependencies`, or `useful_tests`
(a JSON array of objects with a `kind`/`ref` or string entries)
```

A caller who has never read the schema should get it right on the first try.
If the three keys accept different shapes, say which shape belongs to which.

---

## D2 — envelope validation reports one missing field at a time

**Symptom.** Reaching a valid `loop apply` envelope cost one consumer session
**five** round trips, because each error names only the first problem:

```
missing required field "apply_id"      -> then
apply_id must be UUID                  -> then
missing required field "seed"          -> then
writes.plan_changes[0].id is required
```

The minimal working envelope is four fields:

```json
{
  "schema_version": "trace.loop.apply.v1",
  "apply_id": "<uuid v4>",
  "seed": { "task_id": "<uuid>", "goal_id": "<uuid>" },
  "writes": { "plan_changes": [ { "id": "<uuid>", "title": "...", "body": "..." } ] }
}
```

**Fix, either is acceptable — pick one and justify it in the commit:**

(a) **Report every missing/invalid field at once.** One error, all problems.
    Cheapest, no new surface, and it is the fix the original session would
    have wanted.
(b) **Add a `--print-envelope` subcommand** (or an equivalent) that emits a
    valid skeleton for the given seed. Costs a little surface area and removes
    the guessing entirely.

Prefer (a) if you want the smallest diff. Both are legitimate; what is **not**
acceptable is leaving one-field-at-a-time.

---

## D3 — `plan_changes[]` opens the `edit` gate, a valid `reflection[]` does not

**Symptom.** Two gates whose preconditions are documented nowhere and behave
inconsistently:

- `loop gate --for edit` refuses with `reason_code: plan_uncritiqued`.
- Applying a **well-formed** `writes.reflections[]` does **not** clear it.
  Applying `writes.plan_changes[]` **does**.

One session applied both, in that order, and re-gated after each: the
reflection left `plan_uncritiqued` standing; the `plan_changes` row cleared it.

**It is documented — but too late to find.** `docs/rules/agent-loop-protocol.md`
§"Post-bootstrap critique path (Phase 37 R11)" (line ~249) does say to seed
critique via `trace loop apply` with a `plan_changes` envelope, and names
`TestGreenfield_MCPPlanBootstrap_EditGatePasses` in `internal/mcp/mcp_test.go`
as the canonical pattern. That test exists — read it before changing anything.

The defect is **discoverability, not absence**: the rule sits in a section
titled for post-bootstrap critique, ~250 lines below the session-start gate
that emits the error, and nothing in the error points at it. A caller who hits
`plan_uncritiqued` from the gate has no way to know a documented answer exists.
So the fix is to surface the remedy at the point of failure, not to write new
prose where the old prose already is.

`trace_add kind=plan-change` still does **not** set the flag — only
`loop apply` with `writes.plan_changes[]` does. That part is worth stating
somewhere reachable.

**Also, and separately: `done` is gated by `execute_pending`.** The natural
closing sequence — apply your closing reflection, then gate — can walk into a
block, because an apply is itself an implementation signal and re-arms
`execute_pending`. `internal/loop/policy_test.go:81` states the intent plainly:
*"implementation signal must clear execute_pending"*, with the paired assertion
that a **recorded git change** moves the loop to `test_pending`. So the
transition the loop wants after applying work is a *change* (via
`trace verify run` / change recording), not another apply.

The observed consequence for a caller: applying an empty `writes: {}`
envelope returns `"saturated": true` and the `done` gate then opens. Work out
whether that is the intended path or an accident of ordering before you rely on
it — and if it is an accident, say so in your fix rather than codifying it.

**Fix.** Make both preconditions discoverable from where a caller already
looks:

- `loop gate` already returns a good `reason_code`. **Add the remedy** — the
  concrete call that would clear it, ideally pointing at
  `TestGreenfield_MCPPlanBootstrap_EditGatePasses`. A caller reading
  `{"allowed": false, "reason_code": "plan_uncritiqued"}` should not have to
  search the source to learn that `writes.plan_changes[]` is the answer.
- Keep the §"Post-bootstrap critique path" section, and add a pointer to it
  from the session-start gate section — the section a reader reaches *first*.
- Consider a `final: true` flag on the envelope meaning "this is the last
  apply", which would let `done` open without a second, empty call. Judgement
  call — if you skip it, say why in the commit.

**Do not** "fix" this by setting `plan_critiqued` from `trace_add
kind=plan-change`. A prior session concluded that was the fix, shipped a
diagnosis on it, and was wrong. If you believe the current design is
genuinely wrong, that is a human decision with a written rationale.

---

## D4 — `trace_link rel=discovery-plan-change` cannot target a plan scope

**Symptom.** The rel is accepted, but resolving `to` against a **`plan_scopes`**
id — the id `trace plan deep` returns and `set-current` takes — fails with a
bare:

```
store: plan_change %q: sql: no rows in result set
```

It names neither table, so the caller cannot tell which id space is wrong. The
`to` must be a `plan_changes` **row** id, which the MCP surface gives no way
to obtain: a synthesised UUID is accepted as a `writes.plan_changes[0].id`, but
nothing hands that id back for linking.

**Net effect: the documented way to record a critique against a plan scope is
unreachable through the tools.**

**Fix.** Either:

(a) accept a `plan_scopes` id on this rel and resolve it to the scope's
    current plan-change, or
(b) return the created `plan_changes` id from the apply that creates it, so
    the caller has something to link.

(a) is the better fix. Whatever you choose, **name the table in the error** —
`plan_changes` vs `plan_scopes` — and say which id space was expected.

---

## Method: reproduce before you fix

For each defect, in a scratch store — **never** the consumer's real graph:

1. `cp /home/ali/Desktop/mostashar/.trace/trace.db /tmp/trace-probe.db`
   (read the copy only; never write to the original — it holds real case
   memory)
2. Point the probe at the copy and run the failing command
3. Capture the exact failing output
4. Write a Go test that **fails** on current `main`
5. Fix, re-run, confirm the test passes
6. Capture the passing output with the same command

**Never hand-edit a `.trace/trace.db`.** Copy it to `/tmp` to experiment.

**Find what *writes* a flag; do not infer it from its name.** A prior session
read `plan_critiqued`, assumed it meant "a critique was recorded", built a
causal chain on that, reported the phase as blocked, and shipped a wrong
diagnosis that a later session had to formally supersede. Read the code that
assigns the field.

## Hard rules

1. **Never widen a gate to get green.** A gate that opens when it should not is
   worse than one that stays shut. If a precondition is wrong, fix the
   *predicate* or the *message* — never the threshold guarding it. The
   `as_operator` requirement for DONE, and "evidence_ids alone do not authorize
   DONE", are load-bearing and stay exactly as they are.
2. **No behaviour change without a test that failed first.**
3. **For a gate fix, the test is the negative case**: prove the gate still
   refuses what it exists to refuse. A fix verified only by a positive test is
   a new defect. Re-run the sequences from D3 and confirm each still blocks.
4. **No behaviour change without operator sign-off.** These prompts touch
   `as_operator` and the DONE gate — surface the plan to the human before
   implementing, per `docs/rules/agent-loop-protocol.md`.
5. **Work on a branch.** Do not push to `main` without explicit approval.
6. **A defect you cannot fix is recorded, not worked around** — write it up
   with the exact failing output and a reason, in the commit message or a
   `docs/TODO/` row. A silent workaround is the failure mode here.
7. **`gofmt` is a fail-fast CI gate.** Run it before you commit.

## Verify

```bash
go test ./... 2>&1 | tail -20
gofmt -l .            # must print nothing
```

Then, per defect, the before/after pair from **Method**. And re-run the four
gate sequences end to end against the fixed build, confirming each still
blocks what it blocked.

## Out of scope

A **fifth** finding exists but is **not Trace's**. The consumer's own phase
runbook contained a `git checkout -- .` that reverted the file under test, so
a "green" result came from unmodified code and would have closed a finding
having shipped nothing. It is being fixed in that repo with a lint over its
prompt docs. **Do not touch it from here.** It is mentioned only because it is
the most serious of the five, and because it is the clearest illustration of
why a green that nobody can explain is not evidence.

## Report back

Per defect: the failing output before, the passing output after, the same
command. The commit sha, on a branch. Anything you declined to fix, and why.
If a fix changes gate behaviour, state which previously-recorded gate verdicts
were produced under the old behaviour and whether they still hold.
```
