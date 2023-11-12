import {
	type ActionFunctionArgs,
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
	redirect,
} from '@remix-run/node';
import { useRevalidator } from '@remix-run/react';
import { useEffect } from 'react';

import { authLogin } from '@/api/auth';
import { userGet } from '@/api/user';
import { Login } from '@/components/login/Login';
import { hasSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Login | Medama' },
		{ name: 'description', content: 'Login into Medama Analytics.' },
	];
};

export const loader = async ({ request }: LoaderFunctionArgs) => {
	// If the user is already logged in, redirect them to the dashboard.
	if (hasSession(request)) {
		// Check if session hasn't been revoked
		await userGet({ cookie: request.headers.get('Cookie'), noRedirect: true });
		return redirect('/');
	}

	return { status: 200 };
};

export const action = async ({ request }: ActionFunctionArgs) => {
	const body = await request.formData();

	const email = body.get('email')?.toString();
	const password = body.get('password')?.toString();

	if (!email || !password) {
		throw json('Missing email or password', {
			status: 400,
		});
	}

	const { cookie } = await authLogin({
		body: {
			email,
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
	// We want to call the revalidator to trigger the root loader and update the header accordingly.
	/* const revalidator = useRevalidator();
	useEffect(() => {
		revalidator.revalidate();
	}, [revalidator]); */

	return <Login />;
}
