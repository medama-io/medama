(function () {
	if (document) {
		var p = document.currentScript,
			f = document.location.protocol + '//' + p.getAttribute('data-api'),
			q = () => Date.now().toString(36) + Math.random().toString(36).substr(2),
			g = q(),
			e = !1,
			h = !1,
			c = 0,
			k = 0,
			l = !1,
			t = history.pushState,
			u = history.replaceState,
			r = (b) =>
				new Promise((m) => {
					const d = new XMLHttpRequest();
					d.onload = () => {
						m(0 == d.responseText);
					};
					d.open('GET', b);
					d.setRequestHeader('Content-Type', 'text/plain');
					d.send();
				}),
			n = () => {
				e = !1;
				g = q();
				k = c = 0;
				l = !1;
			},
			a = (b) => {
				1 != b ||
					e ||
					r(
						f +
							'/event/ping?u=' +
							encodeURIComponent(location.host + location.pathname)
					).then((m) => {
						h = m;
					});
				l ||
					navigator.sendBeacon(
						f + '/event/hit',
						JSON.stringify(
							1 == b
								? {
										b: g,
										u: location.href,
										r: document.referrer,
										p: e,
										q: h,
										t: Intl.DateTimeFormat().resolvedOptions().timeZone,
								  }
								: { b: g, m: Math.round(performance.now() - c - k) }
						)
					);
				0 == b && (l = !0);
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
					: ((k += performance.now() - c), (c = 0));
			},
			{ capture: !0 }
		);
		r(f + '/event/ping').then((b) => {
			h = e = b;
			a(1);
			p.getAttribute('data-hash')
				? document.addEventListener(
						'hashchange',
						() => {
							a(1);
						},
						{ capture: !0 }
				  )
				: ((history.pushState = function () {
						a(0);
						n();
						t.apply(history, arguments);
						a(1);
				  }),
				  (history.replaceState = function () {
						a(0);
						n();
						u.apply(history, arguments);
						a(1);
				  }),
				  window.addEventListener(
						'popstate',
						() => {
							n();
							a(1);
						},
						{ capture: !0 }
				  ));
		});
	}
})();
