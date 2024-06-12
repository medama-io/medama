import { client, type ClientOptions, type DataResponse } from './client';

const authLogin = async (
	opts: ClientOptions<'AuthLogin'>,
): Promise<DataResponse> => {
	return { res: await client('/auth/login', { method: 'POST', ...opts }) };
};

const authLogout = async (): Promise<DataResponse> => {
	return { res: await client('/auth/logout', { method: 'POST' }) };
};

export { authLogin, authLogout };
