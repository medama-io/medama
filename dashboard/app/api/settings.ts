import { type ClientOptions, client, type DataResponse } from './client';

export const tenantSettingsGet = async (
	opts?: ClientOptions,
): Promise<DataResponse<'TenantSettings'>> => {
	const res = await client('/tenant/settings', { method: 'GET', ...opts });
	return { data: await res.json(), res };
};

export const tenantSettingsUpdate = async (
	opts: ClientOptions<'TenantSettings'>,
): Promise<DataResponse<'TenantSettings'>> => {
	const res = await client('/tenant/settings', { method: 'PATCH', ...opts });
	return { data: await res.json(), res };
};
