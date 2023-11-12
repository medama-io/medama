(function () {
	if (document) {
		var d =
				document.location.protocol +
				'//' +
				document.currentScript.getAttribute('data-api'),
			h = Date.now().toString(36) + Math.random().toString(36).substr(2),
			e = !0,
			c = new XMLHttpRequest();
		c.open('GET', d + '/event/ping');
		c.onreadystatechange = function () {
			c.readyState === XMLHttpRequest.DONE && 304 === c.status && (e = !1);
		};
		var f = 0,
			g = 0,
			b = function (a) {
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
						'unload' === a || 'pagehide' === a
							? self.performance.now() - f
							: void 0,
				};
				navigator.sendBeacon(d + '/event/hit', JSON.stringify(a));
			};
		'onpagehide' in self
			? document.addEventListener(
					'pagehide',
					function () {
						b('pagehide');
					},
					{ capture: !0 }
			  )
			: document.addEventListener(
					'unload',
					function () {
						b('unload');
					},
					{ capture: !0 }
			  );
		document.addEventListener(
			'visibilitychange',
			function () {
				'hidden' === document.visibilityState
					? ((g = self.performance.now()), b('hidden'))
					: ((f += self.performance.now() - g), b('visible'));
			},
			{ capture: !0 }
		);
		document.addEventListener(
			'load',
			function () {
				b('load');
			},
			{ capture: !0 }
		);
	}
})();
