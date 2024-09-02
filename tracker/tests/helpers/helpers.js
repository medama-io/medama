// @ts-check
import { expect, test } from '@playwright/test';
import { loadUnloadTests } from './load-unload';
import { clickEventTests } from './click-events';

/**
 * @typedef {('simple'|'history')} Tests
 */

/**
 * @typedef PostData
 * @property {string} e - Event type
 * Load Events
 * @property {string=} u - URL
 * @property {string=} r - Referrer
 * @property {boolean=} p - Unique visitor
 * @property {boolean=} q - Unique page view
 * Unload Events
 * @property {string=} m - Time spent on page
 * Custom Events
 * @property {string=} g - Custom event groupname
 * @property {Object=} d - Custom event properties
 */

/**
 * @typedef ExpectedRequest
 * @property {string} method
 * @property {string} url
 * @property {number} status
 * @property {PostData=} postData
 * @property {string=} responseBody
 * @property {Array<string>=} ignoreBrowsers
 */

/**
 * @typedef {Object} RequestResponsePair
 * @property {import('@playwright/test').Request} request
 * @property {import('@playwright/test').Response} response
 */

const TIMEOUT_DELAY = 1250;

/**
 * Get the browser name from the page context.
 *
 * @param {import('@playwright/test').Page} page
 * @returns {string}
 */
const getBrowser = (page) => page.context().browser().browserType().name();

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
		}, 4000); // 4 second timeout

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
 * @param {import('@playwright/test').Page} page
 * @param {Promise<RequestResponsePair[]>} responses
 * @param {Array<ExpectedRequest>} expectedRequests
 * @returns {Promise<void>}
 */
const matchRequests = async (page, responses, expectedRequests) => {
	const data = await responses;

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
			} else if (request.method() === 'POST') {
				// Chromium sometimes doesn't return a response for sendBeacon but the others do (???)
				status = 204;
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

	// Get browser name to ignore certain requests]
	const browserName = getBrowser(page);

	const expected = expectedRequests
		.map((req) => {
			if (req.ignoreBrowsers && req.ignoreBrowsers.includes(browserName)) {
				return null;
			}

			return {
				method: req.method,
				url: req.url,
				status: req.status,
				postData: req.method === 'POST' ? req.postData : undefined,
				responseBody: req.status !== 204 ? req.responseBody : undefined,
			};
		})
		.filter((req) => req !== null);

	expect.soft(actualRequests).toEqual(
		expected.map((exp) => ({
			method: exp.method,
			url: expect.stringContaining(exp.url),
			status: exp.status,
			postData: exp.postData
				? expect.objectContaining({
						...exp.postData,
						...(exp.postData.u
							? { u: expect.stringContaining(exp.postData.u) }
							: {}),
					})
				: undefined,
			responseBody: exp.responseBody,
		})),
	);

	expect(actualRequests.length).toBe(expected.length);
};

/**
 * Create relative URL for given test.
 *
 * @param {Tests} name
 * @param {string} path
 * @param {boolean=} relative
 * @returns {string}
 */
const createURL = (name, path, relative = true) =>
	`${relative ? '.' : ''}/${name}/${path}`;

/**
 * Create test block for given URL.
 *
 * @param {Tests} name
 */
const createTests = (name) => {
	test.describe(name + ' load + unload tests', () => {
		loadUnloadTests(name);
	});

	test.describe(name + ' click event tests', () => {
		clickEventTests(name);
	});
};

export {
	addRequestListeners,
	createTests,
	createURL,
	getBrowser,
	matchRequests,
	TIMEOUT_DELAY,
};
