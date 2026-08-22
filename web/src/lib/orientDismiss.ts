export const ORIENT_DISMISS_KEY = 'trace.orient.dismissed'

export function isOrientDismissed(): boolean {
  try {
    return localStorage.getItem(ORIENT_DISMISS_KEY) === '1'
  } catch {
    return false
  }
}

export function dismissOrient(): void {
  try {
    localStorage.setItem(ORIENT_DISMISS_KEY, '1')
  } catch {
    // Storage unavailable — dismiss still hides for this session via parent state.
  }
}
