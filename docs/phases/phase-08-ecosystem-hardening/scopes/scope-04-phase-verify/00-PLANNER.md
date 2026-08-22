# P08 / S04 / 00-PLANNER — Phase 08 VERIFY / compat+security checklist

## Metadata
- id: P08-S04-00
- todo_ids: [P08-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Phase 08 VERIFY commands, evidence table, **compatibility + security checklist** harness path, spawn 01a/b/c on fail, and DR-HANDOFF = **`no successor`** (`A_PROJECT_PLAN` ends at Phase 8) unless Notes promote follow-on. No product Go.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 8
- Pattern: Phase 07 VERIFY [`../../../phase-07-performance-ladder/scopes/scope-03-phase-verify/`](../../../phase-07-performance-ladder/scopes/scope-03-phase-verify/)
- [docs/TODO.md](../../../../TODO.md)
- S01–S03 REVIEW-NOTES under sibling scopes

## Session start
Agent → clarify → Plan → execute (planner).

## Phase defaults already locked (respect — P08-00)
| Item | Value |
|------|-------|
| Validation gate | Compatibility + security checklist |
| Preferred harness | **`evals/compat`** planted checks |
| Preferred names | `TestCompatibilitySecurityChecklist` + `schema-compat.json` v1 + temp `metrics-compat.json` (**finalize here**) |
| S01–S03 | Re-prove plugin API + worktrees + production hardening surfaces |
| Carry-forward | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Successor | **`no successor`** unless Notes promote |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / this checklist |

## Depends note from S03 (updated 2026-08-16 — S03-02 APPROVE)
S03 **shipped + reviewed**: re-prove against live — **`trace migrate status`** (`MigrationStatus` max/embed/pending); **`trace backup`/`restore`** (`VACUUM INTO` `trace.db` snapshot + Abs `root_path` rebind; `HasBlobLikeColumns` false; token excluded by default); optional **`.trace/access.token`** + **`TRACE_ACCESS_TOKEN`** fail-closed on `store.Open` (`ErrUnauthorized`); **no** `011_*`; CLI-primary (**no** new MCP tools; Open inherits gate). Checklist must include these fail-closed proofs plus S02 `trace.lock` + S01 `LanguageAdapterAPIVersion=1`. See [../scope-03-production-hardening/REVIEW-NOTES.md](../scope-03-production-hardening/REVIEW-NOTES.md).

## Planner work
1. Lock checklist harness path + must-pass items (compat API version; no daemon/HTTP primary; migrate/backup/auth fail-closed; G19; no source BLOBs; S02 lock; S01 adapter version).
2. Thicken `01-verify.md` evidence table + spawn convention + no-successor checklist.
3. Thicken `02-scope-review.md` owns DR-HANDOFF completion (`no successor`).
4. SCOPE-TODOS sync.

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Checklist path | **`evals/compat`** / **`TestCompatibilitySecurityChecklist`** |
| Schema / metrics | **`schema-compat.json` v1** + temp **`metrics-compat.json`** (`dry_run:false`) |
| Harness ownership | **S04-01 creates** `evals/compat` as VERIFY work (S01–S03 left none — confirmed absent at plan time). Do **not** block checklist awaiting a prior seed. |
| Allowed Go on VERIFY | **`evals/compat/**` only** (doc.go, schema, planted fixtures, named test). No opportunistic `cmd/` / `internal/` product rewrites unless spawn remediation requires them. |
| Must-pass (S01) | `LanguageAdapterAPIVersion=1`; `TestLanguageAdapterAPIVersion` + `TestBuiltinLanguageAdaptersContributionPath`; static adapters only (no `.so` / megastore) |
| Must-pass (S02) | Path-local `<abs>/.trace/` (no walk-up); exclusive `trace.lock` / `ErrLocked`; `TestProjectBindPathLocalIsolation` + `TestConcurrentStoreOpenFailClosed` (+ CLI lock fail-closed) |
| Must-pass (S03) | `MigrationStatus` / `trace migrate status` (embed max = applied; **no** `011_*`); backup↔restore (`VACUUM INTO` `trace.db`; Abs `root_path` rebind; `HasBlobLikeColumns` false; token excluded by default); local auth fail-closed (`.trace/access.token` + `TRACE_ACCESS_TOKEN` → `ErrUnauthorized`); lock held on backup/restore; **no** new MCP auth/backup tools |
| Laws (checklist) | No daemon / always-on HTTP as primary; G19 (libraries do not import `cmd/trace` / `cmd/trace-mcp`); no source BLOBs; no full-rebuild-on-any-change indexer; CLI-primary for migrate/backup/auth |
| Carry-forward | Gate H (`TestPlantedPerfLadderGateH`); honesty A/B/C + Gate G; Gate E; Gate F; ablation; p0x 7/7; x0; Gate C artifacts `dry_run:false` N=3 |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / **this checklist** |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | Record **`no successor`** (`A_PROJECT_PLAN` ends at Phase 8). Do **not** scaffold Phase 09 unless user Notes explicitly promote. S04-01 starts handoff Notes; **S04-02 owns completion**. |

### Metrics schema shape (v1 — lock fields)

`schema-compat.json` **required** keys (`additionalProperties: true` OK):

| Field | Notes |
|-------|-------|
| `schema_version` | const `1` |
| `gate` | const `"compat-security"` |
| `suite` | const `"compat"` |
| `dry_run` | boolean — must be **`false`** for checklist pass artifact |
| `named_test` | `"TestCompatibilitySecurityChecklist"` |
| `language_adapter_api_version` | integer — must be `1` |
| `language_adapter_api_version_ok` | boolean |
| `path_local_bind_ok` | boolean |
| `trace_lock_ok` | boolean |
| `migrate_status_ok` | boolean |
| `backup_restore_ok` | boolean |
| `no_blob_columns_ok` | boolean |
| `local_auth_fail_closed_ok` | boolean |
| `g19_ok` | boolean |
| `no_daemon_http_primary_ok` | boolean |
| `no_011_mig_ok` | boolean |
| `s01_hooks` / `s02_hooks` / `s03_hooks` | string arrays naming consumed surfaces / tests |

Pass bar: named test PASS + schema-valid temp metrics with all `*_ok` true + `language_adapter_api_version == 1` + `dry_run:false`.

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] Checklist names + no-successor locked
- [x] Board Notes; next `P08-S04-01`

## Out of scope
- Running VERIFY (S04-01)
- Inventing commercial multi-model security theater
- Scaffolding Phase 09 without explicit promotion
