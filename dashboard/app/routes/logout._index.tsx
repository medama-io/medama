import {
	isRouteErrorResponse,
	useLoaderData,
	useRevalidator,
	useRouteError,
	json,
	type MetaFunction,
} from '@remix-run/react';
import { useEffect } from 'react';

import { expireSession, hasSession } from '@/utils/cookies';
import { authLogout } from '@/api/auth';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Logout | Medama' },
		{ name: 'description', content: 'Logout from Medama Analytics.' },
	];
};

export const clientLoader = async () => {
	// If the user is already logged in, expire session cookie with success message.
	if (hasSession()) {
		const { res } = await authLogout();
		if (!res.ok) {
			throw new Error('Failed to logout.');
		}
		expireSession(true);
		return json('You have been sucessfully logged out.');
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
