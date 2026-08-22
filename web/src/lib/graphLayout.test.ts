/**
 * Node strip-types unit checks for graphLayout (no vitest in web/).
 * Run: node --experimental-strip-types --test src/lib/graphLayout.test.ts
 */
import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import {
  COMPACT_NODE_THRESHOLD,
  EDGE_LABEL_MAX,
  computeForceLayout,
  countNodesByKind,
  filterEdgesByNodeIds,
  filterNodesByKinds,
  shouldShowEdgeLabels,
  shouldShowFullNode,
  shouldUseCompactNodes,
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
  it('reveals full cards on zoom, select, or center', () => {
    assert.equal(shouldShowFullNode(true, 0.4, 'n1', null, 'c1'), false)
    assert.equal(shouldShowFullNode(true, 0.9, 'n1', null, 'c1'), true)
    assert.equal(shouldShowFullNode(true, 0.4, 'n1', 'n1', 'c1'), true)
    assert.equal(shouldShowFullNode(true, 0.4, 'c1', null, 'c1'), true)
    assert.equal(shouldShowFullNode(false, 0.4, 'n1', null, 'c1'), true)
  })
})

describe('shouldShowEdgeLabels', () => {
  it('hides labels when zoomed out or edge count is high', () => {
    assert.equal(shouldShowEdgeLabels(10, 0.8), true)
    assert.equal(shouldShowEdgeLabels(10, 0.4), false)
    assert.equal(shouldShowEdgeLabels(EDGE_LABEL_MAX + 1, 1), false)
  })
})
