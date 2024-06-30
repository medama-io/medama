import '@fontsource-variable/inter/wght.css';
// Comman
import '@mantine/core/styles/global.css';
import '@mantine/core/styles/ScrollArea.css';
import '@mantine/core/styles/UnstyledButton.css';
import '@mantine/core/styles/VisuallyHidden.css';
import '@mantine/core/styles/Paper.css';
import '@mantine/core/styles/Popover.css';
import '@mantine/core/styles/Group.css';
import '@mantine/core/styles/Overlay.css';
import '@mantine/core/styles/ModalBase.css';
import '@mantine/core/styles/Input.css';
import '@mantine/core/styles/Flex.css';
import '@mantine/core/styles/InlineInput.css';
// Layout
import '@mantine/core/styles/SimpleGrid.css';
import '@mantine/core/styles/Container.css';
import '@mantine/core/styles/Stack.css';
import '@mantine/core/styles/Center.css';
import '@mantine/core/styles/FloatingIndicator.css';
import '@mantine/core/styles/Drawer.css';
// Buttons
import '@mantine/core/styles/ActionIcon.css';
import '@mantine/core/styles/Burger.css';
import '@mantine/core/styles/CloseButton.css';
// Inputs
import '@mantine/core/styles/Checkbox.css';
import '@mantine/core/styles/Combobox.css';
import '@mantine/core/styles/PasswordInput.css';
// Navigation
import '@mantine/core/styles/Burger.css';
import '@mantine/core/styles/NavLink.css';
import '@mantine/core/styles/Tabs.css';
// Feedback
import '@mantine/core/styles/Skeleton.css';
import '@mantine/core/styles/Tooltip.css';
// Misc
import '@mantine/core/styles/Modal.css';
import '@mantine/core/styles/ColorSwatch.css';
import '@mantine/core/styles/Table.css';
import '@mantine/core/styles/Text.css';
import '@mantine/core/styles/Anchor.css';
import '@mantine/core/styles/Title.css';
import '@mantine/core/styles/Notification.css';
import '@mantine/notifications/styles.css';
import '@mantine/charts/styles.css';
import 'mantine-datatable/styles.css';
import 'react-indiana-drag-scroll/dist/style.css';
import '@/styles/global.css';

import {
	ColorSchemeScript,
	Flex,
	Loader,
	MantineProvider,
} from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import {
	Links,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
	isRouteErrorResponse,
	json,
	useLoaderData,
	useRouteError,
} from '@remix-run/react';

import { API_BASE } from '@/api/client';
import { AppShell } from '@/components/layout/AppShell';
import {
	BadRequestError,
	ForbiddenError,
	InternalServerError,
	NotFoundError,
} from '@/components/layout/Error';
import theme from '@/styles/theme';
import { EXPIRE_LOGGED_IN, hasSession } from '@/utils/cookies';

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
				<link
					rel="apple-touch-icon"
					sizes="180x180"
					href="/apple-touch-icon.png"
				/>
				<link rel="icon" type="image/svg+xml" href="/favicon.svg" />
				<link
					rel="icon"
					sizes="192x192"
					type="image/png"
					href="/favicon-192x192.png"
				/>
				<link
					rel="icon"
					sizes="48x48"
					type="image/png"
					href="/favicon-48x48.png"
				/>
				<link
					rel="icon"
					type="image/png"
					sizes="32x32"
					href="/favicon-32x32.png"
				/>
				<link
					rel="icon"
					type="image/png"
					sizes="16x16"
					href="/favicon-16x16.png"
				/>
				<link rel="manifest" href="/site.webmanifest" />
				<link rel="mask-icon" href="/safari-pinned-tab.svg" color="#17cd8c" />
				<meta name="apple-mobile-web-app-title" content="Medama Analytics" />
				<meta name="application-name" content="Medama Analytics" />
				<meta name="msapplication-TileColor" content="#111111" />
				<meta name="theme-color" content="#111111" />
				<Meta />
				<Links />
				<ColorSchemeScript />
				{isLocalhost && <script defer src={scriptSrc} />}
			</head>
			<body>
				<MantineProvider classNamesPrefix="me" theme={theme}>
					<Notifications />
					<AppShell>{children}</AppShell>
					<ScrollRestoration />
					<Scripts />
				</MantineProvider>
			</body>
		</html>
	);
};

export default function App() {
	// Trigger loader for session check.
	useLoaderData<LoaderData>();
	return (
		<Document>
			<Outlet />
		</Document>
	);
}

export const HydrateFallback = () => {
	return (
		<Document>
			<Flex justify="center" align="center" style={{ height: '90vh' }}>
				<Loader color="#17cd8c" type="bars" />
			</Flex>
		</Document>
	);
};

export const ErrorBoundary = () => {
	const error = useRouteError();
	console.error(error);

	if (isRouteErrorResponse(error)) {
		switch (error.status) {
			case 400: {
				return (
					<Document>
						<BadRequestError />
					</Document>
				);
			}
			case 403: {
				return (
					<Document>
						<ForbiddenError />
					</Document>
				);
			}
			case 404: {
				return (
					<Document>
						<NotFoundError />
					</Document>
				);
			}
		}

		return (
			<Document>
				<InternalServerError error={error.data ?? error.statusText} />
			</Document>
		);
	}

	if (error instanceof Error) {
		// If the error is due to a loader mismatch, reload the page as it may be
		// related to a bad cookie cache from the API restarting. This is probably
		// a bug in Remix SPA mode.
		if (error.message.startsWith('You defined a loader for route "routes')) {
			document.cookie = EXPIRE_LOGGED_IN;
			window.location.reload();
			return HydrateFallback();
		}

		return (
			<Document>
				<InternalServerError error={error.message} />
			</Document>
		);
	}

	return (
		<Document>
			<InternalServerError />
		</Document>
	);
};
