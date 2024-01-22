/**
 * Event types for addEventListener function calls.
 *
 * @remark Closure compiler doesn't understand string unions, so we have to
 * use an enum instead.
 * @enum {string}
 */
const EventType = {
	PAGEHIDE: 'pagehide',
	UNLOAD: 'unload',
	LOAD: 'load',
	VISIBILITYCHANGE: 'visibilitychange',
	// Custom events that are not part of the event listener spec, but is
	// used to determine what state visibilitychange is in.
	HIDDEN: 'hidden',
	VISIBLE: 'visible',
};

/**
 * @typedef {Object} Payload
 * @property {string} b Beacon ID.
 * @property {string} u Page URL.
 * @property {string} r Referrer URL.
 * @property {EventType} e Event type.
 * @property {boolean=} p If the user is unique or not.
 * @property {string=} d Timezone of the user.
 * @property {number=} m Time spent on page. Only sent on unload.
 */
var Payload;

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
	 * Generate a unique ID for linking multiple beacon events together for the same user.
	 * This is necessary for us to determine how long someone has spent on a page.
	 *
	 * @remarks We intentionally use Math.random() instead of the Web Crypto API
	 * because uniqueness against collisions is not a requirement and is worth
	 * the tradeoff for bundle size and performance.
	 */
	let uid = Date.now().toString(36) + Math.random().toString(36).substr(2);

	/**
	 * Whether the user is unique or not.
	 * This is updated when the server checks the ping cache on page load.
	 */
	let isUnique = true;

	/**
	 * Counter for how long a page may have been hidden.
	 * This will then be removed from the total time spent on a page.
	 */
	let hiddenTimeMs = 0;
	/**
	 * The temporary counter is used to keep track of how long a page has been hidden.
	 * Then when the page becomes visible again, we can subtract the hidden time from
	 * the total time spent on a page.
	 */
	let hiddenTimeTemp = 0;

	/**
	 * @remarks We hoist the following variables to the top to let the closure compiler
	 * infer that it can declare these variables together with the other variables in a
	 * single line instead of separately, which saves us a few bytes.
	 */

	/**
	 * XMLHttpRequest object for pinging the server used in load event.
	 */
	const xhr = new XMLHttpRequest();

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
		// Ping cache won't be called again, so we can assume the user is not unique.
		isUnique = false;
		uid = Date.now().toString(36) + Math.random().toString(36).substr(2);
		hiddenTimeMs = 0;
		hiddenTimeTemp = 0;
	};

	/**
	 * Send a beacon event to the server.
	 *
	 * @param {EventType} eventType Event type.
	 * @returns {void}
	 */
	const sendBeacon = (eventType) => {
		/**
		 * Payload to send to the server.
		 * @type {Payload}
		 * @remarks We use string literals for the keys to tell Closure Compiler
		 * to not rename them.
		 */
		// prettier-ignore
		const payload = {
			"b": uid,
			"u": location.href,
			"r": document.referrer,
			/**
			 * Get timezone for country detection.
			 *
			 * @suppress {checkTypes} Compiler throws an error because we don't call
			 * "new" for this even though it is unnecessary.
			 * @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DateTimeFormat/DateTimeFormat#return_value
			 */
			"d": Intl.DateTimeFormat().resolvedOptions().timeZone,
			"p": isUnique,
			"e": eventType,
			"m":
				eventType === EventType.PAGEHIDE ||
				eventType === EventType.HIDDEN ||
				eventType === EventType.UNLOAD
					? self.performance.now() - hiddenTimeMs
					: undefined,
		};

		navigator.sendBeacon(host + '/event/hit', JSON.stringify(payload));

		// If the event is a history change, then we need to reset the id and timers
		// because the page is not actually reloading the script.
		if (eventType === EventType.UNLOAD) {
			cleanup();
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
				sendBeacon(EventType.PAGEHIDE);
			},
			{ capture: true }
		);
	} else {
		// Otherwise, use unload. This will break bfcache, but it's better than nothing.
		// We can also use beforeunload as well to improve reliability, but it isn't
		// worth the extra code deduplicating both events for the slight increase in
		// accuracy. Mobile browsers don't even fire beforeunload.
		document.addEventListener(
			EventType.UNLOAD,
			() => {
				sendBeacon(EventType.UNLOAD);
			},
			{ capture: true }
		);
	}

	// Visibility change events allow us to track whether a user is tabbed out and
	// correct our timings. It is also an additional fallback to unload events.
	document.addEventListener(
		EventType.VISIBILITYCHANGE,
		() => {
			if (document.visibilityState === EventType.HIDDEN) {
				hiddenTimeTemp = self.performance.now();
				sendBeacon(EventType.HIDDEN);
			} else {
				hiddenTimeMs += self.performance.now() - hiddenTimeTemp;
			}
		},
		{ capture: true }
	);

	/**
	 * Ping the server with the cache endpoint and read the last modified header to determine
	 * if the user is unique or not.
	 *
	 * If the response is not cached, then the user is unique. If it is cached, then the
	 * browser will send an If-Modified-Since header indicating the user is not unique.
	 */
	xhr.open('GET', host + '/event/ping');
	xhr.setRequestHeader('Content-Type', 'text/plain');
	xhr.addEventListener(
		EventType.LOAD,
		() => {
			// The server will respond with a 0 or 1 depending on whether an If-Modified-Since
			// header was sent or not. It also considers if the header is more than a day old
			// to reset the cache. If the response is 1 then the user is not unique.
			if (xhr.responseText === '1') {
				isUnique = false;
			}

			// Send the first beacon event to the server.
			sendBeacon(EventType.LOAD);

			// Check if hash mode is enabled. If it is, then we need to send a beacon event
			// when the hash changes. If disabled, it is safe to override the History API.
			if (currentScript.getAttribute('data-hash')) {
				// Hash mode is enabled. Add hashchange event listener.
				document.addEventListener(
					'hashchange',
					() => {
						sendBeacon(EventType.LOAD);
					},
					{
						capture: true,
					}
				);
			} else {
				// Add pushState event listeners to track navigation changes with
				// router libraries that use the History API.
				history.pushState = function () {
					sendBeacon(EventType.UNLOAD);
					historyPush.apply(history, arguments);
					sendBeacon(EventType.LOAD);
				};

				// replaceState is used by some router libraries to replace the current
				// history state instead of pushing a new one.
				history.replaceState = function () {
					sendBeacon(EventType.UNLOAD);
					historyReplace.apply(history, arguments);
					sendBeacon(EventType.LOAD);
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
						sendBeacon(EventType.LOAD);
					},
					{
						capture: true,
					}
				);
			}

			// Add event listeners for all outbound links.
			// Get all elements with an href attribute.
			/* const links = document.getElementsByTagName('a');
			for (let i = 0; i < links.length; i++) {
				// Get the link.
				const link = links[i];

				// Check if the link is outbound.
				if (link.host !== document.location.host && link.host !== '') {
					// Add event listener to the link.
					link.addEventListener('click', () => {
						sendBeacon(EventType.REPLACE);
					});

					// Handle middle click and ctrl/cmd click.
					link.addEventListener(
						'auxclick',
						(event) => {
							if (event.button === 1) {
								sendBeacon(EventType.REPLACE);
							}
						},
						{
							capture: true,
						}
					);
				}
			} */
		},
		{
			capture: true,
			// The load event might be called multiple times, so we only want to
			// listen to it once.
			once: true,
		}
	);
	xhr.send();
})();
