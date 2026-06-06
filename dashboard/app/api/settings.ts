import { type ClientOptions, client, type DataResponse } from './client';

export const systemSettingsGet = async (
	opts?: ClientOptions,
): Promise<DataResponse<'SystemSettings'>> => {
	const res = await client('/system/settings', { method: 'GET', ...opts });
	return { data: await res.json(), res };
};

export const systemSettingsUpdate = async (
	opts: ClientOptions<'SystemSettings'>,
): Promise<DataResponse<'SystemSettings'>> => {
	const res = await client('/system/settings', { method: 'PATCH', ...opts });
	return { data: await res.json(), res };
};
