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

/** Ensure a task has at least one neighbor for select≠expand smoke. */
async function ensureNeighborhood(request: APIRequestContext): Promise<string> {
  const tasks = await apiJson(request, 'GET', '/v1/tasks?limit=20')
  expect(tasks.status).toBe(200)
  const items = (tasks.data as { items: { id: string; goal_id?: string }[] }).items ?? []
  expect(items.length).toBeGreaterThan(0)

  for (const t of items.slice(0, 8)) {
    const g = await apiJson(
      request,
      'GET',
      `/v1/graph?center=${encodeURIComponent(t.id)}&max_nodes=50&depth=2`,
    )
    if (g.status === 200) {
      const nodes = (g.data as { nodes: { id: string }[] }).nodes ?? []
      if (nodes.length >= 2) return t.id
    }
  }

  const seed = items[0]
  const disc = await apiJson(request, 'POST', '/v1/entities', {
    kind: 'discovery',
    title: `e2e-depth-neighbor-${Date.now()}`,
    body: 'select≠expand neighbor',
  })
  expect(disc.status).toBe(201)
  const discId = (disc.data as { id: string }).id
  const link = await apiJson(request, 'POST', '/v1/links', {
    rel: 'discovery-mentions-task',
    from: discId,
    to: seed.id,
    source_type: 'e2e',
  })
  expect(link.status).toBe(201)
  return seed.id
}

test.describe('S03 depth', () => {
  test('overview-on-open shows canvas without pick-first; select≠expand', async ({
    page,
    request,
  }: {
    page: Page
    request: APIRequestContext
  }) => {
    const taskId = await ensureNeighborhood(request)

    await page.goto('/')
    await expect(page.getByRole('heading', { name: 'Explore' })).toBeVisible()

    // Project graph canvas appears without mandatory pick-task click
    await expect(page.getByTestId('graph-overview-loading')).toBeHidden({ timeout: 30_000 })
    await expect(page.getByTestId('graph-canvas')).toBeVisible({ timeout: 30_000 })
    await expect(page.locator('#graph-budget')).toBeVisible()
    await expect(page.getByTestId('graph-center-id')).not.toHaveText('')

    const canvasNodes = page.locator('[data-testid^="graph-canvas-node-"]')
    await expect(canvasNodes.first()).toBeVisible()
    expect(await canvasNodes.count()).toBeGreaterThanOrEqual(1)
    await expect(canvasNodes.first()).toHaveAttribute('data-kind', /.+/ )

    // Manual center lives on Orient
    await page.goto('/orient')
    await expect(page.getByTestId('graph-manual-center')).toBeVisible()

    // Drill into neighborhood via Orient task pick
    await page.getByTestId(`graph-pick-task-${taskId}`).click()
    await expect(page.getByTestId('graph-center-id')).toHaveText(taskId, { timeout: 15_000 })
    const centerBefore = taskId

    await expect(page.getByTestId('graph-canvas')).toBeVisible()
    expect(await canvasNodes.count()).toBeGreaterThanOrEqual(2)

    let other: string | null = null
    const count = await canvasNodes.count()
    for (let i = 0; i < count; i++) {
      const testId = await canvasNodes.nth(i).getAttribute('data-testid')
      const id = testId?.replace('graph-canvas-node-', '') ?? ''
      if (id && id !== centerBefore) {
        other = id
        await canvasNodes.nth(i).click()
        break
      }
    }
    expect(other).toBeTruthy()

    await expect(page.getByTestId('inspector-selected-id')).toHaveText(other!, { timeout: 15_000 })
    await expect(page.locator('[data-section="summary"]')).toBeVisible()
    // Select must not re-center
    await expect(page.getByTestId('graph-center-id')).toHaveText(centerBefore)

    // Canvas click uses the same onSelect path
    await page.getByTestId(`graph-canvas-node-${other}`).click()
    await expect(page.getByTestId('inspector-selected-id')).toHaveText(other!, { timeout: 15_000 })
    await expect(page.getByTestId('graph-center-id')).toHaveText(centerBefore)

    await page.getByTestId('inspector-use-center').click()
    await expect(page.getByTestId('graph-center-id')).toHaveText(other!, { timeout: 15_000 })
  })
})
