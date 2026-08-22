# fixtures/x0 — synthetic P0-X ground truth

Apache-2.0 synthetic mini-repo for Trace Phase 00 P0-X (DR-BENCH / DR-SEED).
**Not** a real OSS corpus. Human-curated seed only — no LLM-generated SoT.

## Layout

| Path | Role |
|------|------|
| `LICENSE` | Apache-2.0 |
| `seed/gt.json` | Seed JSON **v1** (evaluator / harness import only — not an agent answer oracle) |
| `src/greeter.ts` | TypeScript: `greet` + import of `./format` |
| `src/math_util.py` | Python: `add` / `hypotenuse` + `import math` |
| `.gitignore` | Ignores local `.trace/` (never commit store) |

## Architecture (planted code)

- Dual-language so P0-X criterion #2 is non-vacuous: TS/JS **and** Python, each with ≥1 symbol and/or import after `trace index`.
- `greeter.ts` imports `./format` (import row) and exports `greet` (symbol).
- `math_util.py` imports `math` and defines `add` / `hypotenuse`.

## Harness note (seed path)

`trace seed import <file>` resolves relative `<file>` under the **`-C` project root**
(then `filepath.Abs`). Absolute paths are unchanged. Harnesses may still pass an
absolute path to `seed/gt.json` (preferred for evals).

Evaluator UUID / causal map: see `evals/x0/GT-MAP.md` (not agent-facing oracle material).
