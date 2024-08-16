// @ts-check
import { test } from '@playwright/test';
import { addRequestListeners, createURL, matchRequests } from './helpers';

/**
 * Create test block for all tagged events related tests.
 *
 * @param {import('./helpers').Tests} name
 */
const taggedEventTests = (name) => {
	test.describe('button click', () => {
		test('click/auxclick event with data attribute', async ({ page }) => {
			const expectedRequests = [
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'custom',
						d: {
							foo: 'bar',
						},
					},
				},
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'custom',
						d: {
							foo: 'bar',
						},
					},
				},
			];

			await page.goto(createURL(name, 'index.html'), { waitUntil: 'load' });
			await page.waitForLoadState();

			const listeners = addRequestListeners(page, expectedRequests);

			await page.getByTestId('tagged-left-click').click({ button: 'left' });
			await page.getByTestId('tagged-middle-click').click({ button: 'middle' });
			// Should not send event on right click
			await page.getByTestId('tagged-right-click').click({ button: 'right' });

			await matchRequests(page, listeners, expectedRequests);
		});
	});
};

export { taggedEventTests };
