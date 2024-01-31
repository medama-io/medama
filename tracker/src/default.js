/**
 * Event types for addEventListener function calls.
 *
 * @remark Closure compiler doesn't understand string unions, so we have to
 * use an enum instead.
 * @enum {string}
 */
const EventType = {
	UNLOAD: 'unload',
	LOAD: 'load',
	// These events must still be sent with the unload event when calling
	// sendBeacon to ensure there are no duplicate events.
	PAGEHIDE: 'pagehide',
	BEFOREUNLOAD: 'beforeunload',
	// Custom events that are not part of the event listener spec, but is
	// used to determine what state visibilitychange is in.
	VISIBILITYCHANGE: 'visibilitychange',
	HIDDEN: 'hidden',
	VISIBLE: 'visible',
};

/**
 * Event types for beacon function calls.
 *
 * @remark We use a different enum for beacon types to reduce the bundle size
 * since numbers are smaller than strings.
 * @enum {number}
 */
const BeaconType = {
	UNLOAD: 0,
	LOAD: 1,
};

/**
 * @typedef {Object} HitPayload
 * @property {string} b Beacon ID.
 * @property {string} u Page URL.
 * @property {string} r Referrer URL.
 * @property {EventType} e Event type.
 * @property {boolean} p If the user is unique or not.
 * @property {boolean} q If this is the first time the user has visited this specific page.
 * @property {string} t Timezone of the user.
 */
var HitPayload;

/**
 * @typedef {Object} DurationPayload
 * @property {string} b Beacon ID.
 * @property {EventType} e Event type.
 * @property {number} m Time spent on page.
 */

var DurationPayload;

/**
 * Note that we don't try to inline global values such as `self` or `document` because
 * while it does reduce actual bundle size, it is LESS efficient with gzip compression
 * which should be a more practical benchmark for users.
 *
 * @see https://github.com/google/closure-compiler/wiki/FAQ#closure-compiler-inlined-all-my-strings-which-made-my-code-size-bigger-why-did-it-do-that
 */
(function () {
	// If server-side rendering, bail out.
	if (!document) {
		return;
	}

	/**
	 * document.currentScript can only be called when the script is being executed. If
	 * we call the script in an event listener, then it will be null. So we need to
	 * make a copy of the currentScript object to use later.
	 */
	const currentScript = document.currentScript;

	/**
	 * Get API URL from data-host in script tag with the correct protocol.
	 */
	const host =
		document.location.protocol + '//' + currentScript.getAttribute('data-api');

	/**
	 * Generate a unique ID for linking multiple beacon events together for the same page
	 * view. This is necessary for us to determine how long someone has spent on a page.
	 *
	 * @remarks We intentionally use Math.random() instead of the Web Crypto API
	 * because uniqueness against collisions is not a requirement and is worth
	 * the tradeoff for bundle size and performance.
	 */
	const generateUid = () =>
		Date.now().toString(36) + Math.random().toString(36).substr(2);

	/**
	 * Unique ID linking multiple beacon events together for the same page view.
	 */
	let uid = generateUid();

	/**
	 * Whether the user is unique or not.
	 * This is updated when the server checks the ping cache on page load.
	 */
	let isUnique = true;

	/**
	 * Whether the user is visiting this page for the first time.
	 */
	let isFirstVisit = true;

	/**
	 * A temporary variable to store the start time of the page when it is hidden.
	 */
	let hiddenStartTime = 0;

	/**
	 * The total time the user has had the page hidden.
	 */
	let hiddenTotalTime = 0;

	/**
	 * Ensure only the unload beacon is called once.
	 */
	let isUnloadCalled = false;

	/**
	 * @remarks We hoist the following variables to the top to let the closure compiler
	 * infer that it can declare these variables together with the other variables in a
	 * single line instead of separately, which saves us a few bytes.
	 */

	/**
	 * Copy of the original pushState and replaceState functions, used for overriding
	 * the History API to track navigation changes.
	 */
	const historyPush = history.pushState;
	const historyReplace = history.replaceState;

	/**
	 * Ping the server with the cache endpoint and read the last modified header to determine
	 * if the user is unique or not.
	 *
	 * If the response is not cached, then the user is unique. If it is cached, then the
	 * browser will send an If-Modified-Since header indicating the user is not unique.
	 *
	 * @param {string} url URL to ping.
	 * @returns {Promise<boolean>} Is the cache unique or not.
	 */
	const pingCache = (url) =>
		new Promise((resolve) => {
			const xhr = new XMLHttpRequest();
			xhr.onload = () => {
				// @ts-ignore - Double equals reduces bundle size.
				resolve(xhr.responseText == 0);
			};
			xhr.open('GET', url);
			xhr.setRequestHeader('Content-Type', 'text/plain');
			xhr.send();
		});

	/**
	 * Cleanup temporary variables and reset the unique ID.
	 */
	const cleanup = () => {
		// Main ping cache won't be called again, so we can assume the user is not unique.
		// However, isFirstVisit will be called on each page load, so we don't need to reset it.
		isUnique = false;
		uid = generateUid();
		hiddenStartTime = 0;
		hiddenTotalTime = 0;
		isUnloadCalled = false;
	};

	/**
	 * Send a beacon event to the server.
	 *
	 * @param {BeaconType} beaconType Load or unload event type.
	 * @returns {void}
	 */
	const sendBeacon = (beaconType) => {
		if (beaconType == BeaconType.LOAD) {
			// Returns true if it is the user's first visit to page, false if not.
			// The u query parameter is a cache busting parameter which is the page host and path
			// without protocol or query parameters.
			pingCache(
				host +
					'/event/ping?u=' +
					encodeURIComponent(location.host + location.pathname)
			).then((response) => {
				isFirstVisit = response;
			});
		}

		if (!isUnloadCalled) {
			navigator.sendBeacon(
				host + '/event/hit',
				JSON.stringify(
					beaconType == BeaconType.LOAD
						? // prettier-ignore
						  /**
						   * Payload to send to the server.
						   * @type {HitPayload}
						   * @remarks We use string literals for the keys to tell Closure Compiler
						   * to not rename them.
						   */ ({
								"b": uid,
								"u": location.href,
								"r": document.referrer,
								"e": EventType.LOAD,
								"p": isUnique,
								"q": isFirstVisit,
								/**
								 * Get timezone for country detection.
								 *
								 * @suppress {checkTypes} Compiler throws an error because we don't call
								 * "new" for this even though it is unnecessary.
								 * @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DateTimeFormat/DateTimeFormat#return_value
								 */
								"t": Intl.DateTimeFormat().resolvedOptions().timeZone,
						  })
						: // prettier-ignore
						  /**
						   * Payload to send to the server.
						   * @type {DurationPayload}
						   * @remarks We use string literals for the keys to tell Closure Compiler
						   * to not rename them.
						   */
						  ({
								"b": uid,
								"e": EventType.UNLOAD,
								"m": Math.round(
									performance.now() - hiddenStartTime - hiddenTotalTime
								),
						  })
				)
			);
		}

		if (beaconType == BeaconType.UNLOAD) {
			// Ensure unload is only called once.
			isUnloadCalled = true;
		}
	};

	// Prefer pagehide if available because it's more reliable than unload.
	// We also prefer pagehide because it doesn't break bfcache.
	if ('onpagehide' in self) {
		/**
		 * @suppress {checkTypes}
		 */
		document.addEventListener(
			EventType.PAGEHIDE,
			() => {
				sendBeacon(BeaconType.UNLOAD);
			},
			{ capture: true }
		);
	} else {
		// Otherwise, use unload and beforeunload. Using both is significantly more
		// reliable than just one due to browser differences. However, this will break
		// bfcache, but it's better than nothing.
		document.addEventListener(
			EventType.BEFOREUNLOAD,
			() => {
				sendBeacon(BeaconType.UNLOAD);
			},
			{ capture: true }
		);
		document.addEventListener(
			EventType.UNLOAD,
			() => {
				sendBeacon(BeaconType.UNLOAD);
			},
			{ capture: true }
		);
	}

	// Visibility change events allow us to track whether a user is tabbed out and
	// correct our timings.
	document.addEventListener(
		EventType.VISIBILITYCHANGE,
		() => {
			if (document.visibilityState == EventType.HIDDEN) {
				// Page is hidden, record the current time.
				hiddenStartTime = performance.now();
			} else {
				// Page is visible, subtract the hidden time to calculate the total time hidden.
				hiddenTotalTime += performance.now() - hiddenStartTime;
				hiddenStartTime = 0;
			}
		},
		{ capture: true }
	);

	pingCache(host + '/event/ping').then((response) => {
		// The response is a boolean indicating if the user is unique or not.
		isUnique = response;

		// Send the first beacon event to the server.
		sendBeacon(BeaconType.LOAD);

		// Check if hash mode is enabled. If it is, then we need to send a beacon event
		// when the hash changes. If disabled, it is safe to override the History API.
		if (currentScript.getAttribute('data-hash')) {
			// Hash mode is enabled. Add hashchange event listener.
			document.addEventListener(
				'hashchange',
				() => {
					sendBeacon(BeaconType.LOAD);
				},
				{
					capture: true,
				}
			);
		} else {
			// Add pushState event listeners to track navigation changes with
			// router libraries that use the History API.
			history.pushState = function () {
				sendBeacon(BeaconType.UNLOAD);
				// If the event is a history change, then we need to reset the id and timers
				// because the page is not actually reloading the script.
				cleanup();
				historyPush.apply(history, arguments);
				sendBeacon(BeaconType.LOAD);
			};

			// replaceState is used by some router libraries to replace the current
			// history state instead of pushing a new one.
			history.replaceState = function () {
				sendBeacon(BeaconType.UNLOAD);
				cleanup();
				historyReplace.apply(history, arguments);
				sendBeacon(BeaconType.LOAD);
			};

			// popstate is fired when the back or forward button is pressed.
			// We use window instead of document here because the document state
			// doesn't change immediately when the event is fired.
			window.addEventListener(
				'popstate',
				() => {
					// Unfortunately, we can't use unload here because we can't call it before
					// the history change, so cleanup any temporary variables here.
					cleanup();
					sendBeacon(BeaconType.LOAD);
				},
				{
					capture: true,
				}
			);
		}
	});
})();
