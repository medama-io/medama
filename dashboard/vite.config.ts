import fs from 'node:fs';
import { reactRouter } from '@react-router/dev/vite';
import babel from '@rolldown/plugin-babel';
import { reactCompilerPreset } from '@vitejs/plugin-react';
import browserslist from 'browserslist';
import { browserslistToTargets } from 'lightningcss';
import { defineConfig } from 'vite';

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
			? {
					tsconfigPaths: true,
				}
			: {
					tsconfigPaths: true,
					alias: {
						'react-dom/server': 'react-dom/server.node',
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
		babel({
			presets: [reactCompilerPreset()],
		}),
		reactRouter(),
	],
});
