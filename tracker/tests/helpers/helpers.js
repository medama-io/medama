// @ts-check
import { expect, test } from '@playwright/test';
import { clickEventTests } from './click-events';
import { loadUnloadTests } from './load-unload';

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
 * @property {import('@playwright/test').Response|null} response
 */

const TIMEOUT_DELAY = 1250;
const REQUEST_TIMEOUT = 10_000;
const RESPONSE_TIMEOUT = 5_000;

/**
 * Get the browser name from the page context.
 *
 * @param {import('@playwright/test').Page} page
 * @returns {string}
 */
const getBrowser = (page) => page.context().browser().browserType().name();

/**
 * Filter browser-specific expected requests.
 *
 * @param {import('@playwright/test').Page} page
 * @param {Array<ExpectedRequest>} expectedRequests
 * @returns {Array<ExpectedRequest>}
 */
const getExpectedRequests = (page, expectedRequests) => {
	const browserName = getBrowser(page);
	return expectedRequests.filter(
		(req) => !req.ignoreBrowsers?.includes(browserName),
	);
};

/**
 * @param {string} requestURL
 * @returns {string}
 */
const getPathWithSearch = (requestURL) => {
	const url = new URL(requestURL);
	return url.pathname + url.search;
};

/**
 * @param {import('@playwright/test').Request} request
 * @param {ExpectedRequest} expectedRequest
 * @returns {boolean}
 */
const matchesExpectedRequest = (request, expectedRequest) => {
	if (request.method() !== expectedRequest.method) {
		return false;
	}

	const pathWithSearch = getPathWithSearch(request.url());
	return expectedRequest.url.includes('?')
		? pathWithSearch.startsWith(expectedRequest.url)
		: pathWithSearch === expectedRequest.url;
};

/**
 * @param {import('@playwright/test').Request} request
 * @param {ExpectedRequest} expectedRequest
 * @returns {Promise<RequestResponsePair>}
 */
const collectRequestResponse = async (request, expectedRequest) => {
	if (expectedRequest.method === 'POST' && expectedRequest.status === 204) {
		return { request, response: request.existingResponse() };
	}

	let timeoutID;
	const timeout = new Promise((resolve) => {
		timeoutID = setTimeout(() => resolve(null), RESPONSE_TIMEOUT);
	});

	try {
		const response = await Promise.race([request.response(), timeout]);
		clearTimeout(timeoutID);
		return { request, response };
	} catch (error) {
		console.warn(`Failed to collect response: ${error.message}`);
		return { request, response: null };
	}
};

/**
 * Add request and response listeners to the page to track API calls.
 *
 * @param {import('@playwright/test').Page} page
 * @param {Array<ExpectedRequest>} expectedRequests
 * @returns {Promise<RequestResponsePair[]>}
 */
const addRequestListeners = (page, expectedRequests) => {
	const requestsToMatch = getExpectedRequests(page, expectedRequests);

	return new Promise((resolve) => {
		/** @type {Array<Promise<RequestResponsePair>>} */
		const pairs = [];

		const finish = () => {
			clearTimeout(timeoutID);
			page.off('request', onRequest);
			Promise.all(pairs).then(resolve);
		};

		const timeoutID = setTimeout(() => {
			console.warn(
				'Timeout waiting for requests. Resolving with partial data.',
			);
			finish();
		}, REQUEST_TIMEOUT);

		/**
		 * @param {import('@playwright/test').Request} request
		 */
		const onRequest = (request) => {
			const matchingExpected = requestsToMatch.find((expectedRequest) =>
				matchesExpectedRequest(request, expectedRequest),
			);
			if (!matchingExpected) {
				return;
			}

			console.log('>>', request.method(), request.url());
			pairs.push(collectRequestResponse(request, matchingExpected));
			if (pairs.length === requestsToMatch.length) {
				finish();
			}
		};

		page.on('request', onRequest);
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
				url: getPathWithSearch(request.url()),
				status: status,
				postData:
					request.method() === 'POST' ? request.postDataJSON() : undefined,
				responseBody: responseBody,
			};
		}),
	);

	const expected = getExpectedRequests(page, expectedRequests).map((exp) => ({
		method: exp.method,
		url: exp.url.includes('?') ? expect.stringContaining(exp.url) : exp.url,
		status: exp.status,
		postData: exp.postData
			? expect.objectContaining({
					...exp.postData,
					...(exp.postData.u
						? { u: expect.stringContaining(exp.postData.u) }
						: {}),
				})
			: undefined,
		responseBody: exp.status !== 204 ? exp.responseBody : undefined,
	}));

	expect.soft(actualRequests).toEqual(expect.arrayContaining(expected));
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
