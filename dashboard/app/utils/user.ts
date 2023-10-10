import { type GetUser, type PostUser } from './types';

const postUser = async (user: PostUser) => {
	const host = process.env.CORE_API_HOST ?? 'http://localhost:8080';

	const res = await fetch(`${host}/user`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(user),
	});

	if (!res.ok) {
		throw new Error(`HTTP error! status: ${res.status}`);
	}

	// Get cookie from response and set it in the browser
	const cookie = res.headers.get('Set-Cookie');
	if (cookie === null) {
		throw new Error('No cookie found in response');
	}

	const data = (await res.json()) as GetUser;

	return { cookie, data };
};

export { postUser };
