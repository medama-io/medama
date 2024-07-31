import { test, expect } from '@playwright/test';

test.describe('Simple', () => {
	test('Load event unique', async ({ page }) => {
		page.on('request', (request) => {
			const url = request.url();
			if (url.includes('/api/')) {
				console.log('<<', request.method(), request.url());
			}
		});
		await page.goto('/simple/index.html');
		await page.waitForLoadState();
	});
});
