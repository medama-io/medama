const ACCEPTED_PATHS = ['/simple', '/hash', '/history'];

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

		if (url.pathname.endsWith('/script.js')) {
			return new Response(Bun.file(__dirname + '/../../dist/default.js'));
		}

		for (const path of ACCEPTED_PATHS) {
			if (url.pathname.startsWith(path)) {
				const file = __dirname + url.pathname;
				console.log('Serving:', file);
				return new Response(Bun.file(file));
			}
		}

		// reject all other routes
		console.error('Not found:', url.pathname);
		return new Response(null, { status: 404 });
	},
});
