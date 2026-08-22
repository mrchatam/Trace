# P33-S04-02 — Craft review

## Metadata
- id: P33-S04-02
- todo_ids: [P33-S04-02]
- role: reviewer
- skills: [impeccable, ui-ux-pro-max, frontend-design, code-review-and-quality]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Independent Theme A craft + a11y review of S04. impeccable **critique**/audit mindset + ui-ux-pro-max contrast/color-not-only + frontend-design signature check. Spawn remediation if still bland, contrast fails, or labels were removed for color-only. Next **P33-S05-00**.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [`../scope-01-design-ux/DESIGN.md`](../scope-01-design-ux/DESIGN.md)
- [01-implement.md](01-implement.md) application map + exit criteria
- Live: `web/src/styles/tokens.css`, `web/src/styles/app.css`, `Graph.tsx` (hooks only — expect **no** compose/route diff)
- S03 contract: [`../scope-03-explore-graph/REVIEW.md`](../scope-03-explore-graph/REVIEW.md)

## Session start

Fresh subagent. Follow agent-loop-protocol Session start. Load review skills before scoring.

## Preflight

1. Diff vs S03 baseline: styles only? Flag any `overviewCompose.ts` / `App.tsx` / budget edits as **blocker** unless justified.
2. Confirm Explore still `/`; `/overview` unchanged ops screen.
3. Re-run or trust implement evidence: `npm run build`; optional smoke e2e `e2e/s03-depth.spec.ts` if CSS risk to hooks.

## Checklist

### Theme A / craft
- [ ] Color system deliberate — matches S01 forest-moss refine (not reinvented palette)
- [ ] Not gray-monochrome; kind **chroma strip** (3px) + soft fill ~14% + kind label color
- [ ] State secondary channel labeled + `--state-*` (color-not-only)
- [ ] S03 text labels retained — S04 must not remove kind/state text for color-only encoding
- [ ] Center/seed moss emphasis coexists with kind literacy
- [ ] Kind spectrum not rainbow-tipping; capability/assumption muted + labeled
- [ ] Anti-slop: no glow/HUD/purple-cream/broadsheet/slate+#22C55E wholesale/emoji
- [ ] Empty/error surfaces use semantic tokens without vibe-only copy changes (S03 owns copy)

### A11y / motion
- [ ] Contrast ≥4.5:1 text / ≥3:1 strokes+focus evidenced light **and** dark (Notes or artifact)
- [ ] `:focus-visible` / `--focus` ring present; not removed
- [ ] `prefers-reduced-motion` still disables node/inspector transitions
- [ ] ui-ux-pro-max Color Only + Color Contrast satisfied

### Contracts
- [ ] Explore still `/`; overview compose / budgets / routes untouched
- [ ] Law 19 — visual adapter only
- [ ] Skills evidence on implement row (impeccable + ui-ux-pro-max + frontend-design)
- [ ] Optional: Tasks CTA polish is visual-only

### Skills (reviewer)
- [ ] impeccable critique/audit against craft-floor
- [ ] ui-ux-pro-max domain ux/color spot-check
- [ ] frontend-design: signature present; chrome restrained
- [ ] code-review-and-quality: CSS specificity / no dead competing hexes

## Findings

Severity: blocker | high | medium | low | nit. blocker/high → inline fix or spawn `P33-S04-02a`/`02b`.

## Artifact

Write [`REVIEW.md`](REVIEW.md) with PASS/FAIL, confidence, contrast evidence citation, residual risks.

## Exit criteria

- [ ] No open blocker/high without pending follow-up
- [ ] Confidence medium or high with evidence
- [ ] Board Notes; next **P33-S05-00** (or spawn rows if needed)

## Todo updates

Status + notes on **P33-S04-02** only. May thicken upcoming S05 prompts if craft residuals affect polish/docs.
