import { client, type ClientOptions, type DataResponse } from './client';

const authLogin = async (
	opts: ClientOptions<'AuthLogin'>,
): Promise<DataResponse> => {
	return { res: await client('/auth/login', { method: 'POST', ...opts }) };
};

export { authLogin };
