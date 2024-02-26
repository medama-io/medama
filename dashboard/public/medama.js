(function () {
	if (document) {
		var n = document.currentScript,
			e = document.location.protocol + '//' + n.getAttribute('data-api'),
			p = () => Date.now().toString(36) + Math.random().toString(36).substr(2),
			g = p(),
			h = !0,
			q = !0,
			f = 0,
			k = Date.now(),
			l = !1,
			t = history.pushState,
			u = history.replaceState,
			r = (a) =>
				new Promise((v) => {
					const b = new XMLHttpRequest();
					b.onload = () => {
						v(0 == b.responseText);
					};
					b.open('GET', a);
					b.setRequestHeader('Content-Type', 'text/plain');
					b.send();
				}),
			m = () => {
				h = !1;
				g = p();
				f = 0;
				k = Date.now();
				l = !1;
			},
			c = () => {
				r(
					e +
						'/event/ping?u=' +
						encodeURIComponent(location.host + location.pathname)
				).then((a) => {
					q = a;
					navigator.sendBeacon(
						e + '/event/hit',
						JSON.stringify({
							b: g,
							u: location.href,
							r: document.referrer,
							e: 'load',
							p: h,
							q,
							t: Intl.DateTimeFormat().resolvedOptions().timeZone,
						})
					);
				});
			},
			d = () => {
				l ||
					navigator.sendBeacon(
						e + '/event/hit',
						JSON.stringify({ b: g, e: 'unload', m: Date.now() - k })
					);
				l = !0;
			};
		'onpagehide' in self
			? document.addEventListener(
					'pagehide',
					() => {
						d();
					},
					{ capture: !0 }
			  )
			: (document.addEventListener(
					'beforeunload',
					() => {
						d();
					},
					{ capture: !0 }
			  ),
			  document.addEventListener(
					'unload',
					() => {
						d();
					},
					{ capture: !0 }
			  ));
		document.addEventListener(
			'visibilitychange',
			() => {
				'hidden' == document.visibilityState
					? (f = Date.now())
					: ((k += Date.now() - f), (f = 0));
			},
			{ capture: !0 }
		);
		r(e + '/event/ping').then((a) => {
			h = a;
			c();
			n.getAttribute('data-hash')
				? document.addEventListener(
						'hashchange',
						() => {
							c();
						},
						{ capture: !0 }
				  )
				: ((history.pushState = function () {
						d();
						m();
						t.apply(history, arguments);
						c();
				  }),
				  (history.replaceState = function () {
						d();
						m();
						u.apply(history, arguments);
						c();
				  }),
				  window.addEventListener(
						'popstate',
						() => {
							m();
							c();
						},
						{ capture: !0 }
				  ));
		});
	}
})();
