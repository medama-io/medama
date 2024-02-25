(function () {
	if (document) {
		var q = document.currentScript,
			e = document.location.protocol + '//' + q.getAttribute('data-api'),
			r = () => Date.now().toString(36) + Math.random().toString(36).substr(2),
			h = r(),
			k = location.origin + location.pathname,
			l = !0,
			t = !0,
			f = 0,
			m = Date.now(),
			n = !1,
			u = history.pushState,
			v = history.replaceState,
			w = (a) =>
				new Promise((g) => {
					const b = new XMLHttpRequest();
					b.onload = () => {
						g(0 == b.responseText);
					};
					b.open('GET', a);
					b.setRequestHeader('Content-Type', 'text/plain');
					b.send();
				}),
			p = () => {
				l = !1;
				h = r();
				f = 0;
				m = Date.now();
				n = !1;
				k = location.host + location.pathname;
			},
			c = () => {
				w(
					e +
						'/event/ping?u=' +
						encodeURIComponent(location.host + location.pathname)
				).then((a) => {
					t = a;
					navigator.sendBeacon(
						e + '/event/hit',
						JSON.stringify({
							b: h,
							u: k,
							r: document.referrer,
							e: 'load',
							p: l,
							q: t,
							t: Intl.DateTimeFormat().resolvedOptions().timeZone,
						})
					);
				});
			},
			d = () => {
				n ||
					navigator.sendBeacon(
						e + '/event/hit',
						JSON.stringify({ b: h, e: 'unload', m: Date.now() - m })
					);
				n = !0;
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
					: ((m += Date.now() - f), (f = 0));
			},
			{ capture: !0 }
		);
		w(e + '/event/ping').then((a) => {
			l = a;
			c();
			if (q.getAttribute('data-hash'))
				document.addEventListener(
					'hashchange',
					() => {
						c();
					},
					{ capture: !0 }
				);
			else {
				const g = k !== location.origin + location.pathname;
				history.pushState = function () {
					g
						? (d(), p(), u.apply(history, arguments), c())
						: u.apply(history, arguments);
				};
				history.replaceState = function () {
					g
						? (d(), p(), v.apply(history, arguments), c())
						: v.apply(history, arguments);
				};
				window.addEventListener(
					'popstate',
					() => {
						p();
						c();
					},
					{ capture: !0 }
				);
			}
		});
	}
})();
