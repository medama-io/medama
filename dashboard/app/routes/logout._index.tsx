import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

import { getSession } from '@/utils/cookies';
import { SESSION_NAME } from '@/utils/types';

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
	if (getSession(request)) {
		return json<LoaderData>(
			{ loggedOut: true },
			{
				status: 200,
				headers: {
					'Set-Cookie': `${SESSION_NAME}=; HttpOnly; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT`,
				},
			}
		);
	}

	return json<LoaderData>({ loggedOut: false }, { status: 200 });
};

export default function Index() {
	const { loggedOut } = useLoaderData<LoaderData>();
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
