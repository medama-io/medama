// @ts-check
import { test } from '@playwright/test';
import {
	addRequestListeners,
	createURL,
	matchRequests,
	TIMEOUT_DELAY,
} from './helpers';

/**
 * Create test block for all page loading related tests.
 *
 * @param {import('./helpers').Tests} name
 */
const pageTests = (name) => {
	test.describe('load', () => {
		test('unique visitor load event', async ({ page }) => {
			const expectedRequests = [
				{
					method: 'GET',
					url: '/api/event/ping?root',
					status: 200,
					responseBody: '0',
				},
				{
					method: 'GET',
					url: '/api/event/ping?u',
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

			await matchRequests(page, listeners, expectedRequests);
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
					ignoreBrowsers: ['firefox'], // TODO: Investigate why this request is not sent in Firefox (works outside Playwright)
				},
				{
					method: 'GET',
					url: '/api/event/ping?root',
					status: 200,
					responseBody: '1',
				},
				{
					method: 'GET',
					url: '/api/event/ping?u',
					status: 200,
					responseBody: '1',
				},
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'load',
						u: createURL(name, 'index.html', false),
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
			await page.reload({ waitUntil: 'networkidle' });

			await matchRequests(page, listeners, expectedRequests);
		});

		test('unique visitor navigates to a new page', async ({ page }) => {
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
					url: '/api/event/ping?u',
					status: 200,
					responseBody: '0',
				},
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'load',
						u: createURL(name, 'about', false),
						p: false, // Returning visitor
						q: true, // New page view
					},
				},
			];

			// First load, should be a new visitor
			await page.goto(createURL(name, 'index.html'), {
				waitUntil: 'networkidle',
			});

			const listeners = addRequestListeners(page, expectedRequests);

			// Navigate to about page using proper routing
			await page.getByRole('link', { name: 'About' }).click();
			await page.waitForTimeout(TIMEOUT_DELAY);
			await page.waitForLoadState('networkidle');

			await matchRequests(page, listeners, expectedRequests);
		});

		test('returning visitor navigates to visited page', async ({ page }) => {
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
					url: '/api/event/ping?u',
					status: 200,
					responseBody: '1',
				},
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'load',
						u: createURL(name, 'index', false),
						p: false, // Returning visitor
						q: false, // Returning page view
					},
					ignoreBrowsers: ['webkit'],
				},
				/**
				 * @note WebKit browsers do not send If-Modified-Since headers for MPA websites
				 * leading to p=true in this test.
				 * @see https://stackoverflow.com/a/75944210
				 */
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'load',
						u: createURL(name, 'index', false),
						p: name == 'simple' ? true : false, // Returning visitor
						q: false, // Returning page view
					},
					ignoreBrowsers: ['firefox', 'chrome', 'msedge', 'chromium'],
				},
			];

			// First load, should be a new visitor
			await page.goto(createURL(name, 'index.html'), {
				waitUntil: 'networkidle',
			});

			// Navigate to about page to cache visit
			await page.getByRole('link', { name: 'About' }).click();
			await page.waitForTimeout(TIMEOUT_DELAY);
			await page.waitForLoadState('networkidle');

			const listeners = addRequestListeners(page, expectedRequests);

			// Navigate back to home page to test returning visitor
			await page.getByRole('link', { name: 'Home' }).click();
			await page.waitForLoadState('networkidle');

			await matchRequests(page, listeners, expectedRequests);
		});

		test('returning visitors uses back button to visited page', async ({
			page,
		}) => {
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
					url: '/api/event/ping?u',
					status: 200,
					responseBody: '1',
				},
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'load',
						u: createURL(name, 'index', false),
						p: false, // Returning visitor
						q: false, // Returning page view
					},
					ignoreBrowsers: ['webkit'],
				},
				/**
				 * @note WebKit browsers do not send If-Modified-Since headers for MPA websites
				 * leading to p=true in this test.
				 * @see https://stackoverflow.com/a/75944210
				 */
				{
					method: 'POST',
					url: '/api/event/hit',
					status: 204,
					postData: {
						e: 'load',
						u: createURL(name, 'index', false),
						p: name == 'simple' ? true : false, // Returning visitor
						q: false, // Returning page view
					},
					ignoreBrowsers: ['firefox', 'chrome', 'msedge', 'chromium'],
				},
			];

			// First load, should be a new visitor
			await page.goto(createURL(name, 'index.html'), {
				waitUntil: 'networkidle',
			});

			// Navigate to about page to cache visit
			await page.getByRole('link', { name: 'About' }).click();
			await page.waitForTimeout(TIMEOUT_DELAY);
			await page.waitForLoadState('networkidle');

			const listeners = addRequestListeners(page, expectedRequests);

			// Navigate back to home page to test returning visitor
			await page.goBack();
			await page.waitForLoadState('networkidle');

			await matchRequests(page, listeners, expectedRequests);
		});
	});
};

export { pageTests };
