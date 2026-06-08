import { test, expect } from '@playwright/test';

test.describe('Onboarding Flow', () => {
  test('displays welcome screen with brand and features', async ({ page }) => {
    await page.goto('/onboarding');

    // Step 0: Welcome
    await expect(page.getByText('欢迎使用 U-Hermes')).toBeVisible();
    await expect(page.getByText('即插即用')).toBeVisible();
    await expect(page.getByText('数据本地')).toBeVisible();
    await expect(page.getByText('越用越聪明')).toBeVisible();
    await expect(page.getByText('开始配置 →')).toBeVisible();
  });

  test('navigates to model selection step', async ({ page }) => {
    await page.goto('/onboarding');

    // Click through to step 1
    await page.getByText('开始配置 →').click();

    // Step 1: Model selection
    await expect(page.getByText('选择 AI 模型')).toBeVisible();
    await expect(page.getByPlaceholder('sk-...')).toBeVisible();
    await expect(page.getByText('测试连接 →')).toBeDisabled(); // empty API key
  });

  test('enables test button when API key is entered', async ({ page }) => {
    await page.goto('/onboarding');
    await page.getByText('开始配置 →').click();

    const testBtn = page.getByText('测试连接 →');
    await expect(testBtn).toBeDisabled();

    await page.getByPlaceholder('sk-...').fill('sk-test-key-123');
    await expect(testBtn).toBeEnabled();
  });

  test('can navigate back to welcome from model step', async ({ page }) => {
    await page.goto('/onboarding');
    await page.getByText('开始配置 →').click();

    await page.getByText('← 上一步').click();

    await expect(page.getByText('欢迎使用 U-Hermes')).toBeVisible();
  });
});

test.describe('Chat Page', () => {
  test('redirects to onboarding when no config', async ({ page }) => {
    await page.goto('/chat');

    // Without config, should redirect to onboarding
    await expect(page).toHaveURL(/onboarding/);
  });
});

test.describe('Settings Page', () => {
  test('shows settings sections', async ({ page }) => {
    await page.goto('/settings');

    await expect(page.getByText('系统提示词')).toBeVisible();
    await expect(page.getByRole('heading', { name: '⚠️ 恢复出厂设置' })).toBeVisible();
  });

  test('reset section shows warning', async ({ page }) => {
    await page.goto('/settings');

    await expect(page.getByText(/不可撤销|删除所有/)).toBeVisible();
  });
});

test.describe('Health Check', () => {
  test('API health endpoint responds', async ({ request }) => {
    const resp = await request.get('/api/health');
    expect(resp.status()).toBe(200);

    const body = await resp.json();
    expect(body).toHaveProperty('status', 'ok');
    expect(body).toHaveProperty('version');
    expect(body).toHaveProperty('configured');
  });
});
