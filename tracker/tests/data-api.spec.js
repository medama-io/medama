// @ts-check

import { expect, test } from '@playwright/test';

test.describe('data-api', () => {
	test('data-api connects to alternative api', async ({ page }) => {
		const requestPromise = page.waitForRequest((request) =>
			request.url().includes('/api/event'),
		);

		await page.goto('/data-api/index.html', { waitUntil: 'load' });

		const request = await requestPromise;
		expect(request.url()).toContain('example.com');
	});
});
