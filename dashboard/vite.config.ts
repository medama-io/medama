import fs from 'node:fs';
import { vitePlugin as remix } from '@remix-run/dev';
import browserslist from 'browserslist';
import { browserslistToTargets } from 'lightningcss';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

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
				if (id.endsWith('.css')) {
					return `${customMedia}\n${code}`;
				}
			},
		},
		remix({
			ssr: false,
		}),
		tsconfigPaths(),
	],
	server: {
		watch: {
			usePolling: true,
		},
	},
});
