# Scope 00 — board map

**S00 research** — embed + StaticDir + auto-port only. Serial: **P34-S00-00 → P34-S00-01 → P34-S00-02**. Primary artifact: `RESEARCH.md` (written in **S00-01**, reviewed in **S00-02**). Do **not** start S01 until S00-02 PASS. Do **not** write product code. Planner (**S00-00**) does **not** author `RESEARCH.md`.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 593 | P34-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 + this file |
| 594 | P34-S00-01 | [01-research.md](01-research.md) | Implementer | Author `RESEARCH.md` |
| 595 | P34-S00-02 | [02-review.md](02-review.md) | Reviewer | Checklist vs DESIGN-LOCKS + INTAKE |

## Planner-locked live facts (for 01) — verified 2026-08-21

| Area | Fact | Path |
|------|------|------|
| Static order | disk `StaticDir` if `index.html` → else embeddist → else placeholder | `internal/httpapi/static.go` |
| Default StaticDir | `<root>/web/dist` — consumers without `web/` hit **stub** embed | `server.go` / `local_http.go` |
| Embed today | **Stub** only (“Embedded GUI stub…”) | `embeddist/index.html` |
| Embed recipe | Manual npm build + `cp -a dist → embeddist`; README still teaches **two-artifact everyday** | `embeddist/README.md` |
| Makefile/CI | **No** Trace-root make/CI target syncing `web/dist` → `embeddist` (as of planner pass) | repo root |
| Port | P32-PORT fail + `FormatAddrInUseMessage` manual `--addr` example `:7433` | `addr_in_use.go` |
| Default bind | `127.0.0.1:7432` | `httpapi.DefaultAddr` |
| Help | “no auto free-port”; prefer disk `web/dist` | `help.go`, `local_http.go` usage |
| Docs | Disk wins; optional embeddist copy | `docs/gui-quickstart.md` |

## L3 supersession (must document in RESEARCH)

| Prior | Stance | P34 |
|-------|--------|-----|
| P32 RESEARCH | Friendly fail + `--addr`; UA auto-increment optional/#2 deferred | L3 **ships** auto for **default** bind |
| P33 RESEARCH | **Reject** UA auto-port (“conflicts with P32-PORT / multi-project `--addr`”) | L3 **overturns** that reject for default bind only |
| Explicit `--addr` | Fail if busy | **Unchanged** (strict) |

## Planner-locked peer cite (for 01)

| Peer | Cite | Borrow under L3 |
|------|------|-----------------|
| UA viewer | `similar projects/Understand-Anything/understand-anything-plugin/packages/viewer/bin/viewer.mjs` | Default busy → `attemptPort+1`, up to **10** attempts; **skip** hop when `--port` explicit; bind **127.0.0.1**; print URL; best-effort open + `--no-open`. Reject public defaults / daemon. |

## Algorithm candidates 01 must compare

1. UA-style increment from `7432` (max N)
2. Fixed scan range then fail
3. `:0` OS assign then advertise (usually inferior default UX)

Recommend one + `--addr` detection rule (`flag.Changed` vs default-string compare).

## Research rejects (01 must document)

1. Require consumer `web/` / two-artifact as consumer primary.
2. SPA under consumer `.trace/` as primary.
3. Keep no-auto-port on default bind (vs L3).
4. Auto-hop on explicit `--addr`.
5. Public default bind / always-on daemon / SaaS.

## Docs audit seeds (01 expand)

`docs/gui-quickstart.md`, `embeddist/README.md`, `cmd/trace/help.go`, `cmd/trace/local_http.go`, skim `README.md` / `AGENTS.md` / `web/README.md`.

## Out of this scope

- Writing `PLAN.md` (S01), shipping embed (S02), shipping auto-port (S03), docs flip (S04).
- Authoring `RESEARCH.md` in the planner row.
- Reopening L1–L4.
