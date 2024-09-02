const ACCEPTED_PATHS = ['/simple', '/history'];

console.log('Serving on http://localhost:3000');
Bun.serve({
	port: 3000,
	fetch(req) {
		const url = new URL(req.url);

		if (
			url.pathname === '/' ||
			url.pathname === '/index.html' ||
			url.pathname === '/favicon.ico'
		) {
			return new Response(null, { status: 200 });
		}

		for (const path of ACCEPTED_PATHS) {
			if (url.pathname.startsWith(path)) {
				const file = __dirname + url.pathname;
				console.log('Serving:', file);
				return new Response(Bun.file(file));
			}
		}

		// Proxy /script.js to serve local file in dist folder
		if (url.pathname === '/script.js') {
			console.log('Serving:', url.pathname);
			return new Response(
				Bun.file(__dirname + '/../../dist/click-events.page-events.min.js'),
			);
		}

		// Proxy all other routes to API
		console.log('Proxying:', url.pathname);
		const newURL = new URL(req.url);
		newURL.port = '8080';
		return fetch(new Request(newURL, req));
	},
});
