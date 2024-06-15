import { hasSession } from '@/utils/cookies';
import { redirect } from '@remix-run/react';
import { client, type ClientOptions, type DataResponse } from './client';

const userGet = async (
	opts?: ClientOptions,
): Promise<DataResponse<'UserGet'>> => {
	const res = await client('/user', { method: 'GET', ...opts });
	return { data: await res.json(), res };
};

const userUpdate = async (
	opts: ClientOptions<'UserPatch'>,
): Promise<DataResponse<'UserGet'>> => {
	const res = await client('/user', { method: 'PATCH', ...opts });
	return { data: await res.json(), res };
};

const userLoggedIn = async () => {
	// If the user is already logged in, redirect them to the dashboard.
	if (hasSession()) {
		// Check if session hasn't been revoked.
		const res = await client('/user', { method: 'GET' });
		if (res.ok) {
			return;
		}
	}

	throw redirect('/login');
};

export { userGet, userLoggedIn, userUpdate };
