import { useEffect } from 'react';
import { redirect, useRevalidator } from 'react-router';

import { authLogout } from '@/api/auth';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { expireSession, hasSession } from '@/utils/cookies';
import type { Route } from './+types/logout._index';

export const meta: Route.MetaFunction = () => {
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
		return 'You have been successfully logged out.';
	}

	return redirect('/login');
};

export default function Index() {
	// We want to call the revalidator to trigger the root loader and update the header accordingly.
	const revalidator = useRevalidator();
	useEffect(() => {
		revalidator.revalidate();
	}, [revalidator]);

	return (
		<>
			<InnerHeader>
				<h1>Log out</h1>
			</InnerHeader>
			<main>
				<p>You have been logged out.</p>
			</main>
		</>
	);
}
