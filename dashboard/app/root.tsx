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
	isRouteErrorResponse,
	Links,
	LiveReload,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
	useRouteError,
} from '@remix-run/react';

import { AppShell } from '@/components/layout/AppShell';
import theme from '@/styles/theme';
import { hasSession } from '@/utils/cookies';

enableReactUse();

interface LoaderData {
	isLoggedIn: boolean;
}

interface DocumentProps {
	children: React.ReactNode;
}

export const links: LinksFunction = () => [
	...(cssBundleHref ? [{ rel: 'stylesheet', href: cssBundleHref }] : []),
];

export const loader = ({ request }: LoaderFunctionArgs) => {
	const session = hasSession(request);

	return json<LoaderData>({ isLoggedIn: Boolean(session) });
};

export const Document = ({ children }: DocumentProps) => {
	return (
		<html lang="en">
			<head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width,initial-scale=1" />
				<Meta />
				<Links />
				<ColorSchemeScript />
				<script defer data-api="medama-core.fly.dev" src="/medama.js" />
			</head>
			<body>
				<MantineProvider classNamesPrefix="me" theme={theme}>
					<AppShell>{children}</AppShell>
					<ScrollRestoration />
					<Scripts />
					<LiveReload />
				</MantineProvider>
			</body>
		</html>
	);
};

export default function App() {
	return (
		<Document>
			<Outlet />
		</Document>
	);
}

export const ErrorBoundary = () => {
	const error = useRouteError();

	if (isRouteErrorResponse(error)) {
		switch (error.status) {
			case 401: {
				return (
					<Document>
						<p>You don&apos;t have access to this page.</p>
					</Document>
				);
			}
			case 404: {
				return <Document>Page not found!</Document>;
			}
		}

		return (
			<Document>
				Something went wrong: {error.status} {error.statusText}
			</Document>
		);
	}

	if (error instanceof Error) {
		return <Document>Something went wrong: {error.message}</Document>;
	}

	return <Document>Something went wrong: Unknown Error</Document>;
};
