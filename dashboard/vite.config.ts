import { vitePlugin as remix } from '@remix-run/dev';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
	server: {
		watch: {
			usePolling: true,
		},
	},
	plugins: [
		remix({
			ssr: false,
		}),
		tsconfigPaths(),
	],
});
