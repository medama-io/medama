import '@mantine/core/styles.css';
import '@/styles/global.css';
import '@fontsource-variable/inter/wght.css';

import { enableReactUse } from '@legendapp/state/config/enableReactUse';
import { ColorSchemeScript, MantineProvider } from '@mantine/core';
import { cssBundleHref } from '@remix-run/css-bundle';
import {
	json,
	type LinksFunction,
	type LoaderFunctionArgs,
} from '@remix-run/node';
import {
	Links,
	LiveReload,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
	useLoaderData,
} from '@remix-run/react';

import { AppShell } from '@/components/layout/AppShell';
import theme from '@/styles/theme';

import { getSession } from './utils/cookies';

enableReactUse();

export const links: LinksFunction = () => [
	...(cssBundleHref ? [{ rel: 'stylesheet', href: cssBundleHref }] : []),
];

interface LoaderData {
	isLoggedIn: boolean;
}

export const loader = ({ request }: LoaderFunctionArgs) => {
	const session = getSession(request);
	return json<LoaderData>({ isLoggedIn: Boolean(session) });
};

export default function App() {
	const { isLoggedIn } = useLoaderData<LoaderData>();

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
					<AppShell isLoggedIn={isLoggedIn}>
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
