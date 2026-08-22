/**
 * Node strip-types unit checks for pickActiveTask (no vitest in web/).
 * Run: node --experimental-strip-types --test src/lib/pickActiveTask.test.ts
 */
import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import type { TaskRow } from '../api/ops'
import { pickActiveTask } from './pickActiveTask.ts'

/** Narrative feet-seller fixtures (synthetic lists only — never mutate dogfood). */
const STEP1 = '33247e2d-aa10-4b25-b194-4b7afb5a6359'
const LOOP112 = '99d8fb92-65ac-462c-82c4-21bcf198c09e'

function row(id: string, work_state: string, title = id): TaskRow {
  return { id, title, work_state }
}

describe('pickActiveTask', () => {
  it('allDone_picksLast — oldest-first all DONE → last id; ≠ Step1', () => {
    const tasks = [
      row(STEP1, 'DONE', 'Step 1'),
      row('mid-a', 'DONE'),
      row('mid-b', 'DONE'),
      row(LOOP112, 'DONE', 'Loop 112'),
    ]
    const pick = pickActiveTask(tasks)
    assert.ok(pick)
    assert.equal(pick.id, LOOP112)
    assert.notEqual(pick.id, STEP1)
  })

  it('inProgress_wins', () => {
    const tasks = [
      row('d1', 'DONE'),
      row('p1', 'PENDING'),
      row('i1', 'IN_PROGRESS'),
      row('i2', 'IN_PROGRESS'),
    ]
    assert.equal(pickActiveTask(tasks)?.id, 'i1')
  })

  it('nonTerminal_beforeDone', () => {
    const tasks = [
      row('d1', 'DONE'),
      row('s1', 'SKIPPED'),
      row('p1', 'PENDING'),
      row('b1', 'BLOCKED'),
      row('d2', 'DONE'),
    ]
    assert.equal(pickActiveTask(tasks)?.id, 'p1')
  })

  it('empty_returnsNull', () => {
    assert.equal(pickActiveTask([]), null)
  })

  it('truncatedList_honesty — full synthetic >100 → P3a last (not [0])', () => {
    const tasks: TaskRow[] = Array.from({ length: 123 }, (_, i) => {
      if (i === 0) return row(STEP1, 'DONE', 'Step 1')
      if (i === 122) return row(LOOP112, 'DONE', 'Loop 112')
      return row(`t-${i}`, 'DONE')
    })
    // Caller contract: feed the complete list (listTasksForPick), not a truncated page.
    const truncatedPage = tasks.slice(0, 100)
    assert.notEqual(truncatedPage[truncatedPage.length - 1]?.id, LOOP112)

    const pick = pickActiveTask(tasks)
    assert.ok(pick)
    assert.equal(pick.id, LOOP112)
    assert.notEqual(pick.id, STEP1)
    // Silent [0] fallback would wrongly bind Step1 on all-DONE.
    assert.notEqual(pick.id, truncatedPage[0]?.id)
  })

  it('red_to_green_seed — INVESTIGATION assert defaultPick ≠ Step1 on all-DONE oldest-first', () => {
    const tasks = [row(STEP1, 'DONE', 'Step 1'), row(LOOP112, 'DONE', 'Loop 112')]
    const defaultPick = pickActiveTask(tasks)
    assert.ok(defaultPick)
    assert.notEqual(defaultPick.id, STEP1)
    assert.equal(defaultPick.id, LOOP112)
  })
})
