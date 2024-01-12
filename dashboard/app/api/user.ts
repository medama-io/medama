import { client, type ClientOptions, type DataResponse } from './client';

const userGet = async (
	opts: ClientOptions
): Promise<DataResponse<'UserGet'>> => {
	const res = await client('/user', { method: 'GET', ...opts });
	return { data: await res.json(), res };
};

const userUpdate = async (
	opts: ClientOptions<'UserPatch'>
): Promise<DataResponse<'UserGet'>> => {
	const res = await client('/user', { method: 'PATCH', ...opts });
	return { data: await res.json(), res };
};

export { userGet, userUpdate };
