import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import {
	isRouteErrorResponse,
	useLoaderData,
	useRevalidator,
	useRouteError,
} from '@remix-run/react';
import { useEffect } from 'react';

import { EXPIRE_COOKIE, hasSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Logout | Medama' },
		{ name: 'description', content: 'Logout from Medama Analytics.' },
	];
};

export const loader = ({ request }: LoaderFunctionArgs) => {
	// If the user is already logged in, expire session cookie with success message.
	if (hasSession(request)) {
		return json('You have been sucessfully logged out.', {
			status: 200,
			headers: {
				'Set-Cookie': EXPIRE_COOKIE,
			},
		});
	}

	throw json('You are not logged in.', {
		status: 401,
	});
};

export default function Index() {
	// Trigger loader for cookie expiration.
	useLoaderData();

	// We want to call the revalidator to trigger the root loader and update the header accordingly.
	const revalidator = useRevalidator();
	useEffect(() => {
		revalidator.revalidate();
	}, [revalidator]);

	return (
		<div>
			<h1>Logout</h1>
			<p>You have been logged out.</p>
		</div>
	);
}

export const ErrorBoundary = () => {
	const error = useRouteError();

	if (isRouteErrorResponse(error) && error.status === 401) {
		return (
			<div>
				<h1>Logout</h1>
				<p>You are not logged in.</p>
			</div>
		);
	}
};
