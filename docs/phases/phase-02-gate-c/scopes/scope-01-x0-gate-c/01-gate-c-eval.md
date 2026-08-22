# P02 / S01 / 01 — X0 Gate C evaluation

## Metadata
- id: P02-S01-01
- todo_ids: [P02-S01-01]
- role: implementer
- skills: [incremental-implementation, tdd, source-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: mixed

## Objective
Extend **`evals/x0`** beyond Phase 01 dry-run: run Experiment X0 for real (agent **B0** vs **G1**), score **understanding** against human GT on `fixtures/x0`, emit schema-valid metrics with `dry_run:false` and **N≥3** runs per condition, and draft Gate C **Go / No-Go / iterate** evidence. Keep `TestX0DryRunMetricsB0AndG1`, `evals/p0x`, and `evals/honesty` green. CLI path only (DR-AGENT); MCP not required.

**Phase 01 dry-run ≠ Gate C pass.** Do not claim product-thesis success from `dry_run:true` alone.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment X0
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 2 Gate C
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-AGENT, DR-SEED, DR-BENCH
- Live: `evals/x0` (`schema.json` v1, `TestX0DryRunMetricsB0AndG1`), `fixtures/x0`, `evals/p0x` criterion-6 query themes

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Conditions | **B0** = agent + ordinary repo tools (read / git / search) — **must not** call `trace why` / `trace context`. **G1** = agent + `trace` CLI `why`/`context` (may still read repo) |
| MCP | **Not** required (DR-AGENT). Do not gate Gate C on MCP |
| Corpus | Synthetic **`fixtures/x0`** only; abs path to `seed/gt.json` v1; **no** new OSS this scope |
| Stable IDs | Goal `11111111-…111`; Task `22222222-…222`; Decision `33333333-…333`; Discovery `44444444-…444`; PlanChange `55555555-…555` |
| Instrument | Prefer **extend** `evals/x0`; keep dry-run regression green |
| Schema | Prefer keep `schema_version` **1** + document `quality` object shape below. Bump to **2** only if a required field must change (then update dry-run + Gate C validators) |
| N-runs | **≥3** independent agent runs **per condition** (B0 and G1 each) |
| Task family | Primary: **`understanding`**. Optional secondary runs (`implementation`, `honesty`) may appear in metrics but **do not** drive Gate C pass/fail. H5 stay in `evals/honesty` |
| Fixture prep | Same as dry-run: temp-copy `fixtures/x0` → `init` → `seed import <ABS>` → `index` (`CGO_ENABLED=1` build of `./cmd/trace`). Separate work dirs per condition (and prefer per-run isolation) |
| Seeding cost | Human GT seed in `fixtures/x0` is **non-trivial** (multi-entity causal graph + links). Second kill conjunct is **satisfied** for this corpus — do not invent a “trivial seed” escape hatch |
| Evidence report | **`GATE-C-NOTES.md`** in this scope folder (committed) |
| Metrics persistence | Copy schema-valid Gate C metrics into **`docs/verification/gate-c-x0/`** (`metrics-b0.json`, `metrics-g1.json`, optional `pins.md`). Harness may write temp first; **do not** commit `.trace/` under fixtures/evals |
| Kill criteria | If **mean G1 understanding_accuracy ≤ mean B0 understanding_accuracy** (see scoring) **and** seeding non-trivial → **thesis endangered** → decision must be **No-Go** or **iterate** (not silent Go) |
| “Within error” | For N=3 on this fixture: treat as **G1 mean ≤ B0 mean** on `understanding_accuracy`. No fancy significance test required to apply kill; document per-run variance in GATE-C-NOTES |
| Keep green | `TestX0DryRunMetricsB0AndG1`; `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/... ./evals/x0/...` + `./...` |
| G19 | Harness shells out to built `trace`; packages must **not** import `cmd/trace` |
| Out | Slice hardening product fixes (S02); daemon/HTTP/embeddings; declaring A1 VALIDATED without evidence; replacing `evals/p0x` with agent scores; MCP-required path |

### Understanding query bank (locked)

Commit a machine-readable query bank under `evals/x0/` (suggested: `queries.json` or `testdata/gate-c/queries.json`) with **≥5** agent-facing understanding questions. Both conditions get the **same** questions in the **same** order.

Questions must be answerable from fixture GT / seeded graph + repo, covering at least:

1. Task↔goal linkage (IDs or titles matching GT)
2. Decision that constrains the task (Decision `3333…` / “Prefer TypeScript greeter”)
3. Discovery→plan_change chain (`4444…` → `5555…`)
4. Why the task is in progress / causal neighbor presence
5. Bounded context expectation (what belongs in task context vs whole-repo dump) — score against GT entities, not against raw item-count alone

Reuse themes from `evals/p0x` criterion-6 where helpful, but phrase as **agent Q&A** with explicit GT answer keys (IDs, relation names, short expected facts). Do **not** score B0 by secretly calling `why`/`context`.

### Scoring (locked — fairness)

| Metric | Definition |
|--------|------------|
| Per-query grade | `correct` \| `incorrect` \| `critical_miss` |
| `critical_miss` | Asserts a **wrong** causal link / entity role that contradicts GT (e.g. invents decision→task, swaps discovery/plan_change) |
| `understanding_accuracy` | `correct / total_queries` for that run (critical_miss counts as not correct) |
| Primary compare | Mean `understanding_accuracy` across ≥3 runs per condition |
| Secondary | Mean `critical_misses`; optional tokens / latency / `ok` (task success) |
| Parity | Same questions, same fixture SHA, same model pin across B0/G1 unless Notes justify a controlled difference |

**Forbidden:** Contaminating B0 with Trace CLI; giving G1 different/easier questions; hand-waving “G1 felt better”; using dry-run stub `ok:true` as Gate C evidence.

### Quality object shape (schema v1 `quality` field)

Prefer each understanding run’s `quality` to be a JSON object:

```json
{
  "understanding_accuracy": 0.8,
  "correct": 4,
  "total": 5,
  "critical_misses": 0,
  "per_query": [{"id": "q1", "grade": "correct"}]
}
```

`additionalProperties` already allowed on runs — keep validators accepting this without a forced schema bump unless you change required top-level fields.

### Agent run protocol (locked)

1. **Pins** (record in metrics + GATE-C-NOTES): model name/version, fixture path + content hash or git tree SHA for `fixtures/x0`, `trace version`, date, operator/session id, RNG/seed if any.
2. **Answer packs**: For each `(condition, run_id)` produce an answer JSON (agent free-text or structured answers keyed by query id). Store under `evals/x0/testdata/gate-c/` **or** load via env path for live runs.
3. **Runner modes** (either OK):
   - **A — Live:** env such as `TRACE_X0_GATE_C=1` runs/invokes an external agent (document command); still scores via harness.
   - **B — Recorded:** committed answer packs from documented operator/agent sessions (same pins). Acceptable for Gate C if honest and N≥3.
4. **Harness duties:** load queries + answers → grade vs GT → write `metrics-b0.json` / `metrics-g1.json` with `dry_run: false`, `runs.length ≥ 3`, tools_used rules (B0 excludes why/context; G1 includes them when used) → schema-validate → copy to `docs/verification/gate-c-x0/`.
5. **Automated test:** Prefer a test that grades **fixture answer packs** (deterministic CI) plus documents how live runs refresh packs. Do **not** require network LLM in default `go test ./evals/x0/...`.

### Gate C decision artifact (`GATE-C-NOTES.md`)

Must include:

| Section | Content |
|---------|---------|
| Verdict | **Go** \| **No-Go** \| **Iterate** (one) |
| Evidence table | Per condition: N, mean accuracy, critical_misses mean, latency/tokens if present, model pin |
| Kill criteria check | Explicit yes/no on (G1≤B0) and (seeding non-trivial) |
| Issue list | Concrete, measurement-driven gaps for S02 (may be empty with reason) |
| Honesty | Statement that Phase 01 dry-run was **not** treated as pass |
| Residuals | Variance, sample-size limits, any unfairness risks |

**Verdict rules (locked):**

- **No-Go** if kill criteria fire (G1≤B0 and non-trivial seeding), unless Notes record an explicit human override (out of band).
- **Go** only if mean G1 understanding_accuracy **>** mean B0 **and** evidence table complete.
- **Iterate** if results mixed / harness gaps / need more N or query fixes before a clean Go/No-Go — still produce issue list.

### Target tree (expected)

```text
evals/x0/
  doc.go                 # update: Gate C + dry-run
  schema.json            # v1 unless bump justified
  x0_test.go             # keep TestX0DryRunMetricsB0AndG1; add Gate C scoring/aggregation tests
  queries.json           # or testdata/gate-c/queries.json
  testdata/gate-c/       # answer packs + optional helpers
docs/verification/gate-c-x0/
  metrics-b0.json
  metrics-g1.json
  pins.md                # optional but recommended
docs/phases/.../scope-01-x0-gate-c/
  GATE-C-NOTES.md        # decision + evidence + issue list
```

### How to run (Notes must cite)

```bash
CGO_ENABLED=1 go test ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/... ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts. Issue list lives in `GATE-C-NOTES.md` for S02 (reviewer/planner consume it).

## Exit criteria
- [ ] Query bank ≥5 + GT keys committed; same questions for B0/G1
- [ ] Scoring harness grades answer packs; `dry_run:false` metrics with **≥3** runs/condition; schema-valid
- [ ] `docs/verification/gate-c-x0/` metrics (+ pins) present
- [ ] `GATE-C-NOTES.md` with verdict + evidence table + kill check + S02 issue-list shape
- [ ] Dry-run regression + p0x + honesty + `./...` PASS
- [ ] Board Notes updated (status + notes only); no Gate C pass claimed from Phase 01 alone

## Minimal todos
- [ ] Add `evals/x0` query bank + GT answer keys (≥5)
- [ ] Implement grade + aggregate helpers; emit Gate C metrics (`dry_run:false`, N≥3)
- [ ] Provide ≥3 answer packs/condition (live or recorded) with pins
- [ ] Write `docs/verification/gate-c-x0/` artifacts + `GATE-C-NOTES.md`
- [ ] Keep `TestX0DryRunMetricsB0AndG1` green; run p0x + honesty + `./...`
- [ ] Board Notes only
