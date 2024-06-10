import { authLogin } from '@/api/auth';
import { userGet } from '@/api/user';
import { Login } from '@/components/login/Login';
import { hasSession } from '@/utils/cookies';
import {
	type ClientActionFunctionArgs,
	type ClientLoaderFunctionArgs,
	json,
	type MetaFunction,
	redirect,
} from '@remix-run/react';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Login | Medama' },
		{ name: 'description', content: 'Login into Medama Analytics.' },
	];
};

export const clientLoader = async ({ request }: ClientLoaderFunctionArgs) => {
	// If the user is already logged in, redirect them to the dashboard.
	if (hasSession(request)) {
		// Check if session hasn't been revoked
		await userGet({ cookie: request.headers.get('Cookie'), noRedirect: true });
		return redirect('/');
	}

	return { status: 200 };
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

	const { cookie } = await authLogin({
		body: {
			username,
			password,
		},
	});

	if (!cookie) {
		throw json('Failed to login.', {
			status: 401,
		});
	}

	return redirect('/', {
		headers: {
			'Set-Cookie': cookie,
		},
	});
};

export default function Index() {
	return (
		<main>
			<Login />
		</main>
	);
}
