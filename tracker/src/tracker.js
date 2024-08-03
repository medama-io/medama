// @ts-check

/** @typedef {('load'|'unload'|'custom')} EventType */

/**
 * @typedef {Object} HitPayload
 * @property {string} b Beacon ID.
 * @property {EventType} e Event type.
 * @property {string} u Page URL.
 * @property {string} r Referrer URL.
 * @property {boolean} p If the user is unique or not.
 * @property {boolean} q If this is the first time the user has visited this specific page.
 * @property {string} t Timezone of the user.
 */

/**
 * @typedef {Object} DurationPayload
 * @property {string} b Beacon ID.
 * @property {EventType} e Event type.
 * @property {number} m Time spent on page.
 */

/**
 * @typedef {Object} CustomPayload
 * @property {string} g Group name of events. Currently, only uses the hostname.
 * @property {EventType} e Event type.
 * @property {Object} d Event custom properties.
 */

/**
 * Note that we don't try to inline global values such as `self` or `document` because
 * while it does reduce actual bundle size, it is LESS efficient with gzip compression
 * which should be a more practical benchmark for users.
 *
 * @see https://github.com/google/closure-compiler/wiki/FAQ#closure-compiler-inlined-all-my-strings-which-made-my-code-size-bigger-why-did-it-do-that
 */
(function () {
	// If server-side rendering, bail out. We use document instead of window here as Deno does have
	// a window object even on the server.
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
	 * Get API URL from data-api in script tag with the correct protocol.
	 * If the data-api attribute is not set, then we use the current script's
	 * src attribute to determine the host.
	 */
	const host = currentScript.getAttribute('data-api')
		? `${location.protocol}//${currentScript.getAttribute('data-api')}`
		: // @ts-ignore - We know this won't be an SVGScriptElement.
			currentScript.src.replace(/[^\/]+$/, 'api/');

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
	 * A temporary variable to store the start time of the page when it is hidden.
	 */
	let hiddenStartTime = 0;

	/**
	 * The total time the user has had the page hidden.
	 * It also signifies the start epoch time of the page.
	 */
	let hiddenTotalTime = Date.now();

	/**
	 * Ensure only the unload beacon is called once.
	 */
	let isUnloadCalled = false;

	/**
	 * @remarks We hoist the following variables to the top to let terser infer that it
	 * can declare these variables together with the other variables in a single line instead
	 * of separately, which saves us a few bytes.
	 */

	/**
	 * Copy of the original pushState and replaceState functions, used for overriding
	 * the History API to track navigation changes.
	 */
	const historyPush = history.pushState;
	const historyReplace = history.replaceState;

	/**
	 * Cleanup temporary variables and reset the unique ID.
	 */
	const cleanup = () => {
		// Main ping cache won't be called again, so we can assume the user is not unique.
		// However, isFirstVisit will be called on each page load, so we don't need to reset it.
		isUnique = false;
		uid = generateUid();
		hiddenStartTime = 0;
		hiddenTotalTime = Date.now();
		isUnloadCalled = false;
	};

	/**
	 * Wraps a history method with additional tracking events.
	 * @param {!Function} original - The original history method to wrap.
	 * @returns {function(this:History, *, string, (string | URL)=): void} The wrapped history method.
	 */
	const wrapHistoryFunc = (
		original,
		/**
		 * @this {History}
		 * @param {*} _state - The state object.
		 * @param {string} _unused - The title (unused).
		 * @param {(string | URL)=} url - The URL to navigate to.
		 * @returns {void}
		 */
	) =>
		function (_state, _unused, url) {
			if (url && location.pathname !== new URL(url, location.href).pathname) {
				sendUnloadBeacon();
				// If the event is a history change, then we need to reset the id and timers
				// because the page is not actually reloading the script.
				cleanup();
				original.apply(this, arguments);
				sendLoadBeacon();
			} else {
				original.apply(this, arguments);
			}
		};

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
			// We use XHR here because fetch GET request requires a CORS
			// header to be set on the server, which adds additional requests and
			// latency to ping the server.
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
	 * Send a load beacon event to the server when the page is loaded.
	 * @returns {Promise<void>}
	 */
	const sendLoadBeacon = async () => {
		// Returns true if it is the user's first visit to page, false if not.
		// The u query parameter is a cache busting parameter which is the page host and path
		// without protocol or query parameters.
		pingCache(
			host +
				'event/ping?u=' +
				encodeURIComponent(location.host + location.pathname),
		).then((isFirstVisit) => {
			// We use fetch here because it is more reliable than XHR.
			fetch(host + 'event/hit', {
				method: 'POST',
				body: JSON.stringify(
					// biome-ignore format: We use string literals for the keys to tell Closure Compiler to not rename them.
					/**
					 * Payload to send to the server.
					 * @type {HitPayload}
					 */ ({
						"b": uid,
						"e": "load",
						"u": location.href,
						"r": document.referrer,
						"p": isUnique,
						"q": isFirstVisit,
						/**
						 * Get timezone for country detection.
						 *
						 * @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DateTimeFormat/DateTimeFormat#return_value
						 */
						"t": Intl.DateTimeFormat().resolvedOptions().timeZone,
					}),
				),
				// Will make the response opaque, but we don't need it.
				mode: 'no-cors',
			});
		});
	};

	/**
	 * Send an unload beacon event to the server when the page is unloaded.
	 * @returns {void}
	 */
	const sendUnloadBeacon = () => {
		if (!isUnloadCalled) {
			// We use sendBeacon here because it is more reliable than fetch on page unloads.
			// The Fetch API keepalive flag has a few caveats and doesn't work very well on
			// Firefox on top of that. Previous experiements also seemed to indicate that
			// the fetch API doesn't work well on page unloads.
			// See: https://github.com/whatwg/fetch/issues/679
			//
			// Some adblockers block this API directly, but since this is the unload event,
			// it's an optional event to send.
			navigator.sendBeacon(
				host + 'event/hit',
				JSON.stringify(
					// biome-ignore format: We use string literals for the keys to tell Closure Compiler to not rename them.
					/**
					 * Payload to send to the server.
					 * @type {DurationPayload}
					 */
					({
						"b": uid,
						"e": "unload",
						"m": Date.now() - hiddenTotalTime,
					}),
				),
			);
		}

		// Ensure unload is only called once.
		isUnloadCalled = true;
	};

	// @ifdef TAGGED_EVENTS
	/**
	 * Send a custom beacon event to the server.
	 * @param {Object.<string, string>} properties Event custom properties.
	 * @returns {void}
	 */
	const sendCustomBeacon = (properties) => {
		// We use fetch here because it is more reliable than XHR.
		fetch(host + 'event/hit', {
			method: 'POST',
			body: JSON.stringify(
				// biome-ignore format: We use string literals for the keys to tell Closure Compiler to not rename them.
				/**
				 * Payload to send to the server.
				 * @type {CustomPayload}
				 */ ({
						"e": "custom",
						"g": location.hostname,
						"d": properties,
					}),
			),
			// Will make the response opaque, but we don't need it.
			mode: 'no-cors',
		});
	};
	// @endif

	// @ifdef TAGGED_EVENTS
	addEventListener('click', (event) => {
		const target =
			event.target instanceof HTMLElement
				? event.target.closest('[data-medama-*]')
				: null;
		if (target) {
			const attributes = target
				.getAttributeNames()
				.filter((attr) => attr.startsWith('data-medama-'));

			if (attributes.length > 0) {
				/** @type {Object<string, string>} */
				const data = {};
				attributes.forEach((attr) => {
					const key = attr.replace('data-medama-', '');
					const value = target.getAttribute(attr);
					data[key] = value;
				});

				sendCustomBeacon(data);
			}
		}
	});
	// @endif

	// Prefer pagehide if available because it's more reliable than unload.
	// We also prefer pagehide because it doesn't break bfcache.
	if ('onpagehide' in self) {
		addEventListener('pagehide', sendUnloadBeacon, { capture: true });
	} else {
		// Otherwise, use unload and beforeunload. Using both is significantly more
		// reliable than just one due to browser differences. However, this will break
		// bfcache, but it's better than nothing.
		addEventListener('beforeunload', sendUnloadBeacon, {
			capture: true,
		});
		addEventListener('unload', sendUnloadBeacon, { capture: true });
	}

	// Visibility change events allow us to track whether a user is tabbed out and
	// correct our timings.
	addEventListener(
		'visibilitychange',
		() => {
			if (document.visibilityState == 'hidden') {
				// Page is hidden, record the current time.
				hiddenStartTime = Date.now();
			} else {
				// Page is visible, subtract the hidden time to calculate the total time hidden.
				hiddenTotalTime += Date.now() - hiddenStartTime;
				hiddenStartTime = 0;
			}
		},
		{ capture: true },
	);

	pingCache(host + 'event/ping?root').then((response) => {
		// The response is a boolean indicating if the user is unique or not.
		isUnique = response;

		// Send the first beacon event to the server.
		sendLoadBeacon();

		// Check if hash mode is enabled. If it is, then we need to send a beacon event
		// when the hash changes. If disabled, it is safe to override the History API.
		if (currentScript.getAttribute('data-hash')) {
			// Hash mode is enabled. Add hashchange event listener.
			addEventListener('hashchange', sendLoadBeacon, {
				capture: true,
			});
		} else {
			//Add pushState event listeners to track navigation changes with
			//router libraries that use the History API.
			history.pushState = wrapHistoryFunc(historyPush);

			// replaceState is used by some router libraries to replace the current
			// history state instead of pushing a new one.
			history.replaceState = wrapHistoryFunc(historyReplace);

			// popstate is fired when the back or forward button is pressed.
			addEventListener(
				'popstate',
				() => {
					sendUnloadBeacon();
					cleanup();
					sendLoadBeacon();
				},
				{
					capture: true,
				},
			);
		}
	});
})();
