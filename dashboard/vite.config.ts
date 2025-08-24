import fs from 'node:fs';
import { vitePlugin as remix } from '@remix-run/dev';
import browserslist from 'browserslist';
import { browserslistToTargets } from 'lightningcss';
import { defineConfig } from 'vite';
import commonjs from 'vite-plugin-commonjs';
import tsconfigPaths from 'vite-tsconfig-paths';

declare module '@remix-run/react' {
	interface Future {
		v3_singleFetch: true;
	}
}

const targets = browserslistToTargets(
	browserslist([
		'chrome >= 107',
		'edge >= 107',
		'firefox >= 104',
		'safari >= 16',
	]),
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
	resolve:
		process.env.NODE_ENV === 'development'
			? {}
			: {
					alias: {
						'react-dom/server': 'react-dom/server.node',
					},
				},
	plugins: [
		commonjs(),
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
				v3_routeConfig: true,
			},
			// biome-ignore lint/suspicious/noExplicitAny: Issue until we migrate to react-router.
		}) as any,
		tsconfigPaths(),
	],
});
