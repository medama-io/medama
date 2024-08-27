import { vitePlugin as remix } from '@remix-run/dev';
import browserslist from 'browserslist';
import { browserslistToTargets } from 'lightningcss';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

const targets = browserslistToTargets(
	browserslist('defaults and fully supports es6-module'),
);

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
