import '@mantine/core/styles.layer.css';
import 'mantine-datatable/styles.layer.css';
import '@fontsource-variable/inter/wght.css';
import '@/styles/global.css';

import { enableReactUse } from '@legendapp/state/config/enableReactUse';
import { ColorSchemeScript, MantineProvider } from '@mantine/core';
import {
	json,
	type ClientLoaderFunctionArgs,
	isRouteErrorResponse,
	Links,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
	useRouteError,
} from '@remix-run/react';

import { AppShell } from '@/components/layout/AppShell';
import theme from '@/styles/theme';
import { hasSession } from '@/utils/cookies';
import { API_BASE } from '@/api/client';

enableReactUse();

interface LoaderData {
	isLoggedIn: boolean;
}

interface DocumentProps {
	children: React.ReactNode;
}

export const clientLoader = () => {
	return json<LoaderData>({ isLoggedIn: Boolean(hasSession()) });
};

export const Document = ({ children }: DocumentProps) => {
	// While end users will have their API servers at a fixed domain, development mode will have the API server
	// running on localhost.
	let isLocalhost = false;
	let scriptSrc = '/script.js';
	if (API_BASE === 'localhost') {
		isLocalhost = true;
		scriptSrc = 'http://localhost:8080/script.js';
	}

	return (
		<html lang="en">
			<head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width,initial-scale=1" />
				<Meta />
				<Links />
				<ColorSchemeScript />
				{isLocalhost && <script defer src={scriptSrc} />}
			</head>
			<body>
				<MantineProvider classNamesPrefix="me" theme={theme}>
					<AppShell>{children}</AppShell>
					<ScrollRestoration />
					<Scripts />
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

export const HydrateFallback = () => {
	return (
		<Document>
			<p>Loading...</p>
		</Document>
	);
};

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
