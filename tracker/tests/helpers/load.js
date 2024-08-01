// @ts-check
import { test } from '@playwright/test';
import { addRequestListeners, matchRequests } from './helpers';

/**
 * Create test block for all loading related tests.
 *
 * @param {string} baseURL
 * @param {string} testURL
 */
const loadTests = (baseURL, testURL) => {
	test.describe(testURL, () => {
		test('unique visitor load event', async ({ page }) => {
			const expectedRequests = [
				{
					method: 'GET',
					url: '/api/event/ping',
					status: 200,
					responseBody: '0',
				},
				{
					method: 'GET',
					url: '/api/event/ping',
					status: 200,
					responseBody: '1',
				},
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
				},
			];

			const postData = {
				e: 'load',
				u: `${baseURL}${testURL}`,
				r: '',
				p: true,
				q: true,
			};

			const listenerPromise = addRequestListeners(page, expectedRequests);

			await page.goto(`${baseURL}${testURL}`);
			await page.waitForLoadState();

			const data = await listenerPromise;
			await matchRequests(data, expectedRequests, postData);
		});
	});
};

export { loadTests };
