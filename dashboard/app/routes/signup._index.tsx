import {
	type ActionFunctionArgs,
	json,
	type MetaFunction,
} from '@remix-run/node';

import { Signup } from '@/components/login/Login';
import { type PostUser } from '@/utils/types';
import { postUser } from '@/utils/user';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Signup | Medama' },
		{ name: 'description', content: 'Sign up to Medama Analytics.' },
	];
};

export const action = async ({ request }: ActionFunctionArgs) => {
	const body = await request.formData();

	const email = body.get('email')?.toString();
	const password = body.get('password')?.toString();

	if (!email || !password) {
		return new Response('Missing email or password', {
			status: 400,
		});
	}

	const post: PostUser = {
		email,
		password,
	};

	const res = await postUser(post);
	console.log(res);

	if (res.data.email) {
		return json(res.data, {
			status: 200,
			headers: {
				'Set-Cookie': res.cookie,
			},
		});
	}
};

export default function Index() {
	return <Signup />;
}
