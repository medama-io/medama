import { type ClientOptions, type DataResponse, client } from './client';

const usageGet = async (): Promise<DataResponse<'SettingsUsageGet'>> => {
	const res = await client('/settings/usage', {});
	return {
		data: await res.json(),
		res,
	};
};

const usagePatch = async (
	opts: ClientOptions<'SettingsUsagePatch'>,
): Promise<DataResponse> => {
	const res = await client('/settings/usage', { method: 'PATCH', ...opts });
	return { res };
};

export { usageGet, usagePatch };
