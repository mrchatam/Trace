# Phase 33 S01 — DESIGN (tokens + craft)

## POV / mode

**Operate** (impeccable). Trace GUI is a local-first knowledge-graph tool for coding agents — scanability and the Explore canvas outrank brand theater.

**Visual authority:** **Refine** live `web/src/styles/tokens.css` forest-moss / IBM Plex Sans+Mono world (light `#f4f5f3` + moss accent `#2f5d3a`; dark `#121512` + moss `#7cbc8a`). Keep light **and** dark. Do **not** replace with purple/cream/broadsheet, cinema glass/glow, or wholesale slate + `#22C55E` “AI SaaS dark.”

**Color strategy (impeccable colorize):** Roles over swatches — surface hierarchy, action/focus rarity, kind chroma for graph literacy, state chroma for task work. Brand lives in precise moss neutrals + kind borders, not HUD glow.

**S00 cite:** Overview IA is **(D)+(B)+(C)** (seeds → budgeted `getGraph` → progressive expand). This file owns tokens/craft only; Explore open flow lives in [`UX-IA.md`](UX-IA.md). **S04** lands these tokens into CSS; this scope ships markdown only.

## Skills evidence

### ui-ux-pro-max

| Query | Domain / flags | Accepted | Rejected |
|-------|----------------|----------|----------|
| `developer tool knowledge graph local-first coding agent operate dense forest moss` | `--design-system -p Trace --density 8 --motion 3 --variance 4` | Dense Operate dials (density 8, motion 3); light+dark pairing; anti-emoji icons; 150–300ms motion; WCAG AA checklist | **Modern Dark (Cinema)** / glassmorphism / glow / indigo ambient blobs; FAQ landing pattern (wrong product surface); Atkinson font swap; **slate primary `#1E293B` + accent `#22C55E`** wholesale |
| `forest moss green developer tool accessible palette dark light` | `--domain color` | Accessible green-as-success / tool green *as semantic ok*, not brand replacement | Result 1–3 code-dark+#22C55E; Result 4–5 / 7–8 purple / indigo study clusters |
| `color not sole encoder graph nodes accessibility contrast` | `--domain ux` | **Color Only** + **Color Contrast** (text ≥4.5:1; color + label/icon); keyboard/focus; duration 150–300ms | Color-only borders; low-contrast muted-on-muted |
| `dense dashboard operate minimal no glow` | `--domain style` | **Data-Dense Dashboard** scanability / dense type / subtle hover — adapted to graph canvas (not KPI card grids) | Executive KPI pulse; Flat Mobile color-block hero; E-Ink cream paper wholesale; glow-heavy cinema |

### impeccable

- Mode: **Operate** (`reference/operate.md`) — consistency, density, state vocabulary; brand in details.
- Playbooks: **shape** (brief-only, no CSS), **colorize** (roles + AA floors), **craft-floor** (bans: glow halos, glass decoration, emoji icons, color-only meaning).
- Not loaded for UI edit: craft-floor applied as **authoring constraints** for token/craft brief (no product CSS this row).

### frontend-design — signature risk (one)

**Signature:** Explore nodes carry a **kind chroma strip** (left border 3px + kind label in mono uppercase) with restrained tinted fill — the graph reads as Trace’s entity vocabulary at a glance.

**Risk:** A full kind spectrum can tip into “rainbow dashboard.” Mitigation: low-chroma fills (`color-mix` ~12–18% into surface), saturated border/label only, always show kind text (and optional 1-char glyph), never glow/shadow bloom. Restraint elsewhere (chrome stays moss-neutral).

## Palette

Named roles (extend existing CSS variable names where possible):

| Role | Intent |
|------|--------|
| Forest canvas | Quiet moss-tinted neutrals — tool disappears into the graph |
| Moss accent | Primary action, nav active, focus, seed/center emphasis |
| Kind chroma | Entity literacy on nodes/pills (goal…capability) |
| State chroma | Task `work_state` (and review FAIL-ish) without competing with kind |
| Semantic alert | Danger / warn / ok for banners and review outcomes |

Keep **IBM Plex Sans** (UI) + **IBM Plex Mono** (kind labels, ids, budget chrome). Do not introduce display serifs or Inter/Roboto defaults.

## Token table

Contrast floor: **WCAG AA** — body/UI text ≥ **4.5:1** vs paired surface; meaningful chrome / node stroke / focus ring ≥ **3:1**. Hexes are **candidates for S04**; verify with a contrast checker when landing CSS.

### Surfaces, text, accent, focus, semantic

| Token | Role | Light | Dark | Contrast notes |
|-------|------|-------|------|----------------|
| `--bg` | App canvas | `#f4f5f3` | `#121512` | Keep incumbent moss tint |
| `--bg-elevated` | Panels, inspector, nodes base | `#ffffff` | `#1a1e19` | Text on elevated ≥4.5:1 |
| `--bg-muted` | Chrome wells, budget bar | `#e8eae6` | `#242a23` | Secondary text still ≥4.5:1 on muted |
| `--border` | Separators, default node stroke | `#c9cec6` | `#3a4238` | Stroke vs bg ≥3:1 |
| `--text` | Primary copy | `#1a1f18` | `#e6ebe4` | ≥4.5:1 on `--bg` / elevated |
| `--text-muted` | Meta, kind fallback | `#5a6356` | `#9aa494` | ≥4.5:1 on `--bg`; if fails on muted surface, use `--text` |
| `--accent` | CTA, nav active, center | `#2f5d3a` | `#7cbc8a` | Light: accent on white ≥3:1 for chrome; dark: accent text on dark bg ≥4.5:1 |
| `--accent-hover` | Hover | `#244a2e` | `#9ad0a5` | Darker/lighter moss |
| `--accent-subtle` | Selected/center fill | `#dce8df` | `#243328` | Pair with `--text` |
| `--focus` | Focus ring | `#2f5d3a` | `#7cbc8a` | 2px ring; ≥3:1 vs adjacent bg |
| `--danger` / `--danger-bg` | Error | `#8b2e2e` / `#f5e4e4` | `#e08a8a` / `#3a2222` | Text ≥4.5:1 on bg sibling |
| `--warn` / `--warn-bg` | Warning / partial fail | `#7a5a12` / `#f5edd8` | `#e0c06a` / `#3a3220` | Same |
| `--ok` / `--ok-bg` | Success / DONE-ish | `#1f5c38` / `#dceee3` | `#8bc9a0` / `#1e3326` | Distinct from brand accent usage on nodes |

Retain spacing, radius, nav, duration tokens from live `tokens.css` (`--radius: 4px`, `--duration-fast: 140ms`, `--duration-settle: 200ms`, `--ease-out`).

### Kind tokens (`--kind-*`)

Use for node border, kind-label color, and optional soft fill. Labels always visible.

| Token | Role | Light | Dark | Notes |
|-------|------|-------|------|-------|
| `--kind-goal` | goal | `#1d4d6e` | `#7eb6d4` | Cool blue-teal; not indigo-purple |
| `--kind-task` | task | `#2f5d3a` | `#7cbc8a` | Aligns with moss brand |
| `--kind-decision` | decision | `#6b4a1a` | `#d4a857` | Amber earth |
| `--kind-assumption` | assumption | `#5a4a6e` | `#b5a4c9` | Muted plum — **not** vivid violet UI |
| `--kind-discovery` | discovery | `#1a5c5c` | `#6fc0bc` | Teal |
| `--kind-plan-change` | plan-change | `#5c3d2e` | `#c49a7c` | Clay |
| `--kind-claim` | claim | `#3d4a6b` | `#9aabc9` | Slate-blue (kind only; not app chrome) |
| `--kind-evidence` | evidence | `#3a5a2a` | `#9bc47f` | Leaf |
| `--kind-capability` | capability | `#4a3d6b` | `#b0a0d4` | Soft lilac for caps — label required |
| `--kind-review` | review | `#6b2e3a` | `#d4899a` | Rose; pairs with FAIL state |
| `--kind-unknown` | fallback | `#5a6356` | `#9aa494` | Same as muted text |

Soft fill recipe (S04): `color-mix(in srgb, var(--kind-*) 14%, var(--bg-elevated))`. Border: solid `var(--kind-*)` at 1.5–2px (center/selected may thicken).

### State tokens (`--state-*`)

Encode task `work_state` (and FAIL-ish review) as a **second channel**: right-edge tick, status pill, or mono state chip — never replace kind encoding.

| Token | Role | Light | Dark | Notes |
|-------|------|-------|------|-------|
| `--state-in-progress` | IN_PROGRESS | `#2f5d3a` | `#7cbc8a` | Active work; may match accent |
| `--state-pending` | PENDING / other non-terminal | `#5a6356` | `#9aa494` | Quiet |
| `--state-done` | DONE / SKIPPED | `#1f5c38` | `#8bc9a0` | Align with `--ok` |
| `--state-fail` | FAIL / FAIL-ish review | `#8b2e2e` | `#e08a8a` | Align with `--danger` |
| `--state-blocked` | blocked / blocked-ish | `#7a5a12` | `#e0c06a` | Align with `--warn` |

## Kind + state encoding

Rules for nodes and pills (S04 implements; S03 wires classes/data):

1. **Kind (primary):** left border color = `--kind-*`; `.graph-node__kind` text = kind string (mono, uppercase tracking); optional single SVG glyph — **no emoji**.
2. **State (secondary, tasks):** small state chip or right border tick using `--state-*`; chip includes **text** (`IN_PROGRESS`, etc.).
3. **Color-not-only:** never rely on fill alone — kind label always present; state chip always labeled when shown.
4. **Center / seed:** moss `--accent` + `--accent-subtle` may override border for “overview seed / current center”; kind label still shown.
5. **Selected:** focus-adjacent outline via `--focus` / accent; keep kind border visible (inset or dual-ring without glow bloom).
6. **Pills in lists:** same `--kind-*` / `--state-*` border or text color + label; flat, no glass.

## Craft floor

### Anti-slop bans

- No glow-first nodes, neon edges, or zero-offset colored halos.
- No glassmorphism / frosted blur as decoration.
- No emoji-as-icons; SVG only (consistent stroke).
- No purple-on-white, cream+#terracotta broadsheet, or cinema-indigo HUD.
- No wholesale replacement of forest neutrals with slate `#0F172A` / `#1E293B` + `#22C55E`.
- No treating Nav **Overview** (`/overview`) as the Explore craft target.
- Cards only where interaction requires a container; Explore canvas is not a card grid.

### Motion

- Keep `--duration-fast` / `--duration-settle` / `--ease-out`.
- Micro-interactions 140–200ms (selection, hover border).
- No orchestrated page-load choreography; graph appears with skeleton/budget chrome, then nodes.
- `@media (prefers-reduced-motion: reduce)`: disable node transitions (already patterned in `app.css`).

### Focus / a11y

- Visible `:focus-visible` rings using `--focus` (2–3px), never removed.
- Keyboard: tab order matches chrome → canvas controls → inspector.
- Touch targets for chrome controls ≥44×44px where applicable.

## Rejected directions

| Direction | Why |
|-----------|-----|
| Purple / cream / broadsheet clusters | AI-default; conflicts with forest-moss brief |
| HUD / neon glow / cinema glass | Operate tool, not game UI; craft-floor ban |
| Slate + `#22C55E` wholesale | ui-ux-pro-max default for “dev tool”; erases Trace identity |
| `/overview` as Explore | Theme B / RESEARCH: Explore = `/` Graph |
| Three.js / second SPA | DESIGN-LOCKS out of scope |
| Unbounded full-graph dump styling as “wow” | Laws 6–7 |

## S03 vs S04 (color ownership)

| Owner | Owns | Does not own |
|-------|------|--------------|
| **S03** | Explore open composition + Graph UX per UX-IA; may add `data-kind` / `data-state` hooks | Full shell recolor; inventing new palette |
| **S04** | Land this token table into `tokens.css` / `app.css`; colorize chrome, nav, pills, Explore nodes, empty/error surfaces | Re-opening IA or budgets |
