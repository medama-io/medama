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
 * @property {string=} t Title of the page.
 * @property {string=} d Timezone of the user.
 * @property {number=} w Screen width.
 * @property {number=} h Screen height.
 * @property {number=} m Time spent on page. Only sent on unload.
 */
var Payload;

/**
 * Note that we don't try inline global values such as `self` or `document` because
 * while it does reduce actual bundle size, it is LESS efficient with gzip compression
 * which is more practical.
 *
 * @see https://github.com/google/closure-compiler/wiki/FAQ#closure-compiler-inlined-all-my-strings-which-made-my-code-size-bigger-why-did-it-do-that
 */
(function () {
	// If server-side rendering, bail out.
	if (!document) {
		return;
	}

	/**
	 * Get API URL from data-host in script tag with the correct protocol.
	 */
	const host =
		document.location.protocol +
		'//' +
		document.currentScript.getAttribute('data-api');

	/**
	 * Generate a unique ID for linking multiple beacon events together for the same user.
	 * This is necessary for us to determine how long someone has spent on a page.
	 *
	 * @remarks We intentionally use Math.random() instead of the Web Crypto API
	 * because uniqueness against collisions is not a requirement and is worth
	 * the tradeoff for bundle size and performance.
	 */
	const uid = Date.now().toString(36) + Math.random().toString(36).substr(2);

	/**
	 * Ping the server with the cache endpoint and read the last modified header to determine
	 * if the user is unique or not.
	 *
	 * If the response is not cached, then the user is unique. If it is cached, the
	 * last-modified header will return the timestamp of the day at midnight incremented by
	 * how many times the user has visited the site.
	 */
	let isUnique = true;
	let xhr = new XMLHttpRequest();
	xhr.open('GET', host + '/event/ping');
	xhr.onreadystatechange = function () {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			if (xhr.status === 304) {
				isUnique = false;
			}
		}
	};

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
			"t": document.title,
			"w": self.screen.width,
			"h": self.screen.height,
			"e": eventType,
			"m":
				eventType === EventType.PAGEHIDE ||
				eventType === EventType.HIDDEN ||
				eventType === EventType.UNLOAD
					? self.performance.now() - hiddenTimeMs
					: undefined,
		};

		navigator.sendBeacon(host + '/event/hit', JSON.stringify(payload));
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
			{
				capture: true,
			}
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
			{
				capture: true,
			}
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

	// Load event is used to track when a user has loaded in a page.
	document.addEventListener(
		EventType.LOAD,
		() => {
			sendBeacon(EventType.LOAD);
		},
		{
			once: true,
			capture: true,
		}
	);
})();
