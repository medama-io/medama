// @ts-check

import { test, expect } from '@playwright/test';

test.describe('data-api', () => {
	test('data-api connects to alternative api', async ({ page }) => {
		const requests = [];
		await page.on('request', (request) => {
			console.log('>>', request.method(), request.url());
			if (request.url().includes('/api/event')) requests.push(request.url());
		});

		await page.goto('/data-api/index.html', { waitUntil: 'load' });
		await page.waitForLoadState();

		expect(requests.length).toBeGreaterThan(0);
		for (const url of requests) {
			expect(url).toContain('example.com');
		}
	});
});
