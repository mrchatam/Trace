# P11 / S04 / 01 — Capability upsert + hatch vs caps

## Metadata
- id: P11-S04-01
- todo_ids: [P11-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-41, DF-51** per sibling **00-PLANNER** FINAL locks (2026-08-16). Empty-ID capability declare upserts by slug (stable id). `--allow-done` stays independent of missing-caps; thicken hatch WARNING + docs. **No new migration. Gate G + DF-24 fail-closed retained.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-41, DF-51
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md)
- [experiments/_post_p10/bughunt/adv_hatch_caps/](../../../../../experiments/_post_p10/bughunt/adv_hatch_caps/)
- [phase README](../../README.md)
- Live: `internal/domain/{capability,task_state}.go`; `internal/store/capability.go`; `cmd/trace/{capability,transition,help}.go`; `internal/mcp/{tools_parity,tools_write,server}.go`
- Prior: P10 S04 DF-24/26; P11-S02 hatch↔caps independence; P11-S03 lock (no coupling)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. If 00-PLANNER is still DRAFT, stop and return to planner.

## Locked defaults (FINAL — P11-S04-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/domain` (+ thin `cmd/trace` + `internal/mcp` copy only) |
| Migration | **None** |
| DF-41 | Empty ID → resolve by slug → reuse id + update; explicit different-id slug clash still fails |
| DF-51 | Hatch does **not** bypass missing caps; WARNING + help/MCP mention both review hatch and `--allow-missing-caps` |
| Check order | Keep: actor+reason → legal edge → DF-24 caps → →DONE hatch/PASS/operator |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false`; P10 DF-17/18/24/26; P11 DF-40/43/44/47 |
| Forbidden | Mig; auto-`AllowMissingCapabilities` from hatch; daemon/HTTP; board spawn; rewrite `done` history |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/domain/capability.go` | Empty-ID: `GetCapabilityBySlug` → set `ID` before store upsert; preserve validation |
| `internal/domain/capability_test.go` | `TestUpsertCapabilityBySlugUpdatesExisting` (+ keep different-id duplicate fail) |
| `internal/domain/operator_gates_test.go` (or sibling) | `TestAllowDoneDoesNotBypassMissingCaps`; hatch+override OK |
| `cmd/trace/transition.go` | Extend loud `--allow-done` WARNING with missing-caps phrase |
| `cmd/trace/help.go` (+ usage note) | Document allow-done ≠ missing-caps bypass |
| `cmd/trace/cli_test.go` | Update `TestAllowDoneWarnsOnStderr` assert |
| `internal/mcp/tools_write.go` (+ schema/desc) | Extend MCP `warning` + `allow_done` honesty copy |
| `internal/mcp/mcp_test.go` | Update `TestTransitionAllowDoneEmitsWarning` |
| CLI capability (optional) | Usage note omit `--id` updates by slug — no logic fork |

## Role work

1. TDD domain: slug re-declare without id → same id + updated status; different-id clash still fails.
2. Wire empty-ID slug resolve in `UpsertCapability` (domain only; G19).
3. TDD: allow-done alone + missing caps → reject; allow-done + allow-missing-caps → OK.
4. Thicken CLI stderr WARNING + MCP `warning` + help/schema with missing-caps independence phrases (keep DF-26 loud hatch).
5. Run locked verify suite; board **status + Notes only** (cite test names + DF-41/51).

## Algorithm sketch (non-normative — locks above win)

```text
# DF-41
UpsertCapability(in):
  normalize kind/slug/status
  if in.ID == "" {
    if existing := GetCapabilityBySlug(slug); found {
      in.ID = existing.ID
    }
  }
  return store.UpsertCapability(...)  # still ON CONFLICT(id)

# DF-51 (policy unchanged; messaging only)
TransitionTask: caps gate before DONE hatch (unchanged)
if success && AllowDoneWithoutReview:
  WARNING mentions PASS/as_operator bypass AND missing-caps still need allow-missing-caps
```

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Carry-forward must stay green. Gate C `dry_run:false` untouched. Pre-existing FAIL under `similar projects/` space-path (if any) is residual OK — not S04 scope.

## Exit criteria
- [ ] DF-41: slug re-declare without id updates same row; different-id clash still fails
- [ ] DF-51: hatch alone does not bypass missing caps; WARNING/docs mention missing-caps override
- [ ] Locked tests green; carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P11-S04-02** (cite tests + DF-41/51)

## Out of scope
- S05+ Phase 11 product work (DF-49, 35, 48, 42, …)
- Making `--allow-done` imply `--allow-missing-caps`
- Daemon/HTTP/embeddings; full-rebuild indexer; rewriting `done` history
