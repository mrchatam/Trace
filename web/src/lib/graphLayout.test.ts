/**
 * Node strip-types unit checks for graphLayout (no vitest in web/).
 * Run: node --experimental-strip-types --test src/lib/graphLayout.test.ts
 */
import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import {
  COMPACT_NODE_THRESHOLD,
  DOT_LABEL_MIN_ZOOM,
  EDGE_LABEL_MAX,
  EDGE_OVERVIEW_MAX,
  EDGE_SAMPLE_MIN_COUNT,
  EDGE_SAMPLE_RATIO,
  FULL_CARD_MIN_ZOOM,
  FULL_NODE_MIN_ZOOM,
  GOAL_BAND_TINT_COUNT,
  LOD_MINIMAL_MAX_ZOOM,
  PROJECT_FIT_PADDING,
  computeForceLayout,
  computeOverviewPositions,
  computeSemanticLayout,
  countNodesByKind,
  edgeHoverStyle,
  edgeStrokeOpacity,
  filterEdgesByNodeIds,
  filterEdgesForLod,
  filterEdgesForOverview,
  filterNodesByKinds,
  getNodeLod,
  goalBandTintIndex,
  inferGoalBands,
  isEdgeHighlighted,
  kindCssKey,
  SEMANTIC_LANE_WIDTH,
  SEMANTIC_MIN_NODE_DISTANCE,
  shouldRenderEdges,
  shouldShowEdgeLabels,
  shouldShowFullNode,
  shouldShowNodeLabel,
  shouldUseCompactNodes,
  stableEdgeId,
  statusCssKey,
  truncateNodeLabel,
} from './graphLayout.ts'

describe('computeForceLayout', () => {
  it('returns a position for every node id', () => {
    const ids = ['a', 'b', 'c', 'd']
    const edges = [
      { from: 'a', to: 'b' },
      { from: 'b', to: 'c' },
      { from: 'c', to: 'd' },
    ]
    const positions = computeForceLayout(ids, edges, { iterations: 40 })
    assert.equal(positions.size, 4)
    for (const id of ids) {
      const p = positions.get(id)
      assert.ok(p)
      assert.equal(typeof p!.x, 'number')
      assert.equal(typeof p!.y, 'number')
    }
  })

  it('ignores edges whose endpoints are missing', () => {
    const positions = computeForceLayout(['x', 'y'], [{ from: 'x', to: 'missing' }], {
      iterations: 20,
    })
    assert.equal(positions.size, 2)
  })

  it('accepts nodes with kind metadata for soft lane bias', () => {
    const nodes = [
      { id: 'g1', kind: 'goal' },
      { id: 't1', kind: 'task' },
      { id: 'd1', kind: 'decision' },
    ]
    const edges = [
      { from: 'g1', to: 't1' },
      { from: 't1', to: 'd1' },
    ]
    const pos = computeForceLayout(nodes, edges, { iterations: 60 })
    assert.equal(pos.size, 3)
    assert.ok(pos.get('g1')!.x <= pos.get('t1')!.x + 120)
  })
})

describe('countNodesByKind', () => {
  it('aggregates counts by kind', () => {
    const counts = countNodesByKind([
      { kind: 'task' },
      { kind: 'task' },
      { kind: 'goal' },
    ])
    assert.equal(counts.get('task'), 2)
    assert.equal(counts.get('goal'), 1)
  })
})

describe('filterNodesByKinds', () => {
  it('keeps only enabled kinds', () => {
    const nodes = [
      { id: '1', kind: 'task' },
      { id: '2', kind: 'goal' },
    ]
    const filtered = filterNodesByKinds(nodes, new Set(['task']))
    assert.deepEqual(
      filtered.map((n) => n.id),
      ['1'],
    )
  })

  it('returns no nodes when enabled set is empty', () => {
    const nodes = [{ id: '1', kind: 'task' }]
    assert.equal(filterNodesByKinds(nodes, new Set()).length, 0)
  })
})

describe('filterEdgesByNodeIds', () => {
  it('drops edges with filtered-out endpoints', () => {
    const edges = [
      { from: 'a', to: 'b', rel: 'blocks' },
      { from: 'a', to: 'c', rel: 'relates' },
    ]
    const kept = filterEdgesByNodeIds(edges, new Set(['a', 'b']))
    assert.equal(kept.length, 1)
    assert.equal(kept[0]!.to, 'b')
  })
})

describe('shouldUseCompactNodes', () => {
  it('enables compact mode for large project graphs', () => {
    assert.equal(shouldUseCompactNodes(COMPACT_NODE_THRESHOLD, 'project'), true)
    assert.equal(shouldUseCompactNodes(COMPACT_NODE_THRESHOLD - 1, 'project'), false)
    assert.equal(shouldUseCompactNodes(100, 'neighborhood'), false)
  })
})

describe('shouldShowFullNode', () => {
  it('reveals full cards only on focus or very high zoom', () => {
    assert.equal(shouldShowFullNode(true, 0.4, 'n1', null, 'c1'), false)
    assert.equal(shouldShowFullNode(true, 0.9, 'n1', null, 'c1'), false)
    assert.equal(shouldShowFullNode(true, 1.1, 'n1', null, 'c1'), false)
    assert.equal(shouldShowFullNode(true, FULL_CARD_MIN_ZOOM, 'n1', null, 'c1'), true)
    assert.equal(shouldShowFullNode(true, 0.4, 'n1', 'n1', 'c1'), true)
    assert.equal(shouldShowFullNode(true, 0.4, 'c1', null, 'c1'), true)
    // Hover must not promote to full (avoids remount flash)
    assert.equal(shouldShowFullNode(true, 0.4, 'n1', null, 'c1', 'n1'), false)
    assert.equal(shouldShowFullNode(false, 0.4, 'n1', null, 'c1'), false)
    assert.equal(shouldShowFullNode(false, 0.9, 'n1', null, 'c1'), false)
    assert.equal(shouldShowFullNode(false, FULL_CARD_MIN_ZOOM, 'n1', null, 'c1'), true)
  })
})

describe('getNodeLod', () => {
  it('steps through minimal → dot → full by zoom (full only at very high zoom)', () => {
    assert.equal(getNodeLod(true, 0.2, 'n1', null, 'c1'), 'minimal')
    assert.equal(getNodeLod(true, 0.5, 'n1', null, 'c1'), 'dot')
    assert.equal(getNodeLod(true, 0.9, 'n1', null, 'c1'), 'dot')
    assert.equal(getNodeLod(true, 1.1, 'n1', null, 'c1'), 'dot')
    assert.equal(getNodeLod(true, FULL_CARD_MIN_ZOOM, 'n1', null, 'c1'), 'full')
    assert.equal(getNodeLod(false, 0.9, 'n1', null, 'c1'), 'dot')
  })

  it('promotes selected/center to full but ignores hover', () => {
    assert.equal(getNodeLod(true, 0.1, 'n1', 'n1', 'c1'), 'full')
    assert.equal(getNodeLod(true, 0.1, 'c1', null, 'c1'), 'full')
    assert.equal(getNodeLod(true, 0.1, 'h1', null, 'c1', 'h1'), 'minimal')
  })

  it('forces dot/minimal LOD in compact overview until very high zoom', () => {
    assert.equal(getNodeLod(true, 1.1, 'n1', null, 'c1'), 'dot')
    assert.equal(getNodeLod(true, 0.2, 'n1', null, 'c1'), 'minimal')
    assert.equal(getNodeLod(true, FULL_CARD_MIN_ZOOM, 'n1', null, 'c1'), 'full')
  })
})

describe('shouldShowNodeLabel', () => {
  it('shows labels at full LOD and defers dot labels until high zoom (hover via CSS)', () => {
    assert.equal(shouldShowNodeLabel('full', 0.5), true)
    assert.equal(shouldShowNodeLabel('dot', 0.5), false)
    assert.equal(shouldShowNodeLabel('dot', DOT_LABEL_MIN_ZOOM), true)
    assert.equal(shouldShowNodeLabel('minimal', 1.5), false)
    // Hover flag is ignored — CSS :hover reveals labels without rebuild
    assert.equal(shouldShowNodeLabel('dot', 0.5, true), false)
  })
})

describe('truncateNodeLabel', () => {
  it('truncates long labels for dot LOD', () => {
    assert.equal(truncateNodeLabel('Short title'), 'Short title')
    assert.equal(truncateNodeLabel('A'.repeat(40)).endsWith('…'), true)
    assert.ok(truncateNodeLabel('A'.repeat(40)).length <= 28)
  })
})

describe('shouldRenderEdges', () => {
  it('hides edges when zoomed out past overview threshold', () => {
    assert.equal(shouldRenderEdges(0.3), false)
    assert.equal(shouldRenderEdges(0.55), true)
  })
})

describe('filterEdgesForLod', () => {
  it('returns no edges when zoom is below hide threshold', () => {
    const edges = [{ from: 'a', to: 'b' }]
    assert.equal(filterEdgesForLod(edges, 0.1, new Set()).length, 0)
  })

  it('samples dense graphs at low zoom but keeps focus edges', () => {
    const edges = Array.from({ length: EDGE_SAMPLE_MIN_COUNT + 10 }, (_, i) => ({
      from: `n${i}`,
      to: `n${i + 1}`,
    }))
    edges.push({ from: 'center', to: 'n0' })
    const filtered = filterEdgesForLod(edges, 0.52, new Set(['center']))
    assert.ok(filtered.length < edges.length)
    assert.ok(filtered.some((e) => e.from === 'center'))
  })

  it('keeps all edges at high zoom', () => {
    const edges = Array.from({ length: EDGE_SAMPLE_MIN_COUNT + 5 }, (_, i) => ({
      from: `a${i}`,
      to: `b${i}`,
    }))
    assert.equal(filterEdgesForLod(edges, 0.8, new Set()).length, edges.length)
  })
})

describe('filterEdgesForOverview', () => {
  it('shows only focus-incident edges in project mode by default', () => {
    const edges = [
      { from: 'a', to: 'b', rel: 'relates_to' },
      { from: 'b', to: 'c', rel: 'blocks' },
      { from: 'c', to: 'd', rel: 'mentions' },
    ]
    const filtered = filterEdgesForOverview(edges, 0.8, new Set(['b']), false)
    assert.deepEqual(
      filtered.map((e) => `${e.from}-${e.to}`),
      ['a-b', 'b-c'],
    )
  })

  it('caps overview edges at EDGE_OVERVIEW_MAX with priority rels', () => {
    const edges = Array.from({ length: EDGE_OVERVIEW_MAX + 40 }, (_, i) => ({
      from: `n${i}`,
      to: `n${i + 1}`,
      rel: i % 5 === 0 ? 'goal_has_task' : 'misc',
    }))
    const filtered = filterEdgesForOverview(edges, 0.8, new Set(), false)
    assert.ok(filtered.length <= EDGE_OVERVIEW_MAX)
    assert.ok(filtered.some((e) => e.rel === 'goal_has_task'))
  })

  it('delegates to LOD sampling when show all edges is enabled', () => {
    const edges = Array.from({ length: EDGE_SAMPLE_MIN_COUNT + 10 }, (_, i) => ({
      from: `n${i}`,
      to: `n${i + 1}`,
    }))
    const all = filterEdgesForOverview(edges, 0.8, new Set(), true)
    assert.equal(all.length, edges.length)
  })
})

describe('isEdgeHighlighted', () => {
  it('matches edges incident to hovered node', () => {
    assert.equal(isEdgeHighlighted({ from: 'a', to: 'b' }, 'a'), true)
    assert.equal(isEdgeHighlighted({ from: 'a', to: 'b' }, 'c'), false)
    assert.equal(isEdgeHighlighted({ from: 'a', to: 'b' }, null), false)
  })
})

describe('edgeStrokeOpacity', () => {
  it('reduces opacity when zoomed out or edge count is high', () => {
    assert.equal(edgeStrokeOpacity(0.3, 50), 0)
    assert.ok(edgeStrokeOpacity(0.52, EDGE_SAMPLE_MIN_COUNT + 1) < edgeStrokeOpacity(0.9, 50))
    assert.ok(edgeStrokeOpacity(0.52, EDGE_SAMPLE_MIN_COUNT + 1) < edgeStrokeOpacity(0.52, 10))
  })
})

describe('computeOverviewPositions', () => {
  it('places every node id', () => {
    const ids = ['s1', 'a', 'b']
    const positions = computeOverviewPositions(ids, new Set(['s1']))
    assert.equal(positions.size, 3)
    for (const id of ids) assert.ok(positions.get(id))
  })
})

describe('kindCssKey', () => {
  it('maps underscores to hyphens for CSS tokens', () => {
    assert.equal(kindCssKey('plan_change'), 'plan-change')
    assert.equal(kindCssKey('task'), 'task')
  })
})

describe('statusCssKey', () => {
  it('maps work_state to CSS data-status tokens', () => {
    assert.equal(statusCssKey('IN_PROGRESS'), 'in-progress')
    assert.equal(statusCssKey('DONE'), 'done')
    assert.equal(statusCssKey('SKIPPED'), 'skipped')
    assert.equal(statusCssKey(null), undefined)
    assert.equal(statusCssKey(undefined), undefined)
  })
})

describe('goalBandTintIndex', () => {
  it('returns stable tint slots in range for goal ids', () => {
    const a = goalBandTintIndex('goal-alpha')
    const b = goalBandTintIndex('goal-alpha')
    const c = goalBandTintIndex('goal-beta')
    assert.equal(a, b)
    assert.ok(a >= 0 && a < GOAL_BAND_TINT_COUNT)
    assert.ok(c >= 0 && c < GOAL_BAND_TINT_COUNT)
    assert.equal(goalBandTintIndex(null), -1)
    assert.equal(goalBandTintIndex(undefined), -1)
    assert.equal(goalBandTintIndex('__ungrouped__'), -1)
  })
})

describe('stableEdgeId', () => {
  it('builds index-free edge ids', () => {
    assert.equal(stableEdgeId({ from: 'a', to: 'b', rel: 'blocks' }), 'a|blocks|b')
    assert.equal(stableEdgeId({ from: 'a', to: 'b' }), 'a||b')
  })
})

describe('edgeHoverStyle', () => {
  it('boosts stroke without changing membership semantics', () => {
    const hot = edgeHoverStyle({ stroke: 'var(--graph-edge-stroke)', strokeWidth: 1, opacity: 0.42 }, true)
    assert.equal(hot.stroke, 'var(--accent)')
    assert.equal(hot.strokeWidth, 2)
    assert.ok(hot.opacity > 0.42)
    const cold = edgeHoverStyle({ opacity: 0.42 }, false)
    assert.equal(cold.stroke, 'var(--graph-edge-stroke)')
    assert.equal(cold.strokeWidth, 1)
    assert.equal(cold.opacity, 0.42)
  })
})

describe('PROJECT_FIT_PADDING', () => {
  it('uses comfortable content padding (not overly zoomed out)', () => {
    assert.ok(PROJECT_FIT_PADDING >= 0.12 && PROJECT_FIT_PADDING <= 0.18)
  })
})

describe('inferGoalBands', () => {
  it('assigns tasks to goal_id and propagates via edges', () => {
    const nodes = [
      { id: 'g1', kind: 'goal', title: 'G1' },
      { id: 't1', kind: 'task', title: 'T1', goal_id: 'g1' },
      { id: 'd1', kind: 'discovery', title: 'D1' },
    ]
    const edges = [
      { from: 't1', to: 'd1' },
      { from: 'g1', to: 't1' },
    ]
    const bands = inferGoalBands(nodes, edges)
    assert.equal(bands.get('g1'), 'g1')
    assert.equal(bands.get('t1'), 'g1')
    assert.equal(bands.get('d1'), 'g1')
  })
})

describe('computeSemanticLayout', () => {
  it('returns deterministic positions for every node', () => {
    const nodes = [
      { id: 'g1', kind: 'goal', title: 'Alpha goal' },
      { id: 't1', kind: 'task', title: 'Task A', goal_id: 'g1', work_state: 'PENDING' },
      { id: 't2', kind: 'task', title: 'Task B', goal_id: 'g1', work_state: 'IN_PROGRESS' },
      { id: 'd1', kind: 'decision', title: 'Decide' },
    ]
    const edges = [
      { from: 'g1', to: 't1' },
      { from: 'g1', to: 't2' },
      { from: 't1', to: 'd1' },
    ]
    const a = computeSemanticLayout(nodes, edges)
    const b = computeSemanticLayout(nodes, edges)
    assert.equal(a.size, 4)
    for (const id of ['g1', 't1', 't2', 'd1']) {
      assert.deepEqual(a.get(id), b.get(id))
    }
  })

  it('places kind lanes left-to-right (goal before task before decision)', () => {
    const nodes = [
      { id: 'g1', kind: 'goal', title: 'G' },
      { id: 't1', kind: 'task', title: 'T', goal_id: 'g1' },
      { id: 'd1', kind: 'decision', title: 'D' },
    ]
    const edges = [{ from: 'g1', to: 't1' }, { from: 't1', to: 'd1' }]
    const pos = computeSemanticLayout(nodes, edges)
    assert.ok(pos.get('g1')!.x < pos.get('t1')!.x)
    assert.ok(pos.get('t1')!.x < pos.get('d1')!.x)
  })

  it('clusters nodes under the same goal on similar Y bands', () => {
    const nodes = [
      { id: 'g1', kind: 'goal', title: 'G' },
      { id: 't1', kind: 'task', title: 'T', goal_id: 'g1' },
    ]
    const edges = [{ from: 'g1', to: 't1' }]
    const pos = computeSemanticLayout(nodes, edges)
    const yDelta = Math.abs(pos.get('g1')!.y - pos.get('t1')!.y)
    assert.ok(yDelta < 200)
  })

  it('keeps lane centers at least 200px apart on X', () => {
    const nodes = [
      { id: 'g1', kind: 'goal', title: 'G' },
      { id: 't1', kind: 'task', title: 'T', goal_id: 'g1' },
      { id: 'd1', kind: 'decision', title: 'D' },
      { id: 'r1', kind: 'review', title: 'R' },
    ]
    const edges = [{ from: 'g1', to: 't1' }]
    const pos = computeSemanticLayout(nodes, edges)
    const xs = [...pos.values()].map((p) => p.x).sort((a, b) => a - b)
    for (let i = 1; i < xs.length; i++) {
      assert.ok(
        xs[i]! - xs[i - 1]! >= 200 || xs[i] === xs[i - 1],
        `lane X spread too tight: ${xs[i - 1]} → ${xs[i]}`,
      )
    }
    assert.ok(pos.get('g1')!.x + SEMANTIC_LANE_WIDTH <= pos.get('t1')!.x + 1)
  })

  it('does not place distinct nodes closer than SEMANTIC_MIN_NODE_DISTANCE', () => {
    const nodes = Array.from({ length: 80 }, (_, i) => ({
      id: `t${i}`,
      kind: 'task',
      title: `Task ${i}`,
      goal_id: 'g1',
      work_state: i % 2 === 0 ? 'PENDING' : 'IN_PROGRESS',
    }))
    nodes.unshift({ id: 'g1', kind: 'goal', title: 'Goal' })
    const edges = nodes.slice(1).map((n) => ({ from: 'g1', to: n.id }))
    const pos = computeSemanticLayout(nodes, edges)
    const entries = [...pos.entries()]
    for (let i = 0; i < entries.length; i++) {
      for (let j = i + 1; j < entries.length; j++) {
        const a = entries[i]![1]
        const b = entries[j]![1]
        const dist = Math.hypot(a.x - b.x, a.y - b.y)
        assert.ok(
          dist >= SEMANTIC_MIN_NODE_DISTANCE - 0.01,
          `nodes ${entries[i]![0]} and ${entries[j]![0]} too close: ${dist.toFixed(1)}px`,
        )
      }
    }
  })

  it('preserves horizontal spread for large graphs (no pillar collapse)', () => {
    const kinds = ['goal', 'task', 'decision', 'review', 'evidence'] as const
    const nodes = Array.from({ length: 120 }, (_, i) => ({
      id: `n${i}`,
      kind: kinds[i % kinds.length]!,
      title: `Node ${i}`,
      goal_id: i === 0 ? undefined : 'n0',
    }))
    nodes[0] = { id: 'n0', kind: 'goal', title: 'Root goal' }
    const edges = nodes.slice(1, 40).map((n) => ({ from: 'n0', to: n.id }))
    const pos = computeSemanticLayout(nodes, edges)
    const xs = [...pos.values()].map((p) => p.x)
    const xSpread = Math.max(...xs) - Math.min(...xs)
    assert.ok(xSpread >= 800, `X spread ${xSpread}px too narrow — layout collapsed`)
  })
})

describe('shouldShowEdgeLabels', () => {
  it('hides labels when zoomed out or edge count is high', () => {
    assert.equal(shouldShowEdgeLabels(10, 0.8), true)
    assert.equal(shouldShowEdgeLabels(10, 0.4), false)
    assert.equal(shouldShowEdgeLabels(EDGE_LABEL_MAX + 1, 1), false)
  })
})
