import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

import { isLoggedIn$ } from '@/observables';
import { EXPIRE_COOKIE, hasSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Logout | Medama' },
		{ name: 'description', content: 'Logout from Medama Analytics.' },
	];
};

interface LoaderData {
	loggedOut: boolean;
}

export const loader = ({ request }: LoaderFunctionArgs) => {
	// If the user is already logged in, expire session cookie with success message.
	if (hasSession(request)) {
		return json<LoaderData>(
			{ loggedOut: true },
			{
				status: 200,
				headers: {
					'Set-Cookie': EXPIRE_COOKIE,
				},
			}
		);
	}

	return json<LoaderData>({ loggedOut: false }, { status: 200 });
};

export default function Index() {
	const { loggedOut } = useLoaderData<LoaderData>();
	isLoggedIn$.set(false);

	return (
		<div>
			<h1>Logout</h1>
			{loggedOut ? (
				<p>You have been logged out.</p>
			) : (
				<p>You are not logged in.</p>
			)}
		</div>
	);
}
