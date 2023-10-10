import '@mantine/core/styles.css';
import '@/styles/global.module.css';
import '@fontsource-variable/inter/wght.css';

import { enableReactUse } from '@legendapp/state/config/enableReactUse';
import { ColorSchemeScript, MantineProvider } from '@mantine/core';
import { cssBundleHref } from '@remix-run/css-bundle';
import type { LinksFunction } from '@remix-run/node';
import {
	Links,
	LiveReload,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
} from '@remix-run/react';

import { AppShell } from '@/components/layout/AppShell';
import theme from '@/styles/theme';

enableReactUse();

export const links: LinksFunction = () => [
	...(cssBundleHref ? [{ rel: 'stylesheet', href: cssBundleHref }] : []),
];

export default function App() {
	return (
		<html lang="en">
			<head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width,initial-scale=1" />
				<Meta />
				<Links />
				<ColorSchemeScript />
			</head>
			<body>
				<MantineProvider classNamesPrefix="me" theme={theme}>
					<AppShell>
						<Outlet />
					</AppShell>
					<ScrollRestoration />
					<Scripts />
					<LiveReload />
				</MantineProvider>
			</body>
		</html>
	);
}
