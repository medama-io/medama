(function () {
	if (document) {
		var n = document.currentScript,
			f = document.location.protocol + '//' + n.getAttribute('data-api'),
			p = () => Date.now().toString(36) + Math.random().toString(36).substr(2),
			g = p(),
			h = !0,
			q = !0,
			a = 0,
			k = 0,
			l = !1,
			t = history.pushState,
			u = history.replaceState,
			r = (b) =>
				new Promise((v) => {
					const c = new XMLHttpRequest();
					c.onload = () => {
						v(0 == c.responseText);
					};
					c.open('GET', b);
					c.setRequestHeader('Content-Type', 'text/plain');
					c.send();
				}),
			m = () => {
				h = !1;
				g = p();
				k = a = 0;
				l = !1;
			},
			d = () => {
				r(
					f +
						'/event/ping?u=' +
						encodeURIComponent(location.host + location.pathname)
				).then((b) => {
					q = b;
					navigator.sendBeacon(
						f + '/event/hit',
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
			e = () => {
				l ||
					navigator.sendBeacon(
						f + '/event/hit',
						JSON.stringify({
							b: g,
							e: 'unload',
							m: Math.round(performance.now() - a - k),
						})
					);
				l = !0;
			};
		'onpagehide' in self
			? document.addEventListener(
					'pagehide',
					() => {
						e();
					},
					{ capture: !0 }
			  )
			: (document.addEventListener(
					'beforeunload',
					() => {
						e();
					},
					{ capture: !0 }
			  ),
			  document.addEventListener(
					'unload',
					() => {
						e();
					},
					{ capture: !0 }
			  ));
		document.addEventListener(
			'visibilitychange',
			() => {
				'hidden' == document.visibilityState
					? (a = performance.now())
					: ((k += performance.now() - a), (a = 0));
			},
			{ capture: !0 }
		);
		r(f + '/event/ping').then((b) => {
			h = b;
			d();
			n.getAttribute('data-hash')
				? document.addEventListener(
						'hashchange',
						() => {
							d();
						},
						{ capture: !0 }
				  )
				: ((history.pushState = function () {
						e();
						m();
						t.apply(history, arguments);
						d();
				  }),
				  (history.replaceState = function () {
						e();
						m();
						u.apply(history, arguments);
						d();
				  }),
				  window.addEventListener(
						'popstate',
						() => {
							m();
							d();
						},
						{ capture: !0 }
				  ));
		});
	}
})();
