# Scope 03 — board map

**S03 auto free-port.** Serial: **P34-S03-00 → P34-S03-01 → P34-S03-02**. Start only after **P34-S02-02** PASS (confirmed).

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 601 | P34-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock algorithm + flags + thicken 01/02 (this row) |
| 602 | P34-S03-01 | [01-implement.md](01-implement.md) | Implementer | Shared httpapi UA-incr + `flag.Changed` + post-bind print/open + T4–T7/T11 + help |
| 603 | P34-S03-02 | [02-review.md](02-review.md) | Reviewer | L3/L4 + PLAN checklist; no public-bind / embed creep |

## Inputs (verified baseline 2026-08-21)

| Source | Fact |
|--------|------|
| [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) | UA-incr **7432–7441** (max 10); `flag.Changed` strict; shared **httpapi** helper; serve print **post-bind**; T4–T7, T11 |
| L3 / L4 | Default busy → auto free loopback + print/open chosen URL; `--addr` pin strict; one process = one root |
| Live `DefaultAddr` | `internal/httpapi/bind.go` → `127.0.0.1:7432` |
| Live listen | `server.go` `ListenAndServe`: single `net.Listen(s.addr)`; `OnListening(s.addr)`; **no** hop |
| Live CLI | `local_http.go`: shared gui/serve; serve prints **before** listen; gui print/open in `OnListening`; addr-in-use → `FormatAddrInUseMessage` (manual `--addr` only) |
| Live help | `help.go` / usageGUI still teach “no auto free-port” |
| Live tests | `TestGuiAddrInUseFriendlyMessage` / `TestServeAddrInUseFriendlyMessage` = explicit `--addr` busy fail; no default-hop / concurrent-default coverage yet |
| Peer UA | `viewer.mjs`: `listen(port, 10)`, `portExplicit`, host fixed `127.0.0.1` |

## Locked answers (P34-S03-00)

1. **Algorithm:** UA-increment, start **7432**, host **127.0.0.1**, max **10** (`7432`–`7441`); `+1` on `IsAddrInUse` when `!AddrExplicit`.
2. **Helper:** **`internal/httpapi`** (`ListenAndServe` / adjacent helper); CLI passes `flag.Changed` via Options.
3. **Tests:** T4 busy→next; T5 concurrent distinct + chosen open URL; T6 explicit fail; T7 exhausted; T11 free→7432.

## Out of this scope

- Embed/StaticDir (S02 done); full docs/quickstart/AGENTS flip (S04); VERIFY (S05); OS `:0` default; public bind default; Explore redesign.
