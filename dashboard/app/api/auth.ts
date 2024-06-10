import { json } from '@remix-run/react';

import { client, type ClientOptions, type DataResponse } from './client';

const authLogin = async (
	opts: ClientOptions<'AuthLogin'>
): Promise<DataResponse> => {
	const res = await client('/auth/login', { method: 'POST', ...opts });
	const cookie = res.headers.get('Set-Cookie');
	if (!cookie) throw json('Failed to get session.', { status: 401 });
	return { cookie, res };
};

export { authLogin };
