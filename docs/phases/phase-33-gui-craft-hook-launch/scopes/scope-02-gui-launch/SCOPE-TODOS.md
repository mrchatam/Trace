# Scope 02 — board map

**S02 launch** — `trace gui` + PATH. Serial: **P33-S02-00 → P33-S02-01 → P33-S02-02**. Board order places after S01 — **do not start until prior rows done**.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 577 | P33-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock CLI + PATH + land `/` |
| 578 | P33-S02-01 | [01-implement.md](01-implement.md) | Implementer | `gui` + open-browser + tests + PATH teach |
| 579 | P33-S02-02 | [02-review.md](02-review.md) · [REVIEW.md](REVIEW.md) | Reviewer | **PASS** (high) — Theme C + Law 19 + PATH ≠ install + land `/` |

## Planner locks (binding — from S02-00 + S00 RESEARCH + S01 land)

| Lock | Value |
|------|-------|
| CLI | Subcommand **`trace gui`** only (no `-gui` this scope); keep `serve` |
| Flags | Serve set + **`--no-open`**; `-C`/`--root` parity |
| Bind | `127.0.0.1:7432`; refuse remote without allow; **reject UA auto-port** |
| Open | Best-effort after listen; print URL; open fail ≠ listen fail |
| Land | Browser URL **`http://{addr}/`** Explore — **≠** `/overview` |
| PATH #1 | `go install github.com/mrchatam/Trace/cmd/trace@…` |
| PATH ≠ | `trace install` = agents/MCP/hooks only |
| Docs | Help + minimal tip OK; **S05** owns full quickstart primary flip |
| Law 19 | CLI → httpapi only |

## S00 research leans (still binding)

- Prefer UA **open + `--no-open`**; reject UA **auto-port**.
- PATH teach ≠ agents `trace install`.
- Theme C: one-command-from-cwd UX.

## S01 review thickening (still binding)

- `trace gui` must land on Explore **`/`**, not Nav Overview.

## Out of this scope

- Explore UI / seeds (S03), colorize (S04), full docs primary flip (S05), packaging (brew/deb).
