import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  retries: 0,
  use: {
    baseURL: 'http://127.0.0.1:21475',
    screenshot: 'only-on-failure',
    channel: 'msedge', // Use system Edge (Windows 11 built-in)
  },
  // No webServer: start the Go server manually before running tests
});
