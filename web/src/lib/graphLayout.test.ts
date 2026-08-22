/**
 * Node strip-types unit checks for graphLayout (no vitest in web/).
 * Run: node --experimental-strip-types --test src/lib/graphLayout.test.ts
 */
import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import {
  COMPACT_NODE_THRESHOLD,
  EDGE_LABEL_MAX,
  EDGE_SAMPLE_MIN_COUNT,
  EDGE_SAMPLE_RATIO,
  FULL_CARD_MIN_ZOOM,
  FULL_NODE_MIN_ZOOM,
  LOD_MINIMAL_MAX_ZOOM,
  computeForceLayout,
  computeOverviewPositions,
  computeSemanticLayout,
  countNodesByKind,
  edgeStrokeOpacity,
  filterEdgesByNodeIds,
  filterEdgesForLod,
  filterNodesByKinds,
  getNodeLod,
  inferGoalBands,
  kindCssKey,
  shouldRenderEdges,
  shouldShowEdgeLabels,
  shouldShowFullNode,
  shouldShowNodeLabel,
  shouldUseCompactNodes,
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
    assert.equal(shouldShowFullNode(true, 0.4, 'n1', null, 'c1', 'n1'), true)
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

  it('promotes focused nodes to full at any zoom', () => {
    assert.equal(getNodeLod(true, 0.1, 'n1', 'n1', 'c1'), 'full')
    assert.equal(getNodeLod(true, 0.1, 'c1', null, 'c1'), 'full')
    assert.equal(getNodeLod(true, 0.1, 'h1', null, 'c1', 'h1'), 'full')
  })
})

describe('shouldShowNodeLabel', () => {
  it('shows labels at dot and full LOD', () => {
    assert.equal(shouldShowNodeLabel('full'), true)
    assert.equal(shouldShowNodeLabel('dot'), true)
    assert.equal(shouldShowNodeLabel('minimal'), false)
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
    assert.equal(shouldRenderEdges(0.1), false)
    assert.equal(shouldRenderEdges(0.3), true)
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
    const filtered = filterEdgesForLod(edges, LOD_MINIMAL_MAX_ZOOM, new Set(['center']))
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

describe('edgeStrokeOpacity', () => {
  it('reduces opacity when zoomed out or edge count is high', () => {
    assert.equal(edgeStrokeOpacity(0.1, 50), 0)
    assert.ok(edgeStrokeOpacity(0.3, 50) < edgeStrokeOpacity(0.9, 50))
    assert.ok(edgeStrokeOpacity(0.3, EDGE_SAMPLE_MIN_COUNT + 1) < edgeStrokeOpacity(0.3, 10))
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
})

describe('shouldShowEdgeLabels', () => {
  it('hides labels when zoomed out or edge count is high', () => {
    assert.equal(shouldShowEdgeLabels(10, 0.8), true)
    assert.equal(shouldShowEdgeLabels(10, 0.4), false)
    assert.equal(shouldShowEdgeLabels(EDGE_LABEL_MAX + 1, 1), false)
  })
})
