import { Divider, Paper, Text } from '@mantine/core';
import {
	json,
	type MetaFunction,
	redirect,
	type ClientLoaderFunctionArgs,
	isRouteErrorResponse,
	Link,
	useLoaderData,
	useRouteError,
} from '@remix-run/react';

import type { components } from '@/api/types';
import { websiteList } from '@/api/websites';
import { hasSession } from '@/utils/cookies';

interface LoaderData {
	websites: Array<components['schemas']['WebsiteGet']>;
}

export const meta: MetaFunction = () => {
	return [
		{ title: 'Medama | Privacy Focused Web Analytics' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const clientLoader = async ({ request }: ClientLoaderFunctionArgs) => {
	// Check for session cookie and redirect to login if missing
	if (!hasSession(request)) {
		throw redirect('/login');
	}

	const { data } = await websiteList({ cookie: request.headers.get('Cookie') });

	if (!data) {
		throw json('Failed to get websites.', {
			status: 500,
		});
	}

	return json<LoaderData>({ websites: data });
};

export default function Index() {
	const { websites } = useLoaderData<LoaderData>();

	return (
		<main>
			<h1>Websites</h1>
			<Divider mb={30} />
			{websites.map((website) => (
				<Paper
					key={website.hostname}
					withBorder
					w={300}
					p={8}
					radius={8}
					component={Link}
					to={`/${website.hostname}`}
					prefetch="intent"
				>
					<Text>{website.name}</Text>
					<Text size="xs" c="gray">
						{website.hostname}
					</Text>
				</Paper>
			))}
		</main>
	);
}

export const ErrorBoundary = () => {
	const error = useRouteError();

	if (isRouteErrorResponse(error) && error.status === 404) {
		return (
			<main>
				<h1>404</h1>
				<p>No websites found</p>
				<Paper
					withBorder
					w={300}
					p={8}
					radius={8}
					component={Link}
					to="/add"
					prefetch="intent"
				>
					Add Website
				</Paper>
			</main>
		);
	}
};
