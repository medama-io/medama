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
// Buttons
import '@mantine/core/styles/ActionIcon.css';
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
// Misc
import '@mantine/core/styles/Modal.css';
import '@mantine/core/styles/Table.css';
import '@mantine/core/styles/Text.css';
import '@mantine/core/styles/Anchor.css';
import '@mantine/core/styles/Title.css';
import '@/styles/global.css';
import 'mantine-datatable/styles.css';

import { enableReactUse } from '@legendapp/state/config/enableReactUse';
import { ColorSchemeScript, MantineProvider } from '@mantine/core';
import {
	Links,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
	isRouteErrorResponse,
	json,
	useRouteError,
} from '@remix-run/react';

import { API_BASE } from '@/api/client';
import { AppShell } from '@/components/layout/AppShell';
import theme from '@/styles/theme';
import { hasSession } from '@/utils/cookies';
import { InternalServerError, NotFoundError } from '@/components/layout/Error';

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
				<InternalServerError error={error.statusText} />
			</Document>
		);
	}

	if (error instanceof Error) {
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
