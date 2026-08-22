# Master TODO — Trace

**Source of truth for run order.** Phase row tables live in [`docs/TODO/`](TODO/). Run top-to-bottom. One fresh subagent per row. Re-read the **active phase board** after every subagent (spawns may appear).

Protocol: [`docs/rules/agent-loop-protocol.md`](rules/agent-loop-protocol.md)  
Laws: [`docs/rules/project-rules.md`](rules/project-rules.md)  
Planning registers: [`docs/init/README.md`](init/README.md) (design SoT; **not** the run board)

Legend: `pending` | `in_progress` | `done` | `failed` | `blocked` | `skipped`

Orchestrator paste:

```text
Phase 00–42 complete — do not re-run closed rows.
@docs/TODO.md

- Active phase: **none** — remediation G1–G9 complete
- Next runnable: **idle** (optional Phase 43+ residuals — human promotion only)
- Follow docs/rules/agent-loop-protocol.md
```

(Phase 00–42 complete. **P42-S02-02 done** 2026-08-22 — G6+G7 delivered; successor **no successor**.)
---

## Phase boards

| Phase | Title | Status | Next | Board |
|------:|-------|--------|------|-------|
| 00 | Foundation (P0-X) | done | — | [TODO/phase-00.md](TODO/phase-00.md) |
| 01 | X0 readiness | done | — | [TODO/phase-01.md](TODO/phase-01.md) |
| 02 | Gate C evaluation & slice hardening | done | — | [TODO/phase-02.md](TODO/phase-02.md) |
| 03 | Progressive planner (minimal) | done | — | [TODO/phase-03.md](TODO/phase-03.md) |
| 04 | Review depth & evidence policies | done | — | [TODO/phase-04.md](TODO/phase-04.md) |
| 05 | Decision impact & simulation | done | — | [TODO/phase-05.md](TODO/phase-05.md) |
| 06 | Environment / capability graph | done | — | [TODO/phase-06.md](TODO/phase-06.md) |
| 07 | Performance ladder & language plugins | done | — | [TODO/phase-07.md](TODO/phase-07.md) |
| 08 | Ecosystem & hardening | done | — | [TODO/phase-08.md](TODO/phase-08.md) |
| 09 | Dogfood hardening & agent UX | done | — | [TODO/phase-09.md](TODO/phase-09.md) |
| 10 | Integrity surfaces (post-dogfood) | done | — | [TODO/phase-10.md](TODO/phase-10.md) |
| 11 | Residual surfaces (post–P10 open findings) | done | — | [TODO/phase-11.md](TODO/phase-11.md) |
| 12 | Peer honesty surfaces (thin: edge provenance + packet honesty) | done | — | [TODO/phase-12.md](TODO/phase-12.md) |
| 13 | Import resolve & honesty residuals (thin: DF-60…67) | done | — | [TODO/phase-13.md](TODO/phase-13.md) |
| 14 | Peer impact + install gates (thin: research ranks 4–6) | done | — | [TODO/phase-14.md](TODO/phase-14.md) |
| 15 | P14 residual remediation plan (thin) | done | — | [TODO/phase-15.md](TODO/phase-15.md) |
| 16 | Assert root & surfaces (thin: post-P15 open DFs) | done | — | [TODO/phase-16.md](TODO/phase-16.md) |
| 17 | Portable graph via git (seed export + plan tree + convention + idempotent import) | done | — | [TODO/phase-17.md](TODO/phase-17.md) |
| 18 | Context FTS + clone honesty (thin: D40 residuals) | done | — | [TODO/phase-18.md](TODO/phase-18.md) |
| 19 | Loop gap detection | done | — | [TODO/phase-19.md](TODO/phase-19.md) |
| 20 | Cognitive deliberation | done | — | [TODO/phase-20.md](TODO/phase-20.md) |
| 21 | TRACE thoughtprocess completion | done | — | [TODO/phase-21.md](TODO/phase-21.md) |
| 22 | Capability completion | done | — | [TODO/phase-22.md](TODO/phase-22.md) |
| 23 | Agent enforcement choke points + harness install | done | — | [TODO/phase-23.md](TODO/phase-23.md) |
| 24 | Agent effectiveness investigation | done | — | [TODO/phase-24.md](TODO/phase-24.md) |
| 25 | Orchestrator + default gap pass (P25-C) | done | — | [TODO/phase-25.md](TODO/phase-25.md) |
| 26 | Loop investigations, planning & implementation (P25-A + P25-B + installer fix) | done | — | [TODO/phase-26.md](TODO/phase-26.md) |
| 27 | Protocol measurement + graph honesty (P25-D + P25-E) | done | — | [TODO/phase-27.md](TODO/phase-27.md) |
| 28 | Residuals closure + implementation validation | done | — | [TODO/phase-28.md](TODO/phase-28.md) |
| 29 | HTTP API + browser GUI (local-first → cloud path) | done | — | [TODO/phase-29.md](TODO/phase-29.md) |
| 30 | Stray root `trace.db` hygiene | done | — | [TODO/phase-30.md](TODO/phase-30.md) |
| 31 | Stray `trace.db` residual testing | done | — | [TODO/phase-31.md](TODO/phase-31.md) |
| 32 | Graph-first GUI (explorer ambition) | done | — | [TODO/phase-32.md](TODO/phase-32.md) |
| 33 | GUI craft + Explore hook + `trace gui` | done | — | [TODO/phase-33.md](TODO/phase-33.md) |
| 34 | GUI packaging + embed SPA + auto ports | done | — | [TODO/phase-34.md](TODO/phase-34.md) |
| 35 | Active task selection (test + fix) | done | — | [TODO/phase-35.md](TODO/phase-35.md) |
| 36 | Planning model alignment + plan_missing root cause | done | — | [TODO/phase-36.md](TODO/phase-36.md) |
| 37 | Phase 36 residuals closure | done | — | [TODO/phase-37.md](TODO/phase-37.md) |
| 38 | Retrieval & context peer-gap investigation | done | — | [TODO/phase-38.md](TODO/phase-38.md) |
| 39 | Context orient & harness | done | — | [TODO/phase-39.md](TODO/phase-39.md) |
| 40+ | Read surface & retrieval depth | done | — | [TODO/phase-40.md](TODO/phase-40.md) |
| 41+ | Layers & intent | done | — | [TODO/phase-41.md](TODO/phase-41.md) |
| 42+ | Concept & index | done | — | [TODO/phase-42.md](TODO/phase-42.md) |

(P42 complete 2026-08-22 — G6+G7 delivered; DR-HANDOFF CLOSED; successor **no successor**; G1–G9 remediation complete.)

---

## Forward residuals (not a board phase)

| Queue | Status | Board |
|-------|--------|-------|
| Phase 28 R6 / FM gaps | **superseded / closed** | Closed on Phase 28 S06/S07 — [`TODO/phase-28.md`](TODO/phase-28.md); legacy index [`TODO/forward-p28-residuals.md`](TODO/forward-p28-residuals.md) |
| Phase 20 residuals | superseded | Moved into **Phase 21** — [`TODO/phase-21.md`](TODO/phase-21.md); legacy [`forward-p20-residuals.md`](TODO/forward-p20-residuals.md) |

Do not re-run Phase 28 `done` rows. Residual wave closed at `P28-S07-02` with successor **no successor**.

---

## Later developments (not a board phase)

**Not runnable. Not ordered into Phase 16–22.** Phase 19–22 closed historically with `no successor`. Do not spawn hosted MCP from this section. If git-only sharing after P17 dogfood is enough, **do not build a server**.

Optional **hosted Trace** (name TBD): a **separate repo**, not `github.com/mrchatam/Trace`. Trace core stays local-first CLI + library + seed JSON (`trace/graph.json`). A future hosted product would be opt-in plus — secure server, user auth, MCP calling **cloud** instead of local stdio.

| Lock | Value |
|------|-------|
| Contract | **P17 seed/plan JSON** (same import/export). Do **not** copy `internal/store` into the hosted repo as a second SoT |
| This repo’s `trace-mcp` | Local stdio only. **Do not** point it at the internet |
| Brand | Trace = local-first. Hosted is a distinct opt-in product (name TBD) |
| Sequence | Only after P17 dogfood of **two humans on git-only sharing**. If that works, stop |
| Security debts that make “secure server” real | Tenancy; human vs agent credentials; audit; redaction; org vs laptop capabilities; `exported_at_commit` as evidence not identity; privacy/retention |
| Must not | Steal or reorder current board rows; point local `trace-mcp` at the public internet; treat `transition.actor` as auth; ship multi-tenant SaaS inside this phase |
| Phase 29 carve-out | **Opt-in local HTTP** (`trace serve`) + browser GUI allowed; default bind loopback; OpenAPI is the contract for a **future hosted** product (separate deploy) |

DF-86 (`trace install git-hook`) is **promoted in Phase 22 S02** (still local git, still not a hosted service, still must not wrap `git commit`). P17/P18 VERIFY treated absence as non-fail historically.

---

## Historical mapping

Former flat board [`docs/init/B_INITIAL_BOARD.md`](init/B_INITIAL_BOARD.md) (`T001`–`T014`) maps into Phase 00 scopes in [`TODO/phase-00.md`](TODO/phase-00.md). **Do not run T-ids**; run the linked phase boards via this index.
