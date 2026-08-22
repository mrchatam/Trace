import type { TaskRow } from '../api/ops'

/**
 * Shared GUI active-task pick (P1 → P2 → P3a).
 * Callers that auto-pick must feed a complete list (see listTasksForPick).
 */
export function pickActiveTask(tasks: TaskRow[]): TaskRow | null {
  if (tasks.length === 0) return null
  const inProgress = tasks.find((t) => t.work_state === 'IN_PROGRESS')
  if (inProgress) return inProgress
  const nonTerminal = tasks.find(
    (t) => t.work_state !== 'DONE' && t.work_state !== 'SKIPPED',
  )
  if (nonTerminal) return nonTerminal
  return tasks[tasks.length - 1] ?? null
}
