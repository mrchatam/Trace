import { defineConfig, devices } from '@playwright/test'

/**
 * E2E against `trace serve` (default 127.0.0.1:7432) serving web/dist.
 * Start server externally or via web/scripts/e2e-with-serve.sh.
 */
export default defineConfig({
  testDir: './e2e',
  timeout: 60_000,
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: 'list',
  use: {
    baseURL: process.env.TRACE_E2E_BASE ?? 'http://127.0.0.1:7432',
    ...devices['Desktop Chrome'],
    trace: 'on-first-retry',
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
})
