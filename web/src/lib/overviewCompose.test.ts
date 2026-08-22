/**
 * Node strip-types unit checks for overviewCompose (no vitest in web/).
 * Run: node --experimental-strip-types --test src/lib/overviewCompose.test.ts
 */
import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import type { BoundedGraph, TaskRow } from '../api/ops'
import {
  SEED_CAP,
  SEED_TARGET,
  UI_CAP,
  applySeedMeta,
  composeSeedsFromParts,
  mergeOverviewGraphs,
  partitionSettledGraphs,
  prioritizeTaskSeeds,
} from './overviewCompose.ts'

describe('prioritizeTaskSeeds', () => {
  it('prefers IN_PROGRESS then non-terminal; skips DONE unless empty active', () => {
    const tasks = [
      { id: 'd1', title: 'done', work_state: 'DONE' },
      { id: 'p1', title: 'pend', work_state: 'PENDING' },
      { id: 'i1', title: 'wip', work_state: 'IN_PROGRESS' },
      { id: 'b1', title: 'blocked', work_state: 'BLOCKED' },
    ] as TaskRow[]
    const seeds = prioritizeTaskSeeds(tasks, 3)
    assert.deepEqual(
      seeds.map((s) => s.id),
      ['i1', 'p1', 'b1'],
    )
  })

  it('uses DONE only when no active tasks', () => {
    const tasks = [
      { id: 'd1', title: 'done', work_state: 'DONE' },
      { id: 's1', title: 'skip', work_state: 'SKIPPED' },
    ] as TaskRow[]
    const seeds = prioritizeTaskSeeds(tasks, 2)
    assert.equal(seeds.length, 2)
    assert.equal(seeds[0].id, 'd1')
  })
})

describe('composeSeedsFromParts', () => {
  it('caps at SEED_CAP and targets fill from search', () => {
    const tasks = Array.from({ length: 2 }, (_, i) => ({
      id: `t${i}`,
      title: `T${i}`,
      work_state: 'IN_PROGRESS',
    })) as TaskRow[]
    const searchHits = Array.from({ length: 10 }, (_, i) => ({
      id: `g${i}`,
      kind: 'goal',
      title: `G${i}`,
    }))
    const seeds = composeSeedsFromParts(tasks, searchHits, {
      target: SEED_TARGET,
      cap: SEED_CAP,
    })
    assert.ok(seeds.length <= SEED_CAP)
    assert.ok(seeds.length >= Math.min(SEED_TARGET, 2 + 10))
    assert.equal(seeds.filter((s) => s.source === 'task').length, 2)
  })
})

describe('mergeOverviewGraphs', () => {
  it('keeps seeds first and trims to UI_CAP', () => {
    const seedIds = ['s0', 's1']
    const nodes = Array.from({ length: 120 }, (_, i) => ({
      id: i < 2 ? `s${i}` : `n${i}`,
      kind: 'task',
      title: `N${i}`,
    }))
    const edges = nodes.slice(2, 40).map((n) => ({
      from: 's0',
      to: n.id,
      rel: 'mentions',
    }))
    const g: BoundedGraph = {
      center: 's0',
      max_nodes: 40,
      truncated: true,
      nodes,
      edges,
    }
    const merged = mergeOverviewGraphs([g], seedIds, { uiCap: UI_CAP })
    assert.ok(merged.nodes.length <= UI_CAP)
    assert.ok(merged.nodes.some((n) => n.id === 's0'))
    assert.ok(merged.nodes.some((n) => n.id === 's1'))
    assert.ok(merged.omitted > 0)
  })
})

describe('mergeOverviewGraphs', () => {
  it('tolerates null edges from API (legacy null JSON)', () => {
    const graphs = [
      {
        center: 'a',
        max_nodes: 40,
        truncated: false,
        nodes: [{ id: 'a', kind: 'task', title: 'A' }],
        edges: null,
      },
    ] as unknown as BoundedGraph[]
    const merged = mergeOverviewGraphs(graphs, ['a'])
    assert.deepEqual(merged.edges, [])
    assert.equal(merged.nodes.length, 1)
  })
})

describe('partitionSettledGraphs', () => {
  it('keeps successes on partial fail', () => {
    const ok: BoundedGraph = {
      center: 'a',
      max_nodes: 40,
      truncated: false,
      nodes: [{ id: 'a', kind: 'task', title: 'A' }],
      edges: [],
    }
    const results: PromiseSettledResult<BoundedGraph>[] = [
      { status: 'fulfilled', value: ok },
      { status: 'rejected', reason: new Error('boom') },
    ]
    const part = partitionSettledGraphs(results, ['a', 'bad-id'])
    assert.equal(part.graphs.length, 1)
    assert.deepEqual(part.failedSeedIds, ['bad-id'])
    assert.equal(part.hardFail, false)
  })

  it('hardFail when all seeds fail', () => {
    const results: PromiseSettledResult<BoundedGraph>[] = [
      { status: 'rejected', reason: new Error('x') },
    ]
    const part = partitionSettledGraphs(results, ['bad'])
    assert.equal(part.hardFail, true)
  })
})

describe('applySeedMeta', () => {
  it('copies work_state onto matching nodes', () => {
    const nodes = applySeedMeta(
      [{ id: 't1', kind: 'task', title: 'T' }],
      [{ id: 't1', kind: 'task', title: 'T', work_state: 'IN_PROGRESS', source: 'task' }],
    )
    assert.equal(nodes[0].work_state, 'IN_PROGRESS')
  })
})
