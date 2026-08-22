# S02 — Install / capability gates — scope todos

**Depends-on:** P14-S01-02 APPROVE (serial). Owns research **ranks 4–5**.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** (done) |
| 2 | 01-install-capability-gates | implement | **done** |
| 3 | 02-scope-review | review | **APPROVE** (done) |

## FINAL homes (from 00-PLANNER)
- **Install:** `internal/install` + thin `trace install` (detect / `<target>` / uninstall); tiers `STABLE`\|`CONDITIONAL`\|`OPT_IN`; Cursor STABLE + ≥1 CONDITIONAL marker-gated target
- **Audit:** mig **013** + domain `ResolveToolDecision` / `AssertToolAllowed` / decide CLI; statuses `AUTO_ALLOWED`\|`PENDING`\|`ALLOWED`\|`DENIED`; auto-allow = exact builtin MCP slugs only
- **No** new MCP tools; **no** ImpactWalk edits; **no** YOLO/AllowAll

## Depends (from S01 — light)
S01 **APPROVE** (P14-S01-02): `trace impact walk` + named ImpactWalk tests + Gate F / `finding`/`alternative`/`report` must stay green. S02 must **not** regress that CLI surface or S01 named tests; VERIFY (S03) will re-run ImpactWalk + Gate F as carry-forward alongside install/capability bars. Residual (non-blocking): late `allowContainsOut` upgrade — optional S03 spot-check only; do not fix in S02.

## Depends (to S03 — light)
After S02 APPROVE: S03 VERIFY imports S02 named install + decision tests, `TestInstallCursor*`, capability ablation, **and** S01 ImpactWalk + Gate F. Optional S01 residual spot-check only.

## Reminders
- No YOLO / AllowAll; CLI/rules-first; no new MCP install/decide tools
- Do not regress S01 impact walks or Phase 09/11 install Cursor
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P14-S03-00**
