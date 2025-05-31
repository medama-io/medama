import { notifications } from '@mantine/notifications';
import {
	type ClientActionFunctionArgs,
	type ClientLoaderFunctionArgs,
	type MetaFunction,
	data as json,
	redirect,
} from '@remix-run/react';

import { authLogin } from '@/api/auth';
import { userGet } from '@/api/user';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { Login } from '@/components/login/Login';
import { LOGGED_IN_COOKIE, hasSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Login | Medama' },
		{ name: 'description', content: 'Login into Medama Analytics.' },
	];
};

export const clientLoader = async (_: ClientLoaderFunctionArgs) => {
	// If the user is in demo mode (hostname matches demo.medama.io or medama.fly.dev), automatically
	// log them into the demo account.
	const hostname = window.location.hostname;
	const isDemo = hostname === 'demo.medama.io' || hostname === 'medama.fly.dev';
	if (isDemo) {
		const { res } = await authLogin({
			body: {
				username: 'admin',
				password: 'CHANGE_ME_ON_FIRST_LOGIN',
			},
			shouldThrow: false,
		});

		if (!res.ok) {
			throw new Error('Failed to login to demo account.');
		}

		// Set logged in cookie
		document.cookie = LOGGED_IN_COOKIE;

		return redirect('/');
	}

	// If the user is already logged in, redirect them to the dashboard.
	if (hasSession()) {
		// Check if session hasn't been revoked
		await userGet({ shouldRedirect: false });
		return redirect('/');
	}

	return null;
};

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();

	const username = body.get('username')
		? String(body.get('username'))
		: undefined;
	const password = body.get('password')
		? String(body.get('password'))
		: undefined;

	if (!username || !password) {
		throw json('Missing email or password', {
			status: 400,
		});
	}

	const { res } = await authLogin({
		body: {
			username,
			password,
		},
		shouldThrow: false,
	});

	if (!res.ok) {
		if (res.status === 401) {
			notifications.show({
				title: 'Invalid username or password.',
				message: 'Please try again.',
				withBorder: true,
				color: 'red',
			});

			return {
				message: 'Invalid username or password. Please try again.',
			};
		}

		throw new Response('Failed to login.', {
			status: res.status,
		});
	}

	// Set logged in cookie
	document.cookie = LOGGED_IN_COOKIE;

	// If user has redirect query param, redirect them to that URL
	const url = new URL(request.url);
	const redirectPathname = url.searchParams.get('redirect');
	if (redirectPathname) {
		return redirect(redirectPathname);
	}

	return redirect('/');
};

export default function Index() {
	return (
		<>
			<InnerHeader>
				<h1>Log In</h1>
			</InnerHeader>
			<main>
				<Login />
			</main>
		</>
	);
}
