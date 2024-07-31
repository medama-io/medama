const { expect, test } = require('@playwright/test');

const baseURL = 'http://localhost:8080';

/**
 * Helper function to track API calls and verify their responses.
 *
 * @param {import('@playwright/test').Page} page
 * @param {Array<{method: string, url: string, status: number, responseBody: string}>} expectedRequests
 * @param {Object} postData
 * @returns {Promise<void>}
 */
const runApiTest = async (page, expectedRequests, postData) => {
	const requestPromises = expectedRequests.map(async (expectedRequest) => {
		const [request, response] = await Promise.all([
			page.waitForRequest(
				(req) =>
					req.url().includes(expectedRequest.url) &&
					req.method() === expectedRequest.method,
			),
			page.waitForResponse(
				(res) =>
					res.url().includes(expectedRequest.url) &&
					res.status() === expectedRequest.status,
			),
		]);

		console.log('>>', request.method(), request.url());
		console.log('<<', response.status(), response.url());

		if (expectedRequest.method === 'POST') {
			expect(request.postDataJSON()).toMatchObject(postData);
		}

		if (response.status() === 200) {
			const responseText = await response.text();
			expect(responseText).toBe(expectedRequest.responseBody);
		} else if (response.status() !== 204) {
			throw new Error('Unexpected request status');
		}

		return Promise.all(requestPromises);
	});
};

/**
 * Create test block for given URL.
 *
 * @param {string} testUrl
 */
const createTests = (testUrl) => {
	test.describe(testUrl, () => {
		test('Unique visitor load event', async ({ page }) => {
			const apiRequests = [
				{ method: 'GET', url: '/api/', status: 200, responseBody: '0' },
				{ method: 'GET', url: '/api/', status: 200, responseBody: '0' },
				{ method: 'POST', url: '/api/', status: 204, responseBody: '0' },
			];

			const postData = {
				e: 'load',
				u: `${baseURL}${testUrl}`,
				r: '',
				p: true,
				q: true,
			};

			await runApiTest(page, apiRequests, postData);

			await page.goto(`${baseURL}${testUrl}`);
			await page.waitForLoadState();
		});
	});
};

export { createTests };
