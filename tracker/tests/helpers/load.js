// @ts-check
import { test } from '@playwright/test';
import { addRequestListeners, createURL, matchRequests } from './helpers';

/**
 * Create test block for all loading related tests.
 *
 * @param {import('./helpers').Tests} name
 */
const loadTests = (name) => {
	test.describe('load', () => {
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
						u: createURL(name, 'index.html', false),
						r: '',
						p: true,
						q: true,
					},
				},
			];

			const listeners = addRequestListeners(page, expectedRequests);

			await page.goto(createURL(name, 'index.html'), { waitUntil: 'load' });
			await page.waitForLoadState();

			await matchRequests(listeners, expectedRequests);
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
					u: '/index.html',
					r: '',
					p: false, // Returning visitor
					q: false, // Not a new page view
				},
			},
		];

		// First load, should be a new visitor
		await page.goto(createURL(name, 'index.html'), {
			waitUntil: 'networkidle',
		});

		const listeners = addRequestListeners(page, expectedRequests);

		// Refresh page for second load
		await page.reload({ waitUntil: 'load' });

		await matchRequests(listeners, expectedRequests);
	});
};

export { loadTests };
