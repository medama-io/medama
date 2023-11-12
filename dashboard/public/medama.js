(function () {
	if (document) {
		var d =
				document.location.protocol +
				'//' +
				document.currentScript.getAttribute('data-api'),
			h = Date.now().toString(36) + Math.random().toString(36).substr(2),
			e = !0,
			b = new XMLHttpRequest();
		b.open('GET', d + '/event/ping');
		b.onreadystatechange = function () {
			b.readyState === XMLHttpRequest.DONE && 304 === b.status && (e = !1);
		};
		var f = 0,
			g = 0,
			c = function (a) {
				a = {
					b: h,
					u: location.href,
					r: document.referrer,
					d: Intl.DateTimeFormat().resolvedOptions().timeZone,
					p: e,
					t: document.title,
					w: self.screen.width,
					h: self.screen.height,
					e: a,
					m:
						'pagehide' === a || 'hidden' === a || 'unload' === a
							? self.performance.now() - f
							: void 0,
				};
				navigator.sendBeacon(d + '/event/hit', JSON.stringify(a));
			};
		'onpagehide' in self
			? document.addEventListener(
					'pagehide',
					function () {
						c('pagehide');
					},
					{ capture: !0 }
			  )
			: document.addEventListener(
					'unload',
					function () {
						c('unload');
					},
					{ capture: !0 }
			  );
		document.addEventListener(
			'visibilitychange',
			function () {
				'hidden' === document.visibilityState
					? ((g = self.performance.now()), c('hidden'))
					: (f += self.performance.now() - g);
			},
			{ capture: !0 }
		);
		document.addEventListener(
			'load',
			function () {
				c('load');
			},
			{ once: !0, capture: !0 }
		);
	}
})();
