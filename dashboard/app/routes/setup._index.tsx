import {
	type ActionFunctionArgs,
	json,
	type MetaFunction,
	redirect,
} from '@remix-run/node';

import { userCreate } from '@/api/user';
import { Setup } from '@/components/setup/Setup';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Setup | Medama' },
		{ name: 'description', content: 'Sign up to Medama Analytics.' },
	];
};

export const action = async ({ request }: ActionFunctionArgs) => {
	const body = await request.formData();

	const email = body.get('email')?.toString();
	const password = body.get('password')?.toString();

	if (!email || !password) {
		throw json('Missing email or password.', {
			status: 400,
		});
	}

	const { res, cookie } = await userCreate({
		body: {
			email,
			password,
		},
	});

	if (!res.ok || !cookie) {
		throw json('Failed to get session.', {
			status: res.status,
		});
	}

	return redirect('/', {
		headers: {
			'Set-Cookie': cookie,
		},
	});
};

export default function Index() {
	return <Setup />;
}
