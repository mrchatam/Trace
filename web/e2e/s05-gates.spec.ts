import { expect, test, type APIRequestContext, type Page } from '@playwright/test'

async function apiJson(
  request: APIRequestContext,
  method: string,
  path: string,
  body?: unknown,
) {
  const res = await request.fetch(path, {
    method,
    headers: body !== undefined ? { 'Content-Type': 'application/json' } : undefined,
    data: body,
  })
  const text = await res.text()
  let data: unknown = undefined
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = text
    }
  }
  return { status: res.status(), data }
}

test.describe('G-promote', () => {
  test('createEntity + createLink + createTransition promote path', async ({ request }) => {
    const disc = await apiJson(request, 'POST', '/v1/entities', {
      kind: 'discovery',
      title: `e2e-disc-${Date.now()}`,
      body: 'playwright promote',
    })
    expect(disc.status).toBe(201)
    const discId = (disc.data as { id: string }).id
    expect(discId).toBeTruthy()

    const tasks = await apiJson(request, 'GET', '/v1/tasks?limit=5')
    expect(tasks.status).toBe(200)
    const goalId = ((tasks.data as { items: { goal_id?: string }[] }).items ?? []).find(
      (t) => t.goal_id,
    )?.goal_id

    const task = await apiJson(request, 'POST', '/v1/entities', {
      kind: 'task',
      title: `e2e-promoted-${Date.now()}`,
      goal_id: goalId,
      body: 'from discovery',
    })
    expect(task.status).toBe(201)
    const taskId = (task.data as { id: string }).id

    const link = await apiJson(request, 'POST', '/v1/links', {
      rel: 'discovery-mentions-task',
      from: discId,
      to: taskId,
      source_type: 'e2e',
    })
    expect(link.status).toBe(201)

    const tx = await apiJson(request, 'POST', '/v1/transitions', {
      task_id: taskId,
      to_state: 'IN_PROGRESS',
      reason: 'e2e promote transition',
      actor: 'e2e',
    })
    // 200 ok, or envelope on deny — either proves library path
    expect([200, 403, 409]).toContain(tx.status)
    if (tx.status !== 200) {
      const env = tx.data as { error?: { code?: string } }
      expect(env.error?.code).toBeTruthy()
    }
  })

  test('Discoveries UI promote surfaces success', async ({ page }) => {
    await page.goto('/discoveries')
    await expect(page.getByRole('heading', { name: 'Discoveries' })).toBeVisible()

    const title = `ui-disc-${Date.now()}`
    await page.locator('#create-title').fill(title)
    await page.locator('#create-body').fill('ui promote smoke')
    await page.getByRole('button', { name: 'Create', exact: true }).click()
    await expect(page.getByRole('status').filter({ hasText: 'Created discovery' })).toBeVisible({
      timeout: 15_000,
    })

    // Detail opens with promote form for discovery
    await expect(page.getByRole('heading', { name: 'Promote to task' })).toBeVisible()
    await page.locator('#promote-title').fill(`${title}-task`)
    // Prefer a goal when the dogfood store has one listed
    const goal = page.locator('#promote-goal')
    const opts = await goal.locator('option').count()
    if (opts > 1) await goal.selectOption({ index: 1 })
    await page.getByRole('button', { name: 'Promote', exact: true }).click()
    await expect(
      page.getByRole('status').filter({ hasText: /Promoted → task|Task created/ }),
    ).toBeVisible({ timeout: 15_000 })
  })
})

test.describe('G-export', () => {
  test('seed path escape + strict honesty via API', async ({ request }) => {
    const escape = await apiJson(request, 'POST', '/v1/seed/import', {
      input_path: '../outside.json',
    })
    expect(escape.status).toBe(400)
    const escEnv = escape.data as { error?: { code?: string; message?: string } }
    expect(escEnv.error?.code).toBe('VALIDATION_ERROR')
    expect(`${escEnv.error?.message ?? ''}`.toLowerCase()).toMatch(/escap|path|valid/)

    const strict = await apiJson(request, 'POST', '/v1/seed/export', {
      output_path: 'trace/graph-e2e.json',
      strict: true,
    })
    expect(strict.status).toBe(501)
    const stEnv = strict.data as { error?: { code?: string } }
    expect(stEnv.error?.code).toBe('NOT_IMPLEMENTED')
  })

  test('Seed UI shows envelope on path escape', async ({ page }) => {
    await page.goto('/seed')
    await expect(page.getByRole('heading', { name: 'Seed' })).toBeVisible()
    await page.locator('#import-path').fill('../outside.json')
    await page.getByRole('button', { name: 'Import…' }).click()
    await page.getByRole('button', { name: 'Import', exact: true }).click()
    await expect(page.getByRole('alert')).toContainText(/VALIDATION|escap|path/i, {
      timeout: 15_000,
    })
  })
})

test.describe('smoke nav', () => {
  test('Reviews nav + Graph home + /graph alias', async ({ page }: { page: Page }) => {
    await page.goto('/reviews')
    await expect(page.getByRole('heading', { name: 'Reviews' })).toBeVisible()

    await page.goto('/')
    await expect(page.getByRole('heading', { name: 'Graph' })).toBeVisible()
    await expect(page.locator('#graph-budget')).toBeVisible()

    await page.goto('/graph')
    await expect(page.getByRole('heading', { name: 'Graph' })).toBeVisible()
    await expect(page.locator('#graph-budget')).toBeVisible()

    await page.goto('/overview')
    await expect(page.getByRole('heading', { name: 'Overview' })).toBeVisible()
  })
})
