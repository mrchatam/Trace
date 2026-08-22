/**
 * Node strip-types unit checks for graph navigation helpers.
 * Run: node --experimental-strip-types --test src/lib/graphNavigate.test.ts
 */
import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { buildGraphHref, parseGraphSearchParams } from './graphNavigate.ts'

describe('buildGraphHref', () => {
  it('returns / when center is empty', () => {
    assert.equal(buildGraphHref(), '/')
    assert.equal(buildGraphHref('  '), '/')
  })

  it('includes center and neighborhood mode', () => {
    assert.equal(buildGraphHref('task-1'), '/?center=task-1&mode=neighborhood')
  })

  it('encodes center id', () => {
    assert.equal(buildGraphHref('a b'), '/?center=a+b&mode=neighborhood')
  })
})

describe('parseGraphSearchParams', () => {
  it('defaults to project with no params', () => {
    assert.deepEqual(parseGraphSearchParams(new URLSearchParams()), {
      center: null,
      mode: 'project',
    })
  })

  it('center param forces neighborhood mode', () => {
    assert.deepEqual(parseGraphSearchParams(new URLSearchParams('center=abc')), {
      center: 'abc',
      mode: 'neighborhood',
    })
  })

  it('ignores blank center', () => {
    assert.deepEqual(parseGraphSearchParams(new URLSearchParams('center=')), {
      center: null,
      mode: 'project',
    })
  })
})
