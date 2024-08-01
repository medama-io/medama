// @ts-check
import { expect, test } from '@playwright/test';
import { loadTests } from './load';

/**
 * @typedef ExpectedRequest
 * @property {string} method
 * @property {string} url
 * @property {number} status
 * @property {string=} responseBody
 */

/**
 * @typedef {Object} RequestResponsePair
 * @property {import('@playwright/test').Request} request
 * @property {import('@playwright/test').Response} response
 */

/**
 * @typedef PostData
 * @property {string} e - Event type
 * @property {string} u - URL
 * @property {string} r - Referrer
 * @property {boolean} p - Unique visitor
 * @property {boolean} q - Unique page view
 */

const baseURL = 'http://localhost:8080';

/**
 * Add request and response listeners to the page to track API calls.
 *
 * @param {import('@playwright/test').Page} page
 * @param {Array<ExpectedRequest>} expectedRequests
 * @returns {Promise<RequestResponsePair[]>}
 */
const addRequestListeners = (page, expectedRequests) => {
	/** @type {RequestResponsePair[]} */
	const pairs = [];

	return new Promise((resolve) => {
		let matchedCount = 0;
		const timeoutId = setTimeout(() => {
			console.warn(
				'Timeout waiting for requests. Resolving with partial data.',
			);
			resolve(pairs);
		}, 5000); // 5 second timeout

		page.on('request', (request) => {
			const matchingExpected = expectedRequests.find(
				(ereq) =>
					request.url().includes(ereq.url) && request.method() === ereq.method,
			);
			if (matchingExpected) {
				console.log('>>', request.method(), request.url());
				pairs.push({ request, response: null });
			}
		});

		page.on('response', (response) => {
			const matchingExpected = expectedRequests.find(
				(ereq) =>
					response.url().includes(ereq.url) &&
					response.status() === ereq.status,
			);
			if (matchingExpected) {
				console.log('<<', response.status(), response.url());
				const pair = pairs.find(
					(p) => !p.response && p.request.url() === response.url(),
				);
				if (pair) {
					pair.response = response;
					matchedCount++;
					if (matchedCount === expectedRequests.length) {
						clearTimeout(timeoutId);
						resolve(pairs);
					}
				}
			}
		});
	});
};

/**
 * After navigating to a page, wait for the API calls to complete before matching the expected requests.
 *
 * @param {RequestResponsePair[]} data
 * @param {Array<ExpectedRequest>} expectedRequests
 * @param {PostData} postData
 * @returns {Promise<void>}
 */
const matchRequests = async (data, expectedRequests, postData) => {
	for (const [index, pair] of data.entries()) {
		const { request, response } = pair;
		const expected = expectedRequests[index];

		expect(request.method()).toBe(expected.method);
		expect(request.url()).toContain(expected.url);
		expect(response.status()).toBe(expected.status);

		if (expected.method === 'POST') {
			expect(request.postDataJSON()).toMatchObject(postData);
		}

		// No response body for 204 responses
		if (expected.status !== 204 && expected.responseBody !== undefined) {
			const responseText = await response.text();
			expect(responseText).toBe(expected.responseBody);
		}
	}
};

/**
 * Create test block for given URL.
 *
 * @param {string} testURL
 */
const createTests = (testURL) => {
	test.describe(testURL, () => {
		loadTests(baseURL, testURL);
	});
};

export { createTests, addRequestListeners, matchRequests };
