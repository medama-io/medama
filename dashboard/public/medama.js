(function () {
	if (document) {
		var n = document.currentScript,
			e = document.location.protocol + '//' + n.getAttribute('data-api'),
			p = () => Date.now().toString(36) + Math.random().toString(36).substr(2),
			f = p(),
			g = !0,
			q = !0,
			c = 0,
			h = 0,
			k = !1,
			t = history.pushState,
			u = history.replaceState,
			r = (b) =>
				new Promise((l) => {
					const d = new XMLHttpRequest();
					d.onload = () => {
						l(0 == d.responseText);
					};
					d.open('GET', b);
					d.setRequestHeader('Content-Type', 'text/plain');
					d.send();
				}),
			m = () => {
				g = !1;
				f = p();
				h = c = 0;
				k = !1;
			},
			a = (b) => {
				1 == b &&
					r(
						e +
							'/event/ping?u=' +
							encodeURIComponent(location.host + location.pathname)
					).then((l) => {
						q = l;
					});
				k ||
					navigator.sendBeacon(
						e + '/event/hit',
						JSON.stringify(
							1 == b
								? {
										b: f,
										u: location.href,
										r: document.referrer,
										e: 'load',
										p: g,
										q,
										t: Intl.DateTimeFormat().resolvedOptions().timeZone,
								  }
								: {
										b: f,
										e: 'unload',
										m: Math.round(performance.now() - c - h),
								  }
						)
					);
				0 == b && (k = !0);
			};
		'onpagehide' in self
			? document.addEventListener(
					'pagehide',
					() => {
						a(0);
					},
					{ capture: !0 }
			  )
			: (document.addEventListener(
					'beforeunload',
					() => {
						a(0);
					},
					{ capture: !0 }
			  ),
			  document.addEventListener(
					'unload',
					() => {
						a(0);
					},
					{ capture: !0 }
			  ));
		document.addEventListener(
			'visibilitychange',
			() => {
				'hidden' == document.visibilityState
					? (c = performance.now())
					: ((h += performance.now() - c), (c = 0));
			},
			{ capture: !0 }
		);
		r(e + '/event/ping').then((b) => {
			g = b;
			a(1);
			n.getAttribute('data-hash')
				? document.addEventListener(
						'hashchange',
						() => {
							a(1);
						},
						{ capture: !0 }
				  )
				: ((history.pushState = function () {
						a(0);
						m();
						t.apply(history, arguments);
						a(1);
				  }),
				  (history.replaceState = function () {
						a(0);
						m();
						u.apply(history, arguments);
						a(1);
				  }),
				  window.addEventListener(
						'popstate',
						() => {
							m();
							a(1);
						},
						{ capture: !0 }
				  ));
		});
	}
})();
