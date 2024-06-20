import { authLogin } from '@/api/auth';
import { userGet } from '@/api/user';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { Login } from '@/components/login/Login';
import { LOGGED_IN_COOKIE, hasSession } from '@/utils/cookies';
import { notifications } from '@mantine/notifications';
import {
	json,
	redirect,
	useActionData,
	type ClientActionFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Login | Medama' },
		{ name: 'description', content: 'Login into Medama Analytics.' },
	];
};

export const clientLoader = async () => {
	// If the user is already logged in, redirect them to the dashboard.
	if (hasSession()) {
		// Check if session hasn't been revoked
		await userGet({ noRedirect: true });
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
		noThrow: true,
	});

	if (!res.ok) {
		if (res.status === 401) {
			notifications.show({
				title: 'Invalid username or password.',
				message: 'Please try again.',
				withBorder: true,
				color: 'red',
			});
			return json({
				message: 'Invalid username or password. Please try again.',
			});
		}

		throw new Response('Failed to login.', {
			status: res.status,
		});
	}

	// Set logged in cookie
	document.cookie = LOGGED_IN_COOKIE;

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
