# OPEN — Port conflict & multi-project `trace serve`

**Status:** Phase 32 **CLOSED** (`P32-S06-02`). **#1 shipped in S02**; **#3/#4 docs closed in S05**; **VERIFY + DR-HANDOFF** ticked #1 + docs; **#2 auto free-port / `:0` remains deferred** (non-blocking residual — not a phase successor). See `VERIFY-NOTES.md` + `REVIEW-NOTES.md`.  
**Intake:** human 2026-08-21 — light review only; do not treat as full RCA.

## Light review (current behavior)

| Item | Finding |
|------|---------|
| Default | `127.0.0.1:7432` (`httpapi.DefaultAddr`) — unchanged |
| Bind | `net.Listen("tcp", addr)` in `ListenAndServe` — **no** free-port search, **no** `:0` default (**#2 not shipped**) |
| On conflict | **#1 done:** `httpapi.IsAddrInUse` + `FormatAddrInUseMessage`; `serve` friendly stderr + help example `127.0.0.1:7433`; still **fail-on-conflict** |
| Multi-project | One `trace serve` = one `--root`/`-C`. Second project needs distinct `--addr` (message + help now say so) |
| Workaround today | Manual `./bin/trace serve -C <proj> --addr 127.0.0.1:<other>` |

## Must investigate in Phase 32 (S00 note → S02 own)

Pick and ship (or document defer with reason) at least one coherent story:

1. **Clearer failure** — detect `EADDRINUSE` / “address already in use”; stderr message suggests `--addr` examples  
2. **Auto free port** — e.g. try 7432 then 7433… or `--addr 127.0.0.1:0` and print the chosen URL  
3. **Multi-project UX** — docs + optional helper (script/flag) to run N serves on distinct ports; and/or a future single-process project switcher (**larger**; may defer with ADR note)  
4. **GUI/docs** — `gui-quickstart.md` documents multi-project port pattern  

Keep local-first / loopback defaults. Do not invent always-on multi-tenant hosting.

## Owner scopes

- **S00** — confirm light review; list peer patterns (UA tokenized URL, etc.) briefly; recommend option(s) for S02  
- **S02** — **always owns this item**, even when API work is `NO-GAPS.md`. Implement chosen serve/UX measure(s) + tests. Prefer ship at least **#1** (friendly `EADDRINUSE` / in-use message + `--addr` guidance). `NO-PORT-CHANGE.md` only if deliberately deferred with written reason (discouraged).  
- **S05** — docs/quickstart update (multi-project / port pattern) — **done** (P32-S05-01/02)
- **S06** — VERIFY must tick P32-PORT addressed
