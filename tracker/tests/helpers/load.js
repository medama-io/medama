// @ts-check
import { test } from '@playwright/test';
import { addRequestListeners, matchRequests } from './helpers';

const baseURL = 'http://localhost:8080';

/**
 * Create test block for all loading related tests.
 *
 * @param {string} testURL
 */
const loadTests = (testURL) => {
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
					responseBody: '0',
				},
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'load',
						u: testURL,
						r: '',
						p: true,
						q: true,
					},
				},
			];

			const listenerPromise = addRequestListeners(page, expectedRequests);

			await page.goto(`${baseURL}${testURL}`, { waitUntil: 'load' });
			await page.waitForLoadState();

			const data = await listenerPromise;
			await matchRequests(data, expectedRequests);
		});
	});

	test('returning visitor second load event', async ({ page }) => {
		const expectedRequests = [
			{
				method: 'POST',
				url: '/api/event/hit',
				status: 204,
				postData: {
					e: 'unload',
				},
			},
			{
				method: 'GET',
				url: '/api/event/ping',
				status: 200,
				responseBody: '1',
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
				postData: {
					e: 'load',
					u: testURL,
					r: '',
					p: false, // Returning visitor
					q: false, // Not a new page view
				},
			},
		];

		// First load, should be a new visitor
		await page.goto(`${baseURL}${testURL}`, { waitUntil: 'networkidle' });

		const listenerPromise = addRequestListeners(page, expectedRequests);

		// Refresh page for second load
		await page.reload({ waitUntil: 'load' });

		const data = await listenerPromise;
		await matchRequests(data, expectedRequests);
	});
};

export { loadTests };
