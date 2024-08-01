// @ts-check
import { expect, test } from '@playwright/test';
import { loadTests } from './load';

/**
 * @typedef PostData
 * @property {string} e - Event type
 * @property {string=} u - URL
 * @property {string=} r - Referrer
 * @property {boolean=} p - Unique visitor
 * @property {boolean=} q - Unique page view
 */

/**
 * @typedef ExpectedRequest
 * @property {string} method
 * @property {string} url
 * @property {number} status
 * @property {PostData=} postData
 * @property {string=} responseBody
 */

/**
 * @typedef {Object} RequestResponsePair
 * @property {import('@playwright/test').Request} request
 * @property {import('@playwright/test').Response} response
 */

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
 * @returns {Promise<void>}
 */
const matchRequests = async (data, expectedRequests) => {
	// Wait for all requests to complete before mapping them into a format that can be compared.
	const actualRequests = await Promise.all(
		data.map(async (pair) => {
			const { request, response } = pair;
			let responseBody = undefined;
			let status = undefined;

			if (response) {
				status = response.status();
				if (status !== 204) {
					try {
						const responseBuffer = await response.body();
						responseBody = responseBuffer ? await response.text() : null;
					} catch (error) {
						console.warn(`Failed to read response body: ${error.message}`);
						responseBody = null;
					}
				}
			} else {
				console.warn(
					`No response for request: ${request.method()} ${request.url()}`,
				);
			}

			return {
				method: request.method(),
				url: request.url(),
				status: status,
				postData:
					request.method() === 'POST' ? request.postDataJSON() : undefined,
				responseBody: responseBody,
			};
		}),
	);

	const expected = expectedRequests.map((req) => ({
		method: req.method,
		url: req.url,
		status: req.status,
		postData: req.method === 'POST' ? req.postData : undefined,
		responseBody: req.status !== 204 ? req.responseBody : undefined,
	}));

	expect.soft(actualRequests).toEqual(
		expected.map((exp) => ({
			method: exp.method,
			url: expect.stringContaining(exp.url),
			status: exp.status,
			postData: exp.postData
				? expect.objectContaining({
						...exp.postData,
						u: exp.postData.u
							? expect.stringContaining(exp.postData.u)
							: undefined,
					})
				: undefined,
			responseBody: exp.responseBody,
		})),
	);

	expect(actualRequests.length).toBe(expectedRequests.length);
};

/**
 * Create test block for given URL.
 *
 * @param {string} testUrl
 */
const createTests = (testUrl) => {
	test.describe(testUrl, () => {
		loadTests(testUrl);
	});
};

export { createTests, addRequestListeners, matchRequests };
