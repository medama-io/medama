import { Button } from '@mantine/core';
import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
	redirect,
} from '@remix-run/node';
import {
	isRouteErrorResponse,
	NavLink,
	useLoaderData,
	useRouteError,
} from '@remix-run/react';

import { type components } from '@/api/types';
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

export const loader = async ({ request }: LoaderFunctionArgs) => {
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
		<div>
			<h1>Websites</h1>
			{JSON.stringify(websites)}
			<Button component={NavLink} to="/localhost">
				localhost
			</Button>
		</div>
	);
}

export const ErrorBoundary = () => {
	const error = useRouteError();

	if (isRouteErrorResponse(error) && error.status === 404) {
		return (
			<div>
				<h1>404</h1>
				<p>No websites found</p>
				<Button component={NavLink} to="/add">
					Add Website
				</Button>
			</div>
		);
	}
};
