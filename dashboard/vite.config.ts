import { vitePlugin as remix } from '@remix-run/dev';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
	build: {
		cssMinify: 'lightningcss',
	},
	css: {
		transformer: 'lightningcss',
		lightningcss: {
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
