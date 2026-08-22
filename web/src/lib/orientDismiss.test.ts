/**
 * Node strip-types unit checks for orient dismiss helpers.
 * Run: node --experimental-strip-types --test src/lib/orientDismiss.test.ts
 */
import assert from 'node:assert/strict'
import { afterEach, describe, it } from 'node:test'
import {
  ORIENT_DISMISS_KEY,
  dismissOrient,
  isOrientDismissed,
} from './orientDismiss.ts'

const storage = new Map<string, string>()

function installLocalStorageMock() {
  Object.defineProperty(globalThis, 'localStorage', {
    configurable: true,
    value: {
      getItem(key: string) {
        return storage.has(key) ? storage.get(key)! : null
      },
      setItem(key: string, value: string) {
        storage.set(key, value)
      },
      removeItem(key: string) {
        storage.delete(key)
      },
      clear() {
        storage.clear()
      },
    },
  })
}

describe('orientDismiss helpers', () => {
  afterEach(() => {
    storage.clear()
  })

  it('isOrientDismissed returns false when key absent', () => {
    installLocalStorageMock()
    assert.equal(isOrientDismissed(), false)
  })

  it('dismissOrient persists flag and isOrientDismissed reads it', () => {
    installLocalStorageMock()
    dismissOrient()
    assert.equal(localStorage.getItem(ORIENT_DISMISS_KEY), '1')
    assert.equal(isOrientDismissed(), true)
  })

  it('isOrientDismissed returns false when storage throws', () => {
    Object.defineProperty(globalThis, 'localStorage', {
      configurable: true,
      value: {
        getItem() {
          throw new Error('blocked')
        },
        setItem() {
          throw new Error('blocked')
        },
      },
    })
    assert.equal(isOrientDismissed(), false)
    dismissOrient()
    assert.equal(isOrientDismissed(), false)
  })
})
