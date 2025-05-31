import fs from 'node:fs';
import { vitePlugin as remix } from '@remix-run/dev';
import browserslist from 'browserslist';
import { browserslistToTargets } from 'lightningcss';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

declare module '@remix-run/node' {
	interface Future {
		v3_singleFetch: true;
	}
}

const targets = browserslistToTargets(
	browserslist('defaults and fully supports es6-module'),
);

const customMedia = fs.readFileSync('./app/styles/_media.css', 'utf-8');

export default defineConfig({
	build: {
		cssMinify: 'lightningcss',
	},
	css: {
		transformer: 'lightningcss',
		lightningcss: {
			targets,
			drafts: {
				customMedia: true,
			},
		},
	},
	plugins: [
		{
			name: 'css-additional-data',
			enforce: 'pre',
			transform(code, id) {
				if (id.endsWith('.css') || id.endsWith('.css?inline')) {
					return `${customMedia}\n${code}`;
				}
			},
		},
		remix({
			ssr: false,
			future: {
				v3_fetcherPersist: true,
				v3_relativeSplatPath: true,
				v3_throwAbortReason: true,
				v3_lazyRouteDiscovery: true,
				v3_singleFetch: true,
			},
		}),
		tsconfigPaths(),
	],
});
